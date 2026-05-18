package install

import (
	"strings"
	"testing"
)

func TestAuthenticationError_Error(t *testing.T) {
	e := &AuthenticationError{Host: "github.com", Message: "token expired"}
	got := e.Error()
	if !strings.Contains(got, "github.com") || !strings.Contains(got, "token expired") {
		t.Errorf("unexpected error string: %q", got)
	}
}

func TestAuthenticationError_EmptyHost(t *testing.T) {
	e := &AuthenticationError{Host: "", Message: "no token"}
	got := e.Error()
	if !strings.Contains(got, "no token") {
		t.Errorf("expected 'no token' in %q", got)
	}
}

func TestFrozenInstallError_Error_Single(t *testing.T) {
	e := &FrozenInstallError{Changed: []string{"pkg-a"}}
	got := e.Error()
	if !strings.Contains(got, "frozen") || !strings.Contains(got, "1") {
		t.Errorf("unexpected error string: %q", got)
	}
}

func TestFrozenInstallError_Error_Multiple(t *testing.T) {
	e := &FrozenInstallError{Changed: []string{"pkg-a", "pkg-b", "pkg-c"}}
	got := e.Error()
	if !strings.Contains(got, "3") {
		t.Errorf("expected '3' in error string: %q", got)
	}
}

func TestFrozenInstallError_Error_Empty(t *testing.T) {
	e := &FrozenInstallError{Changed: nil}
	got := e.Error()
	if got == "" {
		t.Error("expected non-empty error string")
	}
}

func TestPolicyViolationError_Single(t *testing.T) {
	e := &PolicyViolationError{Violations: []PolicyViolation{{Message: "blocked"}}}
	got := e.Error()
	if !strings.Contains(got, "blocked") {
		t.Errorf("expected 'blocked' in %q", got)
	}
}

func TestPolicyViolationError_Multiple(t *testing.T) {
	e := &PolicyViolationError{Violations: []PolicyViolation{
		{Message: "blocked1"},
		{Message: "blocked2"},
	}}
	got := e.Error()
	if !strings.Contains(got, "2") {
		t.Errorf("expected '2' in %q", got)
	}
}

func TestPolicyViolationError_Empty(t *testing.T) {
	e := &PolicyViolationError{Violations: nil}
	got := e.Error()
	if got == "" {
		t.Error("expected non-empty error string")
	}
}

func TestMapToEntry_AllFields(t *testing.T) {
	m := map[string]string{
		"name": "mypkg",
		"ref":  "v1.0.0",
		"host": "github.com",
		"org":  "myorg",
		"repo": "myrepo",
	}
	e := mapToEntry(m)
	if e.Name != "mypkg" {
		t.Errorf("name: got %q", e.Name)
	}
	if e.Ref != "v1.0.0" {
		t.Errorf("ref: got %q", e.Ref)
	}
	if e.Host != "github.com" {
		t.Errorf("host: got %q", e.Host)
	}
	if e.Org != "myorg" {
		t.Errorf("org: got %q", e.Org)
	}
	if e.Repo != "myrepo" {
		t.Errorf("repo: got %q", e.Repo)
	}
}

func TestMapToEntry_Empty(t *testing.T) {
	e := mapToEntry(map[string]string{})
	if e.Name != "" || e.Ref != "" {
		t.Errorf("expected empty entry, got %+v", e)
	}
}

func TestMergeDependencies_NewEntry(t *testing.T) {
	existing := []DependencyEntry{{Name: "pkg-a"}}
	additions := []DependencyEntry{{Name: "pkg-b"}}
	result := mergeDependencies(existing, additions)
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestMergeDependencies_UpdateExisting(t *testing.T) {
	existing := []DependencyEntry{{Name: "pkg-a", Ref: "v1"}}
	additions := []DependencyEntry{{Name: "pkg-a", Ref: "v2"}}
	result := mergeDependencies(existing, additions)
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
	if result[0].Ref != "v2" {
		t.Errorf("expected ref v2, got %q", result[0].Ref)
	}
}

func TestMergeDependencies_Empty(t *testing.T) {
	result := mergeDependencies(nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestMergeDependencies_EmptyAdditions(t *testing.T) {
	existing := []DependencyEntry{{Name: "pkg-a"}}
	result := mergeDependencies(existing, nil)
	if len(result) != 1 || result[0].Name != "pkg-a" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestMergeDependencies_MultipleAdditions(t *testing.T) {
	existing := []DependencyEntry{{Name: "a"}, {Name: "b"}}
	additions := []DependencyEntry{{Name: "c"}, {Name: "b", Ref: "new"}, {Name: "d"}}
	result := mergeDependencies(existing, additions)
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
}

func TestInstallMode_Constants(t *testing.T) {
	if InstallModeAll != "all" {
		t.Errorf("InstallModeAll: got %q", InstallModeAll)
	}
	if InstallModePrimitives != "primitives" {
		t.Errorf("InstallModePrimitives: got %q", InstallModePrimitives)
	}
	if InstallModeClients != "clients" {
		t.Errorf("InstallModeClients: got %q", InstallModeClients)
	}
}

func TestParseDependencyRefs_WithRef(t *testing.T) {
	entries := parseDependencyRefs([]string{"myorg/myrepo@v2.0.0"})
	if len(entries) != 1 {
		t.Fatalf("expected 1, got %d", len(entries))
	}
	if entries[0].Ref != "v2.0.0" {
		t.Errorf("expected ref v2.0.0, got %q", entries[0].Ref)
	}
}

func TestParseDependencyRefs_Multiple(t *testing.T) {
	entries := parseDependencyRefs([]string{"pkg-a", "myorg/myrepo", "github.com/foo/bar"})
	if len(entries) != 3 {
		t.Fatalf("expected 3, got %d", len(entries))
	}
}

func TestParseDependencyRefs_Empty(t *testing.T) {
	entries := parseDependencyRefs(nil)
	if len(entries) != 0 {
		t.Errorf("expected 0, got %d", len(entries))
	}
}

func TestInstallOptions_Defaults(t *testing.T) {
	opts := InstallOptions{}
	if opts.Frozen || opts.DryRun || opts.Verbose || opts.Force {
		t.Error("expected all bool fields false by default")
	}
	if opts.ConcurrentDL != 0 {
		t.Errorf("expected 0 ConcurrentDL, got %d", opts.ConcurrentDL)
	}
}

func TestInstallResult_Defaults(t *testing.T) {
	r := InstallResult{}
	if r.PackagesInstalled != 0 || r.LockfileUpdated {
		t.Error("expected zero-value InstallResult")
	}
	if len(r.Warnings) != 0 || len(r.Errors) != 0 {
		t.Error("expected empty warnings/errors")
	}
}

func TestDependencyEntry_Fields(t *testing.T) {
	d := DependencyEntry{
		Name: "foo",
		Ref:  "main",
		Host: "gitlab.com",
		Org:  "org",
		Repo: "repo",
	}
	if d.Name != "foo" || d.Ref != "main" || d.Host != "gitlab.com" {
		t.Errorf("unexpected fields: %+v", d)
	}
}

func TestPolicyViolation_Fields(t *testing.T) {
	v := PolicyViolation{
		Package: "bad-pkg",
		Rule:    "no-banned",
		Message: "package is banned",
	}
	if v.Package != "bad-pkg" || v.Rule != "no-banned" || v.Message != "package is banned" {
		t.Errorf("unexpected fields: %+v", v)
	}
}
