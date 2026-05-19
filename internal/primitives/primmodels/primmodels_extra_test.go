package primmodels

import (
	"testing"
)

func TestChatmode_AllFields(t *testing.T) {
	c := &Chatmode{
		Name:        "test-mode",
		FilePath:    "/path/to/file.md",
		Description: "A test chatmode",
		ApplyTo:     "**/*.go",
		Content:     "# Content",
		Author:      "Alice",
		Version:     "1.0.0",
		Source:      "local",
	}
	if c.Name != "test-mode" {
		t.Errorf("unexpected Name %q", c.Name)
	}
	if c.Author != "Alice" {
		t.Errorf("unexpected Author %q", c.Author)
	}
	if c.Version != "1.0.0" {
		t.Errorf("unexpected Version %q", c.Version)
	}
}

func TestInstruction_AllFields(t *testing.T) {
	i := &Instruction{
		Name:        "my-instruction",
		FilePath:    "/instr.md",
		Description: "desc",
		Content:     "some content",
		ApplyTo:     "*.py",
		Author:      "Bob",
		Version:     "2.0",
		Source:      "package",
	}
	errs := i.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid Instruction, got %v", errs)
	}
}

func TestInstruction_MissingBoth(t *testing.T) {
	i := &Instruction{}
	errs := i.Validate()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors for empty Instruction, got %d: %v", len(errs), errs)
	}
}

func TestContext_ValidContent(t *testing.T) {
	c := &Context{
		Name:    "ctx",
		Content: "context content",
	}
	errs := c.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid Context, got %v", errs)
	}
}

func TestContext_EmptyContent(t *testing.T) {
	c := &Context{}
	errs := c.Validate()
	if len(errs) != 1 {
		t.Errorf("expected 1 error for empty Context, got %d", len(errs))
	}
}

func TestSkill_Validate_NoErrors(t *testing.T) {
	s := &Skill{Name: "s", Description: "d", Content: "c"}
	errs := s.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no validation errors for Skill, got %v", errs)
	}
}

func TestSkill_AllFields(t *testing.T) {
	s := &Skill{
		Name:        "skill-one",
		FilePath:    "/skill.md",
		Description: "Does something",
		ApplyTo:     "*.ts",
		Content:     "skill content",
		Author:      "Carol",
		Version:     "0.1",
		Source:      "remote",
	}
	if s.Source != "remote" {
		t.Errorf("unexpected Source %q", s.Source)
	}
}

func TestAgent_ZeroValue(t *testing.T) {
	a := Agent{}
	if a.Name != "" || a.Content != "" {
		t.Error("expected zero value Agent to have empty fields")
	}
}

func TestHook_ZeroValue(t *testing.T) {
	h := Hook{}
	if h.Name != "" || h.Content != "" {
		t.Error("expected zero value Hook to have empty fields")
	}
}

func TestNewConflictIndex_Maps(t *testing.T) {
	ci := NewConflictIndex()
	if ci.Chatmodes == nil || ci.Instructions == nil || ci.Skills == nil || ci.Agents == nil {
		t.Fatal("expected all maps to be initialized")
	}
	if len(ci.Chatmodes) != 0 || len(ci.Instructions) != 0 {
		t.Error("expected empty maps")
	}
}

func TestConflictIndex_InsertAndLookup(t *testing.T) {
	ci := NewConflictIndex()
	cm := &Chatmode{Name: "my-mode", Description: "d", Content: "c"}
	ci.Chatmodes["my-mode"] = cm
	got, ok := ci.Chatmodes["my-mode"]
	if !ok {
		t.Fatal("expected to find my-mode in Chatmodes")
	}
	if got.Name != "my-mode" {
		t.Errorf("unexpected Name %q", got.Name)
	}
}

func TestConflictIndex_MultipleTypes(t *testing.T) {
	ci := NewConflictIndex()
	ci.Chatmodes["cm"] = &Chatmode{Name: "cm"}
	ci.Instructions["instr"] = &Instruction{Name: "instr"}
	ci.Skills["sk"] = &Skill{Name: "sk"}
	ci.Agents["ag"] = &Agent{Name: "ag"}

	if len(ci.Chatmodes) != 1 || len(ci.Instructions) != 1 {
		t.Error("unexpected map lengths")
	}
	if len(ci.Skills) != 1 || len(ci.Agents) != 1 {
		t.Error("unexpected map lengths")
	}
}

func TestPrimitive_Interface(t *testing.T) {
	// All primitive types implement the Primitive interface.
	var _ Primitive = (*Chatmode)(nil)
	var _ Primitive = (*Instruction)(nil)
	var _ Primitive = (*Context)(nil)
	var _ Primitive = (*Skill)(nil)
}
