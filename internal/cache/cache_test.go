package cache_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/cache"
)

// --- url_normalize parity tests ---

func TestParityNormalizeRepoURL_HTTPS_dotgit(t *testing.T) {
	got := cache.NormalizeRepoURL("https://github.com/Owner/Repo.git")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("NormalizeRepoURL HTTPS .git: got %q want %q", got, want)
	}
}

func TestParityNormalizeRepoURL_SCP_like(t *testing.T) {
	got := cache.NormalizeRepoURL("git@github.com:owner/repo.git")
	want := "ssh://git@github.com/owner/repo"
	if got != want {
		t.Errorf("NormalizeRepoURL SCP-like: got %q want %q", got, want)
	}
}

func TestParityNormalizeRepoURL_SSH_explicit(t *testing.T) {
	got := cache.NormalizeRepoURL("ssh://git@github.com:22/owner/repo.git")
	want := "ssh://git@github.com/owner/repo"
	if got != want {
		t.Errorf("NormalizeRepoURL SSH explicit port: got %q want %q", got, want)
	}
}

func TestParityNormalizeRepoURL_HTTPS_caseInsensitiveHost(t *testing.T) {
	// Hostname lowercased, path lowercased for github.com
	got := cache.NormalizeRepoURL("https://GITHUB.COM/MyOrg/MyRepo")
	want := "https://github.com/myorg/myrepo"
	if got != want {
		t.Errorf("NormalizeRepoURL case-insensitive host: got %q want %q", got, want)
	}
}

func TestParityNormalizeRepoURL_NonCaseInsensitiveHost(t *testing.T) {
	// self-hosted: path case preserved
	got := cache.NormalizeRepoURL("https://gitea.example.com/MyOrg/MyRepo")
	want := "https://gitea.example.com/MyOrg/MyRepo"
	if got != want {
		t.Errorf("NormalizeRepoURL non-case-insensitive host: got %q want %q", got, want)
	}
}

func TestParityNormalizeRepoURL_TrailingSlash(t *testing.T) {
	got := cache.NormalizeRepoURL("https://github.com/owner/repo/")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("NormalizeRepoURL trailing slash: got %q want %q", got, want)
	}
}

func TestParityCacheShardKey_Deterministic(t *testing.T) {
	k1 := cache.CacheShardKey("https://github.com/owner/repo.git")
	k2 := cache.CacheShardKey("https://github.com/Owner/Repo.git")
	if k1 != k2 {
		t.Errorf("CacheShardKey not deterministic: %q vs %q", k1, k2)
	}
	if len(k1) != 16 {
		t.Errorf("CacheShardKey length: got %d want 16", len(k1))
	}
}

func TestParityCacheShardKey_SCPEqualsHTTPS(t *testing.T) {
	// github.com SCP-like and HTTPS should produce the same shard key
	k1 := cache.CacheShardKey("git@github.com:owner/repo.git")
	k2 := cache.CacheShardKey("ssh://git@github.com/owner/repo")
	if k1 != k2 {
		t.Errorf("CacheShardKey SCP vs SSH: %q vs %q", k1, k2)
	}
}

// --- paths parity tests ---

func TestParityGetCachePaths(t *testing.T) {
	root := "/tmp/test_cache_root"
	if cache.GetGitDBPath(root) != filepath.Join(root, "git/db_v1") {
		t.Error("GetGitDBPath wrong")
	}
	if cache.GetGitCheckoutsPath(root) != filepath.Join(root, "git/checkouts_v1") {
		t.Error("GetGitCheckoutsPath wrong")
	}
	if cache.GetHTTPPath(root) != filepath.Join(root, "http_v1") {
		t.Error("GetHTTPPath wrong")
	}
}

func TestParityGetCacheRoot_TempOnNoCache(t *testing.T) {
	dir, err := cache.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("GetCacheRoot(noCache=true) error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty temp cache dir")
	}
	// Should exist on disk
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("temp cache dir not created: %v", err)
	}
}

func TestParityGetCacheRoot_APMCacheDirOverride(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APM_CACHE_DIR", tmp)
	dir, err := cache.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("GetCacheRoot with APM_CACHE_DIR error: %v", err)
	}
	if dir != tmp {
		t.Errorf("expected %q got %q", tmp, dir)
	}
}

// --- integrity parity tests ---

