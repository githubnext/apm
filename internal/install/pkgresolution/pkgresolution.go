// Package pkgresolution provides helpers for install-time package reference resolution.
// Migrated from src/apm_cli/install/package_resolution.py
package pkgresolution

import (
	"errors"
	"fmt"
	"strings"
)

// GITParentUserScopeError is returned when a git parent dependency is used at user scope.
const GITParentUserScopeError = "git: parent dependencies are not supported at user scope. " +
	"Use project scope or specify explicit git URL."

// YAMLEntry represents a serialized dependency reference for apm.yml storage.
type YAMLEntry struct {
	Git   string `json:"git"`
	Path  string `json:"path,omitempty"`
	Ref   string `json:"ref,omitempty"`
	Alias string `json:"alias,omitempty"`
}

// DependencyRef is the minimal interface used by resolution helpers.
type DependencyRef interface {
	// ToGitHubURL returns the canonical https://... clone URL.
	ToGitHubURL() string
	// GetVirtualPath returns the sub-path within the repo, or "".
	GetVirtualPath() string
	// GetRef returns the explicit git ref, or "".
	GetRef() string
	// GetAlias returns the user-supplied alias, or "".
	GetAlias() string
	// NeedsGitLabDirectShorthandProbing reports whether this ref requires GitLab probing.
	NeedsGitLabDirectShorthandProbing(raw string) bool
}

// DependencyReferenceToYAMLEntry serializes a dependency reference for apm.yml storage.
func DependencyReferenceToYAMLEntry(dep DependencyRef) YAMLEntry {
	entry := YAMLEntry{Git: dep.ToGitHubURL()}
	if vp := dep.GetVirtualPath(); vp != "" {
		entry.Path = vp
	}
	if ref := dep.GetRef(); ref != "" {
		entry.Ref = ref
	}
	if alias := dep.GetAlias(); alias != "" {
		entry.Alias = alias
	}
	return entry
}

// ResolutionResult is the outcome of resolving a raw package specifier.
type ResolutionResult struct {
	// DepRef is the resolved dependency reference.
	DepRef DependencyRef
	// DirectGitLabVirtualResolved is true when GitLab shorthand probing
	// produced a virtual path entry.
	DirectGitLabVirtualResolved bool
}

// ResolutionError wraps a failure to resolve a package specifier.
type ResolutionError struct {
	Package string
	Cause   error
}

func (e *ResolutionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("could not resolve %q: %s", e.Package, e.Cause.Error())
	}
	return fmt.Sprintf("could not resolve %q", e.Package)
}

func (e *ResolutionError) Unwrap() error { return e.Cause }

// NormalizePackageSpec trims whitespace and normalises the package specifier.
func NormalizePackageSpec(pkg string) string {
	return strings.TrimSpace(pkg)
}

// IsGitParentAtUserScope reports whether dep is a git parent dependency
// being added at user scope, which is not supported.
func IsGitParentAtUserScope(depRef DependencyRef, scope string) bool {
	if scope != "user" {
		return false
	}
	url := depRef.ToGitHubURL()
	return strings.Contains(url, "..") || strings.HasPrefix(url, "../")
}

// ValidateGitParentScope returns an error if a git parent dep is used at user scope.
func ValidateGitParentScope(depRef DependencyRef, scope string) error {
	if IsGitParentAtUserScope(depRef, scope) {
		return errors.New(GITParentUserScopeError)
	}
	return nil
}
