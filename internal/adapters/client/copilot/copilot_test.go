package copilot

import (
	"testing"
)

func TestTranslateEnvPlaceholder(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"${MY_TOKEN}", "${MY_TOKEN}"},
		{"<MY_TOKEN>", "${MY_TOKEN}"},
		{"plain-string", "plain-string"},
		{"", ""},
		{"<TOKEN_A> and <TOKEN_B>", "${TOKEN_A} and ${TOKEN_B}"},
	}
	for _, tc := range cases {
		got := TranslateEnvPlaceholder(tc.in)
		if got != tc.want {
			t.Errorf("TranslateEnvPlaceholder(%q): got %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestHasEnvPlaceholder(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"${MY_TOKEN}", true},
		{"<MY_TOKEN>", true},
		{"plain-string", false},
		{"", false},
		{"prefix${VAR}suffix", true},
		{"prefix<VAR>suffix", true},
	}
	for _, tc := range cases {
		got := HasEnvPlaceholder(tc.in)
		if got != tc.want {
			t.Errorf("HasEnvPlaceholder(%q): got %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestExtractLegacyAngleVars(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"<MY_TOKEN>", []string{"MY_TOKEN"}},
		{"<A> and <B>", []string{"A", "B"}},
		{"${VAR}", nil},
		{"no vars here", nil},
		{"<TOKEN_1> ${VAR2} <TOKEN_3>", []string{"TOKEN_1", "TOKEN_3"}},
	}
	for _, tc := range cases {
		got := ExtractLegacyAngleVars(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("ExtractLegacyAngleVars(%q): got %v, want %v", tc.in, got, tc.want)
			continue
		}
		for i, g := range got {
			if g != tc.want[i] {
				t.Errorf("ExtractLegacyAngleVars(%q)[%d]: got %q, want %q", tc.in, i, g, tc.want[i])
			}
		}
	}
}

func TestNew(t *testing.T) {
	a := New("/repo", false)
	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.TargetName() != "copilot" {
		t.Errorf("TargetName: got %q, want copilot", a.TargetName())
	}
	if a.MCPServersKey() != "mcpServers" {
		t.Errorf("MCPServersKey: got %q", a.MCPServersKey())
	}
	if !a.SupportsUserScope() {
		t.Error("SupportsUserScope should be true")
	}
}

func TestGetConfigPathUserScope(t *testing.T) {
	a := New("/repo", true)
	path := a.GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}
}

func TestGetConfigPathProjectScope(t *testing.T) {
	a := New("/my/project", false)
	path := a.GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}
}

func TestResetInstallRunState(t *testing.T) {
	// Just verify it does not panic
	ResetInstallRunState()
	ResetInstallRunState()
}
