package codexruntime

import (
	"strings"
	"testing"
)

func TestGetRuntimeInfo_CapabilitiesKey(t *testing.T) {
	r := &CodexRuntime{ModelName: "test-model"}
	info := r.GetRuntimeInfo()
	if _, ok := info["capabilities"]; !ok {
		t.Error("expected capabilities key in runtime info")
	}
}

func TestGetRuntimeInfo_CapabilitiesModelExecution(t *testing.T) {
	r := &CodexRuntime{ModelName: "test-model"}
	info := r.GetRuntimeInfo()
	caps, _ := info["capabilities"].(map[string]interface{})
	if caps == nil {
		t.Fatal("capabilities should be a map")
	}
	modelExec, _ := caps["model_execution"].(bool)
	if !modelExec {
		t.Error("model_execution should be true")
	}
}

func TestListAvailableModels_HasProvider(t *testing.T) {
	r := &CodexRuntime{ModelName: ""}
	models := r.ListAvailableModels()
	for k, v := range models {
		m, _ := v.(map[string]string)
		if m == nil {
			t.Errorf("model %q value should be map[string]string", k)
			continue
		}
		if m["provider"] == "" {
			t.Errorf("model %q missing provider field", k)
		}
	}
}

func TestListAvailableModels_HasDescription(t *testing.T) {
	r := &CodexRuntime{ModelName: ""}
	models := r.ListAvailableModels()
	for k, v := range models {
		m, _ := v.(map[string]string)
		if m == nil {
			t.Errorf("model %q value should be map[string]string", k)
			continue
		}
		if m["description"] == "" {
			t.Errorf("model %q missing description", k)
		}
	}
}

func TestString_ContainsModelField(t *testing.T) {
	r := &CodexRuntime{ModelName: "my-model"}
	s := r.String()
	if !strings.Contains(s, "my-model") {
		t.Errorf("String() should contain model name: %q", s)
	}
}

func TestString_FormatCheck(t *testing.T) {
	r := &CodexRuntime{ModelName: "model-x"}
	s := r.String()
	if !strings.HasPrefix(s, "CodexRuntime(") {
		t.Errorf("String() should start with 'CodexRuntime(': %q", s)
	}
	if !strings.HasSuffix(s, ")") {
		t.Errorf("String() should end with ')': %q", s)
	}
}

func TestGetRuntimeInfo_NameIsCodex(t *testing.T) {
	r := &CodexRuntime{}
	info := r.GetRuntimeInfo()
	if info["name"] != "codex" {
		t.Errorf("name should be 'codex': %v", info["name"])
	}
}

func TestGetRuntimeInfo_TypeIsCLI(t *testing.T) {
	r := &CodexRuntime{}
	info := r.GetRuntimeInfo()
	if info["type"] != "codex_cli" {
		t.Errorf("type should be 'codex_cli': %v", info["type"])
	}
}

func TestCodexRuntime_ModelNamePreserved(t *testing.T) {
	r := &CodexRuntime{ModelName: "special-model-123"}
	if r.ModelName != "special-model-123" {
		t.Errorf("ModelName: %q", r.ModelName)
	}
}

func TestListAvailableModels_HasID(t *testing.T) {
	r := &CodexRuntime{}
	models := r.ListAvailableModels()
	for k, v := range models {
		m, _ := v.(map[string]string)
		if m["id"] == "" {
			t.Errorf("model %q missing id field", k)
		}
	}
}

func TestGetRuntimeInfo_VersionFieldPresent(t *testing.T) {
	r := &CodexRuntime{ModelName: ""}
	info := r.GetRuntimeInfo()
	if _, ok := info["version"]; !ok {
		t.Error("version key should be present in runtime info")
	}
}
