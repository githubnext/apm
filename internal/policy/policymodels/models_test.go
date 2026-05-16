package policymodels

import (
	"strings"
	"testing"
)

func TestArtifactForCheck(t *testing.T) {
	if ArtifactForCheck("lockfile-exists") != "apm.lock.yaml" {
		t.Error("lockfile-exists should map to apm.lock.yaml")
	}
	if ArtifactForCheck("dependency-allowlist") != "apm.yml" {
		t.Error("dependency-allowlist should map to apm.yml")
	}
	if ArtifactForCheck("unknown-check-xyz") != "apm.lock.yaml" {
		t.Error("unknown check should default to apm.lock.yaml")
	}
}

func TestCIAuditResult_Passed_AllGreen(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: true},
		{Name: "ref-consistency", Passed: true},
	}}
	if !r.Passed() {
		t.Error("expected Passed() = true")
	}
	if r.HasFailures() {
		t.Error("expected HasFailures() = false")
	}
	if len(r.FailedChecks()) != 0 {
		t.Error("expected no failed checks")
	}
}

func TestCIAuditResult_Passed_WithFailure(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: true},
		{Name: "ref-consistency", Passed: false, Message: "mismatch"},
	}}
	if r.Passed() {
		t.Error("expected Passed() = false")
	}
	if !r.HasFailures() {
		t.Error("expected HasFailures() = true")
	}
	failed := r.FailedChecks()
	if len(failed) != 1 || failed[0].Name != "ref-consistency" {
		t.Errorf("FailedChecks() = %v, want one entry ref-consistency", failed)
	}
}

func TestCIAuditResult_ToJSON(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: true, Message: "ok"},
		{Name: "ref-consistency", Passed: false, Message: "bad"},
	}}
	j := r.ToJSON()
	if j["passed"] != false {
		t.Error("ToJSON: passed should be false")
	}
	summary, ok := j["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("ToJSON: no summary map")
	}
	if summary["total"] != 2 {
		t.Errorf("ToJSON: total = %v, want 2", summary["total"])
	}
}

func TestCIAuditResult_RenderSummary_Passed(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{{Name: "lockfile-exists", Passed: true}}}
	out := r.RenderSummary()
	if !strings.Contains(out, "[+]") {
		t.Error("RenderSummary: passed result should contain [+]")
	}
}

func TestCIAuditResult_RenderSummary_Failed(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "ref-consistency", Passed: false, Message: "hash mismatch"},
	}}
	out := r.RenderSummary()
	if !strings.Contains(out, "[x]") {
		t.Error("RenderSummary: failed result should contain [x]")
	}
	if !strings.Contains(out, "ref-consistency") {
		t.Error("RenderSummary: should show failing check name")
	}
}

func TestCIAuditResult_ToSARIF(t *testing.T) {
	r := &CIAuditResult{Checks: []CheckResult{
		{Name: "lockfile-exists", Passed: false, Message: "missing", Details: []string{"detail1"}},
	}}
	sarif := r.ToSARIF("1.0.0")
	if sarif["version"] != "2.1.0" {
		t.Errorf("ToSARIF: version = %v, want 2.1.0", sarif["version"])
	}
	runs, ok := sarif["runs"].([]interface{})
	if !ok || len(runs) == 0 {
		t.Fatal("ToSARIF: no runs")
	}
}
