package versionpins_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/versionpins"
)

func TestLoadRefPins_Missing(t *testing.T) {
	dir := t.TempDir()
	pins := versionpins.LoadRefPins(dir)
	if len(pins) != 0 {
		t.Errorf("expected empty map for missing file, got %v", pins)
	}
}

func TestLoadRefPins_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "version-pins.json"), []byte("not-json"), 0o644)
	pins := versionpins.LoadRefPins(dir)
	if len(pins) != 0 {
		t.Error("expected empty map for invalid JSON")
	}
}

func TestSaveAndLoadRefPins(t *testing.T) {
	dir := t.TempDir()
	original := map[string]string{
		"marketplace/plugin/1.0": "abc123",
		"other/tool/2.0":        "def456",
	}
	versionpins.SaveRefPins(original, dir)

	loaded := versionpins.LoadRefPins(dir)
	for k, v := range original {
		if loaded[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, loaded[k])
		}
	}
}

func TestSaveRefPins_Atomic(t *testing.T) {
	dir := t.TempDir()
	pins := map[string]string{"k": "v"}
	versionpins.SaveRefPins(pins, dir)

	data, err := os.ReadFile(filepath.Join(dir, "version-pins.json"))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
}

func TestCheckRefPin_NewPin(t *testing.T) {
	dir := t.TempDir()
	warn := versionpins.CheckRefPin("mp", "plugin", "sha1", "1.0", dir)
	if warn != "" {
		t.Errorf("expected no warning for new pin, got: %s", warn)
	}
}

func TestCheckRefPin_SameRef(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mp", "plugin", "sha1", "1.0", dir)
	warn := versionpins.CheckRefPin("mp", "plugin", "sha1", "1.0", dir)
	if warn != "" {
		t.Errorf("expected no warning for same ref, got: %s", warn)
	}
}

func TestCheckRefPin_ChangedRef(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mp", "plugin", "sha1", "1.0", dir)
	warn := versionpins.CheckRefPin("mp", "plugin", "sha2", "1.0", dir)
	if warn == "" {
		t.Error("expected warning when ref changes")
	}
	if warn != "sha1" {
		t.Errorf("expected previous ref 'sha1', got %q", warn)
	}
}

func TestRecordRefPin_Overwrite(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mp", "plugin", "sha1", "1.0", dir)
	versionpins.RecordRefPin("mp", "plugin", "sha2", "1.0", dir)
	pins := versionpins.LoadRefPins(dir)
	for _, v := range pins {
		if v == "sha2" {
			return
		}
	}
	t.Error("expected overwritten pin sha2 to be present")
}
