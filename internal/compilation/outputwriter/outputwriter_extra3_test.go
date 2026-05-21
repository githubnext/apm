package outputwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_NewFileCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, "hello\n"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(path)
	if string(b) != "hello\n" {
		t.Fatalf("unexpected content: %q", string(b))
	}
}

func TestWrite_NestedDirsCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "out.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, "content"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestWrite_ContentPreservedVerbatim(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.md")
	content := "line1\nline2\nline3\n"
	w := &CompiledOutputWriter{}
	_ = w.Write(path, content)
	b, _ := os.ReadFile(path)
	if string(b) != content {
		t.Fatalf("content mismatch: %q vs %q", string(b), content)
	}
}

func TestWrite_OverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.md")
	w := &CompiledOutputWriter{}
	_ = w.Write(path, "first")
	_ = w.Write(path, "second")
	b, _ := os.ReadFile(path)
	if string(b) != "second" {
		t.Fatalf("expected 'second', got %q", string(b))
	}
}

func TestWrite_EmptyStringSucceeds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.md")
	w := &CompiledOutputWriter{}
	if err := w.Write(path, ""); err != nil {
		t.Fatalf("unexpected error writing empty: %v", err)
	}
}

func TestWrite_MultilineContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi.md")
	w := &CompiledOutputWriter{}
	lines := strings.Repeat("x\n", 100)
	if err := w.Write(path, lines); err != nil {
		t.Fatalf("error: %v", err)
	}
	b, _ := os.ReadFile(path)
	if strings.Count(string(b), "x") != 100 {
		t.Fatalf("content count mismatch")
	}
}

func TestWrite_ZeroValueWriter(t *testing.T) {
	var w CompiledOutputWriter
	dir := t.TempDir()
	path := filepath.Join(dir, "zero.md")
	if err := w.Write(path, "ok"); err != nil {
		t.Fatalf("zero value writer error: %v", err)
	}
}

func TestWrite_SpecialCharContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "special.md")
	w := &CompiledOutputWriter{}
	content := "tab:\there\nnewline\n"
	_ = w.Write(path, content)
	b, _ := os.ReadFile(path)
	if !strings.Contains(string(b), "tab:") {
		t.Fatalf("special chars not preserved")
	}
}

func TestWrite_LargeContentSucceeds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.md")
	w := &CompiledOutputWriter{}
	content := strings.Repeat("a", 100000)
	if err := w.Write(path, content); err != nil {
		t.Fatalf("error: %v", err)
	}
}
