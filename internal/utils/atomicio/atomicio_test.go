package atomicio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/atomicio"
)

func TestWriteText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	if err := atomicio.WriteText(path, "hello world", 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestWriteTextOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	atomicio.WriteText(path, "first", 0)
	if err := atomicio.WriteText(path, "second", 0); err != nil {
		t.Fatalf("WriteText overwrite: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "second" {
		t.Errorf("got %q, want second", got)
	}
}

func TestWriteText_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := atomicio.WriteText(path, "", 0); err != nil {
		t.Fatalf("WriteText empty: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty file, got %q", got)
	}
}

func TestWriteText_Unicode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "unicode.txt")
	content := "hello world\nline two\n"

	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != content {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestWriteText_MultipleOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi.txt")

	for i, s := range []string{"a", "bb", "ccc", "dddd", "eeeee"} {
		if err := atomicio.WriteText(path, s, 0); err != nil {
			t.Fatalf("WriteText iteration %d: %v", i, err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != s {
			t.Errorf("iteration %d: got %q, want %q", i, got, s)
		}
	}
}

func TestWriteText_AtomicOnExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.txt")

	// Create initial file
	if err := os.WriteFile(path, []byte("original"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := atomicio.WriteText(path, "replaced", 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "replaced" {
		t.Errorf("got %q, want replaced", got)
	}
}

func TestWriteText_LargeContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.txt")

	// Build a 1 MB string
	chunk := "abcdefghijklmnopqrstuvwxyz0123456789\n"
	var sb []byte
	for len(sb) < 1<<20 {
		sb = append(sb, chunk...)
	}
	content := string(sb)

	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText large: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != content {
		t.Errorf("large content mismatch (lengths: got %d, want %d)", len(got), len(content))
	}
}

func TestWriteText_WithMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "withmode.txt")

	if err := atomicio.WriteText(path, "mode test", 0o600); err != nil {
		t.Fatalf("WriteText with mode: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "mode test" {
		t.Errorf("got %q, want 'mode test'", got)
	}
}
