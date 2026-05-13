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
