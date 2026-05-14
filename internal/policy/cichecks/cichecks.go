// Package cichecks implements baseline CI checks for lockfile consistency.
// These checks run without any policy file, validating on-disk state against
// the lockfile. Mirrors src/apm_cli/policy/ci_checks.py.
package cichecks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CheckResult is the outcome of a single baseline check.
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
	Details []string
}

// HasFailures returns true when the check failed.
func (r CheckResult) HasFailures() bool { return !r.Passed }

// CIAuditResult aggregates multiple check results.
type CIAuditResult struct {
	Checks []CheckResult
}

// HasFailures returns true when any check failed.
func (r CIAuditResult) HasFailures() bool {
	for _, c := range r.Checks {
		if !c.Passed {
			return true
		}
	}
	return false
}

// RenderSummary returns a human-readable summary.
func (r CIAuditResult) RenderSummary() string {
	var sb strings.Builder
	for _, c := range r.Checks {
		sym := "[+]"
		if !c.Passed {
			sym = "[x]"
		}
		sb.WriteString(fmt.Sprintf("%s %s: %s\n", sym, c.Name, c.Message))
		for _, d := range c.Details {
			sb.WriteString("    " + d + "\n")
		}
	}
	return sb.String()
}

// LockedDepInfo is the minimum information about a locked dependency.
type LockedDepInfo struct {
	Key           string
	ResolvedRef   string
	ManifestRef   string // what apm.yml declares
	DeployedFiles []string
	ContentHash   string
}

// DriftFinding describes a single drift between expected and actual state.
type DriftFinding struct {
	DepKey   string
	FilePath string
	Reason   string
}

// CheckManifestParse returns a pass result to indicate the manifest was
// successfully parsed (the parse itself happens at the call site).
func CheckManifestParse() CheckResult {
	return CheckResult{Name: "manifest-parse", Passed: true, Message: "apm.yml parsed successfully"}
}

// CheckManifestParseFailed returns the failure result for a manifest parse error.
func CheckManifestParseFailed(err error) CheckResult {
	return CheckResult{
		Name:    "manifest-parse",
		Passed:  false,
		Message: fmt.Sprintf("apm.yml parse error: %v", err),
		Details: []string{err.Error()},
	}
}

// CheckLockfileExists verifies that apm.lock.yaml is present when needed.
func CheckLockfileExists(projectRoot string, hasDeps bool) CheckResult {
	lockPath := filepath.Join(projectRoot, "apm.lock.yaml")
	if !hasDeps {
		return CheckResult{Name: "lockfile-exists", Passed: true, Message: "No dependencies declared -- lockfile not required"}
	}
	if _, err := os.Stat(lockPath); err == nil {
		return CheckResult{Name: "lockfile-exists", Passed: true, Message: "Lockfile present"}
	}
	return CheckResult{
		Name:    "lockfile-exists",
		Passed:  false,
		Message: "Lockfile missing -- run 'apm install' to generate apm.lock.yaml",
		Details: []string{"apm.yml declares dependencies but apm.lock.yaml is absent"},
	}
}

// CheckLockfileSync verifies that every manifest dependency has a lockfile entry.
func CheckLockfileSync(manifestKeys, lockfileKeys map[string]bool) CheckResult {
	var missing []string
	for k := range manifestKeys {
		if !lockfileKeys[k] {
			missing = append(missing, k)
		}
	}
	if len(missing) == 0 {
		return CheckResult{Name: "lockfile-sync", Passed: true, Message: "Lockfile in sync with manifest"}
	}
	return CheckResult{
		Name:    "lockfile-sync",
		Passed:  false,
		Message: fmt.Sprintf("%d dep(s) in manifest but missing from lockfile", len(missing)),
		Details: missing,
	}
}

// CheckRefConsistency verifies that every dep's manifest ref matches the
// lockfile resolved_ref.
func CheckRefConsistency(deps []LockedDepInfo) CheckResult {
	var mismatches []string
	for _, dep := range deps {
		if dep.ManifestRef != "" && dep.ResolvedRef != "" && dep.ManifestRef != dep.ResolvedRef {
			mismatches = append(mismatches, fmt.Sprintf("%s: manifest=%q lockfile=%q", dep.Key, dep.ManifestRef, dep.ResolvedRef))
		}
	}
	if len(mismatches) == 0 {
		return CheckResult{Name: "ref-consistency", Passed: true, Message: "All dependency refs consistent"}
	}
	return CheckResult{
		Name:    "ref-consistency",
		Passed:  false,
		Message: fmt.Sprintf("%d ref mismatch(es) between manifest and lockfile", len(mismatches)),
		Details: mismatches,
	}
}

// CheckDeployedFilesPresent verifies that every deployed file in the lockfile
// exists on disk.
func CheckDeployedFilesPresent(projectRoot string, deps []LockedDepInfo) CheckResult {
	var missing []string
	for _, dep := range deps {
		for _, rel := range dep.DeployedFiles {
			full := filepath.Join(projectRoot, rel)
			if _, err := os.Stat(full); err != nil {
				missing = append(missing, fmt.Sprintf("%s: %s", dep.Key, rel))
			}
		}
	}
	if len(missing) == 0 {
		return CheckResult{Name: "deployed-files-present", Passed: true, Message: "All deployed files present on disk"}
	}
	return CheckResult{
		Name:    "deployed-files-present",
		Passed:  false,
		Message: fmt.Sprintf("%d deployed file(s) missing from disk", len(missing)),
		Details: missing,
	}
}

// CheckDriftFindings returns a check result based on drift scan findings.
func CheckDriftFindings(findings []DriftFinding) CheckResult {
	if len(findings) == 0 {
		return CheckResult{Name: "content-integrity", Passed: true, Message: "No drift detected"}
	}
	var details []string
	for _, f := range findings {
		details = append(details, fmt.Sprintf("%s / %s: %s", f.DepKey, f.FilePath, f.Reason))
	}
	return CheckResult{
		Name:    "content-integrity",
		Passed:  false,
		Message: fmt.Sprintf("%d file(s) have drifted from the lockfile", len(findings)),
		Details: details,
	}
}

// RunBaselineChecks executes all baseline checks and returns a CIAuditResult.
// manifestParsed is true when apm.yml was found and parsed without error.
// hasDeps is true when the manifest declares APM or MCP dependencies.
func RunBaselineChecks(
	projectRoot string,
	manifestParsed bool,
	manifestParseErr error,
	hasDeps bool,
	manifestKeys map[string]bool,
	lockfileKeys map[string]bool,
	deps []LockedDepInfo,
	driftFindings []DriftFinding,
) CIAuditResult {
	var checks []CheckResult

	if !manifestParsed {
		checks = append(checks, CheckManifestParseFailed(manifestParseErr))
		return CIAuditResult{Checks: checks}
	}
	checks = append(checks, CheckManifestParse())
	checks = append(checks, CheckLockfileExists(projectRoot, hasDeps))
	if hasDeps {
		checks = append(checks, CheckLockfileSync(manifestKeys, lockfileKeys))
		checks = append(checks, CheckRefConsistency(deps))
		checks = append(checks, CheckDeployedFilesPresent(projectRoot, deps))
		checks = append(checks, CheckDriftFindings(driftFindings))
	}
	return CIAuditResult{Checks: checks}
}
