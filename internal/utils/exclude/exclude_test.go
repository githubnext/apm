package exclude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/exclude"
)

func TestValidateExcludePatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		wantErr  bool
		wantLen  int
	}{
		{"nil input", nil, false, 0},
		{"empty", []string{}, false, 0},
		{"simple", []string{"foo/bar"}, false, 1},
		{"collapses consecutive **", []string{"**/**/foo"}, false, 1},
		{"too many stars", []string{"a/**/b/**/c/**/d/**/e/**/f/**/g"}, true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exclude.ValidateExcludePatterns(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err=%v got %v", tt.wantErr, err)
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Fatalf("want len=%d got %d", tt.wantLen, len(got))
			}
		})
	}
}

func TestShouldExclude(t *testing.T) {
	base, err := os.MkdirTemp("", "exclude-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	tests := []struct {
		file     string
		patterns []string
		want     bool
	}{
		{"docs/foo.md", []string{"docs/"}, true},
		{"src/main.go", []string{"docs/"}, false},
		{"build/output.bin", []string{"build/*"}, true},
		{"src/a/b/c.py", []string{"src/**/*.py"}, true},
		{"src/a/b/c.go", []string{"src/**/*.py"}, false},
	}

	for _, tt := range tests {
		full := filepath.Join(base, filepath.FromSlash(tt.file))
		os.MkdirAll(filepath.Dir(full), 0o755)
		if got := exclude.ShouldExclude(full, base, tt.patterns); got != tt.want {
			t.Errorf("ShouldExclude(%q, patterns=%v) = %v, want %v", tt.file, tt.patterns, got, tt.want)
		}
	}
}
