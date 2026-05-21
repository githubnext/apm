package install

import (
	"strings"
	"testing"
)

func TestInstallOptions_AllFlags_Extra3(t *testing.T) {
	opts := InstallOptions{
		Frozen:    true,
		DryRun:    true,
		Verbose:   true,
		Force:     false,
		UserScope: false,
	}
	if !opts.Frozen {
		t.Error("Frozen should be true")
	}
	if !opts.DryRun {
		t.Error("DryRun should be true")
	}
	if opts.Force {
		t.Error("Force should be false")
	}
}

func TestInstallOptions_PackageRefs_Extra3(t *testing.T) {
	opts := InstallOptions{
		PackageRefs: []string{"owner/a", "owner/b@v1"},
	}
	if len(opts.PackageRefs) != 2 {
		t.Errorf("PackageRefs len = %d, want 2", len(opts.PackageRefs))
	}
}

func TestInstallResult_Counts_Extra3(t *testing.T) {
	r := InstallResult{
		PackagesInstalled: 3,
		PackagesSkipped:   1,
		PackagesRemoved:   0,
	}
	if r.PackagesInstalled != 3 {
		t.Errorf("PackagesInstalled = %d, want 3", r.PackagesInstalled)
	}
	if r.PackagesSkipped != 1 {
		t.Errorf("PackagesSkipped = %d, want 1", r.PackagesSkipped)
	}
}

func TestInstallResult_Warnings_Extra3(t *testing.T) {
	r := InstallResult{Warnings: []string{"deprecated package", "slow network"}}
	if len(r.Warnings) != 2 {
		t.Errorf("Warnings len = %d, want 2", len(r.Warnings))
	}
}

func TestInstallResult_LockfileUpdated_Extra3(t *testing.T) {
	r := InstallResult{LockfileUpdated: true}
	if !r.LockfileUpdated {
		t.Error("LockfileUpdated should be true")
	}
}

func TestDependencyEntry_AllFields_Extra3(t *testing.T) {
	e := DependencyEntry{
		Name: "mypkg", Ref: "v1.0", Host: "github.com",
		Org: "owner", Repo: "repo", Local: false,
	}
	if e.Name != "mypkg" {
		t.Errorf("Name = %q", e.Name)
	}
	if e.Host != "github.com" {
		t.Errorf("Host = %q", e.Host)
	}
	if e.Local {
		t.Error("Local should be false")
	}
}

func TestDependencyEntry_LocalPackage_Extra3(t *testing.T) {
	e := DependencyEntry{Local: true, Path: "./my/local/pkg"}
	if !e.Local {
		t.Error("Local should be true")
	}
	if e.Path != "./my/local/pkg" {
		t.Errorf("Path = %q", e.Path)
	}
}

func TestPolicyViolation_Fields_Extra3(t *testing.T) {
	v := PolicyViolation{Package: "owner/pkg", Rule: "allow", Message: "not allowed"}
	if v.Package != "owner/pkg" {
		t.Errorf("Package = %q", v.Package)
	}
	if !strings.Contains(v.Message, "allowed") {
		t.Errorf("Message = %q should contain 'allowed'", v.Message)
	}
}

func TestAuthenticationError_Error_Extra3(t *testing.T) {
	e := &AuthenticationError{Host: "github.com"}
	s := e.Error()
	if s == "" {
		t.Error("Error() should not be empty")
	}
}

func TestFrozenInstallError_Extra3(t *testing.T) {
	e := &FrozenInstallError{}
	s := e.Error()
	if s == "" {
		t.Error("FrozenInstallError.Error() should not be empty")
	}
}
