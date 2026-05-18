package copilot

import (
	"testing"
)

func TestTranslateEnvPlaceholder_MultipleAngles(t *testing.T) {
	got := TranslateEnvPlaceholder("<A> <B> <C>")
	want := "${A} ${B} ${C}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTranslateEnvPlaceholder_AlreadyBraces(t *testing.T) {
	got := TranslateEnvPlaceholder("${ALREADY}")
	if got != "${ALREADY}" {
		t.Errorf("got %q", got)
	}
}

func TestHasEnvPlaceholder_MixedFormats(t *testing.T) {
	if !HasEnvPlaceholder("<TOKEN> and some text") {
		t.Error("expected true for angle-bracket placeholder")
	}
	if !HasEnvPlaceholder("some ${VAR} text") {
		t.Error("expected true for brace placeholder")
	}
}

func TestExtractLegacyAngleVars_EmptyString(t *testing.T) {
	got := ExtractLegacyAngleVars("")
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestExtractLegacyAngleVars_BracesIgnored(t *testing.T) {
	got := ExtractLegacyAngleVars("${VAR1} ${VAR2}")
	if len(got) != 0 {
		t.Errorf("expected no angle vars for brace placeholders, got %v", got)
	}
}

func TestNew_ProjectScope(t *testing.T) {
	a := New("/some/path", false)
	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.TargetName() != "copilot" {
		t.Errorf("TargetName: %q", a.TargetName())
	}
}

func TestNew_MCPServersKey(t *testing.T) {
	a := New("/repo", false)
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("MCPServersKey: %q", a.MCPServersKey())
	}
}

func TestNew_UserScope(t *testing.T) {
	a := New("/repo", true)
	if !a.SupportsUserScope() {
		t.Error("SupportsUserScope should be true")
	}
}

func TestGetConfigPath_NonEmpty(t *testing.T) {
	cases := []struct {
		root      string
		userScope bool
	}{
		{"/project", false},
		{"/home/user", true},
		{"", false},
	}
	for _, tc := range cases {
		a := New(tc.root, tc.userScope)
		path := a.GetConfigPath()
		if path == "" {
			t.Errorf("GetConfigPath returned empty for root=%q userScope=%v", tc.root, tc.userScope)
		}
	}
}

func TestResetInstallRunState_MultipleReset(t *testing.T) {
	ResetInstallRunState()
	ResetInstallRunState()
	ResetInstallRunState()
}

func TestTranslateEnvPlaceholder_NoSpecialChars(t *testing.T) {
	got := TranslateEnvPlaceholder("just a plain string with no vars")
	if got != "just a plain string with no vars" {
		t.Errorf("expected unchanged, got %q", got)
	}
}
