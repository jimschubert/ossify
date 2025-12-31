package model

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T, structure map[string]bool) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "ossify-check-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	for name, isDir := range structure {
		fullPath := filepath.Join(tempDir, name)
		if isDir {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Fatalf("failed to create directory %s: %v", name, err)
			}
		} else {
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				t.Fatalf("failed to create parent dir for %s: %v", name, err)
			}
			if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
				t.Fatalf("failed to create file %s: %v", name, err)
			}
		}
	}

	return tempDir
}

func TestConvention_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		convention    Convention
		structure     map[string]bool // true = directory, false = file
		wantPassCount int
		wantFailCount int
		wantWarnCount int
		wantSkipCount int
	}{
		{
			name: "all required items present",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Required, Type: Directory, Value: "src"},
					{Level: Required, Type: File, Value: "README.md"},
				},
			},
			structure: map[string]bool{
				"src":       true,
				"README.md": false,
			},
			wantPassCount: 2,
			wantFailCount: 0,
		},
		{
			name: "missing required directory",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Required, Type: Directory, Value: "src"},
					{Level: Required, Type: File, Value: "README.md"},
				},
			},
			structure: map[string]bool{
				"README.md": false,
			},
			wantPassCount: 1,
			wantFailCount: 1,
		},
		{
			name: "missing required file",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Required, Type: File, Value: "LICENSE"},
				},
			},
			structure:     map[string]bool{},
			wantPassCount: 0,
			wantFailCount: 1,
		},
		{
			name: "prohibited directory exists",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Prohibited, Type: Directory, Value: "src"},
				},
			},
			structure: map[string]bool{
				"src": true,
			},
			wantPassCount: 0,
			wantFailCount: 1,
		},
		{
			name: "prohibited directory does not exist",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Prohibited, Type: Directory, Value: "src"},
				},
			},
			structure:     map[string]bool{},
			wantPassCount: 1,
			wantFailCount: 0,
		},
		{
			name: "optional items missing",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Optional, Type: Directory, Value: "docs"},
					{Level: Optional, Type: File, Value: "CONTRIBUTING.md"},
				},
			},
			structure:     map[string]bool{},
			wantPassCount: 2, // Optional missing counts as pass
			wantFailCount: 0,
			wantSkipCount: 0,
		},
		{
			name: "preferred items missing",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Preferred, Type: Directory, Value: "docs"},
				},
			},
			structure:     map[string]bool{},
			wantPassCount: 0,
			wantFailCount: 0,
			wantWarnCount: 1,
		},
		{
			name: "pattern matching - files exist",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Required, Type: Pattern, Value: "*.go"},
				},
			},
			structure: map[string]bool{
				"main.go": false,
				"util.go": false,
			},
			wantPassCount: 1,
			wantFailCount: 0,
		},
		{
			name: "pattern matching - no matches",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Required, Type: Pattern, Value: "*.py"},
				},
			},
			structure: map[string]bool{
				"main.go": false,
			},
			wantPassCount: 0,
			wantFailCount: 1,
		},
		{
			name: "pattern matching - prohibited pattern exists",
			convention: Convention{
				Name: "Test",
				Rules: []Rule{
					{Level: Prohibited, Type: Pattern, Value: "*.bak"},
				},
			},
			structure: map[string]bool{
				"file.bak": false,
			},
			wantPassCount: 0,
			wantFailCount: 1,
		},
		{
			name: "mixed rules",
			convention: Convention{
				Name: "Go Project",
				Rules: []Rule{
					{Level: Required, Type: File, Value: "go.mod"},
					{Level: Required, Type: File, Value: "README.md"},
					{Level: Optional, Type: Directory, Value: "cmd"},
					{Level: Prohibited, Type: Directory, Value: "src"},
					{Level: Preferred, Type: File, Value: "LICENSE"},
				},
			},
			structure: map[string]bool{
				"go.mod":    false,
				"README.md": false,
				"cmd":       true,
			},
			wantPassCount: 4, // go.mod, README.md, cmd, src (prohibited & absent)
			wantFailCount: 0,
			wantWarnCount: 1, // LICENSE preferred but missing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.structure)
			defer os.RemoveAll(tempDir)

			result, err := tt.convention.Evaluate(tempDir)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}

			if result.PassCount != tt.wantPassCount {
				t.Errorf("PassCount = %d, want %d", result.PassCount, tt.wantPassCount)
			}
			if result.FailCount != tt.wantFailCount {
				t.Errorf("FailCount = %d, want %d", result.FailCount, tt.wantFailCount)
			}
			if result.WarnCount != tt.wantWarnCount {
				t.Errorf("WarnCount = %d, want %d", result.WarnCount, tt.wantWarnCount)
			}
			if result.SkipCount != tt.wantSkipCount {
				t.Errorf("SkipCount = %d, want %d", result.SkipCount, tt.wantSkipCount)
			}

			// Verify HasFailures
			expectedHasFailures := tt.wantFailCount > 0
			if result.HasFailures() != expectedHasFailures {
				t.Errorf("HasFailures() = %v, want %v", result.HasFailures(), expectedHasFailures)
			}
		})
	}
}

func TestEvaluateRule_DirectoryAsFile(t *testing.T) {
	// Test that a directory is not treated as a file
	tempDir := setupTestDir(t, map[string]bool{
		"mydir": true, // This is a directory
	})
	defer os.RemoveAll(tempDir)

	convention := Convention{
		Name: "Test",
		Rules: []Rule{
			{Level: Required, Type: File, Value: "mydir"}, // Expecting a file, but it's a directory
		},
	}

	result, err := convention.Evaluate(tempDir)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if result.PassCount != 0 || result.FailCount != 1 {
		t.Errorf("Expected file rule to fail when path is directory, got pass=%d fail=%d",
			result.PassCount, result.FailCount)
	}
}

func TestEvaluateRule_FileAsDirectory(t *testing.T) {
	// Test that a file is not treated as a directory
	tempDir := setupTestDir(t, map[string]bool{
		"myfile": false, // This is a file
	})
	defer os.RemoveAll(tempDir)

	convention := Convention{
		Name: "Test",
		Rules: []Rule{
			{Level: Required, Type: Directory, Value: "myfile"}, // Expecting a directory, but it's a file
		},
	}

	result, err := convention.Evaluate(tempDir)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if result.PassCount != 0 || result.FailCount != 1 {
		t.Errorf("Expected directory rule to fail when path is file, got pass=%d fail=%d",
			result.PassCount, result.FailCount)
	}
}

func TestCheckResult_HasFailures(t *testing.T) {
	tests := []struct {
		name      string
		failCount int
		want      bool
	}{
		{"no failures", 0, false},
		{"one failure", 1, true},
		{"multiple failures", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &CheckResult{FailCount: tt.failCount}
			if got := result.HasFailures(); got != tt.want {
				t.Errorf("HasFailures() = %v, want %v", got, tt.want)
			}
		})
	}
}
