package reflink_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/reflink"
)

func TestCloneFile_EmptyFileContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "empty.txt")
	dst := filepath.Join(dir, "empty-dst.txt")
	if err := os.WriteFile(src, []byte(""), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if len(got) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(got))
	}
}

func TestCloneFile_LargeFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "large.bin")
	dst := filepath.Join(dir, "large-dst.bin")
	data := make([]byte, 128*1024)
	for i := range data {
		data[i] = byte(i & 0xff)
	}
	if err := os.WriteFile(src, data, 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if len(got) != len(data) {
		t.Errorf("size mismatch: got %d, want %d", len(got), len(data))
	}
	for i := range data {
		if got[i] != data[i] {
			t.Errorf("data mismatch at byte %d", i)
			break
		}
	}
}

func TestCloneFile_OverwriteDst(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(dst, []byte("old"), 0o644); err != nil {
		t.Fatalf("write existing dst: %v", err)
	}
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if _, err := reflink.CloneFile(src, dst); err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "new" {
		t.Errorf("expected 'new', got %q", string(got))
	}
}

func TestCloneFile_ContentIntegrity(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	content := "line1\nline2\nline3\n"
	if err := os.WriteFile(src, []byte(content), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if _, err := reflink.CloneFile(src, dst); err != nil {
		t.Fatalf("CloneFile error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != content {
		t.Errorf("content mismatch: got %q, want %q", string(got), content)
	}
}

func TestNoReflinkEnvConst(t *testing.T) {
	if reflink.NoReflinkEnv == "" {
		t.Error("NoReflinkEnv constant should not be empty")
	}
}

func TestReflinkSupported_Path(t *testing.T) {
	// ReflinkSupported should not panic on a real path.
	dir := t.TempDir()
	_ = reflink.ReflinkSupported(dir)
}

func TestCloneFile_MultipleTimes(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	if err := os.WriteFile(src, []byte("multi"), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	for i := 0; i < 5; i++ {
		dst := filepath.Join(dir, "dst-"+string(rune('0'+i))+".txt")
		if _, err := reflink.CloneFile(src, dst); err != nil {
			t.Fatalf("CloneFile iter %d: %v", i, err)
		}
		got, _ := os.ReadFile(dst)
		if string(got) != "multi" {
			t.Errorf("iter %d content mismatch: %q", i, string(got))
		}
	}
}
