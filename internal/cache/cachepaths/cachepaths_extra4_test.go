package cachepaths_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/cachepaths"
)

func TestGetCacheRoot_APMCacheEnvSet_UsesIt_Extra4(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APM_CACHE_DIR", dir)
	got, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(got, dir) {
		t.Errorf("expected path under %q, got %q", dir, got)
	}
}

func TestGetCacheRoot_NoCache_NotEmpty_Extra4(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	got, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetCacheRoot_ReturnsDir_Extra4(t *testing.T) {
	t.Setenv("APM_NO_CACHE", "")
	got, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected a non-empty path")
	}
}

func TestGetCacheRoot_APMCacheDir_IsAbsolute_Extra4(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APM_CACHE_DIR", dir)
	got, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
}

func TestGetCacheRoot_NoCacheTrue_DoesNotCreate_Extra4(t *testing.T) {
	t.Setenv("APM_CACHE_DIR", "")
	got, err := cachepaths.GetCacheRoot(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty temp path")
	}
}

func TestGetCacheRoot_CustomDir_Created_Extra4(t *testing.T) {
	parent := t.TempDir()
	custom := filepath.Join(parent, "myapmcache")
	t.Setenv("APM_CACHE_DIR", custom)
	got, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(got); statErr != nil {
		t.Errorf("expected created dir at %q: %v", got, statErr)
	}
}

func TestGetCacheRoot_CallTwice_SameResult_Extra4(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APM_CACHE_DIR", dir)
	a, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if a != b {
		t.Errorf("expected same result on repeated calls: %q vs %q", a, b)
	}
}
