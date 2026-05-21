package sharedclonecache

import (
	"fmt"
	"testing"
)

func TestNew_ReturnsCacheE3(t *testing.T) {
	c := New("")
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestNew_WithBaseDirE3(t *testing.T) {
	c := New("/tmp")
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestGetOrClone_NilFetchFnOk(t *testing.T) {
	c := New("")
	path, err := c.GetOrClone("github.com", "owner", "repo", "main",
		func(clonePath string) error { return nil },
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetOrClone_ErrorPropagated(t *testing.T) {
	c := New("")
	_, err := c.GetOrClone("github.com", "owner2", "repo2", "main",
		func(clonePath string) error { return fmt.Errorf("clone failed") },
		nil,
	)
	if err == nil {
		t.Error("expected error from clone fn")
	}
}

func TestGetOrClone_DiffRefsAreDifferentKeys(t *testing.T) {
	c := New("")
	var calls int
	cloneFn := func(clonePath string) error {
		calls++
		return nil
	}
	_, _ = c.GetOrClone("h", "o", "r", "ref1", cloneFn, nil)
	_, _ = c.GetOrClone("h", "o", "r", "ref2", cloneFn, nil)
	if calls != 2 {
		t.Errorf("expected 2 clone calls for different refs, got %d", calls)
	}
}

func TestGetOrClone_SameRefCached(t *testing.T) {
	c := New("")
	var calls int
	cloneFn := func(clonePath string) error {
		calls++
		return nil
	}
	_, _ = c.GetOrClone("h", "o", "r", "main", cloneFn, nil)
	_, _ = c.GetOrClone("h", "o", "r", "main", cloneFn, nil)
	if calls != 1 {
		t.Errorf("expected 1 clone call (cached), got %d", calls)
	}
}

func TestCleanup_DoesNotPanic(t *testing.T) {
	c := New("")
	c.Cleanup()
}

func TestGetOrClone_IndependentCaches(t *testing.T) {
	c1 := New("")
	c2 := New("")
	var calls int
	cloneFn := func(_ string) error { calls++; return nil }
	_, _ = c1.GetOrClone("h", "o", "r", "ref", cloneFn, nil)
	_, _ = c2.GetOrClone("h", "o", "r", "ref", cloneFn, nil)
	if calls != 2 {
		t.Errorf("expected 2 clone calls for independent caches, got %d", calls)
	}
}

func TestNew_ZeroEntries(t *testing.T) {
	c := New("")
	c2 := New("alt")
	if c == c2 {
		t.Error("expected distinct instances")
	}
}
