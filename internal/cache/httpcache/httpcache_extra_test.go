package httpcache

import (
	"testing"
)

func TestNew_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	hc, err := New(dir)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if hc == nil {
		t.Error("New() returned nil")
	}
}

func TestStore_ThenGetReturnsBody(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	url := "https://api.example.com/resource"
	body := []byte("response body content")
	hc.Store(url, body, 200, nil)
	entry, err := hc.Get(url)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected cache hit, got nil")
	}
	if string(entry.Body) != string(body) {
		t.Errorf("body = %q, want %q", entry.Body, body)
	}
}

func TestStore_ETagPreserved(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	url := "https://example.com/etag"
	headers := map[string]string{"ETag": "\"xyz789\""}
	hc.Store(url, []byte("data"), 200, headers)
	entry, _ := hc.Get(url)
	if entry == nil {
		t.Fatal("expected cache hit")
	}
	if entry.ETag != "\"xyz789\"" {
		t.Errorf("ETag = %q, want \"xyz789\"", entry.ETag)
	}
}

func TestStore_StatusCodePreserved(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	url := "https://example.com/status"
	hc.Store(url, []byte("ok"), 201, nil)
	entry, _ := hc.Get(url)
	if entry == nil {
		t.Fatal("expected cache hit")
	}
	if entry.StatusCode != 201 {
		t.Errorf("StatusCode = %d, want 201", entry.StatusCode)
	}
}

func TestGet_MissingKey(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	entry, err := hc.Get("https://never-stored.example.com/key")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if entry != nil {
		t.Error("expected nil entry for cache miss")
	}
}

func TestGetStats_AfterStore(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	hc.Store("https://example.com/a", []byte("aaa"), 200, nil)
	hc.Store("https://example.com/b", []byte("bbb"), 200, nil)
	stats := hc.GetStats()
	if stats.EntryCount < 2 {
		t.Errorf("EntryCount = %d, want >= 2", stats.EntryCount)
	}
	if stats.TotalSizeBytes <= 0 {
		t.Errorf("TotalSizeBytes = %d, want > 0", stats.TotalSizeBytes)
	}
}

func TestParseTTL_Zero(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	ttl := hc.parseTTL(nil)
	if ttl != 0 {
		t.Errorf("parseTTL(nil) = %f, want 0", ttl)
	}
}

func TestParseTTL_SmallValue(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	ttl := hc.parseTTL(map[string]string{"Cache-Control": "max-age=60"})
	if ttl != 60 {
		t.Errorf("parseTTL(60) = %f, want 60", ttl)
	}
}

func TestParseTTL_Exact24h(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	ttl := hc.parseTTL(map[string]string{"Cache-Control": "max-age=86400"})
	if ttl != 86400 {
		t.Errorf("parseTTL(86400) = %f, want 86400", ttl)
	}
}

func TestParseTTL_Exceeds24h(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	ttl := hc.parseTTL(map[string]string{"Cache-Control": "max-age=99999"})
	if ttl != MaxHTTPCacheTTLSeconds {
		t.Errorf("parseTTL(99999) = %f, want %d (capped)", ttl, MaxHTTPCacheTTLSeconds)
	}
}

func TestStore_MultipleURLs(t *testing.T) {
	dir := t.TempDir()
	hc, _ := New(dir)
	urls := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	}
	for i, u := range urls {
		hc.Store(u, []byte{byte(i)}, 200, nil)
	}
	for _, u := range urls {
		entry, err := hc.Get(u)
		if err != nil {
			t.Fatalf("Get(%s) error: %v", u, err)
		}
		if entry == nil {
			t.Errorf("Get(%s) returned nil", u)
		}
	}
}

func TestMaxHTTPCacheBytes_100MB(t *testing.T) {
	const want = 100 * 1024 * 1024
	if MaxHTTPCacheBytes != want {
		t.Errorf("MaxHTTPCacheBytes = %d, want %d", MaxHTTPCacheBytes, want)
	}
}

func TestMaxHTTPCacheTTLSeconds_24h(t *testing.T) {
	const want = 86400
	if MaxHTTPCacheTTLSeconds != want {
		t.Errorf("MaxHTTPCacheTTLSeconds = %d, want %d", MaxHTTPCacheTTLSeconds, want)
	}
}
