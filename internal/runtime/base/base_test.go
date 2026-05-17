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
