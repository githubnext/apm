// Package cleanuphelper provides a shared helper for removing stale deployed
// files after an APM install.
//
// Mirrors src/apm_cli/integration/cleanup.py.
//
// Safety gates (applied in order):
//  1. Path validation -- reject traversal and paths not under a known prefix.
//  2. Directory rejection -- APM only manages individual files.
//  3. Provenance check -- if APM recorded a hash, the on-disk content must
//     still match. Fails CLOSED on hash-read errors.
package cleanuphelper

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const coworkURIScheme = "cowork://"

// Diagnostic captures a single recoverable warning.
type Diagnostic struct {
	Package string
	Message string
}

// DiagnosticCollector accumulates non-fatal warnings during cleanup.
type DiagnosticCollector struct {
	Warnings []Diagnostic
}

// Warn records a warning associated with a package key.
func (d *DiagnosticCollector) Warn(pkg, msg string) {
	d.Warnings = append(d.Warnings, Diagnostic{Package: pkg, Message: msg})
}

// ValidateDeployPath is the path security gate. It rejects:
//   - paths with ".." components (traversal)
//   - cowork:// URIs (handled separately)
//   - absolute paths
//   - paths not starting with one of the allowed integration prefixes
func ValidateDeployPath(stalePath string, projectRoot string, integrationPrefixes []string) bool {
	if strings.HasPrefix(stalePath, coworkURIScheme) {
		return false
	}
	if filepath.IsAbs(stalePath) {
		return false
	}
	if strings.Contains(stalePath, "..") {
		return false
	}
	for _, prefix := range integrationPrefixes {
		if strings.HasPrefix(stalePath, prefix) {
			return true
		}
	}
	return false
}

// computeFileHash returns the SHA-256 hash of path in the form "sha256:<hex>".
func computeFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}

// stripSHA256Prefix removes the "sha256:" prefix from a hash string, if
// present, for normalised comparison.
func stripSHA256Prefix(h string) string {
	if strings.HasPrefix(h, "sha256:") {
		return h[len("sha256:"):]
	}
	return h
}

// CleanupResult summarises the outcome of a stale-file cleanup pass for a
// single package.
type CleanupResult struct {
	Deleted         []string // workspace-relative paths removed from disk
	Failed          []string // paths that raised during removal (retained for retry)
	SkippedUserEdit []string // paths skipped because the user edited the file
	SkippedUnmanaged []string // paths refused by safety gates
	DeletedTargets  []string // absolute paths of deleted entries
}

// Options configures RemoveStaleDeployedFiles.
type Options struct {
	// DepKey is the unique key of the package (used for diagnostic attribution).
	DepKey string
	// ProjectRoot is the project root directory.
	ProjectRoot string
	// IntegrationPrefixes are the allowed workspace-relative path prefixes.
	IntegrationPrefixes []string
	// RecordedHashes maps rel-path -> "sha256:<hex>" as stored in the
	// previous lockfile. Empty disables provenance checking.
	RecordedHashes map[string]string
	// FailedPathRetained controls the wording of failure diagnostics:
	//   true  = caller will re-insert failed paths into deployed_files (intra-package stale cleanup)
	//   false = package is being removed; failed paths cannot be retained (orphan cleanup)
	FailedPathRetained bool
	// Diagnostics accumulates recoverable warnings.
	Diagnostics *DiagnosticCollector
	// CoworkRootResolver, when non-nil, is called to resolve cowork:// URIs to
	// absolute paths. Return ("", nil) when the cowork root is unavailable.
	CoworkRootResolver func() (string, error)
	// CoworkFromLockfilePath, when non-nil, maps a cowork:// URI + resolved
	// root to an absolute path. Returns an error on containment violations.
	CoworkFromLockfilePath func(uri, coworkRoot string) (string, error)
}

