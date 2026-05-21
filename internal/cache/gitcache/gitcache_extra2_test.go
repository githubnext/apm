package gitcache

import (
	"os"
	"testing"
)

func TestCacheStats_ZeroValue(t *testing.T) {
	var s CacheStats
	if s.DBCount != 0 || s.CheckoutCount != 0 || s.TotalSizeBytes != 0 {
		t.Error("zero CacheStats should have all zero fields")
	}
}

func TestCacheStats_AssignFields(t *testing.T) {
	s := CacheStats{DBCount: 5, CheckoutCount: 12, TotalSizeBytes: 1024}
	if s.DBCount != 5 {
		t.Errorf("expected DBCount=5, got %d", s.DBCount)
	}
	if s.CheckoutCount != 12 {
		t.Errorf("expected CheckoutCount=12, got %d", s.CheckoutCount)
	}
	if s.TotalSizeBytes != 1024 {
		t.Errorf("expected TotalSizeBytes=1024, got %d", s.TotalSizeBytes)
	}
}

func TestNew_RefreshFalse(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil GitCache")
	}
}

func TestNew_RefreshTrue(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil GitCache")
	}
}

func TestNew_InvalidPath(t *testing.T) {
	// Try to create cache in a path that can't be a directory (a file's child)
	f, err := os.CreateTemp("", "gitcache-test-*.tmp")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = New(f.Name()+"/sub", false)
	if err == nil {
		t.Error("expected error when cacheRoot is under a file")
	}
}

func TestGetCacheStats_EmptyCache(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	stats := c.GetCacheStats()
	if stats.DBCount != 0 {
		t.Errorf("expected DBCount=0, got %d", stats.DBCount)
	}
	if stats.CheckoutCount != 0 {
		t.Errorf("expected CheckoutCount=0, got %d", stats.CheckoutCount)
	}
}

func TestCleanAll_EmptyCache(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	c.CleanAll()
	stats := c.GetCacheStats()
	if stats.DBCount != 0 || stats.CheckoutCount != 0 {
		t.Error("stats should be zero after CleanAll on empty cache")
	}
}

func TestPrune_ZeroPrunedOnEmpty(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	pruned := c.Prune(30)
	if pruned != 0 {
		t.Errorf("expected 0 pruned on empty cache, got %d", pruned)
	}
}
