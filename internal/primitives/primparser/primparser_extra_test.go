package primparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/primitives/primparser"
)

func TestParseSkillFile_DefaultName(t *testing.T) {
	// When no "name" field in frontmatter, name defaults to parent dir name.
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "my-skill")
	os.MkdirAll(skillDir, 0o755)
	path := filepath.Join(skillDir, "SKILL.md")
	os.WriteFile(path, []byte("---\ndescription: a skill\n---\n\nContent.\n"), 0o644)
	skill, err := primparser.ParseSkillFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", skill.Name)
	}
}

func TestParseSkillFile_ExplicitName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	os.WriteFile(path, []byte("---\nname: ExplicitSkill\ndescription: explicit\n---\n\nBody.\n"), 0o644)
	skill, err := primparser.ParseSkillFile(path, "dep:pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "ExplicitSkill" {
		t.Errorf("expected 'ExplicitSkill', got %q", skill.Name)
	}
	if skill.Source != "dep:pkg" {
		t.Errorf("expected source 'dep:pkg', got %q", skill.Source)
	}
}

func TestParseSkillFile_MissingFile(t *testing.T) {
	_, err := primparser.ParseSkillFile("/nonexistent/SKILL.md", "local")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParsePrimitiveFile_InstructionFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.instructions.md")
	os.WriteFile(path, []byte("---\napplyTo: \"**\"\ndescription: my instructions\n---\n\nContent.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errs := prim.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got %v", errs)
	}
}

func TestParsePrimitiveFile_AgentFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.agent.md")
	os.WriteFile(path, []byte("---\ndescription: my agent\n---\n\nAgent body.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestParsePrimitiveFile_ContextFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "info.context.md")
	os.WriteFile(path, []byte("Context content.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive for .context.md")
	}
}

func TestParsePrimitiveFile_MemoryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notes.memory.md")
	os.WriteFile(path, []byte("Memory content.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive for .memory.md")
	}
}

func TestValidatePrimitive_ValidChatmode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chat.chatmode.md")
	os.WriteFile(path, []byte("---\ndescription: test\napplyTo: '**'\n---\nBody.\n"), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errs := primparser.ValidatePrimitive(prim)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestParsePrimitiveFile_EmptyBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.instructions.md")
	os.WriteFile(path, []byte(""), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error for empty file: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}
