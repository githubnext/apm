package primparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/primitives/primparser"
)

func TestParseSkillFile_EmptySource_Extra3(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(f, []byte("---\nname: skill-a\n---\ncontent here"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	skill, err := primparser.ParseSkillFile(f, "")
	if err != nil {
		t.Fatalf("ParseSkillFile error: %v", err)
	}
	if skill.Source != "" {
		t.Errorf("Source = %q, want empty", skill.Source)
	}
}

func TestParseSkillFile_Description_Extra3(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "SKILL.md")
	content := "---\nname: my-skill\ndescription: A great skill\n---\nsome content"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	skill, err := primparser.ParseSkillFile(f, "local")
	if err != nil {
		t.Fatalf("ParseSkillFile error: %v", err)
	}
	if skill.Description != "A great skill" {
		t.Errorf("Description = %q, want 'A great skill'", skill.Description)
	}
}

func TestParseSkillFile_FilePath_Extra3(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(f, []byte("---\nname: fp-skill\n---\ncontent"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	skill, err := primparser.ParseSkillFile(f, "local")
	if err != nil {
		t.Fatalf("ParseSkillFile error: %v", err)
	}
	if skill.FilePath != f {
		t.Errorf("FilePath = %q, want %q", skill.FilePath, f)
	}
}

func TestParsePrimitiveFile_MissingFile_Extra3(t *testing.T) {
	_, err := primparser.ParsePrimitiveFile("/no/such/file.chatmode.md", "local")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseSkillFile_NoFrontmatter_Extra3(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(f, []byte("just raw content, no frontmatter"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	skill, err := primparser.ParseSkillFile(f, "local")
	if err != nil {
		t.Fatalf("ParseSkillFile error: %v", err)
	}
	if skill.Content == "" {
		t.Error("Content should not be empty")
	}
}
