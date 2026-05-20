package policymodels

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCheckResult_ZeroValue(t *testing.T) {
	var r CheckResult
	if r.Name != "" || r.Passed || r.Message != "" || r.Details != nil {
		t.Errorf("unexpected zero-value: %+v", r)
	}
}

func TestCIAuditResult_Passed_NoChecks(t *testing.T) {
	r := &CIAuditResult{}
	if !r.Passed() {
		t.Error("empty check list should report Passed=true")
	}
}

func TestCIAuditResult_HasFailures_NoChecks(t *testing.T) {
	r := &CIAuditResult{}
	if r.HasFailures() {
		t.Error("empty check list should have no failures")
	}
}

func TestCIAuditResult_FailedChecks_OnlyFailed(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "c1", Passed: true},
		{Name: "c2", Passed: false, Message: "broken"},
		{Name: "c3", Passed: true},
	}}
	failed := r.FailedChecks()
	if len(failed) != 1 || failed[0].Name != "c2" {
		t.Errorf("unexpected failed checks: %v", failed)
	}
}

func TestCIAuditResult_RenderSummary_AllPassed(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{{Name: "c1", Passed: true}}}
	got := r.RenderSummary()
	if !strings.Contains(got, "All checks passed") {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestCIAuditResult_RenderSummary_HasFailures(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: false, Message: "no lockfile"},
	}}
	got := r.RenderSummary()
	if !strings.Contains(got, "1 check(s) failed") {
		t.Errorf("unexpected summary: %q", got)
	}
	if !strings.Contains(got, "lockfile-exists") {
		t.Errorf("check name missing from summary: %q", got)
	}
}

func TestCIAuditResult_ToJSON_RoundTrip(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "c1", Passed: true, Message: "ok", Details: []string{}},
		{Name: "c2", Passed: false, Message: "fail", Details: []string{"detail1"}},
	}}
	m := r.ToJSON()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("ToJSON marshal error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("json roundtrip error: %v", err)
	}
	if passed, ok := out["passed"].(bool); !ok || passed {
		t.Errorf("expected passed=false, got %v", out["passed"])
	}
}

func TestArtifactForCheck_UnknownFallback(t *testing.T) {
	got := ArtifactForCheck("nonexistent-check")
	if got != "apm.lock.yaml" {
		t.Errorf("expected fallback 'apm.lock.yaml', got %q", got)
	}
}

func TestCIAuditResult_ToSARIF_AllPassed(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{{Name: "c1", Passed: true}}}
	sarif := r.ToSARIF("1.0.0")
	if sarif["version"] != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %v", sarif["version"])
	}
	runs, _ := sarif["runs"].([]interface{})
	if len(runs) == 0 {
		t.Error("expected at least one run in SARIF output")
	}
}
