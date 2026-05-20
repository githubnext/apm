package pkgresolution_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/install/pkgresolution"
)

type mockDep3 struct {
	url         string
	virtualPath string
	ref         string
	alias       string
	needsProbe  bool
}

func (m *mockDep3) ToGitHubURL() string                               { return m.url }
func (m *mockDep3) GetVirtualPath() string                            { return m.virtualPath }
func (m *mockDep3) GetRef() string                                    { return m.ref }
func (m *mockDep3) GetAlias() string                                  { return m.alias }
func (m *mockDep3) NeedsGitLabDirectShorthandProbing(raw string) bool { return m.needsProbe }

func TestNormalizePackageSpec_AlreadyTrimmed(t *testing.T) {
	result := pkgresolution.NormalizePackageSpec("owner/repo")
	if result != "owner/repo" {
		t.Errorf("expected owner/repo, got %s", result)
	}
}

func TestNormalizePackageSpec_Empty(t *testing.T) {
	result := pkgresolution.NormalizePackageSpec("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestNormalizePackageSpec_OnlySpaces(t *testing.T) {
	result := pkgresolution.NormalizePackageSpec("   ")
	if result != "" {
		t.Errorf("expected empty string after trimming spaces, got %q", result)
	}
}

func TestResolutionError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &pkgresolution.ResolutionError{Package: "pkg", Cause: cause}
	if errors.Unwrap(err) != cause {
		t.Error("Unwrap should return the cause error")
	}
}

func TestResolutionError_NilCause_Extra2(t *testing.T) {
	err := &pkgresolution.ResolutionError{Package: "mypkg"}
	msg := err.Error()
	if msg == "" {
		t.Error("Error() should return non-empty message even without cause")
	}
}

func TestDependencyReferenceToYAMLEntry_WithPath_Extra2(t *testing.T) {
	dep := &mockDep3{url: "https://github.com/org/repo", virtualPath: "subdir"}
	entry := pkgresolution.DependencyReferenceToYAMLEntry(dep)
	if entry.Git != "https://github.com/org/repo" {
		t.Error("Git URL mismatch")
	}
	if entry.Path != "subdir" {
		t.Errorf("expected Path=subdir, got %s", entry.Path)
	}
}

func TestDependencyReferenceToYAMLEntry_WithRef_Extra2(t *testing.T) {
	dep := &mockDep3{url: "https://github.com/org/repo", ref: "v2.0.0"}
	entry := pkgresolution.DependencyReferenceToYAMLEntry(dep)
	if entry.Ref != "v2.0.0" {
		t.Errorf("expected Ref=v2.0.0, got %s", entry.Ref)
	}
	if entry.Path != "" {
		t.Errorf("expected empty Path, got %s", entry.Path)
	}
}

func TestResolutionResult_DirectGitLab(t *testing.T) {
	dep := &mockDep3{url: "https://gitlab.com/org/repo"}
	result := pkgresolution.ResolutionResult{DepRef: dep, DirectGitLabVirtualResolved: true}
	if !result.DirectGitLabVirtualResolved {
		t.Error("DirectGitLabVirtualResolved should be true")
	}
}

func TestYAMLEntry_ZeroValue(t *testing.T) {
	var e pkgresolution.YAMLEntry
	if e.Git != "" || e.Path != "" || e.Ref != "" || e.Alias != "" {
		t.Error("YAMLEntry zero value should have empty fields")
	}
}
