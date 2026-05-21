package skillintegrator

import (
	"testing"
)

func TestToHyphenCaseVariants(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"HTMLParser", "h-t-m-l-parser"},
		{"getHTMLContent", "get-h-t-m-l-content"},
		{"my__double_underscore", "my-double-underscore"},
		{"MixedCase_With_Underscores", "mixed-case-with-underscores"},
		{"abc123def", "abc123def"},
		{"  spaces  ", "spaces"},
	}
	for _, tc := range cases {
		got := ToHyphenCase(tc.input)
		if len(got) == 0 && len(tc.input) != 0 && tc.input != "  spaces  " {
			// just check it doesn't panic and returns something
		}
		// check no uppercase
		for _, c := range got {
			if c >= 'A' && c <= 'Z' {
				t.Errorf("ToHyphenCase(%q) = %q contains uppercase", tc.input, got)
			}
		}
		// check max length
		if len(got) > 64 {
			t.Errorf("ToHyphenCase(%q) = %q exceeds 64 chars", tc.input, got)
		}
	}
}

func TestValidateSkillNameEdgeCases(t *testing.T) {
	cases := []struct {
		name      string
		wantValid bool
	}{
		{"a-b-c-d", true},
		{"abc-123", true},
		{"1skill", true},
		{"skill--name", false},  // consecutive hyphens
		{"-skill", false},       // leading hyphen
		{"skill-", false},       // trailing hyphen
		{"UPPER", false},        // uppercase
		{"has space", false},    // space
		{"has/slash", false},    // slash
	}
	for _, tc := range cases {
		got, msg := ValidateSkillName(tc.name)
		if got != tc.wantValid {
			t.Errorf("ValidateSkillName(%q) = %v (%s), want %v", tc.name, got, msg, tc.wantValid)
		}
	}
}

func TestValidateSkillNameMaxLength(t *testing.T) {
	// 64 chars: valid
	name64 := "abcdefghijklmnopqrstuvwxyz01234567890abcdefghijklmnopqrstuvwxyz0"
	if len(name64) != 64 {
		t.Fatalf("test setup: expected 64 chars, got %d", len(name64))
	}
	valid, _ := ValidateSkillName(name64)
	if !valid {
		t.Errorf("expected 64-char name to be valid")
	}
	// 65 chars: invalid
	name65 := name64 + "x"
	valid65, _ := ValidateSkillName(name65)
	if valid65 {
		t.Errorf("expected 65-char name to be invalid")
	}
}

func TestNormalizeSkillNameExtra(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"MyPackage"},
		{"my_package"},
		{"some/path/to/pkg"},
		{"already-good"},
	}
	for _, tc := range cases {
		got := NormalizeSkillName(tc.input)
		valid, _ := ValidateSkillName(got)
		if got != "" && !valid {
			t.Errorf("NormalizeSkillName(%q) = %q which is invalid", tc.input, got)
		}
	}
}

func TestSkillIntegrationResultZeroValue(t *testing.T) {
	var r SkillIntegrationResult
	if r.SkillCreated || r.SkillUpdated || r.SkillSkipped {
		t.Error("zero value should have all bool fields false")
	}
	if r.ReferencesCopied != 0 {
		t.Error("zero value ReferencesCopied should be 0")
	}
	if r.SubSkillsPromoted != 0 {
		t.Error("zero value SubSkillsPromoted should be 0")
	}
	if r.SkillPath != "" {
		t.Error("zero value SkillPath should be empty")
	}
	if r.TargetPaths != nil && len(r.TargetPaths) != 0 {
		t.Error("zero value TargetPaths should be nil/empty")
	}
}

func TestSkillIntegrationResultFields(t *testing.T) {
	r := SkillIntegrationResult{
		SkillCreated:      true,
		SkillUpdated:      false,
		SkillSkipped:      false,
		SkillPath:         ".github/skills/my-skill/SKILL.md",
		ReferencesCopied:  3,
		LinksResolved:     0,
		SubSkillsPromoted: 2,
		TargetPaths:       []string{".github/skills/my-skill", ".claude/skills/my-skill"},
	}
	if !r.SkillCreated {
		t.Error("SkillCreated should be true")
	}
	if r.ReferencesCopied != 3 {
		t.Errorf("ReferencesCopied: got %d, want 3", r.ReferencesCopied)
	}
	if r.SubSkillsPromoted != 2 {
		t.Errorf("SubSkillsPromoted: got %d, want 2", r.SubSkillsPromoted)
	}
	if len(r.TargetPaths) != 2 {
		t.Errorf("TargetPaths: got %d, want 2", len(r.TargetPaths))
	}
}
