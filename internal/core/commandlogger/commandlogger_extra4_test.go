package commandlogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/commandlogger"
)

func TestCommandLogger_Verbose_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("test", true, false)
	if !l.Verbose {
		t.Fatal("expected Verbose=true")
	}
}

func TestCommandLogger_DryRun_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("cmd", false, true)
	if !l.DryRun {
		t.Fatal("expected DryRun=true")
	}
}

func TestShouldExecute_DryRunFalse_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("cmd", false, false)
	if !l.ShouldExecute() {
		t.Fatal("expected ShouldExecute=true when DryRun=false")
	}
}

func TestShouldExecute_DryRunTrue_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("cmd", false, true)
	if l.ShouldExecute() {
		t.Fatal("expected ShouldExecute=false when DryRun=true")
	}
}

func TestStripSourcePrefix_Empty_Extra4(t *testing.T) {
	if commandlogger.StripSourcePrefix("") != "" {
		t.Fatal("expected empty string")
	}
}

func TestStripSourcePrefix_Org_Extra4(t *testing.T) {
	s := commandlogger.StripSourcePrefix("org:myorg")
	if s != "myorg" {
		t.Fatalf("expected myorg, got %s", s)
	}
}

func TestStripSourcePrefix_URL_Extra4(t *testing.T) {
	s := commandlogger.StripSourcePrefix("url:https://example.com/policy")
	if s != "https://example.com/policy" {
		t.Fatalf("unexpected: %s", s)
	}
}

func TestStripSourcePrefix_NoPrefix_Extra4(t *testing.T) {
	s := commandlogger.StripSourcePrefix("bare-value")
	if s != "bare-value" {
		t.Fatalf("expected bare-value, got %s", s)
	}
}

func TestCommandLogger_Fields_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", true, true)
	if l.Command != "install" || !l.Verbose || !l.DryRun {
		t.Fatal("field mismatch")
	}
}

func TestMCPLookupHeartbeat_Zero_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("mcp", false, false)
	l.MCPLookupHeartbeat(0) // must not panic
}

func TestMCPLookupHeartbeat_One_Extra4(t *testing.T) {
	l := commandlogger.NewCommandLogger("mcp", false, false)
	l.MCPLookupHeartbeat(1) // must not panic
}
