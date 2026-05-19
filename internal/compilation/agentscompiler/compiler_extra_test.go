package agentscompiler

import (
	"strings"
	"testing"
)

func TestCompilationConfigDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DryRun {
		t.Error("DryRun should be false by default")
	}
	if cfg.Verbose {
		t.Error("Verbose should be false by default")
	}
	if cfg.Quiet {
		t.Error("Quiet should be false by default")
	}
	if cfg.Chatmode != "" {
		t.Errorf("Chatmode should be empty by default, got %q", cfg.Chatmode)
	}
}

func TestTargetConstants(t *testing.T) {
	targets := []CompileTargetType{
		TargetAll, TargetVSCode, TargetAgents, TargetCopilot,
		TargetClaude, TargetGemini, TargetCursor, TargetOpenCode,
		TargetCodex, TargetWindsurf, TargetMinimal, TargetAgentSkills,
	}
	for _, tgt := range targets {
		if tgt == "" {
			t.Error("target constant should not be empty")
		}
	}
}

func TestStrategyConstants(t *testing.T) {
	if StrategyDistributed == StrategySingleFile {
		t.Error("strategies should be distinct")
	}
	if StrategyDistributed == "" || StrategySingleFile == "" {
		t.Error("strategy constants should not be empty")
	}
}

func TestBuildIDPlaceholder(t *testing.T) {
	if BuildIDPlaceholder == "" {
		t.Error("BuildIDPlaceholder should not be empty")
	}
	if !strings.Contains(BuildIDPlaceholder, "APM_BUILD_ID") {
		t.Errorf("BuildIDPlaceholder %q does not contain APM_BUILD_ID", BuildIDPlaceholder)
	}
}

func TestCopilotRootGeneratedMarker(t *testing.T) {
	if CopilotRootGeneratedMarker == "" {
		t.Error("CopilotRootGeneratedMarker should not be empty")
	}
	if !strings.Contains(CopilotRootGeneratedMarker, "APM") {
		t.Errorf("marker %q should reference APM", CopilotRootGeneratedMarker)
	}
}

func TestCompilationResultOKExtra(t *testing.T) {
	ok := CompilationResult{}
	if !ok.OK() {
		t.Error("zero CompilationResult should have OK=true (no error)")
	}
}

func TestMergedResultOKNoResults(t *testing.T) {
	m := MergedResult{}
	if !m.OK() {
		t.Error("empty MergedResult should be OK")
	}
}

func TestMergedResultNotOKWithError(t *testing.T) {
	import_err := &MergedResult{
		Results: []CompilationResult{
			{Error: nil},
			{Error: &testError{"some error"}},
		},
	}
	if import_err.OK() {
		t.Error("MergedResult with an error result should not be OK")
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestAgentsCompilerNew(t *testing.T) {
	a := New("")
	if a == nil {
		t.Fatal("New('') returned nil")
	}
	// should default to abs path
	if a.baseDir == "" {
		t.Error("baseDir should not be empty")
	}
}

func TestAgentsCompilerNewWithDir(t *testing.T) {
	tmp := t.TempDir()
	a := New(tmp)
	if a.baseDir != tmp {
		t.Errorf("baseDir: got %q, want %q", a.baseDir, tmp)
	}
}

func TestCopilotRootInstructionsPathExtra(t *testing.T) {
	p := CopilotRootInstructionsPath("/some/dir")
	if p == "" {
		t.Error("CopilotRootInstructionsPath returned empty string")
	}
	if !strings.Contains(p, "some") {
		t.Errorf("expected path to contain base dir, got %q", p)
	}
}
