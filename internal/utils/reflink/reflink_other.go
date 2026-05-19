//go:build !linux

package reflink

import "os"

func platformClone(src, dst string) (bool, error) {
	return false, nil // no reflink support on this platform
}

func platformSupported(path string) bool {
	return false
}

func deviceID(path string) (uint64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	_ = info
	return 0, nil
}
