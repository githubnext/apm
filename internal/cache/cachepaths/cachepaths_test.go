package cachepaths_test

import (
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/cachepaths"
)

func TestConstants(t *testing.T) {
	if cachepaths.GitDBBucket == "" {
		t.Error("GitDBBucket must not be empty")
	}
	if cachepaths.GitCheckoutsBucket == "" {
		t.Error("GitCheckoutsBucket must not be empty")
	}
	if cachepaths.HTTPBucket == "" {
		t.Error("HTTPBucket must not be empty")
	}
}

func TestGetCacheRoot_NoCache(t *testing.T) {
	// With noCache=true, should return a temp dir.
	dir, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir")
	}
	if !strings.HasPrefix(dir, os.TempDir()) && !strings.Contains(dir, "apm-cache-") {
		// Just verify it's a valid path
		if _, err2 := os.Stat(dir); err2 != nil {
			t.Errorf("temp dir does not exist: %v", err2)
		}
	}
}

func TestGetCacheRoot_NoCacheEnv(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "1")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir")
	}
}

func TestGetCacheRoot_OverrideEnv(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APM_CACHE_DIR", tmp)
	t.Setenv("APM_NO_CACHE", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != tmp {
		t.Errorf("expected %q, got %q", tmp, dir)
	}
}

func TestGetCacheRoot_NoCacheTrue_Singleton(t *testing.T) {
	// Calling GetCacheRoot(true) twice should return the same temp dir.
	d1, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}
	d2, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}
	if d1 != d2 {
		t.Errorf("expected same singleton dir, got %q and %q", d1, d2)
	}
}

func TestGetCacheRoot_NoCacheEnv_True(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "true")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir")
	}
}

func TestGetCacheRoot_NoCacheEnv_Yes(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "yes")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir")
	}
}

func TestConstantValues(t *testing.T) {
	if cachepaths.GitDBBucket != "git/db_v1" {
		t.Errorf("GitDBBucket = %q, want %q", cachepaths.GitDBBucket, "git/db_v1")
	}
	if cachepaths.GitCheckoutsBucket != "git/checkouts_v1" {
		t.Errorf("GitCheckoutsBucket = %q, want %q", cachepaths.GitCheckoutsBucket, "git/checkouts_v1")
	}
	if cachepaths.HTTPBucket != "http_v1" {
		t.Errorf("HTTPBucket = %q, want %q", cachepaths.HTTPBucket, "http_v1")
	}
}

func TestGetCacheRoot_XDGOverride(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APM_CACHE_DIR", "")
	t.Setenv("APM_NO_CACHE", "")
	t.Setenv("XDG_CACHE_HOME", tmp)
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir with XDG override")
	}
}