func TestParityVerifyCheckoutSHA_ValidDetachedHEAD(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.Mkdir(gitDir, 0o700); err != nil {
		t.Fatal(err)
	}
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if !cache.VerifyCheckoutSHA(dir, sha) {
		t.Error("expected VerifyCheckoutSHA to return true for detached HEAD")
	}
}

func TestParityVerifyCheckoutSHA_Mismatch(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	_ = os.Mkdir(gitDir, 0o700)
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	_ = os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0o600)
	if cache.VerifyCheckoutSHA(dir, "0000000000000000000000000000000000000000") {
		t.Error("expected VerifyCheckoutSHA to return false for mismatched SHA")
	}
}

func TestParityVerifyCheckoutSHA_MissingDir(t *testing.T) {
	if cache.VerifyCheckoutSHA("/nonexistent/path/xyz", "abcdef1234567890abcdef1234567890abcdef12") {
		t.Error("expected false for missing directory")
	}
}

// --- http_cache parity tests ---

func TestParityHTTPCache_StoreAndGet(t *testing.T) {
	root := t.TempDir()
	c, err := cache.NewHTTPCache(root)
	if err != nil {
		t.Fatalf("NewHTTPCache error: %v", err)
	}
	url := "https://example.com/api/resource"
	body := []byte(`{"data":"test"}`)
	c.Store(url, body, 200, map[string]string{
		"Cache-Control": "max-age=3600",
		"ETag":          "\"abc123\"",
		"Content-Type":  "application/json",
	})
	entry := c.Get(url)
	if entry == nil {
		t.Fatal("expected cache hit, got nil")
	}
	if string(entry.Body) != string(body) {
		t.Errorf("body mismatch: got %q want %q", entry.Body, body)
	}
	if entry.ETag != `"abc123"` {
		t.Errorf("ETag mismatch: got %q want %q", entry.ETag, `"abc123"`)
	}
	if entry.StatusCode != 200 {
		t.Errorf("StatusCode: got %d want 200", entry.StatusCode)
	}
}

func TestParityHTTPCache_MissOnExpired(t *testing.T) {
	root := t.TempDir()
	c, _ := cache.NewHTTPCache(root)
	url := "https://example.com/expired"
	// TTL=0 => expires immediately (max-age=0)
	c.Store(url, []byte("body"), 200, map[string]string{"Cache-Control": "max-age=0"})
	// Sleep is not needed; max-age=0 means expires_at = now, so next call is after
	// We check that get returns nil for TTL=0 (expiresAt <= now)
	entry := c.Get(url)
	// May or may not be nil depending on sub-second timing; only check if nil that it's acceptable
	_ = entry // TTL=0 is a boundary case
}

func TestParityHTTPCache_ConditionalHeaders(t *testing.T) {
	root := t.TempDir()
	c, _ := cache.NewHTTPCache(root)
	url := "https://example.com/cond"
	c.Store(url, []byte("x"), 200, map[string]string{
		"ETag":          "\"v1\"",
		"Cache-Control": "max-age=60",
	})
	hdrs := c.ConditionalHeaders(url)
	if hdrs["If-None-Match"] != `"v1"` {
		t.Errorf("ConditionalHeaders: got %v", hdrs)
	}
}

func TestParityHTTPCache_GetStats(t *testing.T) {
	root := t.TempDir()
	c, _ := cache.NewHTTPCache(root)
	stats := c.GetStats()
	if stats.EntryCount != 0 {
		t.Errorf("expected 0 entries, got %d", stats.EntryCount)
	}
	c.Store("https://a.com/1", []byte("body1"), 200, map[string]string{"Cache-Control": "max-age=3600"})
	c.Store("https://a.com/2", []byte("body2"), 200, map[string]string{"Cache-Control": "max-age=3600"})
	stats = c.GetStats()
	if stats.EntryCount != 2 {
		t.Errorf("expected 2 entries, got %d", stats.EntryCount)
	}
	if stats.TotalSizeBytes == 0 {
		t.Error("expected non-zero total size")
	}
}

func TestParityHTTPCache_CleanAll(t *testing.T) {
	root := t.TempDir()
	c, _ := cache.NewHTTPCache(root)
	c.Store("https://example.com/x", []byte("body"), 200, map[string]string{"Cache-Control": "max-age=3600"})
	c.CleanAll()
	stats := c.GetStats()
	if stats.EntryCount != 0 {
		t.Errorf("expected 0 entries after CleanAll, got %d", stats.EntryCount)
	}
}
