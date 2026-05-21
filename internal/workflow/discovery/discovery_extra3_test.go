package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverWorkflows_EmptyDirectory_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	wfs, errs := DiscoverWorkflows(dir)
	if len(wfs) != 0 {
		t.Errorf("expected 0 workflows, got %d", len(wfs))
	}
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d", len(errs))
	}
}

func TestDiscoverWorkflows_SinglePromptMdFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "greet.prompt.md")
	content := "---\ndescription: greet\n---\nHello\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	wfs, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(wfs))
	}
	if wfs[0].Name != "greet" {
		t.Errorf("expected name 'greet', got %q", wfs[0].Name)
	}
}

func TestDiscoverWorkflows_TwoFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.prompt.md", "b.prompt.md"} {
		f := filepath.Join(dir, name)
		if err := os.WriteFile(f, []byte("---\ndescription: x\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	wfs, _ := DiscoverWorkflows(dir)
	if len(wfs) != 2 {
		t.Errorf("expected 2 workflows, got %d", len(wfs))
	}
}

func TestDiscoverWorkflows_IgnoresPlainMdFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "readme.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	wfs, _ := DiscoverWorkflows(dir)
	if len(wfs) != 0 {
		t.Errorf("expected 0 workflows, got %d", len(wfs))
	}
}

func TestDiscoverWorkflows_NestedSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(sub, "nested.prompt.md")
	if err := os.WriteFile(f, []byte("---\ndescription: nested\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	wfs, _ := DiscoverWorkflows(dir)
	if len(wfs) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(wfs))
	}
}

func TestDiscoverWorkflows_NonExistentDirVariant(t *testing.T) {
	wfs, _ := DiscoverWorkflows("/nonexistent/dir/xyz")
	// Should return empty, not panic
	if len(wfs) != 0 {
		t.Errorf("expected 0, got %d", len(wfs))
	}
}

func TestDiscoverWorkflows_FilePathSet(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "myflow.prompt.md")
	if err := os.WriteFile(f, []byte("---\ndescription: x\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	wfs, _ := DiscoverWorkflows(dir)
	if len(wfs) != 1 {
		t.Fatalf("expected 1, got %d", len(wfs))
	}
	if wfs[0].FilePath != f {
		t.Errorf("FilePath mismatch: got %q, want %q", wfs[0].FilePath, f)
	}
}
