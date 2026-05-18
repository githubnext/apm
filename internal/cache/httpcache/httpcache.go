// Package httpcache implements an HTTP response cache with conditional revalidation.
//
// Caches HTTP GET responses using content-addressable storage with support for:
//   - Cache-Control: max-age=N (capped at 24h)
//   - ETag / If-None-Match conditional revalidation
//   - LRU eviction when cache exceeds size limit
//   - Atomic writes (stage-rename pattern)
//   - sha256 body integrity verification on read
package httpcache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/githubnext/apm/internal/cache/cachepaths"
	"github.com/githubnext/apm/internal/cache/locking"
)

// getHTTPPath returns the HTTP cache directory.
func getHTTPPath(cacheRoot string) string {
	return filepath.Join(cacheRoot, cachepaths.HTTPBucket)
}

const (
	// MaxHTTPCacheTTLSeconds caps TTL at 24 hours even if the server says longer.
	MaxHTTPCacheTTLSeconds = 86400
	// MaxHTTPCacheBytes caps total cache size at 100 MB.
	MaxHTTPCacheBytes = 100 * 1024 * 1024
)

var maxAgeRe = regexp.MustCompile(`(?i)max-age=(\d+)`)

// CacheEntry holds an HTTP response retrieved from cache.
type CacheEntry struct {
	Body        []byte
	ETag        string
	ExpiresAt   float64
	ContentType string
	StatusCode  int
}

// GetStats holds aggregate statistics about the HTTP cache.
type GetStats struct {
	EntryCount    int
	TotalSizeBytes int64
}

type entryMeta struct {
	URL         string  `json:"url"`
	ETag        string  `json:"etag"`
	ExpiresAt   float64 `json:"expires_at"`
	ContentType string  `json:"content_type"`
	StatusCode  int     `json:"status_code"`
	StoredAt    float64 `json:"stored_at"`
	BodySHA256  string  `json:"body_sha256"`
}

// HttpCache is a persistent HTTP response cache.
type HttpCache struct {
	cacheDir string
}

// New creates an HttpCache rooted at cacheRoot.
func New(cacheRoot string) (*HttpCache, error) {
	cacheDir := getHTTPPath(cacheRoot)
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		return nil, fmt.Errorf("httpcache: mkdir: %w", err)
	}
	locking.CleanupIncomplete(cacheDir)
	return &HttpCache{cacheDir: cacheDir}, nil
}

// Get returns a cached response for url, or nil if not cached or expired.
// Returns a non-empty ETag if the entry has one (for conditional revalidation).
func (c *HttpCache) Get(url string) (*CacheEntry, error) {
	entryPath := c.entryPath(url)
	metaPath := filepath.Join(entryPath, "meta.json")
	bodyPath := filepath.Join(entryPath, "body")

	raw, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, nil //nolint:nilerr // cache miss
	}
	var meta entryMeta
	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil, nil //nolint:nilerr // corrupt entry
	}

	if float64(time.Now().Unix()) > meta.ExpiresAt {
		// Stale but return with ETag so caller can revalidate.
		if meta.ETag == "" {
			return nil, nil
		}
		return &CacheEntry{
			ETag:        meta.ETag,
			ExpiresAt:   meta.ExpiresAt,
			ContentType: meta.ContentType,
			StatusCode:  meta.StatusCode,
		}, nil
	}

	body, err := os.ReadFile(bodyPath)
	if err != nil {
		return nil, nil //nolint:nilerr
	}

	// Integrity check
	sum := sha256.Sum256(body)
	if hex.EncodeToString(sum[:]) != meta.BodySHA256 {
		_ = os.RemoveAll(entryPath)
		return nil, nil
	}

	// Bump mtime for LRU
	_ = os.Chtimes(entryPath, time.Now(), time.Now())

	return &CacheEntry{
		Body:        body,
		ETag:        meta.ETag,
		ExpiresAt:   meta.ExpiresAt,
		ContentType: meta.ContentType,
		StatusCode:  meta.StatusCode,
	}, nil
}

