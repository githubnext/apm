package constants_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/constants"
)

func TestInstallMode_AllValue_Extra4(t *testing.T) {
	if constants.InstallModeAll != "all" {
		t.Errorf("expected 'all', got %q", constants.InstallModeAll)
	}
}

func TestInstallMode_APMValue_Extra4(t *testing.T) {
	if constants.InstallModeAPM != "apm" {
		t.Errorf("expected 'apm', got %q", constants.InstallModeAPM)
	}
}

func TestInstallMode_MCPValue_Extra4(t *testing.T) {
	if constants.InstallModeMCP != "mcp" {
		t.Errorf("expected 'mcp', got %q", constants.InstallModeMCP)
	}
}

func TestAPMYMLFilename_IsYML_Extra4(t *testing.T) {
	if !strings.HasSuffix(constants.APMYMLFilename, ".yml") {
		t.Errorf("expected .yml suffix, got %q", constants.APMYMLFilename)
	}
}

func TestAPMLockFilename_HasLock_Extra4(t *testing.T) {
	if !strings.Contains(constants.APMLockFilename, "lock") {
		t.Errorf("expected 'lock' in filename, got %q", constants.APMLockFilename)
	}
}

func TestAPMModulesDir_NotEmpty_Extra4(t *testing.T) {
	if constants.APMModulesDir == "" {
		t.Error("expected non-empty APMModulesDir")
	}
}

func TestAPMDir_StartsDot_Extra4(t *testing.T) {
	if !strings.HasPrefix(constants.APMDir, ".") {
		t.Errorf("expected dot prefix for APMDir, got %q", constants.APMDir)
	}
}

func TestSkillMDFilename_EndsMD_Extra4(t *testing.T) {
	if !strings.HasSuffix(constants.SkillMDFilename, ".md") {
		t.Errorf("expected .md suffix, got %q", constants.SkillMDFilename)
	}
}

func TestAgentsMDFilename_EndsMD_Extra4(t *testing.T) {
	if !strings.HasSuffix(constants.AgentsMDFilename, ".md") {
		t.Errorf("expected .md suffix, got %q", constants.AgentsMDFilename)
	}
}

func TestDefaultSkipDirs_GitignoreNotPresent_Extra4(t *testing.T) {
	if _, ok := constants.DefaultSkipDirs[".gitignore"]; ok {
		t.Error("expected .gitignore not in skip dirs")
	}
}

func TestDefaultSkipDirs_NodeModulesPresent_Extra4(t *testing.T) {
	if _, ok := constants.DefaultSkipDirs["node_modules"]; !ok {
		t.Error("expected node_modules in skip dirs")
	}
}

func TestDefaultSkipDirs_PycachePresent_Extra4(t *testing.T) {
	if _, ok := constants.DefaultSkipDirs["__pycache__"]; !ok {
		t.Error("expected __pycache__ in skip dirs")
	}
}

func TestGitHubDir_StartsDot_Extra4(t *testing.T) {
	if !strings.HasPrefix(constants.GitHubDir, ".") {
		t.Errorf("expected .github, got %q", constants.GitHubDir)
	}
}

func TestGitignoreFilename_StartsDot_Extra4(t *testing.T) {
	if !strings.HasPrefix(constants.GitignoreFilename, ".") {
		t.Errorf("expected dot prefix, got %q", constants.GitignoreFilename)
	}
}
