package versionpins_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/versionpins"
)

func TestLoadRefPins_MissingFile(t *testing.T) {
	dir := t.TempDir()
	pins := versionpins.LoadRefPins(filepath.Join(dir, "nonexistent"))
	if len(pins) != 0 {
		t.Errorf("expected empty map for missing file, got %v", pins)
	}
}

func TestLoadRefPins_InvalidJSONContent(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "version-pins.json"), []byte("not json"), 0o644)
	pins := versionpins.LoadRefPins(dir)
	if len(pins) != 0 {
		t.Errorf("expected empty map for invalid JSON, got %v", pins)
	}
}

func TestCheckRefPin_FirstTime(t *testing.T) {
	dir := t.TempDir()
	prev := versionpins.CheckRefPin("mkt", "plugin", "sha1", "v1.0", dir)
	if prev != "" {
		t.Errorf("expected empty on first check, got %q", prev)
	}
}

func TestRecordRefPin_ThenCheck(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "plugin", "sha-abc", "v1.0", dir)
	prev := versionpins.CheckRefPin("mkt", "plugin", "sha-abc", "v1.0", dir)
	if prev != "" {
		t.Errorf("same ref should not return a previous pin, got %q", prev)
	}
}

func TestCheckRefPin_RefSwap(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "plugin", "sha-old", "v1.0", dir)
	prev := versionpins.CheckRefPin("mkt", "plugin", "sha-new", "v1.0", dir)
	if prev != "sha-old" {
		t.Errorf("expected old pin %q, got %q", "sha-old", prev)
	}
}

func TestRecordRefPin_OverwritesSameKey(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "plugin", "sha-1", "v1.0", dir)
	versionpins.RecordRefPin("mkt", "plugin", "sha-2", "v1.0", dir)
	prev := versionpins.CheckRefPin("mkt", "plugin", "sha-2", "v1.0", dir)
	if prev != "" {
		t.Errorf("after overwrite, same ref should not trigger swap, got %q", prev)
	}
}

func TestLoadRefPins_NonStringValues(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "version-pins.json"), []byte(`{"key": 123}`), 0o644)
	pins := versionpins.LoadRefPins(dir)
	if _, ok := pins["key"]; ok {
		t.Error("non-string values should be skipped")
	}
}

func TestCheckRefPin_DifferentVersionsSamePlugin(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "plugin", "sha-v1", "v1.0", dir)
	versionpins.RecordRefPin("mkt", "plugin", "sha-v2", "v2.0", dir)
	prev := versionpins.CheckRefPin("mkt", "plugin", "sha-v1", "v1.0", dir)
	if prev != "" {
		t.Errorf("v1.0 pin unchanged, should not trigger swap, got %q", prev)
	}
}

func TestCheckRefPin_EmptyVersion(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "plugin", "sha-noversion", "", dir)
	prev := versionpins.CheckRefPin("mkt", "plugin", "sha-noversion", "", dir)
	if prev != "" {
		t.Errorf("same ref with empty version should not trigger swap, got %q", prev)
	}
}

func TestSaveRefPins_CreatesDirs(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "a", "b", "c")
	versionpins.SaveRefPins(map[string]string{"k": "v"}, nested)
	pins := versionpins.LoadRefPins(nested)
	if pins["k"] != "v" {
		t.Errorf("expected pin 'v', got %q", pins["k"])
	}
}
