package guards_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/guards"
)

func TestNoMutation(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// No mutation.
	if err := g.Exit(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestModificationDetected(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// Mutate the file.
	if err := os.WriteFile(f, []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected error for modified file")
	}
	var pe *guards.ProtectedPathMutationError
	if ok := errorAs(err, &pe); !ok {
		t.Fatalf("expected ProtectedPathMutationError, got %T", err)
	}
}

func TestDeletionDetected(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	os.Remove(f)
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected error for deleted file")
	}
}

func TestCreationDetected(t *testing.T) {
	dir := t.TempDir()

	g := guards.NewReadOnlyProjectGuard(dir, []string{"."})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	// Create a new file.
	f := filepath.Join(dir, "new.txt")
	if err := os.WriteFile(f, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected error for created file")
	}
}

func TestMissingRootSilentlyIgnored(t *testing.T) {
	dir := t.TempDir()
	g := guards.NewReadOnlyProjectGuard(dir, []string{"nonexistent"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	if err := g.Exit(nil); err != nil {
		t.Fatalf("expected no error for missing root, got: %v", err)
	}
}

func TestOrigErrSuppressesGuardError(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(f, []byte("changed"), 0o644)
	// When there is an original error, guard violation should be suppressed.
	origErr := os.ErrNotExist
	if err := g.Exit(origErr); err != nil {
		t.Fatalf("expected guard to be suppressed, got: %v", err)
	}
}

// errorAs is a minimal errors.As helper to avoid importing errors package.
func errorAs(err error, target **guards.ProtectedPathMutationError) bool {
	if pe, ok := err.(*guards.ProtectedPathMutationError); ok {
		*target = pe
		return true
	}
	return false
}
