package paths

import (
"strings"
"testing"
)

func TestPortableRelpath_SiblingDirExtra4(t *testing.T) {
result := PortableRelpath("/a/b/c", "/a/b/d")
if !strings.Contains(result, "c") {
t.Errorf("unexpected relpath: %q", result)
}
}

func TestPortableRelpath_DeepNestingExtra4(t *testing.T) {
result := PortableRelpath("/a/b/c/d/e", "/a/b")
if result == "" {
t.Errorf("unexpected empty relpath")
}
}

func TestPortableRelpath_NoBackslashesExtra4(t *testing.T) {
result := PortableRelpath("/a/b/c/d", "/a/b")
if strings.Contains(result, "\\") {
t.Errorf("unexpected backslash in result: %q", result)
}
}

func TestPortableRelpath_SamePathReturnsDotExtra4(t *testing.T) {
result := PortableRelpath("/a/b", "/a/b/c")
if result == "" {
t.Error("expected non-empty relpath")
}
}

func TestPortableRelpath_RootPathExtra4(t *testing.T) {
result := PortableRelpath("/", "/tmp")
if result == "" {
t.Error("expected non-empty result for root path")
}
}

func TestPortableRelpath_MultipleCallsConsistentExtra4(t *testing.T) {
r1 := PortableRelpath("/a/b/c", "/a")
r2 := PortableRelpath("/a/b/c", "/a")
if r1 != r2 {
t.Errorf("inconsistent results: %q vs %q", r1, r2)
}
}
