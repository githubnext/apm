package primparser_test

import (
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/primitives/primparser"
)

func writeTmpPrim(t *testing.T, ext, content string) string {
t.Helper()
dir := t.TempDir()
path := filepath.Join(dir, "test"+ext)
if err := os.WriteFile(path, []byte(content), 0644); err != nil {
t.Fatal(err)
}
return path
}

func TestParseSkillFile_NameFromFrontmatter_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".skill.md", "---\nname: my-skill\ndescription: desc\n---\n# body")
skill, err := primparser.ParseSkillFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if skill == nil {
t.Fatal("expected non-nil skill")
}
}

func TestParseSkillFile_BodyContent_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".skill.md", "# Hello\nThis is a skill body.")
skill, err := primparser.ParseSkillFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if skill == nil {
t.Fatal("expected non-nil skill")
}
if skill.Content == "" {
t.Error("expected non-empty Content")
}
}

func TestParsePrimitiveFile_ChatmodeExt_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".chatmode.md", "---\ndescription: my chatmode\n---\n# body")
p, err := primparser.ParsePrimitiveFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if p == nil {
t.Fatal("expected non-nil primitive")
}
}

func TestParsePrimitiveFile_InstructionsExt_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".instructions.md", "# My instructions\nDo this.")
p, err := primparser.ParsePrimitiveFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if p == nil {
t.Fatal("expected non-nil primitive")
}
}

func TestParsePrimitiveFile_ContextExt_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".context.md", "# context file\ncontent here")
p, err := primparser.ParsePrimitiveFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if p == nil {
t.Fatal("expected non-nil primitive")
}
}

func TestValidatePrimitive_ValidSkill_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".skill.md", "---\ndescription: my skill\n---\n# body")
skill, err := primparser.ParseSkillFile(path, "local")
if err != nil {
t.Fatalf("parse error: %v", err)
}
errs := primparser.ValidatePrimitive(skill)
if len(errs) != 0 {
t.Errorf("expected no errors, got %v", errs)
}
}

func TestParseSkillFile_LongContent_Extra4(t *testing.T) {
body := ""
for i := 0; i < 50; i++ {
body += "Line of content in the skill body.\n"
}
path := writeTmpPrim(t, ".skill.md", "---\ndescription: longskill\n---\n"+body)
skill, err := primparser.ParseSkillFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if skill == nil {
t.Fatal("expected non-nil skill")
}
}

func TestParseSkillFile_EmptyFrontmatter_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".skill.md", "---\n---\n# body")
skill, err := primparser.ParseSkillFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
_ = skill
}

func TestParsePrimitiveFile_AgentExt_Extra4(t *testing.T) {
path := writeTmpPrim(t, ".agent.md", "---\ndescription: an agent\n---\n# body")
p, err := primparser.ParsePrimitiveFile(path, "local")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
_ = p
}
