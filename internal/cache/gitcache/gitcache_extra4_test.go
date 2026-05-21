package gitcache_test

import (
	"testing"

	"github.com/githubnext/apm/internal/cache/gitcache"
)

func TestCacheStats_ZeroFields_Extra4(t *testing.T) {
	var s gitcache.CacheStats
	if s.DBCount != 0 || s.CheckoutCount != 0 || s.TotalSizeBytes != 0 {
		t.Fatal("expected zero-value CacheStats")
	}
}

func TestNew_ValidDir_ReturnsCache_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestNew_RefreshTrue_ReturnsCache_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestGetCacheStats_Empty_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stats := c.GetCacheStats()
	if stats.DBCount < 0 {
		t.Fatal("expected non-negative DBCount")
	}
}

func TestCleanAll_Empty_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.CleanAll() // must not panic
}

func TestPrune_ZeroDays_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n := c.Prune(0)
	if n < 0 {
		t.Fatal("expected non-negative count")
	}
}

func TestPrune_LargeDays_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n := c.Prune(9999)
	if n != 0 {
		t.Fatalf("expected 0 pruned in empty cache, got %d", n)
	}
}

func TestEvictCheckout_NonExistent_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := gitcache.New(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.EvictCheckout("/nonexistent/checkout/dir") // must not panic
}

func TestCacheStats_Mutation_Extra4(t *testing.T) {
	s := gitcache.CacheStats{DBCount: 5, CheckoutCount: 3, TotalSizeBytes: 1024}
	if s.DBCount != 5 || s.CheckoutCount != 3 || s.TotalSizeBytes != 1024 {
		t.Fatal("field assignment failed")
	}
}
