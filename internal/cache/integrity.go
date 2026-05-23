package cache

import (
	"os"
	"path/filepath"
	"strings"
)

// VerifyCheckoutSHA verifies that a cached checkout's HEAD matches the expected SHA.
// Reads .git/HEAD (and follows refs / packed-refs as needed) rather than spawning
// "git rev-parse": faster, and cannot be influenced by a poisoned local .git/config.
func VerifyCheckoutSHA(checkoutDir, expectedSHA string) bool {
	if _, err := os.Stat(checkoutDir); err != nil {
		return false
	}
	actualSHA := readHeadSHA(checkoutDir)
	if actualSHA == "" {
		return false
	}
	return actualSHA == strings.TrimSpace(strings.ToLower(expectedSHA))
}

// readHeadSHA returns the resolved 40-char SHA at HEAD, or "" on any failure.
func readHeadSHA(checkoutDir string) string {
	gitPath := filepath.Join(checkoutDir, ".git")

	fi, err := os.Stat(gitPath)
	if err != nil {
		return ""
	}

	var gitDir string
	if fi.Mode().IsRegular() {
		// Worktree pointer: "gitdir: <path>"
		content, err := os.ReadFile(gitPath)
		if err != nil {
			return ""
		}
		line := strings.TrimSpace(string(content))
		if !strings.HasPrefix(line, "gitdir:") {
			return ""
		}
		target := strings.TrimSpace(line[len("gitdir:"):])
		abs := filepath.Join(checkoutDir, target)
		resolved, err := filepath.Abs(abs)
		if err != nil {
			return ""
		}
		gitDir = resolved
	} else if fi.IsDir() {
		gitDir = gitPath
	} else {
		return ""
	}

	headPath := filepath.Join(gitDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	head := strings.TrimSpace(string(headContent))

	if strings.HasPrefix(head, "ref:") {
		refTarget := strings.TrimSpace(head[len("ref:"):])
		refPath := filepath.Join(gitDir, refTarget)
		if data, err := os.ReadFile(refPath); err == nil {
			return strings.TrimSpace(strings.ToLower(string(data)))
		}
		// Try packed-refs
		packedPath := filepath.Join(gitDir, "packed-refs")
		if packed, err := os.ReadFile(packedPath); err == nil {
			for _, raw := range strings.Split(string(packed), "\n") {
				line := strings.TrimSpace(raw)
				if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "^") {
					continue
				}
				parts := strings.SplitN(line, " ", 2)
				if len(parts) == 2 && parts[1] == refTarget {
					return strings.ToLower(parts[0])
				}
			}
		}
		return ""
	}

	// Detached HEAD: should be a 40-char hex SHA
	lower := strings.ToLower(head)
	if len(lower) == 40 && isHex(lower) {
		return lower
	}
	return ""
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
