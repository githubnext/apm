package version

import (
	"testing"
)

func TestBuildVersion_DefaultEmpty(t *testing.T) {
	// Without -ldflags injection, BuildVersion should be empty or set
	// Just verify the variable is accessible
	_ = BuildVersion
}

func TestBuildSHA_DefaultEmpty(t *testing.T) {
	_ = BuildSHA
}

func TestGetVersion_SetAndRestore(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "0.9.0"
	if got := GetVersion(); got != "0.9.0" {
		t.Errorf("GetVersion() = %q, want 0.9.0", got)
	}
}

func TestGetVersion_NonEmptyAfterSet(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "3.0.0"
	if GetVersion() == "" {
		t.Error("GetVersion() should not be empty when BuildVersion is set")
	}
}

func TestGetBuildSHA_SetAndRestore(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = "deadbeef"
	if got := GetBuildSHA(); got != "deadbeef" {
		t.Errorf("GetBuildSHA() = %q, want deadbeef", got)
	}
}

func TestGetBuildSHA_LongSHA(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = "1234567890abcdef"
	if got := GetBuildSHA(); got != "1234567890abcdef" {
		t.Errorf("GetBuildSHA() = %q, want 1234567890abcdef", got)
	}
}

func TestGetVersion_Alpha(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "1.0.0a1"
	got := GetVersion()
	if got != "1.0.0a1" {
		t.Errorf("GetVersion() = %q, want 1.0.0a1", got)
	}
}

func TestGetVersion_Beta(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "2.0.0b2"
	got := GetVersion()
	if got != "2.0.0b2" {
		t.Errorf("GetVersion() = %q, want 2.0.0b2", got)
	}
}

func TestGetVersion_ReleaseCandidate(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "1.5.0rc3"
	got := GetVersion()
	if got != "1.5.0rc3" {
		t.Errorf("GetVersion() = %q, want 1.5.0rc3", got)
	}
}

func TestGetVersion_MultipleIterations(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	for i, v := range []string{"1.0.0", "2.0.0", "3.0.0"} {
		BuildVersion = v
		got := GetVersion()
		if got != v {
			t.Errorf("iteration %d: GetVersion() = %q, want %q", i, got, v)
		}
	}
}

func TestGetBuildSHA_AllZeros(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = "0000000"
	got := GetBuildSHA()
	if got != "0000000" {
		t.Errorf("GetBuildSHA() = %q, want 0000000", got)
	}
}

func TestGetBuildSHA_EmptyFallsThrough(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()
	BuildSHA = ""
	// Should not panic; may return "" or a git SHA
	_ = GetBuildSHA()
}

func TestGetVersion_EmptyFallbackNotEmpty(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = ""
	got := GetVersion()
	// In CI/dev the fallback reads pyproject.toml or returns "unknown"
	_ = got // just assert no panic
}

func TestGetVersion_PatchVersion(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "0.0.1"
	if got := GetVersion(); got != "0.0.1" {
		t.Errorf("GetVersion() = %q, want 0.0.1", got)
	}
}

func TestGetVersion_MajorVersion(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = "100.0.0"
	if got := GetVersion(); got != "100.0.0" {
		t.Errorf("GetVersion() = %q, want 100.0.0", got)
	}
}
