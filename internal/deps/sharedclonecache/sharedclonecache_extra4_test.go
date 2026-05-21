package sharedclonecache

import (
	"testing"
)

func TestNew_ReturnsCacheInstance(t *testing.T) {
	c := New("/tmp/test-cache-e4")
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestNew_DifferentBaseDirs(t *testing.T) {
	c1 := New("/tmp/cache1-e4")
	c2 := New("/tmp/cache2-e4")
	if c1 == c2 {
		t.Error("expected independent instances")
	}
}

func TestGetOrClone_FirstCallInvokesCloneFn(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneCalled := 0
	cloneFn := func(path string) error {
		cloneCalled++
		return nil
	}
	_, err := c.GetOrClone("gh.com", "owner", "repo", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloneCalled != 1 {
		t.Errorf("expected 1 clone call, got %d", cloneCalled)
	}
}

func TestGetOrClone_SecondCallReturnsCached(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneCalled := 0
	cloneFn := func(path string) error {
		cloneCalled++
		return nil
	}
	c.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
	c.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
	if cloneCalled != 1 {
		t.Errorf("expected 1 clone call (cached), got %d", cloneCalled)
	}
}

func TestGetOrClone_DiffRefCallsCloneAgain(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneCalled := 0
	cloneFn := func(path string) error {
		cloneCalled++
		return nil
	}
	c.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
	c.GetOrClone("gh.com", "o", "r", "v1.0.0", cloneFn, nil)
	if cloneCalled != 2 {
		t.Errorf("expected 2 clone calls for different refs, got %d", cloneCalled)
	}
}

func TestGetOrClone_ReturnsNonEmptyPath(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	path, err := c.GetOrClone("gh.com", "o", "r", "main", func(p string) error { return nil }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetOrClone_DifferentOwnersDifferentPaths(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	fn := func(p string) error { return nil }
	p1, _ := c.GetOrClone("gh.com", "owner1", "repo", "main", fn, nil)
	p2, _ := c.GetOrClone("gh.com", "owner2", "repo", "main", fn, nil)
	if p1 == p2 {
		t.Error("expected different paths for different owners")
	}
}
