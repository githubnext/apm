package cursor

import (
	"path/filepath"
	"testing"
)

func TestTargetName(t *testing.T) {
	a := New("/project", false)
	if got := a.TargetName(); got != "cursor" {
		t.Errorf("TargetName() = %q, want %q", got, "cursor")
	}
}

func TestMCPServersKey(t *testing.T) {
	a := New("/project", false)
	if got := a.MCPServersKey(); got != "mcpServers" {
		t.Errorf("MCPServersKey() = %q, want %q", got, "mcpServers")
	}
}

func TestSupportsUserScope(t *testing.T) {
	a := New("/project", false)
	if a.SupportsUserScope() {
		t.Error("SupportsUserScope() = true, want false for cursor")
	}
}

func TestGetConfigPath(t *testing.T) {
	a := New("/myproject", false)
	got := a.GetConfigPath()
	want := filepath.Join("/myproject", ".cursor", "mcp.json")
	if got != want {
		t.Errorf("GetConfigPath() = %q, want %q", got, want)
	}
}

func TestSupportsRuntimeEnvSubstitution(t *testing.T) {
	a := New("/project", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution should be false for cursor")
	}
}
