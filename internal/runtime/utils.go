package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// FindRuntimeBinary resolves the path to a runtime binary.
// Priority:
//  1. ~/.apm/runtimes/<name>  (APM-managed, executable)
//  2. PATH lookup via exec.LookPath
//
// Returns empty string if not found.
// Returns error if name contains path-traversal sequences.
func FindRuntimeBinary(name string) (string, error) {
	// Security: reject names with path separators
	for _, ch := range []string{"/", "\\"} {
		if containsStr(name, ch) {
			return "", fmt.Errorf("invalid runtime name %q: must be a plain binary name without path separators", name)
		}
	}
	if name == "" || name == "." || name == ".." {
		return "", fmt.Errorf("invalid runtime name %q", name)
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		apmRuntimes := filepath.Join(homeDir, ".apm", "runtimes")
		if runtime.GOOS == "windows" {
			candidate := filepath.Join(apmRuntimes, name+".exe")
			if isExecutableFile(candidate, apmRuntimes) {
				return candidate, nil
			}
		}
		candidate := filepath.Join(apmRuntimes, name)
		if isExecutableFile(candidate, apmRuntimes) {
			return candidate, nil
		}
	}

	// PATH fallback
	p, err := exec.LookPath(name)
	if err != nil {
		return "", nil
	}
	return p, nil
}

func isExecutableFile(path, base string) bool {
	// Ensure path is within base (security guard)
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absBase, err := filepath.Abs(base)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absBase, abs)
	if err != nil {
		return false
	}
	if len(rel) >= 2 && rel[:2] == ".." {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	// Check executable bit (non-Windows)
	if runtime.GOOS != "windows" {
		if info.Mode()&0o111 == 0 {
			return false
		}
	}
	return true
}

func containsStr(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, rune(sub[0])))
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
