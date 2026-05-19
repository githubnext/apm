package mkio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mkio"
)

func TestAtomicWrite_ContentIsCorrect(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "file.json")
	content := []byte(`{"key":"value"}`)
	if err := mkio.AtomicWrite(p, content); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestAtomicWriteString_ContentIsCorrect(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	if err := mkio.AtomicWriteString(p, "hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestAtomicWrite_BinaryContent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bin.dat")
	content := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	if err := mkio.AtomicWrite(p, content); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(content) {
		t.Errorf("expected %d bytes, got %d", len(content), len(got))
	}
	for i := range content {
		if got[i] != content[i] {
			t.Errorf("byte %d: got %02x want %02x", i, got[i], content[i])
		}
	}
}

func TestAtomicWrite_MultipleWrites(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "multi.txt")
	for _, s := range []string{"first", "second", "third"} {
		if err := mkio.AtomicWriteString(p, s); err != nil {
			t.Fatalf("write %q: unexpected error: %v", s, err)
		}
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "third" {
		t.Errorf("expected 'third', got %q", string(got))
	}
}

func TestAtomicWrite_NoTmpFileLeftAfterSuccess(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "data.json")
	if err := mkio.AtomicWrite(p, []byte("{}")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "data.json" {
			t.Errorf("unexpected file left in temp dir: %s", e.Name())
		}
	}
}

func TestAtomicWriteString_LargeContent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "large.txt")
	var sb []byte
	for i := 0; i < 10000; i++ {
		sb = append(sb, 'a'+byte(i%26))
	}
	if err := mkio.AtomicWrite(p, sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 10000 {
		t.Errorf("expected 10000 bytes, got %d", len(got))
	}
}

func TestAtomicWrite_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "perms.txt")
	if err := mkio.AtomicWrite(p, []byte("perm test")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	mode := info.Mode().Perm()
	if mode&0o400 == 0 {
		t.Errorf("expected file to be readable by owner, mode=%o", mode)
	}
}

func TestAtomicWriteString_EmptyStringCreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.txt")
	if err := mkio.AtomicWriteString(p, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}
}
