package scope_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/scope"
)

func TestParseScope(t *testing.T) {
	tests := []struct {
		input   string
		want    scope.InstallScope
		ok      bool
	}{
		{"project", scope.ScopeProject, true},
		{"user", scope.ScopeUser, true},
		{"USER", scope.ScopeUser, true},
		{"PROJECT", scope.ScopeProject, true},
		{"", scope.ScopeProject, false},
		{"global", scope.ScopeProject, false},
	}
	for _, tt := range tests {
		got, ok := scope.ParseScope(tt.input)
		if ok != tt.ok {
			t.Errorf("ParseScope(%q) ok=%v, want %v", tt.input, ok, tt.ok)
		}
		if ok && got != tt.want {
			t.Errorf("ParseScope(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestInstallScopeString(t *testing.T) {
	if scope.ScopeProject.String() != "project" {
		t.Errorf("expected 'project', got %q", scope.ScopeProject.String())
	}
	if scope.ScopeUser.String() != "user" {
		t.Errorf("expected 'user', got %q", scope.ScopeUser.String())
	}
}

func TestGetDeployRoot_User(t *testing.T) {
	root, err := scope.GetDeployRoot(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == "" {
		t.Error("expected non-empty user deploy root")
	}
}

func TestGetDeployRoot_Project(t *testing.T) {
	root, err := scope.GetDeployRoot(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == "" {
		t.Error("expected non-empty project deploy root")
	}
}

func TestGetAPMDir_Project(t *testing.T) {
	dir, err := scope.GetAPMDir(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty APM dir")
	}
}

func TestGetModulesDir_Project(t *testing.T) {
	dir, err := scope.GetModulesDir(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty modules dir")
	}
}

func TestGetModulesDir_User(t *testing.T) {
	dir, err := scope.GetModulesDir(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty user modules dir")
	}
}

func TestGetManifestPath_Project(t *testing.T) {
	path, err := scope.GetManifestPath(scope.ScopeProject)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty manifest path")
	}
}

func TestGetManifestPath_User(t *testing.T) {
	path, err := scope.GetManifestPath(scope.ScopeUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty user manifest path")
	}
}

func TestGetLockfileDir_Both(t *testing.T) {
	for _, s := range []scope.InstallScope{scope.ScopeProject, scope.ScopeUser} {
		dir, err := scope.GetLockfileDir(s)
		if err != nil {
			t.Errorf("GetLockfileDir(%v) error: %v", s, err)
		}
		if dir == "" {
			t.Errorf("GetLockfileDir(%v) returned empty string", s)
		}
	}
}

func TestEnsureUserDirs(t *testing.T) {
	root, err := scope.EnsureUserDirs()
	if err != nil {
		t.Fatalf("EnsureUserDirs error: %v", err)
	}
	if root == "" {
		t.Error("EnsureUserDirs returned empty root")
	}
}

func TestScopeScopeString_AllValues(t *testing.T) {
	if scope.ScopeProject.String() == scope.ScopeUser.String() {
		t.Error("ScopeProject and ScopeUser should have different String() values")
	}
}
