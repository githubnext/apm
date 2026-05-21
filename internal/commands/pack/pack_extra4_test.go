package pack

import "testing"

func TestPackOptions_AllFields_Extra4(t *testing.T) {
	opts := PackOptions{
		ProjectRoot:       "/project",
		Format:            FormatAPM,
		Archive:           true,
		OutputDir:         "/out",
		Offline:           true,
		DryRun:            true,
		MarketplaceOutput: "/mkt.json",
		Verbose:           true,
	}
	if opts.ProjectRoot != "/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Format != FormatAPM {
		t.Errorf("Format = %q", opts.Format)
	}
	if !opts.Archive {
		t.Error("Archive should be true")
	}
	if opts.OutputDir != "/out" {
		t.Errorf("OutputDir = %q", opts.OutputDir)
	}
	if !opts.Offline {
		t.Error("Offline should be true")
	}
	if !opts.DryRun {
		t.Error("DryRun should be true")
	}
	if opts.MarketplaceOutput != "/mkt.json" {
		t.Errorf("MarketplaceOutput = %q", opts.MarketplaceOutput)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestPackOptions_FormatPlugin_Extra4(t *testing.T) {
	opts := PackOptions{Format: FormatPlugin}
	if opts.Format != FormatPlugin {
		t.Errorf("Format = %q, want %q", opts.Format, FormatPlugin)
	}
}

func TestPackOptions_ZeroValue_Extra4(t *testing.T) {
	var opts PackOptions
	if opts.Format != "" {
		t.Errorf("zero Format should be empty, got %q", opts.Format)
	}
	if opts.Archive {
		t.Error("zero Archive should be false")
	}
	if opts.DryRun {
		t.Error("zero DryRun should be false")
	}
}

func TestPackResult_ZeroValue_Extra4(t *testing.T) {
	var r PackResult
	if r.DryRun {
		t.Error("zero DryRun should be false")
	}
	if r.OutputPaths != nil {
		t.Errorf("zero OutputPaths should be nil, got %v", r.OutputPaths)
	}
}

func TestPackResult_SingleOutput_Extra4(t *testing.T) {
	r := PackResult{
		OutputPaths: []string{"out.tar.gz"},
		DryRun:      false,
	}
	if len(r.OutputPaths) != 1 {
		t.Errorf("len(OutputPaths) = %d, want 1", len(r.OutputPaths))
	}
	if r.OutputPaths[0] != "out.tar.gz" {
		t.Errorf("OutputPaths[0] = %q, want out.tar.gz", r.OutputPaths[0])
	}
}

func TestPackResult_DryRun_Extra4(t *testing.T) {
	r := PackResult{DryRun: true, OutputPaths: []string{}}
	if !r.DryRun {
		t.Error("DryRun should be true")
	}
	if len(r.OutputPaths) != 0 {
		t.Errorf("OutputPaths should be empty")
	}
}

func TestFormatConstants_Extra4(t *testing.T) {
	if FormatPlugin == FormatAPM {
		t.Error("FormatPlugin and FormatAPM should differ")
	}
	if string(FormatPlugin) == "" {
		t.Error("FormatPlugin should be non-empty")
	}
	if string(FormatAPM) == "" {
		t.Error("FormatAPM should be non-empty")
	}
}

func TestPackOptions_OnlyProjectRoot_Extra4(t *testing.T) {
	opts := PackOptions{ProjectRoot: "/tmp/project"}
	if opts.ProjectRoot != "/tmp/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.DryRun {
		t.Error("DryRun defaults false")
	}
}

func TestPackOptions_OutputDirOnly_Extra4(t *testing.T) {
	opts := PackOptions{OutputDir: "/dist"}
	if opts.OutputDir != "/dist" {
		t.Errorf("OutputDir = %q", opts.OutputDir)
	}
}
