// Package atomicio provides atomic file-write primitives.
// Mirrors src/apm_cli/utils/atomic_io.py.
package atomicio

import (
	"os"
	"path/filepath"
)

// WriteText atomically writes data (UTF-8) to path.
// The temp file is created in path's parent directory so the eventual
// os.Rename is a same-filesystem rename. If newFileMode > 0 and path
// does not yet exist, the temp file's mode bits are set to that value
// before the rename. On any failure, the temp file is removed and the
// original target file (if any) remains untouched.
func WriteText(path string, data string, newFileMode os.FileMode) error {
	dir := filepath.Dir(path)
	existed := fileExists(path)

	f, err := os.CreateTemp(dir, "apm-atomic-")
	if err != nil {
		return err
	}
	tmpName := f.Name()

	cleanup := func() {
		f.Close()
		os.Remove(tmpName)
	}

	if newFileMode > 0 && !existed {
		if err := f.Chmod(newFileMode); err != nil {
			cleanup()
			return err
		}
	}

	if _, err := f.WriteString(data); err != nil {
		cleanup()
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return err
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
