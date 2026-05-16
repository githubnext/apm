package paths

import (
	"path/filepath"
	"testing"
)

func TestPortableRelpath_Basic(t *testing.T) {
	base := "/home/user/project"
	path := "/home/user/project/src/foo.py"
	got := PortableRelpath(path, base)
	want := "src/foo.py"
	if got != want {
		t.Errorf("PortableRelpath(%q, %q) = %q, want %q", path, base, got, want)
	}
}

func TestPortableRelpath_ForwardSlash(t *testing.T) {
	// On any OS the result should use forward slashes
	tmpDir := t.TempDir()
	sub := filepath.Join(tmpDir, "a", "b", "c.txt")
	got := PortableRelpath(sub, tmpDir)
	for _, c := range got {
		if c == '\\' {
			t.Errorf("PortableRelpath returned backslash: %q", got)
			break
		}
	}
}

func TestPortableRelpath_SamePath(t *testing.T) {
	base := "/home/user/project"
	got := PortableRelpath(base, base)
	if got != "." {
		t.Errorf("PortableRelpath(same, same) = %q, want %q", got, ".")
	}
}
