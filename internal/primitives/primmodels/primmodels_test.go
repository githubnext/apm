package primmodels

import (
	"testing"
)

func TestChatmodeValidate(t *testing.T) {
	tests := []struct {
		name     string
		cm       *Chatmode
		wantErrs int
	}{
		{
			name:     "valid chatmode",
			cm:       &Chatmode{Description: "desc", Content: "content"},
			wantErrs: 0,
		},
		{
			name:     "missing description",
			cm:       &Chatmode{Content: "content"},
			wantErrs: 1,
		},
		{
			name:     "missing content",
			cm:       &Chatmode{Description: "desc"},
			wantErrs: 1,
		},
		{
			name:     "missing both",
			cm:       &Chatmode{},
			wantErrs: 2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.cm.Validate()
			if len(errs) != tc.wantErrs {
				t.Errorf("got %d errors, want %d: %v", len(errs), tc.wantErrs, errs)
			}
		})
	}
}

func TestInstructionValidate(t *testing.T) {
	i := &Instruction{Description: "desc", Content: "body"}
	if errs := i.Validate(); len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
	i2 := &Instruction{}
	errs := i2.Validate()
	if len(errs) != 2 {
		t.Errorf("want 2 errors, got %d", len(errs))
	}
}

func TestContextValidate(t *testing.T) {
	c := &Context{Content: "data"}
	if errs := c.Validate(); errs != nil {
		t.Errorf("unexpected errors: %v", errs)
	}
	empty := &Context{}
	errs := empty.Validate()
	if len(errs) != 1 {
		t.Errorf("want 1 error, got %d: %v", len(errs), errs)
	}
}

func TestSkillValidate(t *testing.T) {
	s := &Skill{}
	if errs := s.Validate(); errs != nil {
		t.Errorf("Skill.Validate should return nil, got %v", errs)
	}
}

func TestNewConflictIndex(t *testing.T) {
	ci := NewConflictIndex()
	if ci == nil {
		t.Fatal("NewConflictIndex returned nil")
	}
	if len(ci.Chatmodes) != 0 || len(ci.Instructions) != 0 || len(ci.Skills) != 0 {
		t.Error("expected empty maps")
	}
}

func TestAgentFields(t *testing.T) {
	a := &Agent{
		Name:        "my-agent",
		Description: "does stuff",
		Content:     "## Instructions\n\nDo things.",
		Author:      "alice",
		Version:     "1.0.0",
		Source:      "local",
	}
	if a.Name != "my-agent" {
		t.Errorf("unexpected name: %q", a.Name)
	}
	if a.Description == "" {
		t.Error("description should not be empty")
	}
}

func TestHookFields(t *testing.T) {
	h := &Hook{
		Name:        "pre-commit",
		Description: "runs before commit",
		Content:     "#!/bin/bash\necho hook",
		Author:      "bob",
		Version:     "0.1",
		Source:      "remote",
	}
	if h.Name == "" {
		t.Error("hook name should not be empty")
	}
}

func TestConflictIndexAddAndRetrieve(t *testing.T) {
	ci := NewConflictIndex()
	cm := &Chatmode{Name: "cm1", Description: "d", Content: "c"}
	ci.Chatmodes["cm1"] = cm
	got, ok := ci.Chatmodes["cm1"]
	if !ok || got.Name != "cm1" {
		t.Errorf("expected cm1 in chatmodes, got %v", got)
	}

	inst := &Instruction{Name: "i1", Description: "d", Content: "c"}
	ci.Instructions["i1"] = inst
	gi, ok2 := ci.Instructions["i1"]
	if !ok2 || gi.Name != "i1" {
		t.Errorf("expected i1 in instructions, got %v", gi)
	}

	sk := &Skill{Name: "s1", Description: "d", Content: "c"}
	ci.Skills["s1"] = sk
	gs, ok3 := ci.Skills["s1"]
	if !ok3 || gs.Name != "s1" {
		t.Errorf("expected s1 in skills, got %v", gs)
	}
}

func TestChatmodeAllFields(t *testing.T) {
	cm := &Chatmode{
		Name:        "test-chatmode",
		FilePath:    "/some/path.md",
		Description: "a chatmode",
		ApplyTo:     "*.go",
		Content:     "content here",
		Author:      "alice",
		Version:     "1.2.3",
		Source:      "github.com/org/repo",
	}
	errs := cm.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got: %v", errs)
	}
	if cm.FilePath == "" {
		t.Error("FilePath should not be empty")
	}
	if cm.ApplyTo != "*.go" {
		t.Errorf("ApplyTo mismatch: %q", cm.ApplyTo)
	}
}

func TestInstructionAllFields(t *testing.T) {
	i := &Instruction{
		Name:        "my-inst",
		FilePath:    "/path/inst.md",
		Description: "instruction desc",
		ApplyTo:     "src/**",
		Content:     "do X when Y",
		Author:      "bob",
		Version:     "2.0",
		Source:      "origin",
	}
	if errs := i.Validate(); len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestContextAllFields(t *testing.T) {
	c := &Context{
		Name:        "ctx1",
		FilePath:    "/ctx.md",
		Content:     "some context",
		Description: "context desc",
		Author:      "carol",
		Version:     "1.0",
		Source:      "src",
	}
	if errs := c.Validate(); len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}
