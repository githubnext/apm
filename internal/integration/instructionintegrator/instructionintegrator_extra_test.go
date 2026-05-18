package instructionintegrator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindInstructionFiles_NoDir(t *testing.T) {
	tmp := t.TempDir()
	files, err := FindInstructionFiles(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestFindInstructionFiles_WithFiles(t *testing.T) {
	tmp := t.TempDir()
	instrDir := filepath.Join(tmp, ".apm", "instructions")
	os.MkdirAll(instrDir, 0o755)
	os.WriteFile(filepath.Join(instrDir, "lint.instructions.md"), []byte("content"), 0o644)
	os.WriteFile(filepath.Join(instrDir, "style.instructions.md"), []byte("content"), 0o644)
	os.WriteFile(filepath.Join(instrDir, "notaninstruction.md"), []byte("content"), 0o644)

	files, err := FindInstructionFiles(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 instruction files, got %d: %v", len(files), files)
	}
	for _, f := range files {
		if !strings.HasSuffix(f, ".instructions.md") {
			t.Errorf("non-instruction file included: %s", f)
		}
	}
}

func TestCopyInstruction_Verbatim(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.instructions.md")
	dst := filepath.Join(tmp, "dst.instructions.md")
	content := "---\napplyTo: \"**\"\n---\n\nRule content here.\n"
	os.WriteFile(src, []byte(content), 0o644)

	n, err := CopyInstruction(src, dst, FormatVerbatim)
	if err != nil {
		t.Fatalf("CopyInstruction error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 links, got %d", n)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != content {
		t.Errorf("verbatim copy mismatch: got %q", string(got))
	}
}

func TestCopyInstruction_CursorRules(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.instructions.md")
	dst := filepath.Join(tmp, "dst.mdc")
	content := "---\napplyTo: \"**/*.go\"\n---\n\nGo rules.\n"
	os.WriteFile(src, []byte(content), 0o644)

	_, err := CopyInstruction(src, dst, FormatCursorRules)
	if err != nil {
		t.Fatalf("CopyInstruction cursor error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if !strings.Contains(string(got), "globs") {
		t.Errorf("cursor output should contain 'globs': %q", string(got))
	}
}

func TestCopyInstruction_ClaudeRules(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.instructions.md")
	dst := filepath.Join(tmp, "claude.md")
	content := "---\napplyTo: \"src/**\"\n---\n\nClaude rules.\n"
	os.WriteFile(src, []byte(content), 0o644)

	_, err := CopyInstruction(src, dst, FormatClaudeRules)
	if err != nil {
		t.Fatalf("CopyInstruction claude error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if !strings.Contains(string(got), "paths:") {
		t.Errorf("claude output should contain 'paths:': %q", string(got))
	}
}

func TestCopyInstruction_WindsurfRules(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.instructions.md")
	dst := filepath.Join(tmp, "windsurf.md")
	content := "---\napplyTo: \"**/*.ts\"\n---\n\nWindsurf rules.\n"
	os.WriteFile(src, []byte(content), 0o644)

	_, err := CopyInstruction(src, dst, FormatWindsurfRules)
	if err != nil {
		t.Fatalf("CopyInstruction windsurf error: %v", err)
	}
	got, _ := os.ReadFile(dst)
	if !strings.Contains(string(got), "trigger: glob") {
		t.Errorf("windsurf output should contain 'trigger: glob': %q", string(got))
	}
}

func TestCopyInstruction_MissingSource(t *testing.T) {
	tmp := t.TempDir()
	_, err := CopyInstruction(filepath.Join(tmp, "missing.md"), filepath.Join(tmp, "dst.md"), FormatVerbatim)
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestConvertToCursorRules_Empty(t *testing.T) {
	out := ConvertToCursorRules("")
	if len(out) == 0 {
		t.Error("expected non-empty output for empty input")
	}
}

func TestConvertToClaudeRules_EmptyBody(t *testing.T) {
	out := ConvertToClaudeRules("---\napplyTo: \"**\"\n---\n")
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestConvertToWindsurfRules_NoFrontmatter(t *testing.T) {
	out := ConvertToWindsurfRules("just body text\n")
	if !strings.Contains(out, "trigger: always_on") {
		t.Errorf("expected trigger: always_on, got: %s", out)
	}
}

func TestParseFrontmatter_DescriptionOnly(t *testing.T) {
	input := "---\ndescription: my desc\n---\n\nBody text.\n"
	applyTo, desc, body := parseFrontmatter(input)
	if applyTo != "" {
		t.Errorf("expected empty applyTo, got %q", applyTo)
	}
	if desc != "my desc" {
		t.Errorf("expected desc 'my desc', got %q", desc)
	}
	if !strings.Contains(body, "Body text.") {
		t.Errorf("expected body to contain 'Body text.', got %q", body)
	}
}

func TestParseFrontmatter_ApplyToAndDescription(t *testing.T) {
	input := "---\napplyTo: \"*.go\"\ndescription: go rules\n---\n\nContent.\n"
	applyTo, desc, body := parseFrontmatter(input)
	if applyTo != "*.go" {
		t.Errorf("expected applyTo '*.go', got %q", applyTo)
	}
	if desc != "go rules" {
		t.Errorf("expected desc 'go rules', got %q", desc)
	}
	if !strings.Contains(body, "Content.") {
		t.Errorf("expected body to contain 'Content.', got %q", body)
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	input := "Just plain text, no frontmatter.\n"
	applyTo, desc, body := parseFrontmatter(input)
	if applyTo != "" {
		t.Errorf("expected empty applyTo, got %q", applyTo)
	}
	if desc != "" {
		t.Errorf("expected empty desc, got %q", desc)
	}
	if !strings.Contains(body, "plain text") {
		t.Errorf("expected body to contain text, got %q", body)
	}
}
