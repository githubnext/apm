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
