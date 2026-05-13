// Package fileops provides retry-aware file operations for cross-platform reliability.
//
// On Windows, antivirus and endpoint-protection software briefly lock files
// while scanning them in temp directories. This package provides drop-in
// replacements for common file operations that transparently retry on
// transient lock errors with exponential backoff.
package fileops

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// RetryConfig controls retry behaviour.
type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// DefaultRetryConfig is tuned for AV scan locks (sub-second to ~3 s total wait).
var DefaultRetryConfig = RetryConfig{
	MaxRetries:    5,
	InitialDelay:  100 * time.Millisecond,
	MaxDelay:      2 * time.Second,
	BackoffFactor: 2.0,
}

// isTransientLockError returns true if err looks like a transient file-lock error.
func isTransientLockError(err error) bool {
	if err == nil {
		return false
	}
	if runtime.GOOS == "windows" {
		// Check for Windows sharing violation (winerror 32) or access denied (5)
		// We check the error message as a portable fallback
		msg := err.Error()
		if contains(msg, "The process cannot access the file") ||
			contains(msg, "Access is denied") {
			return true
		}
	}
	// Unix: EBUSY
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		// errno.EBUSY check via the underlying error
		if pathErr.Err.Error() == "device or resource busy" {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && indexStr(s, sub) >= 0)
}

func indexStr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func retryOnLock(cfg RetryConfig, op func() error) error {
	delay := cfg.InitialDelay
	var lastErr error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if err := op(); err == nil {
			return nil
		} else {
			lastErr = err
			if !isTransientLockError(err) || attempt == cfg.MaxRetries {
				return err
			}
		}
		time.Sleep(delay)
		delay = time.Duration(float64(delay) * cfg.BackoffFactor)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}
	return lastErr
}

// RobustRmtree removes path and everything under it, retrying on transient lock errors.
func RobustRmtree(path string) error {
	return retryOnLock(DefaultRetryConfig, func() error {
		err := os.RemoveAll(path)
		if err != nil {
			// Make read-only files writable and retry once
			_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, e error) error {
				if e == nil {
					_ = os.Chmod(p, 0o700)
				}
				return nil
			})
			return os.RemoveAll(path)
		}
		return nil
	})
}

// RobustCopytree copies the directory tree at src to dst, retrying on transient lock errors.
func RobustCopytree(src, dst string) error {
	return retryOnLock(DefaultRetryConfig, func() error {
		return copyDirTree(src, dst)
	})
}

func copyDirTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

// RobustCopy2 copies a single file from src to dst, retrying on transient lock errors.
func RobustCopy2(src, dst string) error {
	return retryOnLock(DefaultRetryConfig, func() error {
		return copyFile(src, dst)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
