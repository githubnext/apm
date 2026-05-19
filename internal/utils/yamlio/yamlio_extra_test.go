package yamlio_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/yamlio"
)

func TestLoadYAML_KeyValueParsed(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "kv.yaml")
	if err := os.WriteFile(p, []byte("name: alice\nage: 30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["name"] != "alice" {
		t.Errorf("name: got %q want %q", m["name"], "alice")
	}
	if m["age"] != "30" {
		t.Errorf("age: got %q want %q", m["age"], "30")
	}
}

func TestLoadYAML_SkipsBlankLines(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "blank.yaml")
	content := "key1: val1\n\nkey2: val2\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
}

func TestLoadYAML_WhitespaceOnlyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ws.yaml")
	if err := os.WriteFile(p, []byte("   \n  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m != nil {
		t.Errorf("expected nil map for whitespace-only file, got %v", m)
	}
}

func TestDumpYAML_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.yaml")
	data := map[string]any{"hello": "world"}
	if err := yamlio.DumpYAML(data, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if !strings.Contains(string(raw), "hello") {
		t.Errorf("expected 'hello' in output, got %q", string(raw))
	}
}

func TestYAMLToStr_EmptyMap(t *testing.T) {
	out, err := yamlio.YAMLToStr(map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string for empty map, got %q", out)
	}
}

func TestYAMLToStr_SingleEntry(t *testing.T) {
	out, err := yamlio.YAMLToStr(map[string]any{"k": "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "k: v") {
		t.Errorf("expected 'k: v' in output, got %q", out)
	}
}

func TestYAMLToStr_IntValue(t *testing.T) {
	out, err := yamlio.YAMLToStr(map[string]any{"count": 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "42") {
		t.Errorf("expected '42' in output, got %q", out)
	}
}

func TestYAMLToStr_BoolValue(t *testing.T) {
	out, err := yamlio.YAMLToStr(map[string]any{"enabled": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "true") {
		t.Errorf("expected 'true' in output, got %q", out)
	}
}

func TestDumpYAML_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "overwrite.yaml")
	if err := os.WriteFile(p, []byte("old: data\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	newData := map[string]any{"new": "content"}
	if err := yamlio.DumpYAML(newData, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "old") {
		t.Errorf("expected old content to be replaced, got %q", string(raw))
	}
	if !strings.Contains(string(raw), "new") {
		t.Errorf("expected 'new' in output, got %q", string(raw))
	}
}

func TestLoadYAML_ValueWithColon(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "colon.yaml")
	if err := os.WriteFile(p, []byte("url: https://example.com\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["url"] != "https://example.com" {
		t.Errorf("url: got %q want %q", m["url"], "https://example.com")
	}
}

func TestLoadYAML_CommentLinesSkipped(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "comments.yaml")
	content := "# top comment\nkey: val\n# another comment\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(m), m)
	}
}
