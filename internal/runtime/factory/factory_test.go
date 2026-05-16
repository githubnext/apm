package factory_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/runtime/factory"
)

// mockAdapter is a minimal ConstructableAdapter for testing.
type mockAdapter struct {
	name      string
	available bool
	failNew   bool
}

func (m *mockAdapter) GetRuntimeName() string { return m.name }
func (m *mockAdapter) IsAvailable() bool       { return m.available }
func (m *mockAdapter) GetRuntimeInfo() factory.RuntimeInfo {
	return factory.RuntimeInfo{Name: m.name, Available: m.available}
}
func (m *mockAdapter) New(modelName string) (factory.RuntimeAdapter, error) {
	if m.failNew {
		return nil, errors.New("init failed")
	}
	return &mockInstance{name: m.name}, nil
}
func (m *mockAdapter) NewDefault() (factory.RuntimeAdapter, error) {
	return m.New("")
}

// mockInstance satisfies RuntimeAdapter.
type mockInstance struct{ name string }

func (i *mockInstance) GetRuntimeName() string                        { return i.name }
func (i *mockInstance) IsAvailable() bool                             { return true }
func (i *mockInstance) GetRuntimeInfo() factory.RuntimeInfo           { return factory.RuntimeInfo{Name: i.name, Available: true} }

// --- Registry tests ---

func TestRegistry_GetAvailableRuntimes_Empty(t *testing.T) {
	r := factory.NewRegistry()
	if got := r.GetAvailableRuntimes(); len(got) != 0 {
		t.Fatalf("expected 0 runtimes, got %d", len(got))
	}
}

func TestRegistry_GetAvailableRuntimes_SkipsUnavailable(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "unavail", available: false})
	if got := r.GetAvailableRuntimes(); len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestRegistry_GetAvailableRuntimes_IncludesAvailable(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "copilot", available: true})
	got := r.GetAvailableRuntimes()
	if len(got) != 1 || got[0].Name != "copilot" {
		t.Fatalf("unexpected runtimes: %v", got)
	}
}

func TestRegistry_GetRuntimeByName_NotFound(t *testing.T) {
	r := factory.NewRegistry()
	_, err := r.GetRuntimeByName("missing", "")
	if err == nil {
		t.Fatal("expected error for missing runtime")
	}
}

func TestRegistry_GetRuntimeByName_Unavailable(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "copilot", available: false})
	_, err := r.GetRuntimeByName("copilot", "")
	if err == nil {
		t.Fatal("expected error for unavailable runtime")
	}
}

func TestRegistry_GetRuntimeByName_Found(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "copilot", available: true})
	rt, err := r.GetRuntimeByName("copilot", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rt.GetRuntimeName() != "copilot" {
		t.Errorf("expected copilot, got %s", rt.GetRuntimeName())
	}
}

func TestRegistry_GetBestAvailableRuntime_NoneAvailable(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "x", available: false})
	_, err := r.GetBestAvailableRuntime("")
	if err == nil {
		t.Fatal("expected error when no runtimes available")
	}
}

func TestRegistry_GetBestAvailableRuntime_ReturnsFirst(t *testing.T) {
	r := factory.NewRegistry(
		&mockAdapter{name: "first", available: true},
		&mockAdapter{name: "second", available: true},
	)
	rt, err := r.GetBestAvailableRuntime("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rt.GetRuntimeName() != "first" {
		t.Errorf("expected first, got %s", rt.GetRuntimeName())
	}
}

func TestRegistry_CreateRuntime_ByName(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "codex", available: true})
	rt, err := r.CreateRuntime("codex", "")
	if err != nil || rt.GetRuntimeName() != "codex" {
		t.Fatalf("CreateRuntime by name failed: %v", err)
	}
}

func TestRegistry_CreateRuntime_BestAvailable(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "llm", available: true})
	rt, err := r.CreateRuntime("", "")
	if err != nil || rt.GetRuntimeName() != "llm" {
		t.Fatalf("CreateRuntime best-available failed: %v", err)
	}
}

func TestRegistry_RuntimeExists(t *testing.T) {
	r := factory.NewRegistry(&mockAdapter{name: "copilot", available: true})
	if !r.RuntimeExists("copilot") {
		t.Fatal("expected RuntimeExists=true for copilot")
	}
	if r.RuntimeExists("nonexistent") {
		t.Fatal("expected RuntimeExists=false for nonexistent")
	}
}
