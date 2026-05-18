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

func TestValidatePathSegments_allowCurrentDir(t *testing.T) {
	// "." is allowed when allowCurrentDir=true.
	if err := pathsecurity.ValidatePathSegments("./foo/bar", "test", false, true); err != nil {
		t.Errorf("unexpected error with allowCurrentDir=true: %v", err)
	}
	// But ".." is still rejected.
	if err := pathsecurity.ValidatePathSegments("./foo/../bar", "test", false, true); err == nil {
		t.Error("expected error for '..' even with allowCurrentDir=true")
	}
}

func TestValidatePathSegments_rejectEmpty(t *testing.T) {
	// Double slash creates empty segment.
	if err := pathsecurity.ValidatePathSegments("foo//bar", "test", true, false); err == nil {
		t.Error("expected error for double slash with rejectEmpty=true")
	}
	if err := pathsecurity.ValidatePathSegments("foo/bar", "test", true, false); err != nil {
		t.Errorf("unexpected error for clean path: %v", err)
	}
}

func TestValidatePathSegments_percentEncoded(t *testing.T) {
	// Percent-encoded ".." should still be rejected.
	if err := pathsecurity.ValidatePathSegments("foo/%2e%2e/bar", "test", false, false); err == nil {
		t.Error("expected error for percent-encoded traversal")
	}
}

func TestIsPathTraversalError(t *testing.T) {
	err := pathsecurity.ValidatePathSegments("../etc", "ctx", false, false)
	if err == nil {
		t.Fatal("expected error")
	}
	if !pathsecurity.IsPathTraversalError(err) {
		t.Errorf("expected PathTraversalError, got %T", err)
	}
}

func TestEnsurePathWithin_deepNesting(t *testing.T) {
	base, err := os.MkdirTemp("", "pathsec-deep")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	deep := filepath.Join(base, "a", "b", "c", "d", "file.txt")
	os.MkdirAll(filepath.Dir(deep), 0o755)
	os.WriteFile(deep, []byte("deep"), 0o644)

	if _, err := pathsecurity.EnsurePathWithin(deep, base); err != nil {
		t.Errorf("expected deep nested path to pass: %v", err)
	}
}
