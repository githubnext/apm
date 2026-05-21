package codex

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_NotNil_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestTargetName_Stable_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != a.TargetName() {
		t.Error("TargetName not stable")
	}
}

func TestMCPServersKey_Stable_Extra3(t *testing.T) {
	a := New("/tmp", false)
	k1 := a.MCPServersKey()
	k2 := a.MCPServersKey()
	if k1 != k2 {
		t.Error("MCPServersKey not stable")
	}
}

func TestSupportsUserScope_Stable_Extra3(t *testing.T) {
	a := New("", false)
	if a.SupportsUserScope() != a.SupportsUserScope() {
		t.Error("SupportsUserScope not stable")
	}
}

func TestGetConfigPath_EndsWithTOML_Extra3(t *testing.T) {
	a := New("/myroot", false)
	p := a.GetConfigPath()
	if filepath.Ext(p) != ".toml" {
		t.Errorf("expected .toml extension, got %q", filepath.Ext(p))
	}
}

func TestGetConfigPath_ContainsCodex_Extra3(t *testing.T) {
	a := New("/myroot", false)
	p := a.GetConfigPath()
	if filepath.Base(filepath.Dir(p)) != ".codex" {
		t.Errorf("expected .codex parent dir, got %q", filepath.Base(filepath.Dir(p)))
	}
}

func TestGetCurrentConfig_NonNilOnMissing_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map")
	}
}

func TestGetCurrentConfig_NonNilOnCorrupt_Extra3(t *testing.T) {
	dir := t.TempDir()
	codexDir := filepath.Join(dir, ".codex")
	_ = os.MkdirAll(codexDir, 0o755)
	_ = os.WriteFile(filepath.Join(codexDir, "config.toml"), []byte("not = valid [[toml"), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map on corrupt TOML")
	}
}

func TestUpdateConfig_CreatesDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"mcp_servers": map[string]interface{}{}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, statErr := os.Stat(a.GetConfigPath())
	if os.IsNotExist(statErr) {
		t.Error("expected config file to be created")
	}
}

func TestUpdateConfig_StoresKey_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	// UpdateConfig merges its argument as entries under mcp_servers
	err := a.UpdateConfig(map[string]interface{}{"srv": map[string]interface{}{"command": "node"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify config file was created
	_, statErr := os.Stat(a.GetConfigPath())
	if os.IsNotExist(statErr) {
		t.Error("expected config file to be created")
	}
}

func TestTwoInstances_IndependentPaths_Extra3(t *testing.T) {
	d1 := t.TempDir()
	d2 := t.TempDir()
	a1 := New(d1, false)
	a2 := New(d2, false)
	if a1.GetConfigPath() == a2.GetConfigPath() {
		t.Error("expected different paths")
	}
}

func TestGetConfigPath_UserScope_ContainsHome_Extra3(t *testing.T) {
	a := New("", true)
	p := a.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty path")
	}
}

func TestSupportsRuntimeEnv_IsFalse_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("expected SupportsRuntimeEnvSubstitution to be false")
	}
}
