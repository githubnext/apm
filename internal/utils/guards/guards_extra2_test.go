package guards_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/guards"
)

func TestGuard_NoMutation_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatalf("Enter: %v", err)
	}
	err := g.Exit(nil)
	if err != nil {
		t.Errorf("expected nil on no mutation, got: %v", err)
	}
}

func TestGuard_FileCreated_ReturnsViolation(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	g := guards.NewReadOnlyProjectGuard(dir, []string{"sub"})
	if err := g.Enter(); err != nil {
		t.Fatalf("Enter: %v", err)
	}
	// Create a new file after Enter
	if err := os.WriteFile(filepath.Join(sub, "new.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected violation error, got nil")
	}
	if !strings.Contains(err.Error(), "created") {
		t.Errorf("expected 'created' in error: %v", err)
	}
}

func TestGuard_FileModified_ReturnsViolation(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "mod.txt")
	if err := os.WriteFile(f, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := guards.NewReadOnlyProjectGuard(dir, []string{"mod.txt"})
	if err := g.Enter(); err != nil {
		t.Fatalf("Enter: %v", err)
	}
	// Modify file after Enter
	if err := os.WriteFile(f, []byte("modified content"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err == nil {
		t.Fatal("expected violation error, got nil")
	}
	if !strings.Contains(err.Error(), "modified") {
		t.Errorf("expected 'modified' in error: %v", err)
	}
}

func TestGuard_OrigErrSuppressesViolation(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := guards.NewReadOnlyProjectGuard(dir, []string{"file.txt"})
	if err := g.Enter(); err != nil {
		t.Fatalf("Enter: %v", err)
	}
	// Modify file after Enter
	if err := os.WriteFile(f, []byte("modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	// When origErr is non-nil, the guard suppresses the mutation error and returns nil
	err := g.Exit(os.ErrNotExist)
	if err != nil {
		t.Errorf("expected nil when origErr is non-nil, got: %v", err)
	}
}

func TestGuard_EmptyProtectedPaths_NoViolation(t *testing.T) {
	dir := t.TempDir()
	g := guards.NewReadOnlyProjectGuard(dir, []string{})
	if err := g.Enter(); err != nil {
		t.Fatalf("Enter: %v", err)
	}
	// Create files; guard watches nothing so no violations
	if err := os.WriteFile(filepath.Join(dir, "x.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := g.Exit(nil)
	if err != nil {
		t.Errorf("expected nil for empty protected paths, got: %v", err)
	}
}

func TestGuard_MissingProtectedPath_NoViolation(t *testing.T) {
	dir := t.TempDir()
	g := guards.NewReadOnlyProjectGuard(dir, []string{"nonexistent_subpath"})
	if err := g.Enter(); err != nil {
		t.Fatalf("Enter: %v", err)
	}
	err := g.Exit(nil)
	if err != nil {
		t.Errorf("expected nil for missing protected path, got: %v", err)
	}
}

func TestProtectedPathMutationError_EmptyViolations(t *testing.T) {
	e := &guards.ProtectedPathMutationError{Violations: []string{}}
	msg := e.Error()
	if msg == "" {
		t.Error("expected non-empty error message even with empty violations")
	}
}

func TestProtectedPathMutationError_MultipleViolations(t *testing.T) {
	violations := []string{"created: /a", "deleted: /b", "modified: /c"}
	e := &guards.ProtectedPathMutationError{Violations: violations}
	msg := e.Error()
	for _, v := range violations {
		if !strings.Contains(msg, v) {
			t.Errorf("expected %q in error message: %q", v, msg)
		}
	}
}
