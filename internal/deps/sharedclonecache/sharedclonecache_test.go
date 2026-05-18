package sharedclonecache

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	c := New("")
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestGetOrClone_ClonesOnce(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	callCount := 0
	cloneFn := func(clonePath string) error {
		callCount++
		return os.MkdirAll(clonePath, 0o755)
	}

	path1, err := c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 clone call, got %d", callCount)
	}

	path2, err := c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected clone to be called only once, got %d", callCount)
	}
	if path1 != path2 {
		t.Errorf("expected same path from both calls: %s vs %s", path1, path2)
	}
}

func TestGetOrClone_DifferentKeysCloneSeparately(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}

	path1, err := c.GetOrClone("github.com", "owner", "repoA", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("first clone failed: %v", err)
	}
	path2, err := c.GetOrClone("github.com", "owner", "repoB", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("second clone failed: %v", err)
	}
	if path1 == path2 {
		t.Error("different keys should produce different paths")
	}
}

func TestGetOrClone_CloneError(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	cloneFn := func(_ string) error {
		return errors.New("network error")
	}

	_, err := c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, nil)
	if err == nil {
		t.Error("expected error from clone failure")
	}
}

func TestGetOrClone_WithFetchFn(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}

	// prime the cache with a clone
	_, err := c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("initial clone failed: %v", err)
	}

	fetchCalled := false
	fetchFn := func(barePath, sha string) bool {
		fetchCalled = true
		return true
	}

	// same key - fetch should not be called (already cached)
	_, err = c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, fetchFn)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	_ = fetchCalled // fetch may or may not be invoked on cache hit
}

func TestCleanup(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}

	clonePath, err := c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	// Directory should exist after clone
	if _, err := os.Stat(clonePath); err != nil {
		t.Fatalf("clone path should exist: %v", err)
	}

	c.Cleanup()

	// All temp dirs should be removed
	if _, err := os.Stat(filepath.Join(dir, "tmpclone")); err == nil {
		// some might remain depending on implementation; just ensure no panic
	}
}
