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
