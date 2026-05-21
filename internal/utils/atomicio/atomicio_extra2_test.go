package atomicio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/atomicio"
)

func TestWriteText_OverwritesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.txt")
	if err := os.WriteFile(path, []byte("old content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := atomicio.WriteText(path, "new content", 0); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "new content" {
		t.Errorf("got %q, want %q", got, "new content")
	}
}

func TestWriteText_EmptyString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := atomicio.WriteText(path, "", 0); err != nil {
		t.Fatalf("WriteText empty: %v", err)
	}
	got, _ := os.ReadFile(path)
	if len(got) != 0 {
		t.Errorf("expected empty file, got %q", got)
	}
}

func TestWriteText_UnicodeContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "uni.txt")
	content := "hello \u00e9\u00e0 world"
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText unicode: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestWriteText_MultilineContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi.txt")
	content := "line1\nline2\nline3\n"
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText multiline: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestWriteText_WithFileModeNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mode.txt")
	if err := atomicio.WriteText(path, "data", 0o600); err != nil {
		t.Fatalf("WriteText with mode: %v", err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if fi.Mode()&0o777 != 0o600 {
		t.Errorf("expected mode 0600, got %o", fi.Mode()&0o777)
	}
}

func TestWriteText_VeryLargeContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.txt")
	content := ""
	for i := 0; i < 1000; i++ {
		content += "line of content with some text padding here\n"
	}
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText large: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("large content mismatch: got len %d, want len %d", len(got), len(content))
	}
}

func TestWriteText_SpecialCharacters(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "special.txt")
	content := "tab:\there\nnull-like: \\0\n"
	if err := atomicio.WriteText(path, content, 0); err != nil {
		t.Fatalf("WriteText special chars: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != content {
		t.Errorf("got %q, want %q", got, content)
	}
}
