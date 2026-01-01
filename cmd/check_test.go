package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jimschubert/ossify/internal/model"
)

// setupTestDirectory creates a temporary directory with the given structure
// structure is a map of path -> isDirectory (true = dir, false = file)
func setupTestDirectory(t *testing.T, structure map[string]bool) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "ossify-cmd-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	for path, isDir := range structure {
		fullPath := filepath.Join(tempDir, path)
		if isDir {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Fatalf("failed to create directory %s: %v", path, err)
			}
		} else {
			// Ensure parent exists
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				t.Fatalf("failed to create parent for %s: %v", path, err)
			}
			if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
				t.Fatalf("failed to create file %s: %v", path, err)
			}
		}
	}

	return tempDir
}

func TestCheckCmd_Help(t *testing.T) {
	// Test that the help flag works
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"check", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Check a directory against") {
		t.Errorf("expected help output to contain description, got: %s", output)
	}
	if !strings.Contains(output, "--convention") {
		t.Errorf("expected help output to contain --convention flag, got: %s", output)
	}
	if !strings.Contains(output, "--file") {
		t.Errorf("expected help output to contain --file flag, got: %s", output)
	}
}

func TestLoadConventionFromFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantName    string
		wantRules   int
		wantErr     bool
		errContains string
	}{
		{
			name: "valid convention",
			content: `{
				"name": "Test Convention",
				"rules": [
					{"level": "required", "type": "file", "value": "README.md"}
				]
			}`,
			wantName:  "Test Convention",
			wantRules: 1,
			wantErr:   false,
		},
		{
			name: "convention with multiple rules",
			content: `{
				"name": "Multi Rule",
				"rules": [
					{"level": "required", "type": "file", "value": "README.md"},
					{"level": "optional", "type": "directory", "value": "src"},
					{"level": "prohibited", "type": "directory", "value": "vendor"}
				]
			}`,
			wantName:  "Multi Rule",
			wantRules: 3,
			wantErr:   false,
		},
		{
			name:        "invalid JSON",
			content:     `{not valid json`,
			wantErr:     true,
			errContains: "invalid JSON",
		},
		{
			name: "empty rules",
			content: `{
				"name": "Empty Rules",
				"rules": []
			}`,
			wantErr:     true,
			errContains: "at least one rule",
		},
		{
			name: "missing name uses filename",
			content: `{
				"rules": [
					{"level": "required", "type": "file", "value": "README.md"}
				]
			}`,
			wantRules: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file with content
			tempFile, err := os.CreateTemp("", "convention-*.json")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer func() { _ = os.Remove(tempFile.Name()) }()

			if _, err := tempFile.WriteString(tt.content); err != nil {
				t.Fatalf("failed to write content: %v", err)
			}
			if err := tempFile.Close(); err != nil {
				t.Fatalf("failed to close temp file: %v", err)
			}

			// Test loading
			convention, err := loadConventionFromFile(tempFile.Name())

			if (err != nil) != tt.wantErr {
				t.Errorf("loadConventionFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want containing %q", err, tt.errContains)
				}
				return
			}

			if tt.wantName != "" && convention.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", convention.Name, tt.wantName)
			}

			if len(convention.Rules) != tt.wantRules {
				t.Errorf("Rules count = %d, want %d", len(convention.Rules), tt.wantRules)
			}
		})
	}
}

func TestLoadConventionFromFile_NotFound(t *testing.T) {
	_, err := loadConventionFromFile("/nonexistent/path/to/file.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("expected 'failed to read file' error, got: %v", err)
	}
}

