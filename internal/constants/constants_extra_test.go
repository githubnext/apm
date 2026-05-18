package constants

import (
	"strings"
	"testing"
)

func TestInstallMode_AllDistinct(t *testing.T) {
	modes := []InstallMode{InstallModeAll, InstallModeAPM, InstallModeMCP}
	seen := map[InstallMode]bool{}
	for _, m := range modes {
		if seen[m] {
			t.Errorf("duplicate InstallMode value: %q", m)
		}
		seen[m] = true
	}
}

func TestInstallMode_NonEmpty(t *testing.T) {
	for _, m := range []InstallMode{InstallModeAll, InstallModeAPM, InstallModeMCP} {
		if string(m) == "" {
			t.Errorf("InstallMode must not be empty string")
		}
	}
}

func TestAPMYMLFilename_Extension(t *testing.T) {
	if !strings.HasSuffix(APMYMLFilename, ".yml") {
		t.Errorf("APMYMLFilename %q should have .yml extension", APMYMLFilename)
	}
}

func TestAPMLockFilename_Extension(t *testing.T) {
	if strings.HasSuffix(APMLockFilename, ".yaml") || strings.HasSuffix(APMLockFilename, ".yml") {
		t.Errorf("APMLockFilename %q should not have yaml extension (it is .lock)", APMLockFilename)
	}
}

func TestAPMDir_Hidden(t *testing.T) {
	if !strings.HasPrefix(APMDir, ".") {
		t.Errorf("APMDir %q should be a hidden directory", APMDir)
	}
}

func TestGitHubDir_Hidden(t *testing.T) {
	if !strings.HasPrefix(GitHubDir, ".") {
		t.Errorf("GitHubDir %q should be a hidden directory", GitHubDir)
	}
}

func TestClaudeDir_Hidden(t *testing.T) {
	if !strings.HasPrefix(ClaudeDir, ".") {
		t.Errorf("ClaudeDir %q should be a hidden directory", ClaudeDir)
	}
}

func TestGitignoreFilename_LeadingDot(t *testing.T) {
	if !strings.HasPrefix(GitignoreFilename, ".") {
		t.Errorf("GitignoreFilename %q should start with dot", GitignoreFilename)
	}
}

func TestAPMModulesGitignorePattern_TrailingSlash(t *testing.T) {
	if !strings.HasSuffix(APMModulesGitignorePattern, "/") {
		t.Errorf("APMModulesGitignorePattern %q should end with /", APMModulesGitignorePattern)
	}
}

func TestSkillMDFilename_IsMD(t *testing.T) {
	if !strings.HasSuffix(SkillMDFilename, ".md") {
		t.Errorf("SkillMDFilename %q should end with .md", SkillMDFilename)
	}
}

func TestAgentsMDFilename_IsMD(t *testing.T) {
	if !strings.HasSuffix(AgentsMDFilename, ".md") {
		t.Errorf("AgentsMDFilename %q should end with .md", AgentsMDFilename)
	}
}

func TestClaudeMDFilename_IsMD(t *testing.T) {
	if !strings.HasSuffix(ClaudeMDFilename, ".md") {
		t.Errorf("ClaudeMDFilename %q should end with .md", ClaudeMDFilename)
	}
}

func TestDefaultSkipDirs_HasGit(t *testing.T) {
	if _, ok := DefaultSkipDirs[".git"]; !ok {
		t.Error("DefaultSkipDirs must contain .git")
	}
}

func TestDefaultSkipDirs_HasNodeModules(t *testing.T) {
	if _, ok := DefaultSkipDirs["node_modules"]; !ok {
		t.Error("DefaultSkipDirs must contain node_modules")
	}
}

func TestDefaultSkipDirs_HasPycache(t *testing.T) {
	if _, ok := DefaultSkipDirs["__pycache__"]; !ok {
		t.Error("DefaultSkipDirs must contain __pycache__")
	}
}

func TestDefaultSkipDirs_HasAPMModules(t *testing.T) {
	if _, ok := DefaultSkipDirs["apm_modules"]; !ok {
		t.Error("DefaultSkipDirs must contain apm_modules")
	}
}

func TestDefaultSkipDirs_NoEmptyKey(t *testing.T) {
	if _, ok := DefaultSkipDirs[""]; ok {
		t.Error("DefaultSkipDirs must not have empty string key")
	}
}

func TestDefaultSkipDirs_AllNonEmpty(t *testing.T) {
	for k := range DefaultSkipDirs {
		if k == "" {
			t.Error("DefaultSkipDirs has empty key")
		}
	}
}

func TestDefaultSkipDirs_HasVenv(t *testing.T) {
	if _, ok := DefaultSkipDirs["venv"]; !ok {
		t.Error("DefaultSkipDirs must contain venv")
	}
}

func TestDefaultSkipDirs_HasDotVenv(t *testing.T) {
	if _, ok := DefaultSkipDirs[".venv"]; !ok {
		t.Error("DefaultSkipDirs must contain .venv")
	}
}

func TestDefaultSkipDirs_HasBuild(t *testing.T) {
	if _, ok := DefaultSkipDirs["build"]; !ok {
		t.Error("DefaultSkipDirs must contain build")
	}
}

func TestDefaultSkipDirs_HasDist(t *testing.T) {
	if _, ok := DefaultSkipDirs["dist"]; !ok {
		t.Error("DefaultSkipDirs must contain dist")
	}
}

func TestDefaultSkipDirs_HasMypyCache(t *testing.T) {
	if _, ok := DefaultSkipDirs[".mypy_cache"]; !ok {
		t.Error("DefaultSkipDirs must contain .mypy_cache")
	}
}

func TestDefaultSkipDirs_HasPytestCache(t *testing.T) {
	if _, ok := DefaultSkipDirs[".pytest_cache"]; !ok {
		t.Error("DefaultSkipDirs must contain .pytest_cache")
	}
}

func TestFileConstants_APMModulesDirMatchesGitignore(t *testing.T) {
	want := APMModulesDir + "/"
	if APMModulesGitignorePattern != want {
		t.Errorf("APMModulesGitignorePattern = %q, want %q", APMModulesGitignorePattern, want)
	}
}

func TestInstallModeAll_Value(t *testing.T) {
	if InstallModeAll != "all" {
		t.Errorf("InstallModeAll = %q, want all", InstallModeAll)
	}
}

func TestInstallModeAPM_Value(t *testing.T) {
	if InstallModeAPM != "apm" {
		t.Errorf("InstallModeAPM = %q, want apm", InstallModeAPM)
	}
}

func TestInstallModeMCP_Value(t *testing.T) {
	if InstallModeMCP != "mcp" {
		t.Errorf("InstallModeMCP = %q, want mcp", InstallModeMCP)
	}
}
