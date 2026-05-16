package constants

import "testing"

func TestInstallModeValues(t *testing.T) {
	if InstallModeAll != "all" {
		t.Errorf("InstallModeAll = %q, want %q", InstallModeAll, "all")
	}
	if InstallModeAPM != "apm" {
		t.Errorf("InstallModeAPM = %q, want %q", InstallModeAPM, "apm")
	}
	if InstallModeMCP != "mcp" {
		t.Errorf("InstallModeMCP = %q, want %q", InstallModeMCP, "mcp")
	}
}

func TestFileConstants(t *testing.T) {
	cases := map[string]string{
		"APMYMLFilename":   APMYMLFilename,
		"APMLockFilename":  APMLockFilename,
		"APMModulesDir":    APMModulesDir,
		"APMDir":           APMDir,
		"SkillMDFilename":  SkillMDFilename,
		"AgentsMDFilename": AgentsMDFilename,
		"ClaudeMDFilename": ClaudeMDFilename,
		"GitHubDir":        GitHubDir,
		"ClaudeDir":        ClaudeDir,
	}
	for name, val := range cases {
		if val == "" {
			t.Errorf("constant %s is empty", name)
		}
	}
	if APMYMLFilename != "apm.yml" {
		t.Errorf("APMYMLFilename = %q, want %q", APMYMLFilename, "apm.yml")
	}
	if APMLockFilename != "apm.lock" {
		t.Errorf("APMLockFilename = %q, want %q", APMLockFilename, "apm.lock")
	}
}

func TestDefaultSkipDirs(t *testing.T) {
	mustSkip := []string{".git", "node_modules", "__pycache__", ".venv", "apm_modules"}
	for _, d := range mustSkip {
		if _, ok := DefaultSkipDirs[d]; !ok {
			t.Errorf("DefaultSkipDirs missing %q", d)
		}
	}
}
