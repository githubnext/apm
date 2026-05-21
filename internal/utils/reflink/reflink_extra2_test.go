package reflink_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/reflink"
)

func TestCloneFile_SimpleCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	content := "hello world"
	if err := os.WriteFile(src, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != content {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestCloneFile_DisabledByEnvVar(t *testing.T) {
	t.Setenv("APM_NO_REFLINK", "1")
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	used, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile with NO_REFLINK: %v", err)
	}
	if used {
		t.Error("expected reflink=false when APM_NO_REFLINK=1")
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "data" {
		t.Errorf("content mismatch: got %q", got)
	}
}

func TestCloneFile_NonExistentSrc(t *testing.T) {
	dir := t.TempDir()
	_, err := reflink.CloneFile(filepath.Join(dir, "nosuchfile"), filepath.Join(dir, "dst"))
	if err == nil {
		t.Error("expected error for non-existent src")
	}
}

func TestCloneFile_BinaryContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	dst := filepath.Join(dir, "dst.bin")
	data := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
	if err := os.WriteFile(src, data, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile binary: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if len(got) != len(data) {
		t.Errorf("binary length mismatch: got %d, want %d", len(got), len(data))
	}
	for i, b := range data {
		if got[i] != b {
			t.Errorf("binary mismatch at byte %d: got %02x, want %02x", i, got[i], b)
		}
	}
}

func TestReflinkSupported_ReturnsBoolean(t *testing.T) {
	dir := t.TempDir()
	// Should not panic; just verify it returns a bool
	_ = reflink.ReflinkSupported(dir)
}

func TestCloneFile_MultilineTextContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	content := "line1\nline2\nline3\nline4\n"
	if err := os.WriteFile(src, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile multiline: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != content {
		t.Errorf("multiline mismatch: got %q, want %q", got, content)
	}
}

func TestCloneFile_DestinationOverwritten(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("new content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old content"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := reflink.CloneFile(src, dst)
	if err != nil {
		t.Fatalf("CloneFile overwrite: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "new content" {
		t.Errorf("expected overwritten content, got %q", got)
	}
}
