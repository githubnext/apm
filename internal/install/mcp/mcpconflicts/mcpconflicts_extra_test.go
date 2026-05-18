package mcpconflicts_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpconflicts"
)

func TestMCPWithVersion(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: true,
		MCPName:    "myserver",
		MCPVersion: "1.2.3",
	}
	if err := mcpconflicts.ValidateMCPConflicts(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMCPWithCommandArgv(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName:  true,
		MCPName:     "myserver",
		CommandArgv: []string{"node", "server.js"},
	}
	if err := mcpconflicts.ValidateMCPConflicts(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMCPGlobalFlag_Fails(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: true,
		MCPName:    "myserver",
		Global:     true,
	}
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --global with --mcp")
	}
}

func TestNoMCPWithRegistryURL_Fails(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName:  false,
		RegistryURL: "https://registry.example.com",
	}
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --registry-url without --mcp")
	}
}

func TestNoMCPWithHeaders_Fails(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: false,
		Headers:    map[string]string{"X-Token": "abc"},
	}
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --header without --mcp")
	}
}

func TestMCPWithRegistryURL(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName:  true,
		MCPName:     "reg-server",
		RegistryURL: "https://registry.example.com",
	}
	if err := mcpconflicts.ValidateMCPConflicts(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConflictConfig_Packages(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: false,
		Packages:   []string{"owner/repo"},
	}
	if len(cfg.Packages) != 1 {
		t.Errorf("Packages length: %d", len(cfg.Packages))
	}
}

func TestMCPWithSSH_Fails(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: true,
		MCPName:    "srv",
		UseSSH:     true,
	}
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --ssh with --mcp")
	}
}

func TestMCPWithHTTPS_Fails(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: true,
		MCPName:    "srv",
		UseHTTPS:   true,
	}
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --https with --mcp")
	}
}

func TestMCPWithUpdate_Fails(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: true,
		MCPName:    "srv",
		Update:     true,
	}
	err := mcpconflicts.ValidateMCPConflicts(cfg)
	if err == nil {
		t.Error("expected error for --update with --mcp")
	}
}

func TestMCPWithOnly(t *testing.T) {
	cfg := mcpconflicts.ConflictConfig{
		HasMCPName: true,
		MCPName:    "srv",
		Only:       "claude",
	}
	if err := mcpconflicts.ValidateMCPConflicts(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
