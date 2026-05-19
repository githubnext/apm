package cursor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_SetsProjectRoot(t *testing.T) {
	a := New("/my/project", false)
	if a.Adapter == nil {
		t.Fatal("Adapter should not be nil")
	}
	if a.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot = %q, want /my/project", a.ProjectRoot)
	}
}

func TestNew_SupportsRuntimeEnvSubstitutionFalse(t *testing.T) {
	a := New("/project", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution should be false for cursor")
	}
}

func TestTargetName_ReturnsConstant(t *testing.T) {
	cases := []string{"/proj", "/other/root", ""}
	for _, root := range cases {
		a := New(root, false)
		if a.TargetName() != "cursor" {
			t.Errorf("TargetName() for root=%q = %q, want cursor", root, a.TargetName())
		}
	}
}

func TestMCPServersKey_ReturnsConstant(t *testing.T) {
	a := New("/project", true)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("MCPServersKey() = %q, want mcpServers", a.MCPServersKey())
	}
}

func TestSupportsUserScope_AlwaysFalse(t *testing.T) {
	for _, userScope := range []bool{true, false} {
		a := New("/proj", userScope)
		if a.SupportsUserScope() {
			t.Errorf("SupportsUserScope() should be false (userScope=%v)", userScope)
		}
	}
}

func TestGetConfigPath_Structure(t *testing.T) {
	a := New("/workspace/proj", false)
	got := a.GetConfigPath()
	if !filepath.IsAbs(got) {
		t.Errorf("GetConfigPath should be absolute: %q", got)
	}
	base := filepath.Base(got)
	if base != "mcp.json" {
		t.Errorf("config file name = %q, want mcp.json", base)
	}
	dir := filepath.Dir(got)
	if filepath.Base(dir) != ".cursor" {
		t.Errorf("parent dir = %q, want .cursor", filepath.Base(dir))
	}
}

func TestGetConfigPath_ContainsProjectRoot(t *testing.T) {
	a := New("/home/user/my-project", false)
	got := a.GetConfigPath()
	if !filepath.HasPrefix(got, "/home/user/my-project") {
		t.Errorf("config path should be under project root: %q", got)
	}
}

func TestGetCurrentConfig_ReturnsEmptyMapNotNil(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should never return nil")
	}
}

func TestGetCurrentConfig_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	cursorDir := filepath.Join(dir, ".cursor")
	os.MkdirAll(cursorDir, 0o755)
	cfgPath := filepath.Join(cursorDir, "mcp.json")
	os.WriteFile(cfgPath, []byte(`{"mcpServers":{"my-server":{"command":"npx"}}}`), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if _, ok := cfg["mcpServers"]; !ok {
		t.Error("expected mcpServers key")
	}
}

func TestUpdateConfig_CursorDirMustExist(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	// No .cursor dir: UpdateConfig should be a no-op (not an error)
	if err := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{}}); err != nil {
		t.Errorf("expected no error when .cursor dir is absent: %v", err)
	}
}

func TestUpdateConfig_CreatesMCPJSON(t *testing.T) {
	dir := t.TempDir()
	cursorDir := filepath.Join(dir, ".cursor")
	os.MkdirAll(cursorDir, 0o755)
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{}})
	if err != nil {
		t.Fatalf("UpdateConfig: %v", err)
	}
	cfgPath := filepath.Join(cursorDir, "mcp.json")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("mcp.json should have been created")
	}
}

func TestNew_MultipleInstances_Independent(t *testing.T) {
	a1 := New("/proj1", false)
	a2 := New("/proj2", false)
	if a1.ProjectRoot == a2.ProjectRoot {
		t.Error("distinct adapters should have distinct ProjectRoot values")
	}
	if a1.GetConfigPath() == a2.GetConfigPath() {
		t.Error("distinct adapters should have distinct config paths")
	}
}
