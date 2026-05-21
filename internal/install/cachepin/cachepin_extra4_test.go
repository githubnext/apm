package cachepin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteMarker_MissingDir_Extra4(t *testing.T) {
	// Should be silent / no-op for missing directory
	WriteMarker("/nonexistent/dir", "abc123")
}

func TestWriteMarker_IsFile_Extra4(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	_ = os.WriteFile(f, []byte("x"), 0o644)
	// Should be silent for non-dir path
	WriteMarker(f, "abc123")
}

func TestWriteMarker_WritesMarker_Extra4(t *testing.T) {
	dir := t.TempDir()
	WriteMarker(dir, "abc123")
	data, err := os.ReadFile(filepath.Join(dir, MarkerFilename))
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["resolved_commit"] != "abc123" {
		t.Errorf("expected abc123, got %v", m["resolved_commit"])
	}
}

func TestVerifyMarker_Missing_Extra4(t *testing.T) {
	err := VerifyMarker("/nonexistent/dir", "abc")
	if err == nil {
		t.Error("expected error for missing marker")
	}
	if !IsCachePinError(err) {
		t.Errorf("expected CachePinError, got %T", err)
	}
}

func TestVerifyMarker_Matches_Extra4(t *testing.T) {
	dir := t.TempDir()
	WriteMarker(dir, "sha-xyz")
	err := VerifyMarker(dir, "sha-xyz")
	if err != nil {
		t.Errorf("expected no error for matching commit, got %v", err)
	}
}

func TestVerifyMarker_Mismatch_Extra4(t *testing.T) {
	dir := t.TempDir()
	WriteMarker(dir, "sha-abc")
	err := VerifyMarker(dir, "sha-xyz")
	if err == nil {
		t.Error("expected error for mismatched commit")
	}
	if !IsCachePinError(err) {
		t.Errorf("expected CachePinError, got %T", err)
	}
}

func TestVerifyMarker_CorruptJSON_Extra4(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, MarkerFilename), []byte("not-json"), 0o644)
	err := VerifyMarker(dir, "sha")
	if err == nil {
		t.Error("expected error for corrupt JSON")
	}
}

func TestIsCachePinError_True_Extra4(t *testing.T) {
	err := &CachePinError{Msg: "test"}
	if !IsCachePinError(err) {
		t.Error("expected IsCachePinError true")
	}
}

func TestIsCachePinError_False_Extra4(t *testing.T) {
	if IsCachePinError(os.ErrNotExist) {
		t.Error("expected IsCachePinError false for os.ErrNotExist")
	}
}

func TestMarkerFilename_NotEmpty_Extra4(t *testing.T) {
	if MarkerFilename == "" {
		t.Error("expected non-empty MarkerFilename")
	}
}

func TestSchemaVersion_Extra4(t *testing.T) {
	if SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion=1, got %d", SchemaVersion)
	}
}

func TestCachePinError_Error_Extra4(t *testing.T) {
	e := &CachePinError{Msg: "my-msg"}
	if e.Error() != "my-msg" {
		t.Errorf("expected my-msg, got %s", e.Error())
	}
}

func TestWriteMarker_Idempotent_Extra4(t *testing.T) {
	dir := t.TempDir()
	WriteMarker(dir, "sha1")
	WriteMarker(dir, "sha2")
	err := VerifyMarker(dir, "sha2")
	if err != nil {
		t.Errorf("expected sha2 after second write, got %v", err)
	}
}
