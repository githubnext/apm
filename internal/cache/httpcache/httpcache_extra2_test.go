package httpcache

import (
	"os"
	"testing"
)

func TestNew_InvalidDir_Errors(t *testing.T) {
	// A file instead of a directory should cause an error.
	f, err := os.CreateTemp("", "httpcache-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()
	// Try to create an httpcache with the file as cacheRoot (http_v1 subdir creation should fail).
	_, err = New(f.Name())
	if err == nil {
		t.Error("expected error when cacheRoot is a file, not a dir")
	}
}

func TestStore_EmptyBody(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	c.Store("https://example.com/empty", []byte{}, 200, map[string]string{})
	entry, err := c.Get("https://example.com/empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected cached entry for empty body")
	}
	if len(entry.Body) != 0 {
		t.Errorf("expected empty body, got %d bytes", len(entry.Body))
	}
}

func TestStore_LargeBody(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	body := make([]byte, 64*1024) // 64KB
	for i := range body {
		body[i] = byte(i % 256)
	}
	c.Store("https://example.com/large", body, 200, map[string]string{})
	entry, err := c.Get("https://example.com/large")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected cached entry")
	}
	if len(entry.Body) != len(body) {
		t.Errorf("expected %d bytes, got %d", len(body), len(entry.Body))
	}
}

func TestRefreshExpiry_DoesNotPanic(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	c.Store("https://example.com/refresh", []byte("hello"), 200, map[string]string{})
	// Should not panic even with empty headers.
	c.RefreshExpiry("https://example.com/refresh", map[string]string{})
}

func TestRefreshExpiry_NonExistentURL(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	// Should not panic for a URL that was never stored.
	c.RefreshExpiry("https://not-stored.example.com/path", map[string]string{})
}

func TestCleanAll_EmptiesCache(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	c.Store("https://example.com/to-clean", []byte("data"), 200, map[string]string{})
	c.CleanAll()
	entry, err := c.Get("https://example.com/to-clean")
	if err != nil {
		t.Fatalf("unexpected error after CleanAll: %v", err)
	}
	if entry != nil {
		t.Error("expected no cached entry after CleanAll")
	}
}

func TestGetStats_EmptyCache(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	stats := c.GetStats()
	if stats.EntryCount != 0 {
		t.Errorf("expected 0 entries for empty cache, got %d", stats.EntryCount)
	}
}

func TestGetStats_AfterCleanAll(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	c.Store("https://example.com/x", []byte("abc"), 200, map[string]string{})
	c.CleanAll()
	stats := c.GetStats()
	if stats.EntryCount != 0 {
		t.Errorf("expected 0 entries after CleanAll, got %d", stats.EntryCount)
	}
}

func TestStore_ContentTypePreserved(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	c.Store("https://example.com/ct", []byte("{}"), 200, map[string]string{
		"Content-Type": "application/json",
	})
	entry, err := c.Get("https://example.com/ct")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry")
	}
	if entry.ContentType != "application/json" {
		t.Errorf("expected application/json, got %q", entry.ContentType)
	}
}

func TestCacheEntry_Fields(t *testing.T) {
	e := CacheEntry{
		Body:        []byte("body"),
		ETag:        "etag123",
		ExpiresAt:   1234567.0,
		ContentType: "text/plain",
		StatusCode:  200,
	}
	if string(e.Body) != "body" {
		t.Error("Body mismatch")
	}
	if e.ETag != "etag123" {
		t.Error("ETag mismatch")
	}
	if e.StatusCode != 200 {
		t.Error("StatusCode mismatch")
	}
}
