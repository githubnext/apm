// Package cachepaths resolves the APM cache root and bucket paths.
package cachepaths

import (
"os"
"path/filepath"
"runtime"
"sync"
)

const (
GitDBBucket        = "git/db_v1"
GitCheckoutsBucket = "git/checkouts_v1"
HTTPBucket         = "http_v1"
)

var (
tempCacheMu  sync.Mutex
tempCacheDir string
)

// GetCacheRoot resolves the cache root directory.
// If noCache is true or APM_NO_CACHE env is set, returns a per-invocation temp dir.
func GetCacheRoot(noCache bool) (string, error) {
if noCache || isNoCacheEnv() {
return getTempCacheDir()
}
if override := os.Getenv("APM_CACHE_DIR"); override != "" {
abs, err := filepath.Abs(override)
if err != nil {
return "", err
}
return abs, os.MkdirAll(abs, 0o700)
}
dir := defaultCacheDir()
return dir, os.MkdirAll(dir, 0o700)
}

func isNoCacheEnv() bool {
v := os.Getenv("APM_NO_CACHE")
return v == "1" || v == "true" || v == "yes"
}

func getTempCacheDir() (string, error) {
tempCacheMu.Lock()
defer tempCacheMu.Unlock()
if tempCacheDir != "" {
return tempCacheDir, nil
}
dir, err := os.MkdirTemp("", "apm-cache-*")
if err != nil {
return "", err
}
tempCacheDir = dir
return dir, nil
}

func defaultCacheDir() string {
switch runtime.GOOS {
case "windows":
local := os.Getenv("LOCALAPPDATA")
if local == "" {
local = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
}
return filepath.Join(local, "apm", "Cache")
case "darwin":
if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
return filepath.Join(xdg, "apm")
}
home, _ := os.UserHomeDir()
return filepath.Join(home, "Library", "Caches", "apm")
default:
if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
return filepath.Join(xdg, "apm")
}
home, _ := os.UserHomeDir()
return filepath.Join(home, ".cache", "apm")
}
}
