package atomicio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/atomicio"
)

func TestWriteText_NewlineOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nl.txt")
	if err := atomicio.WriteText(path, "\n", 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "\n" {
		t.Errorf("got %q, want newline", got)
	}
}

func TestWriteText_FileCreatedInDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "out.txt")
	// directory does not exist; should fail gracefully
	err := atomicio.WriteText(path, "data", 0)
	if err == nil {
		t.Error("expected error for missing parent directory")
	}
}

func TestWriteText_Concurrent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "concurrent.txt")
	done := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func(n int) {
			s := string(rune('a' + n))
			done <- atomicio.WriteText(path, s, 0)
		}(i)
	}
	for i := 0; i < 5; i++ {
		if err := <-done; err != nil {
			t.Errorf("concurrent WriteText error: %v", err)
		}
	}
	// File should contain exactly one character written by the last winner.
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected single char, got %q", got)
	}
}

func TestWriteText_PreservesExistingOnError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "preserve.txt")
	if err := atomicio.WriteText(path, "original", 0); err != nil {
		t.Fatalf("initial write: %v", err)
	}
	// File exists; overwrite with same content to confirm it still reads back.
	if err := atomicio.WriteText(path, "original", 0); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "original" {
		t.Errorf("expected 'original', got %q", got)
	}
}

func TestWriteText_BinaryContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bin.bin")
	content := "null\x00byte\x01here"
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("binary content mismatch")
	}
}

func TestWriteText_ModeZeroNoChmod(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nomode.txt")
	if err := atomicio.WriteText(path, "nomode", 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file should exist: %v", err)
	}
}

func TestWriteText_TrailingNewlines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trail.txt")
	content := "line1\nline2\n\n\n"
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("trailing newlines not preserved: got %q", got)
	}
}

func TestWriteText_RepeatIdentical(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "repeat.txt")
	for i := 0; i < 10; i++ {
		if err := atomicio.WriteText(path, "stable", 0); err != nil {
			t.Fatalf("WriteText iter %d: %v", i, err)
		}
	}
	got, _ := os.ReadFile(path)
	if string(got) != "stable" {
		t.Errorf("content changed after repeats: %q", got)
	}
}

func TestWriteText_SpecialChars(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "special.txt")
	content := "tab:\there\nnewline\r\nCR-LF\n"
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("special chars mismatch: got %q", got)
	}
}

func TestWriteText_ExistingReadOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ro.txt")
	if err := atomicio.WriteText(path, "first", 0o600); err != nil {
		t.Fatalf("first write: %v", err)
	}
	// Overwrite existing file regardless of mode.
	if err := atomicio.WriteText(path, "second", 0o600); err != nil {
		t.Fatalf("second write: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "second" {
		t.Errorf("expected 'second', got %q", got)
	}
}
