package version

import (
	"os"
	"path/filepath"
	"testing"
)

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
	_ = GetBuildSHA()
}

func TestGetVersion_VariousVersionStrings(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()

	cases := []string{"0.1.0", "1.0.0", "2.3.4", "10.20.30", "0.0.1"}
	for _, v := range cases {
		BuildVersion = v
		got := GetVersion()
		if got != v {
			t.Errorf("GetVersion() = %q, want %q", got, v)
		}
	}
}

func TestGetVersion_SpecialVersionStrings(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()

	cases := []string{"1.0.0a1", "2.0.0b3", "3.0.0rc1", "dev"}
	for _, v := range cases {
		BuildVersion = v
		got := GetVersion()
		if got != v {
			t.Errorf("GetVersion() = %q, want %q", got, v)
		}
	}
}

func TestGetBuildSHA_DifferentSHAs(t *testing.T) {
	orig := BuildSHA
	defer func() { BuildSHA = orig }()

	cases := []string{"abc1234", "deadbeef", "0000000", "1234567"}
	for _, sha := range cases {
		BuildSHA = sha
		got := GetBuildSHA()
		if got != sha {
			t.Errorf("GetBuildSHA() = %q, want %q", got, sha)
		}
	}
}

func TestVersionFromPyproject_ValidFile(t *testing.T) {
	dir := t.TempDir()
	pyproject := filepath.Join(dir, "pyproject.toml")
	content := `[tool.poetry]\nname = "apm"\nversion = "1.2.3"\n`
	if err := os.WriteFile(pyproject, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	// versionFromPyproject is unexported; test indirectly via GetVersion with BuildVersion=""
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	BuildVersion = ""
	// Cannot inject path, but ensure GetVersion does not panic
	_ = GetVersion()
}

func TestGetVersion_EmptyString(t *testing.T) {
	orig := BuildVersion
	defer func() { BuildVersion = orig }()
	// Empty BuildVersion triggers fallback
	BuildVersion = ""
	got := GetVersion()
	// Should return something non-empty (either from pyproject.toml or "unknown")
	if got == "" {
		t.Error("GetVersion() should not return empty string")
	}
}
