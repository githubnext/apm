package outputwriter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/outputwriter"
)

func TestWrite_SimpleContent_Extra4(t *testing.T) {
	dir := t.TempDir()
	w := &outputwriter.CompiledOutputWriter{}
	p := filepath.Join(dir, "out.md")
	if err := w.Write(p, "hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if string(data) != "hello world" {
		t.Fatalf("expected 'hello world', got %s", string(data))
	}
}

func TestWrite_NestedDir_Extra4(t *testing.T) {
	dir := t.TempDir()
	w := &outputwriter.CompiledOutputWriter{}
	p := filepath.Join(dir, "a", "b", "c", "out.txt")
	if err := w.Write(p, "nested"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestWrite_Overwrite_Extra4(t *testing.T) {
	dir := t.TempDir()
	w := &outputwriter.CompiledOutputWriter{}
	p := filepath.Join(dir, "file.txt")
	_ = w.Write(p, "first")
	if err := w.Write(p, "second"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if string(data) != "second" {
		t.Fatalf("expected 'second', got %s", string(data))
	}
}

func TestWrite_EmptyContent_Extra4(t *testing.T) {
	dir := t.TempDir()
	w := &outputwriter.CompiledOutputWriter{}
	p := filepath.Join(dir, "empty.txt")
	if err := w.Write(p, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if len(data) != 0 {
		t.Fatalf("expected empty file, got %d bytes", len(data))
	}
}

func TestWrite_WithNewlines_Extra4(t *testing.T) {
	dir := t.TempDir()
	w := &outputwriter.CompiledOutputWriter{}
	p := filepath.Join(dir, "multi.txt")
	content := "line1\nline2\nline3"
	if err := w.Write(p, content); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "line2") {
		t.Fatalf("expected line2 in output, got %s", string(data))
	}
}

func TestWrite_FileIsReadable_Extra4(t *testing.T) {
	dir := t.TempDir()
	w := &outputwriter.CompiledOutputWriter{}
	p := filepath.Join(dir, "readable.txt")
	_ = w.Write(p, "readable content")
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected non-zero file size")
	}
}

func TestCompiledOutputWriter_ZeroValue_Extra4(t *testing.T) {
	var w outputwriter.CompiledOutputWriter
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	if err := w.Write(p, "zero"); err != nil {
		t.Fatalf("unexpected error with zero-value writer: %v", err)
	}
}
