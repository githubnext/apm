package promptintegrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/promptintegrator"
)

func TestSyncIntegration_NilManagedFiles_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	removed, errs := promptintegrator.SyncIntegration(tmpDir, nil)
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
	if errs != 0 {
		t.Errorf("expected 0 errors, got %d", errs)
	}
}

func TestSyncIntegration_NilManagedFiles_RemovesApmPrompts(t *testing.T) {
	tmpDir := t.TempDir()
	promptsDir := filepath.Join(tmpDir, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	f1 := filepath.Join(promptsDir, "my-tool-apm.prompt.md")
	f2 := filepath.Join(promptsDir, "other.prompt.md")
	os.WriteFile(f1, []byte("apm file"), 0o644)
	os.WriteFile(f2, []byte("other file"), 0o644)

	removed, _ := promptintegrator.SyncIntegration(tmpDir, nil)
	if removed != 1 {
		t.Errorf("expected 1 removed (-apm.prompt.md), got %d", removed)
	}
	if _, err := os.Stat(f2); err != nil {
		t.Errorf("non-apm file should not be removed")
	}
}

func TestSyncIntegration_WithManagedFiles_RemovesManagedPrompts(t *testing.T) {
	tmpDir := t.TempDir()
	promptsDir := filepath.Join(tmpDir, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	f1 := filepath.Join(promptsDir, "managed.prompt.md")
	f2 := filepath.Join(promptsDir, "user.prompt.md")
	os.WriteFile(f1, []byte("managed"), 0o644)
	os.WriteFile(f2, []byte("user"), 0o644)

	managed := map[string]bool{
		".github/prompts/managed.prompt.md": true,
	}
	removed, _ := promptintegrator.SyncIntegration(tmpDir, managed)
	if removed != 1 {
		t.Errorf("expected 1 removed (managed file), got %d", removed)
	}
	if _, err := os.Stat(f2); err != nil {
		t.Errorf("user file should not be removed")
	}
}

func TestSyncIntegration_WithManagedFiles_NonPromptPathIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	managed := map[string]bool{
		".github/instructions/rule.instructions.md": true,
	}
	removed, errs := promptintegrator.SyncIntegration(tmpDir, managed)
	if removed != 0 {
		t.Errorf("expected 0 removed for non-prompts path, got %d", removed)
	}
	if errs != 0 {
		t.Errorf("expected 0 errors, got %d", errs)
	}
}

func TestFindPromptFiles_RootPromptMd(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "myprompt.prompt.md")
	os.WriteFile(f, []byte("content"), 0o644)
	files, err := promptintegrator.FindPromptFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

func TestFindPromptFiles_NotPromptMd_Ignored(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("readme"), 0o644)
	files, _ := promptintegrator.FindPromptFiles(tmpDir)
	if len(files) != 0 {
		t.Errorf("expected 0 prompt files, got %d", len(files))
	}
}

func TestGetTargetFilename_Simple(t *testing.T) {
	got := promptintegrator.GetTargetFilename("/some/path/myprompt.prompt.md")
	if got != "myprompt.prompt.md" {
		t.Errorf("expected 'myprompt.prompt.md', got %q", got)
	}
}

func TestGetTargetFilename_NestedPath(t *testing.T) {
	got := promptintegrator.GetTargetFilename("/a/b/c/deep.prompt.md")
	if got != "deep.prompt.md" {
		t.Errorf("expected 'deep.prompt.md', got %q", got)
	}
}

func TestCopyPrompt_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.prompt.md")
	dst := filepath.Join(tmpDir, "dest.prompt.md")
	os.WriteFile(src, []byte("prompt content"), 0o644)

	n, err := promptintegrator.CopyPrompt(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 links resolved, got %d", n)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != "prompt content" {
		t.Errorf("expected copied content, got %q", string(data))
	}
}

func TestCopyPrompt_MissingSource(t *testing.T) {
	_, err := promptintegrator.CopyPrompt("/no/such/file.md", "/tmp/dest.md")
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
