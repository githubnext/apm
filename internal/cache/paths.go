// Package cache provides HTTP and git caching primitives for the APM CLI.
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Bucket layout within cache root.
const (
	GitDBBucket        = "git/db_v1"
	GitCheckoutsBucket = "git/checkouts_v1"
	HTTPBucket         = "http_v1"
)

var (
	tempCacheDir  string
	tempCacheMu   sync.Mutex
	tempCacheOnce sync.Once
)

// GetCacheRoot resolves the cache root directory.
// If noCache is true, returns a temporary directory cleaned up at process exit.
// Honours APM_NO_CACHE and APM_CACHE_DIR environment variables.
func GetCacheRoot(noCache bool) (string, error) {
	if noCache || strings.TrimSpace(os.Getenv("APM_NO_CACHE")) == "1" ||
		strings.TrimSpace(os.Getenv("APM_NO_CACHE")) == "true" ||
		strings.TrimSpace(os.Getenv("APM_NO_CACHE")) == "yes" {
		return getTempCacheRoot()
	}

	override := strings.TrimSpace(os.Getenv("APM_CACHE_DIR"))
	if override != "" {
		return validateAndEnsure(override)
	}

	return validateAndEnsure(platformDefault())
}

// GetGitDBPath returns the git database bucket path (full clones).
func GetGitDBPath(cacheRoot string) string {
	return filepath.Join(cacheRoot, GitDBBucket)
}

// GetGitCheckoutsPath returns the git checkouts bucket path (per-SHA working copies).
func GetGitCheckoutsPath(cacheRoot string) string {
	return filepath.Join(cacheRoot, GitCheckoutsBucket)
}

// GetHTTPPath returns the HTTP cache bucket path.
func GetHTTPPath(cacheRoot string) string {
	return filepath.Join(cacheRoot, HTTPBucket)
}

func platformDefault() string {
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return filepath.Join(localAppData, "apm", "Cache")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Local", "apm", "Cache")
	case "darwin":
		xdg := strings.TrimSpace(os.Getenv("XDG_CACHE_HOME"))
		if xdg != "" {
			return filepath.Join(xdg, "apm")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Caches", "apm")
	default:
		xdg := strings.TrimSpace(os.Getenv("XDG_CACHE_HOME"))
		if xdg != "" {
			return filepath.Join(xdg, "apm")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".cache", "apm")
	}
}

func validateAndEnsure(pathStr string) (string, error) {
	if pathStr == "" {
		return "", fmt.Errorf("cache path must not be empty")
	}
	if strings.Contains(pathStr, "\x00") {
		return "", fmt.Errorf("cache path must not contain NUL bytes")
	}

	expanded := pathStr
	if strings.HasPrefix(expanded, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot expand ~: %w", err)
		}
		expanded = home + expanded[1:]
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", fmt.Errorf("cannot make cache path absolute: %w", err)
	}
	if err := ensureDir(abs); err != nil {
		return "", err
	}
	return abs, nil
}

func ensureDir(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("failed to create cache directory %s: %w", path, err)
	}
	// Best-effort chmod -- no-op on Windows
	_ = os.Chmod(path, 0o700)
	return nil
}

func getTempCacheRoot() (string, error) {
	tempCacheMu.Lock()
	defer tempCacheMu.Unlock()
	if tempCacheDir == "" {
		dir, err := os.MkdirTemp("", "apm_cache_")
		if err != nil {
			return "", fmt.Errorf("failed to create temp cache: %w", err)
		}
		_ = os.Chmod(dir, 0o700)
		tempCacheDir = dir
	}
	return tempCacheDir, nil
}
