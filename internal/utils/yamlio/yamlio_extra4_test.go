package yamlio_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/yamlio"
)

func TestLoadYAML_MultipleKeys_Extra4(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "multi.yaml")
	_ = os.WriteFile(p, []byte("key1: val1\nkey2: val2\n"), 0o644)
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["key1"] != "val1" || m["key2"] != "val2" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestLoadYAML_SkipsComments_Extra4(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "comments.yaml")
	_ = os.WriteFile(p, []byte("# comment\nfoo: bar\n"), 0o644)
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["foo"] != "bar" {
		t.Fatalf("expected bar, got %v", m["foo"])
	}
	if _, ok := m["# comment"]; ok {
		t.Fatal("comment line should not be a key")
	}
}

func TestYAMLToStr_NonMap_Extra4(t *testing.T) {
	out, err := yamlio.YAMLToStr("just a string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "just a string") {
		t.Fatalf("expected string in output, got %s", out)
	}
}

func TestYAMLToStr_EmptyMap_Extra4(t *testing.T) {
	out, err := yamlio.YAMLToStr(map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty output, got %s", out)
	}
}

func TestDumpYAML_MultipleKeys_Extra4(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.yaml")
	m := map[string]any{"alpha": "one", "beta": "two"}
	if err := yamlio.DumpYAML(m, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	s := string(data)
	if !strings.Contains(s, "alpha") || !strings.Contains(s, "beta") {
		t.Fatalf("expected keys in output: %s", s)
	}
}

func TestLoadYAML_WhitespaceOnly_Extra4(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "whitespace.yaml")
	_ = os.WriteFile(p, []byte("   \n\t\n"), 0o644)
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m != nil {
		t.Fatalf("expected nil for whitespace-only file, got %v", m)
	}
}

func TestYAMLToStr_IntValue_Extra4(t *testing.T) {
	out, err := yamlio.YAMLToStr(map[string]any{"count": 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "42") {
		t.Fatalf("expected 42 in output, got %s", out)
	}
}

func TestLoadYAML_ValueWithColon_Extra4(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "colon.yaml")
	_ = os.WriteFile(p, []byte("url: http://example.com\n"), 0o644)
	m, err := yamlio.LoadYAML(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["url"]; !ok {
		t.Fatal("expected url key")
	}
}
