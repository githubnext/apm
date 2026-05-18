package skilltransformer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToHyphenCase(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"mySkill", "my-skill"},
		{"my_skill", "my-skill"},
		{"my skill", "my-skill"},
		{"MySkill", "my-skill"},
		{"already-hyphen", "already-hyphen"},
		{"foo_bar_baz", "foo-bar-baz"},
		{"fooBarBaz", "foo-bar-baz"},
		{"foo--bar", "foo-bar"},
		{"-foo-", "foo"},
	}
	for _, tc := range cases {
		got := ToHyphenCase(tc.in)
		if got != tc.want {
			t.Errorf("ToHyphenCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestTransformToAgent(t *testing.T) {
	dir := t.TempDir()
	st := &SkillTransformer{}
	skill := Skill{
		Name:        "mySkill",
		Description: "does something",
		Content:     "## Instructions\n\nDo the thing.",
		Source:      "local",
	}
	path, err := st.TransformToAgent(skill, dir, false)
	if err != nil {
		t.Fatalf("TransformToAgent error: %v", err)
	}
	if !strings.HasSuffix(path, "my-skill.agent.md") {
		t.Errorf("unexpected path: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "name: mySkill") {
		t.Errorf("missing name in content: %s", content)
	}
	if !strings.Contains(content, "description: does something") {
		t.Errorf("missing description in content: %s", content)
	}
	// local source should not produce Source comment
	if strings.Contains(content, "Source:") {
		t.Errorf("unexpected Source comment for local skill: %s", content)
	}
}

func TestTransformToAgentWithRemoteSource(t *testing.T) {
	dir := t.TempDir()
	st := &SkillTransformer{}
	skill := Skill{
		Name:        "remote-skill",
		Description: "remote skill",
		Content:     "content",
		Source:      "https://example.com/skill.md",
	}
	path, err := st.TransformToAgent(skill, dir, false)
	if err != nil {
		t.Fatalf("TransformToAgent error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "Source: https://example.com/skill.md") {
		t.Errorf("expected Source comment in output: %s", string(data))
	}
}

func TestTransformToAgentDryRun(t *testing.T) {
	dir := t.TempDir()
	st := &SkillTransformer{}
	skill := Skill{Name: "testSkill", Description: "d", Content: "c"}
	path, err := st.TransformToAgent(skill, dir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(path, "test-skill.agent.md") {
		t.Errorf("unexpected path: %s", path)
	}
	// File should not be created in dry-run mode
	if _, err := os.Stat(filepath.Join(dir, ".github", "agents", "test-skill.agent.md")); !os.IsNotExist(err) {
		t.Error("file should not exist in dry-run mode")
	}
}

func TestGetAgentName(t *testing.T) {
	st := &SkillTransformer{}
	cases := []struct {
		name string
		want string
	}{
		{"MySkill", "my-skill"},
		{"code_review", "code-review"},
		{"PR Helper", "pr-helper"},
	}
	for _, tc := range cases {
		got := st.GetAgentName(Skill{Name: tc.name})
		if got != tc.want {
			t.Errorf("GetAgentName(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}
