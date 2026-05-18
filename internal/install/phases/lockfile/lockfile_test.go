package lockfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeployedFileHash(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := DeployedFileHash(f)
	if h == "" {
		t.Error("expected non-empty hash")
	}
	if len(h) < 8 || h[:7] != "sha256:" {
		t.Errorf("expected 'sha256:' prefix, got %q", h)
	}
}

func TestDeployedFileHash_Missing(t *testing.T) {
	h := DeployedFileHash("/nonexistent/file.txt")
	if h != "" {
		t.Errorf("expected empty hash for missing file, got %q", h)
	}
}

func TestComputeDeployedHashes(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := ComputeDeployedHashes(dir, []string{"a.txt", "missing.txt", ""})
	h, ok := result["a.txt"]
	if !ok || h == "" {
		t.Error("expected hash for a.txt")
	}
	if _, ok := result["missing.txt"]; ok {
		t.Error("missing file should not appear in result")
	}
}

func TestWriteIfChanged_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lock.yaml")
	changed, err := WriteIfChanged(path, []byte("content"))
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Error("expected changed=true for new file")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "content" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestWriteIfChanged_SameContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lock.yaml")
	if err := os.WriteFile(path, []byte("same"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := WriteIfChanged(path, []byte("same"))
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Error("expected changed=false when content identical")
	}
}

func TestWriteIfChanged_DifferentContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lock.yaml")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := WriteIfChanged(path, []byte("new"))
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Error("expected changed=true when content differs")
	}
}

func TestSortedDeployedFiles(t *testing.T) {
	files := []string{"c.txt", "a.txt", "b.txt"}
	sorted := SortedDeployedFiles(files)
	if sorted[0] != "a.txt" || sorted[1] != "b.txt" || sorted[2] != "c.txt" {
		t.Errorf("unexpected order: %v", sorted)
	}
	// Ensure original slice is not mutated
	if files[0] != "c.txt" {
		t.Error("original slice should not be mutated")
	}
}

func TestSortedDeployedFiles_Empty(t *testing.T) {
	result := SortedDeployedFiles(nil)
	if len(result) != 0 {
		t.Errorf("want empty, got %v", result)
	}
}
