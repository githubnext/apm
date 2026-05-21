package apmpackage_test

import (
	"testing"

	"github.com/githubnext/apm/internal/models/apmpackage"
)

func TestContentTypeString_Instructions(t *testing.T) {
	if apmpackage.ContentTypeInstructions.String() != "instructions" {
		t.Errorf("got %q", apmpackage.ContentTypeInstructions.String())
	}
}

func TestContentTypeString_Skill(t *testing.T) {
	if apmpackage.ContentTypeSkill.String() != "skill" {
		t.Errorf("got %q", apmpackage.ContentTypeSkill.String())
	}
}

func TestContentTypeString_Hybrid(t *testing.T) {
	if apmpackage.ContentTypeHybrid.String() != "hybrid" {
		t.Errorf("got %q", apmpackage.ContentTypeHybrid.String())
	}
}

func TestContentTypeString_Prompts(t *testing.T) {
	if apmpackage.ContentTypePrompts.String() != "prompts" {
		t.Errorf("got %q", apmpackage.ContentTypePrompts.String())
	}
}

func TestParseContentType_Prompts(t *testing.T) {
	ct, err := apmpackage.ParseContentType("prompts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != apmpackage.ContentTypePrompts {
		t.Errorf("got %v, want ContentTypePrompts", ct)
	}
}

func TestParseContentType_Instructions(t *testing.T) {
	ct, err := apmpackage.ParseContentType("instructions")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != apmpackage.ContentTypeInstructions {
		t.Errorf("got %v, want ContentTypeInstructions", ct)
	}
}

func TestParseContentType_InvalidReturnsError(t *testing.T) {
	_, err := apmpackage.ParseContentType("nonsense")
	if err == nil {
		t.Error("expected error for invalid content type")
	}
}

func TestParseContentType_RoundTrip(t *testing.T) {
	for _, want := range []apmpackage.PackageContentType{
		apmpackage.ContentTypeInstructions,
		apmpackage.ContentTypeSkill,
		apmpackage.ContentTypeHybrid,
		apmpackage.ContentTypePrompts,
	} {
		got, err := apmpackage.ParseContentType(want.String())
		if err != nil {
			t.Errorf("ParseContentType(%q) error: %v", want.String(), err)
		}
		if got != want {
			t.Errorf("round-trip failed: got %v, want %v", got, want)
		}
	}
}
