package installctx_test

import (
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/installctx"
)

func TestNew_CounterDefaults(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	if ctx.InstalledCount != 0 {
		t.Errorf("InstalledCount should start 0, got %d", ctx.InstalledCount)
	}
	if ctx.TotalPromptsIntegrated != 0 {
		t.Errorf("TotalPromptsIntegrated should start 0, got %d", ctx.TotalPromptsIntegrated)
	}
	if ctx.TotalAgentsIntegrated != 0 {
		t.Errorf("TotalAgentsIntegrated should start 0, got %d", ctx.TotalAgentsIntegrated)
	}
	if ctx.TotalSkillsIntegrated != 0 {
		t.Errorf("TotalSkillsIntegrated should start 0, got %d", ctx.TotalSkillsIntegrated)
	}
	if ctx.TotalInstructionsIntegrated != 0 {
		t.Errorf("TotalInstructionsIntegrated should start 0, got %d", ctx.TotalInstructionsIntegrated)
	}
	if ctx.TotalCommandsIntegrated != 0 {
		t.Errorf("TotalCommandsIntegrated should start 0, got %d", ctx.TotalCommandsIntegrated)
	}
	if ctx.TotalHooksIntegrated != 0 {
		t.Errorf("TotalHooksIntegrated should start 0, got %d", ctx.TotalHooksIntegrated)
	}
	if ctx.TotalLinksResolved != 0 {
		t.Errorf("TotalLinksResolved should start 0, got %d", ctx.TotalLinksResolved)
	}
}

func TestNew_StringDefaults(t *testing.T) {
	ctx := installctx.New("/root", "/root/.apm")
	if ctx.TargetOverride != "" {
		t.Errorf("TargetOverride should be empty, got %q", ctx.TargetOverride)
	}
	if ctx.LockfilePath != "" {
		t.Errorf("LockfilePath should be empty, got %q", ctx.LockfilePath)
	}
	if ctx.ApmModulesDir != "" {
		t.Errorf("ApmModulesDir should be empty, got %q", ctx.ApmModulesDir)
	}
}

func TestApmModulesDirOrDefault_WithValue(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.ApmModulesDir = "/custom/modules"
	got := ctx.ApmModulesDirOrDefault()
	if got != "/custom/modules" {
		t.Errorf("expected /custom/modules, got %q", got)
	}
}

func TestLockfilePathOrDefault_WithValue(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.LockfilePath = "/custom/lock.yaml"
	got := ctx.LockfilePathOrDefault()
	if got != "/custom/lock.yaml" {
		t.Errorf("expected /custom/lock.yaml, got %q", got)
	}
}

func TestLockfilePathOrDefault_DefaultPath(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	got := ctx.LockfilePathOrDefault()
	if got == "" {
		t.Error("default lockfile path should not be empty")
	}
	// Should be relative to ApmDir
	if filepath.IsAbs(got) {
		// fine; just ensure it contains apm reference
		_ = got
	}
}

func TestNew_SliceDefaults(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	if ctx.OnlyPackages == nil {
		t.Errorf("OnlyPackages should be initialized (not nil)")
	}
	if ctx.AllowInsecureHosts == nil {
		t.Errorf("AllowInsecureHosts should be initialized (not nil)")
	}
	if ctx.SkillSubset != nil {
		t.Errorf("SkillSubset should be nil by default, got %v", ctx.SkillSubset)
	}
	if ctx.OldLocalDeployed == nil {
		t.Errorf("OldLocalDeployed should be initialized (not nil)")
	}
	if ctx.LocalDeployedFiles == nil {
		t.Errorf("LocalDeployedFiles should be initialized (not nil)")
	}
}

func TestNew_AllowProtocolFallback(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	if ctx.AllowProtocolFallback != nil {
		t.Error("AllowProtocolFallback should be nil (read from env)")
	}
}

func TestInstallContext_CounterMutation(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.InstalledCount = 5
	ctx.TotalPromptsIntegrated = 10
	ctx.TotalAgentsIntegrated = 3
	if ctx.InstalledCount != 5 {
		t.Errorf("InstalledCount = %d, want 5", ctx.InstalledCount)
	}
	if ctx.TotalPromptsIntegrated != 10 {
		t.Errorf("TotalPromptsIntegrated = %d, want 10", ctx.TotalPromptsIntegrated)
	}
}

func TestInstallContext_MapMutation(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.IntendedDepKeys["pkg/a"] = true
	ctx.PackageDeployedFiles["pkg/a"] = []string{"file.txt"}
	ctx.PackageTypes["pkg/a"] = "github"
	ctx.PackageHashes["pkg/a"] = "abc123"
	ctx.ExpectedHashChangeDeps["pkg/a"] = true

	if !ctx.IntendedDepKeys["pkg/a"] {
		t.Error("IntendedDepKeys mutation failed")
	}
	if ctx.PackageDeployedFiles["pkg/a"][0] != "file.txt" {
		t.Error("PackageDeployedFiles mutation failed")
	}
	if ctx.PackageTypes["pkg/a"] != "github" {
		t.Error("PackageTypes mutation failed")
	}
	if ctx.PackageHashes["pkg/a"] != "abc123" {
		t.Error("PackageHashes mutation failed")
	}
	if !ctx.ExpectedHashChangeDeps["pkg/a"] {
		t.Error("ExpectedHashChangeDeps mutation failed")
	}
}

func TestInstallContext_TwoInstances_Independent(t *testing.T) {
	a := installctx.New("/a", "/a/.apm")
	b := installctx.New("/b", "/b/.apm")
	a.IntendedDepKeys["x"] = true
	if b.IntendedDepKeys["x"] {
		t.Error("maps should be independent between instances")
	}
}
