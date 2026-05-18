package lockfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeployedFileHash_nonexistent(t *testing.T) {
	got := DeployedFileHash("/nonexistent/path/file.txt")
	if got != "" {
		t.Errorf("expected empty string for nonexistent file, got %s", got)
	}
}

func TestDeployedFileHash_real(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := DeployedFileHash(f)
	if got == "" {
		t.Error("expected non-empty hash")
	}
	if len(got) < 7 || got[:7] != "sha256:" {
		t.Errorf("expected sha256: prefix, got %s", got)
	}
}

func TestDeployedFileHash_stable(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "stable.txt")
	if err := os.WriteFile(f, []byte("stable content"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1 := DeployedFileHash(f)
	h2 := DeployedFileHash(f)
	if h1 != h2 {
		t.Errorf("hash should be stable: %s vs %s", h1, h2)
	}
}

func TestDeployedFileHash_diffContent(t *testing.T) {
	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "a.txt")
	f2 := filepath.Join(tmp, "b.txt")
	os.WriteFile(f1, []byte("content a"), 0o644)
	os.WriteFile(f2, []byte("content b"), 0o644)
	h1 := DeployedFileHash(f1)
	h2 := DeployedFileHash(f2)
	if h1 == h2 {
		t.Error("different content should produce different hashes")
	}
}

func TestComputeDeployedHashes_skipMissing(t *testing.T) {
	tmp := t.TempDir()
	out := ComputeDeployedHashes(tmp, []string{"nonexistent.md", ""})
	if len(out) != 0 {
		t.Errorf("expected empty map for missing files, got %v", out)
	}
}

func TestComputeDeployedHashes_realFile(t *testing.T) {
	tmp := t.TempDir()
	rel := "foo/bar.md"
	abs := filepath.Join(tmp, rel)
	os.MkdirAll(filepath.Dir(abs), 0o755)
	os.WriteFile(abs, []byte("data"), 0o644)
	out := ComputeDeployedHashes(tmp, []string{rel})
	if _, ok := out[rel]; !ok {
		t.Errorf("expected hash for %s", rel)
	}
}

func TestSortedDeployedFiles_stable(t *testing.T) {
	files := []string{"z.md", "a.md", "m.md"}
	sorted := SortedDeployedFiles(files)
	if sorted[0] != "a.md" || sorted[1] != "m.md" || sorted[2] != "z.md" {
		t.Errorf("unexpected sort order: %v", sorted)
	}
}

func TestSortedDeployedFiles_empty(t *testing.T) {
	got := SortedDeployedFiles(nil)
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestSortedDeployedFiles_noMutate(t *testing.T) {
	orig := []string{"z.md", "a.md"}
	SortedDeployedFiles(orig)
	if orig[0] != "z.md" {
		t.Error("original slice should not be mutated")
	}
}

func TestWriteIfChanged_newFile(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "lock.yaml")
	changed, err := WriteIfChanged(p, []byte("content: 1\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for new file")
	}
}

func TestWriteIfChanged_sameContent(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "lock.yaml")
	content := []byte("content: 1\n")
	os.WriteFile(p, content, 0o644)
	changed, err := WriteIfChanged(p, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false for same content")
	}
}

func TestWriteIfChanged_differentContent(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "lock.yaml")
	os.WriteFile(p, []byte("old"), 0o644)
	changed, err := WriteIfChanged(p, []byte("new"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for different content")
	}
}
