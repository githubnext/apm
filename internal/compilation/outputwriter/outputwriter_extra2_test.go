package outputwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestWrite_CreatesDirectories(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "nested", "dir", "output.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, "content\n"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file should exist after write: %v", err)
	}
}

func TestWrite_BuildIDPlaceholderReplaced(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "out.md")
	w := &CompiledOutputWriter{}
	content := "header\n" + compilationconst.BuildIDPlaceholder + "\nfooter\n"
	if err := w.Write(path, content); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), compilationconst.BuildIDPlaceholder) {
		t.Error("placeholder should have been replaced")
	}
}

func TestWrite_EmptyContentVariant(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "empty.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, ""); err != nil {
		t.Fatalf("expected no error for empty content, got %v", err)
	}
	data, _ := os.ReadFile(path)
	if len(data) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(data))
	}
}

func TestWrite_AtomicOverwrite(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "atomic.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, "first\n"); err != nil {
		t.Fatal(err)
	}
	if err := w.Write(path, "second\n"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "second\n" {
		t.Errorf("expected second write to win, got %q", string(data))
	}
}

func TestWrite_NoBuildIDInSimpleContent(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "simple.md")
	w := &CompiledOutputWriter{}
	content := "no placeholders here\n"
	if err := w.Write(path, content); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != content {
		t.Errorf("expected content unchanged, got %q", string(data))
	}
}

func TestCompiledOutputWriter_ZeroValue(t *testing.T) {
	var w CompiledOutputWriter
	base := t.TempDir()
	path := filepath.Join(base, "zv.md")
	if err := w.Write(path, "test\n"); err != nil {
		t.Fatalf("zero-value writer should work: %v", err)
	}
}
