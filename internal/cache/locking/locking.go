// Package locking provides cross-platform shard locking and atomic landing
// primitives for the APM cache layer.
//
// Atomic landing protocol:
//  1. Stage content into <shard>.incomplete.<pid>.<ts>/
//  2. Acquire shard .lock file
//  3. Re-check final path does not exist (TOCTOU defense)
//  4. os.Rename staged -> final (atomic on same filesystem)
//  5. Release lock
//  6. On cache init, clean up stale *.incomplete.* siblings
package locking

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const defaultLockTimeout = 120 * time.Second

// ShardLock is a per-shard file-based advisory lock implemented with
// a sync.Mutex for in-process concurrency and a sentinel file for
// cross-process coordination (best-effort on platforms without flock).
type ShardLock struct {
	mu       sync.Mutex
	lockFile string
	timeout  time.Duration
}

// NewShardLock creates a ShardLock for the given shard directory.
// The lock file is placed adjacent to (not inside) the shard directory.
func NewShardLock(shardDir string, timeout time.Duration) *ShardLock {
	if timeout == 0 {
		timeout = defaultLockTimeout
	}
	ext := filepath.Ext(shardDir)
	base := strings.TrimSuffix(shardDir, ext)
	return &ShardLock{
		lockFile: base + ".lock",
		timeout:  timeout,
	}
}

// Lock acquires the shard lock. Returns an error on timeout.
func (l *ShardLock) Lock() error {
	deadline := time.Now().Add(l.timeout)
	for {
		l.mu.Lock()
		f, err := os.OpenFile(l.lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			f.Close()
			return nil
		}
		l.mu.Unlock()
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for shard lock: %s", l.lockFile)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Unlock releases the shard lock.
func (l *ShardLock) Unlock() {
	os.Remove(l.lockFile)
	l.mu.Unlock()
}

// StagePath returns a staging path adjacent to finalPath.
// Format: <final>.incomplete.<pid>.<monotonic_ns>
func StagePath(finalPath string) string {
	pid := os.Getpid()
	ts := time.Now().UnixNano()
	base := filepath.Base(finalPath)
	dir := filepath.Dir(finalPath)
	return filepath.Join(dir, fmt.Sprintf("%s.incomplete.%d.%d", base, pid, ts))
}

// AtomicLand atomically moves staged to final under lock.
// Returns true if the landing succeeded, false if another process
// already populated final (TOCTOU defense).
func AtomicLand(staged, final string, lock *ShardLock) (bool, error) {
	if err := lock.Lock(); err != nil {
		SafeRemoveAll(staged)
		return false, err
	}
	defer lock.Unlock()

	if _, err := os.Stat(final); err == nil {
		// Another process already populated the target.
		SafeRemoveAll(staged)
		return false, nil
	}

	if err := os.Rename(staged, final); err != nil {
		SafeRemoveAll(staged)
		return false, fmt.Errorf("atomic rename %s -> %s: %w", staged, final, err)
	}
	return true, nil
}

// CleanupIncomplete removes stale .incomplete.* directories under parent.
// Returns the number of directories removed.
func CleanupIncomplete(parent string) int {
	info, err := os.Stat(parent)
	if err != nil || !info.IsDir() {
		return 0
	}

	entries, err := os.ReadDir(parent)
	if err != nil {
		return 0
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), ".incomplete.") {
			if err := SafeRemoveAll(filepath.Join(parent, entry.Name())); err == nil {
				removed++
			}
		}
	}
	return removed
}

// SafeRemoveAll removes path without following symlinks (best-effort).
func SafeRemoveAll(path string) error {
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(path)
}
