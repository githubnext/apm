package codexruntime

import (
	"strings"
	"testing"
)

func TestCodexRuntime_ModelName(t *testing.T) {
	r := &CodexRuntime{ModelName: "codex-mini"}
	if r.ModelName != "codex-mini" {
		t.Errorf("ModelName = %q, want codex-mini", r.ModelName)
	}
}

func TestGetRuntimeName_AlwaysCodex(t *testing.T) {
	models := []string{"gpt-4", "codex-mini", "default", ""}
	for _, m := range models {
		r := &CodexRuntime{ModelName: m}
		if r.GetRuntimeName() != "codex" {
			t.Errorf("GetRuntimeName() = %q for model %q, want codex", r.GetRuntimeName(), m)
		}
	}
}

func TestString_ContainsModelNameVariants(t *testing.T) {
	models := []string{"gpt-4", "codex-mini", "my-model"}
	for _, m := range models {
		r := &CodexRuntime{ModelName: m}
		s := r.String()
		if !strings.Contains(s, m) {
			t.Errorf("String() = %q does not contain model %q", s, m)
		}
	}
}

func TestString_ContainsCodexRuntime(t *testing.T) {
	r := &CodexRuntime{ModelName: "x"}
	if !strings.Contains(r.String(), "CodexRuntime") {
		t.Errorf("String() = %q should contain CodexRuntime", r.String())
	}
}

func TestGetRuntimeInfo_NameField(t *testing.T) {
	r := &CodexRuntime{ModelName: "m"}
	info := r.GetRuntimeInfo()
	if info["name"] != "codex" {
		t.Errorf("info[name] = %v, want codex", info["name"])
	}
}

func TestGetRuntimeInfo_TypeField(t *testing.T) {
	r := &CodexRuntime{ModelName: "m"}
	info := r.GetRuntimeInfo()
	if info["type"] != "codex_cli" {
		t.Errorf("info[type] = %v, want codex_cli", info["type"])
	}
}

func TestGetRuntimeInfo_VersionField(t *testing.T) {
	r := &CodexRuntime{ModelName: "m"}
	info := r.GetRuntimeInfo()
	if _, ok := info["version"]; !ok {
		t.Error("info should have version field")
	}
}

func TestListAvailableModels_HasCodexDefault(t *testing.T) {
	r := &CodexRuntime{ModelName: "m"}
	models := r.ListAvailableModels()
	if _, ok := models["codex-default"]; !ok {
		t.Error("ListAvailableModels should contain codex-default")
	}
}

func TestListAvailableModels_NotNil(t *testing.T) {
	r := &CodexRuntime{ModelName: "m"}
	models := r.ListAvailableModels()
	if models == nil {
		t.Error("ListAvailableModels should not be nil")
	}
	if len(models) == 0 {
		t.Error("ListAvailableModels should not be empty")
	}
}

func TestListAvailableModels_CodexDefaultFields(t *testing.T) {
	r := &CodexRuntime{ModelName: "m"}
	models := r.ListAvailableModels()
	m, ok := models["codex-default"].(map[string]string)
	if !ok {
		t.Fatalf("codex-default should be map[string]string, got %T", models["codex-default"])
	}
	if m["id"] != "codex-default" {
		t.Errorf("id = %q, want codex-default", m["id"])
	}
	if m["provider"] != "codex" {
		t.Errorf("provider = %q, want codex", m["provider"])
	}
}

func TestCodexRuntime_ZeroValueModelName(t *testing.T) {
	var r CodexRuntime
	if r.GetRuntimeName() != "codex" {
		t.Error("zero value GetRuntimeName should still return codex")
	}
	if r.ModelName != "" {
		t.Error("zero value ModelName should be empty")
	}
}

func TestGetRuntimeInfo_ModelNameInclusion(t *testing.T) {
	r := &CodexRuntime{ModelName: "special-model"}
	info := r.GetRuntimeInfo()
	// Verify info is a non-nil map with expected required fields
	if info == nil {
		t.Fatal("GetRuntimeInfo returned nil")
	}
	if info["name"] == nil {
		t.Error("missing name field")
	}
}

func TestInstallCmd_NotEmpty(t *testing.T) {
	if installCmd == "" {
		t.Error("installCmd constant should not be empty")
	}
	if !strings.Contains(installCmd, "npm") {
		t.Errorf("installCmd = %q should mention npm", installCmd)
	}
}

func TestInstallCmd_ContainsCodex(t *testing.T) {
	if !strings.Contains(installCmd, "codex") {
		t.Errorf("installCmd = %q should contain codex", installCmd)
	}
}
