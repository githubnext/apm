package fileops_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/fileops"
)

func TestRobustRemoveAll(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "f.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "target")
	if err := os.Rename(sub, target); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustRemoveAll(target, false, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Error("directory should have been removed")
	}
}

func TestRobustCopyTree(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "dst")
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", data)
	}
}

func TestRobustCopy2(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fileops.RobustCopy2(src, dst, 0); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Errorf("expected 'content', got %q", data)
	}
}

func TestRobustRemoveAll_Nonexistent(t *testing.T) {
dir := t.TempDir()
nonexistent := dir + "/nonexistent_subdir"
// Should succeed even if path doesn't exist
err := fileops.RobustRemoveAll(nonexistent, false, 0)
if err != nil {
t.Errorf("expected no error for nonexistent path, got: %v", err)
}
}

func TestRobustRemoveAll_IgnoreErrors(t *testing.T) {
// Removing nonexistent with ignoreErrors=true should always succeed
err := fileops.RobustRemoveAll("/tmp/gh-aw/agent/nonexistent-xyz-123", true, 0)
if err != nil {
t.Errorf("ignoreErrors=true should suppress errors, got: %v", err)
}
}

func TestRobustCopyTree_NestedSubdirs(t *testing.T) {
src := t.TempDir()
dst := t.TempDir() + "/dst"
sub := src + "/subdir"
if err := os.MkdirAll(sub, 0o755); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(sub+"/nested.txt", []byte("nested"), 0o644); err != nil {
t.Fatal(err)
}
if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
t.Fatal(err)
}
data, err := os.ReadFile(dst + "/subdir/nested.txt")
if err != nil {
t.Fatal(err)
}
if string(data) != "nested" {
t.Errorf("expected 'nested', got %q", data)
}
}

func TestRobustCopyTree_MultipleFiles(t *testing.T) {
src := t.TempDir()
dst := t.TempDir() + "/dst2"
files := []string{"a.txt", "b.txt", "c.txt"}
for _, f := range files {
if err := os.WriteFile(src+"/"+f, []byte(f), 0o644); err != nil {
t.Fatal(err)
}
}
if err := fileops.RobustCopyTree(src, dst, false, false, 0); err != nil {
t.Fatal(err)
}
for _, f := range files {
data, err := os.ReadFile(dst + "/" + f)
if err != nil {
t.Fatalf("missing file %s: %v", f, err)
}
if string(data) != f {
t.Errorf("file %s: expected %q, got %q", f, f, data)
}
}
}

func TestRobustCopy2_OverwriteExisting(t *testing.T) {
dir := t.TempDir()
src := dir + "/src.txt"
dst := dir + "/dst.txt"
if err := os.WriteFile(dst, []byte("old"), 0o644); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
t.Fatal(err)
}
if err := fileops.RobustCopy2(src, dst, 0); err != nil {
t.Fatal(err)
}
data, err := os.ReadFile(dst)
if err != nil {
t.Fatal(err)
}
if string(data) != "new" {
t.Errorf("expected 'new', got %q", data)
}
}
