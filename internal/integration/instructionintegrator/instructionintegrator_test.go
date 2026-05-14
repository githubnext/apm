package instructionintegrator

import (
	"testing"
)

func TestConvertToCursorRules_WithApplyTo(t *testing.T) {
	input := "---\napplyTo: \"**/*.go\"\ndescription: Go lint rules\n---\n\nContent here.\n"
	out := ConvertToCursorRules(input)
	if !contains(out, `globs: "**/*.go"`) {
		t.Errorf("expected globs field, got: %s", out)
	}
	if !contains(out, "description: Go lint rules") {
		t.Errorf("expected description field, got: %s", out)
	}
}

func TestConvertToCursorRules_NoApplyTo(t *testing.T) {
	input := "# My Rule\n\nDo this.\n"
	out := ConvertToCursorRules(input)
	if !contains(out, "---") {
		t.Errorf("expected frontmatter, got: %s", out)
	}
	if !contains(out, "Do this.") {
		t.Errorf("expected body, got: %s", out)
	}
}

func TestConvertToClaudeRules_WithApplyTo(t *testing.T) {
	input := "---\napplyTo: \"src/**\"\n---\n\nBody.\n"
	out := ConvertToClaudeRules(input)
	if !contains(out, `"src/**"`) {
		t.Errorf("expected path, got: %s", out)
	}
	if !contains(out, "paths:") {
		t.Errorf("expected paths key, got: %s", out)
	}
}

func TestConvertToClaudeRules_NoApplyTo(t *testing.T) {
	input := "---\ndescription: foo\n---\n\nBody.\n"
	out := ConvertToClaudeRules(input)
	if contains(out, "paths:") {
		t.Errorf("unexpected paths key, got: %s", out)
	}
	if !contains(out, "Body.") {
		t.Errorf("expected body, got: %s", out)
	}
}

func TestConvertToWindsurfRules_WithApplyTo(t *testing.T) {
	input := "---\napplyTo: \"**/*.ts\"\n---\n\nBody.\n"
	out := ConvertToWindsurfRules(input)
	if !contains(out, "trigger: glob") {
		t.Errorf("expected trigger: glob, got: %s", out)
	}
	if !contains(out, `globs: "**/*.ts"`) {
		t.Errorf("expected globs field, got: %s", out)
	}
}

func TestConvertToWindsurfRules_NoApplyTo(t *testing.T) {
	input := "Body.\n"
	out := ConvertToWindsurfRules(input)
	if !contains(out, "trigger: always_on") {
		t.Errorf("expected trigger: always_on, got: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
