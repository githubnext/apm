package manager_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/runtime/manager"
)

func TestNew_HasFourRuntimes(t *testing.T) {
	m := manager.New()
	if len(m.SupportedRuntimes) != 4 {
		t.Errorf("expected 4 supported runtimes, got %d", len(m.SupportedRuntimes))
	}
}

func TestNew_HasCopilotRuntime(t *testing.T) {
	m := manager.New()
	if _, ok := m.SupportedRuntimes["copilot"]; !ok {
		t.Error("expected 'copilot' in supported runtimes")
	}
}

func TestNew_HasCodexRuntime(t *testing.T) {
	m := manager.New()
	if _, ok := m.SupportedRuntimes["codex"]; !ok {
		t.Error("expected 'codex' in supported runtimes")
	}
}

func TestNew_HasLLMRuntime(t *testing.T) {
	m := manager.New()
	if _, ok := m.SupportedRuntimes["llm"]; !ok {
		t.Error("expected 'llm' in supported runtimes")
	}
}

func TestNew_HasGeminiRuntime(t *testing.T) {
	m := manager.New()
	if _, ok := m.SupportedRuntimes["gemini"]; !ok {
		t.Error("expected 'gemini' in supported runtimes")
	}
}

func TestListRuntimes_AllFour(t *testing.T) {
	m := manager.New()
	list := m.ListRuntimes()
	for _, name := range []string{"copilot", "codex", "llm", "gemini"} {
		if _, ok := list[name]; !ok {
			t.Errorf("expected runtime %q in ListRuntimes output", name)
		}
	}
}

func TestListRuntimes_DescriptionsNonEmpty(t *testing.T) {
	m := manager.New()
	for name, desc := range m.ListRuntimes() {
		if desc == "" {
			t.Errorf("runtime %q has empty description", name)
		}
	}
}

func TestValidateRuntime_AllKnown(t *testing.T) {
	m := manager.New()
	for _, name := range []string{"copilot", "codex", "llm", "gemini"} {
		if err := m.ValidateRuntime(name); err != nil {
			t.Errorf("ValidateRuntime(%q) unexpected error: %v", name, err)
		}
	}
}

func TestValidateRuntime_UnknownContainsSuggestion(t *testing.T) {
	m := manager.New()
	err := m.ValidateRuntime("bogus")
	if err == nil {
		t.Fatal("expected error for unknown runtime")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("error should mention 'bogus', got: %v", err)
	}
}

func TestGetRuntimeDir_ContainsRuntimeName(t *testing.T) {
	m := manager.New()
	dir := m.GetRuntimeDir("copilot")
	if !strings.Contains(dir, "copilot") {
		t.Errorf("runtime dir should contain 'copilot', got %q", dir)
	}
}

func TestGetRuntimeDir_DifferentRuntimes(t *testing.T) {
	m := manager.New()
	d1 := m.GetRuntimeDir("copilot")
	d2 := m.GetRuntimeDir("codex")
	if d1 == d2 {
		t.Error("different runtimes should have different directories")
	}
}

func TestGetScriptPath_KnownRuntimeNoError(t *testing.T) {
	m := manager.New()
	_, err := m.GetScriptPath("copilot")
	if err != nil {
		t.Fatalf("unexpected error for known runtime: %v", err)
	}
}

func TestGetScriptPath_ContainsRuntimeName(t *testing.T) {
	m := manager.New()
	p, _ := m.GetScriptPath("llm")
	if !strings.Contains(p, "llm") {
		t.Errorf("script path should contain 'llm', got %q", p)
	}
}

func TestSetupEnvironment_WithTokenSetsVars(t *testing.T) {
	m := manager.New()
	env, err := m.SetupEnvironment("copilot", "mytoken123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["GITHUB_TOKEN"] != "mytoken123" {
		t.Errorf("GITHUB_TOKEN: got %q want %q", env["GITHUB_TOKEN"], "mytoken123")
	}
	if env["GH_TOKEN"] != "mytoken123" {
		t.Errorf("GH_TOKEN: got %q want %q", env["GH_TOKEN"], "mytoken123")
	}
}

func TestSetupEnvironment_EmptyTokenNoVars(t *testing.T) {
	m := manager.New()
	env, err := m.SetupEnvironment("llm", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["GITHUB_TOKEN"]; ok {
		t.Error("GITHUB_TOKEN should not be set for empty token")
	}
}

func TestGetCommonScriptPath_NonEmpty(t *testing.T) {
	m := manager.New()
	p := m.GetCommonScriptPath()
	if p == "" {
		t.Error("common script path should be non-empty")
	}
}

func TestGetCommonScriptPath_ContainsCommon(t *testing.T) {
	m := manager.New()
	p := m.GetCommonScriptPath()
	if !strings.Contains(p, "common") {
		t.Errorf("common script path should contain 'common', got %q", p)
	}
}

func TestRuntimeInfo_BinaryNonEmpty(t *testing.T) {
	m := manager.New()
	for name, info := range m.SupportedRuntimes {
		if info.Binary == "" {
			t.Errorf("runtime %q has empty Binary field", name)
		}
	}
}

func TestRuntimeInfo_ScriptNonEmpty(t *testing.T) {
	m := manager.New()
	for name, info := range m.SupportedRuntimes {
		if info.Script == "" {
			t.Errorf("runtime %q has empty Script field", name)
		}
	}
}
