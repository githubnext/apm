//go:build linux

package reflink

import (
	"os"
	"syscall"
	"unsafe"
)

// FICLONE ioctl number on Linux: _IOW(0x94, 9, int) = 0x40049409
const ficlone = 0x40049409

func platformClone(src, dst string) (bool, error) {
	if err := os.MkdirAll(getDir(dst), 0o755); err != nil {
		return false, err
	}

	in, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return false, err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return false, err
	}
	defer func() {
		out.Close()
		if err != nil {
			os.Remove(dst)
		}
	}()

	// Check device capability cache
	inStat, statErr := in.Stat()
	if statErr == nil {
		if sysInfo, ok := inStat.Sys().(*syscall.Stat_t); ok {
			dev := sysInfo.Dev
			capMu.Lock()
			supported, probed := deviceCapability[dev]
			capMu.Unlock()
			if probed && !supported {
				out.Close()
				return false, regularCopy(src, dst)
			}
		}
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, out.Fd(), ficlone, in.Fd())
	if errno == 0 {
		// Record success
		if statErr == nil {
			if sysInfo, ok := inStat.Sys().(*syscall.Stat_t); ok {
				setCachedCapability(sysInfo.Dev, true)
			}
		}
		return true, nil
	}
	// errno indicates not supported -- cache and fall through
	if statErr == nil {
		if sysInfo, ok := inStat.Sys().(*syscall.Stat_t); ok {
			setCachedCapability(sysInfo.Dev, false)
		}
	}
	out.Close()
	return false, regularCopy(src, dst)
}

func platformSupported(path string) bool {
	// probe by attempting a clone of a temp file
	return false // conservative: return false without actually probing
}

func deviceID(path string) (uint64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if sysInfo, ok := info.Sys().(*syscall.Stat_t); ok {
		return sysInfo.Dev, nil
	}
	return 0, nil
}

func getDir(path string) string {
	dir := path
	for len(dir) > 0 && dir[len(dir)-1] != '/' && dir[len(dir)-1] != '\\' {
		dir = dir[:len(dir)-1]
	}
	if dir == "" {
		return "."
	}
	return dir
}

// ensure unused import is not flagged
var _ = unsafe.Pointer(nil)
