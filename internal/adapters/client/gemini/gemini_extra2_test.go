package gemini

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_SupportsRuntimeEnvSubstitutionFalse(t *testing.T) {
	a := New("/tmp", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("expected SupportsRuntimeEnvSubstitution=false for Gemini adapter")
	}
}

func TestNew_ProjectRootPreserved(t *testing.T) {
	a := New("/my/project", false)
	if a.Adapter.ProjectRoot != "/my/project" {
		t.Errorf("expected ProjectRoot=/my/project, got %q", a.Adapter.ProjectRoot)
	}
}

func TestGetConfigPath_ContainsSettingsJSON(t *testing.T) {
	a := New("/myroot", false)
	p := a.GetConfigPath()
	base := filepath.Base(p)
	if base != "settings.json" {
		t.Errorf("expected settings.json, got %q", base)
	}
}

func TestGetConfigPath_ContainsGeminiDir(t *testing.T) {
	a := New("/myroot", false)
	p := a.GetConfigPath()
	dir := filepath.Base(filepath.Dir(p))
	if dir != ".gemini" {
		t.Errorf("expected .gemini parent dir, got %q", dir)
	}
}

func TestGetCurrentConfig_MissingFile(t *testing.T) {
	a := New("/nonexistent/path/xyz123", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map on missing file")
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty map for missing file, got %v", cfg)
	}
}

func TestGetCurrentConfig_ValidFile(t *testing.T) {
	dir := t.TempDir()
	geminiDir := filepath.Join(dir, ".gemini")
	if err := os.MkdirAll(geminiDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"),
		[]byte(`{"mcpServers":{},"theme":"dark"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["mcpServers"]; !ok {
		t.Error("expected mcpServers key in parsed config")
	}
	if cfg["theme"] != "dark" {
		t.Errorf("expected theme=dark, got %v", cfg["theme"])
	}
}

func TestTargetName_StableRepeated(t *testing.T) {
	a := New("/tmp", false)
	for i := 0; i < 5; i++ {
		if a.TargetName() != "gemini" {
			t.Errorf("TargetName() should always return gemini, got %q", a.TargetName())
		}
	}
}

func TestMCPServersKey_StableRepeated(t *testing.T) {
	a := New("/tmp", false)
	for i := 0; i < 5; i++ {
		if a.MCPServersKey() != "mcpServers" {
			t.Errorf("MCPServersKey() should always return mcpServers, got %q", a.MCPServersKey())
		}
	}
}

func TestSupportsUserScope_True(t *testing.T) {
	a := New("/tmp", true)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope=true")
	}
}
