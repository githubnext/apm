package codex

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath_ProjectScope_Extra(t *testing.T) {
	a := New("/my/project", false)
	got := a.GetConfigPath()
	want := filepath.Join("/my/project", ".codex", "config.toml")
	if got != want {
		t.Errorf("GetConfigPath (project) = %q, want %q", got, want)
	}
}

func TestGetConfigPath_UserScope_Extra(t *testing.T) {
	a := New("/my/project", true)
	got := a.GetConfigPath()
	if filepath.Base(got) != "config.toml" {
		t.Errorf("expected config.toml, got %q", got)
	}
	if filepath.Base(filepath.Dir(got)) != ".codex" {
		t.Errorf("expected .codex dir, got %q", filepath.Dir(got))
	}
}

func TestGetConfigPath_EmptyRoot(t *testing.T) {
	a := New("", false)
	got := a.GetConfigPath()
	if filepath.Base(got) != "config.toml" {
		t.Errorf("expected config.toml for empty root, got %q", got)
	}
}

func TestSupportsRuntimeEnvSubstitution_Codex(t *testing.T) {
	a := New("/tmp", false)
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("Codex adapter should NOT support runtime env substitution")
	}
}

func TestGetCurrentConfig_MissingFile_Codex(t *testing.T) {
	a := New("/nonexistent/path/xyz", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should return non-nil map when file is missing")
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty map, got %v", cfg)
	}
}

func TestTargetName_Codex(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "codex" {
		t.Errorf("TargetName = %q, want codex", a.TargetName())
	}
}

func TestMCPServersKey_Codex(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcp_servers" {
		t.Errorf("MCPServersKey = %q, want mcp_servers", a.MCPServersKey())
	}
}

func TestProjectVsUserScope_DifferentPaths_Codex(t *testing.T) {
	project := New("/my/project", false)
	user := New("/my/project", true)
	if project.GetConfigPath() == user.GetConfigPath() {
		t.Errorf("project and user scope should produce different paths")
	}
}

func TestSupportsUserScope_Codex(t *testing.T) {
	a := New("/tmp", false)
	if !a.SupportsUserScope() {
		t.Error("Codex adapter should support user scope")
	}
}

func TestGetCurrentConfig_InvalidTOML(t *testing.T) {
	dir := t.TempDir()
	codexDir := filepath.Join(dir, ".codex")
	_ = os.MkdirAll(codexDir, 0o755)
	cfgPath := filepath.Join(codexDir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte("not valid toml {{{{"), 0o644); err != nil {
		t.Skip("cannot write test file")
	}
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	// Should return nil or empty map, never panic
	_ = cfg
}
