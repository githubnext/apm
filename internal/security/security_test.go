package security_test

import (
	"testing"

	"github.com/githubnext/apm/internal/security"
)

// TestParityScanPolicyBlock verifies BlockPolicy mirrors Python BLOCK_POLICY.
func TestParityScanPolicyBlock(t *testing.T) {
	p := security.BlockPolicy
	if p.OnCritical != "block" {
		t.Errorf("expected on_critical=block")
	}
	if !p.ForceOverrides {
		t.Errorf("expected force_overrides=true")
	}
}

// TestParityScanPolicyWarn verifies WarnPolicy mirrors Python WARN_POLICY.
func TestParityScanPolicyWarn(t *testing.T) {
	p := security.WarnPolicy
	if p.OnCritical != "warn" {
		t.Errorf("expected on_critical=warn")
	}
	if p.ForceOverrides {
		t.Errorf("expected force_overrides=false")
	}
}

// TestParityScanPolicyReport verifies ReportPolicy mirrors Python REPORT_POLICY.
func TestParityScanPolicyReport(t *testing.T) {
	p := security.ReportPolicy
	if p.OnCritical != "ignore" {
		t.Errorf("expected on_critical=ignore")
	}
}

// TestParityEffectiveBlockNoForce verifies block=true when force=false.
func TestParityEffectiveBlockNoForce(t *testing.T) {
	if !security.BlockPolicy.EffectiveBlock(false) {
		t.Errorf("block policy should block when force=false")
	}
}

// TestParityEffectiveBlockWithForce verifies force overrides block.
func TestParityEffectiveBlockWithForce(t *testing.T) {
	if security.BlockPolicy.EffectiveBlock(true) {
		t.Errorf("block policy with force_overrides=true should not block when force=true")
	}
}

// TestParityEffectiveBlockWarnPolicy verifies warn policy never blocks.
func TestParityEffectiveBlockWarnPolicy(t *testing.T) {
	if security.WarnPolicy.EffectiveBlock(false) {
		t.Errorf("warn policy should not block")
	}
	if security.WarnPolicy.EffectiveBlock(true) {
		t.Errorf("warn policy should not block with force")
	}
}

// TestParityScanVerdictHasFindings verifies finding detection.
func TestParityScanVerdictHasFindings(t *testing.T) {
	v := &security.ScanVerdict{}
	if v.HasFindings() {
		t.Errorf("empty verdict should have no findings")
	}
	v2 := &security.ScanVerdict{
		FindingsByFile: map[string][]security.ScanFinding{
			"file.py": {{Rule: "hardcoded-secret", Severity: "critical"}},
		},
	}
	if !v2.HasFindings() {
		t.Errorf("verdict with findings should return true")
	}
}

// TestParityScanVerdictAllFindings flattens findings across files.
func TestParityScanVerdictAllFindings(t *testing.T) {
	v := &security.ScanVerdict{
		FindingsByFile: map[string][]security.ScanFinding{
			"a.py": {{Rule: "r1", Severity: "critical"}, {Rule: "r2", Severity: "warning"}},
			"b.py": {{Rule: "r3", Severity: "info"}},
		},
	}
	all := v.AllFindings()
	if len(all) != 3 {
		t.Errorf("expected 3 findings, got %d", len(all))
	}
}

// TestParityDefaultContentPatterns verifies built-in patterns are present.
func TestParityDefaultContentPatterns(t *testing.T) {
	patterns := security.DefaultContentPatterns()
	if len(patterns) == 0 {
		t.Errorf("expected at least one default pattern")
	}
	ids := map[string]bool{}
	for _, p := range patterns {
		ids[p.ID] = true
	}
	for _, want := range []string{"hardcoded-secret", "private-key", "aws-key", "github-token"} {
		if !ids[want] {
			t.Errorf("expected pattern %s to be in default patterns", want)
		}
	}
}

// TestParityAuditSeverityValues verifies severity constants.
func TestParityAuditSeverityValues(t *testing.T) {
	cases := map[security.AuditSeverity]string{
		security.SeverityCritical: "critical",
		security.SeverityHigh:     "high",
		security.SeverityMedium:   "medium",
		security.SeverityLow:      "low",
		security.SeverityInfo:     "info",
	}
	for sev, want := range cases {
		if string(sev) != want {
			t.Errorf("expected %s, got %s", want, sev)
		}
	}
}

// TestParityAuditReportCritical filters critical entries.
func TestParityAuditReportCritical(t *testing.T) {
	r := &security.AuditReport{
		Entries: []security.AuditReportEntry{
			{CheckName: "a", Severity: security.SeverityCritical},
			{CheckName: "b", Severity: security.SeverityHigh},
			{CheckName: "c", Severity: security.SeverityInfo},
		},
	}
	crit := r.Critical()
	if len(crit) != 1 {
		t.Errorf("expected 1 critical entry, got %d", len(crit))
	}
	if crit[0].CheckName != "a" {
		t.Errorf("expected check a to be critical")
	}
}

// TestParityAuditReportHasBlockers verifies critical/high detection.
func TestParityAuditReportHasBlockers(t *testing.T) {
	r := &security.AuditReport{}
	if r.HasBlockers() {
		t.Errorf("empty report should have no blockers")
	}
	r2 := &security.AuditReport{
		Entries: []security.AuditReportEntry{
			{CheckName: "a", Severity: security.SeverityHigh},
		},
	}
	if !r2.HasBlockers() {
		t.Errorf("high severity entry should count as blocker")
	}
}
