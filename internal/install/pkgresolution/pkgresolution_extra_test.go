package pkgresolution_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/pkgresolution"
)

// mockDep2 is an additional mock for extra tests (avoid redeclaration).
type mockDep2 struct {
	url         string
	virtualPath string
	ref         string
	alias       string
}

func (m *mockDep2) ToGitHubURL() string                               { return m.url }
func (m *mockDep2) GetVirtualPath() string                            { return m.virtualPath }
func (m *mockDep2) GetRef() string                                    { return m.ref }
func (m *mockDep2) GetAlias() string                                  { return m.alias }
func (m *mockDep2) NeedsGitLabDirectShorthandProbing(raw string) bool { return false }

func TestNormalizePackageSpec_Whitespace(t *testing.T) {
	cases := []struct{ in, want string }{
		{"\t owner/repo \t", "owner/repo"},
		{"  ", ""},
		{"\n\n", ""},
		{"owner/repo\n", "owner/repo"},
	}
	for _, tc := range cases {
		got := pkgresolution.NormalizePackageSpec(tc.in)
		if got != tc.want {
			t.Errorf("NormalizePackageSpec(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestResolutionError_WithCause(t *testing.T) {
	cause := errors.New("network error")
	err := &pkgresolution.ResolutionError{Package: "owner/pkg", Cause: cause}
	if !strings.Contains(err.Error(), "owner/pkg") {
		t.Errorf("error should contain package name: %q", err.Error())
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("error should contain cause: %q", err.Error())
	}
	if err.Unwrap() != cause {
		t.Error("Unwrap should return original cause")
	}
}

func TestResolutionError_NoCause(t *testing.T) {
	err := &pkgresolution.ResolutionError{Package: "my/pkg"}
	msg := err.Error()
	if !strings.Contains(msg, "my/pkg") {
		t.Errorf("error without cause should contain package name: %q", msg)
	}
	if err.Unwrap() != nil {
		t.Error("Unwrap with no cause should return nil")
	}
}

func TestYAMLEntry_AllFields(t *testing.T) {
	e := pkgresolution.YAMLEntry{
		Git:   "https://github.com/owner/repo",
		Path:  "sub/dir",
		Ref:   "v1.2.3",
		Alias: "my-alias",
	}
	if e.Git != "https://github.com/owner/repo" {
		t.Errorf("Git = %q", e.Git)
	}
	if e.Path != "sub/dir" {
		t.Errorf("Path = %q", e.Path)
	}
	if e.Ref != "v1.2.3" {
		t.Errorf("Ref = %q", e.Ref)
	}
	if e.Alias != "my-alias" {
		t.Errorf("Alias = %q", e.Alias)
	}
}

func TestDependencyReferenceToYAMLEntry_NoPath(t *testing.T) {
	dep := &mockDep2{url: "https://github.com/a/b", ref: "main"}
	entry := pkgresolution.DependencyReferenceToYAMLEntry(dep)
	if entry.Path != "" {
		t.Errorf("Path should be empty when virtualPath is empty, got %q", entry.Path)
	}
	if entry.Ref != "main" {
		t.Errorf("Ref = %q, want main", entry.Ref)
	}
}

func TestDependencyReferenceToYAMLEntry_WithAlias(t *testing.T) {
	dep := &mockDep2{url: "https://github.com/x/y", alias: "renamed"}
	entry := pkgresolution.DependencyReferenceToYAMLEntry(dep)
	if entry.Alias != "renamed" {
		t.Errorf("Alias = %q, want renamed", entry.Alias)
	}
}

func TestValidateGitParentScope_NotUserScope(t *testing.T) {
	dep := &mockDep2{url: "../parent"}
	for _, scope := range []string{"project", "global", "org"} {
		err := pkgresolution.ValidateGitParentScope(dep, scope)
		if err != nil {
			t.Errorf("scope=%q: expected no error, got %v", scope, err)
		}
	}
}

func TestValidateGitParentScope_UserScopeError(t *testing.T) {
	dep := &mockDep2{url: "../sibling"}
	err := pkgresolution.ValidateGitParentScope(dep, "user")
	if err == nil {
		t.Error("expected error for git parent at user scope")
	}
}

func TestIsGitParentAtUserScope_NonParent(t *testing.T) {
	dep := &mockDep2{url: "https://github.com/owner/repo"}
	if pkgresolution.IsGitParentAtUserScope(dep, "user") {
		t.Error("absolute URL should not be considered a git parent")
	}
}

func TestIsGitParentAtUserScope_ProjectScope(t *testing.T) {
	dep := &mockDep2{url: "../parent-pkg"}
	if pkgresolution.IsGitParentAtUserScope(dep, "project") {
		t.Error("project scope: parent dep is allowed")
	}
}

func TestGITParentUserScopeError_IsString(t *testing.T) {
	if pkgresolution.GITParentUserScopeError == "" {
		t.Error("GITParentUserScopeError constant should not be empty")
	}
}

func TestResolutionResult_Fields(t *testing.T) {
	dep := &mockDep2{url: "https://github.com/a/b"}
	res := pkgresolution.ResolutionResult{
		DepRef:                      dep,
		DirectGitLabVirtualResolved: true,
	}
	if !res.DirectGitLabVirtualResolved {
		t.Error("expected DirectGitLabVirtualResolved=true")
	}
	if res.DepRef == nil {
		t.Error("DepRef should not be nil")
	}
}
