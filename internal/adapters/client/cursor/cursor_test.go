package cursor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetCurrentConfig_Missing(t *testing.T) {
	a := New(t.TempDir(), false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should return empty map, not nil")
	}
	if len(cfg) != 0 {
		t.Errorf("GetCurrentConfig on missing file: want empty, got %v", cfg)
	}
}

func TestGetCurrentConfig_WithFile(t *testing.T) {
	dir := t.TempDir()
	cursorDir := filepath.Join(dir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cursorDir, "mcp.json")
	data := map[string]interface{}{"mcpServers": map[string]interface{}{}}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(cfgPath, b, 0o644); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should return non-nil for existing file")
	}
	if _, ok := cfg["mcpServers"]; !ok {
		t.Error("expected mcpServers key in config")
	}
}

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
