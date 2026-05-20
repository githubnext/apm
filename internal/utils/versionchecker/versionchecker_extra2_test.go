package versionchecker_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/versionchecker"
)

func TestIsNewerVersion_SameMajor_NewerMinor(t *testing.T) {
	if !versionchecker.IsNewerVersion("1.0.0", "1.1.0") {
		t.Error("1.1.0 should be newer than 1.0.0")
	}
}

func TestIsNewerVersion_SameMajor_NewerPatch(t *testing.T) {
	if !versionchecker.IsNewerVersion("1.0.0", "1.0.1") {
		t.Error("1.0.1 should be newer than 1.0.0")
	}
}

func TestIsNewerVersion_Identical(t *testing.T) {
	if versionchecker.IsNewerVersion("1.2.3", "1.2.3") {
		t.Error("identical versions should not be considered newer")
	}
}

func TestIsNewerVersion_OlderThanCurrent(t *testing.T) {
	if versionchecker.IsNewerVersion("2.0.0", "1.9.9") {
		t.Error("1.9.9 is not newer than 2.0.0")
	}
}

func TestIsNewerVersion_InvalidInputs(t *testing.T) {
	if versionchecker.IsNewerVersion("not-a-version", "1.0.0") {
		t.Error("invalid current should return false")
	}
	if versionchecker.IsNewerVersion("1.0.0", "not-a-version") {
		t.Error("invalid latest should return false")
	}
}

func TestParseVersion_StableRelease(t *testing.T) {
	v := versionchecker.ParseVersion("2.4.1")
	if v == nil {
		t.Fatal("expected non-nil for stable version")
	}
	if v.Major != 2 || v.Minor != 4 || v.Patch != 1 {
		t.Errorf("unexpected parsed version: %+v", v)
	}
	if v.Prerelease != "" {
		t.Errorf("expected empty prerelease, got %q", v.Prerelease)
	}
}

func TestParseVersion_Alpha(t *testing.T) {
	v := versionchecker.ParseVersion("1.0.0a1")
	if v == nil {
		t.Fatal("expected non-nil for alpha version")
	}
	if v.Prerelease != "a1" {
		t.Errorf("expected prerelease 'a1', got %q", v.Prerelease)
	}
}

func TestParseVersion_Beta(t *testing.T) {
	v := versionchecker.ParseVersion("1.0.0b2")
	if v == nil {
		t.Fatal("expected non-nil for beta version")
	}
	if v.Prerelease != "b2" {
		t.Errorf("expected prerelease 'b2', got %q", v.Prerelease)
	}
}

func TestParseVersion_InvalidInputs(t *testing.T) {
	cases := []string{"", "1.2", "1.2.3.4", "v1.2.3", "abc"}
	for _, c := range cases {
		if v := versionchecker.ParseVersion(c); v != nil {
			t.Errorf("ParseVersion(%q) expected nil, got %+v", c, v)
		}
	}
}

func TestIsNewerVersion_StableNewerThanPrerelease(t *testing.T) {
	// Stable 1.0.0 should be newer than 1.0.0a1
	if !versionchecker.IsNewerVersion("1.0.0a1", "1.0.0") {
		t.Error("stable 1.0.0 should be considered newer than 1.0.0a1")
	}
}

func TestIsNewerVersion_PrereleaseNotNewerThanStable(t *testing.T) {
	if versionchecker.IsNewerVersion("1.0.0", "1.0.0b1") {
		t.Error("prerelease should not be newer than stable of same version")
	}
}
