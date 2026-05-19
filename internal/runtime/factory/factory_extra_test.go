package factory

import (
	"errors"
	"testing"
)

// mockAdapter is a test double for ConstructableAdapter.
type mockAdapter struct {
	name      string
	available bool
	info      RuntimeInfo
	newErr    error
}

func (m *mockAdapter) GetRuntimeName() string    { return m.name }
func (m *mockAdapter) IsAvailable() bool         { return m.available }
func (m *mockAdapter) GetRuntimeInfo() RuntimeInfo {
	if m.info.Name == "" {
		return RuntimeInfo{Name: m.name, Available: m.available}
	}
	return m.info
}
func (m *mockAdapter) New(modelName string) (RuntimeAdapter, error) {
	if m.newErr != nil {
		return nil, m.newErr
	}
	return &mockRuntimeAdapter{name: m.name, info: m.info}, nil
}
func (m *mockAdapter) NewDefault() (RuntimeAdapter, error) {
	if m.newErr != nil {
		return nil, m.newErr
	}
	return &mockRuntimeAdapter{name: m.name, info: m.info}, nil
}

type mockRuntimeAdapter struct {
	name string
	info RuntimeInfo
}

func (r *mockRuntimeAdapter) GetRuntimeName() string    { return r.name }
func (r *mockRuntimeAdapter) IsAvailable() bool         { return true }
func (r *mockRuntimeAdapter) GetRuntimeInfo() RuntimeInfo {
	if r.info.Name == "" {
		return RuntimeInfo{Name: r.name, Available: true}
	}
	return r.info
}

func TestRegistry_GetAvailableRuntimes_NewDefaultError(t *testing.T) {
	reg := NewRegistry(&mockAdapter{
		name:      "faulty",
		available: true,
		newErr:    errors.New("init failure"),
	})
	infos := reg.GetAvailableRuntimes()
	if len(infos) != 1 {
		t.Fatalf("expected 1 runtime, got %d", len(infos))
	}
	if infos[0].Error == "" {
		t.Error("expected Error to be set when NewDefault fails")
	}
}

func TestRegistry_GetRuntimeByName_WithModel(t *testing.T) {
	reg := NewRegistry(&mockAdapter{name: "my-rt", available: true})
	rt, err := reg.GetRuntimeByName("my-rt", "gpt-4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rt == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestRegistry_GetAvailableRuntimes_MultipleAdapters(t *testing.T) {
	reg := NewRegistry(
		&mockAdapter{name: "a", available: true},
		&mockAdapter{name: "b", available: false},
		&mockAdapter{name: "c", available: true},
	)
	infos := reg.GetAvailableRuntimes()
	if len(infos) != 2 {
		t.Fatalf("expected 2 available, got %d", len(infos))
	}
	names := map[string]bool{}
	for _, i := range infos {
		names[i.Name] = true
	}
	if !names["a"] || !names["c"] {
		t.Error("expected adapters 'a' and 'c'")
	}
}

func TestRuntimeInfo_Fields(t *testing.T) {
	ri := RuntimeInfo{Name: "test-rt", Available: true, Error: "some err"}
	if ri.Name != "test-rt" {
		t.Errorf("Name mismatch: %q", ri.Name)
	}
	if !ri.Available {
		t.Error("Available should be true")
	}
	if ri.Error != "some err" {
		t.Errorf("Error mismatch: %q", ri.Error)
	}
}

func TestRuntimeInfo_ZeroValue(t *testing.T) {
	var ri RuntimeInfo
	if ri.Name != "" || ri.Available || ri.Error != "" {
		t.Error("zero-value RuntimeInfo should have empty fields")
	}
}

func TestRegistry_GetRuntimeByName_NewErr(t *testing.T) {
	reg := NewRegistry(&mockAdapter{name: "broken", available: true, newErr: errors.New("fail")})
	_, err := reg.GetRuntimeByName("broken", "")
	if err == nil {
		t.Error("expected error from GetRuntimeByName when New fails")
	}
}

func TestNewRegistry_Empty(t *testing.T) {
	reg := NewRegistry()
	infos := reg.GetAvailableRuntimes()
	if len(infos) != 0 {
		t.Errorf("expected empty registry to have 0 runtimes, got %d", len(infos))
	}
}
