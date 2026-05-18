package windsurf_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/adapters/windsurf"
)

func TestNew_Returns_NonNil(t *testing.T) {
	a := windsurf.New()
	if a == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGetConfigPath_NotEmpty(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if p == "" {
		t.Error("GetConfigPath() returned empty string")
	}
}

func TestGetConfigPath_ContainsCodeium(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.Contains(p, ".codeium") {
		t.Errorf("GetConfigPath() should contain .codeium, got %q", p)
	}
}

func TestGetConfigPath_ContainsWindsurfDir(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.Contains(p, "windsurf") {
		t.Errorf("GetConfigPath() should contain windsurf dir, got %q", p)
	}
}

func TestGetConfigPath_EndsMCPConfig(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, "mcp_config.json") {
		t.Errorf("GetConfigPath() should end with mcp_config.json, got %q", p)
	}
}

func TestIsAvailable_AlwaysTrue(t *testing.T) {
	for i := 0; i < 3; i++ {
		a := windsurf.New()
		if !a.IsAvailable() {
			t.Error("IsAvailable() must always return true")
		}
	}
}

func TestGetRuntimeName_IsWindsurf(t *testing.T) {
	a := windsurf.New()
	if a.GetRuntimeName() != "windsurf" {
		t.Errorf("GetRuntimeName() = %q, want windsurf", a.GetRuntimeName())
	}
}

func TestClientLabel_Exact(t *testing.T) {
	a := windsurf.New()
	if a.ClientLabel != "Windsurf" {
		t.Errorf("ClientLabel = %q, want Windsurf", a.ClientLabel)
	}
}

func TestMCPServersKey_CamelCase(t *testing.T) {
	a := windsurf.New()
	if a.MCPServersKey != "mcpServers" {
		t.Errorf("MCPServersKey = %q, want mcpServers", a.MCPServersKey)
	}
}

func TestSupportsUserScope_True(t *testing.T) {
	a := windsurf.New()
	if !a.SupportsUserScope {
		t.Error("SupportsUserScope must be true")
	}
}

func TestSupportsRuntimeEnvSubstitution_False(t *testing.T) {
	a := windsurf.New()
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution must be false")
	}
}

func TestTargetName_Windsurf(t *testing.T) {
	a := windsurf.New()
	if a.TargetName != "windsurf" {
		t.Errorf("TargetName = %q, want windsurf", a.TargetName)
	}
}

func TestGetRuntimeName_ConsistentWithTargetName(t *testing.T) {
	a := windsurf.New()
	if a.GetRuntimeName() != a.TargetName {
		t.Errorf("GetRuntimeName() %q != TargetName %q", a.GetRuntimeName(), a.TargetName)
	}
}

func TestMultipleInstances_Independent(t *testing.T) {
	a1 := windsurf.New()
	a2 := windsurf.New()
	a1.ClientLabel = "Modified"
	if a2.ClientLabel == "Modified" {
		t.Error("instances should be independent")
	}
}
