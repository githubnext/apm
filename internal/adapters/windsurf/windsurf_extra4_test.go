package windsurf

import (
	"strings"
	"testing"
)

func TestNew_NotNil_Extra4(t *testing.T) {
	a := New()
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestGetRuntimeName_Extra4(t *testing.T) {
	a := New()
	if a.GetRuntimeName() != "windsurf" {
		t.Errorf("expected windsurf, got %s", a.GetRuntimeName())
	}
}

func TestIsAvailable_Extra4(t *testing.T) {
	a := New()
	if !a.IsAvailable() {
		t.Error("windsurf should always be available")
	}
}

func TestGetConfigPath_ContainsWindsurf_Extra4(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	if !strings.Contains(p, "windsurf") {
		t.Errorf("expected windsurf in path, got %s", p)
	}
}

func TestGetConfigPath_EndsWithJSON_Extra4(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, ".json") {
		t.Errorf("expected .json suffix, got %s", p)
	}
}

func TestGetConfigPath_ContainsCodium_Extra4(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	if !strings.Contains(p, "codeium") {
		t.Errorf("expected codeium in path, got %s", p)
	}
}

func TestTargetName_Extra4(t *testing.T) {
	a := New()
	if a.TargetName != "windsurf" {
		t.Errorf("expected windsurf, got %s", a.TargetName)
	}
}

func TestMCPServersKey_Extra4(t *testing.T) {
	a := New()
	if a.MCPServersKey != "mcpServers" {
		t.Errorf("expected mcpServers, got %s", a.MCPServersKey)
	}
}

func TestSupportsUserScope_Extra4(t *testing.T) {
	a := New()
	if !a.SupportsUserScope {
		t.Error("expected windsurf to support user scope")
	}
}

func TestClientLabel_Extra4(t *testing.T) {
	a := New()
	if a.ClientLabel == "" {
		t.Error("expected non-empty ClientLabel")
	}
}

func TestRuntimeEnvSubstitution_False_Extra4(t *testing.T) {
	a := New()
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("expected SupportsRuntimeEnvSubstitution=false")
	}
}
