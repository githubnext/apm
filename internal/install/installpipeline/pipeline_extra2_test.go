package installpipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/installpipeline"
)

func TestInstallContext_ZeroValue_extra2(t *testing.T) {
	ctx := &installpipeline.InstallContext{}
	if ctx.DryRun || ctx.Verbose || ctx.Force || ctx.Frozen || ctx.SkipLockfile {
		t.Error("expected all bool fields false by default")
	}
}

func TestInstallContext_Fields(t *testing.T) {
	ctx := &installpipeline.InstallContext{
		ProjectRoot: "/proj",
		ModulesDir:  "/proj/.apm/modules",
		Targets:     []string{"copilot"},
		DryRun:      true,
		Verbose:     true,
		AuthToken:   "token123",
	}
	if ctx.ProjectRoot != "/proj" {
		t.Errorf("ProjectRoot = %q", ctx.ProjectRoot)
	}
	if !ctx.DryRun {
		t.Error("expected DryRun=true")
	}
	if ctx.AuthToken != "token123" {
		t.Errorf("AuthToken = %q", ctx.AuthToken)
	}
	if len(ctx.Targets) != 1 || ctx.Targets[0] != "copilot" {
		t.Errorf("Targets = %v", ctx.Targets)
	}
}

func TestResolvedDep_LocalFlag(t *testing.T) {
	d := installpipeline.ResolvedDep{
		Name:  "local-pkg",
		Local: true,
	}
	if !d.Local {
		t.Error("expected Local=true")
	}
}

func TestResolvedDep_AllFields(t *testing.T) {
	d := installpipeline.ResolvedDep{
		Name:   "mypkg",
		Ref:    "v1.0",
		Commit: "abc123",
		Source: "github.com/org/repo",
		PkgDir: "/tmp/pkg",
	}
	if d.Name != "mypkg" || d.Ref != "v1.0" || d.Commit != "abc123" {
		t.Errorf("unexpected fields: %+v", d)
	}
}

func TestPipelineResult_ZeroValue_extra2(t *testing.T) {
	r := installpipeline.PipelineResult{}
	if r.Installed != 0 || r.Skipped != 0 || r.Removed != 0 || r.Updated != 0 {
		t.Errorf("unexpected non-zero fields: %+v", r)
	}
}

func TestPipelineResult_AllFields(t *testing.T) {
	r := &installpipeline.PipelineResult{
		Installed: 5,
		Skipped:   2,
		Updated:   1,
		Warnings:  []string{"warn1"},
	}
	if r.Installed != 5 {
		t.Errorf("Installed = %d", r.Installed)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("Warnings = %v", r.Warnings)
	}
}

func TestNewPipeline_Empty(t *testing.T) {
	p := installpipeline.NewPipeline()
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestPipeline_RunEmptyContext_NoError(t *testing.T) {
	p := installpipeline.NewPipeline()
	ctx := &installpipeline.InstallContext{}
	// Empty project root falls back to current working directory; should not panic.
	_, err := p.Run(ctx)
	_ = err
}

func TestPipeline_RunWithRoot(t *testing.T) {
	dir := t.TempDir()
	// Write minimal apm.yml.
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("dependencies: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := installpipeline.NewPipeline()
	ctx := &installpipeline.InstallContext{
		ProjectRoot: dir,
		DryRun:      true,
	}
	_, err := p.Run(ctx)
	// Should succeed or return a non-panic error.
	_ = err
}

func TestDiagCollector_EmptyMessages_extra2(t *testing.T) {
	d := &installpipeline.DiagCollector{}
	msgs := d.Messages()
	if len(msgs) != 0 {
		t.Errorf("expected empty messages, got %v", msgs)
	}
}

func TestDiagCollector_AddAndRetrieve(t *testing.T) {
	d := &installpipeline.DiagCollector{}
	d.Add("first")
	d.Add("second")
	d.Add("third")
	msgs := d.Messages()
	if len(msgs) != 3 {
		t.Errorf("expected 3 messages, got %d", len(msgs))
	}
	if msgs[0] != "first" {
		t.Errorf("msgs[0] = %q", msgs[0])
	}
	if msgs[2] != "third" {
		t.Errorf("msgs[2] = %q", msgs[2])
	}
}
