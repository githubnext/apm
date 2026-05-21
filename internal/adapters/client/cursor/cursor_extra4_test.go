package cursor

import (
	"strings"
	"testing"
)

func TestTargetName_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "cursor" {
		t.Errorf("expected cursor, got %s", a.TargetName())
	}
}

func TestMCPServersKey_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %s", a.MCPServersKey())
	}
}

func TestSupportsUserScope_False_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.SupportsUserScope() {
		t.Error("cursor should not support user scope")
	}
}

func TestGetConfigPath_ContainsCursor_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, ".cursor") {
		t.Errorf("expected .cursor in path, got %s", p)
	}
}

func TestGetConfigPath_EndsWithJSON_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, ".json") {
		t.Errorf("expected .json suffix, got %s", p)
	}
}

func TestGetConfigPath_EmptyRoot_Extra4(t *testing.T) {
	a := New("", false)
	p := a.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetCurrentConfig_Missing_Extra4(t *testing.T) {
	a := New("/nonexistent/xyzabc", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for missing config")
	}
}

func TestNew_NotNil_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestNew_RuntimeEnvSubstitutionFalse_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("cursor should have SupportsRuntimeEnvSubstitution=false")
	}
}

func TestServerKeyFor_Name_Extra4(t *testing.T) {
	k := serverKeyFor("owner/pkg", "myname")
	if k != "myname" {
		t.Errorf("expected myname, got %s", k)
	}
}

func TestServerKeyFor_Slash_Extra4(t *testing.T) {
	k := serverKeyFor("owner/pkg", "")
	if k != "pkg" {
		t.Errorf("expected pkg, got %s", k)
	}
}

func TestServerKeyFor_Simple_Extra4(t *testing.T) {
	k := serverKeyFor("simplepkg", "")
	if k != "simplepkg" {
		t.Errorf("expected simplepkg, got %s", k)
	}
}
