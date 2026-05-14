package primparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/primitives/primparser"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestParseFrontmatterNoFM(t *testing.T) {
	path := writeTmp(t, "just content\nno frontmatter\n")
	// Rename to .instructions.md so ParsePrimitiveFile picks it up.
	newPath := filepath.Join(filepath.Dir(path), "foo.instructions.md")
	os.Rename(path, newPath)
	prim, err := primparser.ParsePrimitiveFile(newPath, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prim == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestParseFrontmatterWithFM(t *testing.T) {
	content := "---\nname: TestSkill\ndescription: A test skill\n---\n# Body\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	os.WriteFile(path, []byte(content), 0o644)
	skill, err := primparser.ParseSkillFile(path, "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "TestSkill" {
		t.Errorf("expected name 'TestSkill', got %q", skill.Name)
	}
	if skill.Description != "A test skill" {
		t.Errorf("expected description 'A test skill', got %q", skill.Description)
	}
	if !contains(skill.Content, "# Body") {
		t.Errorf("expected content to contain '# Body', got %q", skill.Content)
	}
}

func TestParseChatmode(t *testing.T) {
	content := "---\ndescription: My chatmode\napplyTo: '**'\n---\nChatmode body\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "test.chatmode.md")
	os.WriteFile(path, []byte(content), 0o644)
	prim, err := primparser.ParsePrimitiveFile(path, "dep:pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	errs := prim.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got %v", errs)
	}
}

func TestParseUnknownType(t *testing.T) {
	path := writeTmp(t, "content")
	_, err := primparser.ParsePrimitiveFile(path, "local")
	if err == nil {
		t.Fatal("expected error for unknown primitive type")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
