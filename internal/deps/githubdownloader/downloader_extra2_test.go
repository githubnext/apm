package githubdownloader

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Options struct fields
// ---------------------------------------------------------------------------

func TestOptions_Fields(t *testing.T) {
	opts := Options{
		CacheDir:    "/tmp/cache",
		Concurrency: 4,
		TimeoutSecs: 30.5,
	}
	if opts.CacheDir != "/tmp/cache" {
		t.Errorf("expected /tmp/cache, got %q", opts.CacheDir)
	}
	if opts.Concurrency != 4 {
		t.Errorf("expected 4, got %d", opts.Concurrency)
	}
	if opts.TimeoutSecs != 30.5 {
		t.Errorf("expected 30.5, got %f", opts.TimeoutSecs)
	}
}

func TestDefaultOptions_TimeoutPositive(t *testing.T) {
	opts := DefaultOptions()
	if opts.TimeoutSecs <= 0 {
		t.Errorf("expected positive TimeoutSecs, got %f", opts.TimeoutSecs)
	}
}

// ---------------------------------------------------------------------------
// TransportPlan
// ---------------------------------------------------------------------------

func TestBuildTransportPlan_HTTPSOnly(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferHTTPS, false)
	if plan.Primary != "https" {
		t.Errorf("expected primary=https, got %q", plan.Primary)
	}
	if len(plan.Fallbacks) != 0 {
		t.Errorf("expected no fallbacks when allowFallback=false, got %v", plan.Fallbacks)
	}
}

func TestBuildTransportPlan_SSHOnly(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferSSH, false)
	if plan.Primary != "ssh" {
		t.Errorf("expected primary=ssh, got %q", plan.Primary)
	}
}

func TestBuildTransportPlan_HTTPSWithFallback(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferHTTPS, true)
	if plan.Primary != "https" {
		t.Errorf("expected primary=https, got %q", plan.Primary)
	}
	if len(plan.Fallbacks) == 0 {
		t.Error("expected fallbacks when allowFallback=true")
	}
}

func TestBuildTransportPlan_SSHWithFallback(t *testing.T) {
	plan := BuildTransportPlan(ProtocolPreferSSH, true)
	if plan.Primary != "ssh" {
		t.Errorf("expected primary=ssh, got %q", plan.Primary)
	}
	if len(plan.Fallbacks) == 0 {
		t.Error("expected fallbacks when allowFallback=true for SSH")
	}
}

// ---------------------------------------------------------------------------
// ProtocolPreference constants
// ---------------------------------------------------------------------------

func TestProtocolPreference_Values(t *testing.T) {
	if ProtocolPreferHTTPS == ProtocolPreferSSH {
		t.Error("expected distinct protocol preference constants")
	}
}

// ---------------------------------------------------------------------------
// SanitizeGitError
// ---------------------------------------------------------------------------

func TestSanitizeGitError_TokenRedacted(t *testing.T) {
	msg := "fatal: https://x-access-token:ghp_secret123@github.com/org/repo.git: not found"
	got := SanitizeGitError(msg)
	if strings.Contains(got, "ghp_secret123") {
		t.Errorf("expected token to be redacted, got %q", got)
	}
}

func TestSanitizeGitError_PlainMessage(t *testing.T) {
	msg := "fatal: repository not found"
	got := SanitizeGitError(msg)
	if got == "" {
		t.Error("expected non-empty sanitized message")
	}
}

func TestSanitizeGitError_Empty(t *testing.T) {
	got := SanitizeGitError("")
	_ = got // should not panic
}

// ---------------------------------------------------------------------------
// BareCloneURL
// ---------------------------------------------------------------------------

func TestBareCloneURL_NonEmpty(t *testing.T) {
	got := BareCloneURL("/tmp/cache", "https://github.com/org/repo")
	if got == "" {
		t.Error("expected non-empty bare clone URL")
	}
}

func TestBareCloneURL_ContainsCache(t *testing.T) {
	got := BareCloneURL("/my/cache", "https://github.com/org/repo")
	if !strings.Contains(got, "/my/cache") {
		t.Errorf("expected cache dir in URL, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// DownloadResult
// ---------------------------------------------------------------------------

func TestDownloadResult_ZeroValue(t *testing.T) {
	var r DownloadResult
	if r.DestDir != "" {
		t.Error("expected empty DestDir in zero value")
	}
	if r.SHA != "" {
		t.Error("expected empty SHA in zero value")
	}
}

// ---------------------------------------------------------------------------
// RawFileResult
// ---------------------------------------------------------------------------

func TestRawFileResult_ZeroValue(t *testing.T) {
	var r RawFileResult
	if r.Content != nil {
		t.Error("expected nil Content in zero value")
	}
}
