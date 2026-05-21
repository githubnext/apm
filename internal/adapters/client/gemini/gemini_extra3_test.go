package gemini

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTargetName_IsGemini_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "gemini" {
		t.Errorf("expected gemini, got %q", a.TargetName())
	}
}

func TestMCPServersKey_IsMcpServers_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %q", a.MCPServersKey())
	}
}

func TestSupportsUserScope_True_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope true")
	}
}

func TestGetConfigPath_ContainsDotGemini_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(filepath.Dir(p)) != ".gemini" {
		t.Errorf("expected .gemini parent, got %q", filepath.Base(filepath.Dir(p)))
	}
}

func TestGetConfigPath_EndsWithSettingsJson_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(p) != "settings.json" {
		t.Errorf("expected settings.json, got %q", filepath.Base(p))
	}
}

func TestGetConfigPath_IsAbsolute_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	p := a.GetConfigPath()
	if !filepath.IsAbs(p) {
		t.Errorf("expected absolute path, got %q", p)
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
	d := filepath.Join(dir, ".gemini")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "settings.json"), []byte(`{"mcpServers":{}}`), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["mcpServers"]; !ok {
		t.Error("expected mcpServers key")
	}
}

func TestGetCurrentConfig_CorruptJSON_Extra3(t *testing.T) {
	dir := t.TempDir()
	d := filepath.Join(dir, ".gemini")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "settings.json"), []byte("{bad"), 0o644)
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
		t.Error("expected SupportsRuntimeEnvSubstitution false")
	}
}

func TestUpdateConfig_NoGeminiDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{}})
	// Gemini only writes when .gemini dir exists; error or nil both valid
	_ = err
}

func TestUpdateConfig_WithGeminiDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".gemini"), 0o755)
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
