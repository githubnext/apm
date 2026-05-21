package copilot_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

func TestTranslateEnvPlaceholder_BracketForm_Extra3(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("${env:TOKEN}")
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestTranslateEnvPlaceholder_LongName_Extra3(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("<VERY_LONG_ENV_NAME_HERE>")
	if got == "" {
		t.Error("expected non-empty")
	}
}

func TestTranslateEnvPlaceholder_TwoAngle_Extra3(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("<A> and <B>")
	_ = got // just ensure no panic
}

func TestExtractLegacyAngleVars_Single_Extra3(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("<FOO>")
	if len(vars) != 1 || vars[0] != "FOO" {
		t.Errorf("expected [FOO], got %v", vars)
	}
}

func TestExtractLegacyAngleVars_Multiple_Extra3(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("use <A> and <B>")
	if len(vars) != 2 {
		t.Errorf("expected 2, got %d", len(vars))
	}
}

func TestExtractLegacyAngleVars_Lowercase_Extra3(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("<lower>")
	if len(vars) != 0 {
		t.Errorf("expected no match for lowercase, got %v", vars)
	}
}

func TestHasEnvPlaceholder_AngleBracket_Extra3(t *testing.T) {
	if !copilot.HasEnvPlaceholder("<TOKEN>") {
		t.Error("expected true for <TOKEN>")
	}
}

func TestHasEnvPlaceholder_DollarBrace_Extra3(t *testing.T) {
	if !copilot.HasEnvPlaceholder("${MY_VAR}") {
		t.Error("expected true for ${MY_VAR}")
	}
}

func TestHasEnvPlaceholder_PlainString_Extra3(t *testing.T) {
	if copilot.HasEnvPlaceholder("no placeholders here") {
		t.Error("expected false")
	}
}

func TestNew_TargetName_Extra3(t *testing.T) {
	a := copilot.New("/root", false)
	if a.TargetName() != "copilot" {
		t.Errorf("expected copilot, got %q", a.TargetName())
	}
}

func TestNew_MCPServersKey_Extra3(t *testing.T) {
	a := copilot.New("/root", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %q", a.MCPServersKey())
	}
}

func TestNew_SupportsUserScope_Extra3(t *testing.T) {
	a := copilot.New("", false)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope true")
	}
}

func TestGetConfigPath_NotEmpty_Extra3(t *testing.T) {
	a := copilot.New("/root", false)
	if a.GetConfigPath() == "" {
		t.Error("expected non-empty config path")
	}
}

func TestResetInstallRunState_NoError_Extra3(t *testing.T) {
	// Ensure ResetInstallRunState doesn't panic
	copilot.ResetInstallRunState()
	copilot.ResetInstallRunState()
}

func TestExtractLegacyAngleVars_WithDigits_Extra3(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("<VAR_1>")
	if len(vars) == 0 {
		t.Error("expected match for <VAR_1>")
	}
}

func TestHasEnvPlaceholder_EnvPrefix_Extra3(t *testing.T) {
	if !copilot.HasEnvPlaceholder("${env:SECRET}") {
		t.Error("expected true for ${env:SECRET}")
	}
}
