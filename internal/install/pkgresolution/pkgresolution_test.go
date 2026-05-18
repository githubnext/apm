package pkgresolution_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/pkgresolution"
)

// mockDep implements DependencyRef for testing.
type mockDep struct {
	url         string
	virtualPath string
	ref         string
	alias       string
	needsProbe  bool
}

func (m *mockDep) ToGitHubURL() string                               { return m.url }
func (m *mockDep) GetVirtualPath() string                            { return m.virtualPath }
func (m *mockDep) GetRef() string                                    { return m.ref }
func (m *mockDep) GetAlias() string                                  { return m.alias }
func (m *mockDep) NeedsGitLabDirectShorthandProbing(raw string) bool { return m.needsProbe }

func TestNormalizePackageSpec(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  owner/repo  ", "owner/repo"},
		{"owner/repo", "owner/repo"},
		{"\towner/repo\n", "owner/repo"},
		{"", ""},
	}
	for _, tt := range tests {
		got := pkgresolution.NormalizePackageSpec(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizePackageSpec(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsGitParentAtUserScope(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		scope    string
		expected bool
	}{
		{"parent dep at user scope", "../sibling", "user", true},
		{"parent dep at project scope", "../sibling", "project", false},
		{"regular url at user scope", "https://github.com/owner/repo", "user", false},
		{"relative with ..", "../../pkg", "user", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &mockDep{url: tt.url}
			got := pkgresolution.IsGitParentAtUserScope(dep, tt.scope)
			if got != tt.expected {
				t.Errorf("IsGitParentAtUserScope() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateGitParentScope(t *testing.T) {
	dep := &mockDep{url: "../sibling"}

	err := pkgresolution.ValidateGitParentScope(dep, "user")
	if err == nil {
		t.Error("expected error for git parent at user scope")
	}
	if err.Error() != pkgresolution.GITParentUserScopeError {
		t.Errorf("unexpected error message: %q", err.Error())
	}

	err2 := pkgresolution.ValidateGitParentScope(dep, "project")
	if err2 != nil {
		t.Errorf("expected nil error for project scope, got %v", err2)
	}
}

func TestDependencyReferenceToYAMLEntry(t *testing.T) {
	t.Run("full entry", func(t *testing.T) {
		dep := &mockDep{
			url:         "https://github.com/owner/repo",
			virtualPath: "subdir",
			ref:         "main",
			alias:       "my-alias",
		}
		entry := pkgresolution.DependencyReferenceToYAMLEntry(dep)
		if entry.Git != "https://github.com/owner/repo" {
			t.Errorf("Git = %q", entry.Git)
		}
		if entry.Path != "subdir" {
			t.Errorf("Path = %q", entry.Path)
		}
		if entry.Ref != "main" {
			t.Errorf("Ref = %q", entry.Ref)
		}
		if entry.Alias != "my-alias" {
			t.Errorf("Alias = %q", entry.Alias)
		}
	})

	t.Run("minimal entry", func(t *testing.T) {
		dep := &mockDep{url: "https://github.com/a/b"}
		entry := pkgresolution.DependencyReferenceToYAMLEntry(dep)
		if entry.Git != "https://github.com/a/b" {
			t.Errorf("Git = %q", entry.Git)
		}
		if entry.Path != "" {
			t.Errorf("Path should be empty, got %q", entry.Path)
		}
		if entry.Ref != "" {
			t.Errorf("Ref should be empty, got %q", entry.Ref)
		}
	})
}

func TestResolutionError(t *testing.T) {
	err := &pkgresolution.ResolutionError{Package: "owner/repo"}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
	if err.Unwrap() != nil {
		t.Error("expected nil cause")
	}
}
