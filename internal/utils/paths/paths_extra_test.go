package paths

import (
"path/filepath"
"strings"
"testing"
)

func TestPortableRelpath_AbsoluteFallback(t *testing.T) {
// When path is not under base we still get a non-empty forward-slash string.
got := PortableRelpath("/a/b/c", "/x/y/z")
if got == "" {
t.Error("PortableRelpath with disjoint paths should not return empty string")
}
if strings.ContainsRune(got, '\\') {
t.Errorf("result contains backslash: %q", got)
}
}

func TestPortableRelpath_SingleComponent(t *testing.T) {
tmpDir := t.TempDir()
child := filepath.Join(tmpDir, "file.go")
got := PortableRelpath(child, tmpDir)
if got != "file.go" {
t.Errorf("got %q, want file.go", got)
}
}

func TestPortableRelpath_TrailingSlashBase(t *testing.T) {
tmpDir := t.TempDir()
child := filepath.Join(tmpDir, "sub", "x.py")
// base with trailing separator — filepath.Abs cleans it.
got := PortableRelpath(child, tmpDir+string(filepath.Separator))
if strings.ContainsRune(got, '\\') {
t.Errorf("result contains backslash: %q", got)
}
if !strings.HasSuffix(got, "x.py") {
t.Errorf("expected result to end with x.py, got %q", got)
}
}

func TestPortableRelpath_MultiLevelReturn(t *testing.T) {
tmpDir := t.TempDir()
child := filepath.Join(tmpDir, "a", "b")
// path is tmpDir, base is child -- should traverse up
got := PortableRelpath(tmpDir, child)
if strings.ContainsRune(got, '\\') {
t.Errorf("result contains backslash: %q", got)
}
if !strings.HasPrefix(got, "..") {
t.Errorf("expected relative upward traversal, got %q", got)
}
}

func TestPortableRelpath_HiddenFile(t *testing.T) {
tmpDir := t.TempDir()
hidden := filepath.Join(tmpDir, ".hidden", "secret.txt")
got := PortableRelpath(hidden, tmpDir)
want := ".hidden/secret.txt"
if got != want {
t.Errorf("got %q, want %q", got, want)
}
}

func TestPortableRelpath_LongPath(t *testing.T) {
tmpDir := t.TempDir()
deep := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f", "g.txt")
got := PortableRelpath(deep, tmpDir)
want := "a/b/c/d/e/f/g.txt"
if got != want {
t.Errorf("got %q, want %q", got, want)
}
}

func TestPortableRelpath_SiblingDir(t *testing.T) {
tmpDir := t.TempDir()
// path is in sibling dir relative to base's parent
baseDir := filepath.Join(tmpDir, "base")
otherDir := filepath.Join(tmpDir, "other", "file.txt")
got := PortableRelpath(otherDir, baseDir)
if strings.ContainsRune(got, '\\') {
t.Errorf("result contains backslash: %q", got)
}
}

func TestPortableRelpath_DotExtension(t *testing.T) {
tmpDir := t.TempDir()
child := filepath.Join(tmpDir, ".env")
got := PortableRelpath(child, tmpDir)
if got != ".env" {
t.Errorf("got %q, want .env", got)
}
}

func TestPortableRelpath_ReturnsSameForwardSlash(t *testing.T) {
// Calling twice returns the same value.
tmpDir := t.TempDir()
child := filepath.Join(tmpDir, "x", "y.go")
got1 := PortableRelpath(child, tmpDir)
got2 := PortableRelpath(child, tmpDir)
if got1 != got2 {
t.Errorf("PortableRelpath not deterministic: %q vs %q", got1, got2)
}
}

func TestPortableRelpath_WindowsBackslashInInput(t *testing.T) {
// Input with backslashes in the path component should be cleaned up by Abs.
tmpDir := t.TempDir()
child := filepath.Join(tmpDir, "sub", "file.txt")
// Convert to backslashes to simulate Windows-style input on Linux.
childWin := strings.ReplaceAll(child, "/", "\\")
// On Linux filepath.Abs won't resolve backslash paths the same way, but the
// function should still return without panicking.
got := PortableRelpath(childWin, tmpDir)
_ = got // just verify no panic
}
