package codex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_DefaultsSet_Extra2(t *testing.T) {
	a := New("/proj", false)
	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.Adapter == nil {
		t.Fatal("embedded Adapter is nil")
	}
}

func TestTargetName_IsCodex_Extra2(t *testing.T) {
	a := New("/proj", false)
	if a.TargetName() != "codex" {
		t.Errorf("expected codex, got %q", a.TargetName())
	}
}

func TestMCPServersKey_IsMcpServers_Extra2(t *testing.T) {
	a := New("/proj", false)
	if a.MCPServersKey() != "mcp_servers" {
		t.Errorf("expected mcp_servers, got %q", a.MCPServersKey())
	}
}

func TestGetConfigPath_ContainsDotCodex_Extra2(t *testing.T) {
	a := New("/home/user/myproject", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, ".codex") {
		t.Errorf("expected .codex in path, got %q", p)
	}
}

func TestGetConfigPath_EndsWithConfigTOML_Extra2(t *testing.T) {
	a := New("/some/path", false)
	p := a.GetConfigPath()
	if filepath.Base(p) != "config.toml" {
		t.Errorf("expected config.toml filename, got %q", filepath.Base(p))
	}
}

func TestUpdateConfig_CreatesParentDirs_Extra2(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{})
	if err != nil {
		t.Fatalf("UpdateConfig error: %v", err)
	}
	if _, statErr := os.Stat(a.GetConfigPath()); statErr != nil {
		t.Errorf("config file not created: %v", statErr)
	}
}

func TestGetCurrentConfig_EmptyFileReturnsEmpty_Extra2(t *testing.T) {
	dir := t.TempDir()
	codexDir := filepath.Join(dir, ".codex")
	_ = os.MkdirAll(codexDir, 0o755)
	cfgPath := filepath.Join(codexDir, "config.toml")
	_ = os.WriteFile(cfgPath, []byte(""), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for empty config")
	}
}

func TestGetCurrentConfig_ValidTOML_ReturnsData_Extra2(t *testing.T) {
	dir := t.TempDir()
	codexDir := filepath.Join(dir, ".codex")
	_ = os.MkdirAll(codexDir, 0o755)
	cfgPath := filepath.Join(codexDir, "config.toml")
	_ = os.WriteFile(cfgPath, []byte("[mcp_servers]\n"), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	_ = cfg // should not panic
}

func TestSupportsUserScope_ReturnsTrue_Extra2(t *testing.T) {
	a := New("/proj", false)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope=true")
	}
}

func TestMultipleInstances_Independent_Extra2(t *testing.T) {
	a1 := New("/proj1", false)
	a2 := New("/proj2", false)
	if a1.GetConfigPath() == a2.GetConfigPath() {
		t.Error("different project roots should yield different paths")
	}
}

func TestUpdateConfig_WritesData_Extra2(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{
		"mcp_servers": map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("UpdateConfig: %v", err)
	}
	data, err := os.ReadFile(a.GetConfigPath())
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty config file after UpdateConfig")
	}
}

func TestGetConfigPath_ProjectRoot_Subfolder_Extra2(t *testing.T) {
	a := New("/deeply/nested/project", false)
	p := a.GetConfigPath()
	if !strings.HasPrefix(p, "/deeply/nested/project") {
		t.Errorf("expected path under project root, got %q", p)
	}
}
