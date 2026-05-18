package promptintegrator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetTargetFilename(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"/some/path/foo.prompt.md", "foo.prompt.md"},
		{"bar.prompt.md", "bar.prompt.md"},
		{"/deep/nested/dir/review.prompt.md", "review.prompt.md"},
	}
	for _, tc := range cases {
		got := GetTargetFilename(tc.in)
		if got != tc.want {
			t.Errorf("GetTargetFilename(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFindPromptFiles(t *testing.T) {
	dir := t.TempDir()
	// Create prompt files in root
	for _, name := range []string{"a.prompt.md", "b.prompt.md", "notprompt.md"} {
		os.WriteFile(filepath.Join(dir, name), []byte("content"), 0o644)
	}
	// Create one in .apm/prompts/
	apmDir := filepath.Join(dir, ".apm", "prompts")
	os.MkdirAll(apmDir, 0o755)
	os.WriteFile(filepath.Join(apmDir, "c.prompt.md"), []byte("content"), 0o644)

	files, err := FindPromptFiles(dir)
	if err != nil {
		t.Fatalf("FindPromptFiles error: %v", err)
	}
	// Should find a.prompt.md, b.prompt.md, c.prompt.md (not notprompt.md)
	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d: %v", len(files), files)
	}
	for _, f := range files {
		if !strings.HasSuffix(f, ".prompt.md") {
			t.Errorf("unexpected file: %s", f)
		}
	}
}

func TestFindPromptFilesEmpty(t *testing.T) {
	dir := t.TempDir()
	files, err := FindPromptFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestCopyPrompt(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.prompt.md")
	dst := filepath.Join(dir, "dst.prompt.md")
	content := "# Test Prompt\n\nContent here."
	os.WriteFile(src, []byte(content), 0o644)

	links, err := CopyPrompt(src, dst)
	if err != nil {
		t.Fatalf("CopyPrompt error: %v", err)
	}
	if links != 0 {
		t.Errorf("expected 0 links, got %d", links)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != content {
		t.Errorf("copied content mismatch: got %q", string(data))
	}
}

func TestCopyPromptMissingSource(t *testing.T) {
	dir := t.TempDir()
	_, err := CopyPrompt(filepath.Join(dir, "missing.md"), filepath.Join(dir, "dst.md"))
	if err == nil {
		t.Error("expected error for missing source")
	}
}
