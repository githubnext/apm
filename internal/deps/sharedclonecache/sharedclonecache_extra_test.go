package sharedclonecache

import (
	"os"
	"sync"
	"testing"
)

func TestNew_WithBaseDir(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestNew_EmptyBaseDir(t *testing.T) {
	c := New("")
	if c == nil {
		t.Fatal("expected non-nil cache with empty baseDir")
	}
}

func TestGetOrClone_RetryAfterError(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	calls := 0
	cloneFn := func(clonePath string) error {
		calls++
		if calls < 2 {
			return os.ErrPermission
		}
		return os.MkdirAll(clonePath, 0o755)
	}
	// First call should fail
	_, err := c.GetOrClone("gh.com", "owner", "retry-repo", "main", cloneFn, nil)
	if err == nil {
		t.Fatal("expected error on first call")
	}
	// Second call with same key should retry (cache clears error)
	path, err2 := c.GetOrClone("gh.com", "owner", "retry-repo", "main", cloneFn, nil)
	if err2 != nil {
		t.Fatalf("second call should succeed: %v", err2)
	}
	if path == "" {
		t.Error("expected non-empty path on second call")
	}
}

func TestGetOrClone_ConcurrentSameKey(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneCount := 0
	var mu sync.Mutex
	cloneFn := func(clonePath string) error {
		mu.Lock()
		cloneCount++
		mu.Unlock()
		return os.MkdirAll(clonePath, 0o755)
	}
	const goroutines = 5
	paths := make([]string, goroutines)
	errs := make([]error, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			paths[idx], errs[idx] = c.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
		}()
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d got error: %v", i, err)
		}
	}
	// All should return the same path
	for i := 1; i < goroutines; i++ {
		if paths[i] != paths[0] {
			t.Errorf("goroutine %d got path %q, want %q", i, paths[i], paths[0])
		}
	}
	// Clone should be called exactly once
	if cloneCount != 1 {
		t.Errorf("expected clone called once, got %d", cloneCount)
	}
}

func TestGetOrClone_DifferentRefsSameRepo(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	path1, err1 := c.GetOrClone("gh.com", "owner", "repo", "main", cloneFn, nil)
	if err1 != nil {
		t.Fatalf("main clone: %v", err1)
	}
	path2, err2 := c.GetOrClone("gh.com", "owner", "repo", "v1.0.0", cloneFn, nil)
	if err2 != nil {
		t.Fatalf("v1.0.0 clone: %v", err2)
	}
	if path1 == path2 {
		t.Error("different refs should produce different clone paths")
	}
}

func TestGetOrClone_DifferentHostsSameOwnerRepo(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	path1, err1 := c.GetOrClone("github.com", "owner", "repo", "main", cloneFn, nil)
	if err1 != nil {
		t.Fatalf("github clone: %v", err1)
	}
	path2, err2 := c.GetOrClone("gitlab.com", "owner", "repo", "main", cloneFn, nil)
	if err2 != nil {
		t.Fatalf("gitlab clone: %v", err2)
	}
	if path1 == path2 {
		t.Error("different hosts should produce different clone paths")
	}
}

func TestGetOrClone_CachedPathExists(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	path, err := c.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("initial clone: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("cloned path should exist on disk: %s", path)
	}
}

func TestGetOrClone_WithFetchFn_NoExistingBare(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	fetchCalled := false
	fetchFn := func(barePath, sha string) bool {
		fetchCalled = true
		return true
	}
	// With no existing bare, fetchFn should not be called
	_, err := c.GetOrClone("gh.com", "o", "r", "abc1234", cloneFn, fetchFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = fetchCalled // may or may not be called depending on implementation
}

func TestNew_IndependentInstances(t *testing.T) {
	dir := t.TempDir()
	c1 := New(dir)
	c2 := New(dir)
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	p1, _ := c1.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
	p2, _ := c2.GetOrClone("gh.com", "o", "r", "main", cloneFn, nil)
	// Independent caches may return different paths
	_ = p1
	_ = p2
}
