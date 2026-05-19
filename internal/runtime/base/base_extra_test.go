package base

import (
	"errors"
	"testing"
)

// errorAdapter always returns errors.
type errorAdapter struct{}

func (e *errorAdapter) ExecutePrompt(_ string, _ map[string]any) (string, error) {
	return "", errors.New("runtime error")
}
func (e *errorAdapter) ListAvailableModels() map[string]any { return map[string]any{} }
func (e *errorAdapter) GetRuntimeInfo() map[string]any      { return map[string]any{"error": true} }
func (e *errorAdapter) IsAvailable() bool                   { return false }
func (e *errorAdapter) GetRuntimeName() string              { return "error-adapter" }

func TestErrorAdapterExecutePrompt(t *testing.T) {
	var a RuntimeAdapter = &errorAdapter{}
	_, err := a.ExecutePrompt("prompt", nil)
	if err == nil {
		t.Error("expected error from errorAdapter")
	}
}

func TestErrorAdapterGetRuntimeName(t *testing.T) {
	a := &errorAdapter{}
	if a.GetRuntimeName() != "error-adapter" {
		t.Errorf("unexpected runtime name: %q", a.GetRuntimeName())
	}
}

func TestErrorAdapterListModels_Empty(t *testing.T) {
	a := &errorAdapter{}
	m := a.ListAvailableModels()
	if m == nil {
		t.Error("expected non-nil map (empty is ok)")
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestNamedAdapter_ExecutePromptWithArgs(t *testing.T) {
	a := &namedAdapter{name: "test", available: true}
	args := map[string]any{
		"key1": "value1",
		"key2": 42,
	}
	resp, err := a.ExecutePrompt("prompt with args", args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestNamedAdapter_MultipleNames(t *testing.T) {
	names := []string{"openai", "claude", "gemini", "mistral", "llama"}
	for _, name := range names {
		a := &namedAdapter{name: name, available: true}
		if a.GetRuntimeName() != name {
			t.Errorf("GetRuntimeName() = %q, want %q", a.GetRuntimeName(), name)
		}
	}
}

func TestRuntimeAdapter_PolymorphicSlice(t *testing.T) {
	adapters := []RuntimeAdapter{
		&mockAdapter{},
		&errorAdapter{},
		&namedAdapter{name: "a", available: true},
		&namedAdapter{name: "b", available: false},
	}
	for _, a := range adapters {
		name := a.GetRuntimeName()
		if name == "" {
			t.Error("GetRuntimeName should never return empty string")
		}
		// IsAvailable should not panic
		_ = a.IsAvailable()
		// ListAvailableModels should not panic
		_ = a.ListAvailableModels()
		// GetRuntimeInfo should not panic
		_ = a.GetRuntimeInfo()
	}
}

func TestMockAdapter_GetRuntimeInfo_Nil(t *testing.T) {
	m := &mockAdapter{}
	info := m.GetRuntimeInfo()
	if info != nil {
		t.Errorf("expected nil info from mockAdapter, got %v", info)
	}
}

func TestNamedAdapter_RuntimeInfoContainsName(t *testing.T) {
	a := &namedAdapter{name: "myruntime", available: true}
	info := a.GetRuntimeInfo()
	if info == nil {
		t.Fatal("expected non-nil info")
	}
	if info["name"] != "myruntime" {
		t.Errorf("info[name] = %v, want myruntime", info["name"])
	}
}
