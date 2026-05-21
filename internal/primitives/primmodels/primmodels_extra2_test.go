package primmodels_test

import (
	"testing"

	"github.com/githubnext/apm/internal/primitives/primmodels"
)

func TestChatmode_ValidateNoErrors(t *testing.T) {
	c := &primmodels.Chatmode{Description: "desc", Content: "body"}
	errs := c.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestChatmode_ValidateMissingDescription(t *testing.T) {
	c := &primmodels.Chatmode{Content: "body"}
	errs := c.Validate()
	if len(errs) == 0 {
		t.Error("expected error for missing description")
	}
}

func TestChatmode_ValidateMissingContent(t *testing.T) {
	c := &primmodels.Chatmode{Description: "desc"}
	errs := c.Validate()
	if len(errs) == 0 {
		t.Error("expected error for empty content")
	}
}

func TestInstruction_ValidateNoErrors(t *testing.T) {
	i := &primmodels.Instruction{Description: "desc", Content: "body"}
	errs := i.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestInstruction_ValidateTwoErrors(t *testing.T) {
	i := &primmodels.Instruction{}
	errs := i.Validate()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestContext_ValidateWithContent(t *testing.T) {
	c := &primmodels.Context{Content: "data"}
	errs := c.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestContext_ValidateEmptyContent(t *testing.T) {
	c := &primmodels.Context{}
	errs := c.Validate()
	if len(errs) == 0 {
		t.Error("expected error for empty content")
	}
}

func TestSkill_ValidateAlwaysNil(t *testing.T) {
	s := &primmodels.Skill{}
	errs := s.Validate()
	if errs != nil {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestAgent_ZeroFields(t *testing.T) {
	var a primmodels.Agent
	if a.Name != "" || a.FilePath != "" || a.Content != "" {
		t.Error("zero value should have empty fields")
	}
}

func TestHook_ZeroFields(t *testing.T) {
	var h primmodels.Hook
	if h.Name != "" || h.Content != "" {
		t.Error("zero value should have empty fields")
	}
}

func TestNewConflictIndex_AllMapsInitialized(t *testing.T) {
	idx := primmodels.NewConflictIndex()
	if idx.Chatmodes == nil {
		t.Error("Chatmodes map should be initialized")
	}
	if idx.Instructions == nil {
		t.Error("Instructions map should be initialized")
	}
	if idx.Skills == nil {
		t.Error("Skills map should be initialized")
	}
	if idx.Agents == nil {
		t.Error("Agents map should be initialized")
	}
}

func TestConflictIndex_InsertChatmode(t *testing.T) {
	idx := primmodels.NewConflictIndex()
	c := &primmodels.Chatmode{Name: "test"}
	idx.Chatmodes["test"] = c
	if idx.Chatmodes["test"] != c {
		t.Error("inserted chatmode not found")
	}
}
