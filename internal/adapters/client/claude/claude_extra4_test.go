package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/claude"
)

func TestTargetName_Extra4(t *testing.T) {
	a := claude.New("", false)
	if a.TargetName() != "claude" {
		t.Fatalf("expected claude, got %s", a.TargetName())
	}
}

func TestMCPServersKey_Extra4(t *testing.T) {
	a := claude.New("", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Fatalf("expected mcpServers, got %s", a.MCPServersKey())
	}
}

func TestSupportsUserScope_Extra4(t *testing.T) {
	a := claude.New("", false)
	if !a.SupportsUserScope() {
		t.Fatal("expected SupportsUserScope to be true")
	}
}

func TestGetConfigPath_ProjectScope_Extra4(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(p) != ".mcp.json" {
		t.Fatalf("expected .mcp.json, got %s", filepath.Base(p))
	}
}

func TestGetConfigPath_UserScope_Extra4(t *testing.T) {
	a := claude.New("", true)
	p := a.GetConfigPath()
	if filepath.Base(p) != ".claude.json" {
		t.Fatalf("expected .claude.json, got %s", filepath.Base(p))
	}
}

func TestGetCurrentConfig_MissingFile_Extra4(t *testing.T) {
	a := claude.New("/nonexistent/path/xyz", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Fatal("expected non-nil empty map")
	}
	if len(cfg) != 0 {
		t.Fatalf("expected empty map, got %v", cfg)
	}
}

func TestUpdateConfig_CreatesFile_Extra4(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"myServer": map[string]interface{}{"url": "http://test"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := a.GetConfigPath()
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("config file not created: %v", err)
	}
}

func TestUpdateConfig_RoundTrip_Extra4(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"srv1": map[string]interface{}{"url": "http://srv1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["mcpServers"]; !ok {
		t.Fatal("expected mcpServers key in config")
	}
}

func TestNew_UserScope_SetsUserScope_Extra4(t *testing.T) {
	a := claude.New("", true)
	if !a.UserScope {
		t.Fatal("expected UserScope to be true")
	}
}

func TestNew_ProjectScope_SetsUserScope_Extra4(t *testing.T) {
	a := claude.New("", false)
	if a.UserScope {
		t.Fatal("expected UserScope to be false")
	}
}
