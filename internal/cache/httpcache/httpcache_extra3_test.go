package httpcache

import (
	"os"
	"testing"
)

func TestCacheEntry_ZeroValue_Extra3(t *testing.T) {
	var e CacheEntry
	if e.Body != nil || e.StatusCode != 0 || e.ETag != "" {
		t.Error("zero CacheEntry should have nil body, zero status, empty ETag")
	}
}

func TestGetStats_ZeroValue_Extra3(t *testing.T) {
	var s GetStats
	if s.EntryCount != 0 || s.TotalSizeBytes != 0 {
		t.Error("zero GetStats should have all-zero fields")
	}
}

func TestGetStats_AssignFields_Extra3(t *testing.T) {
	s := GetStats{EntryCount: 5, TotalSizeBytes: 1024}
	if s.EntryCount != 5 {
		t.Errorf("expected EntryCount=5, got %d", s.EntryCount)
	}
	if s.TotalSizeBytes != 1024 {
		t.Errorf("expected TotalSizeBytes=1024, got %d", s.TotalSizeBytes)
	}
}

func TestNew_ValidDir_CreatesCache_Extra3(t *testing.T) {
	tmp := t.TempDir()
	hc, err := New(tmp)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if hc == nil {
		t.Fatal("New should return non-nil HttpCache")
	}
}

func TestGetStats_EmptyCache_ZeroStats_Extra3(t *testing.T) {
	tmp := t.TempDir()
	hc, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	s := hc.GetStats()
	if s.EntryCount != 0 || s.TotalSizeBytes != 0 {
		t.Errorf("empty cache should have zero stats, got entryCount=%d totalSize=%d", s.EntryCount, s.TotalSizeBytes)
	}
}

func TestGet_MissingURL_ReturnsNil_Extra3(t *testing.T) {
	tmp := t.TempDir()
	hc, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	entry, err := hc.Get("https://example.com/nonexistent")
	if err != nil {
		t.Fatalf("unexpected error on cache miss: %v", err)
	}
	if entry != nil {
		t.Error("cache miss should return nil entry")
	}
}

func TestCleanAll_EmptyCache_NoError_Extra3(t *testing.T) {
	tmp := t.TempDir()
	hc, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	// Should not panic.
	hc.CleanAll()
}

func TestStore_And_Get_RoundTrip_Extra3(t *testing.T) {
	tmp := t.TempDir()
	hc, err := New(tmp)
	if err != nil {
		t.Fatal(err)
	}
	url := "https://example.com/resource"
	body := []byte("hello world")
	hc.Store(url, body, 200, map[string]string{"cache-control": "max-age=3600"})

	entry, err := hc.Get(url)
	if err != nil {
		t.Fatalf("Get after Store failed: %v", err)
	}
	if entry == nil {
		t.Fatal("expected cached entry, got nil")
	}
	if string(entry.Body) != "hello world" {
		t.Errorf("body mismatch: got %q", string(entry.Body))
	}
	_ = os.Getenv // suppress unused import
}
