package version

import (
	"strings"
	"testing"
)

func TestGetVersion_ReturnsNonEmpty(t *testing.T) {
	got := GetVersion()
	if strings.TrimSpace(got) == "" {
		t.Error("GetVersion should not return empty string")
	}
}

func TestGetVersion_BuildVersionOverrides(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "3.1.4"
	if got := GetVersion(); got != "3.1.4" {
		t.Errorf("expected 3.1.4, got %q", got)
	}
}

func TestGetVersion_BuildVersionSemver(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "0.0.1"
	if got := GetVersion(); got != "0.0.1" {
		t.Errorf("expected 0.0.1, got %q", got)
	}
}

func TestGetVersion_BuildVersionPreRelease(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "1.0.0a1"
	if got := GetVersion(); got != "1.0.0a1" {
		t.Errorf("expected 1.0.0a1, got %q", got)
	}
}

func TestGetBuildSHA_ReturnsBuildSHAWhenSet(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = "deadbeef"
	if got := GetBuildSHA(); got != "deadbeef" {
		t.Errorf("expected deadbeef, got %q", got)
	}
}

func TestGetBuildSHA_EmptyBuildSHAFallsThrough(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = ""
	// Just ensure we don't panic; result is environment-dependent
	_ = GetBuildSHA()
}

func TestGetBuildSHA_NonEmpty_IfBuildSHASet(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = "abc12345"
	got := GetBuildSHA()
	if got == "" {
		t.Error("expected non-empty when BuildSHA is set")
	}
}

func TestBuildVersion_CanBeSet(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "99.0.0"
	if BuildVersion != "99.0.0" {
		t.Error("BuildVersion should be assignable")
	}
}

func TestBuildSHA_CanBeSet(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = "cafebabe"
	if BuildSHA != "cafebabe" {
		t.Error("BuildSHA should be assignable")
	}
}
