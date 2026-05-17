package pack

import (
	"testing"
)

func TestPackResultFields(t *testing.T) {
	r := &PackResult{
		OutputPaths: []string{"/out/pkg.apm", "/out/pkg.tar.gz"},
		DryRun:      false,
	}
	if len(r.OutputPaths) != 2 {
		t.Errorf("OutputPaths len = %d, want 2", len(r.OutputPaths))
	}
	if r.DryRun {
		t.Error("DryRun should be false")
	}
}

func TestPackResult_DryRun(t *testing.T) {
	r := &PackResult{DryRun: true}
	if !r.DryRun {
		t.Error("DryRun should be true")
	}
	if len(r.OutputPaths) != 0 {
		t.Errorf("OutputPaths should be empty in dry-run, got %v", r.OutputPaths)
	}
}

func TestPackOptionsAllFields(t *testing.T) {
	opts := PackOptions{
		ProjectRoot:       "/proj",
		Format:            FormatAPM,
		Archive:           true,
		OutputDir:         "/out",
		Offline:           true,
		DryRun:            false,
		MarketplaceOutput: "/out/marketplace.json",
		Verbose:           true,
	}
	if opts.Format != FormatAPM {
		t.Errorf("Format = %q, want %q", opts.Format, FormatAPM)
	}
	if !opts.Archive {
		t.Error("Archive should be true")
	}
	if opts.OutputDir != "/out" {
		t.Errorf("OutputDir = %q, want /out", opts.OutputDir)
	}
	if !opts.Offline {
		t.Error("Offline should be true")
	}
}

func TestUnpackOptionsFields(t *testing.T) {
	opts := UnpackOptions{
		BundlePath: "/tmp/pkg.apm",
		DestDir:    "/tmp/dest",
		DryRun:     true,
	}
	if !opts.DryRun {
		t.Error("expected DryRun true")
	}
	if opts.BundlePath == "" {
		t.Error("expected non-empty BundlePath")
	}
}

func TestFormatConstants(t *testing.T) {
	tests := []struct {
		f    Format
		want Format
	}{
		{FormatPlugin, "plugin"},
		{FormatAPM, "apm"},
	}
	for _, tc := range tests {
		if tc.f != tc.want {
			t.Errorf("Format %q != %q", tc.f, tc.want)
		}
	}
}

func TestPackOptionsDefaults(t *testing.T) {
	opts := PackOptions{
		ProjectRoot: "/tmp/pkg",
		Format:      FormatPlugin,
		DryRun:      true,
	}
	if opts.ProjectRoot == "" {
		t.Error("expected non-empty ProjectRoot")
	}
	if opts.Format != FormatPlugin {
		t.Errorf("unexpected Format %q", opts.Format)
	}
	if !opts.DryRun {
		t.Error("expected DryRun true")
	}
}

func TestUnpackOptions(t *testing.T) {
	opts := UnpackOptions{
		BundlePath: "/tmp/pkg.apm",
		DestDir:    "/tmp/dest",
		DryRun:     true,
	}
	if !opts.DryRun {
		t.Error("expected DryRun true")
	}
	if opts.BundlePath == "" {
		t.Error("expected non-empty BundlePath")
	}
}

