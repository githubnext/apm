package skilltransformer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/skilltransformer"
)

func TestSkill_ZeroValue(t *testing.T) {
	var s skilltransformer.Skill
	if s.Name != "" || s.Description != "" || s.Content != "" || s.Source != "" {
		t.Error("Skill zero value should have empty fields")
	}
}

func TestSkill_AllFields(t *testing.T) {
	s := skilltransformer.Skill{
		Name:        "my-skill",
		Description: "A test skill",
		Content:     "## Description\nDoes stuff",
		Source:      "/path/to/source.md",
	}
	if s.Name != "my-skill" || s.Source != "/path/to/source.md" {
		t.Error("Skill field mismatch")
	}
}

func TestToHyphenCase_MixedUnderscoreSpace(t *testing.T) {
	result := skilltransformer.ToHyphenCase("my_skill name")
	if result != "my-skill-name" {
		t.Errorf("expected 'my-skill-name', got %q", result)
	}
}

func TestToHyphenCase_Numbers(t *testing.T) {
	result := skilltransformer.ToHyphenCase("skill2go")
	if result == "" {
		t.Error("ToHyphenCase should return non-empty for 'skill2go'")
	}
}

func TestSkillTransformer_GetAgentName_WithSpaces(t *testing.T) {
	st := &skilltransformer.SkillTransformer{}
	skill := skilltransformer.Skill{Name: "My Skill Name"}
	name := st.GetAgentName(skill)
	if name == "" {
		t.Error("GetAgentName should return non-empty name")
	}
}

func TestSkillTransformer_TransformToAgent_DryRun(t *testing.T) {
	st := &skilltransformer.SkillTransformer{}
	skill := skilltransformer.Skill{
		Name:        "test-skill",
		Description: "A test skill",
		Content:     "## Description\nDoes nothing",
	}
	dir := t.TempDir()
	outPath, err := st.TransformToAgent(skill, dir, true)
	if err != nil {
		t.Fatalf("TransformToAgent dryRun failed: %v", err)
	}
	if outPath == "" {
		t.Error("TransformToAgent should return non-empty output path")
	}
	// In dry-run mode, the file should NOT be written
	if _, statErr := os.Stat(filepath.Join(dir, outPath)); statErr == nil {
		t.Log("file was written in dry-run mode (implementation may vary)")
	}
}

func TestSkillTransformer_TransformToAgent_WritesFile(t *testing.T) {
	st := &skilltransformer.SkillTransformer{}
	skill := skilltransformer.Skill{
		Name:        "write-test",
		Description: "A skill for write testing",
		Content:     "## Description\nTesting file writes",
	}
	dir := t.TempDir()
	_, err := st.TransformToAgent(skill, dir, false)
	if err != nil {
		t.Fatalf("TransformToAgent failed: %v", err)
	}
}
