package factory

import (
	"testing"
)

func TestRuntimeInfo_AvailableField(t *testing.T) {
	ri := RuntimeInfo{Name: "copilot", Available: true}
	if !ri.Available {
		t.Error("expected available=true")
	}
}

func TestRuntimeInfo_ErrorField(t *testing.T) {
	ri := RuntimeInfo{Name: "copilot", Error: "not installed"}
	if ri.Error != "not installed" {
		t.Errorf("expected error, got %q", ri.Error)
	}
}

func TestNewRegistry_SingleAdapter(t *testing.T) {
	a := &mockAdapter{name: "test", available: true}
	reg := NewRegistry(a)
	if reg == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRuntimeExists_Present(t *testing.T) {
	a := &mockAdapter{name: "myruntime", available: true}
	reg := NewRegistry(a)
	if !reg.RuntimeExists("myruntime") {
		t.Error("expected myruntime to exist")
	}
}

func TestRuntimeExists_Absent(t *testing.T) {
	a := &mockAdapter{name: "myruntime", available: true}
	reg := NewRegistry(a)
	if reg.RuntimeExists("other") {
		t.Error("other should not exist")
	}
}

func TestRuntimeExists_Empty(t *testing.T) {
	reg := NewRegistry()
	if reg.RuntimeExists("anything") {
		t.Error("empty registry should not have any runtime")
	}
}

func TestGetAvailableRuntimes_NoneAvailable(t *testing.T) {
	a := &mockAdapter{name: "unavailable", available: false}
	reg := NewRegistry(a)
	runtimes := reg.GetAvailableRuntimes()
	if len(runtimes) != 0 {
		t.Errorf("expected no available runtimes, got %d", len(runtimes))
	}
}

func TestGetBestAvailableRuntime_NoneAvailable(t *testing.T) {
	a := &mockAdapter{name: "rt", available: false}
	reg := NewRegistry(a)
	_, err := reg.GetBestAvailableRuntime("model")
	if err == nil {
		t.Error("expected error when no runtime available")
	}
}

func TestCreateRuntime_UnknownRuntime(t *testing.T) {
	a := &mockAdapter{name: "known", available: true}
	reg := NewRegistry(a)
	_, err := reg.CreateRuntime("unknown", "model")
	if err == nil {
		t.Error("expected error for unknown runtime")
	}
}
