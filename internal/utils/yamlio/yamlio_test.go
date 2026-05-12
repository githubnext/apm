package yamlio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/yamlio"
)

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
