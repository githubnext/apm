package constants

import (
	"strings"
	"testing"
)

func TestInstallMode_StringValues(t *testing.T) {
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
			t.Errorf("expected %q, got %q", c.want, string(c.mode))
		}
	}
}

func TestAPMModulesDir_NotHidden(t *testing.T) {
	if strings.HasPrefix(APMModulesDir, ".") {
		t.Errorf("APMModulesDir should not be hidden, got %q", APMModulesDir)
	}
}

func TestDefaultSkipDirs_APMModules(t *testing.T) {
	if _, ok := DefaultSkipDirs["apm_modules"]; !ok {
		t.Error("expected 'apm_modules' in DefaultSkipDirs")
	}
}

func TestDefaultSkipDirs_Venv(t *testing.T) {
	if _, ok := DefaultSkipDirs["venv"]; !ok {
		t.Error("expected 'venv' in DefaultSkipDirs")
	}
}

func TestDefaultSkipDirs_DotVenv(t *testing.T) {
	if _, ok := DefaultSkipDirs[".venv"]; !ok {
		t.Error("expected '.venv' in DefaultSkipDirs")
	}
}

func TestDefaultSkipDirs_Build(t *testing.T) {
	if _, ok := DefaultSkipDirs["build"]; !ok {
		t.Error("expected 'build' in DefaultSkipDirs")
	}
}

func TestDefaultSkipDirs_Dist(t *testing.T) {
	if _, ok := DefaultSkipDirs["dist"]; !ok {
		t.Error("expected 'dist' in DefaultSkipDirs")
	}
}

func TestDefaultSkipDirs_NonPresent(t *testing.T) {
	if _, ok := DefaultSkipDirs["src"]; ok {
		t.Error("'src' should not be in DefaultSkipDirs")
	}
}

func TestAPMYMLFilename_NonEmpty(t *testing.T) {
	if APMYMLFilename == "" {
		t.Error("APMYMLFilename should not be empty")
	}
}

func TestAPMLockFilename_NonEmpty(t *testing.T) {
	if APMLockFilename == "" {
		t.Error("APMLockFilename should not be empty")
	}
}

func TestInstallModeAll_NotAPM(t *testing.T) {
	if InstallModeAll == InstallModeAPM {
		t.Error("InstallModeAll and InstallModeAPM should be distinct")
	}
}

func TestInstallModeMCP_NotAPM(t *testing.T) {
	if InstallModeMCP == InstallModeAPM {
		t.Error("InstallModeMCP and InstallModeAPM should be distinct")
	}
}

func TestDefaultSkipDirs_IsMap(t *testing.T) {
	if DefaultSkipDirs == nil {
		t.Error("DefaultSkipDirs should not be nil")
	}
}

func TestDefaultSkipDirs_MypyCache(t *testing.T) {
	if _, ok := DefaultSkipDirs[".mypy_cache"]; !ok {
		t.Error("expected '.mypy_cache' in DefaultSkipDirs")
	}
}
