package cachepaths_test

import (
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/cachepaths"
)

func TestConstants_GitDBBucketNonEmpty_Extra3(t *testing.T) {
	if cachepaths.GitDBBucket == "" {
		t.Error("GitDBBucket must not be empty")
	}
}

func TestConstants_HTTPBucketNonEmpty_Extra3(t *testing.T) {
	if cachepaths.HTTPBucket == "" {
		t.Error("HTTPBucket must not be empty")
	}
}

func TestConstants_GitCheckoutsBucketNonEmpty_Extra3(t *testing.T) {
	if cachepaths.GitCheckoutsBucket == "" {
		t.Error("GitCheckoutsBucket must not be empty")
	}
}

func TestConstants_BucketsAreValidDirNames_Extra3(t *testing.T) {
	for _, b := range []string{cachepaths.GitDBBucket, cachepaths.GitCheckoutsBucket, cachepaths.HTTPBucket} {
		if strings.Contains(b, " ") {
			t.Errorf("bucket %q must not contain spaces", b)
		}
		if strings.HasPrefix(b, "/") {
			t.Errorf("bucket %q must be relative", b)
		}
	}
}

func TestGetCacheRoot_APMNoCacheEnvYes_ReturnsDir_Extra3(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "yes")
	t.Setenv("APM_CACHE_DIR", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty temp dir")
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Errorf("temp dir should exist: %v", statErr)
	}
}

func TestGetCacheRoot_APMNoCacheEnv1_ReturnsDir_Extra3(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "1")
	t.Setenv("APM_CACHE_DIR", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty temp dir")
	}
}

func TestGetCacheRoot_NoCacheParamTrue_IgnoresAPMCacheDir_Extra3(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	t.Setenv("APM_CACHE_DIR", "/some/explicit/path")
	dir, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "/some/explicit/path" {
		t.Error("noCache=true should bypass APM_CACHE_DIR")
	}
}

func TestGetCacheRoot_APMCacheDirTakesPrecedenceOverDefault_Extra3(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APM_CACHE_DIR", tmp)
	t.Setenv("APM_NO_CACHE", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != tmp {
		t.Errorf("expected APM_CACHE_DIR=%q, got %q", tmp, dir)
	}
}
