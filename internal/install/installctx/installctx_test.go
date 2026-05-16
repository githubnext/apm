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
