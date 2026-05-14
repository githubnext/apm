// Package securityscan provides the pre-deploy security scan helper for the install pipeline.
// Migrated from src/apm_cli/install/helpers/security_scan.py
//
// Wraps the SecurityGate scanner used by the install pipeline. The scan detects
// hidden characters (zero-width joiners, bidirectional overrides, etc.) that could
// be used to smuggle malicious payloads into prompts, skills, or agent definitions.
package securityscan

import (
	"fmt"
	"os"
	"path/filepath"
)

// Finding represents a single security finding in a file.
type Finding struct {
	// FilePath is the file where the finding was detected.
	FilePath string
	// Description describes the hidden-character pattern found.
	Description string
	// Line is the 1-based line number (0 = unknown).
	Line int
}

// ScanResult holds the outcome of a pre-deploy security scan.
type ScanResult struct {
	// HasFindings is true when at least one finding was detected.
	HasFindings bool
	// ShouldBlock is true when the finding severity warrants blocking install.
	ShouldBlock bool
	// Findings is the list of detected findings.
	Findings []Finding
	// FilesScanned is the number of files that were examined.
	FilesScanned int
}

// scannerFunc is the function signature for the security gate scan.
// Provided as a variable so tests can inject a stub.
var scannerFunc func(root string, force bool) (*ScanResult, error) = defaultScanner

// defaultScanner performs a simple hidden-character scan on all files under root.
// This is a lightweight stdlib-only implementation; the full Python SecurityGate
// uses a richer classification engine (separate migration).
func defaultScanner(root string, force bool) (*ScanResult, error) {
	result := &ScanResult{}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if d.IsDir() {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		result.FilesScanned++

		findings := scanBytes(path, data)
		if len(findings) > 0 {
			result.HasFindings = true
			result.ShouldBlock = true
			result.Findings = append(result.Findings, findings...)
		}
		return nil
	})
	return result, err
}

// hiddenPatterns are Unicode code-points considered dangerous in prompt/skill files.
var hiddenPatterns = []rune{
	'\u200B', // zero-width space
	'\u200C', // zero-width non-joiner
	'\u200D', // zero-width joiner
	'\u2028', // line separator
	'\u2029', // paragraph separator
	'\u202A', // left-to-right embedding
	'\u202B', // right-to-left embedding
	'\u202C', // pop directional formatting
	'\u202D', // left-to-right override
	'\u202E', // right-to-left override (most dangerous)
	'\uFEFF', // byte order mark (mid-file)
	'\u00AD', // soft hyphen
}

func scanBytes(path string, data []byte) []Finding {
	var findings []Finding
	text := string(data)
	for _, r := range hiddenPatterns {
		for i, c := range text {
			if c == r {
				findings = append(findings, Finding{
					FilePath:    path,
					Description: fmt.Sprintf("hidden character U+%04X at byte offset %d", r, i),
				})
				break // one finding per pattern per file
			}
		}
	}
	return findings
}

// PreDeploySecurityScan scans package source files for hidden characters BEFORE deployment.
//
// Returns true if deployment should proceed, false to block.
// When force is true the scan still runs but never returns false (block is suppressed).
func PreDeploySecurityScan(installPath string, packageName string, force bool) (bool, *ScanResult) {
	result, err := scannerFunc(installPath, force)
	if err != nil || result == nil {
		// Scan error -- allow deployment to proceed (fail-open)
		return true, &ScanResult{}
	}
	if !result.HasFindings {
		return true, result
	}
	if force || !result.ShouldBlock {
		return true, result
	}
	return false, result
}
