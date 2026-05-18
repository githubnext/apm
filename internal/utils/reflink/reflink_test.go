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

func TestCloneFile_MissingSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "nonexistent.txt")
	dst := filepath.Join(dir, "dst.txt")
	_, err := reflink.CloneFile(src, dst)
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestCloneFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "empty.txt")
	dst := filepath.Join(dir, "dst_empty.txt")
	if err := os.WriteFile(src, []byte(""), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile error on empty file: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(got))
	}
}

func TestCloneFile_LargeContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "large.txt")
	dst := filepath.Join(dir, "large_dst.txt")
	data := make([]byte, 64*1024)
	for i := range data {
		data[i] = byte(i % 251)
	}
	if err := os.WriteFile(src, data, 0o644); err != nil {
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
	if len(got) != len(data) {
		t.Errorf("size mismatch: got %d, want %d", len(got), len(data))
	}
}

func TestCloneFile_DisabledPreservesContent(t *testing.T) {
	t.Setenv(reflink.NoReflinkEnv, "1")
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	dst := filepath.Join(dir, "dst.bin")
	data := []byte("binary\x00content\xff")
	if err := os.WriteFile(src, data, 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if _, err := reflink.CloneFile(src, dst); err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(got) != string(data) {
		t.Error("content mismatch after fallback copy")
	}
}

func TestReflinkSupported_Normal(t *testing.T) {
	dir := t.TempDir()
	// Just verify it doesn't panic; result depends on filesystem
	_ = reflink.ReflinkSupported(dir)
}

func TestReflinkSupported_MissingDir(t *testing.T) {
	_ = reflink.ReflinkSupported("/nonexistent/path/for/reflink/test")
}
