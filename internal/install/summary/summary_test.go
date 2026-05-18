package summary

import (
	"strings"
	"testing"
)

func TestFormatSummary_basic(t *testing.T) {
	r := SummaryResult{ApmCount: 3, McpCount: 2}
	got := FormatSummary(r)
	if !strings.Contains(got, "3 APM package(s)") {
		t.Errorf("unexpected output: %q", got)
	}
	if !strings.Contains(got, "2 MCP server(s)") {
		t.Errorf("unexpected output: %q", got)
	}
	if !strings.HasSuffix(got, ".") {
		t.Errorf("expected trailing period: %q", got)
	}
}

func TestFormatSummary_withErrors(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 0, Errors: 2}
	got := FormatSummary(r)
	if !strings.Contains(got, "2 error(s)") {
		t.Errorf("expected errors in output: %q", got)
	}
}

func TestFormatSummary_withStales(t *testing.T) {
	r := SummaryResult{ApmCount: 0, McpCount: 0, StalesCleaned: 5}
	got := FormatSummary(r)
	if !strings.Contains(got, "5 stale artifact(s)") {
		t.Errorf("expected stales in output: %q", got)
	}
}

func TestFormatSummary_withElapsed(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 1, ElapsedSecs: 3.14}
	got := FormatSummary(r)
	if !strings.Contains(got, "3.1s") {
		t.Errorf("expected elapsed time: %q", got)
	}
}

func TestFormatSummary_noElapsed(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 0, ElapsedSecs: 0}
	got := FormatSummary(r)
	if strings.Contains(got, "in 0") {
		t.Errorf("should not include zero elapsed: %q", got)
	}
}

func TestHasCriticalSecurityError(t *testing.T) {
	if !HasCriticalSecurityError(true, false) {
		t.Error("expected true: critical=true, force=false")
	}
	if HasCriticalSecurityError(true, true) {
		t.Error("expected false: critical=true, force=true")
	}
	if HasCriticalSecurityError(false, false) {
		t.Error("expected false: critical=false, force=false")
	}
	if HasCriticalSecurityError(false, true) {
		t.Error("expected false: critical=false, force=true")
	}
}

func TestFormatSummary_Zero(t *testing.T) {
	r := SummaryResult{}
	got := FormatSummary(r)
	if !strings.Contains(got, "0 APM package(s)") {
		t.Errorf("expected 0 APM packages, got %q", got)
	}
	if !strings.Contains(got, "0 MCP server(s)") {
		t.Errorf("expected 0 MCP servers, got %q", got)
	}
}

func TestFormatSummary_AllFields(t *testing.T) {
	r := SummaryResult{ApmCount: 2, McpCount: 3, Errors: 1, StalesCleaned: 4, ElapsedSecs: 10.5}
	got := FormatSummary(r)
	if !strings.Contains(got, "2 APM package(s)") {
		t.Errorf("expected APM count, got %q", got)
	}
	if !strings.Contains(got, "3 MCP server(s)") {
		t.Errorf("expected MCP count, got %q", got)
	}
	if !strings.Contains(got, "1 error(s)") {
		t.Errorf("expected errors, got %q", got)
	}
	if !strings.Contains(got, "4 stale artifact(s)") {
		t.Errorf("expected stales, got %q", got)
	}
	if !strings.Contains(got, "10.5s") {
		t.Errorf("expected elapsed, got %q", got)
	}
}

func TestFormatSummary_NoErrors(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 1, Errors: 0}
	got := FormatSummary(r)
	if strings.Contains(got, "error") {
		t.Errorf("should not contain error when Errors=0: %q", got)
	}
}

func TestFormatSummary_NoStales(t *testing.T) {
	r := SummaryResult{ApmCount: 1, McpCount: 0, StalesCleaned: 0}
	got := FormatSummary(r)
	if strings.Contains(got, "stale") {
		t.Errorf("should not contain stale when StalesCleaned=0: %q", got)
	}
}

func TestFormatSummary_EndsWithPeriod(t *testing.T) {
	cases := []SummaryResult{
		{},
		{ApmCount: 1, McpCount: 2, Errors: 3, StalesCleaned: 4, ElapsedSecs: 5.0},
	}
	for _, r := range cases {
		got := FormatSummary(r)
		if !strings.HasSuffix(got, ".") {
			t.Errorf("FormatSummary should end with period: %q", got)
		}
	}
}
