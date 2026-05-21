package vscode

import (
	"strings"
	"testing"
)

func TestTargetName_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "vscode" {
		t.Errorf("expected vscode, got %s", a.TargetName())
	}
}

func TestMCPServersKey_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "servers" {
		t.Errorf("expected servers, got %s", a.MCPServersKey())
	}
}

func TestSupportsUserScope_False_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.SupportsUserScope() {
		t.Error("vscode should not support user scope")
	}
}

func TestGetConfigPath_ContainsVSCode_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, "vscode") && !strings.Contains(p, ".vscode") {
		t.Errorf("expected vscode in path, got %s", p)
	}
}

func TestGetConfigPath_EndsWithJSON_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, ".json") {
		t.Errorf("expected .json suffix, got %s", p)
	}
}

func TestGetCurrentConfig_Missing_Extra4(t *testing.T) {
	a := New("/nonexistent/xyzabc", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for missing config")
	}
}

func TestNew_NotNil_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a == nil {
		t.Error("expected non-nil adapter")
	}
}

func TestTranslateEnvValueForVSCode_DollarBrace_Extra4(t *testing.T) {
	got := translateEnvValueForVSCode("${MY_VAR}")
	if got == "" {
		t.Error("expected non-empty translation")
	}
}

func TestTranslateEnvValueForVSCode_Plain_Extra4(t *testing.T) {
	got := translateEnvValueForVSCode("plain-value")
	if got != "plain-value" {
		t.Errorf("expected plain-value, got %s", got)
	}
}

func TestTranslateEnvVarsForVSCode_Empty_Extra4(t *testing.T) {
	out := translateEnvVarsForVSCode(map[string]interface{}{})
	if out == nil {
		t.Error("expected non-nil map")
	}
}

func TestTranslateEnvVarsForVSCode_StringVal_Extra4(t *testing.T) {
	env := map[string]interface{}{"KEY": "value"}
	out := translateEnvVarsForVSCode(env)
	if out["KEY"] == nil {
		t.Error("expected KEY in output")
	}
}

func TestExtractInputVariables_Empty_Extra4(t *testing.T) {
	vars := extractInputVariables(map[string]interface{}{}, "srv")
	_ = vars // may return nil, just verify no panic
}

func TestFilterOut_RemovesAll_Extra4(t *testing.T) {
	out := filterOut([]string{"-y", "pkg", "-y"}, "-y")
	for _, v := range out {
		if v == "-y" {
			t.Error("expected -y removed")
		}
	}
}

func TestExtractPackageArgs_Empty_Extra4(t *testing.T) {
	args := extractPackageArgs(map[string]interface{}{})
	_ = args // may return nil, just verify no panic
}

func TestStrField_Missing_Extra4(t *testing.T) {
	v := strField(map[string]interface{}{}, "nonexistent")
	if v != "" {
		t.Errorf("expected empty string, got %s", v)
	}
}

func TestStrField_Present_Extra4(t *testing.T) {
	v := strField(map[string]interface{}{"k": "v"}, "k")
	if v != "v" {
		t.Errorf("expected v, got %s", v)
	}
}
