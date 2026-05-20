package cursor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTargetName_IsCursor_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "cursor" {
		t.Errorf("expected cursor, got %q", a.TargetName())
	}
}

func TestMCPServersKey_IsMcpServers_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %q", a.MCPServersKey())
	}
}

func TestSupportsUserScope_IsFalse_Extra3(t *testing.T) {
	a := New("/tmp", false)
	if a.SupportsUserScope() {
		t.Error("expected false")
	}
}

func TestGetConfigPath_HasMCPJson_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(p) != "mcp.json" {
		t.Errorf("expected mcp.json base, got %q", filepath.Base(p))
	}
}

func TestGetConfigPath_HasVSCodeDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(filepath.Dir(p)) != ".cursor" {
		t.Errorf("expected .cursor parent, got %q", filepath.Base(filepath.Dir(p)))
	}
}

func TestGetCurrentConfig_ReturnsMap_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil")
	}
}

func TestGetCurrentConfig_CorruptJSON_Extra3(t *testing.T) {
	dir := t.TempDir()
	d := filepath.Join(dir, ".cursor")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "mcp.json"), []byte("{bad"), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil on corrupt JSON")
	}
}

func TestUpdateConfig_CreatesFile_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	_ = os.MkdirAll(filepath.Join(dir, ".cursor"), 0o755)
	err := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTwoInstances_DifferentRoots_Extra3(t *testing.T) {
	d1 := t.TempDir()
	d2 := t.TempDir()
	a1 := New(d1, false)
	a2 := New(d2, false)
	if a1.GetConfigPath() == a2.GetConfigPath() {
		t.Error("expected different paths")
	}
}

func TestGetConfigPath_IsAbsolute_Extra3(t *testing.T) {
	a := New(t.TempDir(), false)
	p := a.GetConfigPath()
	if !filepath.IsAbs(p) {
		t.Errorf("expected absolute path, got %q", p)
	}
}

func TestGetCurrentConfig_ValidJSON_Extra3(t *testing.T) {
	dir := t.TempDir()
	d := filepath.Join(dir, ".cursor")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "mcp.json"), []byte(`{"mcpServers":{"s":{"command":"cmd"}}}`), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["mcpServers"]; !ok {
		t.Error("expected mcpServers key")
	}
}

func TestSupportsRuntimeEnvSubstitution_False_Extra3(t *testing.T) {
	a := New("/root", false)
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("expected false")
	}
}
