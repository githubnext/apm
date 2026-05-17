package pack

import "testing"

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

