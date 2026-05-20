package windsurf_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/windsurf"
)

func TestNew_NotNil_Extra3(t *testing.T) {
	a := windsurf.New()
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestTargetName_IsWindsurf_Extra3(t *testing.T) {
	a := windsurf.New()
	if a.TargetName != "windsurf" {
		t.Errorf("expected windsurf, got %q", a.TargetName)
	}
}

func TestMCPServersKey_IsMcpServers_Extra3(t *testing.T) {
	a := windsurf.New()
	if a.MCPServersKey != "mcpServers" {
		t.Errorf("expected mcpServers, got %q", a.MCPServersKey)
	}
}

func TestClientLabel_IsWindsurf_Extra3(t *testing.T) {
	a := windsurf.New()
	if a.ClientLabel == "" {
		t.Error("expected non-empty ClientLabel")
	}
}

func TestSupportsUserScope_True_Extra3(t *testing.T) {
	a := windsurf.New()
	if !a.SupportsUserScope {
		t.Error("expected SupportsUserScope true")
	}
}

func TestSupportsRuntimeEnvSubstitution_False_Extra3(t *testing.T) {
	a := windsurf.New()
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("expected false")
	}
}

func TestGetConfigPath_NotEmpty_Extra3(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty config path")
	}
}

func TestGetConfigPath_ContainsWindsurf_Extra3(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	found := false
	for _, seg := range []string{"windsurf", "codeium"} {
		if len(p) > len(seg) {
			for i := 0; i <= len(p)-len(seg); i++ {
				if p[i:i+len(seg)] == seg {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Errorf("expected windsurf or codeium in path %q", p)
	}
}

func TestGetRuntimeName_IsWindsurf_Extra3(t *testing.T) {
	a := windsurf.New()
	if a.GetRuntimeName() != "windsurf" {
		t.Errorf("expected windsurf, got %q", a.GetRuntimeName())
	}
}

func TestIsAvailable_True_Extra3(t *testing.T) {
	a := windsurf.New()
	if !a.IsAvailable() {
		t.Error("expected IsAvailable true")
	}
}

func TestNew_TwoInstances_SameDefaults_Extra3(t *testing.T) {
	a1 := windsurf.New()
	a2 := windsurf.New()
	if a1.TargetName != a2.TargetName {
		t.Error("expected same TargetName for two instances")
	}
	if a1.MCPServersKey != a2.MCPServersKey {
		t.Error("expected same MCPServersKey for two instances")
	}
}

func TestGetConfigPath_ContainsMCPConfig_Extra3(t *testing.T) {
	a := windsurf.New()
	p := a.GetConfigPath()
	const suffix = "mcp_config.json"
	if len(p) < len(suffix) || p[len(p)-len(suffix):] != suffix {
		t.Errorf("expected mcp_config.json suffix, got %q", p)
	}
}
