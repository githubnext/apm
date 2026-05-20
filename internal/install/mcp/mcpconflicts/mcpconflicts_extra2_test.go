package mcpconflicts

import (
	"strings"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	e := &ValidationError{Message: "test error"}
	if e.Error() != "test error" {
		t.Errorf("unexpected Error() %q", e.Error())
	}
}

func TestValidationError_EmptyMessage(t *testing.T) {
	e := &ValidationError{}
	if e.Error() != "" {
		t.Errorf("expected empty error, got %q", e.Error())
	}
}

func TestConflictConfig_ZeroValue(t *testing.T) {
	var cfg ConflictConfig
	if cfg.HasMCPName {
		t.Error("expected HasMCPName false")
	}
	if cfg.Global {
		t.Error("expected Global false")
	}
}

func TestValidateMCPConflicts_NoFlags_OK(t *testing.T) {
	cfg := ConflictConfig{}
	if err := ValidateMCPConflicts(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMCPConflicts_E10_Transport(t *testing.T) {
	cfg := ConflictConfig{Transport: "sse"}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for transport without --mcp")
	}
	if !strings.Contains(err.Error(), "--transport") {
		t.Errorf("expected --transport in error, got: %v", err)
	}
}

func TestValidateMCPConflicts_E10_URL(t *testing.T) {
	cfg := ConflictConfig{URL: "https://example.com"}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --url without --mcp")
	}
}

func TestValidateMCPConflicts_E10_Env(t *testing.T) {
	cfg := ConflictConfig{Env: map[string]string{"KEY": "val"}}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --env without --mcp")
	}
}

func TestValidateMCPConflicts_E7_EmptyName(t *testing.T) {
	cfg := ConflictConfig{HasMCPName: true, MCPName: ""}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for empty mcp name")
	}
}

func TestValidateMCPConflicts_E8_DashName(t *testing.T) {
	cfg := ConflictConfig{HasMCPName: true, MCPName: "-bad"}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for name starting with -")
	}
}

func TestValidateMCPConflicts_E2_GlobalFails(t *testing.T) {
	cfg := ConflictConfig{HasMCPName: true, MCPName: "myserver", Global: true}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --global with --mcp")
	}
}

func TestValidateMCPConflicts_ValidMCP_OK(t *testing.T) {
	cfg := ConflictConfig{
		HasMCPName:  true,
		MCPName:     "myserver",
		CommandArgv: []string{"node", "server.js"},
	}
	err := ValidateMCPConflicts(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConflictConfig_PreDashPackages(t *testing.T) {
	cfg := ConflictConfig{
		HasMCPName:      true,
		MCPName:         "server",
		PreDashPackages: []string{"pkg1"},
	}
	err := ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for positional packages with --mcp")
	}
}
