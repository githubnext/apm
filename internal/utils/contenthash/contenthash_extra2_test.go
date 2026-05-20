package contenthash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeFileHash_ValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "testfile.txt")
	if err := os.WriteFile(p, []byte("test content"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := ComputeFileHash(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
	if !strings.HasPrefix(hash, "sha256:") {
		t.Errorf("hash should start with sha256:, got %q", hash)
	}
}

func TestComputeFileHash_Deterministic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "det.txt")
	if err := os.WriteFile(p, []byte("deterministic content"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash1, err := ComputeFileHash(p)
	if err != nil {
		t.Fatal(err)
	}
	hash2, err := ComputeFileHash(p)
	if err != nil {
		t.Fatal(err)
	}
	if hash1 != hash2 {
		t.Errorf("hashes differ: %q vs %q", hash1, hash2)
	}
}

func TestComputeFileHash_DifferentContent(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "f1.txt")
	p2 := filepath.Join(dir, "f2.txt")
	if err := os.WriteFile(p1, []byte("content A"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, []byte("content B"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, _ := ComputeFileHash(p1)
	h2, _ := ComputeFileHash(p2)
	if h1 == h2 {
		t.Error("different files should have different hashes")
	}
}

func TestVerifyPackageHash_CorrectHash(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "pkg.txt")
	if err := os.WriteFile(p, []byte("pkg content"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := VerifyPackageHash(dir, hash)
	if err != nil {
		t.Fatalf("VerifyPackageHash error: %v", err)
	}
	if !ok {
		t.Error("expected hash verification to pass")
	}
}

func TestVerifyPackageHash_WrongHash(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "pkg.txt")
	if err := os.WriteFile(p, []byte("pkg content"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := VerifyPackageHash(dir, "sha256:wronghashvalue")
	if err != nil {
		t.Fatalf("VerifyPackageHash error: %v", err)
	}
	if ok {
		t.Error("wrong hash should not verify")
	}
}

func TestComputePackageHash_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	hash, err := ComputePackageHash(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash even for empty dir")
	}
}

func TestComputePackageHash_DotGitExcluded(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "source.txt"), []byte("source content"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, err := ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	// Changing the .git file shouldn't change the hash
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/other"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Error("hash should not change when only .git files change")
	}
}

func TestComputePackageHash_FileAdded(t *testing.T) {
	dir := t.TempDir()
	h1, err := ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("new file"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Error("adding a file should change the hash")
	}
}

func TestComputeFileHash_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(p, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := ComputeFileHash(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("empty file should still produce a hash")
	}
}
