package cichecks

import (
	"strings"
	"testing"
)

func TestCheckLockfileSync_InSync(t *testing.T) {
	manifest := map[string]bool{"pkg-a": true, "pkg-b": true}
	lockfile := map[string]bool{"pkg-a": true, "pkg-b": true}
	result := CheckLockfileSync(manifest, lockfile)
	if !result.Passed {
		t.Errorf("expected passed, got: %v", result)
	}
}

func TestCheckLockfileSync_MissingFromLockfile(t *testing.T) {
	manifest := map[string]bool{"pkg-a": true, "pkg-b": true, "pkg-c": true}
	lockfile := map[string]bool{"pkg-a": true}
	result := CheckLockfileSync(manifest, lockfile)
	if result.Passed {
		t.Error("expected failure when manifest has entries missing from lockfile")
	}
	if len(result.Details) != 2 {
		t.Errorf("expected 2 missing details, got %d", len(result.Details))
	}
}

func TestCheckLockfileSync_EmptyBoth(t *testing.T) {
	result := CheckLockfileSync(nil, nil)
	if !result.Passed {
		t.Error("expected passed for both empty")
	}
}

func TestCheckLockfileSync_LockfileHasExtra(t *testing.T) {
	// extra entries in lockfile are OK -- only manifest->lockfile direction matters
	manifest := map[string]bool{"pkg-a": true}
	lockfile := map[string]bool{"pkg-a": true, "pkg-b": true}
	result := CheckLockfileSync(manifest, lockfile)
	if !result.Passed {
		t.Error("expected passed when lockfile has extra entries")
	}
}

func TestCheckRefConsistency_AllMatch(t *testing.T) {
	deps := []LockedDepInfo{
		{Key: "pkg-a", ManifestRef: "v1.0", ResolvedRef: "v1.0"},
		{Key: "pkg-b", ManifestRef: "main", ResolvedRef: "main"},
	}
	result := CheckRefConsistency(deps)
	if !result.Passed {
		t.Errorf("expected passed, got: %v", result)
	}
}

func TestCheckRefConsistency_Mismatch(t *testing.T) {
	deps := []LockedDepInfo{
		{Key: "pkg-a", ManifestRef: "v1.0", ResolvedRef: "v2.0"},
	}
	result := CheckRefConsistency(deps)
	if result.Passed {
		t.Error("expected failure for ref mismatch")
	}
	if len(result.Details) != 1 {
		t.Errorf("expected 1 mismatch detail, got %d", len(result.Details))
	}
}

func TestCheckRefConsistency_EmptyRefSkipped(t *testing.T) {
	deps := []LockedDepInfo{
		{Key: "pkg-a", ManifestRef: "", ResolvedRef: "abc123"},
	}
	result := CheckRefConsistency(deps)
	if !result.Passed {
		t.Error("expected passed when ManifestRef is empty (no pinning required)")
	}
}

func TestCheckRefConsistency_Empty(t *testing.T) {
	result := CheckRefConsistency(nil)
	if !result.Passed {
		t.Error("expected passed for empty deps")
	}
}

func TestCIAuditResult_HasFailures_True(t *testing.T) {
	r := CIAuditResult{
		Checks: []CheckResult{
			{Passed: true},
			{Passed: false},
		},
	}
	if !r.HasFailures() {
		t.Error("expected HasFailures to be true")
	}
}

func TestCIAuditResult_HasFailures_False(t *testing.T) {
	r := CIAuditResult{
		Checks: []CheckResult{
			{Passed: true},
			{Passed: true},
		},
	}
	if r.HasFailures() {
		t.Error("expected HasFailures to be false when all passed")
	}
}

func TestCIAuditResult_HasFailures_Empty(t *testing.T) {
	r := CIAuditResult{}
	if r.HasFailures() {
		t.Error("expected HasFailures false for empty audit")
	}
}

func TestCIAuditResult_RenderSummary_PassedAndFailed(t *testing.T) {
	r := CIAuditResult{
		Checks: []CheckResult{
			{Name: "check-a", Passed: true, Message: "all good"},
			{Name: "check-b", Passed: false, Message: "something failed", Details: []string{"detail1"}},
		},
	}
	summary := r.RenderSummary()
	if !strings.Contains(summary, "[+]") {
		t.Error("expected [+] for passed check")
	}
	if !strings.Contains(summary, "[x]") {
		t.Error("expected [x] for failed check")
	}
	if !strings.Contains(summary, "detail1") {
		t.Error("expected detail1 in summary")
	}
}

func TestCIAuditResult_RenderSummary_Empty(t *testing.T) {
	r := CIAuditResult{}
	summary := r.RenderSummary()
	if summary != "" {
		t.Errorf("expected empty summary for empty result, got %q", summary)
	}
}

func TestCheckResult_HasFailures(t *testing.T) {
	r := CheckResult{Passed: false}
	if !r.HasFailures() {
		t.Error("expected HasFailures true for failed check")
	}
	r2 := CheckResult{Passed: true}
	if r2.HasFailures() {
		t.Error("expected HasFailures false for passed check")
	}
}

func TestLockedDepInfo_Fields(t *testing.T) {
	d := LockedDepInfo{
		Key:         "pkg-a",
		ManifestRef: "v1.0",
		ResolvedRef: "v1.0",
	}
	if d.Key != "pkg-a" || d.ManifestRef != "v1.0" {
		t.Errorf("unexpected fields: %+v", d)
	}
}

func TestDriftFinding_Fields(t *testing.T) {
	df := DriftFinding{
		DepKey:   "pkg-a",
		FilePath: "/path/to/file",
		Reason:   "file modified",
	}
	if df.DepKey != "pkg-a" || df.FilePath != "/path/to/file" || df.Reason != "file modified" {
		t.Errorf("unexpected fields: %+v", df)
	}
}
