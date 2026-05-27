package runtime_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/runtime"
)

// TestParityRuntimeAdapterInterface verifies the RuntimeAdapter interface is defined.
func TestParityRuntimeAdapterInterface(t *testing.T) {
	// RuntimeAdapter is an interface -- verify it exists as a type
	var _ runtime.RuntimeAdapter = (*mockAdapter)(nil)
}

type mockAdapter struct{}

func (m *mockAdapter) ExecutePrompt(prompt string, opts map[string]interface{}) (string, error) {
	return "ok", nil
}
func (m *mockAdapter) ListAvailableModels() (map[string]interface{}, error) {
	return map[string]interface{}{"default": true}, nil
}
func (m *mockAdapter) GetRuntimeInfo() map[string]interface{} {
	return map[string]interface{}{"name": "mock"}
}
func (m *mockAdapter) IsAvailable() bool    { return true }
func (m *mockAdapter) GetRuntimeName() string { return "mock" }

// TestParityRuntimeErrors verifies sentinel errors exist.
func TestParityRuntimeErrors(t *testing.T) {
	if runtime.ErrRuntimeNotFound == nil {
		t.Fatal("ErrRuntimeNotFound should not be nil")
	}
	if runtime.ErrRuntimeUnavailable == nil {
		t.Fatal("ErrRuntimeUnavailable should not be nil")
	}
}

// TestParityRuntimeInfo verifies RuntimeInfo struct fields.
func TestParityRuntimeInfo(t *testing.T) {
	info := runtime.RuntimeInfo{
		Name:      "llm",
		Available: true,
		Type:      "llm_library",
	}
	if info.Name != "llm" {
		t.Fatalf("expected llm got %s", info.Name)
	}
	if !info.Available {
		t.Fatal("expected Available=true")
	}
}

// TestParitySupportedRuntimes verifies the four known runtimes are registered.
func TestParitySupportedRuntimes(t *testing.T) {
	for _, name := range []string{"copilot", "codex", "llm", "gemini"} {
		if !runtime.IsKnownRuntime(name) {
			t.Errorf("expected %s to be a known runtime", name)
		}
	}
}

// TestParityIsKnownRuntimeCaseInsensitive verifies case-insensitive lookup.
func TestParityIsKnownRuntimeCaseInsensitive(t *testing.T) {
	if !runtime.IsKnownRuntime("LLM") {
		t.Fatal("IsKnownRuntime should be case-insensitive")
	}
}

// TestParityGetSupportedRuntimeNames verifies names are returned.
func TestParityGetSupportedRuntimeNames(t *testing.T) {
	names := runtime.GetSupportedRuntimeNames()
	if len(names) < 4 {
		t.Fatalf("expected at least 4 runtime names, got %d", len(names))
	}
}

// TestParityFindRuntimeBinaryRejectsPathTraversal verifies security guard.
func TestParityFindRuntimeBinaryRejectsPathTraversal(t *testing.T) {
	_, err := runtime.FindRuntimeBinary("../etc/passwd")
	if err == nil {
		t.Fatal("expected error for path-traversal name")
	}
	if !strings.Contains(err.Error(), "path separator") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestParityFindRuntimeBinaryRejectsEmpty verifies empty name is rejected.
func TestParityFindRuntimeBinaryRejectsEmpty(t *testing.T) {
	_, err := runtime.FindRuntimeBinary("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

// TestParityFindRuntimeBinaryUnknown returns empty string (not error) for unknown.
func TestParityFindRuntimeBinaryUnknown(t *testing.T) {
	path, err := runtime.FindRuntimeBinary("definitely-not-a-real-binary-apm-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Fatalf("expected empty path for unknown binary, got %s", path)
	}
}

// TestParityGetAvailableRuntimes returns a slice (may be empty in CI).
func TestParityGetAvailableRuntimes(t *testing.T) {
	infos := runtime.GetAvailableRuntimes()
	_ = infos // may be empty in CI -- just verify it doesn't panic
}

// TestParityGetRuntimeByNameUnknown returns ErrRuntimeNotFound.
func TestParityGetRuntimeByNameUnknown(t *testing.T) {
	_, err := runtime.GetRuntimeByName("nonexistent-runtime")
	if err == nil {
		t.Fatal("expected error for unknown runtime")
	}
}

// TestParityIsRuntimeBinaryAvailableFalse returns false for nonexistent binary.
func TestParityIsRuntimeBinaryAvailableFalse(t *testing.T) {
	if runtime.IsRuntimeBinaryAvailable("definitely-not-installed-apm-test") {
		t.Fatal("expected false for unavailable binary")
	}
}
