package pathsecurity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/pathsecurity"
)

func TestValidatePathSegments(t *testing.T) {
	tests := []struct {
		path     string
		wantErr  bool
	}{
		{"foo/bar/baz", false},
		{"../etc/passwd", true},
		{"foo/../etc", true},
		{"./relative", true},
		{"foo/bar", false},
	}
	for _, tt := range tests {
		err := pathsecurity.ValidatePathSegments(tt.path, "test", false, false)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePathSegments(%q) err=%v, wantErr=%v", tt.path, err, tt.wantErr)
		}
	}
}

func TestEnsurePathWithin(t *testing.T) {
	base, err := os.MkdirTemp("", "pathsec-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	safe := filepath.Join(base, "subdir", "file.txt")
	os.MkdirAll(filepath.Dir(safe), 0o755)
	os.WriteFile(safe, []byte("x"), 0o644)

	if _, err := pathsecurity.EnsurePathWithin(safe, base); err != nil {
		t.Errorf("expected safe path to pass, got err: %v", err)
	}

	if _, err := pathsecurity.EnsurePathWithin("/etc/passwd", base); err == nil {
		t.Error("expected /etc/passwd to fail containment check")
	}
}
