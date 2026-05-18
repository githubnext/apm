package skilltransformer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/integration/skilltransformer"
)

func TestToHyphenCase_Underscores(t *testing.T) {
	got := skilltransformer.ToHyphenCase("my_skill_name")
	if got != "my-skill-name" {
		t.Errorf("expected 'my-skill-name', got %q", got)
	}
}

func TestToHyphenCase_Spaces(t *testing.T) {
	got := skilltransformer.ToHyphenCase("my skill name")
	if got != "my-skill-name" {
		t.Errorf("expected 'my-skill-name', got %q", got)
	}
}

func TestToHyphenCase_CamelCase(t *testing.T) {
	got := skilltransformer.ToHyphenCase("mySkillName")
	if got != "my-skill-name" {
		t.Errorf("expected 'my-skill-name', got %q", got)
	}
}

func TestToHyphenCase_AlreadyHyphenated(t *testing.T) {
	got := skilltransformer.ToHyphenCase("my-skill")
	if got != "my-skill" {
		t.Errorf("expected 'my-skill', got %q", got)
	}
}

func TestToHyphenCase_LowerCased(t *testing.T) {
	got := skilltransformer.ToHyphenCase("MYSCILL")
	if strings.ToLower(got) != got {
		t.Errorf("expected all lowercase, got %q", got)
	}
}

func TestToHyphenCase_Empty(t *testing.T) {
	got := skilltransformer.ToHyphenCase("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestToHyphenCase_ConsecutiveUnderscores(t *testing.T) {
	got := skilltransformer.ToHyphenCase("a__b")
	// double underscore becomes double hyphen, then collapsed to single
	if got != "a-b" {
		t.Errorf("expected 'a-b', got %q", got)
	}
}

func TestToHyphenCase_LeadingTrailingHyphens(t *testing.T) {
	got := skilltransformer.ToHyphenCase("_skill_")
	if strings.HasPrefix(got, "-") {
		t.Errorf("expected no leading hyphen, got %q", got)
	}
	if strings.HasSuffix(got, "-") {
		t.Errorf("expected no trailing hyphen, got %q", got)
	}
}

func TestGetAgentName_Simple(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "my-skill"}
	got := tr.GetAgentName(s)
	if got != "my-skill" {
		t.Errorf("expected 'my-skill', got %q", got)
	}
}

func TestGetAgentName_CamelCase(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "mySkill"}
	got := tr.GetAgentName(s)
	if got != "my-skill" {
		t.Errorf("expected 'my-skill', got %q", got)
	}
}

func TestTransformToAgent_DryRun_ReturnsPath(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "test-skill", Description: "A test", Content: "content"}
	path, err := tr.TransformToAgent(s, "/tmp/fake-dir", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(path, "test-skill.agent.md") {
		t.Errorf("expected path to contain 'test-skill.agent.md', got %q", path)
	}
	if !strings.Contains(path, ".github") {
		t.Errorf("expected path to contain '.github', got %q", path)
	}
}

func TestTransformToAgent_DryRun_NoFileCreated(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "dry-skill", Description: "d", Content: "c"}
	tmpDir := t.TempDir()
	path, err := tr.TransformToAgent(s, tmpDir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Error("expected file NOT to exist in dry-run mode")
	}
}

func TestTransformToAgent_WritesFile(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "real-skill", Description: "desc", Content: "body content"}
	tmpDir := t.TempDir()
	path, err := tr.TransformToAgent(s, tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created at %s: %v", path, err)
	}
	content := string(data)
	if !strings.Contains(content, "body content") {
		t.Errorf("expected 'body content' in file, got %q", content)
	}
	if !strings.Contains(content, "real-skill") {
		t.Errorf("expected 'real-skill' in file, got %q", content)
	}
}

func TestTransformToAgent_ContentHasFrontmatter(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "fm-skill", Description: "my desc", Content: "## Overview\n\ncontent"}
	tmpDir := t.TempDir()
	path, err := tr.TransformToAgent(s, tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Errorf("expected frontmatter at start, got %q", content[:min5(len(content), 20)])
	}
	if !strings.Contains(content, "my desc") {
		t.Errorf("expected description in frontmatter, got %q", content)
	}
}

func TestTransformToAgent_WithSource(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "src-skill", Description: "d", Content: "content", Source: "owner/repo"}
	tmpDir := t.TempDir()
	path, err := tr.TransformToAgent(s, tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "owner/repo") {
		t.Errorf("expected source 'owner/repo' in content, got %q", string(data))
	}
}

func TestTransformToAgent_LocalSourceNotIncluded(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "local-skill", Description: "d", Content: "content", Source: "local"}
	tmpDir := t.TempDir()
	path, err := tr.TransformToAgent(s, tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "<!-- Source:") {
		t.Errorf("local source should not emit source comment, got %q", string(data))
	}
}

func TestTransformToAgent_AgentFileNameFromHyphenCase(t *testing.T) {
	tr := &skilltransformer.SkillTransformer{}
	s := skilltransformer.Skill{Name: "MyGreatSkill", Description: "d", Content: "c"}
	tmpDir := t.TempDir()
	path, err := tr.TransformToAgent(s, tmpDir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := filepath.Base(path)
	if base != "my-great-skill.agent.md" {
		t.Errorf("expected 'my-great-skill.agent.md', got %q", base)
	}
}

func min5(a, b int) int {
	if a < b {
		return a
	}
	return b
}
