package guards_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/guards"
)

func TestProtectedPathMutationError_Message(t *testing.T) {
	e := &guards.ProtectedPathMutationError{Violations: []string{"modified: /tmp/a", "deleted: /tmp/b"}}
	msg := e.Error()
	if !strings.Contains(msg, "modified: /tmp/a") {
		t.Errorf("expected violation in error message: %q", msg)
	}
	if !strings.Contains(msg, "deleted: /tmp/b") {
		t.Errorf("expected second violation in error message: %q", msg)
	}
}

func TestProtectedPathMutationError_SingleViolation(t *testing.T) {
	e := &guards.ProtectedPathMutationError{Violations: []string{"created: /tmp/x"}}
	if !strings.Contains(e.Error(), "created: /tmp/x") {
		t.Errorf("unexpected message: %q", e.Error())
	}
}

func TestGuard_NestedDirectoryDetection(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(sub, "nested.txt")
	if err := os.WriteFile(f, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"sub"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// Modify nested file
	if err := os.WriteFile(f, []byte("changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected error for nested file modification")
	}
}

func TestGuard_MultipleProtectedRoots(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(f1, []byte("aaa"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("bbb"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"a.txt", "b.txt"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// Modify only b.txt
	if err := os.WriteFile(f2, []byte("changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected error for modified b.txt")
	}
}

func TestGuard_UnmodifiedNestedDir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(sub, "file.txt")
	if err := os.WriteFile(f, []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"sub"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// No changes
	if err := g.Exit(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGuard_EmptyProtectedRoots(t *testing.T) {
	dir := t.TempDir()
	g := guards.NewReadOnlyProjectGuard(dir, []string{})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// Create a file - should NOT be detected since no roots
	f := filepath.Join(dir, "new.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := g.Exit(nil); err != nil {
		t.Fatalf("empty roots should not guard anything: %v", err)
	}
}

func TestGuard_OrigErrWrapped(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(f, []byte("modified"), 0o644)
	// Wrap an error
	origErr := errors.New("some upstream error")
	if err := g.Exit(origErr); err != nil {
		t.Errorf("guard violation should be suppressed when origErr != nil, got: %v", err)
	}
}

func TestGuard_ViolationsAreSorted(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"z.txt", "a.txt", "m.txt"} {
		f := filepath.Join(dir, name)
		if err := os.WriteFile(f, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"."})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// Delete all files
	for _, name := range []string{"z.txt", "a.txt", "m.txt"} {
		os.Remove(filepath.Join(dir, name))
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	pe, ok := err.(*guards.ProtectedPathMutationError)
	if !ok {
		t.Fatalf("expected ProtectedPathMutationError, got %T", err)
	}
	// Violations should be sorted
	for i := 1; i < len(pe.Violations); i++ {
		if pe.Violations[i-1] > pe.Violations[i] {
			t.Errorf("violations not sorted: %v", pe.Violations)
			break
		}
	}
}

func TestNewReadOnlyProjectGuard_RelativePathHandled(t *testing.T) {
	dir := t.TempDir()
	// NewReadOnlyProjectGuard converts to abs; using a relative path should work
	g := guards.NewReadOnlyProjectGuard(dir, []string{"nonexistent"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	if err := g.Exit(nil); err != nil {
		t.Fatalf("missing protected root should not cause error: %v", err)
	}
}
