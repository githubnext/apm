package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MaxHTTPCacheTTLSeconds caps server-provided TTL at 24 hours.
const MaxHTTPCacheTTLSeconds = 86400

// MaxHTTPCacheBytes caps total HTTP cache at 100 MB.
const MaxHTTPCacheBytes = 100 * 1024 * 1024

var maxAgeRE = regexp.MustCompile(`(?i)max-age=(\d+)`)

// CacheEntry represents a cached HTTP response.
type CacheEntry struct {
	Body        []byte
	ETag        string
	ExpiresAt   float64
	ContentType string
	StatusCode  int
}

// CacheStats holds cache statistics.
type CacheStats struct {
	EntryCount    int
	TotalSizeBytes int64
}

// shardMu provides in-process per-shard mutex to avoid concurrent writes to the same entry.
var (
	shardMuMap sync.Map // key: entry path string -> *sync.Mutex
)

func shardMutex(entryPath string) *sync.Mutex {
	v, _ := shardMuMap.LoadOrStore(entryPath, &sync.Mutex{})
	return v.(*sync.Mutex)
}

// HTTPCache is an HTTP response cache with conditional revalidation.
type HTTPCache struct {
	cacheDir string
}

// NewHTTPCache creates a new HTTPCache rooted at cacheRoot.
func NewHTTPCache(cacheRoot string) (*HTTPCache, error) {
	cacheDir := GetHTTPPath(cacheRoot)
	if err := ensureDir(cacheDir); err != nil {
		return nil, err
	}
	cleanupIncomplete(cacheDir)
	return &HTTPCache{cacheDir: cacheDir}, nil
}

// Get looks up a cached response for url.
// Returns nil if the entry is missing, expired, or fails integrity check.
func (c *HTTPCache) Get(rawURL string) *CacheEntry {
	entryPath := c.entryPath(rawURL)
	metaPath := filepath.Join(entryPath, "meta.json")
	bodyPath := filepath.Join(entryPath, "body")

	if !fileExists(metaPath) || !fileExists(bodyPath) {
		return nil
	}

	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return nil
	}
	var meta map[string]any
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return nil
	}

	expiresAt, _ := meta["expires_at"].(float64)
	if float64(time.Now().Unix()) > expiresAt {
		return nil
	}

	body, err := os.ReadFile(bodyPath)
	if err != nil {
		return nil
	}

	// Integrity check
	if recorded, ok := meta["body_sha256"].(string); ok && recorded != "" {
		actual := fmt.Sprintf("%x", sha256.Sum256(body))
		if actual != recorded {
			_ = os.RemoveAll(entryPath)
			return nil
		}
	}

	etag, _ := meta["etag"].(string)
	contentType, _ := meta["content_type"].(string)
	statusCode := 200
	if sc, ok := meta["status_code"].(float64); ok {
		statusCode = int(sc)
	}

	return &CacheEntry{
		Body:        body,
		ETag:        etag,
		ExpiresAt:   expiresAt,
		ContentType: contentType,
		StatusCode:  statusCode,
	}
}

// ConditionalHeaders returns If-None-Match headers for revalidation if an ETag is cached.
func (c *HTTPCache) ConditionalHeaders(rawURL string) map[string]string {
	metaPath := filepath.Join(c.entryPath(rawURL), "meta.json")
	if !fileExists(metaPath) {
		return map[string]string{}
	}
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return map[string]string{}
	}
	var meta map[string]any
	if err := json.Unmarshal(data, &meta); err != nil {
		return map[string]string{}
	}
	etag, _ := meta["etag"].(string)
	if etag != "" {
		return map[string]string{"If-None-Match": etag}
	}
	return map[string]string{}
}

// Store caches an HTTP response.
func (c *HTTPCache) Store(rawURL string, body []byte, statusCode int, headers map[string]string) {
	ttl := parseTTL(headers)
	etag := headerGet(headers, "ETag")
	contentType := headerGet(headers, "Content-Type")

	entryPath := c.entryPath(rawURL)
	if err := ensurePathWithin(entryPath, c.cacheDir); err != nil {
		return
	}

	now := float64(time.Now().Unix())
	meta := map[string]any{
		"url":          rawURL,
		"etag":         etag,
		"expires_at":   now + ttl,
		"content_type": contentType,
		"status_code":  statusCode,
		"stored_at":    now,
		"body_sha256":  fmt.Sprintf("%x", sha256.Sum256(body)),
	}

	// Atomic stage-rename
	staged := stagedPath(entryPath)
	if err := ensurePathWithin(staged, c.cacheDir); err != nil {
		return
	}
	if err := os.MkdirAll(staged, 0o700); err != nil {
		return
	}
	_ = os.Chmod(staged, 0o700)

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		_ = os.RemoveAll(staged)
		return
	}
	if err := os.WriteFile(filepath.Join(staged, "meta.json"), metaBytes, 0o600); err != nil {
		_ = os.RemoveAll(staged)
		return
	}
	if err := os.WriteFile(filepath.Join(staged, "body"), body, 0o600); err != nil {
		_ = os.RemoveAll(staged)
		return
	}

	mu := shardMutex(entryPath)
	mu.Lock()
	_ = os.RemoveAll(entryPath)
	if err := os.Rename(staged, entryPath); err != nil {
		_ = os.RemoveAll(staged)
	}
	mu.Unlock()

	_ = os.Chtimes(entryPath, time.Now(), time.Now())
	c.enforceSizeCap()
}

