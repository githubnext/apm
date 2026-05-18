package codexruntime

import (
	"strings"
	"testing"
)

func TestGetRuntimeName(t *testing.T) {
	r := &CodexRuntime{ModelName: "gpt-4"}
	if got := r.GetRuntimeName(); got != "codex" {
		t.Errorf("GetRuntimeName() = %q, want %q", got, "codex")
	}
}

func TestString(t *testing.T) {
	r := &CodexRuntime{ModelName: "gpt-4"}
	s := r.String()
	if !strings.Contains(s, "gpt-4") {
		t.Errorf("String() = %q, want to contain model name", s)
	}
	if !strings.Contains(s, "CodexRuntime") {
		t.Errorf("String() = %q, want to contain CodexRuntime", s)
	}
}

func TestGetRuntimeInfo(t *testing.T) {
	r := &CodexRuntime{ModelName: "gpt-4"}
	info := r.GetRuntimeInfo()
	if info["name"] != "codex" {
		t.Errorf("GetRuntimeInfo()['name'] = %v, want %q", info["name"], "codex")
	}
	if info["type"] != "codex_cli" {
		t.Errorf("GetRuntimeInfo()['type'] = %v, want %q", info["type"], "codex_cli")
	}
}

func TestListAvailableModels(t *testing.T) {
	r := &CodexRuntime{ModelName: "default"}
	models := r.ListAvailableModels()
	if len(models) == 0 {
		t.Error("ListAvailableModels() returned empty map")
	}
}

func TestNewDefaultModelName(t *testing.T) {
	// IsAvailable returns false in sandbox; verify New() sets default model.
	if !IsAvailable() {
		r := &CodexRuntime{ModelName: ""}
		_, err := New("")
		if err == nil {
			t.Error("Expected error when codex not available")
		}
		_ = r
		return
	}
	r, err := New("")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if r.ModelName != "default" {
		t.Errorf("ModelName = %q, want %q", r.ModelName, "default")
	}
}

func TestCodexRuntime_ZeroValue(t *testing.T) {
	r := &CodexRuntime{}
	if r.GetRuntimeName() != "codex" {
		t.Errorf("GetRuntimeName() = %q, want codex", r.GetRuntimeName())
	}
}

func TestGetRuntimeInfo_Keys(t *testing.T) {
	r := &CodexRuntime{ModelName: "gpt-4"}
	info := r.GetRuntimeInfo()
	if _, ok := info["name"]; !ok {
		t.Error("GetRuntimeInfo() should have 'name' key")
	}
	if _, ok := info["type"]; !ok {
		t.Error("GetRuntimeInfo() should have 'type' key")
	}
}

func TestListAvailableModels_NonEmpty(t *testing.T) {
	r := &CodexRuntime{ModelName: "any"}
	models := r.ListAvailableModels()
	if len(models) == 0 {
		t.Error("ListAvailableModels() should return non-empty map")
	}
}

func TestString_ContainsModelName(t *testing.T) {
	r := &CodexRuntime{ModelName: "o1-mini"}
	s := r.String()
	if !strings.Contains(s, "o1-mini") {
		t.Errorf("String() = %q, expected to contain model name", s)
	}
}

func TestNewDefault_WhenUnavailable(t *testing.T) {
	if IsAvailable() {
		t.Skip("codex available, skipping unavailable test")
	}
	_, err := NewDefault()
	if err == nil {
		t.Error("NewDefault() should return error when codex unavailable")
	}
}

func TestGetRuntimeName_Const(t *testing.T) {
	r1 := &CodexRuntime{ModelName: "a"}
	r2 := &CodexRuntime{ModelName: "b"}
	if r1.GetRuntimeName() != r2.GetRuntimeName() {
		t.Error("GetRuntimeName() should be the same across instances")
	}
}
