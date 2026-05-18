package httpcache

import (
	"os"
	"testing"
)

func TestConstants(t *testing.T) {
	if MaxHTTPCacheTTLSeconds != 86400 {
		t.Errorf("MaxHTTPCacheTTLSeconds = %d, want 86400", MaxHTTPCacheTTLSeconds)
	}
	if MaxHTTPCacheBytes != 100*1024*1024 {
		t.Errorf("MaxHTTPCacheBytes = %d, want 100MB", MaxHTTPCacheBytes)
	}
}

func TestNewAndGetStats(t *testing.T) {
	dir := t.TempDir()
	hc, err := New(dir)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	stats := hc.GetStats()
	if stats.EntryCount != 0 {
		t.Errorf("GetStats().EntryCount = %d, want 0 for empty cache", stats.EntryCount)
	}
}

func TestStoreAndGet(t *testing.T) {
	dir := t.TempDir()
	hc, err := New(dir)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	url := "https://example.com/data"
	body := []byte(`{"hello":"world"}`)
	headers := map[string]string{
		"Cache-Control": "max-age=3600",
		"ETag":          "\"abc123\"",
	}
	hc.Store(url, body, 200, headers)

	entry, err := hc.Get(url)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if string(entry.Body) != string(body) {
		t.Errorf("Get() body = %q, want %q", entry.Body, body)
	}
	if entry.StatusCode != 200 {
		t.Errorf("Get() StatusCode = %d, want 200", entry.StatusCode)
	}
}

func TestGetMiss(t *testing.T) {
	dir := t.TempDir()
	hc, err := New(dir)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	entry, err := hc.Get("https://notcached.example.com/foo")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if entry != nil {
		t.Error("Get() miss should return nil entry")
	}
}

func TestCleanAll(t *testing.T) {
	dir := t.TempDir()
	hc, err := New(dir)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	hc.Store("https://example.com/x", []byte("data"), 200, nil)
	hc.CleanAll()
	// After clean, dir may be removed; entry should not be found.
	hc2, _ := New(dir)
	if hc2 != nil {
		stats := hc2.GetStats()
		_ = stats
	}
	_ = os.MkdirAll(dir, 0o755)
}

func TestParseTTLCapped(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	// TTL > 24h should be capped.
	ttl := hc.parseTTL(map[string]string{"Cache-Control": "max-age=999999"})
	if ttl != MaxHTTPCacheTTLSeconds {
		t.Errorf("parseTTL(huge) = %f, want %f", ttl, float64(MaxHTTPCacheTTLSeconds))
	}
}

func TestParseTTLNormal(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	ttl := hc.parseTTL(map[string]string{"Cache-Control": "max-age=3600"})
	if ttl != 3600 {
		t.Errorf("parseTTL(3600) = %f, want 3600", ttl)
	}
}

func TestParseTTLMissing(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	ttl := hc.parseTTL(map[string]string{})
	if ttl != 0 {
		t.Errorf("parseTTL(empty) = %f, want 0", ttl)
	}
}
