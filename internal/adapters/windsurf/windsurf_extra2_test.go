package windsurf

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_TargetName_IsWindsurf(t *testing.T) {
	a := New()
	if a.TargetName != "windsurf" {
		t.Errorf("expected TargetName windsurf, got %q", a.TargetName)
	}
}

func TestNew_SupportsRuntimeEnvSubstitution_False(t *testing.T) {
	a := New()
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("expected SupportsRuntimeEnvSubstitution to be false")
	}
}

func TestNew_SupportsUserScope_True(t *testing.T) {
	a := New()
	if !a.SupportsUserScope {
		t.Error("expected SupportsUserScope to be true")
	}
}

func TestGetConfigPath_Separator(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	// path must have at least one separator
	if !strings.ContainsAny(p, string(filepath.Separator)+"/") {
		t.Errorf("expected path separator in %q", p)
	}
}

func TestGetConfigPath_MCPConfigJSON(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	if filepath.Base(p) != "mcp_config.json" {
		t.Errorf("expected filename mcp_config.json, got %q", filepath.Base(p))
	}
}

func TestGetConfigPath_WindsurfDir(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	dir := filepath.Base(filepath.Dir(p))
	if dir != "windsurf" {
		t.Errorf("expected parent dir windsurf, got %q", dir)
	}
}

func TestGetConfigPath_CodeiumGrandparent(t *testing.T) {
	a := New()
	p := a.GetConfigPath()
	parts := strings.Split(filepath.ToSlash(p), "/")
	found := false
	for _, part := range parts {
		if part == ".codeium" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected .codeium in path %q", p)
	}
}

func TestAdapter_ZeroValue_Fields(t *testing.T) {
	var a Adapter
	if a.ClientLabel != "" {
		t.Errorf("expected empty ClientLabel in zero value")
	}
	if a.TargetName != "" {
		t.Errorf("expected empty TargetName in zero value")
	}
	if a.MCPServersKey != "" {
		t.Errorf("expected empty MCPServersKey in zero value")
	}
}

func TestAdapter_SetFields_Roundtrip(t *testing.T) {
	a := &Adapter{
		ClientLabel:   "MyClient",
		TargetName:    "mytarget",
		MCPServersKey: "servers",
	}
	if a.ClientLabel != "MyClient" {
		t.Errorf("unexpected ClientLabel %q", a.ClientLabel)
	}
	if a.TargetName != "mytarget" {
		t.Errorf("unexpected TargetName %q", a.TargetName)
	}
	if a.MCPServersKey != "servers" {
		t.Errorf("unexpected MCPServersKey %q", a.MCPServersKey)
	}
}

func TestIsAvailable_CustomAdapter_AlwaysTrue(t *testing.T) {
	a := &Adapter{}
	if !a.IsAvailable() {
		t.Error("IsAvailable should always return true")
	}
}

func TestGetRuntimeName_CustomTargetName(t *testing.T) {
	a := &Adapter{TargetName: "custom"}
	if a.GetRuntimeName() != "custom" {
		t.Errorf("expected custom, got %q", a.GetRuntimeName())
	}
}
