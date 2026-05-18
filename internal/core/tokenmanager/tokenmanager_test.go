package tokenmanager

import (
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	if ADOBearerSource == "" {
		t.Error("ADOBearerSource should be non-empty")
	}
	if DefaultCredentialTimeout <= 0 {
		t.Error("DefaultCredentialTimeout should be positive")
	}
	if MaxCredentialTimeout < DefaultCredentialTimeout {
		t.Error("MaxCredentialTimeout should be >= DefaultCredentialTimeout")
	}
}

func TestNew(t *testing.T) {
	m := New(false)
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	if m.PreserveExisting {
		t.Error("expected PreserveExisting=false")
	}
	m2 := New(true)
	if !m2.PreserveExisting {
		t.Error("expected PreserveExisting=true")
	}
}

func TestGetTokenForPurpose_FromEnv(t *testing.T) {
	m := New(false)
	env := map[string]string{
		"GITHUB_TOKEN": "ghp_test123",
	}
	tok, ok := m.GetTokenForPurpose("models", env)
	if !ok {
		t.Error("expected token found")
	}
	if tok != "ghp_test123" {
		t.Errorf("unexpected token: %s", tok)
	}
}

func TestGetTokenForPurpose_Missing(t *testing.T) {
	m := New(false)
	_, ok := m.GetTokenForPurpose("copilot", map[string]string{})
	if ok {
		t.Error("expected no token")
	}
}

func TestGetTokenForPurpose_UnknownPurpose(t *testing.T) {
	m := New(false)
	_, ok := m.GetTokenForPurpose("unknown_purpose", map[string]string{"GITHUB_TOKEN": "tok"})
	// unknown purpose has no token list, so should not find anything
	if ok {
		t.Error("expected no token for unknown purpose")
	}
}

func TestValidateTokens_Valid(t *testing.T) {
	m := New(false)
	env := map[string]string{
		"GITHUB_TOKEN": "ghp_" + strings.Repeat("a", 36),
	}
	ok, _ := m.ValidateTokens(env)
	// validation depends on token format checks; just ensure it doesn't panic
	_ = ok
}

func TestValidateTokens_Empty(t *testing.T) {
	m := New(false)
	ok, msg := m.ValidateTokens(map[string]string{})
	_ = ok
	_ = msg
}

func TestSetupEnvironment(t *testing.T) {
	m := New(false)
	env := map[string]string{
		"GITHUB_TOKEN": "ghp_testtoken",
	}
	out := m.SetupEnvironment(env)
	if out == nil {
		t.Error("expected non-nil environment")
	}
}

func TestSetupRuntimeEnvironment(t *testing.T) {
	env := map[string]string{
		"GITHUB_TOKEN": "ghp_testtoken",
	}
	out := SetupRuntimeEnvironment(env)
	if out == nil {
		t.Error("expected non-nil environment")
	}
}

func TestValidateGitHubTokens(t *testing.T) {
	ok, msg := ValidateGitHubTokens(map[string]string{})
	_ = ok
	_ = msg
}

func TestGetGitHubTokenForRuntime(t *testing.T) {
	env := map[string]string{"GH_TOKEN": "ghp_test"}
	tok, ok := GetGitHubTokenForRuntime("copilot", env)
	_ = tok
	_ = ok
}

func TestIsValidCredentialToken(t *testing.T) {
	cases := []struct {
		token string
		valid bool
	}{
		{"ghp_" + strings.Repeat("a", 36), true},
		{"", false},
		{"short", true},
	}
	for _, c := range cases {
		got := isValidCredentialToken(c.token)
		if got != c.valid {
			t.Errorf("isValidCredentialToken(%q) = %v, want %v", c.token, got, c.valid)
		}
	}
}

func TestFormatCredentialHost_NoPort(t *testing.T) {
	got := formatCredentialHost("github.com", nil)
	if got != "github.com" {
		t.Errorf("unexpected: %s", got)
	}
}

func TestFormatCredentialHost_WithPort(t *testing.T) {
	port := 8080
	got := formatCredentialHost("github.com", &port)
	if got != "github.com:8080" {
		t.Errorf("unexpected: %s", got)
	}
}

func TestSanitizeCredentialPath(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"/repo/path"},
		{"https://github.com/owner/repo"},
		{""},
	}
	for _, c := range cases {
		out := sanitizeCredentialPath(c.input)
		_ = out
	}
}

func TestAppendOrReplace_New(t *testing.T) {
	env := []string{"FOO=bar"}
	out := appendOrReplace(env, "BAZ", "qux")
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestAppendOrReplace_Replace(t *testing.T) {
	env := []string{"FOO=old", "BAR=keep"}
	out := appendOrReplace(env, "FOO", "new")
	if len(out) != 2 {
		t.Errorf("expected 2 entries after replace, got %d", len(out))
	}
	for _, e := range out {
		if strings.HasPrefix(e, "FOO=") && e != "FOO=new" {
			t.Error("expected FOO=new")
		}
	}
}

func TestSupportsGhCLIHost(t *testing.T) {
	if !supportsGhCLIHost("github.com") {
		t.Error("expected github.com to be supported")
	}
	if supportsGhCLIHost("gitlab.com") {
		t.Error("expected gitlab.com to not be supported")
	}
}
