package yamlio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAML_SingleKeyValue(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.yaml")
	if err := os.WriteFile(f, []byte("key: value\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadYAML(f)
	if err != nil {
		t.Fatal(err)
	}
	if m["key"] != "value" {
		t.Errorf("expected value, got %v", m["key"])
	}
}

func TestLoadYAML_IntValue(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "b.yaml")
	if err := os.WriteFile(f, []byte("count: 42\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadYAML(f)
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("expected non-nil map")
	}
}

func TestLoadYAML_MissingFileReturnsError(t *testing.T) {
	_, err := LoadYAML("/nonexistent/path/file.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDumpYAML_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "rt.yaml")
	data := map[string]any{"hello": "world"}
	if err := DumpYAML(data, f); err != nil {
		t.Fatal(err)
	}
	m, err := LoadYAML(f)
	if err != nil {
		t.Fatal(err)
	}
	if m["hello"] != "world" {
		t.Errorf("round-trip failed: got %v", m["hello"])
	}
}

func TestDumpYAML_InvalidPathReturnsError(t *testing.T) {
	err := DumpYAML(map[string]any{"k": "v"}, "/nonexistent/dir/out.yaml")
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestYAMLToStr_EmptyMapReturnsString(t *testing.T) {
	s, err := YAMLToStr(map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	_ = s // may be empty string
}

func TestYAMLToStr_NonMapValue(t *testing.T) {
	s, err := YAMLToStr("plain string")
	if err != nil {
		t.Fatal(err)
	}
	if s == "" {
		t.Error("non-map should produce a non-empty string")
	}
}

func TestYAMLToStr_MultiKey(t *testing.T) {
	m := map[string]any{"a": "1", "b": "2"}
	s, err := YAMLToStr(m)
	if err != nil {
		t.Fatal(err)
	}
	if s == "" {
		t.Error("expected non-empty YAML output")
	}
}

func TestLoadYAML_CommentLineSkipped(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "c.yaml")
	content := "# this is a comment\nkey: val\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadYAML(f)
	if err != nil {
		t.Fatal(err)
	}
	if m["key"] != "val" {
		t.Errorf("expected val, got %v", m["key"])
	}
}
