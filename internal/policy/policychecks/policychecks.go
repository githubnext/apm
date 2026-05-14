// Package policychecks implements organisational governance enforcement checks.
// Mirrors src/apm_cli/policy/policy_checks.py.
package policychecks

import (
	"fmt"
	"os"
	"strings"
)

// CheckResult is the outcome of a single policy check.
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
	Details []string
}

// HasFailures returns true when the result represents a failure.
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

// RenderSummary returns a human-readable summary of all checks.
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

// DependencyPolicy is the minimal policy struct needed by the checks.
type DependencyPolicy struct {
	Allow   []string
	Deny    []string
	Require []string
}

// DependencyRef is a minimal reference to a resolved dependency.
type DependencyRef struct {
	CanonicalString string
	IsLocal         bool
}

// CheckDependencyAllowlist verifies that every dep matches the policy allow list.
func CheckDependencyAllowlist(deps []DependencyRef, policy DependencyPolicy) CheckResult {
	if len(policy.Allow) == 0 {
		return CheckResult{
			Name:    "dependency-allowlist",
			Passed:  true,
			Message: "No dependency allow list configured",
		}
	}
	var violations []string
	for _, dep := range deps {
		if dep.IsLocal {
			continue
		}
		matched := false
		for _, pattern := range policy.Allow {
			if globMatch(pattern, dep.CanonicalString) {
				matched = true
				break
			}
		}
		if !matched {
			violations = append(violations, fmt.Sprintf("%s: not in allowed list", dep.CanonicalString))
		}
	}
	if len(violations) == 0 {
		return CheckResult{Name: "dependency-allowlist", Passed: true, Message: "All dependencies match allow list"}
	}
	return CheckResult{
		Name:    "dependency-allowlist",
		Passed:  false,
		Message: fmt.Sprintf("%d dependency(ies) not in allow list", len(violations)),
		Details: violations,
	}
}

// CheckDependencyDenylist verifies that no dep matches the policy deny list.
func CheckDependencyDenylist(deps []DependencyRef, policy DependencyPolicy) CheckResult {
	if len(policy.Deny) == 0 {
		return CheckResult{Name: "dependency-denylist", Passed: true, Message: "No dependency deny list configured"}
	}
	var violations []string
	for _, dep := range deps {
		if dep.IsLocal {
			continue
		}
		for _, pattern := range policy.Deny {
			if globMatch(pattern, dep.CanonicalString) {
				violations = append(violations, fmt.Sprintf("%s: denied by pattern %q", dep.CanonicalString, pattern))
				break
			}
		}
	}
	if len(violations) == 0 {
		return CheckResult{Name: "dependency-denylist", Passed: true, Message: "No dependencies match deny list"}
	}
	return CheckResult{
		Name:    "dependency-denylist",
		Passed:  false,
		Message: fmt.Sprintf("%d dependency(ies) match deny list", len(violations)),
		Details: violations,
	}
}

// CheckRequiredPackages verifies every required package is in the manifest.
func CheckRequiredPackages(deps []DependencyRef, policy DependencyPolicy) CheckResult {
	if len(policy.Require) == 0 {
		return CheckResult{Name: "required-packages", Passed: true, Message: "No required packages configured"}
	}
	depNames := map[string]bool{}
	for _, d := range deps {
		base := strings.SplitN(d.CanonicalString, "#", 2)[0]
		depNames[base] = true
	}
	var missing []string
	for _, req := range policy.Require {
		pkgName := strings.SplitN(req, "#", 2)[0]
		if !depNames[pkgName] {
			missing = append(missing, pkgName)
		}
	}
	if len(missing) == 0 {
		return CheckResult{Name: "required-packages", Passed: true, Message: "All required packages present in manifest"}
	}
	return CheckResult{
		Name:    "required-packages",
		Passed:  false,
		Message: fmt.Sprintf("%d required package(s) missing from manifest", len(missing)),
		Details: missing,
	}
}

// CheckCompilationTarget verifies the apm.yml compilation target matches
// the policy-required value.
func CheckCompilationTarget(actualTarget string, requiredTarget string) CheckResult {
	if requiredTarget == "" {
		return CheckResult{Name: "compilation-target", Passed: true, Message: "No compilation target required by policy"}
	}
	if actualTarget == requiredTarget {
		return CheckResult{Name: "compilation-target", Passed: true, Message: fmt.Sprintf("Compilation target matches policy: %q", requiredTarget)}
	}
	return CheckResult{
		Name:    "compilation-target",
		Passed:  false,
		Message: fmt.Sprintf("Compilation target mismatch: got %q, policy requires %q", actualTarget, requiredTarget),
	}
}

// CheckExtensionsPresent verifies required apm.yml extension keys are present.
func CheckExtensionsPresent(presentExtensions map[string]bool, requiredExtensions []string) CheckResult {
	if len(requiredExtensions) == 0 {
		return CheckResult{Name: "extensions-present", Passed: true, Message: "No extensions required by policy"}
	}
	var missing []string
	for _, ext := range requiredExtensions {
		if !presentExtensions[ext] {
			missing = append(missing, ext)
		}
	}
	if len(missing) == 0 {
		return CheckResult{Name: "extensions-present", Passed: true, Message: "All required extensions present"}
	}
	return CheckResult{
		Name:    "extensions-present",
		Passed:  false,
		Message: fmt.Sprintf("%d required extension(s) missing", len(missing)),
		Details: missing,
	}
}

// LoadRawApmYML reads apm.yml at projectRoot as raw key-value pairs.
// Returns nil when the file is absent, unreadable, or malformed.
func LoadRawApmYML(projectRoot string) map[string]interface{} {
	path := projectRoot + "/apm.yml"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	// Minimal YAML key scanner -- extracts top-level keys only.
	result := map[string]interface{}{}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "#") || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" || strings.Contains(key, " ") {
			continue
		}
		val := strings.TrimSpace(parts[1])
		result[key] = val
	}
	return result
}

// globMatch is a minimal glob pattern matcher supporting * and ? wildcards.
func globMatch(pattern, str string) bool {
	if pattern == "" {
		return str == ""
	}
	if pattern == "*" {
		return true
	}
	// Simple recursive match -- sufficient for dep pattern matching.
	if pattern[0] == '*' {
		for i := 0; i <= len(str); i++ {
			if globMatch(pattern[1:], str[i:]) {
				return true
			}
		}
		return false
	}
	if len(str) == 0 {
		return false
	}
	if pattern[0] == '?' || pattern[0] == str[0] {
		return globMatch(pattern[1:], str[1:])
	}
	return false
}
