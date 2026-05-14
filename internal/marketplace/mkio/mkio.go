// Package mkio provides shared I/O helpers for marketplace modules.
// Migrated from src/apm_cli/marketplace/_io.py
package mkio

import (
	"os"
	"path/filepath"
)

// AtomicWrite writes content to path atomically via tmp + rename.
// The caller sees either the complete new content or the previous
// content -- never a partial write.
func AtomicWrite(path string, content []byte) error {
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	tmpPath := path[:len(path)-len(ext)] + ext + ".tmp"

	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	_, writeErr := f.Write(content)
	syncErr := f.Sync()
	closeErr := f.Close()

	if writeErr != nil {
		os.Remove(tmpPath)
		return writeErr
	}
	if syncErr != nil {
		os.Remove(tmpPath)
		return syncErr
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return closeErr
	}

	_ = dir // dir used implicitly via tmpPath construction
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}

// AtomicWriteString writes string content to path atomically.
func AtomicWriteString(path, content string) error {
	return AtomicWrite(path, []byte(content))
}
