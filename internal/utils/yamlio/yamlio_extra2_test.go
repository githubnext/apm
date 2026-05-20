package yamlio_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/yamlio"
)

func TestLoadYAML_FileNotExist(t *testing.T) {
	_, err := yamlio.LoadYAML("/nonexistent/path/file.yaml")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestLoadYAML_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(p, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m != nil {
		t.Errorf("expected nil for empty YAML, got %v", m)
	}
}

func TestLoadYAML_CommentOnly(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "comment.yaml")
	if err := os.WriteFile(p, []byte("# This is a comment\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Comment-only YAML returns nil or empty map (implementation detail)
	if m != nil && len(m) != 0 {
		t.Errorf("expected nil or empty map for comment-only YAML, got %v", m)
	}
}

func TestLoadYAML_BooleanValue(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bool.yaml")
	if err := os.WriteFile(p, []byte("enabled: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["enabled"] != "true" {
		t.Errorf("expected 'true', got %q", m["enabled"])
	}
}

func TestYAMLToStr_SimpleMap(t *testing.T) {
	data := map[string]any{"key": "value"}
	out, err := yamlio.YAMLToStr(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "key: value") {
		t.Errorf("expected 'key: value' in output, got %q", out)
	}
}

func TestYAMLToStr_NonMapInput(t *testing.T) {
	out, err := yamlio.YAMLToStr("plain string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "plain string") {
		t.Errorf("expected string content in output, got %q", out)
	}
}

func TestDumpYAML_WritesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.yaml")
	data := map[string]any{"foo": "bar"}
	if err := yamlio.DumpYAML(data, p); err != nil {
		t.Fatalf("DumpYAML: %v", err)
	}
	content, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(content), "foo: bar") {
		t.Errorf("expected 'foo: bar' in file, got %q", string(content))
	}
}

func TestDumpYAML_BadPath(t *testing.T) {
	err := yamlio.DumpYAML(map[string]any{"k": "v"}, "/nonexistent/dir/file.yaml")
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestLoadYAML_MultipleKeys(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "multi.yaml")
	content := "alpha: 1\nbeta: 2\ngamma: 3\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 3 {
		t.Errorf("expected 3 keys, got %d", len(m))
	}
	if m["alpha"] != "1" || m["beta"] != "2" || m["gamma"] != "3" {
		t.Errorf("unexpected values: %v", m)
	}
}

func TestYAMLToStr_EmptyMapOutput(t *testing.T) {
	data := map[string]any{}
	out, err := yamlio.YAMLToStr(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty output for empty map, got %q", out)
	}
}
