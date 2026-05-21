package guards

import (
"os"
"path/filepath"
"testing"
)

func TestGuard_MultipleProtectedPaths_NoViolationExtra4(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, "a.txt"), []byte("aaa"), 0o644)
os.WriteFile(filepath.Join(dir, "b.txt"), []byte("bbb"), 0o644)
g := NewReadOnlyProjectGuard(dir, []string{"a.txt", "b.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := g.Exit(nil); err != nil {
t.Errorf("expected nil, got %v", err)
}
}

func TestGuard_DeletedFile_ReturnsViolationExtra4(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "del.txt")
os.WriteFile(path, []byte("data"), 0o644)
g := NewReadOnlyProjectGuard(dir, []string{"del.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
os.Remove(path)
if err := g.Exit(nil); err == nil {
t.Error("expected violation for deleted file")
}
}

func TestGuard_EmptyProtectedPaths_NoViolationExtra4(t *testing.T) {
dir := t.TempDir()
g := NewReadOnlyProjectGuard(dir, []string{})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := g.Exit(nil); err != nil {
t.Errorf("expected nil for empty paths, got %v", err)
}
}

func TestGuard_OrigErrSuppressesViolationExtra4(t *testing.T) {
// When origErr != nil, violations are suppressed (origErr wins).
dir := t.TempDir()
path := filepath.Join(dir, "file.txt")
os.WriteFile(path, []byte("before"), 0o644)
g := NewReadOnlyProjectGuard(dir, []string{"file.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
os.WriteFile(path, []byte("after longer content here"), 0o644)
origErr := &ProtectedPathMutationError{Violations: []string{"sentinel"}}
err := g.Exit(origErr)
// With origErr set, Exit returns nil (violations suppressed)
if err != nil {
t.Errorf("expected nil when origErr is set, got %v", err)
}
}

func TestGuard_NonExistentProtectedPath_NoErrorExtra4(t *testing.T) {
dir := t.TempDir()
g := NewReadOnlyProjectGuard(dir, []string{"nonexistent.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := g.Exit(nil); err != nil {
t.Errorf("expected nil for nonexistent path, got %v", err)
}
}

func TestProtectedPathMutationError_MessageNonEmptyExtra4(t *testing.T) {
e := &ProtectedPathMutationError{Violations: []string{"/some/path/file.go"}}
msg := e.Error()
if msg == "" {
t.Error("expected non-empty error message")
}
}
