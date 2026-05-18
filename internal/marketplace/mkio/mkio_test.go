package mkio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mkio"
)

func TestAtomicWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	if err := mkio.AtomicWrite(p, []byte("hello")); err != nil {
		t.Fatalf("AtomicWrite error: %v", err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("content mismatch: %q", data)
	}
}

func TestAtomicWrite_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	_ = os.WriteFile(p, []byte("old"), 0o644)
	if err := mkio.AtomicWrite(p, []byte("new")); err != nil {
		t.Fatalf("AtomicWrite error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if string(data) != "new" {
		t.Fatalf("expected 'new', got %q", data)
	}
}

func TestAtomicWriteString(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	if err := mkio.AtomicWriteString(p, "test-content"); err != nil {
		t.Fatalf("AtomicWriteString error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if string(data) != "test-content" {
		t.Fatalf("content mismatch: %q", data)
	}
}

func TestAtomicWrite_NoTmpLeftover(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	_ = mkio.AtomicWrite(p, []byte("data"))
	// Tmp file should not exist after successful write
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("tmp file not cleaned up: %s", e.Name())
		}
	}
}

func TestAtomicWrite_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.txt")
	if err := mkio.AtomicWrite(p, []byte{}); err != nil {
		t.Fatalf("AtomicWrite empty error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if len(data) != 0 {
		t.Errorf("expected empty file, got %q", data)
	}
}

func TestAtomicWrite_LargeContent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "big.txt")
	content := make([]byte, 64*1024)
	for i := range content {
		content[i] = byte(i % 256)
	}
	if err := mkio.AtomicWrite(p, content); err != nil {
		t.Fatalf("AtomicWrite large error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if len(data) != len(content) {
		t.Errorf("size mismatch: want %d, got %d", len(content), len(data))
	}
}

func TestAtomicWriteString_EmptyString(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.txt")
	if err := mkio.AtomicWriteString(p, ""); err != nil {
		t.Fatalf("AtomicWriteString empty error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if len(data) != 0 {
		t.Errorf("expected empty file, got %q", data)
	}
}

func TestAtomicWrite_SubdirMissing(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "nonexistent", "out.txt")
	// Should error since directory does not exist
	err := mkio.AtomicWrite(p, []byte("data"))
	if err == nil {
		t.Error("expected error when parent dir missing")
	}
}

func TestAtomicWrite_Idempotent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	for i := 0; i < 3; i++ {
		if err := mkio.AtomicWrite(p, []byte("same")); err != nil {
			t.Fatalf("iteration %d: AtomicWrite error: %v", i, err)
		}
	}
	data, _ := os.ReadFile(p)
	if string(data) != "same" {
		t.Errorf("final content = %q, want 'same'", data)
	}
}
