package pack

import (
	"testing"
)

func TestFormat_String(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{FormatPlugin, "plugin"},
		{FormatAPM, "apm"},
		{Format("custom"), "custom"},
	}
	for _, tc := range tests {
		if string(tc.f) != tc.want {
			t.Errorf("Format(%q) string = %q, want %q", tc.f, string(tc.f), tc.want)
		}
	}
}

func TestPackOptions_ZeroValue(t *testing.T) {
	var opts PackOptions
	if opts.ProjectRoot != "" {
		t.Error("zero value ProjectRoot should be empty")
	}
	if opts.Format != "" {
		t.Error("zero value Format should be empty")
	}
	if opts.Archive {
		t.Error("zero value Archive should be false")
	}
	if opts.DryRun {
		t.Error("zero value DryRun should be false")
	}
}

func TestPackOptions_FullFields(t *testing.T) {
	opts := PackOptions{
		ProjectRoot:       "/some/path",
		Format:            FormatAPM,
		Archive:           true,
		OutputDir:         "/output",
		Offline:           true,
		DryRun:            true,
		MarketplaceOutput: "/out/marketplace.json",
		Verbose:           true,
	}
	if opts.ProjectRoot != "/some/path" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
	if opts.MarketplaceOutput != "/out/marketplace.json" {
		t.Errorf("MarketplaceOutput = %q", opts.MarketplaceOutput)
	}
}

func TestUnpackOptions_AllFields(t *testing.T) {
	opts := UnpackOptions{
		BundlePath: "/tmp/bundle.apm",
		DestDir:    "/tmp/dest",
		DryRun:     false,
	}
	if opts.BundlePath != "/tmp/bundle.apm" {
		t.Errorf("BundlePath = %q", opts.BundlePath)
	}
	if opts.DestDir != "/tmp/dest" {
		t.Errorf("DestDir = %q", opts.DestDir)
	}
	if opts.DryRun {
		t.Error("DryRun should be false")
	}
}

func TestPackResult_MultipleOutputs(t *testing.T) {
	r := &PackResult{
		OutputPaths: []string{"/a.apm", "/b.tar.gz", "/c.json"},
		DryRun:      false,
	}
	if len(r.OutputPaths) != 3 {
		t.Errorf("expected 3 outputs, got %d", len(r.OutputPaths))
	}
	if r.OutputPaths[2] != "/c.json" {
		t.Errorf("OutputPaths[2] = %q", r.OutputPaths[2])
	}
}

func TestPackResult_EmptyOutputs(t *testing.T) {
	r := &PackResult{
		OutputPaths: nil,
		DryRun:      true,
	}
	if len(r.OutputPaths) != 0 {
		t.Errorf("expected empty outputs, got %d", len(r.OutputPaths))
	}
}

func TestFormatPlugin_IsPlugin(t *testing.T) {
	if FormatPlugin != "plugin" {
		t.Errorf("FormatPlugin = %q, want plugin", FormatPlugin)
	}
}

func TestFormatAPM_IsAPM(t *testing.T) {
	if FormatAPM != "apm" {
		t.Errorf("FormatAPM = %q, want apm", FormatAPM)
	}
}

func TestFormatEquality(t *testing.T) {
	a := FormatPlugin
	b := Format("plugin")
	if a != b {
		t.Errorf("Format equality: %q != %q", a, b)
	}
}

func TestPackOptions_Offline(t *testing.T) {
	opts := PackOptions{
		ProjectRoot: "/p",
		Offline:     true,
	}
	if !opts.Offline {
		t.Error("Offline should be true")
	}
}

func TestPackOptions_OutputDir(t *testing.T) {
	opts := PackOptions{
		ProjectRoot: "/p",
		OutputDir:   "/custom/output",
	}
	if opts.OutputDir != "/custom/output" {
		t.Errorf("OutputDir = %q", opts.OutputDir)
	}
}

func TestUnpackOptions_ZeroValue(t *testing.T) {
	var opts UnpackOptions
	if opts.BundlePath != "" {
		t.Error("zero BundlePath should be empty")
	}
	if opts.DestDir != "" {
		t.Error("zero DestDir should be empty")
	}
	if opts.DryRun {
		t.Error("zero DryRun should be false")
	}
}
