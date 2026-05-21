package opencode_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/opencode"
)

func TestNew_NotNil_Extra3(t *testing.T) {
	a := opencode.New("/tmp")
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestConfigPath_NotEmpty_Extra3(t *testing.T) {
	a := opencode.New(t.TempDir())
	if a.ConfigPath() == "" {
		t.Error("expected non-empty config path")
	}
}

func TestConfigPath_EndsWithOpencodeJson_Extra3(t *testing.T) {
	a := opencode.New(t.TempDir())
	p := a.ConfigPath()
	if filepath.Base(p) != "opencode.json" {
		t.Errorf("expected opencode.json, got %q", filepath.Base(p))
	}
}

func TestIsOptedIn_NoDir_Extra3(t *testing.T) {
	a := opencode.New(t.TempDir())
	if a.IsOptedIn() {
		t.Error("expected false when .opencode dir does not exist")
	}
}

func TestIsOptedIn_WithDir_Extra3(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".opencode"), 0o755)
	a := opencode.New(dir)
	if !a.IsOptedIn() {
		t.Error("expected true when .opencode dir exists")
	}
}

func TestGetCurrentConfig_NonNilOnMissing_Extra3(t *testing.T) {
	a := opencode.New(t.TempDir())
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map")
	}
}

func TestToOpenCodeFormat_LocalType_Extra3(t *testing.T) {
	e := opencode.CopilotEntry{Command: "node", Args: []string{"server.js"}}
	r := opencode.ToOpenCodeFormat(e, true)
	if r.Type != "local" {
		t.Errorf("expected local, got %q", r.Type)
	}
}

func TestToOpenCodeFormat_EnabledFalse_Extra3(t *testing.T) {
	e := opencode.CopilotEntry{Command: "node"}
	r := opencode.ToOpenCodeFormat(e, false)
	if r.Enabled {
		t.Error("expected Enabled false")
	}
}

func TestToOpenCodeFormat_EmptyEntry_Extra3(t *testing.T) {
	e := opencode.CopilotEntry{}
	r := opencode.ToOpenCodeFormat(e, true)
	if r.Type == "" {
		t.Error("expected non-empty Type")
	}
}

func TestToOpenCodeFormat_WithEnv_Extra3(t *testing.T) {
	e := opencode.CopilotEntry{
		Command: "cmd",
		Env:     map[string]string{"KEY": "val"},
	}
	r := opencode.ToOpenCodeFormat(e, true)
	if r.Environment["KEY"] != "val" {
		t.Errorf("expected env KEY=val, got %v", r.Environment)
	}
}

func TestServerEntry_ZeroValue_Extra3(t *testing.T) {
	var e opencode.ServerEntry
	if e.Enabled {
		t.Error("expected zero Enabled")
	}
	if e.Type != "" {
		t.Error("expected zero Type")
	}
}

func TestCopilotEntry_ZeroValue_Extra3(t *testing.T) {
	var e opencode.CopilotEntry
	if e.Command != "" {
		t.Error("expected zero Command")
	}
}
