package version

import "testing"

func TestGetVersion_BuildVersion(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()

	BuildVersion = "1.2.3"
	if got := GetVersion(); got != "1.2.3" {
		t.Errorf("GetVersion() = %q, want %q", got, "1.2.3")
	}
}

func TestGetVersion_Fallback(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()

	BuildVersion = ""
	got := GetVersion()
	// In test mode, either parses from pyproject.toml or returns "unknown"
	if got == "" {
		t.Error("GetVersion() should not be empty")
	}
}

func TestGetBuildSHA_BuildSHA(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()

	BuildSHA = "abc1234"
	if got := GetBuildSHA(); got != "abc1234" {
		t.Errorf("GetBuildSHA() = %q, want %q", got, "abc1234")
	}
}

func TestGetBuildSHA_Fallback(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()

	BuildSHA = ""
	// In a git repo this should return something or empty; just should not panic
	_ = GetBuildSHA()
}
