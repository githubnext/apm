package llmruntime

import (
	"strings"
	"testing"
)

func TestLLMRuntime_FieldsDirectAccess(t *testing.T) {
	r := &LLMRuntime{ModelName: "gpt-4"}
	if r.ModelName != "gpt-4" {
		t.Errorf("expected ModelName gpt-4, got %q", r.ModelName)
	}
}

func TestLLMRuntime_ZeroValue(t *testing.T) {
	var r LLMRuntime
	if r.ModelName != "" {
		t.Errorf("zero-value ModelName should be empty, got %q", r.ModelName)
	}
}

func TestGetRuntimeName_IsLLM(t *testing.T) {
	cases := []string{"gpt-4", "claude-3", "", "my-model"}
	for _, m := range cases {
		r := &LLMRuntime{ModelName: m}
		got := r.GetRuntimeName()
		if got != "llm" {
			t.Errorf("GetRuntimeName() with model %q: expected 'llm', got %q", m, got)
		}
	}
}

func TestString_Format(t *testing.T) {
	r := &LLMRuntime{ModelName: "claude-3"}
	s := r.String()
	if !strings.Contains(s, "claude-3") {
		t.Errorf("String() should contain model name 'claude-3', got %q", s)
	}
}

func TestGetRuntimeInfo_HasCapabilities(t *testing.T) {
	r := &LLMRuntime{ModelName: "test-model"}
	info := r.GetRuntimeInfo()
	if _, ok := info["capabilities"]; !ok {
		t.Error("GetRuntimeInfo should include capabilities key")
	}
}

func TestGetRuntimeInfo_HasType(t *testing.T) {
	r := &LLMRuntime{ModelName: "test-model"}
	info := r.GetRuntimeInfo()
	if _, ok := info["type"]; !ok {
		t.Error("GetRuntimeInfo should include type key")
	}
}

func TestGetRuntimeInfo_HasDescription(t *testing.T) {
	r := &LLMRuntime{ModelName: "test-model"}
	info := r.GetRuntimeInfo()
	if _, ok := info["description"]; !ok {
		t.Error("GetRuntimeInfo should include description key")
	}
}

func TestListAvailableModels_NoError(t *testing.T) {
	r := &LLMRuntime{ModelName: "test-model"}
	models := r.ListAvailableModels()
	// Should return a map (possibly empty) without panicking
	if models == nil {
		t.Log("ListAvailableModels returned nil (acceptable when llm binary not present)")
	}
}

func TestGetRuntimeInfo_TypeValue(t *testing.T) {
	r := &LLMRuntime{}
	info := r.GetRuntimeInfo()
	typ, ok := info["type"]
	if !ok {
		t.Fatal("missing 'type' key")
	}
	if typ == nil {
		t.Error("type value should not be nil")
	}
}

func TestString_ContainsRuntimeLabel(t *testing.T) {
	r := &LLMRuntime{ModelName: "my-model"}
	s := r.String()
	lower := strings.ToLower(s)
	if !strings.Contains(lower, "llm") && !strings.Contains(lower, "runtime") {
		t.Errorf("String() should mention llm or runtime, got %q", s)
	}
}

func TestNewWhenAvailable_FallsThrough(t *testing.T) {
	// When llm binary is not installed, New returns an error.
	// When it is installed, it returns a valid runtime.
	// Either outcome is acceptable in CI.
	r, err := New("some-model")
	if err != nil {
		// llm binary not present -- expected in CI
		return
	}
	if r == nil {
		t.Error("expected non-nil runtime when error is nil")
	}
}

func TestNewDefault_ReturnTypeOrError(t *testing.T) {
	r, err := NewDefault()
	if err != nil {
		return // acceptable
	}
	if r == nil {
		t.Error("NewDefault: non-nil error but nil runtime")
	}
	if r.ModelName != "" {
		t.Log("NewDefault ModelName:", r.ModelName)
	}
}
