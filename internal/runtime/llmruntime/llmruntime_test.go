package llmruntime

import (
	"strings"
	"testing"
)

func TestGetRuntimeInfo_Capabilities(t *testing.T) {
	r := &LLMRuntime{ModelName: "gpt-4"}
	info := r.GetRuntimeInfo()
	caps, ok := info["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatalf("capabilities field not a map: %T", info["capabilities"])
	}
	if caps["model_execution"] != true {
		t.Error("model_execution capability should be true")
	}
}

func TestGetRuntimeInfo_Type(t *testing.T) {
	r := &LLMRuntime{ModelName: "claude"}
	info := r.GetRuntimeInfo()
	if info["type"] != "llm_library" {
		t.Errorf("type = %v, want llm_library", info["type"])
	}
}

func TestGetRuntimeInfo_Description(t *testing.T) {
	r := &LLMRuntime{}
	info := r.GetRuntimeInfo()
	desc, ok := info["description"].(string)
	if !ok || desc == "" {
		t.Error("description should be non-empty string")
	}
}

func TestString_EmptyModel(t *testing.T) {
	r := &LLMRuntime{ModelName: ""}
	s := r.String()
	if !strings.Contains(s, "LLMRuntime") {
		t.Errorf("String() = %q, want to contain LLMRuntime", s)
	}
}

func TestLLMRuntimeStruct_Fields(t *testing.T) {
	r := &LLMRuntime{ModelName: "gemini-pro"}
	if r.ModelName != "gemini-pro" {
		t.Errorf("ModelName = %q, want gemini-pro", r.ModelName)
	}
}

func TestGetRuntimeName(t *testing.T) {
	r := &LLMRuntime{ModelName: "gpt-4"}
	if got := r.GetRuntimeName(); got != "llm" {
		t.Errorf("GetRuntimeName() = %q, want %q", got, "llm")
	}
}

func TestGetRuntimeInfo(t *testing.T) {
	r := &LLMRuntime{ModelName: "gpt-4"}
	info := r.GetRuntimeInfo()
	if info["name"] != "llm" {
		t.Errorf("GetRuntimeInfo()['name'] = %v, want %q", info["name"], "llm")
	}
	if info["current_model"] != "gpt-4" {
		t.Errorf("GetRuntimeInfo()['current_model'] = %v, want %q", info["current_model"], "gpt-4")
	}
}

func TestGetRuntimeInfoDefaultModel(t *testing.T) {
	r := &LLMRuntime{ModelName: ""}
	info := r.GetRuntimeInfo()
	if info["current_model"] != "default" {
		t.Errorf("GetRuntimeInfo()['current_model'] = %v, want %q", info["current_model"], "default")
	}
}

func TestString(t *testing.T) {
	r := &LLMRuntime{ModelName: "claude-3"}
	s := r.String()
	if !strings.Contains(s, "claude-3") {
		t.Errorf("String() = %q, want to contain model name", s)
	}
	if !strings.Contains(s, "LLMRuntime") {
		t.Errorf("String() = %q, want to contain LLMRuntime", s)
	}
}

func TestNewWhenNotAvailable(t *testing.T) {
	if IsAvailable() {
		t.Skip("llm binary present, skipping unavailability test")
	}
	_, err := New("gpt-4")
	if err == nil {
		t.Error("Expected error when llm not available")
	}
}
