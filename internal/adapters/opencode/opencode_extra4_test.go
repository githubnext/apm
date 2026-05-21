package opencode

import (
	"testing"
)

func TestToOpenCodeFormat_LocalType_Extra4(t *testing.T) {
	entry := CopilotEntry{Command: "npx", Args: []string{"-y", "pkg"}}
	result := ToOpenCodeFormat(entry, true)
	if result.Type != "local" {
		t.Errorf("expected local, got %s", result.Type)
	}
}

func TestToOpenCodeFormat_Enabled_Extra4(t *testing.T) {
	entry := CopilotEntry{Command: "npx"}
	result := ToOpenCodeFormat(entry, true)
	if !result.Enabled {
		t.Error("expected enabled=true")
	}
}

func TestToOpenCodeFormat_Disabled_Extra4(t *testing.T) {
	entry := CopilotEntry{Command: "npx"}
	result := ToOpenCodeFormat(entry, false)
	if result.Enabled {
		t.Error("expected enabled=false")
	}
}

func TestToOpenCodeFormat_CommandCombined_Extra4(t *testing.T) {
	entry := CopilotEntry{Command: "npx", Args: []string{"-y", "pkg"}}
	result := ToOpenCodeFormat(entry, true)
	if len(result.Command) == 0 {
		t.Error("expected non-empty command")
	}
	if result.Command[0] != "npx" {
		t.Errorf("expected npx, got %s", result.Command[0])
	}
}

func TestToOpenCodeFormat_WithEnv_Extra4(t *testing.T) {
	entry := CopilotEntry{Command: "cmd", Env: map[string]string{"KEY": "val"}}
	result := ToOpenCodeFormat(entry, true)
	if result.Environment["KEY"] != "val" {
		t.Errorf("expected val, got %s", result.Environment["KEY"])
	}
}

func TestToOpenCodeFormat_URLType_Extra4(t *testing.T) {
	entry := CopilotEntry{URL: "https://example.com/sse"}
	result := ToOpenCodeFormat(entry, true)
	if result.URL != "https://example.com/sse" {
		t.Errorf("expected url preserved, got %s", result.URL)
	}
}

func TestToOpenCodeFormat_EmptyEntry_Extra4(t *testing.T) {
	entry := CopilotEntry{}
	result := ToOpenCodeFormat(entry, false)
	if result.Type != "local" {
		t.Errorf("expected local type, got %s", result.Type)
	}
}

func TestNew_NotNil_Extra4(t *testing.T) {
	a := New("/tmp")
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestConfigPath_NotEmpty_Extra4(t *testing.T) {
	a := New("/myproject")
	p := a.ConfigPath()
	if p == "" {
		t.Error("expected non-empty config path")
	}
}

func TestIsOptedIn_FalseForMissing_Extra4(t *testing.T) {
	a := New("/nonexistent/xyzabc")
	if a.IsOptedIn() {
		t.Error("expected not opted in for missing directory")
	}
}

func TestGetCurrentConfig_Missing_Extra4(t *testing.T) {
	a := New("/nonexistent/xyzabc")
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map")
	}
}

func TestServerEntry_Fields_Extra4(t *testing.T) {
	se := ServerEntry{Type: "local", Enabled: true}
	if se.Type != "local" {
		t.Error("expected local")
	}
	if !se.Enabled {
		t.Error("expected enabled")
	}
}

func TestCopilotEntry_Fields_Extra4(t *testing.T) {
	ce := CopilotEntry{Command: "npx", Args: []string{"-y"}}
	if ce.Command != "npx" {
		t.Error("expected npx")
	}
	if len(ce.Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(ce.Args))
	}
}
