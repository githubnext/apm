package copilot_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

func TestTranslateEnvPlaceholder_BracesPassthrough(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("${FOO}")
	if got != "${FOO}" {
		t.Errorf("expected ${FOO}, got %q", got)
	}
}

func TestTranslateEnvPlaceholder_EmptyString(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestTranslateEnvPlaceholder_SingleAngle(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("<TOKEN>")
	if got == "" {
		t.Error("expected non-empty result for angle-bracket placeholder")
	}
}

func TestTranslateEnvPlaceholder_NoPlaceholder(t *testing.T) {
	got := copilot.TranslateEnvPlaceholder("hello world")
	if got != "hello world" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestExtractLegacyAngleVars_MultipleVars(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("hello <FOO> and <BAR>")
	if len(vars) != 2 {
		t.Errorf("expected 2 vars, got %d: %v", len(vars), vars)
	}
}

func TestExtractLegacyAngleVars_NoVars(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("no angle vars here")
	if len(vars) != 0 {
		t.Errorf("expected 0 vars, got %d", len(vars))
	}
}

func TestExtractLegacyAngleVars_NoBraces(t *testing.T) {
	vars := copilot.ExtractLegacyAngleVars("${BRACE} only")
	if len(vars) != 0 {
		t.Errorf("expected 0 angle vars for brace-only, got %d", len(vars))
	}
}

func TestHasEnvPlaceholder_BothFormats(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"${VAR}", true},
		{"<VAR>", true},
		{"plain text", false},
		{"", false},
	}
	for _, c := range cases {
		got := copilot.HasEnvPlaceholder(c.in)
		if got != c.want {
			t.Errorf("HasEnvPlaceholder(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestNew_SupportsUserScope(t *testing.T) {
	a := copilot.New("/tmp/proj", false)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope=true")
	}
}

func TestNew_TargetName(t *testing.T) {
	a := copilot.New("/tmp/proj", false)
	if a.TargetName() != "copilot" {
		t.Errorf("expected TargetName=copilot, got %q", a.TargetName())
	}
}

func TestNew_MCPServersKey(t *testing.T) {
	a := copilot.New("/tmp/proj", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("expected mcpServers, got %q", a.MCPServersKey())
	}
}

func TestGetCurrentConfig_MissingDir(t *testing.T) {
	a := copilot.New("/nonexistent/path/xyz", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil config map even for missing dir")
	}
}

func TestResetInstallRunState_NoopOnSecondCall(t *testing.T) {
	copilot.ResetInstallRunState()
	copilot.ResetInstallRunState()
}

func TestGetConfigPath_Contains_CopilotOrVSCode(t *testing.T) {
	a := copilot.New("/tmp/myproject", false)
	p := a.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty config path")
	}
}

func TestGetConfigPath_UserScope_ContainsHome(t *testing.T) {
	a := copilot.New("", true)
	p := a.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty user-scope config path")
	}
}

func TestNew_ProjectRoot_Stored(t *testing.T) {
	a := copilot.New("/my/project", false)
	if a.ProjectRoot != "/my/project" {
		t.Errorf("expected ProjectRoot=/my/project, got %q", a.ProjectRoot)
	}
}

func TestNew_UserScope_Stored(t *testing.T) {
	a := copilot.New("/x", true)
	if !a.UserScope {
		t.Error("expected UserScope=true")
	}
}

func TestGetCurrentConfig_EmptyDirReturnsMap(t *testing.T) {
	a := copilot.New("", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map")
	}
}
