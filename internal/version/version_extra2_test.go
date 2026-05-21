package version

import (
	"testing"
)

func TestGetVersion_ReturnsString(t *testing.T) {
	// GetVersion should always return a non-empty string.
	v := GetVersion()
	if v == "" {
		t.Error("GetVersion should not return empty string")
	}
}

func TestGetBuildSHA_ReturnsString(t *testing.T) {
	// GetBuildSHA should not panic; may return empty string in sandbox.
	_ = GetBuildSHA()
}

func TestBuildVersion_Assignable(t *testing.T) {
	// Verify BuildVersion is a settable package-level var.
	old := BuildVersion
	BuildVersion = "1.2.3-test"
	got := GetVersion()
	if got != "1.2.3-test" {
		t.Errorf("expected BuildVersion override, got %q", got)
	}
	BuildVersion = old
}

func TestBuildSHA_Assignable(t *testing.T) {
	// Verify BuildSHA is a settable package-level var.
	old := BuildSHA
	BuildSHA = "abc1234"
	got := GetBuildSHA()
	if got != "abc1234" {
		t.Errorf("expected BuildSHA override, got %q", got)
	}
	BuildSHA = old
}

func TestGetVersion_FallbackNotEmpty_WhenBuildVersionEmpty(t *testing.T) {
	old := BuildVersion
	BuildVersion = ""
	v := GetVersion()
	// Should return pyproject.toml version or "unknown", never empty.
	if v == "" {
		t.Error("GetVersion should never return empty string")
	}
	BuildVersion = old
}
