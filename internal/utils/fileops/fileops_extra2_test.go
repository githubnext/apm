package fileops_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/fileops"
)

func TestRobustCopy2_BasicCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopy2(src, dst, 0); err != nil {
		t.Fatalf("RobustCopy2: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestRobustCopy2_NonExistentSrc(t *testing.T) {
	dir := t.TempDir()
	err := fileops.RobustCopy2(filepath.Join(dir, "nosrc.txt"), filepath.Join(dir, "dst.txt"), 0)
	if err == nil {
		t.Error("expected error for non-existent src")
	}
}

func TestRobustCopy2_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopy2(src, dst, 0); err != nil {
		t.Fatalf("RobustCopy2: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "new" {
		t.Errorf("got %q, want %q", got, "new")
	}
}

func TestRobustRemoveAll_NonExistentPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does_not_exist")
	// Should succeed (like os.RemoveAll) for non-existent paths
	err := fileops.RobustRemoveAll(path, false, 0)
	if err != nil {
		t.Errorf("expected no error for non-existent path, got: %v", err)
	}
}

func TestRobustRemoveAll_NestedTree(t *testing.T) {
	dir := t.TempDir()
	tree := filepath.Join(dir, "tree")
	sub := filepath.Join(tree, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(sub, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := fileops.RobustRemoveAll(tree, false, 0); err != nil {
		t.Fatalf("RobustRemoveAll tree: %v", err)
	}
	if _, err := os.Stat(tree); !os.IsNotExist(err) {
		t.Error("tree should have been removed")
	}
}

func TestRobustCopyTree_SingleFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "srcdir")
	dst := filepath.Join(dir, "dstdir")
	if err := os.Mkdir(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "file.txt"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatalf("RobustCopyTree: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dst, "file.txt"))
	if string(got) != "content" {
		t.Errorf("got %q, want %q", got, "content")
	}
}

func TestRobustCopyTree_NestedDirs(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	sub := filepath.Join(src, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "deep.txt"), []byte("deep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatalf("RobustCopyTree nested: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dst, "a", "b", "deep.txt"))
	if string(got) != "deep" {
		t.Errorf("got %q, want %q", got, "deep")
	}
}
