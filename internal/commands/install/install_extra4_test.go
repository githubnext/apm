package install

import "testing"

func TestInstallOptions_AllFields_Extra4(t *testing.T) {
	opts := InstallOptions{
		ProjectRoot:  "/project",
		PackageRefs:  []string{"pkg1", "pkg2"},
		Targets:      []string{"copilot"},
		Frozen:       true,
		DryRun:       false,
		Verbose:      true,
		Force:        false,
		UserScope:    true,
		NoProgress:   false,
		SkipLockfile: true,
		Mode:         InstallModeAll,
		AuthToken:    "tok123",
		ConcurrentDL: 4,
	}
	if opts.ProjectRoot != "/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if len(opts.PackageRefs) != 2 {
		t.Errorf("PackageRefs len = %d", len(opts.PackageRefs))
	}
	if !opts.Frozen {
		t.Error("Frozen should be true")
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
	if !opts.UserScope {
		t.Error("UserScope should be true")
	}
	if !opts.SkipLockfile {
		t.Error("SkipLockfile should be true")
	}
	if opts.Mode != InstallModeAll {
		t.Errorf("Mode = %q", opts.Mode)
	}
	if opts.ConcurrentDL != 4 {
		t.Errorf("ConcurrentDL = %d", opts.ConcurrentDL)
	}
}

func TestInstallOptions_ZeroValue_Extra4(t *testing.T) {
	var opts InstallOptions
	if opts.Frozen {
		t.Error("zero Frozen should be false")
	}
	if opts.ConcurrentDL != 0 {
		t.Errorf("zero ConcurrentDL = %d", opts.ConcurrentDL)
	}
}

func TestInstallMode_Constants_Extra4(t *testing.T) {
	if InstallModeAll == InstallModePrimitives {
		t.Error("InstallModeAll and InstallModePrimitives should differ")
	}
	if InstallModeAll == InstallModeClients {
		t.Error("InstallModeAll and InstallModeClients should differ")
	}
	if InstallModePrimitives == InstallModeClients {
		t.Error("InstallModePrimitives and InstallModeClients should differ")
	}
}

func TestInstallResult_ZeroValue_Extra4(t *testing.T) {
	var r InstallResult
	if r.PackagesInstalled != 0 {
		t.Errorf("zero PackagesInstalled = %d", r.PackagesInstalled)
	}
	if r.LockfileUpdated {
		t.Error("zero LockfileUpdated should be false")
	}
}

func TestInstallResult_Fields_Extra4(t *testing.T) {
	r := InstallResult{
		PackagesInstalled: 3,
		PackagesSkipped:   1,
		PackagesRemoved:   0,
		FilesWritten:      []string{"a.yml", "b.yml"},
		LockfileUpdated:   true,
		DurationSeconds:   2.5,
		Warnings:          []string{"warn1"},
		Errors:            []string{},
	}
	if r.PackagesInstalled != 3 {
		t.Errorf("PackagesInstalled = %d", r.PackagesInstalled)
	}
	if !r.LockfileUpdated {
		t.Error("LockfileUpdated should be true")
	}
	if len(r.FilesWritten) != 2 {
		t.Errorf("FilesWritten len = %d", len(r.FilesWritten))
	}
	if r.DurationSeconds != 2.5 {
		t.Errorf("DurationSeconds = %f", r.DurationSeconds)
	}
}

func TestDependencyEntry_Fields_Extra4(t *testing.T) {
	d := DependencyEntry{Name: "copilot-kit", Ref: "v1.2.3"}
	if d.Name != "copilot-kit" {
		t.Errorf("Name = %q", d.Name)
	}
	if d.Ref != "v1.2.3" {
		t.Errorf("Ref = %q", d.Ref)
	}
}
