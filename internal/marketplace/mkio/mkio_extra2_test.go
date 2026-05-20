package mkio_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mkio"
)

func TestAtomicWrite_SameContentTwice(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "file.txt")
	content := []byte("hello world")
	if err := mkio.AtomicWrite(p, content); err != nil {
		t.Fatalf("first write: %v", err)
	}
	if err := mkio.AtomicWrite(p, content); err != nil {
		t.Fatalf("second write: %v", err)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got, content) {
		t.Errorf("expected %q got %q", content, got)
	}
}

func TestAtomicWrite_UpdateContent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "upd.txt")
	mkio.AtomicWrite(p, []byte("old"))
	if err := mkio.AtomicWrite(p, []byte("new")); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, _ := os.ReadFile(p)
	if string(got) != "new" {
		t.Errorf("expected 'new' got %q", got)
	}
}

func TestAtomicWriteString_Unicode(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "uni.txt")
	s := "hello\nworld\n"
	if err := mkio.AtomicWriteString(p, s); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, _ := os.ReadFile(p)
	if string(got) != s {
		t.Errorf("expected %q got %q", s, got)
	}
}

func TestAtomicWrite_NoBytesLost(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "data.bin")
	data := make([]byte, 4096)
	for i := range data {
		data[i] = byte(i % 256)
	}
	if err := mkio.AtomicWrite(p, data); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got, data) {
		t.Error("binary content mismatch")
	}
}

func TestAtomicWriteString_Newlines(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "nl.txt")
	s := "line1\nline2\nline3\n"
	mkio.AtomicWriteString(p, s)
	got, _ := os.ReadFile(p)
	if string(got) != s {
		t.Errorf("newlines not preserved: %q", got)
	}
}

func TestAtomicWrite_NestedSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b")
	os.MkdirAll(sub, 0o755)
	p := filepath.Join(sub, "file.txt")
	if err := mkio.AtomicWrite(p, []byte("nested")); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, _ := os.ReadFile(p)
	if string(got) != "nested" {
		t.Errorf("expected 'nested' got %q", got)
	}
}

func TestAtomicWriteString_ReturnNilOnSuccess(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ok.txt")
	if err := mkio.AtomicWriteString(p, "data"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestAtomicWrite_ReturnNilOnSuccess(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ok.bin")
	if err := mkio.AtomicWrite(p, []byte{0x01, 0x02}); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