func TestFindConventionByName(t *testing.T) {
	conventions := []model.Convention{
		{Name: "Go", Rules: []model.Rule{{Level: model.Required, Type: model.File, Value: "go.mod"}}},
		{Name: "Standard Distribution", Rules: []model.Rule{{Level: model.Required, Type: model.Directory, Value: "src"}}},
		{Name: "Node.js", Rules: []model.Rule{{Level: model.Required, Type: model.File, Value: "package.json"}}},
	}

	tests := []struct {
		name      string
		search    string
		wantFound bool
		wantName  string
	}{
		{"exact match", "Go", true, "Go"},
		{"case insensitive", "go", true, "Go"},
		{"case insensitive upper", "GO", true, "Go"},
		{"multi-word", "Standard Distribution", true, "Standard Distribution"},
		{"multi-word case insensitive", "standard distribution", true, "Standard Distribution"},
		{"not found", "Python", false, ""},
		{"partial match fails", "Stand", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findConventionByName(conventions, tt.search)

			if tt.wantFound {
				if result == nil {
					t.Errorf("expected to find convention %q, got nil", tt.search)
					return
				}
				if result.Name != tt.wantName {
					t.Errorf("Name = %q, want %q", result.Name, tt.wantName)
				}
			} else {
				if result != nil {
					t.Errorf("expected nil, got convention %q", result.Name)
				}
			}
		})
	}
}

func TestCheckCmd_WithConventionFile(t *testing.T) {
	// Create a test directory structure
	testDir := setupTestDirectory(t, map[string]bool{
		"README.md": false,
		"src":       true,
	})
	defer func() { _ = os.RemoveAll(testDir) }()

	// Create a convention file
	convention := model.Convention{
		Name: "Test",
		Rules: []model.Rule{
			{Level: model.Required, Type: model.File, Value: "README.md"},
			{Level: model.Required, Type: model.Directory, Value: "src"},
		},
	}
	convData, _ := json.Marshal(convention)

	convFile, err := os.CreateTemp("", "test-convention-*.json")
	if err != nil {
		t.Fatalf("failed to create convention file: %v", err)
	}
	defer func() { _ = os.Remove(convFile.Name()) }()
	if _, err := convFile.Write(convData); err != nil {
		t.Fatalf("failed to write convention data: %v", err)
	}
	if err := convFile.Close(); err != nil {
		t.Fatalf("failed to close convention file: %v", err)
	}

	// Reset flags before test
	checkFlags = &CheckFlags{}

	// Run the check command
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"check", "-f", convFile.Name(), "-d", testDir})

	// The command will call os.Exit on failure, so we need to handle that
	// For now, just verify the setup works
	// In a real test, you'd want to refactor to avoid os.Exit or use a test helper

	// Test the underlying function directly instead
	conv, err := loadConventionFromFile(convFile.Name())
	if err != nil {
		t.Fatalf("failed to load convention: %v", err)
	}

	result, err := conv.Evaluate(testDir)
	if err != nil {
		t.Fatalf("failed to evaluate: %v", err)
	}

	if result.FailCount != 0 {
		t.Errorf("expected 0 failures, got %d", result.FailCount)
	}
	if result.PassCount != 2 {
		t.Errorf("expected 2 passes, got %d", result.PassCount)
	}
}

func TestCheckCmd_Integration(t *testing.T) {
	// Skip in short mode as this is an integration test
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create test directory
	testDir := setupTestDirectory(t, map[string]bool{
		"README.md": false,
		"LICENSE":   false,
		"docs":      true,
	})
	defer func() { _ = os.RemoveAll(testDir) }()

	// Create convention file
	convFile, err := os.CreateTemp("", "conv-*.json")
	if err != nil {
		t.Fatalf("failed to create convention file: %v", err)
	}
	defer func() { _ = os.Remove(convFile.Name()) }()
	if _, err := convFile.WriteString(`{
		"name": "Integration Test",
		"rules": [
			{"level": "required", "type": "file", "value": "README.md"},
			{"level": "required", "type": "file", "value": "LICENSE"},
			{"level": "required", "type": "directory", "value": "docs"},
			{"level": "prohibited", "type": "directory", "value": "vendor"}
		]
	}`); err != nil {
		t.Fatalf("failed to write convention: %v", err)
	}
	if err := convFile.Close(); err != nil {
		t.Fatalf("failed to close convention file: %v", err)
	}

	// Load and evaluate
	conv, err := loadConventionFromFile(convFile.Name())
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	result, err := conv.Evaluate(testDir)
	if err != nil {
		t.Fatalf("failed to evaluate: %v", err)
	}

	// All should pass
	if result.HasFailures() {
		t.Errorf("expected no failures, got %d", result.FailCount)
	}
	if result.PassCount != 4 {
		t.Errorf("expected 4 passes, got %d", result.PassCount)
	}
}
