package githubdownloader

import (
	"strings"
	"testing"
)

func TestParseLsRemoteOutput_basic(t *testing.T) {
	input := "abc123\trefs/heads/main\ndef456\trefs/tags/v1.0.0\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
	if refs[0].SHA != "abc123" || refs[0].Name != "refs/heads/main" {
		t.Errorf("unexpected ref[0]: %+v", refs[0])
	}
	if refs[1].SHA != "def456" || refs[1].Name != "refs/tags/v1.0.0" {
		t.Errorf("unexpected ref[1]: %+v", refs[1])
	}
}

func TestParseLsRemoteOutput_empty(t *testing.T) {
	refs := ParseLsRemoteOutput("")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for empty input, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_skips_malformed(t *testing.T) {
	input := "abc123\trefs/heads/main\nmalformed_line\ndef456\trefs/tags/v2.0.0\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d", len(refs))
	}
}

func TestSemverSortKey_valid(t *testing.T) {
	tests := []struct {
		name     string
		expected [4]int
	}{
		{"v1.2.3", [4]int{1, 2, 3, 0}},
		{"2.10.5", [4]int{2, 10, 5, 0}},
		{"v0.0.1", [4]int{0, 0, 1, 0}},
		{"v1.0.0-alpha", [4]int{1, 0, 0, -1}},
	}
	for _, tc := range tests {
		got := SemverSortKey(tc.name)
		if got != tc.expected {
			t.Errorf("SemverSortKey(%q) = %v, want %v", tc.name, got, tc.expected)
		}
	}
}

func TestSemverSortKey_invalid(t *testing.T) {
	key := SemverSortKey("not-a-version")
	if key[0] != -1 {
		t.Errorf("expected -1 for non-semver, got %v", key)
	}
}

func TestSortRemoteRefs_ordering(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.0", SHA: "a"},
		{Name: "v2.0.0", SHA: "b"},
		{Name: "v1.5.0", SHA: "c"},
	}
	sorted := SortRemoteRefs(refs)
	if sorted[0].Name != "v2.0.0" {
		t.Errorf("expected v2.0.0 first, got %s", sorted[0].Name)
	}
	if sorted[1].Name != "v1.5.0" {
		t.Errorf("expected v1.5.0 second, got %s", sorted[1].Name)
	}
}

func TestBareCloneURL_sanitizes(t *testing.T) {
	url := BareCloneURL("/cache", "https://github.com/owner/repo")
	if !strings.HasPrefix(url, "/cache/") {
		t.Errorf("expected path under /cache, got %s", url)
	}
	if !strings.HasSuffix(url, ".git") {
		t.Errorf("expected .git suffix, got %s", url)
	}
	// should not contain ://
	if strings.Contains(url, "://") {
		t.Errorf("URL should be sanitized, got %s", url)
	}
}

func TestSanitizeGitError_redacts_token(t *testing.T) {
	msg := "error: https://x-access-token:ghp_SECRETTOKEN@github.com/owner/repo"
	sanitized := SanitizeGitError(msg)
	if strings.Contains(sanitized, "SECRETTOKEN") {
		t.Errorf("token should be redacted, got: %s", sanitized)
	}
	if !strings.Contains(sanitized, "[REDACTED]") {
		t.Errorf("expected [REDACTED] in output, got: %s", sanitized)
	}
}

func TestSanitizeGitError_no_token(t *testing.T) {
	msg := "fatal: repository not found"
	sanitized := SanitizeGitError(msg)
	if sanitized != msg {
		t.Errorf("message without token should be unchanged, got: %s", sanitized)
	}
}

func TestBuildTransportPlan_https_only(t *testing.T) {
	plan := BuildTransportPlan(ProtocolHTTPSOnly, true)
	if plan.Primary != "https" {
		t.Errorf("expected https primary, got %s", plan.Primary)
	}
	if len(plan.Fallbacks) != 0 {
		t.Errorf("HTTPS-only should have no fallbacks, got %v", plan.Fallbacks)
	}
}

func TestBuildTransportPlan_ssh_only(t *testing.T) {
	plan := BuildTransportPlan(ProtocolSSHOnly, true)
	if plan.Primary != "ssh" {
		t.Errorf("expected ssh primary, got %s", plan.Primary)
	}
}

func TestBuildTransportPlan_prefer_https_with_fallback(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferHTTPS, true)
	if plan.Primary != "https" {
		t.Errorf("expected https primary")
	}
	if len(plan.Fallbacks) == 0 || plan.Fallbacks[0] != "ssh" {
		t.Errorf("expected ssh fallback, got %v", plan.Fallbacks)
	}
}

func TestBuildTransportPlan_prefer_ssh_with_fallback(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferSSH, true)
	if plan.Primary != "ssh" {
		t.Errorf("expected ssh primary")
	}
	if len(plan.Fallbacks) == 0 || plan.Fallbacks[0] != "https" {
		t.Errorf("expected https fallback, got %v", plan.Fallbacks)
	}
}

func TestBuildTransportPlan_no_fallback(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferHTTPS, false)
	if len(plan.Fallbacks) != 0 {
		t.Errorf("no fallback expected, got %v", plan.Fallbacks)
	}
}
