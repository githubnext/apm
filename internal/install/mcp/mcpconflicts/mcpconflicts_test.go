package mcpconflicts_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpconflicts"
)

func ok(t *testing.T, cfg mcpconflicts.ConflictConfig) {
	t.Helper()
	if err := mcpconflicts.ValidateMCPConflicts(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func fail(t *testing.T, cfg mcpconflicts.ConflictConfig, substr string) {
	t.Helper()
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Errorf("expected error containing %q, got nil", substr)
		return
	}
	if substr != "" {
		if ve, ok2 := err.(*mcpconflicts.ValidationError); !ok2 {
			t.Errorf("expected *ValidationError, got %T", err)
		} else if len(ve.Message) == 0 {
			t.Error("empty validation message")
		}
	}
}

func TestNoMCPName_NoFlags(t *testing.T) {
	ok(t, mcpconflicts.ConflictConfig{HasMCPName: false})
}

func TestNoMCPName_WithTransport_Fails(t *testing.T) {
	fail(t, mcpconflicts.ConflictConfig{HasMCPName: false, Transport: "stdio"}, "--transport requires --mcp")
}

func TestNoMCPName_WithURL_Fails(t *testing.T) {
	fail(t, mcpconflicts.ConflictConfig{HasMCPName: false, URL: "https://x.com"}, "--url requires --mcp")
}

func TestEmptyMCPName_Fails(t *testing.T) {
	fail(t, mcpconflicts.ConflictConfig{HasMCPName: true, MCPName: ""}, "empty")
}

func TestMCPNameStartsDash_Fails(t *testing.T) {
	fail(t, mcpconflicts.ConflictConfig{HasMCPName: true, MCPName: "-flag"}, "start with '-'")
}

func TestPositionalPackagesMixedWithMCP_Fails(t *testing.T) {
	fail(t, mcpconflicts.ConflictConfig{HasMCPName: true, MCPName: "srv", PreDashPackages: []string{"pkg"}}, "cannot mix")
}

func TestValidMCPWithStdio(t *testing.T) {
	ok(t, mcpconflicts.ConflictConfig{HasMCPName: true, MCPName: "myserver", Transport: "stdio"})
}
