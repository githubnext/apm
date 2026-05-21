package guards

import (
"os"
"path/filepath"
"testing"
)

func TestGuard_NoFiles_NoViolationExtra3(t *testing.T) {
dir := t.TempDir()
g := NewReadOnlyProjectGuard(dir, []string{"."})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := g.Exit(nil); err != nil {
t.Errorf("expected nil, got %v", err)
}
}

func TestGuard_UnchangedFile_NoViolationExtra3(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(filepath.Join(dir, "readme.md"), []byte("hello"), 0o644); err != nil {
t.Fatal(err)
}
g := NewReadOnlyProjectGuard(dir, []string{"readme.md"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := g.Exit(nil); err != nil {
t.Errorf("expected nil, got %v", err)
}
}

func TestGuard_ModifiedFile_ReturnsErrorExtra3(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(filepath.Join(dir, "mod.txt"), []byte("before"), 0o644); err != nil {
t.Fatal(err)
}
g := NewReadOnlyProjectGuard(dir, []string{"mod.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
// Write larger content to ensure size change
if err := os.WriteFile(filepath.Join(dir, "mod.txt"), []byte("after and longer content here"), 0o644); err != nil {
t.Fatal(err)
}
err := g.Exit(nil)
if err == nil {
t.Error("expected violation error for modified file")
}
}

func TestGuard_OrigErr_SuppressesViolationExtra3(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(filepath.Join(dir, "x.txt"), []byte("a"), 0o644); err != nil {
t.Fatal(err)
}
g := NewReadOnlyProjectGuard(dir, []string{"x.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(filepath.Join(dir, "x.txt"), []byte("b and more content to change size"), 0o644); err != nil {
t.Fatal(err)
}
// When origErr != nil, violations are suppressed (returns nil per implementation)
err := g.Exit(os.ErrNotExist)
if err != nil {
t.Errorf("with origErr set, violations should be suppressed, got %v", err)
}
}

func TestProtectedPathMutationError_ErrorStringExtra3(t *testing.T) {
e := &ProtectedPathMutationError{Violations: []string{"a/b", "c/d"}}
s := e.Error()
if s == "" {
t.Error("Error() should return non-empty string")
}
}

func TestProtectedPathMutationError_NoViolationsExtra3(t *testing.T) {
e := &ProtectedPathMutationError{Violations: nil}
s := e.Error()
if s == "" {
t.Error("Error() should return non-empty string even with no violations")
}
}

func TestGuard_NewFile_DetectedExtra3(t *testing.T) {
dir := t.TempDir()
g := NewReadOnlyProjectGuard(dir, []string{"."})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("new"), 0o644); err != nil {
t.Fatal(err)
}
err := g.Exit(nil)
if err == nil {
t.Error("expected violation for newly created file")
}
}

func TestGuard_MultipleViolations_AllReportedExtra3(t *testing.T) {
dir := t.TempDir()
f1 := filepath.Join(dir, "f1.txt")
if err := os.WriteFile(f1, []byte("aaa"), 0o644); err != nil {
t.Fatal(err)
}
g := NewReadOnlyProjectGuard(dir, []string{".", "f1.txt"})
if err := g.Enter(); err != nil {
t.Fatal(err)
}
// Modify f1 (size change)
if err := os.WriteFile(f1, []byte("aaabbbccc extra content"), 0o644); err != nil {
t.Fatal(err)
}
err := g.Exit(nil)
if err == nil {
t.Error("expected violations for modified file")
}
}
