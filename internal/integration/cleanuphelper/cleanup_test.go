package cleanuphelper_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/cleanuphelper"
)

// --- ValidateDeployPath tests ---

func TestValidateDeployPath_ValidPrefix(t *testing.T) {
	ok := cleanuphelper.ValidateDeployPath(".github/prompts/foo.md", "/project", []string{".github/prompts"})
	if !ok {
		t.Fatal("expected true for valid prefix")
	}
}

func TestValidateDeployPath_TraversalRejected(t *testing.T) {
	ok := cleanuphelper.ValidateDeployPath(".github/../secret", "/project", []string{".github"})
	if ok {
		t.Fatal("expected false for traversal path")
	}
}

func TestValidateDeployPath_AbsoluteRejected(t *testing.T) {
	ok := cleanuphelper.ValidateDeployPath("/etc/passwd", "/project", []string{".github"})
	if ok {
		t.Fatal("expected false for absolute path")
	}
}

func TestValidateDeployPath_CoworkURIRejected(t *testing.T) {
	ok := cleanuphelper.ValidateDeployPath("cowork://some/path", "/project", []string{".github"})
	if ok {
		t.Fatal("expected false for cowork:// URI")
	}
}

func TestValidateDeployPath_NoPrefixMatch(t *testing.T) {
	ok := cleanuphelper.ValidateDeployPath("other/file.md", "/project", []string{".github"})
	if ok {
		t.Fatal("expected false for no matching prefix")
	}
}

// --- RemoveStaleDeployedFiles tests ---

func TestRemoveStaleDeployedFiles_DeletesFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".github")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	fp := filepath.Join(sub, "prompt.md")
	if err := os.WriteFile(fp, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := cleanuphelper.RemoveStaleDeployedFiles(
		[]string{".github/prompt.md"},
		cleanuphelper.Options{
			DepKey:              "pkg",
			ProjectRoot:         dir,
			IntegrationPrefixes: []string{".github"},
		},
	)
	if len(result.Deleted) != 1 || result.Deleted[0] != ".github/prompt.md" {
		t.Fatalf("unexpected deleted: %v", result.Deleted)
	}
	if _, err := os.Stat(fp); !os.IsNotExist(err) {
		t.Fatal("file should have been removed")
	}
}

func TestRemoveStaleDeployedFiles_SkipsUnmanaged(t *testing.T) {
	dir := t.TempDir()
	result := cleanuphelper.RemoveStaleDeployedFiles(
		[]string{"outside/file.md"},
		cleanuphelper.Options{
			DepKey:              "pkg",
			ProjectRoot:         dir,
			IntegrationPrefixes: []string{".github"},
		},
	)
	if len(result.SkippedUnmanaged) != 1 {
		t.Fatalf("expected 1 skipped, got %v", result.SkippedUnmanaged)
	}
}

func TestRemoveStaleDeployedFiles_SkipsDirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, ".github", "prompts")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	diag := &cleanuphelper.DiagnosticCollector{}
	result := cleanuphelper.RemoveStaleDeployedFiles(
		[]string{".github/prompts"},
		cleanuphelper.Options{
			DepKey:              "pkg",
			ProjectRoot:         dir,
			IntegrationPrefixes: []string{".github"},
			Diagnostics:         diag,
		},
	)
	if len(result.SkippedUnmanaged) != 1 {
		t.Fatalf("expected 1 directory-skipped, got %v", result.SkippedUnmanaged)
	}
}

func TestRemoveStaleDeployedFiles_AlreadyGone(t *testing.T) {
	dir := t.TempDir()
	result := cleanuphelper.RemoveStaleDeployedFiles(
		[]string{".github/gone.md"},
		cleanuphelper.Options{
			DepKey:              "pkg",
			ProjectRoot:         dir,
			IntegrationPrefixes: []string{".github"},
		},
	)
	// Already-gone is a no-op -- not an error.
	if len(result.Failed) != 0 {
		t.Fatalf("expected no failures for already-gone file, got %v", result.Failed)
	}
}

func TestDiagnosticCollector_Warn(t *testing.T) {
	d := &cleanuphelper.DiagnosticCollector{}
	d.Warn("pkg", "something happened")
	if len(d.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(d.Warnings))
	}
	if d.Warnings[0].Package != "pkg" || d.Warnings[0].Message != "something happened" {
		t.Errorf("unexpected warning: %+v", d.Warnings[0])
	}
}
