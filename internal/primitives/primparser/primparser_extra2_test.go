package primparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/primitives/primparser"
)

func TestParseSkillFile_ExplicitNameField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	os.WriteFile(path, []byte("---\nname: my-skill\ndescription: test\n---\n\nContent.\n"), 0o644)
	skill, err := primparser.ParseSkillFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", skill.Name)
	}
}

func TestParseSkillFile_SourcePreserved(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	os.WriteFile(path, []byte("---\n---\n\nContent.\n"), 0o644)
	skill, err := primparser.ParseSkillFile(path, "dep:mypkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Source != "dep:mypkg" {
		t.Errorf("expected source 'dep:mypkg', got %q", skill.Source)
	}
}

func TestParseSkillFile_ContentExtracted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	os.WriteFile(path, []byte("---\n---\n\nSkill body text.\n"), 0o644)
	skill, err := primparser.ParseSkillFile(path, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestParsePrimitiveFile_InstructionsMd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.instructions.md")
	os.WriteFile(path, []byte("---\ndescription: instructions\n---\n\nBody.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestParsePrimitiveFile_ChatmodeMd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "assistant.chatmode.md")
	os.WriteFile(path, []byte("---\ndescription: a chatmode\n---\n\nBody.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestParsePrimitiveFile_ContextMd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ctx.context.md")
	os.WriteFile(path, []byte("---\n---\n\nContext.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestParsePrimitiveFile_UnknownType(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "unknown.txt")
	os.WriteFile(path, []byte("content\n"), 0o644)
	_, err := primparser.ParsePrimitiveFile(path, "local")
	if err == nil {
		t.Error("expected error for unknown file type")
	}
}

func TestValidatePrimitive_Missing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "foo.instructions.md")
	os.WriteFile(path, []byte("---\n---\n\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	errs := primparser.ValidatePrimitive(prim)
	// No description and empty content: expect validation errors
	if len(errs) == 0 {
		t.Error("expected validation errors for missing description and empty content")
	}
}
