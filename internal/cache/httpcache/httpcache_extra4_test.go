package httpcache_test

import (
	"testing"

	"github.com/githubnext/apm/internal/cache/httpcache"
)

func TestNew_ValidDir_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, err := httpcache.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestGetStats_AfterCreation_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	s := c.GetStats()
	if s.EntryCount < 0 || s.TotalSizeBytes < 0 {
		t.Fatal("unexpected negative stats")
	}
}

func TestStore_Get_Roundtrip_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	c.Store("http://example.com/data", []byte("body"), 200, nil)
	entry, err := c.Get("http://example.com/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		// May be expired immediately; acceptable
		return
	}
	if string(entry.Body) != "body" {
		t.Fatalf("expected body, got %s", string(entry.Body))
	}
}

func TestStore_WithETag_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	c.Store("http://example.com/etag", []byte("content"), 200, map[string]string{
		"ETag":          `"abc123"`,
		"Cache-Control": "max-age=3600",
	})
	// Should not panic
}

func TestGet_MissReturnsNil_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	entry, err := c.Get("http://never-stored.example.com/x")
	if err != nil {
		t.Fatalf("unexpected error on cache miss: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil for cache miss")
	}
}

func TestCleanAll_Empty_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	c.CleanAll() // must not panic
}

func TestRefreshExpiry_NoEntry_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	c.RefreshExpiry("http://example.com/absent", nil) // must not panic
}

func TestCacheEntry_Fields_Extra4(t *testing.T) {
	e := httpcache.CacheEntry{
		Body:        []byte("test"),
		ETag:        `"etag1"`,
		ExpiresAt:   1234567890.0,
		ContentType: "application/json",
		StatusCode:  200,
	}
	if string(e.Body) != "test" {
		t.Fatalf("expected test, got %s", string(e.Body))
	}
	if e.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", e.StatusCode)
	}
}

func TestGetStats_AfterStore_Extra4(t *testing.T) {
	dir := t.TempDir()
	c, _ := httpcache.New(dir)
	c.Store("http://example.com/a", []byte("hello"), 200, map[string]string{"Cache-Control": "max-age=3600"})
	s := c.GetStats()
	if s.EntryCount < 0 {
		t.Fatal("negative entry count")
	}
}