// Store caches an HTTP response for url.
func (c *HttpCache) Store(url string, body []byte, statusCode int, headers map[string]string) {
	ttl := c.parseTTL(headers)
	etag := headers["ETag"]
	if etag == "" {
		etag = headers["etag"]
	}
	ct := headers["Content-Type"]
	if ct == "" {
		ct = headers["content-type"]
	}

	sum := sha256.Sum256(body)
	meta := entryMeta{
		URL:         url,
		ETag:        etag,
		ExpiresAt:   float64(time.Now().Unix()) + ttl,
		ContentType: ct,
		StatusCode:  statusCode,
		StoredAt:    float64(time.Now().Unix()),
		BodySHA256:  hex.EncodeToString(sum[:]),
	}

	entryPath := c.entryPath(url)
	staged := locking.StagePath(entryPath)

	if err := os.MkdirAll(staged, 0o700); err != nil {
		return
	}

	metaBytes, _ := json.Marshal(meta)
	if err := os.WriteFile(filepath.Join(staged, "meta.json"), metaBytes, 0o600); err != nil {
		_ = os.RemoveAll(staged)
		return
	}
	if err := os.WriteFile(filepath.Join(staged, "body"), body, 0o600); err != nil {
		_ = os.RemoveAll(staged)
		return
	}

	lock := locking.NewShardLock(entryPath, 0)
	if entryPath != "" {
		_ = os.RemoveAll(entryPath)
	}
	_, _ = locking.AtomicLand(staged, entryPath, lock)
	_ = os.Chtimes(entryPath, time.Now(), time.Now())

	c.enforceSizeCap()
}

// RefreshExpiry updates the TTL for a cached entry on 304 Not Modified.
func (c *HttpCache) RefreshExpiry(url string, headers map[string]string) {
	entryPath := c.entryPath(url)
	metaPath := filepath.Join(entryPath, "meta.json")
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		return
	}
	var meta entryMeta
	if err := json.Unmarshal(raw, &meta); err != nil {
		return
	}
	ttl := c.parseTTL(headers)
	meta.ExpiresAt = float64(time.Now().Unix()) + ttl
	if newEtag := headers["ETag"]; newEtag != "" {
		meta.ETag = newEtag
	}
	metaBytes, _ := json.Marshal(meta)
	_ = os.WriteFile(metaPath, metaBytes, 0o600)
}

// GetStats returns aggregate cache statistics.
func (c *HttpCache) GetStats() GetStats {
	var stats GetStats
	entries, _ := os.ReadDir(c.cacheDir)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		stats.EntryCount++
		sub := filepath.Join(c.cacheDir, e.Name())
		files, _ := os.ReadDir(sub)
		for _, f := range files {
			if info, err := f.Info(); err == nil {
				stats.TotalSizeBytes += info.Size()
			}
		}
	}
	return stats
}

// CleanAll removes all HTTP cache content.
func (c *HttpCache) CleanAll() {
	_ = os.RemoveAll(c.cacheDir)
	_ = os.MkdirAll(c.cacheDir, 0o700)
}

func (c *HttpCache) entryPath(url string) string {
	sum := sha256.Sum256([]byte(url))
	key := hex.EncodeToString(sum[:])
	return filepath.Join(c.cacheDir, key[:2], key)
}

func (c *HttpCache) parseTTL(headers map[string]string) float64 {
	cc := headers["Cache-Control"]
	if cc == "" {
		cc = headers["cache-control"]
	}
	if m := maxAgeRe.FindStringSubmatch(cc); len(m) == 2 {
		if n, err := strconv.ParseFloat(m[1], 64); err == nil {
			if n > MaxHTTPCacheTTLSeconds {
				return MaxHTTPCacheTTLSeconds
			}
			return n
		}
	}
	return 0
}

func (c *HttpCache) enforceSizeCap() {
	stats := c.GetStats()
	if stats.TotalSizeBytes <= MaxHTTPCacheBytes {
		return
	}

	// Collect entries sorted by mtime (LRU eviction)
	type entry struct {
		path  string
		mtime time.Time
		size  int64
	}
	var entries []entry
	_ = filepath.WalkDir(c.cacheDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == c.cacheDir {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if strings.Count(path[len(c.cacheDir):], string(os.PathSeparator)) == 2 {
			size := int64(0)
			files, _ := os.ReadDir(path)
			for _, f := range files {
				if fi, err := f.Info(); err == nil {
					size += fi.Size()
				}
			}
			entries = append(entries, entry{path: path, mtime: info.ModTime(), size: size})
		}
		return nil
	})
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].mtime.Before(entries[j].mtime)
	})

	total := stats.TotalSizeBytes
	for _, e := range entries {
		if total <= MaxHTTPCacheBytes {
			break
		}
		_ = os.RemoveAll(e.path)
		total -= e.size
	}
}
