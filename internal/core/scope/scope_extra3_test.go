package scope_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/scope"
)

func TestInstallScope_String_UserExtra3(t *testing.T) {
	if scope.ScopeUser.String() != "user" {
		t.Errorf("expected 'user', got %q", scope.ScopeUser.String())
	}
}

func TestInstallScope_String_ProjectExtra3(t *testing.T) {
	if scope.ScopeProject.String() != "project" {
		t.Errorf("expected 'project', got %q", scope.ScopeProject.String())
	}
}

func TestParseScope_UserString(t *testing.T) {
	s, ok := scope.ParseScope("user")
	if !ok {
		t.Error("expected ok for 'user'")
	}
	if s != scope.ScopeUser {
		t.Error("expected ScopeUser")
	}
}

func TestParseScope_InvalidString(t *testing.T) {
	_, ok := scope.ParseScope("global")
	if ok {
		t.Error("expected not ok for 'global'")
	}
}

func TestParseScope_MixedCase(t *testing.T) {
	s, ok := scope.ParseScope("USER")
	if !ok {
		t.Error("expected ok for 'USER'")
	}
	if s != scope.ScopeUser {
		t.Error("expected ScopeUser for 'USER'")
	}
}

func TestGetDeployRoot_ProjectScope(t *testing.T) {
	root, err := scope.GetDeployRoot(scope.ScopeProject)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if root == "" {
		t.Error("expected non-empty deploy root")
	}
}

func TestGetDeployRoot_UserScope(t *testing.T) {
	root, err := scope.GetDeployRoot(scope.ScopeUser)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if root == "" {
		t.Error("expected non-empty user home dir")
	}
}

func TestGetAPMDir_ProjectScope(t *testing.T) {
	dir, err := scope.GetAPMDir(scope.ScopeProject)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty APM dir")
	}
}
