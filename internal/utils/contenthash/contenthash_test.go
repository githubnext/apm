package contenthash_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/contenthash"
)

func TestComputePackageHash_empty(t *testing.T) {
	dir := t.TempDir()
	h, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := "sha256:" + fmt.Sprintf("%x", sha256.Sum256([]byte{}))
	if h != want {
		t.Errorf("empty dir: got %s, want %s", h, want)
	}
}

func TestComputePackageHash_nonexistent(t *testing.T) {
	h, err := contenthash.ComputePackageHash("/nonexistent/path/xyz")
	if err != nil {
		t.Fatal(err)
	}
	want := "sha256:" + fmt.Sprintf("%x", sha256.Sum256([]byte{}))
	if h != want {
		t.Errorf("nonexistent: got %s, want %s", h, want)
	}
}

func TestComputePackageHash_deterministic(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := contenthash.ComputePackageHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Errorf("not deterministic: %s != %s", h1, h2)
	}
}

func TestComputePackageHash_excludesMarker(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, _ := contenthash.ComputePackageHash(dir)

	if err := os.WriteFile(filepath.Join(dir, ".apm-pin"), []byte("marker"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, _ := contenthash.ComputePackageHash(dir)
	if h1 != h2 {
		t.Errorf("marker should not affect hash: %s vs %s", h1, h2)
	}
}

func TestComputeFileHash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	h, err := contenthash.ComputeFileHash(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte("content"))
	want := fmt.Sprintf("sha256:%x", sum)
	if h != want {
		t.Errorf("got %s, want %s", h, want)
	}
}

func TestVerifyPackageHash(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	h, _ := contenthash.ComputePackageHash(dir)
	ok, err := contenthash.VerifyPackageHash(dir, h)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected hash to verify")
	}
	ok2, _ := contenthash.VerifyPackageHash(dir, "sha256:wrong")
	if ok2 {
		t.Error("expected mismatch to fail")
	}
}
