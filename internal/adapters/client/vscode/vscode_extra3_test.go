package vscode

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTargetName_IsVSCode_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "vscode" {
		t.Errorf("expected vscode, got %q", a.TargetName())
	}
}

func TestMCPServersKey_IsServers_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "servers" {
		t.Errorf("expected servers, got %q", a.MCPServersKey())
	}
}

func TestSupportsUserScope_IsFalse_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.SupportsUserScope() {
		t.Error("expected false")
	}
}

func TestGetConfigPath_ContainsDotVSCode_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(filepath.Dir(p)) != ".vscode" {
		t.Errorf("expected .vscode parent, got %q", filepath.Base(filepath.Dir(p)))
	}
}

func TestGetConfigPath_EndsMCPJson_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	p := a.GetConfigPath()
	if filepath.Base(p) != "mcp.json" {
		t.Errorf("expected mcp.json, got %q", filepath.Base(p))
	}
}

func TestGetConfigPath_IsAbsolute_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	p := a.GetConfigPath()
	if !filepath.IsAbs(p) {
		t.Errorf("expected absolute, got %q", p)
	}
}

func TestGetCurrentConfig_NonNilOnMissing_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map")
	}
}

func TestGetCurrentConfig_ValidJSON_Extra3(t *testing.T) {
	dir := t.TempDir()
	d := filepath.Join(dir, ".vscode")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "mcp.json"), []byte(`{"servers":{}}`), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["servers"]; !ok {
		t.Error("expected servers key")
	}
}

func TestGetCurrentConfig_CorruptJSON_Extra3(t *testing.T) {
	dir := t.TempDir()
	d := filepath.Join(dir, ".vscode")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "mcp.json"), []byte("{bad"), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil on corrupt JSON")
	}
}

func TestTwoInstances_DifferentPaths_Extra3(t *testing.T) {
	d1 := t.TempDir()
	d2 := t.TempDir()
	a1 := New(d1, false)
	a2 := New(d2, false)
	if a1.GetConfigPath() == a2.GetConfigPath() {
		t.Error("expected different paths")
	}
}

func TestSupportsRuntimeEnvSubstitution_False_Extra3(t *testing.T) {
	a := New("/root", false)
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("expected false")
	}
}

func TestUpdateConfig_WithVSCodeDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".vscode"), 0o755)
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"servers": map[string]interface{}{}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTranslateEnvValueForVSCode_EnvPrefix_Extra3(t *testing.T) {
	got := translateEnvValueForVSCode("${env:TOKEN}")
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestTranslateEnvValueForVSCode_NoChange_Extra3(t *testing.T) {
	got := translateEnvValueForVSCode("literal-value")
	if got != "literal-value" {
		t.Errorf("expected unchanged, got %q", got)
	}
}
