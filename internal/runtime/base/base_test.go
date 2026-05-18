package base

import "testing"

// mockAdapter is a minimal implementation of RuntimeAdapter for compilation testing.
type mockAdapter struct{}

func (m *mockAdapter) ExecutePrompt(_ string, _ map[string]any) (string, error) { return "", nil }
func (m *mockAdapter) ListAvailableModels() map[string]any                      { return nil }
func (m *mockAdapter) GetRuntimeInfo() map[string]any                           { return nil }
func (m *mockAdapter) IsAvailable() bool                                        { return false }
func (m *mockAdapter) GetRuntimeName() string                                   { return "mock" }

func TestRuntimeAdapterInterface(t *testing.T) {
	var adapter RuntimeAdapter = &mockAdapter{}
	if adapter.GetRuntimeName() != "mock" {
		t.Errorf("unexpected runtime name")
	}
	if adapter.IsAvailable() {
		t.Error("expected IsAvailable false")
	}
	models := adapter.ListAvailableModels()
	if models != nil {
		t.Error("expected nil models")
	}
}

// namedAdapter implements RuntimeAdapter with a configurable name.
type namedAdapter struct {
	name      string
	available bool
}

func (n *namedAdapter) ExecutePrompt(_ string, _ map[string]any) (string, error) {
	return "response", nil
}
func (n *namedAdapter) ListAvailableModels() map[string]any {
	return map[string]any{"default": "gpt-4"}
}
func (n *namedAdapter) GetRuntimeInfo() map[string]any {
	return map[string]any{"name": n.name}
}
func (n *namedAdapter) IsAvailable() bool    { return n.available }
func (n *namedAdapter) GetRuntimeName() string { return n.name }

func TestNamedAdapterAvailable(t *testing.T) {
	a := &namedAdapter{name: "openai", available: true}
	var iface RuntimeAdapter = a
	if iface.GetRuntimeName() != "openai" {
		t.Errorf("GetRuntimeName = %q, want openai", iface.GetRuntimeName())
	}
	if !iface.IsAvailable() {
		t.Error("expected IsAvailable true")
	}
}

func TestNamedAdapterUnavailable(t *testing.T) {
	a := &namedAdapter{name: "anthropic", available: false}
	var iface RuntimeAdapter = a
	if iface.IsAvailable() {
		t.Error("expected IsAvailable false")
	}
}

func TestNamedAdapterListModels(t *testing.T) {
	a := &namedAdapter{name: "gemini", available: true}
	models := a.ListAvailableModels()
	if models == nil {
		t.Fatal("expected non-nil models map")
	}
	if _, ok := models["default"]; !ok {
		t.Error("expected 'default' key in models map")
	}
}

func TestNamedAdapterGetRuntimeInfo(t *testing.T) {
	a := &namedAdapter{name: "claude", available: true}
	info := a.GetRuntimeInfo()
	if info == nil {
		t.Fatal("expected non-nil runtime info")
	}
	if info["name"] != "claude" {
		t.Errorf("runtime info name = %q, want claude", info["name"])
	}
}

func TestNamedAdapterExecutePrompt(t *testing.T) {
	a := &namedAdapter{name: "test", available: true}
	resp, err := a.ExecutePrompt("hello", nil)
	if err != nil {
		t.Fatalf("ExecutePrompt returned error: %v", err)
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestMockAdapterExecutePrompt(t *testing.T) {
	m := &mockAdapter{}
	resp, err := m.ExecutePrompt("test prompt", map[string]any{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "" {
		t.Errorf("expected empty string, got %q", resp)
	}
}

func TestInterfaceSlice(t *testing.T) {
	adapters := []RuntimeAdapter{
		&mockAdapter{},
		&namedAdapter{name: "openai", available: true},
		&namedAdapter{name: "anthropic", available: false},
	}
	names := map[string]bool{}
	for _, a := range adapters {
		names[a.GetRuntimeName()] = true
	}
	for _, want := range []string{"mock", "openai", "anthropic"} {
		if !names[want] {
			t.Errorf("missing adapter with name %q", want)
		}
	}
}
