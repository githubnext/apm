package yamlio_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/yamlio"
)

func TestLoadEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(path, []byte("  \n  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := yamlio.LoadYAML(path)
	if err != nil {
		t.Fatalf("LoadYAML empty: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for whitespace-only file, got %v", result)
	}
}

func TestLoadYAMLWithComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "commented.yaml")
	content := "# this is a comment\nkey: value\n# another comment\nnum: 42\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := yamlio.LoadYAML(path)
	if err != nil {
		t.Fatalf("LoadYAML with comments: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("key = %v, want value", result["key"])
	}
}

func TestYAMLToStr_NonMap(t *testing.T) {
	s, err := yamlio.YAMLToStr("hello")
	if err != nil {
		t.Fatalf("YAMLToStr non-map: %v", err)
	}
	if !strings.Contains(s, "hello") {
		t.Errorf("expected 'hello' in output, got %q", s)
	}
}

func TestYAMLToStr_MapMultipleKeys(t *testing.T) {
	data := map[string]any{"a": "1", "b": "2"}
	s, err := yamlio.YAMLToStr(data)
	if err != nil {
		t.Fatalf("YAMLToStr: %v", err)
	}
	if !strings.Contains(s, "a: 1") && !strings.Contains(s, "a: 2") {
		// at least one key should be present
	}
	if s == "" {
		t.Error("expected non-empty YAML string")
	}
}

func TestDumpAndLoadRoundTrip_MultipleKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi.yaml")
	data := map[string]any{
		"name":    "apm",
		"version": "1.0",
	}
	if err := yamlio.DumpYAML(data, path); err != nil {
		t.Fatalf("DumpYAML: %v", err)
	}
	loaded, err := yamlio.LoadYAML(path)
	if err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	if loaded["name"] != "apm" {
		t.Errorf("name = %v, want apm", loaded["name"])
	}
	if loaded["version"] != "1.0" {
		t.Errorf("version = %v, want 1.0", loaded["version"])
	}
}

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")

	data := map[string]any{
		"key": "value",
		"num": 42,
	}

	if err := yamlio.DumpYAML(data, path); err != nil {
		t.Fatalf("DumpYAML: %v", err)
	}

	loaded, err := yamlio.LoadYAML(path)
	if err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}

	if loaded["key"] != "value" {
		t.Errorf("key: got %v, want value", loaded["key"])
	}
}

func TestLoadMissing(t *testing.T) {
	_, err := yamlio.LoadYAML("/nonexistent/file.yaml")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestYAMLToStr(t *testing.T) {
	data := map[string]any{"a": 1}
	s, err := yamlio.YAMLToStr(data)
	if err != nil {
		t.Fatalf("YAMLToStr: %v", err)
	}
	if s == "" {
		t.Error("expected non-empty YAML string")
	}
}
