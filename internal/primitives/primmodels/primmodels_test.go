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
