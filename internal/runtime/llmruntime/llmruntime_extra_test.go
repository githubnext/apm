package llmruntime

import (
"strings"
"testing"
)

func TestGetRuntimeInfo_HasName(t *testing.T) {
r := &LLMRuntime{ModelName: "gpt-4o"}
info := r.GetRuntimeInfo()
name, ok := info["name"].(string)
if !ok || name != "llm" {
t.Errorf("GetRuntimeInfo name = %v, want 'llm'", info["name"])
}
}

func TestGetRuntimeInfo_CurrentModelSet(t *testing.T) {
r := &LLMRuntime{ModelName: "claude-3-opus"}
info := r.GetRuntimeInfo()
if info["current_model"] != "claude-3-opus" {
t.Errorf("current_model = %v, want claude-3-opus", info["current_model"])
}
}

func TestGetRuntimeInfo_EmptyModelDefaultsToDefault(t *testing.T) {
r := &LLMRuntime{ModelName: ""}
info := r.GetRuntimeInfo()
if info["current_model"] != "default" {
t.Errorf("current_model = %v, want 'default'", info["current_model"])
}
}

func TestGetRuntimeInfo_TypeIsLLMLibrary(t *testing.T) {
r := &LLMRuntime{}
info := r.GetRuntimeInfo()
if info["type"] != "llm_library" {
t.Errorf("type = %v, want llm_library", info["type"])
}
}

func TestGetRuntimeInfo_DescriptionNonEmpty(t *testing.T) {
r := &LLMRuntime{ModelName: "x"}
info := r.GetRuntimeInfo()
desc, ok := info["description"].(string)
if !ok || desc == "" {
t.Error("description should be non-empty string")
}
}

func TestGetRuntimeInfo_CapabilitiesMap(t *testing.T) {
r := &LLMRuntime{ModelName: "m"}
info := r.GetRuntimeInfo()
caps, ok := info["capabilities"].(map[string]interface{})
if !ok {
t.Fatalf("capabilities should be map, got %T", info["capabilities"])
}
if caps["model_execution"] != true {
t.Error("model_execution capability should be true")
}
}

func TestGetRuntimeName_AlwaysLLM(t *testing.T) {
for _, model := range []string{"", "gpt-4", "claude", "gemini-pro"} {
r := &LLMRuntime{ModelName: model}
if got := r.GetRuntimeName(); got != "llm" {
t.Errorf("GetRuntimeName(%q) = %q, want llm", model, got)
}
}
}

func TestString_ContainsModelName(t *testing.T) {
r := &LLMRuntime{ModelName: "my-special-model"}
s := r.String()
if !strings.Contains(s, "my-special-model") {
t.Errorf("String() = %q, should contain model name", s)
}
}

func TestString_ContainsLLMRuntime(t *testing.T) {
r := &LLMRuntime{ModelName: ""}
s := r.String()
if !strings.Contains(s, "LLMRuntime") {
t.Errorf("String() = %q, should contain 'LLMRuntime'", s)
}
}

func TestLLMRuntime_MultipleInstances(t *testing.T) {
r1 := &LLMRuntime{ModelName: "a"}
r2 := &LLMRuntime{ModelName: "b"}
if r1.GetRuntimeName() != r2.GetRuntimeName() {
t.Error("GetRuntimeName should be the same for all instances")
}
if r1.ModelName == r2.ModelName {
t.Error("instances should have independent ModelName fields")
}
}

func TestNewDefault_NotAvailable(t *testing.T) {
if IsAvailable() {
t.Skip("llm binary present on PATH, skipping unavailability test")
}
_, err := NewDefault()
if err == nil {
t.Error("NewDefault should return error when llm not available")
}
}
