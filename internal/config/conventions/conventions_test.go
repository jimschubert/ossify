package conventions

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jimschubert/ossify/internal/config"
	"github.com/jimschubert/ossify/internal/model"
	"github.com/jimschubert/ossify/internal/util"
)

func Test_indexOf(t *testing.T) {
	type args struct {
		data   []string
		search string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "first index", args: args{data: []string{"asdf", "jkl;"}, search: "asdf"}, want: 0},
		{name: "last index", args: args{data: []string{"asdf", "jkl;"}, search: "jkl;"}, want: 1},
		{name: "other index", args: args{data: []string{"asdf", "aaaa", "bbbb", "cccc", "jkl;"}, search: "bbbb"}, want: 2},
		{name: "not found", args: args{data: []string{"asdf", "aaaa", "bbbb", "cccc"}, search: "jkl;"}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.StringSearch(tt.args.data, tt.args.search); got != tt.want {
				t.Errorf("indexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_MarshalJSON(t *testing.T) {
	type fields struct {
		Level model.StrictnessLevel
		Type  model.RuleType
		Value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"valid input 0/0/.bak", fields{model.StrictnessLevel(0), model.RuleType(0), ".bak"}, `{"level":"prohibited","type":"unspecified","value":".bak"}`, false},
		{"valid input 1/1/src", fields{model.StrictnessLevel(1), model.RuleType(1), "src"}, `{"level":"optional","type":"directory","value":"src"}`, false},
		{"valid input 2/2/LICENSE", fields{model.StrictnessLevel(2), model.RuleType(2), "LICENSE"}, `{"level":"preferred","type":"file","value":"LICENSE"}`, false},
		{"valid input mixed/other", fields{model.StrictnessLevel(3), model.RuleType(1), "tools"}, `{"level":"required","type":"directory","value":"tools"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &model.Rule{
				Level: tt.fields.Level,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			got, err := r.MarshalJSON()
			actual := string(got[:])
			if (err != nil) != tt.wantErr {
				t.Errorf("Rule.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr {
				return
			}
			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("Rule.MarshalJSON() = %v, want %v", actual, tt.want)
			}
		})
	}
}

func TestRule_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Level model.StrictnessLevel
		Type  model.RuleType
		Value string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"valid input 0/0/.bak", fields{model.StrictnessLevel(0), model.RuleType(0), ".bak"}, args{[]byte(`{"level":"prohibited","type":"unspecified","value":".bak"}`)}, false},
		{"invalid input 0/0/.bak", fields{model.StrictnessLevel(1), model.RuleType(1), ".bak"}, args{[]byte(`{"level":"unspecified","type":"unspecified","value":".bak"}`)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &model.Rule{
				Level: tt.fields.Level,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			unMarshaled := &model.Rule{}
			err := unMarshaled.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rule.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr {
				return
			}
			if !reflect.DeepEqual(unMarshaled, r) {
				t.Errorf("Rule.UnmarshalJSON() = %v, want %v", unMarshaled, r)
			}
		})
	}
}

// helper to create a temp directory with convention files
func setupTestConventionsDir(t *testing.T, conventions map[string]interface{}) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "ossify-conventions-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	for filename, content := range conventions {
		filePath := filepath.Join(tempDir, filename)
		switch v := content.(type) {
		case model.Convention:
			data, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("failed to marshal convention: %v", err)
			}
			if err := os.WriteFile(filePath, data, 0644); err != nil {
				t.Fatalf("failed to write convention file: %v", err)
			}
		case string:
			// For invalid JSON or other content
			if err := os.WriteFile(filePath, []byte(v), 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		case nil:
			// Create a directory instead of a file
			if err := os.Mkdir(filePath, 0755); err != nil {
				t.Fatalf("failed to create directory: %v", err)
			}
		}
	}

	return tempDir
}

// helper to save and restore the original ConfigManager
func mockConfigManager(t *testing.T, loadFn config.LoadConfig) func() {
	t.Helper()
	originalManager := config.ConfigManager
	config.ConfigManager = &config.Manager{
		Load: loadFn,
		Save: func(c *config.Config) error { return nil },
	}
	return func() {
		config.ConfigManager = originalManager
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		setupDir       map[string]interface{} // filename -> Convention, string (invalid), or nil (directory)
		configOverride func(dir string) config.LoadConfig
		wantCount      int // number of conventions expected
		wantNames      []string
		wantErr        bool
		errContains    string
	}{
		{
			name:     "returns default conventions when no custom conventions exist",
			setupDir: nil, // no files
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return &config.Config{ConventionPath: dir}, nil
				}
			},
			wantCount: 2, // Standard Distribution + Go
			wantNames: []string{"Standard Distribution", "Go"},
			wantErr:   false,
		},
		{
			name: "loads custom conventions from files",
			setupDir: map[string]interface{}{
				"nodejs.json": model.Convention{
					Name: "Node.js",
					Rules: []model.Rule{
						{Level: model.Required, Type: model.File, Value: "package.json"},
					},
				},
				"python.json": model.Convention{
					Name: "Python",
					Rules: []model.Rule{
						{Level: model.Required, Type: model.File, Value: "setup.py"},
					},
				},
			},
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return &config.Config{ConventionPath: dir}, nil
				}
			},
			wantCount: 4, // 2 defaults + 2 custom
			wantNames: []string{"Standard Distribution", "Go", "Node.js", "Python"},
			wantErr:   false,
		},
		{
			name:     "returns error for empty convention path",
			setupDir: nil,
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return &config.Config{ConventionPath: ""}, nil
				}
			},
			wantCount:   0,
			wantErr:     true,
			errContains: "invalid convention path",
		},
		{
			name: "skips invalid JSON files",
			setupDir: map[string]interface{}{
				"valid.json":   model.Convention{Name: "Valid", Rules: []model.Rule{{Level: model.Required, Type: model.File, Value: "README.md"}}},
				"invalid.json": "this is not valid json {{{",
			},
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return &config.Config{ConventionPath: dir}, nil
				}
			},
			wantCount: 3, // 2 defaults + 1 valid custom
			wantNames: []string{"Standard Distribution", "Go", "Valid"},
			wantErr:   false,
		},
		{
			name: "skips directories in conventions folder",
			setupDir: map[string]interface{}{
				"valid.json": model.Convention{Name: "Valid", Rules: []model.Rule{{Level: model.Required, Type: model.File, Value: "README.md"}}},
				"subdir":     nil, // this creates a directory
			},
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return &config.Config{ConventionPath: dir}, nil
				}
			},
			wantCount: 3, // 2 defaults + 1 valid custom (directory skipped)
			wantNames: []string{"Standard Distribution", "Go", "Valid"},
			wantErr:   false,
		},
		{
			name:     "returns defaults when conventions directory does not exist",
			setupDir: nil,
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return &config.Config{ConventionPath: "/nonexistent/path/that/does/not/exist"}, nil
				}
			},
			wantCount: 2, // only defaults
			wantNames: []string{"Standard Distribution", "Go"},
			wantErr:   false,
		},
		{
			name:     "returns error when config manager fails to load",
			setupDir: nil,
			configOverride: func(dir string) config.LoadConfig {
				return func() (*config.Config, error) {
					return nil, errors.New("config load failed")
				}
			},
			wantCount:   0,
			wantErr:     true,
			errContains: "config load failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tempDir string
			if tt.setupDir != nil {
				tempDir = setupTestConventionsDir(t, tt.setupDir)
				defer os.RemoveAll(tempDir)
			} else {
				// Create an empty temp dir for tests that need a valid path
				var err error
				tempDir, err = os.MkdirTemp("", "ossify-empty-*")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempDir)
			}

			restore := mockConfigManager(t, tt.configOverride(tempDir))
			defer restore()

			got, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && (err == nil || !contains(err.Error(), tt.errContains)) {
					t.Errorf("Load() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if got == nil {
				t.Errorf("Load() returned nil, expected conventions")
				return
			}

			if len(*got) != tt.wantCount {
				t.Errorf("Load() returned %d conventions, want %d", len(*got), tt.wantCount)
			}

			// Check that expected convention names are present
			gotNames := make([]string, len(*got))
			for i, c := range *got {
				gotNames[i] = c.Name
			}

			for _, wantName := range tt.wantNames {
				found := false
				for _, gotName := range gotNames {
					if gotName == wantName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Load() missing expected convention %q, got names: %v", wantName, gotNames)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
