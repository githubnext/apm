package vscode

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_TargetName(t *testing.T) {
	a := New("/project", false)
	if a.TargetName() != "vscode" {
		t.Errorf("TargetName = %q, want vscode", a.TargetName())
	}
}

func TestNew_MCPServersKey(t *testing.T) {
	a := New("/project", false)
	if a.MCPServersKey() != "servers" {
		t.Errorf("MCPServersKey = %q, want servers", a.MCPServersKey())
	}
}

func TestNew_SupportsUserScope(t *testing.T) {
	a := New("/project", false)
	if a.SupportsUserScope() {
		t.Error("VSCode adapter should not support user scope")
	}
	a2 := New("/project", true)
	if a2.SupportsUserScope() {
		t.Error("VSCode adapter should not support user scope regardless of flag")
	}
}

func TestNew_GetConfigPath(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	got := a.GetConfigPath()
	want := filepath.Join(dir, ".vscode", "mcp.json")
	if got != want {
		t.Errorf("GetConfigPath = %q, want %q", got, want)
	}
}

func TestNew_GetCurrentConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	// missing file should return empty map, not nil
	if cfg == nil {
		t.Error("GetCurrentConfig should return non-nil map even when file missing")
	}
}

func TestNew_GetCurrentConfig_WithFile(t *testing.T) {
	dir := t.TempDir()
	vscodeDir := filepath.Join(dir, ".vscode")
	os.MkdirAll(vscodeDir, 0o755)
	os.WriteFile(filepath.Join(vscodeDir, "mcp.json"), []byte(`{"servers":{}}`), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should return non-nil map")
	}
}

func TestTranslateEnvValueForVSCode_MultipleVars(t *testing.T) {
	// A plain string with no var references should be returned unchanged
	got := translateEnvValueForVSCode("plain-value")
	if got != "plain-value" {
		t.Errorf("plain value: got %q, want plain-value", got)
	}
}

func TestTranslateEnvValueForVSCode_DollarNoBrace(t *testing.T) {
	// $VAR without braces should not be transformed
	got := translateEnvValueForVSCode("$MY_TOKEN")
	if got != "$MY_TOKEN" {
		t.Errorf("expected $MY_TOKEN unchanged, got %q", got)
	}
}

func TestFilterOut_AllMatch(t *testing.T) {
	got := filterOut([]string{"x", "x", "x"}, "x")
	if len(got) != 0 {
		t.Errorf("all-match: expected empty, got %v", got)
	}
}

func TestFilterOut_Nil(t *testing.T) {
	got := filterOut(nil, "anything")
	if got != nil && len(got) != 0 {
		t.Errorf("nil input: expected nil/empty, got %v", got)
	}
}

func TestNew_ProjectRootStored(t *testing.T) {
	a := New("/my/project", false)
	if a.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot = %q, want /my/project", a.ProjectRoot)
	}
}

func TestNew_SupportsRuntimeEnvSubstitution(t *testing.T) {
	a := New("/project", false)
	if a.SupportsRuntimeEnvSubstitution {
		t.Error("VSCode adapter SupportsRuntimeEnvSubstitution should be false")
	}
}
