// Package filescanner provides lockfile-driven file scanning for content integrity checks.
//
// Mirrors src/apm_cli/security/file_scanner.py.
//
// Extracted from commands/audit so the policy module can call ScanLockfilePackages
// without importing from the command layer.
package filescanner

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanFinding represents a single security finding in a file.
type ScanFinding struct {
	Type    string
	Message string
	Line    int
}

// ScanResult holds findings for a single file path label.
type ScanResult struct {
	Label    string
	Findings []ScanFinding
}

// isSafeLockfilePath returns true if a relative path from the lockfile is safe to read.
// Rejects paths containing ".." or that escape the project root.
func isSafeLockfilePath(relPath string, projectRoot string) bool {
	if strings.Contains(relPath, "..") {
		return false
	}
	abs, err := filepath.Abs(filepath.Join(projectRoot, relPath))
	if err != nil {
		return false
	}
	rootAbs, err := filepath.Abs(projectRoot)
	if err != nil {
		return false
	}
	return strings.HasPrefix(abs, rootAbs+string(os.PathSeparator)) || abs == rootAbs
}

// LockedDependency represents a dependency entry from apm.lock.yaml.
type LockedDependency struct {
	DeployedFiles []string
}

// LockFileData holds parsed lock file dependency data.
type LockFileData struct {
	Dependencies map[string]LockedDependency
}

// ScanDeployedFiles scans the deployed files listed in a lock file for security findings.
// It accepts a lockData map and a projectRoot path. packageFilter optionally restricts
// scanning to a single package key.
//
// Returns (findings by file label, total files scanned).
func ScanDeployedFiles(lockData LockFileData, projectRoot, packageFilter string) (map[string][]ScanFinding, int) {
	allFindings := map[string][]ScanFinding{}
	filesScanned := 0

	for depKey, dep := range lockData.Dependencies {
		if packageFilter != "" && depKey != packageFilter {
			continue
		}

		for _, relPath := range dep.DeployedFiles {
			safe := isSafeLockfilePath(strings.TrimRight(relPath, "/"), projectRoot)
			if !safe {
				continue
			}

			absPath := filepath.Join(projectRoot, relPath)
			info, err := os.Stat(absPath)
			if err != nil {
				continue
			}

			if info.IsDir() {
				dirFindings, dirCount := scanDir(absPath, strings.TrimRight(relPath, "/"))
				filesScanned += dirCount
				for label, findings := range dirFindings {
					allFindings[label] = findings
				}
				continue
			}

			filesScanned++
			findings := scanFile(absPath)
			if len(findings) > 0 {
				allFindings[relPath] = findings
			}
		}
	}

	return allFindings, filesScanned
}

// scanDir recursively scans all files under a directory for security findings.
func scanDir(dirPath, baseLabel string) (map[string][]ScanFinding, int) {
	findings := map[string][]ScanFinding{}
	count := 0

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return findings, count
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		label := baseLabel + "/" + entry.Name()

		if entry.IsDir() {
			sub, subCount := scanDir(fullPath, label)
			count += subCount
			for k, v := range sub {
				findings[k] = v
			}
			continue
		}

		count++
		f := scanFile(fullPath)
		if len(f) > 0 {
			findings[label] = f
		}
	}

	return findings, count
}

// scanFile scans a single file for suspicious patterns.
// Returns findings if any suspicious content is detected.
func scanFile(path string) []ScanFinding {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return detectSuspiciousBytes(data)
}

// suspiciousRunes contains Unicode codepoints that are suspicious in source files.
var suspiciousRunes = []struct {
	codepoint rune
	name      string
}{
	{0x200B, "zero-width space"},
	{0x200C, "zero-width non-joiner"},
	{0x200D, "zero-width joiner"},
	{0x202A, "left-to-right embedding"},
	{0x202B, "right-to-left embedding"},
	{0x202C, "pop directional formatting"},
	{0x202D, "left-to-right override"},
	{0x202E, "right-to-left override"},
	{0x2066, "left-to-right isolate"},
	{0x2067, "right-to-left isolate"},
	{0x2068, "first strong isolate"},
	{0x2069, "pop directional isolate"},
	{0xFEFF, "byte order mark / zero-width no-break space"},
}

func detectSuspiciousBytes(data []byte) []ScanFinding {
	var findings []ScanFinding
	content := string(data)
	for _, sr := range suspiciousRunes {
		if strings.ContainsRune(content, sr.codepoint) {
			findings = append(findings, ScanFinding{
				Type:    "hidden-character",
				Message: "file contains " + sr.name,
			})
		}
	}
	return findings
}
