package sharedclonecache

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestNew_NilFields(t *testing.T) {
	c := New(t.TempDir())
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestCleanup_EmptyCache(t *testing.T) {
	c := New(t.TempDir())
	// Should not panic on empty cleanup.
	c.Cleanup()
}

func TestCleanup_AfterClone(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	_, err := c.GetOrClone("gh.com", "owner", "cleanup-repo", "main", cloneFn, nil)
	if err != nil {
		t.Fatalf("GetOrClone failed: %v", err)
	}
	// Cleanup removes the temp dirs.
	c.Cleanup()
}

func TestGetOrClone_SameKeyReturnsCache(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	calls := 0
	cloneFn := func(clonePath string) error {
		calls++
		return os.MkdirAll(clonePath, 0o755)
	}
	path1, err1 := c.GetOrClone("gh.com", "owner", "idempotent-repo", "main", cloneFn, nil)
	if err1 != nil {
		t.Fatalf("first call error: %v", err1)
	}
	path2, err2 := c.GetOrClone("gh.com", "owner", "idempotent-repo", "main", cloneFn, nil)
	if err2 != nil {
		t.Fatalf("second call error: %v", err2)
	}
	if path1 != path2 {
		t.Errorf("expected same path on second call; got %q and %q", path1, path2)
	}
	if calls != 1 {
		t.Errorf("expected cloneFn called once, got %d", calls)
	}
}

func TestGetOrClone_DifferentRefsProduceDifferentEntries(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	calls := 0
	cloneFn := func(clonePath string) error {
		calls++
		return os.MkdirAll(clonePath, 0o755)
	}
	path1, _ := c.GetOrClone("gh.com", "owner", "multi-ref", "v1", cloneFn, nil)
	path2, _ := c.GetOrClone("gh.com", "owner", "multi-ref", "v2", cloneFn, nil)
	if path1 == path2 {
		t.Errorf("expected different paths for different refs, both got %q", path1)
	}
	if calls != 2 {
		t.Errorf("expected 2 clone calls, got %d", calls)
	}
}

func TestGetOrClone_ConcurrentCallsSameKey_OnlyOneClone(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	var mu sync.Mutex
	cloneCount := 0
	cloneFn := func(clonePath string) error {
		mu.Lock()
		cloneCount++
		mu.Unlock()
		return os.MkdirAll(clonePath, 0o755)
	}

	const N = 8
	paths := make([]string, N)
	errs := make([]error, N)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			paths[idx], errs[idx] = c.GetOrClone("gh.com", "owner", "concurrent-repo", "main", cloneFn, nil)
		}(i)
	}
	wg.Wait()

	for i, e := range errs {
		if e != nil {
			t.Errorf("goroutine %d error: %v", i, e)
		}
	}
	for i := 1; i < N; i++ {
		if paths[i] != paths[0] {
			t.Errorf("goroutine %d got path %q, expected %q", i, paths[i], paths[0])
		}
	}
	if cloneCount != 1 {
		t.Errorf("expected 1 clone, got %d", cloneCount)
	}
}

func TestGetOrClone_FetchFn_UsesExistingBare(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)

	// First clone for v1 to establish the bare.
	cloneFn := func(clonePath string) error {
		return os.MkdirAll(clonePath, 0o755)
	}
	_, err := c.GetOrClone("gh.com", "owner", "fetch-repo", "v1", cloneFn, nil)
	if err != nil {
		t.Fatalf("first clone: %v", err)
	}

	fetchCalled := false
	fetchFn := func(barePath string, sha string) bool {
		fetchCalled = true
		return filepath.IsAbs(barePath)
	}
	// Second call for same repo, different ref, with fetchFn.
	_, err2 := c.GetOrClone("gh.com", "owner", "fetch-repo", "v2", cloneFn, fetchFn)
	if err2 != nil {
		t.Fatalf("second clone: %v", err2)
	}
	if !fetchCalled {
		t.Error("expected fetchFn to be called when existing bare is available")
	}
}
