package versionchecker_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/versionchecker"
)

func TestParseVersion_RCPrerelease(t *testing.T) {
	v := versionchecker.ParseVersion("2.0.0rc1")
	if v == nil {
		t.Fatal("expected non-nil for rc prerelease")
	}
	if v.Major != 2 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("unexpected components: %+v", v)
	}
	if v.Prerelease != "rc1" {
		t.Errorf("Prerelease = %q, want rc1", v.Prerelease)
	}
}

func TestParseVersion_AlphaPrerelease(t *testing.T) {
	v := versionchecker.ParseVersion("1.3.0a2")
	if v == nil {
		t.Fatal("expected non-nil for alpha prerelease")
	}
	if v.Prerelease != "a2" {
		t.Errorf("Prerelease = %q, want a2", v.Prerelease)
	}
}

func TestIsNewerVersion_BothPrerelease_BetaVsRC(t *testing.T) {
	// rc > b lexicographically
	if !versionchecker.IsNewerVersion("1.0.0b1", "1.0.0rc1") {
		t.Error("rc1 should be newer than b1")
	}
}

func TestIsNewerVersion_PreReleaseVsStable(t *testing.T) {
	// stable 1.0.0 is newer than 1.0.0b1
	if !versionchecker.IsNewerVersion("1.0.0b1", "1.0.0") {
		t.Error("stable should be newer than prerelease of same version")
	}
}

func TestIsNewerVersion_SamePatchDifferentPrerelease(t *testing.T) {
	// b2 > b1
	if !versionchecker.IsNewerVersion("1.0.0b1", "1.0.0b2") {
		t.Error("b2 should be newer than b1")
	}
}

func TestIsNewerVersion_MajorTakesPrecedence(t *testing.T) {
	if !versionchecker.IsNewerVersion("1.99.99", "2.0.0") {
		t.Error("major bump should dominate minor/patch")
	}
	if versionchecker.IsNewerVersion("2.0.0", "1.99.99") {
		t.Error("older major should not be newer")
	}
}

func TestIsNewerVersion_MinorTakesPrecedence(t *testing.T) {
	if !versionchecker.IsNewerVersion("1.0.99", "1.1.0") {
		t.Error("minor bump should dominate patch")
	}
}

func TestIsNewerVersion_BothInvalid(t *testing.T) {
	if versionchecker.IsNewerVersion("garbage", "also-garbage") {
		t.Error("both invalid versions should not be considered newer")
	}
}

func TestVersionComponents_Fields(t *testing.T) {
	v := versionchecker.ParseVersion("3.14.159")
	if v == nil {
		t.Fatal("expected non-nil")
	}
	if v.Major != 3 {
		t.Errorf("Major = %d, want 3", v.Major)
	}
	if v.Minor != 14 {
		t.Errorf("Minor = %d, want 14", v.Minor)
	}
	if v.Patch != 159 {
		t.Errorf("Patch = %d, want 159", v.Patch)
	}
	if v.Prerelease != "" {
		t.Errorf("Prerelease = %q, want empty", v.Prerelease)
	}
}

func TestParseVersion_LeadingVPrefix(t *testing.T) {
	// v-prefix should not be accepted by the regex
	v := versionchecker.ParseVersion("v1.2.3")
	if v != nil {
		t.Error("v-prefixed version should not parse (no v prefix support)")
	}
}

func TestIsNewerVersion_ZeroVersions(t *testing.T) {
	if versionchecker.IsNewerVersion("0.0.0", "0.0.0") {
		t.Error("equal zero versions should not be newer")
	}
	if !versionchecker.IsNewerVersion("0.0.0", "0.0.1") {
		t.Error("0.0.1 should be newer than 0.0.0")
	}
}
