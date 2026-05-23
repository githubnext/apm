// Package constants_test provides parity tests for constants.
package constants_test

import (
	"testing"

	"github.com/githubnext/apm/internal/constants"
)

// TestParityConstantsFilenames verifies file and directory name constants
// match the Python source (src/apm_cli/constants.py).
func TestParityConstantsFilenames(t *testing.T) {
	cases := []struct {
		name     string
		got      string
		expected string
	}{
		{"APMYMLFilename", constants.APMYMLFilename, "apm.yml"},
		{"APMLockFilename", constants.APMLockFilename, "apm.lock"},
		{"APMModulesDir", constants.APMModulesDir, "apm_modules"},
		{"APMDir", constants.APMDir, ".apm"},
		{"SkillMDFilename", constants.SkillMDFilename, "SKILL.md"},
		{"AgentsMDFilename", constants.AgentsMDFilename, "AGENTS.md"},
		{"ClaudeMDFilename", constants.ClaudeMDFilename, "CLAUDE.md"},
		{"GithubDir", constants.GithubDir, ".github"},
		{"ClaudeDir", constants.ClaudeDir, ".claude"},
		{"GitignoreFilename", constants.GitignoreFilename, ".gitignore"},
		{"APMModulesGitignorePattern", constants.APMModulesGitignorePattern, "apm_modules/"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.expected {
				t.Errorf("got %q, want %q", tc.got, tc.expected)
			}
		})
	}
}

// TestParityConstantsInstallMode verifies InstallMode values match Python.
func TestParityConstantsInstallMode(t *testing.T) {
	if constants.InstallModeAll != "all" {
		t.Errorf("InstallModeAll: got %q, want %q", constants.InstallModeAll, "all")
	}
	if constants.InstallModeAPM != "apm" {
		t.Errorf("InstallModeAPM: got %q, want %q", constants.InstallModeAPM, "apm")
	}
	if constants.InstallModeMCP != "mcp" {
		t.Errorf("InstallModeMCP: got %q, want %q", constants.InstallModeMCP, "mcp")
	}
}

// TestParityConstantsDefaultSkipDirs verifies the default skip dirs set
// matches Python's DEFAULT_SKIP_DIRS frozenset.
func TestParityConstantsDefaultSkipDirs(t *testing.T) {
	expected := []string{
		".git", "node_modules", "__pycache__", ".pytest_cache",
		".venv", "venv", ".tox", "build", "dist", ".mypy_cache", "apm_modules",
	}
	for _, dir := range expected {
		if !constants.DefaultSkipDirs[dir] {
			t.Errorf("DefaultSkipDirs missing %q", dir)
		}
	}
	// .apm must NOT be in skip dirs
	if constants.DefaultSkipDirs[".apm"] {
		t.Error("DefaultSkipDirs must not contain .apm")
	}
}
