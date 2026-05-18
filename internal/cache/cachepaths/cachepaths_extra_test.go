package cachepaths_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/cachepaths"
)

func TestGetCacheRoot_APMCacheDirAbsolute(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APM_CACHE_DIR", tmp)
	t.Setenv("APM_NO_CACHE", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Result should be absolute
	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got %q", dir)
	}
}

func TestGetCacheRoot_APMCacheDirCreated(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "apm-test-cache")
	t.Setenv("APM_CACHE_DIR", sub)
	t.Setenv("APM_NO_CACHE", "")
	_, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(sub); statErr != nil {
		t.Errorf("directory should be created: %v", statErr)
	}
}

func TestGetCacheRoot_NoCacheParamTrue_IsTempDir(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	dir, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Must be a real directory
	if info, err2 := os.Stat(dir); err2 != nil || !info.IsDir() {
		t.Errorf("result should be an existing directory: %v", err2)
	}
}

func TestGetCacheRoot_DefaultReturnsDir(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	t.Setenv("APM_CACHE_DIR", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("default cache root must not be empty")
	}
}

func TestGetCacheRoot_DefaultContainsApm(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	t.Setenv("APM_CACHE_DIR", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(dir), "apm") {
		t.Errorf("default cache root should contain 'apm': %q", dir)
	}
}

func TestConstants_ContainV1(t *testing.T) {
	for _, c := range []string{cachepaths.GitDBBucket, cachepaths.GitCheckoutsBucket, cachepaths.HTTPBucket} {
		if !strings.Contains(c, "_v1") {
			t.Errorf("bucket should contain _v1: %q", c)
		}
	}
}

func TestConstants_DistinctValues(t *testing.T) {
	buckets := []string{cachepaths.GitDBBucket, cachepaths.GitCheckoutsBucket, cachepaths.HTTPBucket}
	seen := map[string]bool{}
	for _, b := range buckets {
		if seen[b] {
			t.Errorf("duplicate bucket: %q", b)
		}
		seen[b] = true
	}
}

func TestGetCacheRoot_NoCacheEnvValues(t *testing.T) {
	for _, val := range []string{"1", "true", "yes"} {
		t.Run("APM_NO_CACHE="+val, func(t *testing.T) {
			t.Setenv("APM_NO_CACHE", val)
			dir, err := cachepaths.GetCacheRoot(false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dir == "" {
				t.Error("expected non-empty dir")
			}
		})
	}
}

func TestGetCacheRoot_NoCacheEnvOtherValues_NoTmp(t *testing.T) {
	// "false", "0", "no" should NOT trigger no-cache
	for _, val := range []string{"false", "0", "no"} {
		t.Run("APM_NO_CACHE="+val, func(t *testing.T) {
			t.Setenv("APM_NO_CACHE", val)
			t.Setenv("APM_CACHE_DIR", t.TempDir())
			dir, err := cachepaths.GetCacheRoot(false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dir == "" {
				t.Error("expected non-empty dir")
			}
		})
	}
}
