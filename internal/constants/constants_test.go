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

func TestDefaultSkipDirs_extraEntries(t *testing.T) {
	extras := []string{"venv", ".tox", "build", "dist", ".mypy_cache", ".pytest_cache"}
	for _, d := range extras {
		if _, ok := DefaultSkipDirs[d]; !ok {
			t.Errorf("DefaultSkipDirs missing %q", d)
		}
	}
}

func TestInstallMode_stringConversion(t *testing.T) {
	cases := []struct {
		mode InstallMode
		want string
	}{
		{InstallModeAll, "all"},
		{InstallModeAPM, "apm"},
		{InstallModeMCP, "mcp"},
	}
	for _, c := range cases {
		if string(c.mode) != c.want {
			t.Errorf("InstallMode %q: string() = %q, want %q", c.mode, string(c.mode), c.want)
		}
	}
}

func TestFileConstants_gitignore(t *testing.T) {
	if GitignoreFilename != ".gitignore" {
		t.Errorf("GitignoreFilename = %q, want .gitignore", GitignoreFilename)
	}
	if APMModulesGitignorePattern != "apm_modules/" {
		t.Errorf("APMModulesGitignorePattern = %q, want apm_modules/", APMModulesGitignorePattern)
	}
}

func TestFileConstants_dirs(t *testing.T) {
	if APMDir != ".apm" {
		t.Errorf("APMDir = %q, want .apm", APMDir)
	}
	if GitHubDir != ".github" {
		t.Errorf("GitHubDir = %q, want .github", GitHubDir)
	}
	if ClaudeDir != ".claude" {
		t.Errorf("ClaudeDir = %q, want .claude", ClaudeDir)
	}
	if APMModulesDir != "apm_modules" {
		t.Errorf("APMModulesDir = %q, want apm_modules", APMModulesDir)
	}
}

func TestFileConstants_markdownFiles(t *testing.T) {
	for name, val := range map[string]string{
		"SkillMDFilename":  SkillMDFilename,
		"AgentsMDFilename": AgentsMDFilename,
		"ClaudeMDFilename": ClaudeMDFilename,
	} {
		if val == "" {
			t.Errorf("%s is empty", name)
		}
	}
	if SkillMDFilename != "SKILL.md" {
		t.Errorf("SkillMDFilename = %q, want SKILL.md", SkillMDFilename)
	}
}
