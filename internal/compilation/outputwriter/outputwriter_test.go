package outputwriter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, "# content\n"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty file")
	}
}

func TestWrite_CreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "AGENTS.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, "hello"); err != nil {
		t.Fatalf("unexpected error creating nested path: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestWrite_Idempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	w := &CompiledOutputWriter{}
	content := "# hello\nworld\n"
	if err := w.Write(path, content); err != nil {
		t.Fatal(err)
	}
	if err := w.Write(path, content); err != nil {
		t.Fatalf("second write failed: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) == "" {
		t.Error("file should not be empty after idempotent write")
	}
}
