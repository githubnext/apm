package commandlogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/commandlogger"
)

func TestNewCommandLogger_Defaults(t *testing.T) {
	l := commandlogger.NewCommandLogger("audit", false, false)
	if l.Command != "audit" {
		t.Errorf("expected Command='audit', got %q", l.Command)
	}
	if l.Verbose {
		t.Error("expected Verbose=false")
	}
	if l.DryRun {
		t.Error("expected DryRun=false")
	}
}

func TestNewCommandLogger_DryRunVerbose(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", true, true)
	if !l.Verbose {
		t.Error("expected Verbose=true")
	}
	if !l.DryRun {
		t.Error("expected DryRun=true")
	}
	if l.ShouldExecute() {
		t.Error("expected ShouldExecute()=false for dry-run")
	}
}

func TestStripSourcePrefix_OrgWithPath(t *testing.T) {
	got := commandlogger.StripSourcePrefix("org:mycompany/subgroup")
	if got != "mycompany/subgroup" {
		t.Errorf("got %q, want %q", got, "mycompany/subgroup")
	}
}

func TestStripSourcePrefix_ShortOrg(t *testing.T) {
	// "org:" with empty suffix should be returned unchanged
	got := commandlogger.StripSourcePrefix("org:")
	if got != "org:" {
		t.Errorf("got %q, want %q", got, "org:")
	}
}

func TestCommandLogger_PolicyDiscoveryMiss_Absent(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("absent", "org:myorg", "", "myorg")
}

func TestCommandLogger_PolicyDiscoveryMiss_Empty(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("empty", "org:myorg", "", "")
}

func TestCommandLogger_PolicyDiscoveryMiss_Malformed(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("malformed", "org:myorg", "unexpected key", "")
}

func TestCommandLogger_PolicyDiscoveryMiss_CacheMissFetchFail(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("cache_miss_fetch_fail", "org:myorg", "connection refused", "")
}

func TestCommandLogger_PolicyDiscoveryMiss_HashMismatch(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("hash_mismatch", "org:myorg", "abc123 != def456", "")
}

func TestCommandLogger_PolicyDiscoveryMiss_Default(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDiscoveryMiss("some_unknown_outcome", "", "something went wrong", "")
}

func TestCommandLogger_PolicyViolation_Block(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyViolation("org/bad-pkg#v1.0.0", "disallowed package", "block", "org:myorg")
}

func TestCommandLogger_PolicyViolation_NoSource(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyViolation("org/pkg#v1", "reason", "block", "")
}

func TestCommandLogger_PolicyDisabled(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PolicyDisabled("--no-policy flag")
}

func TestCommandLogger_AuthStep_Success(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", true, false)
	l.AuthStep("resolve token", true, "github.com")
}

func TestCommandLogger_AuthStep_Failure(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", true, false)
	l.AuthStep("resolve token", false, "")
}

func TestCommandLogger_AuthStep_NotVerbose(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	// Should be a no-op when not verbose
	l.AuthStep("resolve token", true, "detail")
}

func TestCommandLogger_PackageInlineWarning_Verbose(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", true, false)
	l.PackageInlineWarning("package warning message")
}

func TestCommandLogger_PackageInlineWarning_NotVerbose(t *testing.T) {
	l := commandlogger.NewCommandLogger("install", false, false)
	l.PackageInlineWarning("package warning message")
}
