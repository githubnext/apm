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
