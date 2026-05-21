package view

import (
	"testing"
)

func TestViewOptions_ZeroValue_Extra2(t *testing.T) {
	var opts ViewOptions
	if opts.ProjectRoot != "" || opts.Package != "" || opts.Verbose {
		t.Error("zero-value ViewOptions should have empty/false fields")
	}
}

func TestViewOptions_AllFields_Extra2(t *testing.T) {
	opts := ViewOptions{
		ProjectRoot: "/my/project",
		Package:     "mypkg",
		Field:       "versions",
		Format:      "json",
		Verbose:     true,
	}
	if opts.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Package != "mypkg" {
		t.Errorf("Package = %q", opts.Package)
	}
	if opts.Field != "versions" {
		t.Errorf("Field = %q", opts.Field)
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q", opts.Format)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestPackageInfo_ZeroValue_Extra2(t *testing.T) {
	var info PackageInfo
	if info.Name != "" || info.InstalledPath != "" {
		t.Error("zero-value PackageInfo should have empty fields")
	}
}

func TestPackageInfo_Fields_Extra2(t *testing.T) {
	info := PackageInfo{
		Name:          "mypkg",
		InstalledPath: "/path/to/pkg",
		Ref:           "v1.0.0",
		Commit:        "abc123",
		Source:        "github",
	}
	if info.Name != "mypkg" {
		t.Errorf("Name = %q", info.Name)
	}
	if info.InstalledPath != "/path/to/pkg" {
		t.Errorf("InstalledPath = %q", info.InstalledPath)
	}
	if info.Ref != "v1.0.0" {
		t.Errorf("Ref = %q", info.Ref)
	}
}

func TestParseSimpleYAML_KeyValue_Extra2(t *testing.T) {
	data := []byte("key: value\nother: 123\n")
	out := make(map[string]interface{})
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	if out["key"] != "value" {
		t.Errorf("key = %v, want value", out["key"])
	}
	if out["other"] != "123" {
		t.Errorf("other = %v, want 123", out["other"])
	}
}

func TestParseSimpleYAML_EmptyInput_Extra2(t *testing.T) {
	out := make(map[string]interface{})
	if err := parseSimpleYAML([]byte(""), &out); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestParseSimpleYAML_CommentLine_Extra2(t *testing.T) {
	data := []byte("# comment\nname: foo\n")
	out := make(map[string]interface{})
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	if out["name"] != "foo" {
		t.Errorf("name = %v, want foo", out["name"])
	}
}

func TestParseSimpleYAML_ColonInValue_Extra2(t *testing.T) {
	data := []byte("url: http://example.com/path\n")
	out := make(map[string]interface{})
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	v, ok := out["url"]
	if !ok {
		t.Error("expected key 'url'")
	}
	s, _ := v.(string)
	if s != "http://example.com/path" {
		t.Errorf("url = %q, want http://example.com/path", s)
	}
}

func TestRun_MissingProjectRoot_Extra2(t *testing.T) {
	err := Run(ViewOptions{Package: "nonexistent-pkg"})
	if err == nil {
		t.Error("expected error when ProjectRoot/Package is invalid")
	}
}

func TestPackageInfo_Versions_Extra2(t *testing.T) {
	info := PackageInfo{
		Name:     "mypkg",
		Versions: []string{"1.0.0", "1.1.0", "2.0.0"},
	}
	if len(info.Versions) != 3 {
		t.Errorf("Versions len = %d, want 3", len(info.Versions))
	}
}

func TestPackageInfo_Files_Extra2(t *testing.T) {
	info := PackageInfo{
		Name:  "mypkg",
		Files: []string{"README.md", "main.py"},
	}
	if len(info.Files) != 2 {
		t.Errorf("Files len = %d, want 2", len(info.Files))
	}
}
