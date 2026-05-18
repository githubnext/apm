package agentscompiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.OutputPath != "AGENTS.md" {
		t.Errorf("OutputPath: got %q, want AGENTS.md", cfg.OutputPath)
	}
	if !cfg.ResolveLinks {
		t.Error("ResolveLinks should be true")
	}
	if !cfg.WithConstitution {
		t.Error("WithConstitution should be true")
	}
	if cfg.Target != TargetAll {
		t.Errorf("Target: got %q, want %q", cfg.Target, TargetAll)
	}
	if cfg.Strategy != StrategyDistributed {
		t.Errorf("Strategy: got %q, want %q", cfg.Strategy, StrategyDistributed)
	}
}

func TestNew(t *testing.T) {
	a := New("")
	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.baseDir == "" {
		t.Error("baseDir should not be empty")
	}

	tmp := t.TempDir()
	a2 := New(tmp)
	if a2.baseDir != tmp {
		t.Errorf("baseDir: got %q, want %q", a2.baseDir, tmp)
	}
}

func TestResolveTargets(t *testing.T) {
	a := New(".")
	cases := []struct {
		in   CompileTargetType
		want []CompileTargetType
	}{
		{TargetAll, []CompileTargetType{TargetVSCode, TargetClaude, TargetGemini}},
		{TargetAgents, []CompileTargetType{TargetVSCode}},
		{TargetCopilot, []CompileTargetType{TargetVSCode}},
		{TargetClaude, []CompileTargetType{TargetClaude}},
		{TargetGemini, []CompileTargetType{TargetGemini}},
		{TargetCursor, []CompileTargetType{TargetCursor}},
	}
	for _, tc := range cases {
		got := a.resolveTargets(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("resolveTargets(%q): got %v, want %v", tc.in, got, tc.want)
			continue
		}
		for i, g := range got {
			if g != tc.want[i] {
				t.Errorf("resolveTargets(%q)[%d]: got %q, want %q", tc.in, i, g, tc.want[i])
			}
		}
	}
}

func TestCompilationResultOK(t *testing.T) {
	r := CompilationResult{}
	if !r.OK() {
		t.Error("empty result should be OK")
	}
	r.Error = os.ErrNotExist
	if r.OK() {
		t.Error("result with error should not be OK")
	}
}

func TestMergedResultOK(t *testing.T) {
	m := MergedResult{}
	if !m.OK() {
		t.Error("empty MergedResult should be OK")
	}
	m.Results = []CompilationResult{{Error: nil}, {Error: nil}}
	if !m.OK() {
		t.Error("results with no errors should be OK")
	}
	m.Results = append(m.Results, CompilationResult{Error: os.ErrPermission})
	if m.OK() {
		t.Error("results with an error should not be OK")
	}
}

func TestCompileStats(t *testing.T) {
	if s := CompileStats(nil); s != "no results" {
		t.Errorf("nil: got %q", s)
	}
	m := &MergedResult{
		Results: []CompilationResult{
			{Target: TargetClaude, Error: nil, ElapsedMS: 10},
			{Target: TargetGemini, Error: os.ErrNotExist, ElapsedMS: 5},
		},
	}
	s := CompileStats(m)
	if !strings.Contains(s, TargetClaude) || !strings.Contains(s, TargetGemini) {
		t.Errorf("stats missing targets: %q", s)
	}
	if !strings.Contains(s, "error") {
		t.Errorf("stats should include error: %q", s)
	}
}

func TestFinalizeBuildID(t *testing.T) {
	a := New(".")
	content := "header\n" + BuildIDPlaceholder + "\nfooter"
	out := a.finalizeBuildID(content)
	if strings.Contains(out, BuildIDPlaceholder) {
		t.Error("placeholder should have been replaced")
	}
	if !strings.Contains(out, "<!-- Build ID:") {
		t.Error("output should contain Build ID comment")
	}
	// Deterministic
	out2 := a.finalizeBuildID(content)
	if out != out2 {
		t.Error("finalizeBuildID should be deterministic")
	}
}

