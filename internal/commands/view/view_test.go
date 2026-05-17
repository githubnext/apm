package view

import "testing"

func TestParseSimpleYAML(t *testing.T) {
	data := []byte("name: mypkg\nversion: v1.0.0\n# comment\n\ndescription: A test package\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["name"] != "mypkg" {
		t.Errorf("name = %v, want %q", out["name"], "mypkg")
	}
	if out["version"] != "v1.0.0" {
		t.Errorf("version = %v, want %q", out["version"], "v1.0.0")
	}
	if out["description"] != "A test package" {
		t.Errorf("description = %v, want %q", out["description"], "A test package")
	}
}

func TestParseSimpleYAMLEmpty(t *testing.T) {
	var out map[string]interface{}
	if err := parseSimpleYAML([]byte(""), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestParseSimpleYAMLNoColon(t *testing.T) {
	data := []byte("justtext\nkey: value\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["justtext"]; ok {
		t.Error("should not have parsed 'justtext' as a key")
	}
	if out["key"] != "value" {
		t.Errorf("key = %v, want %q", out["key"], "value")
	}
}

func TestViewOptions(t *testing.T) {
	opts := ViewOptions{
		Package: "mypkg",
		Format:  "text",
	}
	if opts.Package != "mypkg" {
		t.Errorf("unexpected Package %q", opts.Package)
	}
}

func TestViewOptionsAllFields(t *testing.T) {
	opts := ViewOptions{
		ProjectRoot: "/home/user/project",
		Package:     "owner/repo",
		Field:       "versions",
		Format:      "json",
		Verbose:     true,
	}
	if opts.ProjectRoot != "/home/user/project" {
		t.Errorf("ProjectRoot mismatch: %q", opts.ProjectRoot)
	}
	if opts.Field != "versions" {
		t.Errorf("Field mismatch: %q", opts.Field)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestParseSimpleYAMLMultipleValues(t *testing.T) {
	data := []byte("key1: val1\nkey2: val2\nkey3: val3\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
	if out["key3"] != "val3" {
		t.Errorf("key3 = %v, want %q", out["key3"], "val3")
	}
}

func TestParseSimpleYAMLColonInValue(t *testing.T) {
	data := []byte("url: https://example.com\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["url"] != "https://example.com" {
		t.Errorf("url = %v, want %q", out["url"], "https://example.com")
	}
}

func TestParseSimpleYAMLOnlyComments(t *testing.T) {
	data := []byte("# this is a comment\n# another comment\n")
	var out map[string]interface{}
	if err := parseSimpleYAML(data, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map for comment-only input, got %v", out)
	}
}

func TestPackageInfoFields(t *testing.T) {
	info := PackageInfo{
		Name:          "my-pkg",
		InstalledPath: "/path/to/.apm_modules/my-pkg",
		Ref:           "v1.2.3",
		Commit:        "deadbeef",
		Source:        "https://github.com/owner/my-pkg",
		Files:         []string{"SKILL.md", "apm.yml"},
		Versions:      []string{"v1.0.0", "v1.2.3"},
	}
	if info.Name != "my-pkg" {
		t.Errorf("Name mismatch: %q", info.Name)
	}
	if len(info.Files) != 2 {
		t.Errorf("Files length: got %d, want 2", len(info.Files))
	}
	if info.Versions[1] != "v1.2.3" {
		t.Errorf("Versions[1]: got %q, want %q", info.Versions[1], "v1.2.3")
	}
}
