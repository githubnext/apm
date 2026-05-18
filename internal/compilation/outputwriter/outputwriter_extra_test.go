package outputwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_ContentPreserved(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.md")
	w := &CompiledOutputWriter{}
	content := "# Heading\n\nParagraph text.\n"
	if err := w.Write(path, content); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "Heading") {
		t.Errorf("expected content preserved, got %q", string(data))
	}
}

func TestWrite_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	w := &CompiledOutputWriter{}
	for i, name := range []string{"a.md", "b.md", "c.md"} {
		content := strings.Repeat("line\n", i+1)
		if err := w.Write(filepath.Join(dir, name), content); err != nil {
			t.Fatalf("write %s failed: %v", name, err)
		}
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 3 {
		t.Errorf("expected 3 files, got %d", len(entries))
	}
}

func TestWrite_LargeContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.md")
	w := &CompiledOutputWriter{}
	content := strings.Repeat("a", 100*1024) // 100 KB
	if err := w.Write(path, content); err != nil {
		t.Fatalf("large write failed: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("expected non-empty large file")
	}
}

func TestWrite_SpecialChars(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "special.md")
	w := &CompiledOutputWriter{}
	content := "line1\nline2\ttabbed\n"
	if err := w.Write(path, content); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != content {
		t.Errorf("content mismatch: got %q want %q", string(data), content)
	}
}

func TestWrite_NewInstance(t *testing.T) {
	// Each call to Write with a fresh writer struct must work independently.
	dir := t.TempDir()
	for i := range []int{0, 1, 2} {
		_ = i
		w := &CompiledOutputWriter{}
		path := filepath.Join(dir, "file.md")
		if err := w.Write(path, "content"); err != nil {
			t.Fatalf("write failed: %v", err)
		}
	}
}
