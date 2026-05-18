package manager

import (
	"testing"
)

func TestValidateRuntime_Known(t *testing.T) {
	m := New()
	for _, name := range []string{"copilot", "codex", "llm", "gemini"} {
		if err := m.ValidateRuntime(name); err != nil {
			t.Errorf("ValidateRuntime(%q) = %v; want nil", name, err)
		}
	}
}

func TestValidateRuntime_Unknown(t *testing.T) {
	m := New()
	if err := m.ValidateRuntime("unknown-runtime"); err == nil {
		t.Error("expected error for unknown runtime")
	}
}

func TestGetRuntimeDir(t *testing.T) {
	m := New()
	dir := m.GetRuntimeDir("copilot")
	if dir == "" {
		t.Error("GetRuntimeDir returned empty string")
	}
	// Should end with the runtime name
	if len(dir) < len("copilot") {
		t.Errorf("dir too short: %q", dir)
	}
}

func TestGetScriptPath_Known(t *testing.T) {
	m := New()
	path, err := m.GetScriptPath("copilot")
	if err != nil {
		t.Fatalf("GetScriptPath: %v", err)
	}
	if path == "" {
		t.Error("GetScriptPath returned empty path")
	}
}

func TestGetScriptPath_Unknown(t *testing.T) {
	m := New()
	_, err := m.GetScriptPath("nonexistent")
	if err == nil {
		t.Error("expected error for unknown runtime")
	}
}

func TestListRuntimes(t *testing.T) {
	m := New()
	runtimes := m.ListRuntimes()
	if len(runtimes) == 0 {
		t.Error("ListRuntimes returned empty map")
	}
	for name, desc := range runtimes {
		if name == "" || desc == "" {
			t.Errorf("empty name or description in ListRuntimes: %q -> %q", name, desc)
		}
	}
	// Verify required runtimes are present
	for _, required := range []string{"copilot", "codex", "llm", "gemini"} {
		if _, ok := runtimes[required]; !ok {
			t.Errorf("runtime %q missing from ListRuntimes", required)
		}
	}
}

func TestSetupEnvironment_WithToken(t *testing.T) {
	m := New()
	env, err := m.SetupEnvironment("copilot", "mytoken")
	if err != nil {
		t.Fatalf("SetupEnvironment: %v", err)
	}
	if env["GITHUB_TOKEN"] != "mytoken" {
		t.Errorf("GITHUB_TOKEN = %q; want mytoken", env["GITHUB_TOKEN"])
	}
	if env["GH_TOKEN"] != "mytoken" {
		t.Errorf("GH_TOKEN = %q; want mytoken", env["GH_TOKEN"])
	}
}

func TestSetupEnvironment_EmptyToken(t *testing.T) {
	m := New()
	env, err := m.SetupEnvironment("codex", "")
	if err != nil {
		t.Fatalf("SetupEnvironment: %v", err)
	}
	if _, ok := env["GITHUB_TOKEN"]; ok {
		t.Error("expected no GITHUB_TOKEN when token is empty")
	}
}

func TestSetupEnvironment_UnknownRuntime(t *testing.T) {
	m := New()
	_, err := m.SetupEnvironment("bogus", "token")
	if err == nil {
		t.Error("expected error for unknown runtime")
	}
}

func TestGetCommonScriptPath(t *testing.T) {
	m := New()
	p := m.GetCommonScriptPath()
	if p == "" {
		t.Error("GetCommonScriptPath returned empty string")
	}
}

func TestIsWindows(t *testing.T) {
	m := New()
	// Just verify it returns a bool without panic
	_ = m.IsWindows()
}

func TestNew_RuntimeDir(t *testing.T) {
	m := New()
	if m.RuntimeDir == "" {
		t.Error("RuntimeDir is empty")
	}
	if len(m.SupportedRuntimes) == 0 {
		t.Error("SupportedRuntimes is empty")
	}
}
