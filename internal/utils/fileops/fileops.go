// Package fileops provides retry-aware file operations for cross-platform
// reliability.
//
// On Windows, antivirus and endpoint-protection software briefly lock files
// while scanning them in temp directories. This package provides drop-in
// replacements for os.RemoveAll, filepath.WalkDir-based copy, and os.Copy
// that transparently retry on transient lock errors with exponential backoff.
package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultMaxRetries    = 5
	defaultInitialDelay  = 100 * time.Millisecond
	defaultMaxDelay      = 2 * time.Second
	defaultBackoffFactor = 2.0
)

// isTransientLockError returns true when err looks like a transient file-lock
// error. Platform-specific detection is in lock_unix.go / lock_windows.go.
// The function defined here handles the Unix EBUSY case; the build-tag files
// add Windows winerror 32/5 detection.

// retryOnLock executes op, retrying on transient lock errors.
func retryOnLock(op func() error, desc string, maxRetries int, initial, max time.Duration, backoff float64, beforeRetry func()) error {
	delay := initial
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := op()
		if err == nil {
			return nil
		}
		lastErr = err
		if !isTransientLockError(err) || attempt == maxRetries {
			return err
		}
		debugFileOp(fmt.Sprintf("%s: transient lock (attempt %d/%d), retrying in %s -- %v",
			desc, attempt+1, maxRetries, delay, err))
		if beforeRetry != nil {
			beforeRetry()
		}
		time.Sleep(delay)
		next := time.Duration(float64(delay) * backoff)
		if next > max {
			next = max
		}
		delay = next
	}
	return lastErr
}

// debugFileOp prints debug output when APM_DEBUG is set.
func debugFileOp(msg string) {
	if os.Getenv("APM_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", msg)
	}
}

// RobustRemoveAll removes a directory tree, retrying on transient lock errors.
// If ignoreErrors is true, any error after retries is silently discarded.
func RobustRemoveAll(path string, ignoreErrors bool, maxRetries int) error {
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	err := retryOnLock(func() error {
		return removeAllWritable(path)
	}, "rmtree "+path, maxRetries, defaultInitialDelay, defaultMaxDelay, defaultBackoffFactor, nil)
	if err != nil && ignoreErrors {
		return nil
	}
	return err
}

// removeAllWritable removes path, chmod-ing read-only files writable first.
func removeAllWritable(path string) error {
	// chmod all files writable so rmtree succeeds on read-only trees (e.g. git pack).
	_ = filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		_ = os.Chmod(p, 0o666)
		return nil
	})
	return os.RemoveAll(path)
}

// RobustCopyTree copies a directory tree from src to dst, retrying on
// transient lock errors. Any partial dst is removed before each retry
// unless dirsExistOK is true.
func RobustCopyTree(src, dst string, symlinks, dirsExistOK bool, maxRetries int) error {
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	var beforeRetry func()
	if !dirsExistOK {
		beforeRetry = func() {
			_ = os.RemoveAll(dst)
		}
	}
	return retryOnLock(func() error {
		return copyTree(src, dst, symlinks, dirsExistOK)
	}, fmt.Sprintf("copytree %s -> %s", src, dst), maxRetries, defaultInitialDelay, defaultMaxDelay, defaultBackoffFactor, beforeRetry)
}

// copyTree is the inner copy-tree implementation (no retry).
func copyTree(src, dst string, symlinks, dirsExistOK bool) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if mkErr := os.MkdirAll(target, 0o755); mkErr != nil && !dirsExistOK {
				return mkErr
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			if symlinks {
				link, readErr := os.Readlink(path)
				if readErr != nil {
					return readErr
				}
				return os.Symlink(link, target)
			}
			// Dereference symlink: stat the real file.
			info, statErr := os.Stat(path)
			if statErr != nil || !info.Mode().IsRegular() {
				return nil
			}
		}
		return copyFile(path, target)
	})
}

// RobustCopy2 copies a single file with metadata, retrying on transient lock
// errors.
func RobustCopy2(src, dst string, maxRetries int) error {
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	return retryOnLock(func() error {
		return copyFile(src, dst)
	}, fmt.Sprintf("copy2 %s -> %s", src, dst), maxRetries, defaultInitialDelay, defaultMaxDelay, defaultBackoffFactor, nil)
}

// copyFile copies src to dst, preserving permissions.
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
