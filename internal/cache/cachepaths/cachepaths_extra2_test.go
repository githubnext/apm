package cachepaths_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/cachepaths"
)

func TestConstants_GitDBBucketHasSlash(t *testing.T) {
	if !strings.Contains(cachepaths.GitDBBucket, "/") {
		t.Errorf("GitDBBucket should have a slash separator: %q", cachepaths.GitDBBucket)
	}
}

func TestConstants_GitCheckoutsBucketHasSlash(t *testing.T) {
	if !strings.Contains(cachepaths.GitCheckoutsBucket, "/") {
		t.Errorf("GitCheckoutsBucket should have a slash separator: %q", cachepaths.GitCheckoutsBucket)
	}
}

func TestConstants_HTTPBucketNoSlash(t *testing.T) {
	if strings.Contains(cachepaths.HTTPBucket, "/") {
		t.Errorf("HTTPBucket should not have a slash: %q", cachepaths.HTTPBucket)
	}
}

func TestConstants_GitBucketsSharePrefix(t *testing.T) {
	if !strings.HasPrefix(cachepaths.GitDBBucket, "git/") {
		t.Errorf("GitDBBucket should start with git/: %q", cachepaths.GitDBBucket)
	}
	if !strings.HasPrefix(cachepaths.GitCheckoutsBucket, "git/") {
		t.Errorf("GitCheckoutsBucket should start with git/: %q", cachepaths.GitCheckoutsBucket)
	}
}

func TestGetCacheRoot_RelativePath_Resolved(t *testing.T) {
	tmp := t.TempDir()
	sub := "relcache"
	// Use a relative path under tmp by constructing relative to cwd.
	// Set APM_CACHE_DIR to an absolute path but verify result is absolute.
	absPath := filepath.Join(tmp, sub)
	t.Setenv("APM_CACHE_DIR", absPath)
	t.Setenv("APM_NO_CACHE", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("GetCacheRoot should return absolute path, got %q", dir)
	}
}

func TestGetCacheRoot_NoCache_ReturnsExistingDir(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	dir, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fi, statErr := os.Stat(dir)
	if statErr != nil {
		t.Fatalf("expected dir to exist: %v", statErr)
	}
	if !fi.IsDir() {
		t.Errorf("expected directory, got %q", dir)
	}
}

func TestGetCacheRoot_APMCacheDir_DirCreated(t *testing.T) {
	tmp := t.TempDir()
	newDir := filepath.Join(tmp, "nested", "cache")
	t.Setenv("APM_CACHE_DIR", newDir)
	t.Setenv("APM_NO_CACHE", "")
	_, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error creating nested dir: %v", err)
	}
	if _, statErr := os.Stat(newDir); statErr != nil {
		t.Errorf("nested dir should be created: %v", statErr)
	}
}

func TestGetCacheRoot_NoCacheParam_IsolatedFromEnv(t *testing.T) {
	// Even if APM_CACHE_DIR is set, noCache=true should bypass it.
	t.Setenv("APM_CACHE_DIR", "/tmp/should-not-be-used")
	t.Setenv("APM_NO_CACHE", "")
	dir1, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Result should not be the APM_CACHE_DIR value.
	if dir1 == "/tmp/should-not-be-used" {
		t.Error("noCache=true should not return APM_CACHE_DIR")
	}
}

func TestGetCacheRoot_DefaultExists(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	t.Setenv("APM_CACHE_DIR", "")
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Fatal("expected non-empty dir")
	}
}

func TestConstants_AllDistinct(t *testing.T) {
	buckets := []string{
		cachepaths.GitDBBucket,
		cachepaths.GitCheckoutsBucket,
		cachepaths.HTTPBucket,
	}
	seen := map[string]bool{}
	for _, b := range buckets {
		if seen[b] {
			t.Errorf("duplicate bucket constant: %q", b)
		}
		seen[b] = true
	}
}

func TestGetCacheRoot_XDGCacheHomeDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APM_CACHE_DIR", "")
	t.Setenv("APM_NO_CACHE", "")
	t.Setenv("XDG_CACHE_HOME", tmp)
	dir, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(dir, tmp) {
		t.Errorf("expected dir under XDG_CACHE_HOME %q, got %q", tmp, dir)
	}
}
