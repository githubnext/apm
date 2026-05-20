package pluginexporter

import (
	"testing"
)

func TestPackResult_ZeroValue(t *testing.T) {
	var r PackResult
	if r.BundlePath != "" {
		t.Errorf("expected empty BundlePath, got %q", r.BundlePath)
	}
	if r.MappedCount != 0 {
		t.Errorf("expected MappedCount 0, got %d", r.MappedCount)
	}
	if r.LockfileEnriched {
		t.Error("expected LockfileEnriched false")
	}
}

func TestPackResult_FieldAssignment(t *testing.T) {
	r := PackResult{
		BundlePath:      "/out/bundle",
		MappedCount:     5,
		LockfileEnriched: true,
		Files:           []string{"a.md", "b.md"},
	}
	if r.BundlePath != "/out/bundle" {
		t.Errorf("unexpected BundlePath %q", r.BundlePath)
	}
	if r.MappedCount != 5 {
		t.Errorf("unexpected MappedCount %d", r.MappedCount)
	}
	if !r.LockfileEnriched {
		t.Error("expected LockfileEnriched true")
	}
	if len(r.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(r.Files))
	}
}

func TestExportOptions_ZeroValue(t *testing.T) {
	var o ExportOptions
	if o.ProjectRoot != "" {
		t.Errorf("expected empty ProjectRoot")
	}
	if o.DryRun {
		t.Error("expected DryRun false")
	}
	if o.Force {
		t.Error("expected Force false")
	}
}

func TestExportOptions_DryRun(t *testing.T) {
	o := ExportOptions{DryRun: true, ProjectRoot: "/p", OutputDir: "/o"}
	if !o.DryRun {
		t.Error("expected DryRun true")
	}
}

func TestPluginJSON_ZeroValue(t *testing.T) {
	var p PluginJSON
	if p.Name != "" {
		t.Errorf("expected empty Name")
	}
	if p.Version != "" {
		t.Errorf("expected empty Version")
	}
}

func TestPluginJSON_FieldRoundtrip(t *testing.T) {
	p := PluginJSON{Name: "myplugin", Version: "1.0.0"}
	if p.Name != "myplugin" {
		t.Errorf("unexpected Name %q", p.Name)
	}
	if p.Version != "1.0.0" {
		t.Errorf("unexpected Version %q", p.Version)
	}
}

func TestSanitizeBundleName_LeadingTrailingHyphens(t *testing.T) {
	result := sanitizeBundleName("!hello!")
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
	if result[0] == '-' || result[len(result)-1] == '-' {
		t.Errorf("result should not have leading/trailing hyphens: %q", result)
	}
}

func TestSanitizeBundleName_Empty_ReturnsUnnamed(t *testing.T) {
	result := sanitizeBundleName("!!!")
	if result != "unnamed" {
		t.Errorf("expected unnamed, got %q", result)
	}
}

func TestValidateOutputRel_DotDot(t *testing.T) {
	if validateOutputRel("../escape") {
		t.Error("expected dotdot to be rejected")
	}
}

func TestValidateOutputRel_AbsoluteUnix(t *testing.T) {
	if validateOutputRel("/absolute/path") {
		t.Error("expected absolute path to be rejected")
	}
}

func TestValidateOutputRel_Normal(t *testing.T) {
	if !validateOutputRel("subdir/file.md") {
		t.Error("expected normal relative path to be valid")
	}
}

func TestExportPluginBundle_NonExistentOutputDir(t *testing.T) {
	_, err := ExportPluginBundle(ExportOptions{
		ProjectRoot: "/nonexistent/project/root/xyz",
		OutputDir:   "/nonexistent/output/dir/xyz",
	})
	if err == nil {
		t.Error("expected error for non-existent directories")
	}
}
