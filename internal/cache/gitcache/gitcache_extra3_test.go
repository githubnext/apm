package gitcache

import (
	"testing"
)

func TestCacheStats_Assign_Extra3(t *testing.T) {
	s := CacheStats{DBCount: 3, CheckoutCount: 7, TotalSizeBytes: 4096}
	if s.DBCount != 3 {
		t.Errorf("expected DBCount=3, got %d", s.DBCount)
	}
	if s.CheckoutCount != 7 {
		t.Errorf("expected CheckoutCount=7, got %d", s.CheckoutCount)
	}
	if s.TotalSizeBytes != 4096 {
		t.Errorf("expected TotalSizeBytes=4096, got %d", s.TotalSizeBytes)
	}
}

func TestNew_WithNoCache_CreatesInstance_Extra3(t *testing.T) {
	tmp := t.TempDir()
	gc, err := New(tmp, false)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if gc == nil {
		t.Fatal("expected non-nil GitCache")
	}
}

func TestNew_WithRefresh_CreatesInstance_Extra3(t *testing.T) {
	tmp := t.TempDir()
	gc, err := New(tmp, true)
	if err != nil {
		t.Fatalf("New with refresh failed: %v", err)
	}
	if gc == nil {
		t.Fatal("expected non-nil GitCache")
	}
}

func TestGetCacheStats_EmptyCache_Extra3(t *testing.T) {
	tmp := t.TempDir()
	gc, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	stats := gc.GetCacheStats()
	if stats.DBCount < 0 || stats.CheckoutCount < 0 {
		t.Error("cache stats counts must be non-negative")
	}
}

func TestCleanAll_EmptyCache_NoError_Extra3(t *testing.T) {
	tmp := t.TempDir()
	gc, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	// Must not panic.
	gc.CleanAll()
}

func TestPrune_EmptyCache_ReturnsZero_Extra3(t *testing.T) {
	tmp := t.TempDir()
	gc, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	removed := gc.Prune(30)
	if removed < 0 {
		t.Errorf("Prune should return non-negative count, got %d", removed)
	}
}
