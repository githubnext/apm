package codex

import (
	"testing"
)

func TestTargetName(t *testing.T) {
	a := New("/project", false)
	if got := a.TargetName(); got != "codex" {
		t.Errorf("TargetName() = %q, want %q", got, "codex")
	}
}

func TestMCPServersKey(t *testing.T) {
	a := New("/project", false)
	if got := a.MCPServersKey(); got != "mcp_servers" {
		t.Errorf("MCPServersKey() = %q, want %q", got, "mcp_servers")
	}
}

func TestSupportsUserScope(t *testing.T) {
	a := New("/project", false)
	if !a.SupportsUserScope() {
		t.Error("SupportsUserScope() = false, want true")
	}
}

func TestGetConfigPathProjectScope(t *testing.T) {
	a := New("/myproject", false)
	got := a.GetConfigPath()
	if got == "" {
		t.Error("GetConfigPath() returned empty string")
	}
}

func TestGetConfigPathUserScope(t *testing.T) {
	a := New("/myproject", true)
	got := a.GetConfigPath()
	if got == "" {
		t.Error("GetConfigPath() returned empty string for user scope")
	}
}

func TestSupportsRuntimeEnvSubstitution(t *testing.T) {
	a := New("/project", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution should be false for codex")
	}
}
