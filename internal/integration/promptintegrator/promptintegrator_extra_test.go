package promptintegrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIntegratePackagePrompts_EmptyPackage(t *testing.T) {
	dir := t.TempDir()
	res, err := IntegratePackagePrompts(dir, dir, false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FilesIntegrated != 0 {
		t.Errorf("expected 0 integrated files, got %d", res.FilesIntegrated)
	}
}

func TestIntegratePackagePrompts_SingleFile(t *testing.T) {
	pkg := t.TempDir()
	proj := t.TempDir()
	promptsDir := filepath.Join(proj, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(pkg, "review.prompt.md"), []byte("# Review"), 0o644)
	res, err := IntegratePackagePrompts(pkg, proj, false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FilesIntegrated < 1 {
		t.Errorf("expected at least 1 integrated file, got %d", res.FilesIntegrated)
	}
}

func TestGetTargetFilename_WithDir(t *testing.T) {
	if got := GetTargetFilename("/a/b/c.prompt.md"); got != "c.prompt.md" {
		t.Errorf("got %q, want c.prompt.md", got)
	}
}

func TestCopyPrompt_LargeContent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "big.prompt.md")
	dst := filepath.Join(dir, "out.prompt.md")
	big := make([]byte, 64*1024)
	for i := range big {
		big[i] = 'a'
	}
	os.WriteFile(src, big, 0o644)
	_, err := CopyPrompt(src, dst)
	if err != nil {
		t.Fatalf("CopyPrompt failed: %v", err)
	}
	data, _ := os.ReadFile(dst)
	if len(data) != len(big) {
		t.Errorf("expected %d bytes, got %d", len(big), len(data))
	}
}

func TestFindPromptFiles_ApmPromptsDeep(t *testing.T) {
	dir := t.TempDir()
	deepDir := filepath.Join(dir, ".apm", "prompts", "sub")
	os.MkdirAll(deepDir, 0o755)
	os.WriteFile(filepath.Join(deepDir, "nested.prompt.md"), []byte("# nested"), 0o644)
	files, err := FindPromptFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) < 1 {
		t.Error("expected at least 1 prompt file in deep subdir")
	}
}

func TestFindPromptFiles_IgnoresNonPromptFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# readme"), 0o644)
	os.WriteFile(filepath.Join(dir, "instructions.md"), []byte("# instructions"), 0o644)
	os.WriteFile(filepath.Join(dir, "real.prompt.md"), []byte("# real"), 0o644)
	files, err := FindPromptFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 prompt file, got %d: %v", len(files), files)
	}
}

func TestCopyPrompt_DestMissingDir(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.prompt.md")
	os.WriteFile(src, []byte("content"), 0o644)
	dst := filepath.Join(dir, "nonexistent", "out.prompt.md")
	_, err := CopyPrompt(src, dst)
	if err == nil {
		t.Error("expected error when destination directory does not exist")
	}
}
