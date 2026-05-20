package base

import (
	"testing"
)

// multiAdapter is a test-only adapter that returns configurable values.
type multiAdapter struct {
	name      string
	available bool
	info      map[string]any
	models    map[string]any
	execErr   error
	execOut   string
}

func (m *multiAdapter) ExecutePrompt(_ string, _ map[string]any) (string, error) {
	return m.execOut, m.execErr
}
func (m *multiAdapter) ListAvailableModels() map[string]any {
	if m.models == nil {
		return map[string]any{}
	}
	return m.models
}
func (m *multiAdapter) GetRuntimeInfo() map[string]any {
	if m.info == nil {
		return map[string]any{}
	}
	return m.info
}
func (m *multiAdapter) IsAvailable() bool    { return m.available }
func (m *multiAdapter) GetRuntimeName() string { return m.name }

func TestMultiAdapter_AvailableFalse(t *testing.T) {
	var a RuntimeAdapter = &multiAdapter{name: "x", available: false}
	if a.IsAvailable() {
		t.Error("expected IsAvailable false")
	}
}

func TestMultiAdapter_GetRuntimeName(t *testing.T) {
	a := &multiAdapter{name: "claude"}
	if a.GetRuntimeName() != "claude" {
		t.Errorf("expected claude, got %q", a.GetRuntimeName())
	}
}

func TestMultiAdapter_ExecutePrompt_Success(t *testing.T) {
	a := &multiAdapter{execOut: "result", available: true}
	out, err := a.ExecutePrompt("prompt", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out != "result" {
		t.Errorf("expected result, got %q", out)
	}
}

func TestMultiAdapter_GetRuntimeInfo_HasEntries(t *testing.T) {
	a := &multiAdapter{info: map[string]any{"version": "1.0"}}
	info := a.GetRuntimeInfo()
	if info["version"] != "1.0" {
		t.Errorf("expected version 1.0, got %v", info["version"])
	}
}

func TestMultiAdapter_ListModels_NonNil(t *testing.T) {
	a := &multiAdapter{models: map[string]any{"model-a": true}}
	m := a.ListAvailableModels()
	if m == nil {
		t.Error("expected non-nil models")
	}
	if _, ok := m["model-a"]; !ok {
		t.Error("expected model-a in models")
	}
}

func TestRuntimeAdapter_Interface_Assignment(t *testing.T) {
	var a RuntimeAdapter = &multiAdapter{name: "test"}
	if a.GetRuntimeName() != "test" {
		t.Errorf("unexpected name %q", a.GetRuntimeName())
	}
}

func TestMultiAdapter_ExecutePrompt_WithArgs(t *testing.T) {
	a := &multiAdapter{execOut: "ok"}
	out, err := a.ExecutePrompt("p", map[string]any{"k": "v"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out != "ok" {
		t.Errorf("expected ok, got %q", out)
	}
}

func TestMultiAdapter_GetRuntimeInfo_Empty(t *testing.T) {
	a := &multiAdapter{}
	info := a.GetRuntimeInfo()
	if info == nil {
		t.Error("expected non-nil map")
	}
	if len(info) != 0 {
		t.Errorf("expected empty map, got %v", info)
	}
}
