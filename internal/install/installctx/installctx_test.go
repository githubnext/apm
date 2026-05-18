package installctx_test

import (
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/installctx"
)

func TestNew(t *testing.T) {
	ctx := installctx.New("/project", "/project/.apm")
	if ctx.ProjectRoot != "/project" {
		t.Errorf("ProjectRoot: got %q, want %q", ctx.ProjectRoot, "/project")
	}
	if ctx.ApmDir != "/project/.apm" {
		t.Errorf("ApmDir: got %q, want %q", ctx.ApmDir, "/project/.apm")
	}
	if ctx.ParallelDownloads != 4 {
		t.Errorf("ParallelDownloads: got %d, want 4", ctx.ParallelDownloads)
	}
	if ctx.IntendedDepKeys == nil {
		t.Error("IntendedDepKeys should be initialized")
	}
	if ctx.PackageDeployedFiles == nil {
		t.Error("PackageDeployedFiles should be initialized")
	}
	if ctx.PackageTypes == nil {
		t.Error("PackageTypes should be initialized")
	}
	if ctx.PackageHashes == nil {
		t.Error("PackageHashes should be initialized")
	}
	if ctx.ExpectedHashChangeDeps == nil {
		t.Error("ExpectedHashChangeDeps should be initialized")
	}
}

func TestApmModulesDirOrDefault(t *testing.T) {
	ctx := installctx.New("/project", "/project/.apm")

	// When ApmModulesDir is not set, returns default
	got := ctx.ApmModulesDirOrDefault()
	want := filepath.Join("/project", "apm_modules")
	if got != want {
		t.Errorf("ApmModulesDirOrDefault: got %q, want %q", got, want)
	}

	// When ApmModulesDir is set, returns it
	ctx.ApmModulesDir = "/custom/modules"
	got = ctx.ApmModulesDirOrDefault()
	if got != "/custom/modules" {
		t.Errorf("ApmModulesDirOrDefault with custom: got %q, want %q", got, "/custom/modules")
	}
}

func TestLockfilePathOrDefault(t *testing.T) {
	ctx := installctx.New("/project", "/project/.apm")

	// When LockfilePath is not set, returns default
	got := ctx.LockfilePathOrDefault()
	want := filepath.Join("/project", "apm.lock.yaml")
	if got != want {
		t.Errorf("LockfilePathOrDefault: got %q, want %q", got, want)
	}

	// When LockfilePath is set, returns it
	ctx.LockfilePath = "/custom/apm.lock.yaml"
	got = ctx.LockfilePathOrDefault()
	if got != "/custom/apm.lock.yaml" {
		t.Errorf("LockfilePathOrDefault with custom: got %q, want %q", got, "/custom/apm.lock.yaml")
	}
}

func TestNew_BoolDefaults(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	if ctx.UpdateRefs {
		t.Error("UpdateRefs should default to false")
	}
	if ctx.DryRun {
		t.Error("DryRun should default to false")
	}
	if ctx.Force {
		t.Error("Force should default to false")
	}
	if ctx.Verbose {
		t.Error("Verbose should default to false")
	}
	if ctx.AllowInsecure {
		t.Error("AllowInsecure should default to false")
	}
}

func TestNew_EmptySlices(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	if ctx.AllowInsecureHosts == nil {
		t.Error("AllowInsecureHosts should not be nil")
	}
	if ctx.OnlyPackages == nil {
		t.Error("OnlyPackages should not be nil")
	}
	if ctx.OldLocalDeployed == nil {
		t.Error("OldLocalDeployed should not be nil")
	}
	if ctx.LocalDeployedFiles == nil {
		t.Error("LocalDeployedFiles should not be nil")
	}
}

func TestNew_ApmDir(t *testing.T) {
	ctx := installctx.New("/workspace", "/workspace/.apm")
	if ctx.ApmDir != "/workspace/.apm" {
		t.Errorf("ApmDir: got %q", ctx.ApmDir)
	}
}

func TestApmModulesDirOrDefault_Empty(t *testing.T) {
	ctx := installctx.New("/root", "/root/.apm")
	ctx.ApmModulesDir = ""
	got := ctx.ApmModulesDirOrDefault()
	want := filepath.Join("/root", "apm_modules")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLockfilePathOrDefault_Empty(t *testing.T) {
	ctx := installctx.New("/root", "/root/.apm")
	ctx.LockfilePath = ""
	got := ctx.LockfilePathOrDefault()
	want := filepath.Join("/root", "apm.lock.yaml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInstallContext_PolicyFields(t *testing.T) {
	ctx := installctx.New("/p", "/p/.apm")
	if ctx.PolicyEnforcementActive {
		t.Error("PolicyEnforcementActive should default to false")
	}
	if ctx.NoPolicy {
		t.Error("NoPolicy should default to false")
	}
	ctx.PolicyEnforcementActive = true
	ctx.NoPolicy = true
	if !ctx.PolicyEnforcementActive {
		t.Error("PolicyEnforcementActive should be settable")
	}
}
