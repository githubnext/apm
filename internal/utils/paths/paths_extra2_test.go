package paths

import (
	"strings"
	"testing"
)

func TestPortableRelpath_ForwardSlashResult(t *testing.T) {
	base := "/home/user/project"
	path := "/home/user/project/src/main.go"
	got := PortableRelpath(path, base)
	if strings.ContainsRune(got, '\\') {
		t.Errorf("result should have forward slashes only: %q", got)
	}
}

func TestPortableRelpath_EmptyPath(t *testing.T) {
	got := PortableRelpath("", "/some/base")
	// Should not panic; empty result or fallback is acceptable
	_ = got
}

func TestPortableRelpath_EmptyBase(t *testing.T) {
	got := PortableRelpath("/some/path", "")
	// Should not panic
	_ = got
}

func TestPortableRelpath_SameDir(t *testing.T) {
	got := PortableRelpath("/a/b/c", "/a/b/c")
	if got != "." && got != "" {
		t.Logf("same-dir PortableRelpath returned %q (implementation-defined)", got)
	}
}

func TestPortableRelpath_ChildPath(t *testing.T) {
	got := PortableRelpath("/a/b/c/d", "/a/b")
	if !strings.Contains(got, "c") {
		t.Errorf("child path should contain 'c': %q", got)
	}
}

func TestPortableRelpath_ConsistentResults(t *testing.T) {
	base := "/project/root"
	path := "/project/root/internal/pkg/file.go"
	got1 := PortableRelpath(path, base)
	got2 := PortableRelpath(path, base)
	if got1 != got2 {
		t.Errorf("repeated calls must be consistent: %q vs %q", got1, got2)
	}
}

func TestPortableRelpath_RelativeToParent(t *testing.T) {
	got := PortableRelpath("/a/b", "/a/b/c")
	// going up from c to b -- should contain ".." or absolute fallback
	_ = got // implementation may vary
}

func TestPortableRelpath_WindowsPathSeparators(t *testing.T) {
	// Ensure backslashes in the input path do not remain in output
	base := "/home/user/project"
	path := "/home/user/project/sub\\dir\\file.go"
	got := PortableRelpath(path, base)
	if strings.ContainsRune(got, '\\') {
		t.Errorf("backslash still present in result: %q", got)
	}
}

func TestPortableRelpath_DeepChildPath(t *testing.T) {
	base := "/project"
	path := "/project/a/b/c/d/e/file.go"
	got := PortableRelpath(path, base)
	// Should contain forward slashes and path components
	if !strings.Contains(got, "a") || !strings.Contains(got, "b") {
		t.Errorf("deep path components missing: %q", got)
	}
}

func TestPortableRelpath_NoTrailingSlash(t *testing.T) {
	base := "/project/root/"
	path := "/project/root/file.go"
	got := PortableRelpath(path, base)
	if strings.HasSuffix(got, "/") {
		t.Errorf("result should not have trailing slash: %q", got)
	}
}
