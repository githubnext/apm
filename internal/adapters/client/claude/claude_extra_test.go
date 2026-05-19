package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/claude"
)

func TestGetConfigPath_ProjectScope(t *testing.T) {
	a := claude.New("/my/project", false)
	got := a.GetConfigPath()
	want := filepath.Join("/my/project", ".mcp.json")
	if got != want {
		t.Errorf("GetConfigPath (project) = %q, want %q", got, want)
	}
}

func TestGetConfigPath_UserScope(t *testing.T) {
	a := claude.New("/my/project", true)
	got := a.GetConfigPath()
	if filepath.Base(got) != ".claude.json" {
		t.Errorf("GetConfigPath (user) should end in .claude.json, got %q", got)
	}
}

func TestGetConfigPath_EmptyRootProjectScope(t *testing.T) {
	a := claude.New("", false)
	got := a.GetConfigPath()
	if filepath.Base(got) != ".mcp.json" {
		t.Errorf("expected .mcp.json for empty root, got %q", got)
	}
}

func TestSupportsRuntimeEnvSubstitution_False(t *testing.T) {
	a := claude.New("/tmp", false)
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("Claude adapter should NOT support runtime env substitution")
	}
}

func TestGetCurrentConfig_MissingFile(t *testing.T) {
	a := claude.New("/nonexistent/path/xyz", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should return non-nil map when file is missing")
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty map, got %v", cfg)
	}
}

func TestTargetName_IsLowerCase(t *testing.T) {
	a := claude.New("/tmp", false)
	if a.TargetName() != "claude" {
		t.Errorf("TargetName = %q, want claude", a.TargetName())
	}
}

func TestMCPServersKey_IsCamelCase(t *testing.T) {
	a := claude.New("/tmp", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("MCPServersKey = %q, want mcpServers", a.MCPServersKey())
	}
}

func TestProjectScopeVsUserScope_DifferentPaths(t *testing.T) {
	project := claude.New("/my/project", false)
	user := claude.New("/my/project", true)
	if project.GetConfigPath() == user.GetConfigPath() {
		t.Errorf("project and user scope should produce different paths")
	}
}

func TestSupportsUserScope_True(t *testing.T) {
	a := claude.New("/tmp", false)
	if !a.SupportsUserScope() {
		t.Error("Claude adapter should support user scope")
	}
}

func TestGetCurrentConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(mcpPath, []byte("{not valid json"), 0o644); err != nil {
		t.Skip("cannot write test file")
	}
	a := claude.New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for malformed JSON")
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty map for malformed JSON, got %v", cfg)
	}
}

func TestGetCurrentConfig_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, ".mcp.json")
	data := []byte(`{"mcpServers":{"myserver":{"command":"npx"}}}`)
	if err := os.WriteFile(mcpPath, data, 0o644); err != nil {
		t.Skip("cannot write test file")
	}
	a := claude.New(dir, false)
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["mcpServers"]; !ok {
		t.Errorf("expected mcpServers key, got %v", cfg)
	}
}
