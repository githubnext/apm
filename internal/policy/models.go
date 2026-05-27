// Package policy provides data models for CI/policy audit checks.
// It mirrors the Python apm_cli.policy.models module.
package policy

// CheckResult is the result of a single CI check.
type CheckResult struct {
	Name    string   // e.g., "lockfile-exists"
	Passed  bool
	Message string   // human-readable description
	Details []string // individual violations
}

// CIAuditResult is the aggregate result of all CI checks.
type CIAuditResult struct {
	Checks []CheckResult
}

// Passed returns true if all checks passed.
func (r *CIAuditResult) Passed() bool {
	for _, c := range r.Checks {
		if !c.Passed {
			return false
		}
	}
	return true
}

// FailedChecks returns only the checks that did not pass.
func (r *CIAuditResult) FailedChecks() []CheckResult {
	var out []CheckResult
	for _, c := range r.Checks {
		if !c.Passed {
			out = append(out, c)
		}
	}
	return out
}

// ToJSON serializes to a JSON-compatible map.
func (r *CIAuditResult) ToJSON() map[string]interface{} {
	checks := make([]map[string]interface{}, len(r.Checks))
	passed, failed := 0, 0
	for i, c := range r.Checks {
		checks[i] = map[string]interface{}{
			"name":    c.Name,
			"passed":  c.Passed,
			"message": c.Message,
			"details": c.Details,
		}
		if c.Passed {
			passed++
		} else {
			failed++
		}
	}
	return map[string]interface{}{
		"passed": r.Passed(),
		"checks": checks,
		"summary": map[string]interface{}{
			"total":  len(r.Checks),
			"passed": passed,
			"failed": failed,
		},
	}
}

// CheckArtifactMap maps check names to their most relevant artifact.
var CheckArtifactMap = map[string]string{
	"lockfile-exists":           "apm.lock.yaml",
	"ref-consistency":           "apm.lock.yaml",
	"deployed-files-present":    "apm.lock.yaml",
	"no-orphaned-packages":      "apm.lock.yaml",
	"config-consistency":        "apm.lock.yaml",
	"content-integrity":         "apm.lock.yaml",
	"dependency-allowlist":      "apm.yml",
	"dependency-denylist":       "apm.yml",
	"required-packages":         "apm.yml",
	"required-packages-deployed": "apm.lock.yaml",
	"required-package-version":  "apm.lock.yaml",
	"transitive-depth":          "apm.lock.yaml",
	"mcp-allowlist":             "apm.yml",
	"mcp-denylist":              "apm.yml",
	"mcp-transport":             "apm.yml",
	"mcp-self-defined":          "apm.yml",
	"compilation-target":        "apm.yml",
	"compilation-strategy":      "apm.yml",
	"source-attribution":        "apm.yml",
	"required-manifest-fields":  "apm.yml",
	"scripts-policy":            "apm.yml",
	"unmanaged-files":           "apm.yml",
	"manifest-parse":            "apm.yml",
	"manifest-missing":          "apm.yml",
}
