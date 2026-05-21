package manager_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/runtime/manager"
)

func TestNew_ReturnsNonNil(t *testing.T) {
	m := manager.New()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestNew_SupportedRuntimesNotNil(t *testing.T) {
	m := manager.New()
	if m.SupportedRuntimes == nil {
		t.Error("SupportedRuntimes should not be nil")
	}
}

func TestListRuntimes_IncludesCopilot(t *testing.T) {
	m := manager.New()
	runtimes := m.ListRuntimes()
	if _, ok := runtimes["copilot"]; !ok {
		t.Error("expected 'copilot' runtime in list")
	}
}

func TestListRuntimes_AllDescriptionsNonEmpty(t *testing.T) {
	m := manager.New()
	for name, desc := range m.ListRuntimes() {
		if desc == "" {
			t.Errorf("runtime %q has empty description", name)
		}
	}
}

func TestValidateRuntime_CopilotKnown(t *testing.T) {
	m := manager.New()
	if err := m.ValidateRuntime("copilot"); err != nil {
		t.Errorf("copilot should be valid: %v", err)
	}
}

func TestValidateRuntime_UnknownReturnsError(t *testing.T) {
	m := manager.New()
	if err := m.ValidateRuntime("notaruntime"); err == nil {
		t.Error("expected error for unknown runtime")
	}
}

func TestValidateRuntime_ErrorMentionsSuggestion(t *testing.T) {
	m := manager.New()
	err := m.ValidateRuntime("notaruntime")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "copilot") && !strings.Contains(err.Error(), "notaruntime") {
		t.Errorf("error should mention runtimes: %q", err.Error())
	}
}

func TestGetRuntimeDir_ContainsName(t *testing.T) {
	m := manager.New()
	dir := m.GetRuntimeDir("copilot")
	if !strings.Contains(dir, "copilot") {
		t.Errorf("runtime dir should contain 'copilot': %q", dir)
	}
}

func TestGetScriptPath_KnownRuntime(t *testing.T) {
	m := manager.New()
	p, err := m.GetScriptPath("copilot")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == "" {
		t.Error("expected non-empty path")
	}
}

func TestGetScriptPath_UnknownRuntime(t *testing.T) {
	m := manager.New()
	_, err := m.GetScriptPath("nosuchruntime")
	if err == nil {
		t.Error("expected error for unknown runtime")
	}
}

func TestGetCommonScriptPath_NonEmptyResult(t *testing.T) {
	m := manager.New()
	p := m.GetCommonScriptPath()
	if p == "" {
		t.Error("expected non-empty common script path")
	}
}

func TestRuntimeInfo_ZeroValue(t *testing.T) {
	var r manager.RuntimeInfo
	if r.Script != "" || r.Description != "" || r.Binary != "" {
		t.Error("zero value should have empty fields")
	}
}