// RemoveStaleDeployedFiles removes APM-deployed files that are no longer
// produced by opts.DepKey.
//
// stalePaths contains workspace-relative paths flagged as stale. The function
// applies three safety gates before deleting each file. See the package-level
// documentation for the gate ordering.
func RemoveStaleDeployedFiles(stalePaths []string, opts Options) CleanupResult {
	if opts.Diagnostics == nil {
		opts.Diagnostics = &DiagnosticCollector{}
	}
	if opts.RecordedHashes == nil {
		opts.RecordedHashes = map[string]string{}
	}

	sorted := make([]string, len(stalePaths))
	copy(sorted, stalePaths)
	sort.Strings(sorted)

	result := CleanupResult{}

	var coworkRootResolved bool
	var coworkRootCached string
	var coworkOrphansSkipped int
	var coworkResolveErrors int

	for _, stalePath := range sorted {
		// ── Cowork:// paths ────────────────────────────────────────────────
		if strings.HasPrefix(stalePath, coworkURIScheme) {
			if strings.Contains(stalePath, "..") {
				result.SkippedUnmanaged = append(result.SkippedUnmanaged, stalePath)
				continue
			}
			hasPrefix := false
			for _, prefix := range opts.IntegrationPrefixes {
				if strings.HasPrefix(stalePath, prefix) {
					hasPrefix = true
					break
				}
			}
			if !hasPrefix {
				result.SkippedUnmanaged = append(result.SkippedUnmanaged, stalePath)
				continue
			}
			// Resolve cowork:// URI.
			if !coworkRootResolved && opts.CoworkRootResolver != nil {
				root, err := opts.CoworkRootResolver()
				if err == nil {
					coworkRootCached = root
				}
				coworkRootResolved = true
			}
			if coworkRootCached == "" {
				coworkOrphansSkipped++
				result.Failed = append(result.Failed, stalePath)
				continue
			}
			if opts.CoworkFromLockfilePath == nil {
				coworkResolveErrors++
				result.Failed = append(result.Failed, stalePath)
				continue
			}
			staleTarget, err := opts.CoworkFromLockfilePath(stalePath, coworkRootCached)
			if err != nil {
				coworkResolveErrors++
				result.Failed = append(result.Failed, stalePath)
				continue
			}
			// Fall through to common delete logic below using staleTarget.
			if err := deleteFile(staleTarget, stalePath, opts, &result); err != nil {
				// handled inside deleteFile
				_ = err
			}
			continue
		}

		// ── Non-cowork paths ────────────────────────────────────────────────
		if !ValidateDeployPath(stalePath, opts.ProjectRoot, opts.IntegrationPrefixes) {
			result.SkippedUnmanaged = append(result.SkippedUnmanaged, stalePath)
			continue
		}
		staleTarget := filepath.Join(opts.ProjectRoot, stalePath)

		info, err := os.Lstat(staleTarget)
		if os.IsNotExist(err) {
			// Already gone -- treat as cleaned.
			continue
		}
		if err != nil {
			result.Failed = append(result.Failed, stalePath)
			continue
		}

		// Gate 2: directory rejection.
		if info.IsDir() {
			result.SkippedUnmanaged = append(result.SkippedUnmanaged, stalePath)
			opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
				"Refused to remove directory entry %s: APM only deletes individual files. "+
					"If this entry was added by a malicious or corrupt lockfile, remove it manually "+
					"from apm.lock.yaml.",
				stalePath,
			))
			continue
		}

		// Gate 3: provenance check.
		if expectedHash, ok := opts.RecordedHashes[stalePath]; ok && expectedHash != "" {
			actualHash, err := computeFileHash(staleTarget)
			if err != nil {
				result.SkippedUserEdit = append(result.SkippedUserEdit, stalePath)
				opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
					"Skipped removing %s: could not verify file content (%v). "+
						"Inspect the file and delete it manually if no longer needed.",
					stalePath, err,
				))
				continue
			}
			if stripSHA256Prefix(actualHash) != stripSHA256Prefix(expectedHash) {
				result.SkippedUserEdit = append(result.SkippedUserEdit, stalePath)
				opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
					"Skipped removing %s: file has been edited since APM deployed it. "+
						"Delete it manually if you no longer need it, or ignore this warning to keep your changes.",
					stalePath,
				))
				continue
			}
		}

		// All gates passed -- safe to delete.
		if err := os.Remove(staleTarget); err != nil {
			result.Failed = append(result.Failed, stalePath)
			if opts.FailedPathRetained {
				opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
					"Could not remove stale file %s: %v. "+
						"Path retained in lockfile; will retry on next 'apm install'.",
					stalePath, err,
				))
			} else {
				opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
					"Could not remove orphaned file %s: %v. "+
						"The owning package is no longer in apm.yml -- delete the file manually.",
					stalePath, err,
				))
			}
		} else {
			result.Deleted = append(result.Deleted, stalePath)
			result.DeletedTargets = append(result.DeletedTargets, staleTarget)
		}
	}

	// One-time warnings for cowork edge cases.
	if coworkOrphansSkipped > 0 {
		noun := "entry"
		if coworkOrphansSkipped != 1 {
			noun = "entries"
		}
		opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
			"Cowork: skipping %d stale lockfile %s -- OneDrive path not detected.\n"+
				"Run: apm config set copilot-cowork-skills-dir <path>  "+
				"(or set APM_COPILOT_COWORK_SKILLS_DIR)\n"+
				"to clean up these entries on the next install/uninstall.",
			coworkOrphansSkipped, noun,
		))
	}
	if coworkResolveErrors > 0 {
		noun := "entry"
		if coworkResolveErrors != 1 {
			noun = "entries"
		}
		opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
			"Cowork: %d lockfile %s failed path resolution "+
				"(containment violation or malformed path). Paths retained for manual inspection.",
			coworkResolveErrors, noun,
		))
	}

	return result
}

// deleteFile is a helper used for the cowork branch to apply gate 2, gate 3,
// and the actual removal using an already-resolved absolute target path.
func deleteFile(staleTarget, stalePath string, opts Options, result *CleanupResult) error {
	info, err := os.Lstat(staleTarget)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		result.Failed = append(result.Failed, stalePath)
		return err
	}
	if info.IsDir() {
		result.SkippedUnmanaged = append(result.SkippedUnmanaged, stalePath)
		opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
			"Refused to remove directory entry %s: APM only deletes individual files.",
			stalePath,
		))
		return nil
	}
	if expectedHash, ok := opts.RecordedHashes[stalePath]; ok && expectedHash != "" {
		actualHash, err := computeFileHash(staleTarget)
		if err != nil {
			result.SkippedUserEdit = append(result.SkippedUserEdit, stalePath)
			opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
				"Skipped removing %s: could not verify file content (%v).", stalePath, err,
			))
			return nil
		}
		if stripSHA256Prefix(actualHash) != stripSHA256Prefix(expectedHash) {
			result.SkippedUserEdit = append(result.SkippedUserEdit, stalePath)
			opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
				"Skipped removing %s: file has been edited since APM deployed it.", stalePath,
			))
			return nil
		}
	}
	if err := os.Remove(staleTarget); err != nil {
		result.Failed = append(result.Failed, stalePath)
		opts.Diagnostics.Warn(opts.DepKey, fmt.Sprintf(
			"Could not remove stale file %s: %v.", stalePath, err,
		))
		return err
	}
	result.Deleted = append(result.Deleted, stalePath)
	result.DeletedTargets = append(result.DeletedTargets, staleTarget)
	return nil
}
