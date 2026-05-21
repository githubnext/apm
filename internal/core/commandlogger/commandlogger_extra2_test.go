package commandlogger

import (
	"testing"
)

func TestStripSourcePrefix_URLPrefix(t *testing.T) {
	result := StripSourcePrefix("url:https://example.com/policy.yml")
	if result != "https://example.com/policy.yml" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestStripSourcePrefix_OrgPrefix(t *testing.T) {
	result := StripSourcePrefix("org:myorg")
	if result != "myorg" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestStripSourcePrefix_Empty(t *testing.T) {
	result := StripSourcePrefix("")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestStripSourcePrefix_NoPrefix(t *testing.T) {
	result := StripSourcePrefix("plainvalue")
	if result != "plainvalue" {
		t.Errorf("expected unchanged, got %q", result)
	}
}

func TestCommandLogger_Fields(t *testing.T) {
	l := NewCommandLogger("install", true, false)
	if l.Command != "install" {
		t.Errorf("expected Command=install, got %q", l.Command)
	}
	if !l.Verbose {
		t.Error("expected Verbose=true")
	}
	if l.DryRun {
		t.Error("expected DryRun=false")
	}
}

func TestCommandLogger_ShouldExecute_DryRun(t *testing.T) {
	l := NewCommandLogger("install", false, true)
	if l.ShouldExecute() {
		t.Error("ShouldExecute should be false in dry-run mode")
	}
}

func TestCommandLogger_ShouldExecute_Live(t *testing.T) {
	l := NewCommandLogger("install", false, false)
	if !l.ShouldExecute() {
		t.Error("ShouldExecute should be true when not dry-run")
	}
}

func TestInstallLogger_Fields(t *testing.T) {
	l := NewInstallLogger(false, false, false)
	if l.CommandLogger == nil {
		t.Error("expected CommandLogger to be non-nil")
	}
	if l.Partial {
		t.Error("expected Partial=false")
	}
}

func TestInstallLogger_Partial(t *testing.T) {
	l := NewInstallLogger(false, false, true)
	if !l.Partial {
		t.Error("expected Partial=true")
	}
}

func TestCommandLogger_MCPLookupHeartbeat_NoPanic(t *testing.T) {
	l := NewCommandLogger("install", false, false)
	l.MCPLookupHeartbeat(0)
	l.MCPLookupHeartbeat(1)
	l.MCPLookupHeartbeat(10)
}

func TestCommandLogger_PolicyDiscoveryMiss_AbsentVerbose(t *testing.T) {
	l := NewCommandLogger("install", true, false)
	// Should not panic
	l.PolicyDiscoveryMiss("absent", "org:myorg", "", "myorg")
}

func TestCommandLogger_PolicyDiscoveryMiss_CachedStale(t *testing.T) {
	l := NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("cached_stale", "", "refresh failed", "")
}

func TestCommandLogger_PolicyDiscoveryMiss_GarbageResponse(t *testing.T) {
	l := NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("garbage_response", "url:http://x", "not yaml", "")
}

func TestCommandLogger_PolicyDiscoveryMiss_NoGitRemote(t *testing.T) {
	l := NewCommandLogger("install", true, false)
	l.PolicyDiscoveryMiss("no_git_remote", "", "", "")
}

func TestCommandLogger_InstallSummary_NoPanic(t *testing.T) {
	l := NewCommandLogger("install", false, false)
	l.InstallSummary(1, 0, 0, 0, 1.5, true)
	l.InstallSummary(0, 1, 0, 0, 0.5, false)
	l.InstallSummary(2, 3, 1, 0, 2.0, true)
	l.InstallSummary(0, 0, 2, 0, 0.0, false)
}

func TestCommandLogger_InstallInterrupted_NoPanic(t *testing.T) {
	l := NewCommandLogger("install", false, false)
	l.InstallInterrupted(3.7)
}

func TestInstallLogger_ValidationStart_NoPanic(t *testing.T) {
	l := NewInstallLogger(false, false, false)
	l.ValidationStart(0)
	l.ValidationStart(1)
	l.ValidationStart(5)
}

func TestInstallLogger_ResolutionStart_NoPanic(t *testing.T) {
	l := NewInstallLogger(false, false, false)
	l.ResolutionStart(3, 5)
}

func TestInstallLogger_ResolutionStart_Partial(t *testing.T) {
	l := NewInstallLogger(true, false, true)
	l.ResolutionStart(1, 0)
	l.ResolutionStart(2, 4)
}

func TestInstallLogger_NothingToInstall_NoPanic(t *testing.T) {
	l := NewInstallLogger(false, false, false)
	l.NothingToInstall(true, false)
	l.NothingToInstall(false, true)
}
