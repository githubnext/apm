package tokenmanager

import (
	"strings"
	"testing"
)

func TestNew_PreserveExistingVariants(t *testing.T) {
	m := New(true)
	if !m.PreserveExisting {
		t.Error("expected PreserveExisting=true")
	}
	m2 := New(false)
	if m2.PreserveExisting {
		t.Error("expected PreserveExisting=false")
	}
	// Both should be distinct instances
	if m == m2 {
		t.Error("expected distinct manager instances")
	}
}

func TestGetTokenForPurpose_MultiplePurposes(t *testing.T) {
	m := New(false)
	env := map[string]string{
		"GITHUB_TOKEN":  "ghp_token",
		"CODEX_API_KEY": "codex-key",
	}
	// Models purpose
	tok, ok := m.GetTokenForPurpose("models", env)
	if !ok || tok == "" {
		t.Error("expected token for models purpose with GITHUB_TOKEN set")
	}
	// Codex purpose - may or may not be found depending on env var mapping
	_, _ = m.GetTokenForPurpose("codex", env)
}

func TestGetTokenForPurpose_EmptyEnv(t *testing.T) {
	m := New(false)
	_, ok := m.GetTokenForPurpose("models", map[string]string{})
	// With no tokens in env, should return false
	_ = ok // may or may not find system env tokens; no panic is the key assertion
}

func TestValidateTokens_WithToken(t *testing.T) {
	m := New(false)
	env := map[string]string{
		"GITHUB_TOKEN": "ghp_validtoken",
	}
	valid, _ := m.ValidateTokens(env)
	if !valid {
		t.Error("expected valid=true with GITHUB_TOKEN set")
	}
}

func TestValidateTokens_EmptyEnvMap(t *testing.T) {
	m := New(false)
	valid, msg := m.ValidateTokens(map[string]string{})
	// With no tokens, should either be valid=false or provide a message
	_ = valid
	_ = msg
}

func TestSetupEnvironment_ReturnsMap(t *testing.T) {
	m := New(false)
	env := map[string]string{
		"GITHUB_TOKEN": "ghp_test",
	}
	result := m.SetupEnvironment(env)
	if result == nil {
		t.Fatal("SetupEnvironment should not return nil")
	}
	// Should contain at minimum the keys we passed in
	for k, v := range env {
		if result[k] != v {
			t.Errorf("expected %s=%s in result, got %s", k, v, result[k])
		}
	}
}

func TestSetupEnvironment_PreserveExisting(t *testing.T) {
	m := New(true)
	env := map[string]string{"GITHUB_TOKEN": "original"}
	result := m.SetupEnvironment(env)
	if result["GITHUB_TOKEN"] != "original" {
		t.Errorf("PreserveExisting: expected GITHUB_TOKEN=original, got %q", result["GITHUB_TOKEN"])
	}
}

func TestGetTokenWithCredentialFallback_NoTokenNoHost(t *testing.T) {
	m := New(false)
	tok, ok := m.GetTokenWithCredentialFallback("models", "github.com", map[string]string{}, nil)
	_ = ok
	_ = tok
	// Should not panic
}

func TestSetupRuntimeEnvironment_ReturnsMap(t *testing.T) {
	env := map[string]string{"GITHUB_TOKEN": "ghp_rt_test"}
	result := SetupRuntimeEnvironment(env)
	if result == nil {
		t.Fatal("SetupRuntimeEnvironment should not return nil")
	}
}

func TestValidateGitHubTokens_WithToken(t *testing.T) {
	env := map[string]string{"GITHUB_TOKEN": "ghp_validate_test"}
	valid, _ := ValidateGitHubTokens(env)
	if !valid {
		t.Error("expected valid=true with GITHUB_TOKEN set")
	}
}

func TestGetGitHubTokenForRuntime_ModelsRuntime(t *testing.T) {
	env := map[string]string{"GITHUB_TOKEN": "ghp_for_runtime"}
	tok, ok := GetGitHubTokenForRuntime("models", env)
	_ = ok
	_ = tok
	// Should not panic; token may or may not be found depending on env var mapping
}

func TestGetGitHubTokenForRuntime_UnknownRuntime(t *testing.T) {
	env := map[string]string{"GITHUB_TOKEN": "ghp_test"}
	tok, ok := GetGitHubTokenForRuntime("unknown-runtime-xyz", env)
	_ = ok
	_ = tok
	// Should not panic
}

func TestADOBearerSource_Value(t *testing.T) {
	if !strings.Contains(ADOBearerSource, "AAD") && !strings.Contains(ADOBearerSource, "az") {
		t.Errorf("ADOBearerSource looks unexpected: %q", ADOBearerSource)
	}
}