func TestExtractBuildID(t *testing.T) {
	a := New(".")
	content := "<!-- Build ID: abc123 -->\nsome content"
	id := a.extractBuildID(content)
	if id != "<!-- Build ID: abc123 -->" {
		t.Errorf("got %q", id)
	}
	if a.extractBuildID("no build id here") != "" {
		t.Error("should return empty string when no build id")
	}
}

func TestCompileWithEmptyDir(t *testing.T) {
	tmp := t.TempDir()
	a := New(tmp)
	cfg := CompilationConfig{
		OutputPath: "out.md",
		Target:     TargetClaude,
		DryRun:     true,
	}
	result, err := a.Compile(cfg)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if len(result.Results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestValidatePrimitivesNoDir(t *testing.T) {
	tmp := t.TempDir()
	a := New(tmp)
	errs := a.ValidatePrimitives()
	if len(errs) == 0 {
		t.Error("expected error for missing .apm dir")
	}
}

func TestValidatePrimitivesWithDir(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, ".apm"), 0o755); err != nil {
		t.Fatal(err)
	}
	a := New(tmp)
	errs := a.ValidatePrimitives()
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestWriteDistributedFiles(t *testing.T) {
	tmp := t.TempDir()
	files := []DistributedFile{
		{Path: filepath.Join(tmp, "a.md"), Content: "hello"},
		{Path: filepath.Join(tmp, "sub", "b.md"), Content: "world"},
	}
	if err := WriteDistributedFiles(files, false); err != nil {
		t.Fatalf("WriteDistributedFiles: %v", err)
	}
	for _, f := range files {
		data, err := os.ReadFile(f.Path)
		if err != nil {
			t.Errorf("read %q: %v", f.Path, err)
		}
		if string(data) != f.Content {
			t.Errorf("%q: got %q, want %q", f.Path, data, f.Content)
		}
	}
}

func TestWriteDistributedFilesDryRun(t *testing.T) {
	tmp := t.TempDir()
	files := []DistributedFile{
		{Path: filepath.Join(tmp, "dry.md"), Content: "should not be written"},
	}
	if err := WriteDistributedFiles(files, true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(files[0].Path); !os.IsNotExist(err) {
		t.Error("dry run should not create files")
	}
}

func TestCopilotRootInstructionsPath(t *testing.T) {
	path := CopilotRootInstructionsPath("/repo")
	if !strings.HasSuffix(path, ".github/copilot-instructions.md") {
		t.Errorf("unexpected path: %q", path)
	}
}

func TestCleanupCopilotRootInstructions(t *testing.T) {
	tmp := t.TempDir()
	// No file -- should not error
	if err := CleanupCopilotRootInstructions(tmp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// With generated marker
	p := CopilotRootInstructionsPath(tmp)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	content := CopilotRootGeneratedMarker + "\nsome content"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CleanupCopilotRootInstructions(tmp); err != nil {
		t.Errorf("cleanup error: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("file should have been deleted")
	}
}

func TestCleanupCopilotRootInstructionsNoMarker(t *testing.T) {
	tmp := t.TempDir()
	p := CopilotRootInstructionsPath(tmp)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "user-written file without marker"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CleanupCopilotRootInstructions(tmp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// File should still exist (no marker)
	if _, err := os.Stat(p); err != nil {
		t.Error("file without marker should not be deleted")
	}
}

func TestCompileAgentsMDConvenienceFunc(t *testing.T) {
	tmp := t.TempDir()
	cfg := CompilationConfig{
		OutputPath: "AGENTS.md",
		Target:     TargetClaude,
		DryRun:     true,
	}
	result, err := CompileAgentsMD(tmp, cfg)
	if err != nil {
		t.Fatalf("CompileAgentsMD: %v", err)
	}
	if result == nil {
		t.Fatal("nil result")
	}
}
