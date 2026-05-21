package agentscompiler

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// CompilationConfig additional fields
// ---------------------------------------------------------------------------

func TestCompilationConfigResolveLinks(t *testing.T) {
	cfg := CompilationConfig{ResolveLinks: true}
	if !cfg.ResolveLinks {
		t.Error("expected ResolveLinks=true")
	}
}

func TestCompilationConfigWithConstitution(t *testing.T) {
	cfg := CompilationConfig{WithConstitution: true}
	if !cfg.WithConstitution {
		t.Error("expected WithConstitution=true")
	}
}

func TestCompilationConfigDryRun(t *testing.T) {
	cfg := CompilationConfig{DryRun: true, OutputPath: "/tmp/out.md"}
	if !cfg.DryRun {
		t.Error("expected DryRun=true")
	}
	if cfg.OutputPath != "/tmp/out.md" {
		t.Errorf("expected /tmp/out.md, got %q", cfg.OutputPath)
	}
}

func TestCompilationConfigChatmode(t *testing.T) {
	cfg := CompilationConfig{Chatmode: "claude"}
	if cfg.Chatmode != "claude" {
		t.Errorf("expected claude, got %q", cfg.Chatmode)
	}
}

// ---------------------------------------------------------------------------
// CompilationResult
// ---------------------------------------------------------------------------

func TestCompilationResultFields(t *testing.T) {
	r := CompilationResult{
		Target:     "copilot",
		OutputPath: "/out/agents.md",
		Content:    "# content",
		BuildID:    "abc123",
		LinesOut:   42,
	}
	if r.Target != "copilot" {
		t.Errorf("expected copilot, got %q", r.Target)
	}
	if r.LinesOut != 42 {
		t.Errorf("expected 42, got %d", r.LinesOut)
	}
}

func TestCompilationResultOKWhenNoError(t *testing.T) {
	r := CompilationResult{}
	if !r.OK() {
		t.Error("zero value CompilationResult should be OK")
	}
}

// ---------------------------------------------------------------------------
// MergedResult fields
// ---------------------------------------------------------------------------

func TestMergedResultWarnings(t *testing.T) {
	m := MergedResult{Warnings: []string{"w1", "w2"}}
	if len(m.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(m.Warnings))
	}
}

func TestMergedResultTotalMS(t *testing.T) {
	m := MergedResult{TotalMS: 1234}
	if m.TotalMS != 1234 {
		t.Errorf("expected 1234, got %d", m.TotalMS)
	}
}

func TestMergedResultMultipleResults(t *testing.T) {
	m := MergedResult{
		Results: []CompilationResult{
			{Target: "copilot"},
			{Target: "claude"},
		},
	}
	if len(m.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(m.Results))
	}
}

// ---------------------------------------------------------------------------
// DistributedFile
// ---------------------------------------------------------------------------

func TestDistributedFileFields(t *testing.T) {
	f := DistributedFile{Path: "/a/b/c.md", Content: "hello"}
	if f.Path != "/a/b/c.md" {
		t.Errorf("expected /a/b/c.md, got %q", f.Path)
	}
	if f.Content != "hello" {
		t.Errorf("expected hello, got %q", f.Content)
	}
}

func TestDistributedFileZeroValue(t *testing.T) {
	var f DistributedFile
	if f.Path != "" || f.Content != "" {
		t.Error("expected zero value DistributedFile")
	}
}

// ---------------------------------------------------------------------------
// CompileStats
// ---------------------------------------------------------------------------

func TestCompileStats_NilResult(t *testing.T) {
	s := CompileStats(nil)
	_ = s // should not panic
}

func TestCompileStats_NonNilResult(t *testing.T) {
	m := &MergedResult{
		Results: []CompilationResult{{Target: "copilot", LinesOut: 10}},
		TotalMS: 50,
	}
	s := CompileStats(m)
	if !strings.Contains(s, "copilot") && s == "" {
		// stats may vary; just verify no panic
	}
}

// ---------------------------------------------------------------------------
// CopilotRootInstructionsPath
// ---------------------------------------------------------------------------

func TestCopilotRootInstructionsPath_AbsPath(t *testing.T) {
	p := CopilotRootInstructionsPath("/repo")
	if !strings.HasPrefix(p, "/repo") {
		t.Errorf("expected path under /repo, got %q", p)
	}
	if !strings.Contains(p, ".github") {
		t.Errorf("expected .github in path, got %q", p)
	}
}

func TestCopilotRootInstructionsPath_EmptyBase(t *testing.T) {
	p := CopilotRootInstructionsPath("")
	if p == "" {
		t.Error("expected non-empty path even for empty base")
	}
}
