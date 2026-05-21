package gemini

import (
	"strings"
	"testing"
)

func TestTargetName_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "gemini" {
		t.Errorf("expected gemini, got %s", a.TargetName())
	}
}

func TestMCPServersKey_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %s", a.MCPServersKey())
	}
}

func TestSupportsUserScope_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if !a.SupportsUserScope() {
		t.Error("gemini should support user scope")
	}
}

func TestGetConfigPath_ContainsGemini_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, ".gemini") {
		t.Errorf("expected .gemini in path, got %s", p)
	}
}

func TestGetConfigPath_EndsWithJSON_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, ".json") {
		t.Errorf("expected .json suffix, got %s", p)
	}
}

func TestGetConfigPath_EmptyRoot_Extra4(t *testing.T) {
	a := New("", false)
	p := a.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetCurrentConfig_Missing_Extra4(t *testing.T) {
	a := New("/nonexistent/xyzabc", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for missing config")
	}
}

func TestNew_NotNil_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestNew_RuntimeEnvSubstitutionFalse_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("gemini should have SupportsRuntimeEnvSubstitution=false")
	}
}

func TestGetConfigPath_SettingsJSON_Extra4(t *testing.T) {
	a := New("/proj", false)
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, "settings.json") {
		t.Errorf("expected settings.json, got %s", p)
	}
}

func TestServerKeyFor_NamePriority_Extra4(t *testing.T) {
	k := serverKeyFor("a/b", "myname")
	if k != "myname" {
		t.Errorf("expected myname, got %s", k)
	}
}

func TestServerKeyFor_SlashPkg_Extra4(t *testing.T) {
	k := serverKeyFor("owner/pkg", "")
	if k != "pkg" {
		t.Errorf("expected pkg, got %s", k)
	}
}

func TestSelectRemoteWithURL_Empty_Extra4(t *testing.T) {
	r := selectRemoteWithURL(nil)
	if r != nil {
		t.Errorf("expected nil for empty remotes, got %v", r)
	}
}

func TestSelectRemoteWithURL_WithURL_Extra4(t *testing.T) {
	remotes := []map[string]interface{}{
		{"url": "https://example.com/sse"},
	}
	r := selectRemoteWithURL(remotes)
	if r == nil {
		t.Error("expected non-nil for remote with url")
	}
}

func TestSelectRemoteWithURL_NoURL_Extra4(t *testing.T) {
	remotes := []map[string]interface{}{
		{"other": "val"},
	}
	r := selectRemoteWithURL(remotes)
	if r != nil {
		t.Errorf("expected nil when no url field, got %v", r)
	}
}
