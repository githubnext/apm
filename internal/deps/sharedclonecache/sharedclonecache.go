// Package sharedclonecache implements a per-run shared clone cache for
// subdirectory dependency deduplication.
// Ported from src/apm_cli/deps/shared_clone_cache.py
package sharedclonecache

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// CloneFn is called to perform the initial clone into the given directory.
type CloneFn func(clonePath string) error

// FetchFn is an optional callback that tries to fetch a missing SHA into an
// already-cloned bare for the same repo. Returns false to signal failure.
type FetchFn func(barePath string, sha string) bool

type cacheEntry struct {
	mu    sync.Mutex
	path  string
	err   error
}

// SharedCloneCache is a thread-safe per-run cache of shared Git clones.
// Keys are (host, owner, repo, ref) tuples. The first caller for a given
// key performs the clone; concurrent callers block until the clone completes
// and then reuse the result.
type SharedCloneCache struct {
	baseDir string

	mu              sync.Mutex
	entries         map[[4]string]*cacheEntry
	tempDirs        []string
	repoBares       map[[3]string][]repoBareEntry
	bareFetchLocks  map[string]*sync.Mutex
}

type repoBareEntry struct {
	ref  string
	path string
}

// New creates a new SharedCloneCache.
// If baseDir is empty, the system temp directory is used.
func New(baseDir string) *SharedCloneCache {
	return &SharedCloneCache{
		baseDir:        baseDir,
		entries:        make(map[[4]string]*cacheEntry),
		repoBares:      make(map[[3]string][]repoBareEntry),
		bareFetchLocks: make(map[string]*sync.Mutex),
	}
}

// GetOrClone returns a path to a shared clone, cloning on first access.
// clone_fn is called at most once per unique (host, owner, repo, ref) key.
func (c *SharedCloneCache) GetOrClone(
	host, owner, repo, ref string,
	cloneFn CloneFn,
	fetchFn FetchFn,
) (string, error) {
	key := [4]string{host, owner, repo, ref}
	entry := c.getOrCreateEntry(key)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.path != "" {
		return entry.path, nil
	}
	if entry.err != nil {
		entry.err = nil
	}

	// Tier-0: try fetching the SHA into an existing bare for the same repo.
	if ref != "" && fetchFn != nil {
		if existingBare := c.findRepoBare(host, owner, repo); existingBare != "" {
			bareLock := c.getBareFetchLock(existingBare)
			bareLock.Lock()
			ok := fetchFn(existingBare, ref)
			bareLock.Unlock()
			if ok {
				entry.path = existingBare
				c.mu.Lock()
				repoKey := [3]string{host, owner, repo}
				c.repoBares[repoKey] = append(c.repoBares[repoKey], repoBareEntry{ref: ref, path: existingBare})
				c.mu.Unlock()
				return existingBare, nil
			}
		}
	}

	// First caller: perform the clone.
	prefix := fmt.Sprintf("apm_shared_%s_%s_", owner, repo)
	var tempDir string
	var err error
	if c.baseDir != "" {
		tempDir, err = os.MkdirTemp(c.baseDir, prefix)
	} else {
		tempDir, err = os.MkdirTemp("", prefix)
	}
	if err != nil {
		entry.err = err
		return "", err
	}
	c.mu.Lock()
	c.tempDirs = append(c.tempDirs, tempDir)
	c.mu.Unlock()

	clonePath := filepath.Join(tempDir, "bare")
	if err := cloneFn(clonePath); err != nil {
		entry.err = err
		return "", err
	}

	// Debug-mode shape invariant: clone_fn MUST produce a bare repo.
	if os.Getenv("APM_DEBUG") != "" {
		headFile := filepath.Join(clonePath, "HEAD")
		gitDir := filepath.Join(clonePath, ".git")
		headInfo, headErr := os.Stat(headFile)
		_, gitDirErr := os.Stat(gitDir)
		headPresent := headErr == nil && !headInfo.IsDir()
		gitDirPresent := gitDirErr == nil
		if !headPresent || gitDirPresent {
			err := fmt.Errorf(
				"SharedCloneCache invariant violated: %s is not a bare repo "+
					"(HEAD file present: %v, .git/ present: %v)",
				clonePath, headPresent, gitDirPresent,
			)
			entry.err = err
			return "", err
		}
	}

	entry.path = clonePath
	c.mu.Lock()
	repoKey := [3]string{host, owner, repo}
	c.repoBares[repoKey] = append(c.repoBares[repoKey], repoBareEntry{ref: ref, path: clonePath})
	c.mu.Unlock()
	return clonePath, nil
}

// findRepoBare returns an existing bare path for the same repo (any ref), or "".
func (c *SharedCloneCache) findRepoBare(host, owner, repo string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	entries := c.repoBares[[3]string{host, owner, repo}]
	if len(entries) > 0 {
		return entries[0].path
	}
	return ""
}

// getOrCreateEntry retrieves or creates a cache entry (thread-safe).
func (c *SharedCloneCache) getOrCreateEntry(key [4]string) *cacheEntry {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.entries[key]; ok {
		return e
	}
	e := &cacheEntry{}
	c.entries[key] = e
	return e
}

// getBareFetchLock returns the per-bare-path lock.
func (c *SharedCloneCache) getBareFetchLock(barePath string) *sync.Mutex {
	c.mu.Lock()
	defer c.mu.Unlock()
	if l, ok := c.bareFetchLocks[barePath]; ok {
		return l
	}
	l := &sync.Mutex{}
	c.bareFetchLocks[barePath] = l
	return l
}

// Cleanup removes all temporary clone directories.
func (c *SharedCloneCache) Cleanup() {
	c.mu.Lock()
	dirs := make([]string, len(c.tempDirs))
	copy(dirs, c.tempDirs)
	c.tempDirs = nil
	c.entries = make(map[[4]string]*cacheEntry)
	c.repoBares = make(map[[3]string][]repoBareEntry)
	c.bareFetchLocks = make(map[string]*sync.Mutex)
	c.mu.Unlock()
	for _, d := range dirs {
		if err := os.RemoveAll(d); err != nil {
			log.Printf("Failed to clean shared clone dir: %s: %v", d, err)
		}
	}
}
