package summary

import (
	"strings"
	"testing"
)

func TestFormatSummary_ZeroAll(t *testing.T) {
	got := FormatSummary(SummaryResult{})
	if !strings.Contains(got, "0 APM package(s)") {
		t.Errorf("unexpected output: %q", got)
	}
	if !strings.Contains(got, "0 MCP server(s)") {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestFormatSummary_WithErrors(t *testing.T) {
	got := FormatSummary(SummaryResult{ApmCount: 1, McpCount: 0, Errors: 3})
	if !strings.Contains(got, "3 error(s)") {
		t.Errorf("expected error count in output: %q", got)
	}
}

func TestFormatSummary_WithStalesCleaned(t *testing.T) {
	got := FormatSummary(SummaryResult{ApmCount: 2, McpCount: 1, StalesCleaned: 5})
	if !strings.Contains(got, "5 stale artifact(s)") {
		t.Errorf("expected stales in output: %q", got)
	}
}

func TestFormatSummary_WithElapsed(t *testing.T) {
	got := FormatSummary(SummaryResult{ApmCount: 0, McpCount: 0, ElapsedSecs: 2.5})
	if !strings.Contains(got, "2.5s") {
		t.Errorf("expected elapsed time in output: %q", got)
	}
}

func TestFormatSummary_NoErrorsWhenZero(t *testing.T) {
	got := FormatSummary(SummaryResult{ApmCount: 1, McpCount: 1, Errors: 0})
	if strings.Contains(got, "error") {
		t.Errorf("should not mention errors when Errors=0: %q", got)
	}
}

func TestFormatSummary_NoStalesWhenZero(t *testing.T) {
	got := FormatSummary(SummaryResult{ApmCount: 1, McpCount: 1, StalesCleaned: 0})
	if strings.Contains(got, "stale") {
		t.Errorf("should not mention stales when StalesCleaned=0: %q", got)
	}
}

func TestFormatSummary_SuffixIsPeriod(t *testing.T) {
	got := FormatSummary(SummaryResult{ApmCount: 2, McpCount: 2})
	if !strings.HasSuffix(got, ".") {
		t.Errorf("summary should end with '.', got: %q", got)
	}
}

func TestHasCriticalSecurityError_ForceSuppresses(t *testing.T) {
	if HasCriticalSecurityError(true, true) {
		t.Error("force=true should suppress critical security error")
	}
}

func TestHasCriticalSecurityError_TrueWhenNotForced(t *testing.T) {
	if !HasCriticalSecurityError(true, false) {
		t.Error("hasCriticalSecurity=true, force=false should return true")
	}
}

func TestHasCriticalSecurityError_FalseWhenNoSecurity(t *testing.T) {
	if HasCriticalSecurityError(false, false) {
		t.Error("hasCriticalSecurity=false should return false regardless of force")
	}
}
