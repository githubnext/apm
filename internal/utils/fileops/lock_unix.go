//go:build !windows

package fileops

import (
	"errors"
	"syscall"
)

// isTransientLockError returns true for EBUSY on Unix.
func isTransientLockError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno == syscall.EBUSY
	}
	return false
}
