// Package reflink provides copy-on-write file cloning (reflinks) for fast
// large-tree materialisation.
//
// Modern filesystems (APFS on macOS, btrfs and XFS on Linux) support
// copy-on-write clones. This package attempts reflinks where possible and
// falls back to regular file copies transparently.
//
// API:
//   - CloneFile: attempt to reflink one file; return true on success
//   - ReflinkSupported: best-effort runtime probe
package reflink

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

// NoReflinkEnv disables reflinks when set to "1".
const NoReflinkEnv = "APM_NO_REFLINK"

// deviceCapability caches per-device reflink support (st_dev -> bool).
var (
	deviceCapability = map[uint64]bool{}
	capMu            sync.Mutex
)

// CloneFile attempts to create a reflink clone of src at dst.
// Falls back to a regular copy if reflinks are not supported.
// Returns true if a reflink was used, false if a regular copy was used.
func CloneFile(src, dst string) (bool, error) {
	if os.Getenv(NoReflinkEnv) == "1" {
		return false, regularCopy(src, dst)
	}

	// Try platform-specific reflink
	ok, err := platformClone(src, dst)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	// Fall back to copy
	return false, regularCopy(src, dst)
}

// ReflinkSupported returns true if reflinks are likely supported on the filesystem
// containing path.
func ReflinkSupported(path string) bool {
	if os.Getenv(NoReflinkEnv) == "1" {
		return false
	}
	dev, err := deviceID(path)
	if err != nil {
		return false
	}
	capMu.Lock()
	supported, probed := deviceCapability[dev]
	capMu.Unlock()
	if probed {
		return supported
	}
	return platformSupported(path)
}

func regularCopy(src, dst string) error {
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
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// setCachedCapability records whether the device supports reflinks.
func setCachedCapability(dev uint64, supported bool) {
	capMu.Lock()
	deviceCapability[dev] = supported
	capMu.Unlock()
}
