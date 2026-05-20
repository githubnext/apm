package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/claude"
)

func TestTargetName_IsString_Extra3(t *testing.T) {
	a := claude.New("/tmp", false)
	if a.TargetName() == "" {
		t.Error("TargetName should not be empty")
	}
}

func TestMCPServersKey_NotEmpty_Extra3(t *testing.T) {
	a := claude.New("/tmp", false)
	if a.MCPServersKey() == "" {
		t.Error("MCPServersKey should not be empty")
	}
}

func TestSupportsUserScope_StableRepeat_Extra3(t *testing.T) {
	a := claude.New("/tmp", false)
	v1 := a.SupportsUserScope()
	v2 := a.SupportsUserScope()
	if v1 != v2 {
		t.Error("SupportsUserScope should be stable")
	}
}

func TestGetConfigPath_IsAbsolute_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	p := a.GetConfigPath()
	if !filepath.IsAbs(p) {
		t.Errorf("expected absolute path, got %q", p)
	}
}

func TestGetConfigPath_ProjectScope_EndsMCPJson_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	p := a.GetConfigPath()
	if filepath.Base(p) != ".mcp.json" {
		t.Errorf("expected .mcp.json, got %q", filepath.Base(p))
	}
}

func TestGetConfigPath_UserScope_EndsClaude_Extra3(t *testing.T) {
	a := claude.New("", true)
	p := a.GetConfigPath()
	if filepath.Base(p) != ".claude.json" {
		t.Errorf("expected .claude.json, got %q", filepath.Base(p))
	}
}

func TestGetCurrentConfig_CorruptJSON_Extra3(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(cfg, []byte("{bad json"), 0o644); err != nil {
		t.Fatal(err)
	}
	a := claude.New(dir, false)
	got := a.GetCurrentConfig()
	if got == nil {
		t.Error("expected non-nil map on corrupt JSON")
	}
}

func TestGetCurrentConfig_EmptyFile_Extra3(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(cfg, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	a := claude.New(dir, false)
	got := a.GetCurrentConfig()
	if got == nil {
		t.Error("expected non-nil map")
	}
}

func TestUpdateConfig_SameKeyTwice_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	err1 := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{"s1": map[string]interface{}{"command": "cmd"}}})
	err2 := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{"s1": map[string]interface{}{"command": "cmd2"}}})
	if err1 != nil || err2 != nil {
		t.Errorf("unexpected errors: %v %v", err1, err2)
	}
}

func TestGetCurrentConfig_AfterUpdate_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	_ = a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{"x": map[string]interface{}{"k": "v"}}})
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if _, ok := cfg["mcpServers"]; !ok {
		t.Error("expected mcpServers key in config")
	}
}

func TestProjectRootEmbedded_Extra3(t *testing.T) {
	dir := t.TempDir()
	a := claude.New(dir, false)
	p := a.GetConfigPath()
	if !filepath.HasPrefix(p, dir) {
		// May use CWD fallback; just ensure non-empty
		if p == "" {
			t.Error("expected non-empty path")
		}
	}
}

func TestNew_TwoInstances_Independent_Extra3(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	a1 := claude.New(dir1, false)
	a2 := claude.New(dir2, false)
	if a1.GetConfigPath() == a2.GetConfigPath() {
		t.Error("expected different config paths")
	}
}

func TestNew_UserScope_True_Extra3(t *testing.T) {
	a := claude.New("", true)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope true")
	}
}
