package cachepin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/cachepin"
)

func TestWriteAndVerifyMarker_Happy(t *testing.T) {
	dir := t.TempDir()
	sha := "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	cachepin.WriteMarker(dir, sha)

	if err := cachepin.VerifyMarker(dir, sha); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyMarker_Missing(t *testing.T) {
	dir := t.TempDir()
	err := cachepin.VerifyMarker(dir, "abc")
	if err == nil {
		t.Fatal("expected error for missing marker")
	}
	if !cachepin.IsCachePinError(err) {
		t.Errorf("expected CachePinError, got %T", err)
	}
}

func TestVerifyMarker_Mismatch(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1")
	err := cachepin.VerifyMarker(dir, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb2")
	if err == nil {
		t.Fatal("expected mismatch error")
	}
	if !cachepin.IsCachePinError(err) {
		t.Errorf("expected CachePinError, got %T", err)
	}
}

func TestVerifyMarker_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, cachepin.MarkerFilename), []byte("{not json}"), 0o644)
	err := cachepin.VerifyMarker(dir, "sha")
	if err == nil || !cachepin.IsCachePinError(err) {
		t.Fatalf("expected CachePinError for malformed JSON, got %v", err)
	}
}

func TestVerifyMarker_WrongSchemaVersion(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, cachepin.MarkerFilename), []byte(`{"schema_version":99,"resolved_commit":"sha"}`), 0o644)
	err := cachepin.VerifyMarker(dir, "sha")
	if err == nil || !cachepin.IsCachePinError(err) {
		t.Fatalf("expected CachePinError for wrong schema_version, got %v", err)
	}
}

func TestVerifyMarker_MissingResolvedCommit(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, cachepin.MarkerFilename), []byte(`{"schema_version":1,"resolved_commit":""}`), 0o644)
	err := cachepin.VerifyMarker(dir, "sha")
	if err == nil || !cachepin.IsCachePinError(err) {
		t.Fatalf("expected CachePinError for missing resolved_commit, got %v", err)
	}
}

func TestWriteMarker_NonexistentDir_IsNoop(t *testing.T) {
	cachepin.WriteMarker("/tmp/nonexistent-dir-xyz-autoloop", "sha")
	// should not panic
}

func TestConstants(t *testing.T) {
	if cachepin.MarkerFilename == "" {
		t.Fatal("MarkerFilename must not be empty")
	}
	if cachepin.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion expected 1, got %d", cachepin.SchemaVersion)
	}
}
