package reflink_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/reflink"
)

func TestCloneFile_BasicCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	content := "hello reflink test"
	if err := os.WriteFile(src, []byte(content), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(got) != content {
		t.Errorf("content mismatch: got %q, want %q", string(got), content)
	}
}

func TestCloneFile_DisabledByEnv(t *testing.T) {
	t.Setenv(reflink.NoReflinkEnv, "1")
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	reflinkUsed, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	if reflinkUsed {
		t.Error("expected reflink to be disabled by env var")
	}
}

func TestReflinkSupported_DisabledByEnv(t *testing.T) {
	t.Setenv(reflink.NoReflinkEnv, "1")
	dir := t.TempDir()
	if reflink.ReflinkSupported(dir) {
		t.Error("ReflinkSupported should return false when env disabled")
	}
}

func TestCloneFile_CreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "nested", "deep", "dst.txt")
	if err := os.WriteFile(src, []byte("content"), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if _, err := reflink.CloneFile(src, dst); err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("dst not created: %v", err)
	}
}
