//go:build windows

package fileops

import (
	"errors"
	"strings"
	"syscall"
)

// isTransientLockError returns true for Windows winerror 32 or 5.
func isTransientLockError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno == syscall.ERROR_SHARING_VIOLATION || errno == syscall.ERROR_ACCESS_DENIED
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "used by another process") || strings.Contains(s, "access is denied")
}
