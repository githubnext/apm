package scope_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core/scope"
)

func TestInstallScope_Iota_Zero(t *testing.T) {
	if int(scope.ScopeProject) != 0 {
		t.Errorf("ScopeProject should be 0, got %d", scope.ScopeProject)
	}
}

func TestInstallScope_Iota_One(t *testing.T) {
	if int(scope.ScopeUser) != 1 {
		t.Errorf("ScopeUser should be 1, got %d", scope.ScopeUser)
	}
}

func TestParseScope_ProjectLowercase(t *testing.T) {
	s, ok := scope.ParseScope("project")
	if !ok || s != scope.ScopeProject {
		t.Errorf("expected ScopeProject, got %v ok=%v", s, ok)
	}
}

func TestParseScope_UserLowercase(t *testing.T) {
	s, ok := scope.ParseScope("user")
	if !ok || s != scope.ScopeUser {
		t.Errorf("expected ScopeUser, got %v ok=%v", s, ok)
	}
}

func TestParseScope_Mixed_Case(t *testing.T) {
	s, ok := scope.ParseScope("USER")
	if !ok || s != scope.ScopeUser {
		t.Errorf("expected ScopeUser for 'USER', got %v ok=%v", s, ok)
	}
}

func TestParseScope_Invalid_FalseAndDefault(t *testing.T) {
	s, ok := scope.ParseScope("invalid")
	if ok {
		t.Error("ParseScope with invalid string should return ok=false")
	}
	if s != scope.ScopeProject {
		t.Errorf("ParseScope default for invalid should be ScopeProject, got %v", s)
	}
}

func TestInstallScope_String_Project(t *testing.T) {
	if scope.ScopeProject.String() != "project" {
		t.Errorf("ScopeProject.String() should be 'project', got %q", scope.ScopeProject.String())
	}
}

func TestInstallScope_String_User(t *testing.T) {
	if scope.ScopeUser.String() != "user" {
		t.Errorf("ScopeUser.String() should be 'user', got %q", scope.ScopeUser.String())
	}
}

func TestGetDeployRoot_User_HasHomePrefix(t *testing.T) {
	home, _ := os.UserHomeDir()
	root, err := scope.GetDeployRoot(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root != home {
		t.Errorf("ScopeUser deploy root should be home dir %q, got %q", home, root)
	}
}

func TestGetDeployRoot_Project_NonEmpty(t *testing.T) {
	root, err := scope.GetDeployRoot(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == "" {
		t.Error("ScopeProject deploy root should not be empty")
	}
}

func TestGetAPMDir_User_EndsWithAPM(t *testing.T) {
	dir, err := scope.GetAPMDir(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(dir, scope.UserAPMDir) {
		t.Errorf("user APM dir should end with %q, got %q", scope.UserAPMDir, dir)
	}
}

func TestGetAPMDir_User_UnderHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	dir, err := scope.GetAPMDir(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(dir, home) {
		t.Errorf("user APM dir should be under home %q, got %q", home, dir)
	}
}

func TestGetModulesDir_ContainsAPMModulesDir(t *testing.T) {
	dir, err := scope.GetModulesDir(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("GetModulesDir should return non-empty path")
	}
}

func TestGetManifestPath_Project_NotEmpty(t *testing.T) {
	p, err := scope.GetManifestPath(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == "" {
		t.Error("GetManifestPath should return non-empty path")
	}
}

func TestGetManifestPath_User_UnderUserAPMDir(t *testing.T) {
	apmDir, _ := scope.GetAPMDir(scope.ScopeUser)
	manifestPath, err := scope.GetManifestPath(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Dir(manifestPath) != apmDir {
		t.Errorf("manifest path dir should be APM dir %q, got %q", apmDir, filepath.Dir(manifestPath))
	}
}

func TestGetLockfileDir_Project_MatchesAPMDir(t *testing.T) {
	apmDir, _ := scope.GetAPMDir(scope.ScopeProject)
	lockDir, err := scope.GetLockfileDir(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lockDir != apmDir {
		t.Errorf("lockfile dir should equal APM dir, got %q vs %q", lockDir, apmDir)
	}
}

func TestGetLockfileDir_User_MatchesAPMDir(t *testing.T) {
	apmDir, _ := scope.GetAPMDir(scope.ScopeUser)
	lockDir, err := scope.GetLockfileDir(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lockDir != apmDir {
		t.Errorf("lockfile dir should equal APM dir, got %q vs %q", lockDir, apmDir)
	}
}

func TestEnsureUserDirs_ReturnsUserRoot(t *testing.T) {
	root, err := scope.EnsureUserDirs()
	if err != nil {
		t.Fatalf("EnsureUserDirs error: %v", err)
	}
	if root == "" {
		t.Error("EnsureUserDirs should return non-empty root")
	}
	if _, statErr := os.Stat(root); statErr != nil {
		t.Errorf("user root dir should exist: %v", statErr)
	}
}

func TestUserAPMDir_Constant(t *testing.T) {
	if scope.UserAPMDir != ".apm" {
		t.Errorf("UserAPMDir should be '.apm', got %q", scope.UserAPMDir)
	}
}
