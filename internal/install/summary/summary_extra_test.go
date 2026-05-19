package summary

import (
	"strings"
	"testing"
)

func TestFormatSummary_ExactFormat(t *testing.T) {
	r := SummaryResult{ApmCount: 5, McpCount: 3}
	got := FormatSummary(r)
	want := "Installed 5 APM package(s), 3 MCP server(s)."
	if got != want {
		t.Errorf("FormatSummary(%+v) = %q, want %q", r, got, want)
	}
}

func TestFormatSummary_WithErrorsOnly(t *testing.T) {
	r := SummaryResult{ApmCount: 0, McpCount: 0, Errors: 1}
	got := FormatSummary(r)
	if !strings.Contains(got, "1 error(s)") {
		t.Errorf("expected error count in output: %q", got)
	}
}

func TestFormatSummary_ElapsedPrecision(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 0, ElapsedSecs: 1.234}
	got := FormatSummary(r)
	// %.1f rounds 1.234 to "1.2"
	if !strings.Contains(got, "1.2s") {
		t.Errorf("expected 1.2s in output: %q", got)
	}
}

func TestFormatSummary_ElapsedSmall(t *testing.T) {
	r := SummaryResult{ApmCount: 0, McpCount: 0, ElapsedSecs: 0.5}
	got := FormatSummary(r)
	if !strings.Contains(got, "0.5s") {
		t.Errorf("expected 0.5s in output: %q", got)
	}
}

func TestFormatSummary_StalesPluralLabel(t *testing.T) {
	r := SummaryResult{ApmCount: 0, McpCount: 0, StalesCleaned: 1}
	got := FormatSummary(r)
	if !strings.Contains(got, "cleaned 1 stale artifact(s)") {
		t.Errorf("expected stale label in output: %q", got)
	}
}

func TestFormatSummary_ZeroElapsed_NoTimeClause(t *testing.T) {
	r := SummaryResult{ApmCount: 2, McpCount: 1, ElapsedSecs: 0.0}
	got := FormatSummary(r)
	if strings.Contains(got, "in") {
		t.Errorf("zero elapsed should omit 'in ...' clause: %q", got)
	}
}

func TestFormatSummary_MultipleFields_Order(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 1, Errors: 1, StalesCleaned: 2, ElapsedSecs: 5.0}
	got := FormatSummary(r)
	// errors should come before stales
	errIdx := strings.Index(got, "error")
	staleIdx := strings.Index(got, "stale")
	if errIdx == -1 || staleIdx == -1 {
		t.Errorf("both errors and stale should be in output: %q", got)
	}
	if errIdx > staleIdx {
		t.Errorf("errors should appear before stales in output: %q", got)
	}
}

func TestHasCriticalSecurityError_Matrix(t *testing.T) {
	cases := []struct {
		critical bool
		force    bool
		want     bool
	}{
		{true, false, true},
		{true, true, false},
		{false, false, false},
		{false, true, false},
	}
	for _, tc := range cases {
		got := HasCriticalSecurityError(tc.critical, tc.force)
		if got != tc.want {
			t.Errorf("HasCriticalSecurityError(%v, %v) = %v, want %v", tc.critical, tc.force, got, tc.want)
		}
	}
}

func TestFormatSummary_NegativeElapsed_Omitted(t *testing.T) {
	// Negative elapsed should be treated as non-positive and omitted
	r := SummaryResult{ApmCount: 1, McpCount: 0, ElapsedSecs: -1.0}
	got := FormatSummary(r)
	if strings.Contains(got, "in") {
		t.Errorf("negative elapsed should not add time clause: %q", got)
	}
}

func TestFormatSummary_LargeValues(t *testing.T) {
	r := SummaryResult{ApmCount: 1000, McpCount: 999, Errors: 50, StalesCleaned: 200, ElapsedSecs: 3600.0}
	got := FormatSummary(r)
	if !strings.Contains(got, "1000 APM package(s)") {
		t.Errorf("expected 1000 packages: %q", got)
	}
	if !strings.Contains(got, "999 MCP server(s)") {
		t.Errorf("expected 999 servers: %q", got)
	}
	if !strings.Contains(got, "3600.0s") {
		t.Errorf("expected 3600.0s: %q", got)
	}
}

func TestSummaryResult_ZeroValue(t *testing.T) {
	var r SummaryResult
	got := FormatSummary(r)
	if !strings.HasSuffix(got, ".") {
		t.Errorf("zero-value result should still end with period: %q", got)
	}
}
