package primitives_test

import (
	"testing"

	"github.com/githubnext/apm/internal/primitives"
)

// TestParityPrimitiveTypeValues mirrors PrimitiveType string constants
func TestParityPrimitiveTypeValues(t *testing.T) {
	if primitives.PrimitiveTypeChatmode != "chatmode" {
		t.Error("PrimitiveTypeChatmode should be 'chatmode'")
	}
	if primitives.PrimitiveTypeInstruction != "instruction" {
		t.Error("PrimitiveTypeInstruction should be 'instruction'")
	}
	if primitives.PrimitiveTypeContext != "context" {
		t.Error("PrimitiveTypeContext should be 'context'")
	}
	if primitives.PrimitiveTypeSkill != "skill" {
		t.Error("PrimitiveTypeSkill should be 'skill'")
	}
}

// TestParityChatmodeStruct mirrors Chatmode dataclass fields
func TestParityChatmodeStruct(t *testing.T) {
	cm := primitives.Chatmode{
		Name:        "test-chatmode",
		FilePath:    "/some/path/test-chatmode.chatmode.md",
		Description: "A test chatmode",
		ApplyTo:     "**/*.go",
		Content:     "# Test",
		Author:      "acme",
		Version:     "1.0",
		Source:      "owner/repo",
	}
	if cm.Name != "test-chatmode" {
		t.Errorf("Chatmode.Name = %s, want 'test-chatmode'", cm.Name)
	}
	if cm.ApplyTo != "**/*.go" {
		t.Errorf("Chatmode.ApplyTo = %s, want '**/*.go'", cm.ApplyTo)
	}
}

// TestParityInstructionStruct mirrors Instruction dataclass fields
func TestParityInstructionStruct(t *testing.T) {
	inst := primitives.Instruction{
		Name:     "testing",
		FilePath: "testing.instructions.md",
		ApplyTo:  "**/*.go",
		Content:  "content",
	}
	if inst.Name != "testing" {
		t.Errorf("Instruction.Name = %s, want 'testing'", inst.Name)
	}
}

// TestParityContextStruct mirrors Context dataclass fields
func TestParityContextStruct(t *testing.T) {
	ctx := primitives.Context{
		Name:    "my-context",
		Scope:   "workspace",
		Content: "context content",
	}
	if ctx.Scope != "workspace" {
		t.Errorf("Context.Scope = %s, want 'workspace'", ctx.Scope)
	}
}

// TestParitySkillStruct mirrors Skill dataclass fields
func TestParitySkillStruct(t *testing.T) {
	skill := primitives.Skill{
		Name:        "my-skill",
		FilePath:    "skills/my-skill/SKILL.md",
		Description: "Does something useful",
		Content:     "# My Skill",
	}
	if skill.Name != "my-skill" {
		t.Errorf("Skill.Name = %s, want 'my-skill'", skill.Name)
	}
}

// TestParityPrimitiveConflict mirrors PrimitiveConflict dataclass fields
func TestParityPrimitiveConflict(t *testing.T) {
	conflict := primitives.PrimitiveConflict{
		Type:   primitives.PrimitiveTypeSkill,
		Name:   "duplicate-skill",
		Path1:  "/a/SKILL.md",
		Path2:  "/b/SKILL.md",
		Reason: "duplicate name",
	}
	if conflict.Type != primitives.PrimitiveTypeSkill {
		t.Errorf("PrimitiveConflict.Type = %s, want 'skill'", conflict.Type)
	}
}

// TestParityPrimitiveCollectionTotalCount mirrors PrimitiveCollection total count
func TestParityPrimitiveCollectionTotalCount(t *testing.T) {
	pc := primitives.PrimitiveCollection{
		Chatmodes:    []primitives.Chatmode{{Name: "c1"}, {Name: "c2"}},
		Instructions: []primitives.Instruction{{Name: "i1"}},
		Contexts:     []primitives.Context{},
		Skills:       []primitives.Skill{{Name: "s1"}, {Name: "s2"}, {Name: "s3"}},
	}
	if pc.TotalCount() != 6 {
		t.Errorf("TotalCount() = %d, want 6", pc.TotalCount())
	}
}

// TestParityPrimitiveCollectionHasConflicts mirrors PrimitiveCollection conflict detection
func TestParityPrimitiveCollectionHasConflicts(t *testing.T) {
	pcNoConflicts := primitives.PrimitiveCollection{}
	if pcNoConflicts.HasConflicts() {
		t.Error("empty collection should have no conflicts")
	}

	pcWithConflicts := primitives.PrimitiveCollection{
		Conflicts: []primitives.PrimitiveConflict{{Name: "dup"}},
	}
	if !pcWithConflicts.HasConflicts() {
		t.Error("collection with conflicts should return true")
	}
}

// TestParityFileNameWithoutExt mirrors file stem extraction
func TestParityFileNameWithoutExt(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"SKILL.md", "SKILL"},
		{"test-chatmode.chatmode.md", "test-chatmode.chatmode"},
		{"testing.instructions.md", "testing.instructions"},
		{"/some/path/my-skill/SKILL.md", "SKILL"},
		{"noext", "noext"},
	}
	for _, c := range cases {
		got := primitives.FileNameWithoutExt(c.input)
		if got != c.want {
			t.Errorf("FileNameWithoutExt(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
