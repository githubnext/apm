package pack

import (
	"testing"
)

func TestFormatPlugin_Value_Extra2(t *testing.T) {
	if FormatPlugin != "plugin" {
		t.Errorf("FormatPlugin = %q, want plugin", FormatPlugin)
	}
}

func TestFormatAPM_Value_Extra2(t *testing.T) {
	if FormatAPM != "apm" {
		t.Errorf("FormatAPM = %q, want apm", FormatAPM)
	}
}

func TestFormat_Distinct_Extra2(t *testing.T) {
	if FormatPlugin == FormatAPM {
		t.Error("FormatPlugin and FormatAPM should be distinct")
	}
}

func TestPackOptions_ZeroValue_Extra2(t *testing.T) {
	var opts PackOptions
	if opts.ProjectRoot != "" || opts.DryRun || opts.Offline || opts.Archive {
		t.Error("zero-value PackOptions should have false/empty fields")
	}
}

func TestPackOptions_AllFields_Extra2(t *testing.T) {
	opts := PackOptions{
		ProjectRoot:  "/my/project",
		Format:       FormatPlugin,
		Archive:      true,
		OutputDir:    "/my/out",
		Offline:      true,
		DryRun:       true,
		Verbose:      true,
		MarketplaceOutput: "/my/mkt",
	}
	if opts.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Format != FormatPlugin {
		t.Errorf("Format = %q", opts.Format)
	}
	if !opts.Archive {
		t.Error("Archive should be true")
	}
	if !opts.DryRun {
		t.Error("DryRun should be true")
	}
	if opts.MarketplaceOutput != "/my/mkt" {
		t.Errorf("MarketplaceOutput = %q", opts.MarketplaceOutput)
	}
}

func TestPackResult_ZeroValue_Extra2(t *testing.T) {
	var r PackResult
	if len(r.OutputPaths) != 0 || r.DryRun {
		t.Error("zero-value PackResult should have empty fields")
	}
}

func TestPackResult_Fields_Extra2(t *testing.T) {
	r := PackResult{
		OutputPaths: []string{"/out/a", "/out/b"},
		DryRun:      true,
	}
	if len(r.OutputPaths) != 2 {
		t.Errorf("OutputPaths len = %d, want 2", len(r.OutputPaths))
	}
	if !r.DryRun {
		t.Error("DryRun should be true")
	}
}

func TestUnpackOptions_ZeroValue_Extra2(t *testing.T) {
	var opts UnpackOptions
	if opts.BundlePath != "" || opts.DestDir != "" || opts.DryRun {
		t.Error("zero-value UnpackOptions should have empty fields")
	}
}

func TestUnpackOptions_AllFields_Extra2(t *testing.T) {
	opts := UnpackOptions{
		BundlePath:  "/some/bundle.tar.gz",
		DestDir:     "/some/dest",
		ProjectRoot: "/some/project",
		DryRun:      true,
		Verbose:     true,
	}
	if opts.BundlePath != "/some/bundle.tar.gz" {
		t.Errorf("BundlePath = %q", opts.BundlePath)
	}
	if !opts.DryRun {
		t.Error("DryRun should be true")
	}
}

func TestRun_MissingProjectRoot_Extra2(t *testing.T) {
	_, err := Run(PackOptions{})
	if err == nil {
		t.Error("expected error when ProjectRoot is empty or invalid")
	}
}

func TestRunUnpack_MissingBundle_Extra2(t *testing.T) {
	err := RunUnpack(UnpackOptions{BundlePath: "/nonexistent/bundle.tar.gz", DestDir: t.TempDir()})
	if err == nil {
		t.Error("expected error when bundle file does not exist")
	}
}

func TestFormat_IsString_Extra2(t *testing.T) {
	f := FormatPlugin
	s := string(f)
	if s != "plugin" {
		t.Errorf("string(FormatPlugin) = %q", s)
	}
}
