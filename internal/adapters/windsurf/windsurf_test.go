package windsurf_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/adapters/windsurf"
)

func TestNew_Defaults(t *testing.T) {
	a := windsurf.New()
	if a.ClientLabel != "Windsurf" {
		t.Errorf("ClientLabel = %q, want Windsurf", a.ClientLabel)
	}
	if a.TargetName != "windsurf" {
		t.Errorf("TargetName = %q, want windsurf", a.TargetName)
	}
	if a.MCPServersKey != "mcpServers" {
		t.Errorf("MCPServersKey = %q, want mcpServers", a.MCPServersKey)
	}
	if !a.SupportsUserScope {
		t.Error("SupportsUserScope should be true")
	}
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution should be false")
	}
}

func TestGetConfigPath_ContainsWindsurf(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.Contains(p, "windsurf") {
		t.Errorf("GetConfigPath should contain 'windsurf', got %q", p)
	}
	if !strings.HasSuffix(p, "mcp_config.json") {
		t.Errorf("GetConfigPath should end with mcp_config.json, got %q", p)
	}
}

func TestGetConfigPath_ContainsCodium(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.Contains(p, ".codeium") {
		t.Errorf("expected .codeium in path, got %q", p)
	}
}

func TestGetRuntimeName(t *testing.T) {
	a := windsurf.New()
	if a.GetRuntimeName() != "windsurf" {
		t.Errorf("GetRuntimeName() = %q, want windsurf", a.GetRuntimeName())
	}
}

func TestIsAvailable(t *testing.T) {
	a := windsurf.New()
	if !a.IsAvailable() {
		t.Error("IsAvailable() should return true")
	}
}

func TestAdapter_Fields(t *testing.T) {
	a := windsurf.New()
	if a.ClientLabel == "" {
		t.Error("ClientLabel should not be empty")
	}
	if a.TargetName == "" {
		t.Error("TargetName should not be empty")
	}
	if a.MCPServersKey == "" {
		t.Error("MCPServersKey should not be empty")
	}
}

func TestGetConfigPath_Absolute(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "~") {
		t.Errorf("GetConfigPath() should be absolute or home-relative, got %q", p)
	}
}

func TestGetConfigPath_MCPJson(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, ".json") {
		t.Errorf("GetConfigPath() should end with .json, got %q", p)
	}
}

func TestNew_SupportsUserScope(t *testing.T) {
	a := windsurf.New()
	if !a.SupportsUserScope {
		t.Error("SupportsUserScope should be true for Windsurf global adapter")
	}
}

func TestNew_NoRuntimeEnvSubstitution(t *testing.T) {
	a := windsurf.New()
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution should be false for Windsurf")
	}
}

func TestGetRuntimeName_MatchesTargetName(t *testing.T) {
	a := windsurf.New()
	if a.GetRuntimeName() != a.TargetName {
		t.Errorf("GetRuntimeName() %q != TargetName %q", a.GetRuntimeName(), a.TargetName)
	}
}

func TestNew_MCPServersKeyFormat(t *testing.T) {
	a := windsurf.New()
	if a.MCPServersKey != "mcpServers" {
		t.Errorf("MCPServersKey = %q, want mcpServers", a.MCPServersKey)
	}
}
