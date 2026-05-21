package tokenmanager

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

func TestADOBearerSource(t *testing.T) {
	if ADOBearerSource == "" {
		t.Error("ADOBearerSource should not be empty")
	}
}

func TestDefaultCredentialTimeout_Positive(t *testing.T) {
	if DefaultCredentialTimeout <= 0 {
		t.Errorf("expected positive DefaultCredentialTimeout, got %d", DefaultCredentialTimeout)
	}
}

func TestMaxCredentialTimeout_GreaterThanDefault(t *testing.T) {
	if MaxCredentialTimeout <= DefaultCredentialTimeout {
		t.Errorf("expected MaxCredentialTimeout > DefaultCredentialTimeout, got %d <= %d", MaxCredentialTimeout, DefaultCredentialTimeout)
	}
}

// ---------------------------------------------------------------------------
// formatCredentialHost
// ---------------------------------------------------------------------------

func TestFormatCredentialHostNilPort2(t *testing.T) {
	got := formatCredentialHost("github.com", nil)
	if got == "" {
		t.Error("expected non-empty credential host")
	}
	if strings.Contains(got, ":") && !strings.Contains(got, "github.com") {
		t.Errorf("unexpected format: %q", got)
	}
}

func TestFormatCredentialHostWithPort8080(t *testing.T) {
	port := 8080
	got := formatCredentialHost("github.company.com", &port)
	if !strings.Contains(got, "8080") {
		t.Errorf("expected port 8080 in result, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// sanitizeCredentialPath
// ---------------------------------------------------------------------------

func TestSanitizeCredentialPath_Empty(t *testing.T) {
	got := sanitizeCredentialPath("")
	_ = got // should not panic
}

func TestSanitizeCredentialPath_WithSlash(t *testing.T) {
	got := sanitizeCredentialPath("/org/repo")
	_ = got
}

// ---------------------------------------------------------------------------
// isValidCredentialToken
// ---------------------------------------------------------------------------

func TestIsValidCredentialToken_ValidToken(t *testing.T) {
	if !isValidCredentialToken("ghp_abc123DEF456") {
		t.Error("expected valid token")
	}
}

func TestIsValidCredentialToken_Empty(t *testing.T) {
	if isValidCredentialToken("") {
		t.Error("empty token should not be valid")
	}
}

func TestIsValidCredentialToken_Whitespace(t *testing.T) {
	if isValidCredentialToken("   ") {
		t.Error("whitespace-only token should not be valid")
	}
}

// ---------------------------------------------------------------------------
// supportsGhCLIHost
// ---------------------------------------------------------------------------

func TestSupportsGhCLIHost_GitHub(t *testing.T) {
	if !supportsGhCLIHost("github.com") {
		t.Error("expected github.com to be supported")
	}
}

func TestSupportsGhCLIHost_Unknown(t *testing.T) {
	// Unknown host may or may not be supported; just ensure no panic
	_ = supportsGhCLIHost("internal.company.com")
}

// ---------------------------------------------------------------------------
// GetGitHubTokenForRuntime
// ---------------------------------------------------------------------------

func TestGetGitHubTokenForRuntime_NoToken(t *testing.T) {
	_, ok := GetGitHubTokenForRuntime("copilot", map[string]string{})
	// no token in empty env — result may vary; just verify no panic
	_ = ok
}

func TestGetGitHubTokenForRuntime_UnknownRuntimeExtra2(t *testing.T) {
	_, ok := GetGitHubTokenForRuntime("unknown-runtime", map[string]string{})
	if ok {
		t.Error("expected false for unknown runtime with no token")
	}
}

// ---------------------------------------------------------------------------
// appendOrReplace
// ---------------------------------------------------------------------------

func TestAppendOrReplace_Append(t *testing.T) {
	env := []string{"A=1", "B=2"}
	got := appendOrReplace(env, "C", "3")
	found := false
	for _, e := range got {
		if e == "C=3" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected C=3 in result, got %v", got)
	}
}

func TestAppendOrReplace_ReplaceExtra2(t *testing.T) {
	env := []string{"A=1", "B=2"}
	got := appendOrReplace(env, "A", "99")
	count := 0
	for _, e := range got {
		if strings.HasPrefix(e, "A=") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one A= entry, got %d in %v", count, got)
	}
}
