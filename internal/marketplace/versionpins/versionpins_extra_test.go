package versionpins_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/versionpins"
)

func TestLoadRefPins_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "version-pins.json"), []byte("{}"), 0o644)
	pins := versionpins.LoadRefPins(dir)
	if len(pins) != 0 {
		t.Errorf("expected empty map for empty JSON object, got %v", pins)
	}
}

func TestSaveRefPins_EmptyMap(t *testing.T) {
	dir := t.TempDir()
	versionpins.SaveRefPins(map[string]string{}, dir)
	pins := versionpins.LoadRefPins(dir)
	if len(pins) != 0 {
		t.Errorf("expected empty pins after saving empty map, got %v", pins)
	}
}

func TestSaveAndLoadRefPins_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	original := map[string]string{
		"mp/a/1.0": "sha-aaa",
		"mp/b/2.0": "sha-bbb",
		"mp/c/3.0": "sha-ccc",
	}
	versionpins.SaveRefPins(original, dir)
	loaded := versionpins.LoadRefPins(dir)
	if len(loaded) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(loaded))
	}
	for k, v := range original {
		if loaded[k] != v {
			t.Errorf("key %q: got %q, want %q", k, loaded[k], v)
		}
	}
}

func TestCheckRefPin_AfterRecord(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("market", "plugin", "abc123", "1.0", dir)
	warn := versionpins.CheckRefPin("market", "plugin", "abc123", "1.0", dir)
	if warn != "" {
		t.Errorf("expected no warning for matching ref, got %q", warn)
	}
}

func TestCheckRefPin_DifferentVersion_NewPin(t *testing.T) {
	dir := t.TempDir()
	// Record for version 1.0 with sha1
	versionpins.RecordRefPin("market", "plugin", "sha1", "1.0", dir)
	// Check for version 2.0 (different key) -- should be new pin
	warn := versionpins.CheckRefPin("market", "plugin", "sha2", "2.0", dir)
	if warn != "" {
		t.Errorf("expected no warning for different version (new key), got %q", warn)
	}
}

func TestRecordRefPin_IdempotentSameRef(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "pkg", "refX", "1.0", dir)
	versionpins.RecordRefPin("mkt", "pkg", "refX", "1.0", dir)
	pins := versionpins.LoadRefPins(dir)
	count := 0
	for _, v := range pins {
		if v == "refX" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 entry for refX, got %d", count)
	}
}

func TestCheckRefPin_ReturnsOldRef(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt", "pkg", "old-sha", "1.0", dir)
	warn := versionpins.CheckRefPin("mkt", "pkg", "new-sha", "1.0", dir)
	if warn != "old-sha" {
		t.Errorf("expected old-sha warning, got %q", warn)
	}
}

func TestSaveRefPins_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	versionpins.SaveRefPins(map[string]string{"k": "v1"}, dir)
	versionpins.SaveRefPins(map[string]string{"k": "v2"}, dir)
	loaded := versionpins.LoadRefPins(dir)
	if loaded["k"] != "v2" {
		t.Errorf("expected overwritten value v2, got %q", loaded["k"])
	}
}

func TestRecordAndCheckPinDifferentMarketplaces(t *testing.T) {
	dir := t.TempDir()
	versionpins.RecordRefPin("mkt1", "plugin", "sha-a", "1.0", dir)
	versionpins.RecordRefPin("mkt2", "plugin", "sha-b", "1.0", dir)
	// Checking mkt1 should return sha-a (no warning since it matches)
	warn1 := versionpins.CheckRefPin("mkt1", "plugin", "sha-a", "1.0", dir)
	if warn1 != "" {
		t.Errorf("mkt1 warn: expected empty, got %q", warn1)
	}
	// Checking mkt2 should also return no warning since sha-b matches
	warn2 := versionpins.CheckRefPin("mkt2", "plugin", "sha-b", "1.0", dir)
	if warn2 != "" {
		t.Errorf("mkt2 warn: expected empty, got %q", warn2)
	}
}
