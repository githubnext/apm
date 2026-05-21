package constants

import (
	"strings"
	"testing"
)

func TestInstallModeAll_IsString(t *testing.T) {
	if string(InstallModeAll) == "" {
		t.Fatal("InstallModeAll should not be empty")
	}
}

func TestInstallModeAPM_IsString(t *testing.T) {
	if string(InstallModeAPM) == "" {
		t.Fatal("InstallModeAPM should not be empty")
	}
}

func TestInstallModeMCP_IsString(t *testing.T) {
	if string(InstallModeMCP) == "" {
		t.Fatal("InstallModeMCP should not be empty")
	}
}

func TestAllInstallModesDistinct2(t *testing.T) {
	modes := []InstallMode{InstallModeAll, InstallModeAPM, InstallModeMCP}
	seen := map[InstallMode]bool{}
	for _, m := range modes {
		if seen[m] {
			t.Fatalf("duplicate install mode: %q", m)
		}
		seen[m] = true
	}
}

func TestAPMYMLFilename_YMLExtension(t *testing.T) {
	if !strings.HasSuffix(APMYMLFilename, ".yml") {
		t.Fatalf("expected .yml extension, got %q", APMYMLFilename)
	}
}

func TestAPMLockFilename_NotEmpty(t *testing.T) {
	if APMLockFilename == "" {
		t.Fatal("APMLockFilename should not be empty")
	}
}

func TestAPMModulesDir_Value(t *testing.T) {
	if APMModulesDir == "" {
		t.Fatal("APMModulesDir should not be empty")
	}
}

func TestDefaultSkipDirs_NotEmpty(t *testing.T) {
	if len(DefaultSkipDirs) == 0 {
		t.Fatal("DefaultSkipDirs should not be empty")
	}
}

func TestDefaultSkipDirs_ContainsPycache(t *testing.T) {
	if _, ok := DefaultSkipDirs["__pycache__"]; !ok {
		t.Fatal("expected __pycache__ in DefaultSkipDirs")
	}
}

func TestDefaultSkipDirs_ContainsDotVenv2(t *testing.T) {
	if _, ok := DefaultSkipDirs[".venv"]; !ok {
		t.Fatal("expected .venv in DefaultSkipDirs")
	}
}

func TestClaudeMDFilename_HasMDExtension(t *testing.T) {
	if !strings.HasSuffix(ClaudeMDFilename, ".md") {
		t.Fatalf("expected .md extension, got %q", ClaudeMDFilename)
	}
}

func TestAgentsMDFilename_HasMDExtension(t *testing.T) {
	if !strings.HasSuffix(AgentsMDFilename, ".md") {
		t.Fatalf("expected .md extension, got %q", AgentsMDFilename)
	}
}

func TestSkillMDFilename_HasMDExtension2(t *testing.T) {
	if !strings.HasSuffix(SkillMDFilename, ".md") {
		t.Fatalf("expected .md extension, got %q", SkillMDFilename)
	}
}

func TestAPMDir_StartsWithDot2(t *testing.T) {
	if !strings.HasPrefix(APMDir, ".") {
		t.Fatalf("expected dot prefix for APMDir, got %q", APMDir)
	}
}

func TestGitignoreFilename_HasLeadingDot(t *testing.T) {
	if !strings.HasPrefix(GitignoreFilename, ".") {
		t.Fatalf("expected leading dot in gitignore filename, got %q", GitignoreFilename)
	}
}
