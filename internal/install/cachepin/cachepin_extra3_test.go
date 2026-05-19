package cachepin_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/cachepin"
)

func TestVerifyMarker_MalformedJSONIsError(t *testing.T) {
	dir := t.TempDir()
	markerPath := filepath.Join(dir, ".apm-pin")
	_ = os.WriteFile(markerPath, []byte("not-json"), 0o600)
	err := cachepin.VerifyMarker(dir, "sha")
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
	if !cachepin.IsCachePinError(err) {
		t.Error("expected IsCachePinError=true for malformed JSON")
	}
}

func TestVerifyMarker_WrongSchemaVersionIsError(t *testing.T) {
	dir := t.TempDir()
	payload := map[string]interface{}{
		"schema_version":  99,
		"resolved_commit": "abc",
	}
	data, _ := json.Marshal(payload)
	_ = os.WriteFile(filepath.Join(dir, ".apm-pin"), data, 0o600)
	err := cachepin.VerifyMarker(dir, "abc")
	if err == nil {
		t.Error("expected error for unsupported schema_version")
	}
	if !cachepin.IsCachePinError(err) {
		t.Error("expected IsCachePinError=true")
	}
}

func TestVerifyMarker_EmptyResolvedCommitInFile(t *testing.T) {
	dir := t.TempDir()
	payload := map[string]interface{}{
		"schema_version":  1,
		"resolved_commit": "",
	}
	data, _ := json.Marshal(payload)
	_ = os.WriteFile(filepath.Join(dir, ".apm-pin"), data, 0o600)
	err := cachepin.VerifyMarker(dir, "abc")
	if err == nil {
		t.Error("expected error when resolved_commit field is empty")
	}
}

func TestMarkerFilename_Constant(t *testing.T) {
	if cachepin.MarkerFilename != ".apm-pin" {
		t.Errorf("expected MarkerFilename='.apm-pin', got %q", cachepin.MarkerFilename)
	}
}

func TestSchemaVersion_IsOne(t *testing.T) {
	if cachepin.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion=1, got %d", cachepin.SchemaVersion)
	}
}

func TestWriteAndVerify_UnicodeCommit(t *testing.T) {
	dir := t.TempDir()
	commit := "sha-with-utf8-safe-ascii-only-012345678901234"
	cachepin.WriteMarker(dir, commit)
	if err := cachepin.VerifyMarker(dir, commit); err != nil {
		t.Errorf("expected verification pass: %v", err)
	}
}

func TestVerifyMarker_FilePermissionUnreadable(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root can read any file")
	}
	dir := t.TempDir()
	markerPath := filepath.Join(dir, ".apm-pin")
	_ = os.WriteFile(markerPath, []byte(`{"schema_version":1,"resolved_commit":"sha"}`), 0o600)
	_ = os.Chmod(markerPath, 0o000)
	defer func() { _ = os.Chmod(markerPath, 0o600) }()
	err := cachepin.VerifyMarker(dir, "sha")
	if err == nil {
		t.Error("expected error for unreadable marker file")
	}
}

func TestWriteMarker_DirectoryPathNotFile(t *testing.T) {
	// Passing a dir path that exists but is itself a directory should be silent
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "sha-ok")
}

func TestIsCachePinError_ErrorInterfaceOnly(t *testing.T) {
	var err error
	if cachepin.IsCachePinError(err) {
		t.Error("nil error should return false")
	}
}
