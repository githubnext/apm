package primmodels_test

import (
	"testing"

	"github.com/githubnext/apm/internal/primitives/primmodels"
)

func TestChatmode_AllFields_Extra3(t *testing.T) {
	c := primmodels.Chatmode{
		Name:        "my-chatmode",
		FilePath:    "/path/to/file",
		Description: "desc",
		ApplyTo:     "**/*.go",
		Content:     "content here",
		Author:      "author",
		Version:     "1.0",
		Source:      "local",
	}
	if c.Name != "my-chatmode" {
		t.Errorf("Name = %q", c.Name)
	}
	if c.Source != "local" {
		t.Errorf("Source = %q", c.Source)
	}
}

func TestInstruction_AllFields_Extra3(t *testing.T) {
	i := primmodels.Instruction{
		Name:        "my-instruction",
		FilePath:    "/path/to/file",
		Description: "desc",
		ApplyTo:     "**/*.py",
		Content:     "content",
		Author:      "author",
		Version:     "2.0",
		Source:      "dep:pkg",
	}
	if i.ApplyTo != "**/*.py" {
		t.Errorf("ApplyTo = %q", i.ApplyTo)
	}
	if i.Version != "2.0" {
		t.Errorf("Version = %q", i.Version)
	}
}

func TestSkill_AllFields_Extra3(t *testing.T) {
	s := primmodels.Skill{
		Name:        "my-skill",
		Description: "a skill",
		Content:     "skill content",
		Author:      "author",
		Version:     "0.1",
		Source:      "marketplace",
	}
	if s.Source != "marketplace" {
		t.Errorf("Source = %q", s.Source)
	}
	errs := s.Validate()
	if len(errs) != 0 {
		t.Errorf("Skill.Validate() should return nil, got %v", errs)
	}
}

func TestAgent_AllFields_Extra3(t *testing.T) {
	a := primmodels.Agent{
		Name:        "my-agent",
		FilePath:    "/agents/agent.md",
		Description: "an agent",
		Content:     "agent body",
		Author:      "owner",
		Version:     "3.0",
		Source:      "local",
	}
	if a.Name != "my-agent" {
		t.Errorf("Name = %q", a.Name)
	}
}

func TestHook_AllFields_Extra3(t *testing.T) {
	h := primmodels.Hook{
		Name:        "pre-push",
		FilePath:    "/hooks/pre-push.md",
		Description: "runs before push",
		Content:     "hook content",
	}
	if h.Description != "runs before push" {
		t.Errorf("Description = %q", h.Description)
	}
}

func TestConflictIndex_InsertAll_Extra3(t *testing.T) {
	idx := primmodels.NewConflictIndex()
	idx.Chatmodes["chat1"] = &primmodels.Chatmode{Name: "chat1"}
	idx.Instructions["instr1"] = &primmodels.Instruction{Name: "instr1"}
	idx.Skills["skill1"] = &primmodels.Skill{Name: "skill1"}
	idx.Agents["agent1"] = &primmodels.Agent{Name: "agent1"}
	if len(idx.Chatmodes) != 1 {
		t.Errorf("Chatmodes count = %d, want 1", len(idx.Chatmodes))
	}
	if len(idx.Agents) != 1 {
		t.Errorf("Agents count = %d, want 1", len(idx.Agents))
	}
}

func TestChatmode_Implements_Primitive_Extra3(t *testing.T) {
	var _ primmodels.Primitive = &primmodels.Chatmode{}
	var _ primmodels.Primitive = &primmodels.Instruction{}
	var _ primmodels.Primitive = &primmodels.Context{}
	var _ primmodels.Primitive = &primmodels.Skill{}
}
