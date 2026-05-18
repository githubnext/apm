package opencode_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/opencode"
)

func TestToOpenCodeFormat_CommandEntry(t *testing.T) {
	entry := opencode.CopilotEntry{
		Command: "npx",
		Args:    []string{"-y", "some-pkg"},
		Env:     map[string]string{"KEY": "val"},
	}
	got := opencode.ToOpenCodeFormat(entry, true)
	if got.Type != "local" {
		t.Errorf("Type = %q, want local", got.Type)
	}
	if !got.Enabled {
		t.Error("Enabled should be true")
	}
	if len(got.Command) < 1 || got.Command[0] != "npx" {
		t.Errorf("Command[0] = %q, want npx", got.Command[0])
	}
	if got.Environment["KEY"] != "val" {
		t.Errorf("Environment[KEY] = %q, want val", got.Environment["KEY"])
	}
}

func TestToOpenCodeFormat_URLEntry(t *testing.T) {
	entry := opencode.CopilotEntry{
		URL:     "http://localhost:3000",
		Headers: map[string]string{"Auth": "token"},
	}
	got := opencode.ToOpenCodeFormat(entry, false)
	if got.URL != "http://localhost:3000" {
		t.Errorf("URL = %q", got.URL)
	}
	if got.Enabled {
		t.Error("Enabled should be false")
	}
}

func TestToOpenCodeFormat_Disabled(t *testing.T) {
	entry := opencode.CopilotEntry{Command: "cmd", Args: []string{}}
	got := opencode.ToOpenCodeFormat(entry, false)
	if got.Enabled {
		t.Error("Enabled should be false when disabled=false")
	}
}

func TestNew_ConfigPath(t *testing.T) {
	a := opencode.New("/some/project")
	want := filepath.Join("/some/project", "opencode.json")
	if a.ConfigPath() != want {
		t.Errorf("ConfigPath() = %q, want %q", a.ConfigPath(), want)
	}
}

func TestIsOptedIn_False(t *testing.T) {
	tmp := t.TempDir()
	a := opencode.New(tmp)
	if a.IsOptedIn() {
		t.Error("expected IsOptedIn=false when .opencode/ does not exist")
	}
}

func TestIsOptedIn_True(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, ".opencode"), 0o755); err != nil {
		t.Fatal(err)
	}
	a := opencode.New(tmp)
	if !a.IsOptedIn() {
		t.Error("expected IsOptedIn=true when .opencode/ exists")
	}
}

func TestGetCurrentConfig_NoFile(t *testing.T) {
	tmp := t.TempDir()
	a := opencode.New(tmp)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map on missing file")
	}
}