// RefreshExpiry refreshes TTL for a cached entry (called on 304 Not Modified).
func (c *HTTPCache) RefreshExpiry(rawURL string, headers map[string]string) {
	metaPath := filepath.Join(c.entryPath(rawURL), "meta.json")
	if !fileExists(metaPath) {
		return
	}
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return
	}
	var meta map[string]any
	if err := json.Unmarshal(data, &meta); err != nil {
		return
	}
	ttl := parseTTL(headers)
	meta["expires_at"] = float64(time.Now().Unix()) + ttl
	if newEtag := headerGet(headers, "ETag"); newEtag != "" {
		meta["etag"] = newEtag
	}
	updated, err := json.Marshal(meta)
	if err != nil {
		return
	}
	_ = os.WriteFile(metaPath, updated, 0o600)
	ep := c.entryPath(rawURL)
	_ = os.Chtimes(ep, time.Now(), time.Now())
}

// CleanAll removes all HTTP cache entries.
func (c *HTTPCache) CleanAll() {
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			_ = os.RemoveAll(filepath.Join(c.cacheDir, e.Name()))
		}
	}
}

// GetStats returns cache statistics.
func (c *HTTPCache) GetStats() CacheStats {
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return CacheStats{}
	}
	var stats CacheStats
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		stats.EntryCount++
		subDir := filepath.Join(c.cacheDir, e.Name())
		files, err := os.ReadDir(subDir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() {
				if fi, err := f.Info(); err == nil {
					stats.TotalSizeBytes += fi.Size()
				}
			}
		}
	}
	return stats
}

// entryPath derives the cache entry directory path for a URL.
func (c *HTTPCache) entryPath(rawURL string) string {
	h := sha256.Sum256([]byte(rawURL))
	urlHash := fmt.Sprintf("%x", h)[:16]
	entry := filepath.Join(c.cacheDir, urlHash)
	return entry
}

func parseTTL(headers map[string]string) float64 {
	cc := headerGet(headers, "Cache-Control")
	if m := maxAgeRE.FindStringSubmatch(cc); m != nil {
		n, err := strconv.Atoi(m[1])
		if err == nil {
			if n > MaxHTTPCacheTTLSeconds {
				return float64(MaxHTTPCacheTTLSeconds)
			}
			return float64(n)
		}
	}
	return 300.0
}

func headerGet(headers map[string]string, key string) string {
	lower := strings.ToLower(key)
	for k, v := range headers {
		if strings.ToLower(k) == lower {
			return v
		}
	}
	return ""
}

func ensurePathWithin(child, parent string) error {
	rel, err := filepath.Rel(parent, child)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path %s escapes cache root %s", child, parent)
	}
	return nil
}

func stagedPath(entryPath string) string {
	return entryPath + fmt.Sprintf(".incomplete.%d", time.Now().UnixNano())
}

func cleanupIncomplete(cacheDir string) {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".incomplete.") {
			_ = os.RemoveAll(filepath.Join(cacheDir, e.Name()))
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (c *HTTPCache) enforceSizeCap() {
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return
	}

	type entryInfo struct {
		mtime   time.Time
		path    string
		size    int64
	}
	var infos []entryInfo
	var totalSize int64

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dirPath := filepath.Join(c.cacheDir, e.Name())
		fi, err := os.Stat(dirPath)
		if err != nil {
			continue
		}
		var sz int64
		files, _ := os.ReadDir(dirPath)
		for _, f := range files {
			if !f.IsDir() {
				if ffi, err := f.Info(); err == nil {
					sz += ffi.Size()
				}
			}
		}
		infos = append(infos, entryInfo{mtime: fi.ModTime(), path: dirPath, size: sz})
		totalSize += sz
	}

	if totalSize <= MaxHTTPCacheBytes {
		return
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].mtime.Before(infos[j].mtime)
	})

	for _, info := range infos {
		if totalSize <= MaxHTTPCacheBytes {
			break
		}
		_ = os.RemoveAll(info.path)
		totalSize -= info.size
	}
}
