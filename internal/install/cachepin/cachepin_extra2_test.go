package cachepin_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/cachepin"
)

func TestWriteAndVerify_LongSHA(t *testing.T) {
	dir := t.TempDir()
	sha := "aabbccddeeff0011223344556677889900112233445566778899aabbccddeeff00"
	cachepin.WriteMarker(dir, sha)
	if err := cachepin.VerifyMarker(dir, sha); err != nil {
		t.Errorf("expected verification to pass for long SHA: %v", err)
	}
}

func TestVerifyMarker_WrongCommit(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "sha-one")
	err := cachepin.VerifyMarker(dir, "sha-two")
	if err == nil {
		t.Error("expected error for mismatched commit")
	}
}

func TestIsCachePinError_DirectError(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "commit-a")
	err := cachepin.VerifyMarker(dir, "commit-b")
	if err == nil {
		t.Skip("no mismatch error generated")
	}
	if !cachepin.IsCachePinError(err) {
		t.Error("expected IsCachePinError=true for mismatch error")
	}
}

func TestIsCachePinError_WrappedError(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "sha-a")
	inner := cachepin.VerifyMarker(dir, "sha-b")
	if inner == nil {
		t.Skip("no mismatch error generated")
	}
	wrapped := errors.Join(errors.New("outer"), inner)
	if !cachepin.IsCachePinError(wrapped) {
		t.Error("expected IsCachePinError=true for wrapped mismatch error")
	}
}

func TestWriteMarker_CreatesPinFile(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "abc123")
	pinFile := filepath.Join(dir, ".apm-pin")
	if _, err := os.Stat(pinFile); os.IsNotExist(err) {
		t.Error("expected .apm-pin file to be created")
	}
}

func TestVerifyMarker_MissingFile(t *testing.T) {
	dir := t.TempDir()
	err := cachepin.VerifyMarker(dir, "sha1")
	if err == nil {
		t.Error("expected error when .apm-pin file is missing")
	}
}

func TestVerifyMarker_EmptyExpected(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "sha1")
	err := cachepin.VerifyMarker(dir, "")
	if err == nil {
		t.Error("expected error when expected commit is empty")
	}
}

func TestWriteMarker_EmptyCommit(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "")
}

func TestWriteMarker_IdempotentOverwrite(t *testing.T) {
	dir := t.TempDir()
	cachepin.WriteMarker(dir, "first-sha")
	cachepin.WriteMarker(dir, "second-sha")
	if err := cachepin.VerifyMarker(dir, "second-sha"); err != nil {
		t.Errorf("expected second-sha after overwrite: %v", err)
	}
}

func TestIsCachePinError_StandardError(t *testing.T) {
	err := errors.New("some random error")
	if cachepin.IsCachePinError(err) {
		t.Error("expected IsCachePinError=false for standard error")
	}
}

func TestIsCachePinError_NilInput(t *testing.T) {
	if cachepin.IsCachePinError(nil) {
		t.Error("expected IsCachePinError=false for nil error")
	}
}

func TestWriteMarker_NonexistentDir_NoopOrCreate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent_subdir_xyz")
	cachepin.WriteMarker(dir, "sha")
}
