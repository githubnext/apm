package copilot

import (
	"strings"
	"testing"
)

func TestTranslateEnvPlaceholder_AngleBracket_Extra4(t *testing.T) {
	got := TranslateEnvPlaceholder("<MY_VAR>")
	if got != "${MY_VAR}" {
		t.Errorf("expected ${MY_VAR}, got %s", got)
	}
}

func TestTranslateEnvPlaceholder_DollarBrace_Extra4(t *testing.T) {
	got := TranslateEnvPlaceholder("${MY_VAR}")
	if got != "${MY_VAR}" {
		t.Errorf("expected ${MY_VAR}, got %s", got)
	}
}

func TestTranslateEnvPlaceholder_EnvColon_Extra4(t *testing.T) {
	got := TranslateEnvPlaceholder("${env:MY_VAR}")
	if got != "${MY_VAR}" {
		t.Errorf("expected ${MY_VAR}, got %s", got)
	}
}

func TestTranslateEnvPlaceholder_NoChange_Extra4(t *testing.T) {
	got := TranslateEnvPlaceholder("plain-string")
	if got != "plain-string" {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestTranslateEnvPlaceholder_Empty_Extra4(t *testing.T) {
	got := TranslateEnvPlaceholder("")
	if got != "" {
		t.Errorf("expected empty, got %s", got)
	}
}

func TestExtractLegacyAngleVars_Single_Extra4(t *testing.T) {
	vars := ExtractLegacyAngleVars("use <API_KEY> here")
	if len(vars) != 1 || vars[0] != "API_KEY" {
		t.Errorf("expected [API_KEY], got %v", vars)
	}
}

func TestExtractLegacyAngleVars_None_Extra4(t *testing.T) {
	vars := ExtractLegacyAngleVars("no vars here")
	if len(vars) != 0 {
		t.Errorf("expected empty, got %v", vars)
	}
}

func TestExtractLegacyAngleVars_Multiple_Extra4(t *testing.T) {
	vars := ExtractLegacyAngleVars("<A> and <B_C>")
	if len(vars) != 2 {
		t.Errorf("expected 2, got %v", vars)
	}
}

func TestHasEnvPlaceholder_AngleBracket_Extra4(t *testing.T) {
	if !HasEnvPlaceholder("<MY_VAR>") {
		t.Error("expected true for angle bracket var")
	}
}

func TestHasEnvPlaceholder_DollarBrace_Extra4(t *testing.T) {
	if !HasEnvPlaceholder("${MY_VAR}") {
		t.Error("expected true for dollar brace var")
	}
}

func TestHasEnvPlaceholder_Plain_Extra4(t *testing.T) {
	if HasEnvPlaceholder("plain") {
		t.Error("expected false for plain string")
	}
}

func TestHasEnvPlaceholder_Empty_Extra4(t *testing.T) {
	if HasEnvPlaceholder("") {
		t.Error("expected false for empty")
	}
}

func TestNew_TargetName_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "copilot" {
		t.Errorf("expected copilot, got %s", a.TargetName())
	}
}

func TestNew_MCPServersKey_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %s", a.MCPServersKey())
	}
}

func TestNew_SupportsUserScope_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope true")
	}
}

func TestGetConfigPath_ContainsCopilot_Extra4(t *testing.T) {
	a := New("/tmp", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, "copilot") && !strings.Contains(p, ".copilot") {
		t.Errorf("expected copilot in path, got %s", p)
	}
}

func TestGetConfigPath_EndsWithJSON_Extra4(t *testing.T) {
	a := New("/tmp", false)
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, ".json") {
		t.Errorf("expected .json suffix, got %s", p)
	}
}

func TestGetCurrentConfig_MissingDir_Extra4(t *testing.T) {
	a := New("/nonexistent/xyz/abc", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for missing config")
	}
}

func TestResetInstallRunState_NoError_Extra4(t *testing.T) {
	ResetInstallRunState()
}

func TestFormatResolveEnv_EmptyOverrides_Extra4(t *testing.T) {
	a := New("/tmp", false)
	env := map[string]interface{}{"KEY": "val"}
	resolved := a.FormatResolveEnv(env, nil)
	if _, ok := resolved["KEY"]; !ok {
		t.Error("expected KEY in resolved map")
	}
}

func TestFormatResolveEnv_Nil_Extra4(t *testing.T) {
	a := New("/tmp", false)
	resolved := a.FormatResolveEnv(nil, nil)
	if resolved == nil {
		t.Error("expected non-nil map")
	}
}

func TestFormatProcessArgs_Empty_Extra4(t *testing.T) {
	a := New("/tmp", false)
	out := a.FormatProcessArgs(nil, nil, nil)
	if out == nil {
		t.Error("expected non-nil slice")
	}
}

func TestFormatProcessArgs_Passthrough_Extra4(t *testing.T) {
	a := New("/tmp", false)
	args := []string{"--flag", "value"}
	out := a.FormatProcessArgs(args, nil, nil)
	if len(out) != 2 {
		t.Errorf("expected 2 args, got %d", len(out))
	}
}
