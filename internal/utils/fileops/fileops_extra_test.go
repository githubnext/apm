package fileops_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/fileops"
)

func TestRobustRemoveAll_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty")
	if err := os.MkdirAll(empty, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustRemoveAll(empty, false, 0); err != nil {
		t.Fatalf("remove empty dir: %v", err)
	}
	if _, err := os.Stat(empty); !os.IsNotExist(err) {
		t.Error("empty dir should have been removed")
	}
}

func TestRobustRemoveAll_SingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lone.txt")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustRemoveAll(path, false, 0); err != nil {
		t.Fatalf("remove file: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should have been removed")
	}
}

func TestRobustCopyTree_Empty(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "dst-empty")
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatalf("copy empty tree: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("dst should exist after empty copy: %v", err)
	}
}

func TestRobustCopyTree_PreservesContent(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "dst-content")
	content := "hello fileops"
	if err := os.WriteFile(filepath.Join(src, "file.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dst, "file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content {
		t.Errorf("expected %q, got %q", content, string(got))
	}
}

func TestRobustCopy2_LargeFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "large.bin")
	dst := filepath.Join(dir, "large-dst.bin")
	data := strings.Repeat("abcdefghij", 10000)
	if err := os.WriteFile(src, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopy2(src, dst, 0); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != data {
		t.Errorf("large file content mismatch (lengths: src=%d dst=%d)", len(data), len(got))
	}
}

func TestRobustCopy2_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "empty.txt")
	dst := filepath.Join(dir, "empty-dst.txt")
	if err := os.WriteFile(src, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopy2(src, dst, 0); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(dst)
	if len(got) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(got))
	}
}

func TestRobustRemoveAll_MaxRetries(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "retrydir")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	// Higher maxRetries should still succeed on a normal dir
	if err := fileops.RobustRemoveAll(sub, false, 10); err != nil {
		t.Fatalf("remove with high retry count: %v", err)
	}
}

func TestRobustCopyTree_DeeplyNested(t *testing.T) {
	src := t.TempDir()
	deep := filepath.Join(src, "a", "b", "c", "d")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(deep, "leaf.txt"), []byte("leaf"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), "nested-dst")
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dst, "a", "b", "c", "d", "leaf.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "leaf" {
		t.Errorf("expected 'leaf', got %q", got)
	}
}

func TestRobustCopyTree_ManyFiles(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "many-dst")
	for i := 0; i < 20; i++ {
		name := filepath.Join(src, strings.Repeat("x", i+1)+".txt")
		if err := os.WriteFile(name, []byte(name), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 20 {
		t.Errorf("expected 20 files, got %d", len(entries))
	}
}
