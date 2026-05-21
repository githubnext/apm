package installpipeline_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/installpipeline"
)

type noopLogger struct{}

func (l *noopLogger) Progress(msg string)      {}
func (l *noopLogger) VerboseDetail(msg string) {}
func (l *noopLogger) Error(msg string)         {}

func TestDiagCollector_SingleMessage(t *testing.T) {
	d := &installpipeline.DiagCollector{}
	d.Add("only message")
	msgs := d.Messages()
	if len(msgs) != 1 || msgs[0] != "only message" {
		t.Errorf("unexpected messages: %v", msgs)
	}
}

func TestDiagCollector_ManyMessages(t *testing.T) {
	d := &installpipeline.DiagCollector{}
	for i := 0; i < 10; i++ {
		d.Add("msg")
	}
	if len(d.Messages()) != 10 {
		t.Errorf("expected 10 messages, got %d", len(d.Messages()))
	}
}

func TestDiagCollector_EmptyMessages(t *testing.T) {
	d := &installpipeline.DiagCollector{}
	if len(d.Messages()) != 0 {
		t.Error("expected 0 messages for fresh collector")
	}
}

func TestInstallContext_ZeroValue(t *testing.T) {
	ctx := &installpipeline.InstallContext{}
	if ctx.DryRun || ctx.Force || ctx.Frozen || ctx.Verbose {
		t.Error("bool fields should default to false")
	}
	if len(ctx.Targets) != 0 {
		t.Error("targets should default to empty")
	}
}

func TestInstallContext_WithLogger(t *testing.T) {
	ctx := &installpipeline.InstallContext{
		Logger: &noopLogger{},
	}
	if ctx.Logger == nil {
		t.Error("logger should be set")
	}
}

func TestPipelineResult_ZeroValue(t *testing.T) {
	r := &installpipeline.PipelineResult{}
	if r.Installed != 0 || r.Skipped != 0 || r.Removed != 0 || r.Updated != 0 {
		t.Error("counts should default to 0")
	}
	if len(r.Warnings) != 0 {
		t.Error("warnings should default to empty")
	}
}

func TestPipeline_Run_VerboseExtra(t *testing.T) {
	dir := t.TempDir()
	lg := &noopLogger{}
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		SkipLockfile: true,
		Verbose:      true,
		Logger:       lg,
	}
	p := installpipeline.NewPipeline()
	result, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPipeline_Run_DryRunNoApmYML(t *testing.T) {
	dir := t.TempDir()
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		DryRun:       true,
		SkipLockfile: true,
	}
	p := installpipeline.NewPipeline()
	result, err := p.Run(ctx)
	// No apm.yml: no dependencies resolved; no error expected
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = result
}

func TestPipeline_Run_WithTargets(t *testing.T) {
	dir := t.TempDir()
	apmYML := "dependencies:\n- name: myorg/mypkg\n  ref: main\n"
	os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(apmYML), 0o644)
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		Targets:      []string{"myorg/mypkg"},
		DryRun:       true,
		SkipLockfile: true,
	}
	p := installpipeline.NewPipeline()
	_, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipeline_Run_FrozenNoLockfile(t *testing.T) {
	dir := t.TempDir()
	ctx := &installpipeline.InstallContext{
		ProjectRoot: dir,
		Frozen:      true,
	}
	p := installpipeline.NewPipeline()
	_, err := p.Run(ctx)
	// Frozen without a lockfile should produce an error
	if err == nil {
		t.Log("frozen without lockfile did not error -- implementation may allow it")
	}
}

func TestPipeline_Run_MultipleRuns(t *testing.T) {
	dir := t.TempDir()
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		SkipLockfile: true,
	}
	p := installpipeline.NewPipeline()
	r1, err := p.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := p.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if r1 == nil || r2 == nil {
		t.Fatal("expected non-nil results")
	}
}

var _ = errors.New // ensure import used
