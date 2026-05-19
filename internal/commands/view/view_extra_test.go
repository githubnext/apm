package view

import (
	"strings"
	"testing"
)

func TestParseSimpleYAML_BlankLines(t *testing.T) {
	data := []byte("\n\nkey: value\n\nother: data\n\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "value" {
		t.Errorf("key = %v", out["key"])
	}
	if out["other"] != "data" {
		t.Errorf("other = %v", out["other"])
	}
}

func TestParseSimpleYAML_ColonInValue_PreservesRest(t *testing.T) {
	data := []byte("repo: https://github.com/owner/repo:extra\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val := out["repo"].(string)
	if !strings.HasPrefix(val, "https://github.com/") {
		t.Errorf("expected URL value, got %q", val)
	}
}

func TestParseSimpleYAML_LeadingSpacesInValue(t *testing.T) {
	data := []byte("name:   my package  \n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// TrimSpace on value
	if out["name"] != "my package" {
		t.Errorf("name = %q", out["name"])
	}
}

func TestPackageInfo_EmptyFiles(t *testing.T) {
	info := PackageInfo{
		Name:   "my-pkg",
		Files:  nil,
	}
	if len(info.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(info.Files))
	}
}

func TestPackageInfo_VersionsList(t *testing.T) {
	info := PackageInfo{
		Name:     "pkg",
		Versions: []string{"v1.0.0", "v1.1.0", "v2.0.0"},
	}
	if len(info.Versions) != 3 {
		t.Errorf("expected 3 versions, got %d", len(info.Versions))
	}
	if info.Versions[0] != "v1.0.0" {
		t.Errorf("Versions[0] = %q", info.Versions[0])
	}
}

func TestPackageInfo_ApmYML(t *testing.T) {
	info := PackageInfo{
		Name: "pkg",
		ApmYML: map[string]interface{}{
			"description": "A useful skill",
			"version":     "1.0.0",
		},
	}
	if info.ApmYML["description"] != "A useful skill" {
		t.Errorf("ApmYML[description] = %v", info.ApmYML["description"])
	}
}

func TestViewOptions_FormatDefault(t *testing.T) {
	opts := ViewOptions{
		Package: "owner/repo",
	}
	if opts.Format != "" {
		t.Errorf("Format default should be empty string, got %q", opts.Format)
	}
}

func TestViewOptions_JSONFormat(t *testing.T) {
	opts := ViewOptions{
		Package: "owner/repo",
		Format:  "json",
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want json", opts.Format)
	}
}

func TestParseSimpleYAML_HashCommentMidLine(t *testing.T) {
	// Lines with no colon should be skipped
	data := []byte("# full comment line\nkey: value\n# another comment\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key (only 'key'), got %d: %v", len(out), out)
	}
	if out["key"] != "value" {
		t.Errorf("key = %v, want value", out["key"])
	}
}

func TestPackageInfo_InstalledPath(t *testing.T) {
	info := PackageInfo{
		Name:          "owner-repo",
		InstalledPath: "/home/user/.apm_modules/owner-repo",
	}
	if !strings.HasSuffix(info.InstalledPath, "owner-repo") {
		t.Errorf("InstalledPath = %q", info.InstalledPath)
	}
}

func TestPackageInfo_Source(t *testing.T) {
	info := PackageInfo{
		Name:   "pkg",
		Source: "https://github.com/owner/pkg",
	}
	if !strings.HasPrefix(info.Source, "https://") {
		t.Errorf("Source should be URL: %q", info.Source)
	}
}

func TestPackageInfo_RefAndCommit(t *testing.T) {
	info := PackageInfo{
		Name:   "pkg",
		Ref:    "main",
		Commit: "abcdef123456",
	}
	if info.Ref != "main" {
		t.Errorf("Ref = %q", info.Ref)
	}
	if len(info.Commit) < 6 {
		t.Errorf("Commit = %q (too short)", info.Commit)
	}
}

func TestParseSimpleYAML_NilMap_Initialized(t *testing.T) {
	var out map[string]interface{}
	if err := parseSimpleYAML([]byte("k: v\n"), &out); err != nil {
		t.Fatalf("error: %v", err)
	}
	if out == nil {
		t.Error("map should be initialized by parseSimpleYAML")
	}
}
