package policymodels

import (
	"strings"
	"testing"
)

func TestArtifactForCheck_AllKnownChecks(t *testing.T) {
	apmLockChecks := []string{
		"lockfile-exists", "ref-consistency", "deployed-files-present",
		"no-orphaned-packages", "config-consistency", "content-integrity",
		"required-packages-deployed", "required-package-version",
		"transitive-depth",
	}
	for _, name := range apmLockChecks {
		if ArtifactForCheck(name) != "apm.lock.yaml" {
			t.Errorf("check %q should map to apm.lock.yaml", name)
		}
	}

	apmYmlChecks := []string{
		"dependency-allowlist", "dependency-denylist", "required-packages",
		"mcp-allowlist", "mcp-denylist", "mcp-transport",
		"mcp-self-defined", "compilation-target", "compilation-strategy",
		"source-attribution", "required-manifest-fields",
		"scripts-policy", "unmanaged-files", "manifest-parse",
	}
	for _, name := range apmYmlChecks {
		if ArtifactForCheck(name) != "apm.yml" {
			t.Errorf("check %q should map to apm.yml", name)
		}
	}
}

func TestCIAuditResult_EmptyChecks(t *testing.T) {
	r := &CIAuditResult{}
	if !r.Passed() {
		t.Error("empty checks should count as passed")
	}
	if r.HasFailures() {
		t.Error("empty checks should have no failures")
	}
	if len(r.FailedChecks()) != 0 {
		t.Error("empty checks: FailedChecks should be empty")
	}
}

func TestCIAuditResult_MultipleFailures(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "a", Passed: false},
		{Name: "b", Passed: true},
		{Name: "c", Passed: false},
	}}
	if r.Passed() {
		t.Error("expected Passed() = false")
	}
	failed := r.FailedChecks()
	if len(failed) != 2 {
		t.Errorf("expected 2 failures, got %d", len(failed))
	}
}

func TestCIAuditResult_ToJSON_AllPassed(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: true},
	}}
	j := r.ToJSON()
	if j["passed"] != true {
		t.Error("expected passed=true")
	}
	summary := j["summary"].(map[string]interface{})
	if summary["failed"] != 0 {
		t.Errorf("expected failed=0, got %v", summary["failed"])
	}
}

func TestCIAuditResult_ToJSON_SummaryFields(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "a", Passed: true},
		{Name: "b", Passed: false},
		{Name: "c", Passed: false},
	}}
	j := r.ToJSON()
	summary := j["summary"].(map[string]interface{})
	if summary["total"] != 3 {
		t.Errorf("total=%v, want 3", summary["total"])
	}
	if summary["passed"] != 1 {
		t.Errorf("passed=%v, want 1", summary["passed"])
	}
	if summary["failed"] != 2 {
		t.Errorf("failed=%v, want 2", summary["failed"])
	}
}

func TestCIAuditResult_RenderSummary_MultipleChecks(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: true},
		{Name: "ref-consistency", Passed: false, Message: "hash mismatch"},
	}}
	out := r.RenderSummary()
	// RenderSummary only lists failing checks
	if !strings.Contains(out, "ref-consistency") {
		t.Error("should contain failing check name")
	}
	if !strings.Contains(out, "[x]") {
		t.Error("should contain failure marker")
	}
}

func TestCIAuditResult_ToSARIF_OnlyFailuresInResults(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: true},
		{Name: "ref-consistency", Passed: false, Message: "bad", Details: []string{"detail"}},
	}}
	sarif := r.ToSARIF("2.0.0")
	runs := sarif["runs"].([]interface{})
	run := runs[0].(map[string]interface{})
	results := run["results"].([]interface{})
	if len(results) != 1 {
		t.Errorf("expected 1 SARIF result (only failures), got %d", len(results))
	}
}

func TestCIAuditResult_ToSARIF_EmptyVersion(t *testing.T) {
	r := &CIAuditResult{}
	sarif := r.ToSARIF("")
	runs := sarif["runs"].([]interface{})
	run := runs[0].(map[string]interface{})
	tool := run["tool"].(map[string]interface{})
	driver := tool["driver"].(map[string]interface{})
	if driver["version"] != "0.0.0" {
		t.Errorf("expected default version 0.0.0, got %v", driver["version"])
	}
}

func TestCheckResult_DetailsNilSafe(t *testing.T) {
	c := CheckResult{Name: "x", Passed: false, Details: nil}
	r := &CIAuditResult{Checks: []CheckResult{c}}
	// ToJSON should not panic with nil details
	j := r.ToJSON()
	if j == nil {
		t.Error("expected non-nil ToJSON result")
	}
}
