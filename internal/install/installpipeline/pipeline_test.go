package installpipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/installpipeline"
)

// captureLogger records messages.
type captureLogger struct {
	progress []string
	verbose  []string
	errors   []string
}

func (l *captureLogger) Progress(msg string)      { l.progress = append(l.progress, msg) }
func (l *captureLogger) VerboseDetail(msg string) { l.verbose = append(l.verbose, msg) }
func (l *captureLogger) Error(msg string)         { l.errors = append(l.errors, msg) }

func TestDiagCollector(t *testing.T) {
	d := &installpipeline.DiagCollector{}
	if len(d.Messages()) != 0 {
		t.Error("expected empty messages")
	}
	d.Add("msg1")
	d.Add("msg2")
	msgs := d.Messages()
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0] != "msg1" || msgs[1] != "msg2" {
		t.Errorf("unexpected messages: %v", msgs)
	}
	// Verify isolation (returned slice is a copy).
	msgs[0] = "mutated"
	if d.Messages()[0] != "msg1" {
		t.Error("Messages() should return an independent copy")
	}
}

func TestPipeline_Run_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		SkipLockfile: true,
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

func TestPipeline_Run_WithApmYML(t *testing.T) {
	dir := t.TempDir()
	// Use one-space indent so the dash is at column 1 after trimLeft
	apmYML := "dependencies:\n- name: pkg-one\n  ref: main\n"
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(apmYML), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		DryRun:       true,
		SkipLockfile: true,
		Logger:       &captureLogger{},
	}
	p := installpipeline.NewPipeline()
	result, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Installed+result.Skipped < 1 {
		t.Errorf("expected >=1 total (dry-run), installed=%d skipped=%d", result.Installed, result.Skipped)
	}
}

func TestPipeline_Run_Verbose(t *testing.T) {
	dir := t.TempDir()
	log := &captureLogger{}
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		SkipLockfile: true,
		Verbose:      true,
		Logger:       log,
	}
	p := installpipeline.NewPipeline()
	_, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.verbose) == 0 {
		t.Error("expected verbose timing messages")
	}
}

func TestPipeline_Run_Frozen_NoLockfile(t *testing.T) {
	dir := t.TempDir()
	ctx := &installpipeline.InstallContext{
		ProjectRoot: dir,
		Frozen:      true,
	}
	p := installpipeline.NewPipeline()
	_, err := p.Run(ctx)
	if err == nil {
		t.Error("expected error: frozen without lockfile")
	}
}

func TestPipeline_AddCustomPhase(t *testing.T) {
	called := false
	type customPhase struct{}
	_ = called

	dir := t.TempDir()
	ctx := &installpipeline.InstallContext{
		ProjectRoot:  dir,
		SkipLockfile: true,
	}
	p := installpipeline.NewPipeline()
	p.AddPhase(&testPhase{name: "custom", fn: func(c *installpipeline.InstallContext) error {
		called = true
		return nil
	}})
	_, err := p.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("custom phase was not called")
	}
}

// testPhase is a minimal Phase implementation for testing.
type testPhase struct {
	name string
	fn   func(*installpipeline.InstallContext) error
}

func (tp *testPhase) Name() string { return tp.name }
func (tp *testPhase) Run(ctx *installpipeline.InstallContext) error {
	if tp.fn != nil {
		return tp.fn(ctx)
	}
	return nil
}

func TestPipeline_Run_WithLockfileWrite(t *testing.T) {
	dir := t.TempDir()
	apmYML := "dependencies:\n  - name: mypkg\n    ref: v1.2.3\n"
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(apmYML), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := &installpipeline.InstallContext{
		ProjectRoot: dir,
		DryRun:      false,
	}
	p := installpipeline.NewPipeline()
	_, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lockPath := filepath.Join(dir, "apm.lock.yaml")
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("expected apm.lock.yaml to be written")
	}
}
