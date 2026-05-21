package view

import "testing"

func TestViewOptions_VerboseTrue_Extra3(t *testing.T) {
	opts := ViewOptions{Verbose: true, Format: "text"}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
	if opts.Format != "text" {
		t.Errorf("Format = %q, want text", opts.Format)
	}
}

func TestViewOptions_JSONFormat_Extra3(t *testing.T) {
	opts := ViewOptions{Format: "json"}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want json", opts.Format)
	}
}

func TestViewOptions_FieldFilter_Extra3(t *testing.T) {
	opts := ViewOptions{Field: "versions"}
	if opts.Field != "versions" {
		t.Errorf("Field = %q, want versions", opts.Field)
	}
}

func TestPackageInfo_NameOnly_Extra3(t *testing.T) {
	info := PackageInfo{Name: "mypkg"}
	if info.Name != "mypkg" {
		t.Errorf("Name = %q, want mypkg", info.Name)
	}
}

func TestPackageInfo_Versions_Extra3(t *testing.T) {
	info := PackageInfo{
		Name:     "mypkg",
		Versions: []string{"v1.0.0", "v1.1.0", "v2.0.0"},
	}
	if len(info.Versions) != 3 {
		t.Errorf("Versions len = %d, want 3", len(info.Versions))
	}
	if info.Versions[2] != "v2.0.0" {
		t.Errorf("Versions[2] = %q, want v2.0.0", info.Versions[2])
	}
}

func TestPackageInfo_Files_Extra3(t *testing.T) {
	info := PackageInfo{
		Files: []string{"a.go", "b.go"},
	}
	if len(info.Files) != 2 {
		t.Errorf("Files len = %d, want 2", len(info.Files))
	}
}

func TestPackageInfo_RefAndCommit_Extra3(t *testing.T) {
	info := PackageInfo{Ref: "main", Commit: "abc123"}
	if info.Ref != "main" {
		t.Errorf("Ref = %q, want main", info.Ref)
	}
	if info.Commit != "abc123" {
		t.Errorf("Commit = %q, want abc123", info.Commit)
	}
}

func TestParseSimpleYAML_SingleKey_Extra3(t *testing.T) {
	data := []byte("name: mypkg\n")
	m := make(map[string]interface{})
	if err := parseSimpleYAML(data, &m); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	if m["name"] != "mypkg" {
		t.Errorf("name = %v, want mypkg", m["name"])
	}
}

func TestParseSimpleYAML_MultiKey_Extra3(t *testing.T) {
	data := []byte("name: foo\nversion: 1.2.3\n")
	m := make(map[string]interface{})
	if err := parseSimpleYAML(data, &m); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	if m["name"] != "foo" || m["version"] != "1.2.3" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestParseSimpleYAML_EmptyData_Extra3(t *testing.T) {
	data := []byte("")
	m := make(map[string]interface{})
	if err := parseSimpleYAML(data, &m); err != nil {
		t.Fatalf("parseSimpleYAML error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestPackageInfo_Source_Extra3(t *testing.T) {
	info := PackageInfo{Source: "https://github.com/owner/repo"}
	if info.Source != "https://github.com/owner/repo" {
		t.Errorf("Source = %q", info.Source)
	}
}
