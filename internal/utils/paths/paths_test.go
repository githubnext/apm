package paths

import (
	"path/filepath"
	"strings"
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

func TestPortableRelpath_DeepNesting(t *testing.T) {
	tmpDir := t.TempDir()
	sub := filepath.Join(tmpDir, "a", "b", "c", "d", "e.go")
	got := PortableRelpath(sub, tmpDir)
	if strings.Contains(got, "\\") {
		t.Errorf("result contains backslash: %q", got)
	}
	want := "a/b/c/d/e.go"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPortableRelpath_ParentDir(t *testing.T) {
	tmpDir := t.TempDir()
	// path is the parent of base
	sub := filepath.Join(tmpDir, "child")
	got := PortableRelpath(tmpDir, sub)
	if strings.Contains(got, "\\") {
		t.Errorf("result contains backslash: %q", got)
	}
	// Should be ".." since tmpDir is parent of sub
	if got != ".." {
		t.Errorf("got %q, want ..", got)
	}
}

func TestPortableRelpath_RealPaths(t *testing.T) {
	tmpDir := t.TempDir()
	cases := []struct {
		rel  string
		want string
	}{
		{"foo.py", "foo.py"},
		{filepath.Join("src", "bar.go"), "src/bar.go"},
		{filepath.Join("tests", "unit", "test_x.py"), "tests/unit/test_x.py"},
	}
	for _, c := range cases {
		full := filepath.Join(tmpDir, c.rel)
		got := PortableRelpath(full, tmpDir)
		if got != c.want {
			t.Errorf("PortableRelpath(%q) = %q, want %q", c.rel, got, c.want)
		}
	}
}

func TestPortableRelpath_NoBackslashInResult(t *testing.T) {
	tmpDir := t.TempDir()
	paths := []string{
		filepath.Join(tmpDir, "a.go"),
		filepath.Join(tmpDir, "sub", "b.go"),
		filepath.Join(tmpDir, "x", "y", "z.go"),
	}
	for _, p := range paths {
		got := PortableRelpath(p, tmpDir)
		if strings.ContainsRune(got, '\\') {
			t.Errorf("PortableRelpath(%q) contains backslash: %q", p, got)
		}
	}
}
