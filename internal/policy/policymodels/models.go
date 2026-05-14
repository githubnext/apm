// Package policymodels provides data models for CI/policy audit checks.
//
// Mirrors src/apm_cli/policy/models.py.
package policymodels

import (
	"encoding/json"
	"fmt"
)

// checkArtifactMap maps check names to their most relevant artifact for SARIF
// location reporting.
var checkArtifactMap = map[string]string{
	"lockfile-exists":            "apm.lock.yaml",
	"ref-consistency":            "apm.lock.yaml",
	"deployed-files-present":     "apm.lock.yaml",
	"no-orphaned-packages":       "apm.lock.yaml",
	"config-consistency":         "apm.lock.yaml",
	"content-integrity":          "apm.lock.yaml",
	"dependency-allowlist":       "apm.yml",
	"dependency-denylist":        "apm.yml",
	"required-packages":          "apm.yml",
	"required-packages-deployed": "apm.lock.yaml",
	"required-package-version":   "apm.lock.yaml",
	"transitive-depth":           "apm.lock.yaml",
	"mcp-allowlist":              "apm.yml",
	"mcp-denylist":               "apm.yml",
	"mcp-transport":              "apm.yml",
	"mcp-self-defined":           "apm.yml",
	"compilation-target":         "apm.yml",
	"compilation-strategy":       "apm.yml",
	"source-attribution":         "apm.yml",
	"required-manifest-fields":   "apm.yml",
	"scripts-policy":             "apm.yml",
	"unmanaged-files":            "apm.yml",
	"manifest-parse":             "apm.yml",
}

// ArtifactForCheck returns the most relevant artifact filename for a check name.
// Falls back to "apm.lock.yaml" for unknown checks.
func ArtifactForCheck(checkName string) string {
	if artifact, ok := checkArtifactMap[checkName]; ok {
		return artifact
	}
	return "apm.lock.yaml"
}

// CheckResult holds the result of a single CI check.
type CheckResult struct {
	Name    string   // e.g. "lockfile-exists"
	Passed  bool
	Message string   // human-readable description
	Details []string // individual violations
}

// CIAuditResult is the aggregate result of all CI checks.
type CIAuditResult struct {
	Checks []CheckResult
}

// Passed returns true when all checks passed.
func (r *CIAuditResult) Passed() bool {
	for i := range r.Checks {
		if !r.Checks[i].Passed {
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

// HasFailures returns true if any check failed.
func (r *CIAuditResult) HasFailures() bool {
	return len(r.FailedChecks()) > 0
}

// checkJSON is the JSON shape for a single check.
type checkJSON struct {
	Name    string   `json:"name"`
	Passed  bool     `json:"passed"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

// ToJSON serialises the result to a JSON-compatible map.
func (r *CIAuditResult) ToJSON() map[string]interface{} {
	checks := make([]checkJSON, len(r.Checks))
	passed := 0
	failed := 0
	for i, c := range r.Checks {
		details := c.Details
		if details == nil {
			details = []string{}
		}
		checks[i] = checkJSON{Name: c.Name, Passed: c.Passed, Message: c.Message, Details: details}
		if c.Passed {
			passed++
		} else {
			failed++
		}
	}
	b, _ := json.Marshal(checks)
	var checksSlice []interface{}
	_ = json.Unmarshal(b, &checksSlice)
	return map[string]interface{}{
		"passed": r.Passed(),
		"checks": checksSlice,
		"summary": map[string]interface{}{
			"total":  len(r.Checks),
			"passed": passed,
			"failed": failed,
		},
	}
}

// sarifResult is one SARIF result entry.
type sarifResult struct {
	RuleID  string                 `json:"ruleId"`
	Level   string                 `json:"level"`
	Message map[string]string      `json:"message"`
	Locations []map[string]interface{} `json:"locations"`
}

// ToSARIF serialises the result to SARIF v2.1.0 format for GitHub Code Scanning.
func (r *CIAuditResult) ToSARIF(toolVersion string) map[string]interface{} {
	if toolVersion == "" {
		toolVersion = "0.0.0"
	}

	var results []sarifResult
	var rules []map[string]interface{}

	for _, check := range r.Checks {
		if check.Passed {
			continue
		}
		artifact := ArtifactForCheck(check.Name)
		details := check.Details
		if len(details) == 0 {
			details = []string{check.Message}
		}
		for _, detail := range details {
			results = append(results, sarifResult{
				RuleID:  check.Name,
				Level:   "error",
				Message: map[string]string{"text": detail},
				Locations: []map[string]interface{}{
					{
						"physicalLocation": map[string]interface{}{
							"artifactLocation": map[string]interface{}{
								"uri": artifact,
							},
						},
					},
				},
			})
		}
		rules = append(rules, map[string]interface{}{
			"id": check.Name,
			"shortDescription": map[string]string{"text": check.Message},
		})
	}

	if results == nil {
		results = []sarifResult{}
	}
	if rules == nil {
		rules = []map[string]interface{}{}
	}

	b, _ := json.Marshal(results)
	var resultsSlice []interface{}
	_ = json.Unmarshal(b, &resultsSlice)

	return map[string]interface{}{
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
		"version": "2.1.0",
		"runs": []interface{}{
			map[string]interface{}{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":            "apm-audit",
						"version":         toolVersion,
						"informationUri":  "https://github.com/microsoft/apm",
						"rules":           rules,
					},
				},
				"results": resultsSlice,
			},
		},
	}
}

// RenderSummary returns a human-readable summary of failed checks.
func (r *CIAuditResult) RenderSummary() string {
	if r.Passed() {
		return "[+] All checks passed"
	}
	failed := r.FailedChecks()
	out := fmt.Sprintf("[x] %d check(s) failed:\n", len(failed))
	for _, c := range failed {
		out += fmt.Sprintf("  - %s: %s\n", c.Name, c.Message)
	}
	return out
}
