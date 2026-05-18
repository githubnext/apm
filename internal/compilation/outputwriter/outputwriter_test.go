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

func TestWrite_Overwrite(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "AGENTS.md")
w := &CompiledOutputWriter{}
if err := w.Write(path, "# first\n"); err != nil {
t.Fatalf("first write: %v", err)
}
if err := w.Write(path, "# second\n"); err != nil {
t.Fatalf("second write: %v", err)
}
data, _ := os.ReadFile(path)
if string(data) != "# second\n" {
t.Errorf("overwrite failed: got %q", string(data))
}
}

func TestWrite_EmptyContent(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "AGENTS.md")
w := &CompiledOutputWriter{}
if err := w.Write(path, ""); err != nil {
t.Fatalf("write empty: %v", err)
}
data, _ := os.ReadFile(path)
if string(data) != "" {
t.Errorf("expected empty file, got %q", string(data))
}
}

func TestWrite_DeepNestedPath(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "a", "b", "c", "AGENTS.md")
w := &CompiledOutputWriter{}
if err := w.Write(path, "nested"); err != nil {
t.Fatalf("deep nested write: %v", err)
}
if _, err := os.Stat(path); err != nil {
t.Errorf("nested file not found: %v", err)
}
}
