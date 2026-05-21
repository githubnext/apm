package paths

import (
	"testing"
)

func TestPortableRelpath_NoBackslashes(t *testing.T) {
	got := PortableRelpath("/a/b/c", "/a/b")
	for _, ch := range got {
		if ch == '\\' {
			t.Errorf("result must not contain backslash, got %q", got)
		}
	}
}

func TestPortableRelpath_ChildIsRelative(t *testing.T) {
	got := PortableRelpath("/a/b/c/d", "/a/b")
	if got != "c/d" {
		t.Errorf("expected c/d, got %q", got)
	}
}

func TestPortableRelpath_SamePathIsDot(t *testing.T) {
	got := PortableRelpath("/a/b", "/a/b")
	if got != "." {
		t.Errorf("expected ., got %q", got)
	}
}

func TestPortableRelpath_ParentReturnsDoubleDot(t *testing.T) {
	got := PortableRelpath("/a", "/a/b")
	if got != ".." {
		t.Errorf("expected .., got %q", got)
	}
}

func TestPortableRelpath_SiblingDirVariant(t *testing.T) {
	got := PortableRelpath("/a/c", "/a/b")
	if got != "../c" {
		t.Errorf("expected ../c, got %q", got)
	}
}

func TestPortableRelpath_DeepNestingVariant(t *testing.T) {
	got := PortableRelpath("/a/b/c/d/e", "/a/b")
	if got != "c/d/e" {
		t.Errorf("expected c/d/e, got %q", got)
	}
}

func TestPortableRelpath_ReturnsString(t *testing.T) {
	got := PortableRelpath("/foo/bar", "/foo")
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestPortableRelpath_ConsistentMultipleCalls(t *testing.T) {
	a := PortableRelpath("/a/b/c", "/a")
	b := PortableRelpath("/a/b/c", "/a")
	if a != b {
		t.Error("PortableRelpath must be deterministic")
	}
}
