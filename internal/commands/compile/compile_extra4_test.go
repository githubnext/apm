package compile

import "testing"

func TestCompileOptions_AllFields_Extra4(t *testing.T) {
	opts := CompileOptions{
		ProjectRoot: "/project",
		Output:      "/out/AGENTS.md",
		DryRun:      true,
		Watch:       false,
		Force:       true,
		Strict:      true,
		Verbose:     false,
	}
	if opts.ProjectRoot != "/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Output != "/out/AGENTS.md" {
		t.Errorf("Output = %q", opts.Output)
	}
	if !opts.DryRun {
		t.Error("DryRun should be true")
	}
	if opts.Watch {
		t.Error("Watch should be false")
	}
	if !opts.Force {
		t.Error("Force should be true")
	}
	if !opts.Strict {
		t.Error("Strict should be true")
	}
}

func TestCompileOptions_ZeroValue_Extra4(t *testing.T) {
	var opts CompileOptions
	if opts.ProjectRoot != "" {
		t.Errorf("zero ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.DryRun {
		t.Error("zero DryRun should be false")
	}
}

func TestCompileStats_ZeroValue_Extra4(t *testing.T) {
	var s CompileStats
	if s.Instructions != 0 {
		t.Errorf("zero Instructions = %d", s.Instructions)
	}
	if s.Primitives != 0 {
		t.Errorf("zero Primitives = %d", s.Primitives)
	}
	if s.Warnings != nil {
		t.Error("zero Warnings should be nil")
	}
}

func TestCompileStats_Fields_Extra4(t *testing.T) {
	s := CompileStats{
		Instructions: 5,
		Contexts:     2,
		Chatmodes:    1,
		Primitives:   8,
		Warnings:     []string{"warn1"},
	}
	if s.Instructions != 5 {
		t.Errorf("Instructions = %d", s.Instructions)
	}
	if s.Primitives != 8 {
		t.Errorf("Primitives = %d", s.Primitives)
	}
	if len(s.Warnings) != 1 {
		t.Errorf("Warnings len = %d", len(s.Warnings))
	}
}

func TestCompileResult_ZeroValue_Extra4(t *testing.T) {
	var r CompileResult
	if r.OutputPath != "" {
		t.Errorf("zero OutputPath = %q", r.OutputPath)
	}
	if r.DryRun {
		t.Error("zero DryRun should be false")
	}
}

func TestCompileResult_Fields_Extra4(t *testing.T) {
	r := CompileResult{
		OutputPath:       "/out/AGENTS.md",
		ConstitutionHash: "abc123",
		Status:           "ok",
		DryRun:           true,
	}
	if r.OutputPath != "/out/AGENTS.md" {
		t.Errorf("OutputPath = %q", r.OutputPath)
	}
	if r.ConstitutionHash != "abc123" {
		t.Errorf("ConstitutionHash = %q", r.ConstitutionHash)
	}
	if r.Status != "ok" {
		t.Errorf("Status = %q", r.Status)
	}
}

func TestPrimitiveSection_ZeroValue_Extra4(t *testing.T) {
	var ps PrimitiveSection
	if ps.Title != "" {
		t.Errorf("zero Title = %q", ps.Title)
	}
}

func TestWatchOptions_Fields_Extra4(t *testing.T) {
	wo := WatchOptions{
		CompileOptions: CompileOptions{ProjectRoot: "/w", Watch: true},
		Interval:       500,
	}
	if wo.CompileOptions.ProjectRoot != "/w" {
		t.Errorf("ProjectRoot = %q", wo.CompileOptions.ProjectRoot)
	}
	if wo.Interval != 500 {
		t.Errorf("Interval = %d", wo.Interval)
	}
}
