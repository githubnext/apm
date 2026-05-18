package instructionintegrator_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/integration/instructionintegrator"
)

func TestConvertToCursorRules_DescriptionFromBody(t *testing.T) {
	// No frontmatter -- description is extracted from first non-empty line of body
	content := "# My Rule Title\n\nSome rule content."
	got := instructionintegrator.ConvertToCursorRules(content)
	if !strings.Contains(got, "My Rule Title") {
		t.Errorf("expected description from body heading, got %q", got)
	}
}

func TestConvertToCursorRules_DescriptionFromFrontmatter(t *testing.T) {
	content := "---\ndescription: My Description\n---\n\nBody content."
	got := instructionintegrator.ConvertToCursorRules(content)
	if !strings.Contains(got, "My Description") {
		t.Errorf("expected description from frontmatter, got %q", got)
	}
}

func TestConvertToCursorRules_WithApplyTo(t *testing.T) {
	content := "---\napplyTo: '**/*.go'\ndescription: Go rules\n---\n\nBody."
	got := instructionintegrator.ConvertToCursorRules(content)
	if !strings.Contains(got, "**/*.go") {
		t.Errorf("expected globs pattern in output, got %q", got)
	}
	if !strings.Contains(got, "globs:") {
		t.Errorf("expected 'globs:' in cursor rules output, got %q", got)
	}
}

func TestConvertToCursorRules_NoApplyTo(t *testing.T) {
	content := "---\ndescription: Global rule\n---\n\nBody."
	got := instructionintegrator.ConvertToCursorRules(content)
	if strings.Contains(got, "globs:") {
		t.Errorf("expected no 'globs:' when no applyTo, got %q", got)
	}
}

func TestConvertToCursorRules_HasFrontmatterDelimiters(t *testing.T) {
	content := "# Rule\n\nbody"
	got := instructionintegrator.ConvertToCursorRules(content)
	if !strings.HasPrefix(got, "---\n") {
		t.Errorf("expected output to start with '---', got %q", got[:min6(len(got), 20)])
	}
}

func TestConvertToCursorRules_BodyPreserved(t *testing.T) {
	content := "---\ndescription: D\n---\n\nImportant rule body here."
	got := instructionintegrator.ConvertToCursorRules(content)
	if !strings.Contains(got, "Important rule body here.") {
		t.Errorf("expected body content preserved, got %q", got)
	}
}

func TestConvertToClaudeRules_WithApplyTo(t *testing.T) {
	content := "---\napplyTo: '**/*.py'\n---\n\nPython rules."
	got := instructionintegrator.ConvertToClaudeRules(content)
	if !strings.Contains(got, "paths:") {
		t.Errorf("expected 'paths:' in Claude rules output, got %q", got)
	}
	if !strings.Contains(got, "**/*.py") {
		t.Errorf("expected pattern in paths list, got %q", got)
	}
}

func TestConvertToClaudeRules_NoApplyTo_NoFrontmatter(t *testing.T) {
	content := "---\ndescription: Global\n---\n\nGlobal rules."
	got := instructionintegrator.ConvertToClaudeRules(content)
	if strings.Contains(got, "paths:") {
		t.Errorf("expected no 'paths:' when no applyTo, got %q", got)
	}
	if !strings.Contains(got, "Global rules.") {
		t.Errorf("expected body content, got %q", got)
	}
}

func TestConvertToClaudeRules_PlainContent(t *testing.T) {
	content := "No frontmatter here, just plain text."
	got := instructionintegrator.ConvertToClaudeRules(content)
	if !strings.Contains(got, "plain text") {
		t.Errorf("expected plain text in output, got %q", got)
	}
}

func TestConvertToWindsurfRules_WithApplyTo(t *testing.T) {
	content := "---\napplyTo: '**/*.ts'\n---\n\nTS rules."
	got := instructionintegrator.ConvertToWindsurfRules(content)
	if !strings.Contains(got, "trigger: glob") {
		t.Errorf("expected 'trigger: glob' in windsurf output, got %q", got)
	}
	if !strings.Contains(got, "**/*.ts") {
		t.Errorf("expected pattern in globs, got %q", got)
	}
}

func TestConvertToWindsurfRules_NoApplyTo(t *testing.T) {
	content := "---\ndescription: Always on\n---\n\nRules."
	got := instructionintegrator.ConvertToWindsurfRules(content)
	if !strings.Contains(got, "trigger: always_on") {
		t.Errorf("expected 'trigger: always_on' when no applyTo, got %q", got)
	}
}

func TestConvertToWindsurfRules_BodyPreserved(t *testing.T) {
	content := "---\napplyTo: '*.md'\n---\n\nMarkdown rules body."
	got := instructionintegrator.ConvertToWindsurfRules(content)
	if !strings.Contains(got, "Markdown rules body.") {
		t.Errorf("expected body preserved, got %q", got)
	}
}

func TestConvertToWindsurfRules_HasFrontmatter(t *testing.T) {
	content := "plain content"
	got := instructionintegrator.ConvertToWindsurfRules(content)
	if !strings.HasPrefix(got, "---\n") {
		t.Errorf("expected frontmatter in windsurf output, got %q", got[:min6(len(got), 20)])
	}
}

func TestConvertToCursorRules_DescriptionFromSentence(t *testing.T) {
	// Body has a sentence, description is first sentence
	content := "First sentence. Second sentence."
	got := instructionintegrator.ConvertToCursorRules(content)
	if !strings.Contains(got, "First sentence") {
		t.Errorf("expected first sentence as description, got %q", got)
	}
}

func TestFindInstructionFiles_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	files, err := instructionintegrator.FindInstructionFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected no files in empty dir, got %v", files)
	}
}

func min6(a, b int) int {
	if a < b {
		return a
	}
	return b
}
