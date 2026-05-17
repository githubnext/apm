package cachepin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/cachepin"
)

func TestWriteMarker_ThenReadBack(t *testing.T) {
	dir := t.TempDir()
	commit := "1234567890abcdef1234567890abcdef12345678"
	cachepin.WriteMarker(dir, commit)

	markerPath := filepath.Join(dir, cachepin.MarkerFilename)
	data, err := os.ReadFile(markerPath)
	if err != nil {
		t.Fatalf("marker file not found after write: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("marker file is empty")
	}
}

func TestVerifyMarker_EmptyCommit(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "realcommit")
	err := cachepin.VerifyMarker(dir, "")
	if err == nil {
		t.Fatal("expected error when expected commit is empty but stored is non-empty")
	}
}

func TestVerifyMarker_ExactMatch(t *testing.T) {
	dir := t.TempDir()
	commit := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	cachepin.WriteMarker(dir, commit)
	if err := cachepin.VerifyMarker(dir, commit); err != nil {
		t.Fatalf("expected no error for exact match, got %v", err)
	}
}

func TestCachePinError_ErrorString(t *testing.T) {
	dir := t.TempDir()
	err := cachepin.VerifyMarker(dir, "sha")
	if err == nil {
		t.Fatal("expected error for missing marker")
	}
	if err.Error() == "" {
		t.Fatal("error string should not be empty")
	}
}

func TestWriteMarker_OverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	commit1 := "1111111111111111111111111111111111111111"
	commit2 := "2222222222222222222222222222222222222222"
	cachepin.WriteMarker(dir, commit1)
	cachepin.WriteMarker(dir, commit2)
	if err := cachepin.VerifyMarker(dir, commit2); err != nil {
		t.Fatalf("expected commit2 after overwrite, got error: %v", err)
	}
	if err := cachepin.VerifyMarker(dir, commit1); err == nil {
		t.Fatal("expected error for old commit after overwrite")
	}
}

func TestIsCachePinError_WithOtherError(t *testing.T) {
	err := os.ErrNotExist
	if cachepin.IsCachePinError(err) {
		t.Error("os.ErrNotExist should not be a CachePinError")
	}
}
