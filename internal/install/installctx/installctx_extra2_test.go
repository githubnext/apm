package installctx_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/installctx"
)

func TestInstallContext_BoolFlags_DefaultFalse(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	if ctx.UpdateRefs {
		t.Error("UpdateRefs should default false")
	}
	if ctx.AllowInsecure {
		t.Error("AllowInsecure should default false")
	}
	if ctx.DryRun {
		t.Error("DryRun should default false")
	}
	if ctx.Force {
		t.Error("Force should default false")
	}
	if ctx.Verbose {
		t.Error("Verbose should default false")
	}
	if ctx.Dev {
		t.Error("Dev should default false")
	}
}

func TestInstallContext_PolicyEnforcement(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.PolicyEnforcementActive = true
	ctx.NoPolicy = false
	if !ctx.PolicyEnforcementActive {
		t.Error("PolicyEnforcementActive should be set")
	}
}

func TestInstallContext_SkillSubset(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.SkillSubset = []string{"skill1", "skill2"}
	ctx.SkillSubsetFromCLI = true
	if len(ctx.SkillSubset) != 2 {
		t.Errorf("expected 2 skill subsets, got %d", len(ctx.SkillSubset))
	}
	if !ctx.SkillSubsetFromCLI {
		t.Error("expected SkillSubsetFromCLI=true")
	}
}

func TestInstallContext_LocalContentTracking(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.OldLocalDeployed = []string{"old1.txt"}
	ctx.LocalDeployedFiles = []string{"new1.txt", "new2.txt"}
	ctx.LocalContentErrorsBefore = 3
	if len(ctx.OldLocalDeployed) != 1 {
		t.Errorf("OldLocalDeployed: expected 1, got %d", len(ctx.OldLocalDeployed))
	}
	if len(ctx.LocalDeployedFiles) != 2 {
		t.Errorf("LocalDeployedFiles: expected 2, got %d", len(ctx.LocalDeployedFiles))
	}
	if ctx.LocalContentErrorsBefore != 3 {
		t.Errorf("LocalContentErrorsBefore = %d, want 3", ctx.LocalContentErrorsBefore)
	}
}

func TestInstallContext_AllowInsecureHosts(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.AllowInsecureHosts = []string{"internal.corp.com", "dev.corp.com"}
	if len(ctx.AllowInsecureHosts) != 2 {
		t.Errorf("expected 2 insecure hosts, got %d", len(ctx.AllowInsecureHosts))
	}
}

func TestInstallContext_IntegrationCounters(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	ctx.TotalSubSkillsPromoted = 7
	ctx.TotalLinksResolved = 3
	if ctx.TotalSubSkillsPromoted != 7 {
		t.Errorf("TotalSubSkillsPromoted = %d", ctx.TotalSubSkillsPromoted)
	}
	if ctx.TotalLinksResolved != 3 {
		t.Errorf("TotalLinksResolved = %d", ctx.TotalLinksResolved)
	}
}

func TestInstallContext_ApmModulesDir_Default(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	got := ctx.ApmModulesDirOrDefault()
	if got == "" {
		t.Error("expected non-empty ApmModulesDir default")
	}
}

func TestInstallContext_LockfilePath_Default(t *testing.T) {
	ctx := installctx.New("/proj", "/proj/.apm")
	got := ctx.LockfilePathOrDefault()
	if got == "" {
		t.Error("expected non-empty LockfilePath default")
	}
}
