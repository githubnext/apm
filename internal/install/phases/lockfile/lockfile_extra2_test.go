package lockfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSortedDeployedFiles_AlreadySorted(t *testing.T) {
	files := []string{"a.txt", "b.txt", "c.txt"}
	got := SortedDeployedFiles(files)
	for i, f := range got {
		if f != files[i] {
			t.Errorf("index %d: expected %q, got %q", i, files[i], f)
		}
	}
}

func TestSortedDeployedFiles_UnsortedInput(t *testing.T) {
	files := []string{"z.txt", "a.txt", "m.txt"}
	got := SortedDeployedFiles(files)
	if got[0] != "a.txt" || got[1] != "m.txt" || got[2] != "z.txt" {
		t.Errorf("unexpected sort order: %v", got)
	}
}

func TestSortedDeployedFiles_EmptySlice(t *testing.T) {
	got := SortedDeployedFiles(nil)
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestSortedDeployedFiles_DoesNotMutateInput(t *testing.T) {
	files := []string{"z.txt", "a.txt"}
	SortedDeployedFiles(files)
	if files[0] != "z.txt" {
		t.Error("SortedDeployedFiles should not mutate the input slice")
	}
}

func TestWriteIfChanged_WritesNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lock.yaml")
	changed, err := WriteIfChanged(path, []byte("content: v1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for new file")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "content: v1" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestWriteIfChanged_NoWriteWhenSameContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lock.yaml")
	_ = os.WriteFile(path, []byte("same"), 0o644)
	changed, err := WriteIfChanged(path, []byte("same"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false when content is identical")
	}
}

func TestWriteIfChanged_WritesWhenContentDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lock.yaml")
	_ = os.WriteFile(path, []byte("old"), 0o644)
	changed, err := WriteIfChanged(path, []byte("new"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true when content differs")
	}
}

func TestLockfileEntry_ZeroValue(t *testing.T) {
	var e LockfileEntry
	if e.DepKey != "" || e.RepoURL != "" || e.ContentHash != "" {
		t.Error("zero-value LockfileEntry fields should be empty strings")
	}
	if e.DeployedFiles != nil {
		t.Error("DeployedFiles should be nil in zero-value")
	}
}

func TestComputeDeployedHashes_SkipsEmptyPaths(t *testing.T) {
	dir := t.TempDir()
	result := ComputeDeployedHashes(dir, []string{"", ""})
	if len(result) != 0 {
		t.Errorf("expected empty map for empty paths, got %v", result)
	}
}
