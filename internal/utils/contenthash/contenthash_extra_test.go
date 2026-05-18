package contenthash_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/contenthash"
)

func TestComputePackageHash_FileChange(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(fp, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fp, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Error("hash should change when file content changes")
	}
}

func TestComputePackageHash_SubdirIncluded(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	h1, _ := contenthash.ComputePackageHash(dir)
	if err := os.WriteFile(filepath.Join(subdir, "f.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Error("hash should differ when subdir file is added")
	}
}

func TestComputePackageHash_GitDirExcluded(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "code.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, _ := contenthash.ComputePackageHash(dir)
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Error(".git directory should be excluded from hash")
	}
}

func TestComputePackageHash_StartsWithSha256(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "x.go"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	h, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(h) < 7 || h[:7] != "sha256:" {
		t.Errorf("hash should start with 'sha256:', got %q", h)
	}
}

func TestComputeFileHash_MissingFile(t *testing.T) {
	h, err := contenthash.ComputeFileHash("/nonexistent/path/file.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "sha256:" + fmt.Sprintf("%x", sha256.Sum256([]byte{}))
	if h != want {
		t.Errorf("missing file: got %s, want %s", h, want)
	}
}

func TestComputeFileHash_Directory(t *testing.T) {
	dir := t.TempDir()
	h, err := contenthash.ComputeFileHash(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "sha256:" + fmt.Sprintf("%x", sha256.Sum256([]byte{}))
	if h != want {
		t.Errorf("directory: got %s, want %s", h, want)
	}
}

func TestVerifyPackageHash_Mismatch(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := contenthash.VerifyPackageHash(dir, "sha256:badhash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for wrong hash")
	}
}

func TestComputePackageHash_PycacheExcluded(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "mod.py"), []byte("x=1"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, _ := contenthash.ComputePackageHash(dir)
	pycDir := filepath.Join(dir, "__pycache__")
	if err := os.MkdirAll(pycDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pycDir, "mod.pyc"), []byte("bytecode"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Error("__pycache__ should be excluded from hash")
	}
}
