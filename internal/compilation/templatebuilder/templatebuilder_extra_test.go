package templatebuilder_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/templatebuilder"
)

func identity(inst templatebuilder.Instruction) []string {
	return []string{inst.Content, ""}
}

func TestRenderInstructionsBlock_MixedGlobalAndScoped(t *testing.T) {
	instructions := []templatebuilder.Instruction{
		{Name: "global", Content: "global-content"},
		{Name: "scoped", Content: "scoped-content", ApplyTo: "**/*.py"},
	}
	lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", identity)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "## Global Instructions") {
		t.Error("expected global heading")
	}
	if !strings.Contains(joined, "## Files matching") {
		t.Error("expected scoped heading")
	}
	if !strings.Contains(joined, "global-content") {
		t.Error("expected global content")
	}
	if !strings.Contains(joined, "scoped-content") {
		t.Error("expected scoped content")
	}
}

func TestRenderInstructionsBlock_MultiplePatterns(t *testing.T) {
	instructions := []templatebuilder.Instruction{
		{Name: "b", Content: "b-content", ApplyTo: "b-pattern"},
		{Name: "a", Content: "a-content", ApplyTo: "a-pattern"},
		{Name: "c", Content: "c-content", ApplyTo: "c-pattern"},
	}
	lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", identity)
	joined := strings.Join(lines, "\n")
	aIdx := strings.Index(joined, "a-pattern")
	bIdx := strings.Index(joined, "b-pattern")
	cIdx := strings.Index(joined, "c-pattern")
	if aIdx > bIdx || bIdx > cIdx {
		t.Error("patterns should be sorted alphabetically: a < b < c")
	}
}

func TestRenderInstructionsBlock_AllScopedNoGlobal(t *testing.T) {
	instructions := []templatebuilder.Instruction{
		{Name: "x", Content: "x-content", ApplyTo: "**/*.go"},
	}
	lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", identity)
	joined := strings.Join(lines, "\n")
	if strings.Contains(joined, "## Global Instructions") {
		t.Error("should not have global heading when no global instructions")
	}
	if !strings.Contains(joined, "x-content") {
		t.Error("expected scoped content")
	}
}

func TestRenderInstructionsBlock_AllEmptySkipped(t *testing.T) {
	instructions := []templatebuilder.Instruction{
		{Name: "empty1", Content: "", ApplyTo: "**/*.ts"},
		{Name: "empty2", Content: ""},
	}
	lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", identity)
	if len(lines) != 0 {
		t.Errorf("expected no output for all-empty instructions, got %v", lines)
	}
}

func TestRenderInstructionsBlock_SamePatternGrouped(t *testing.T) {
	instructions := []templatebuilder.Instruction{
		{Name: "first", Content: "first-content", ApplyTo: "**/*.go"},
		{Name: "second", Content: "second-content", ApplyTo: "**/*.go"},
	}
	lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", identity)
	joined := strings.Join(lines, "\n")
	// Only one heading for the pattern
	count := strings.Count(joined, "## Files matching `**/*.go`")
	if count != 1 {
		t.Errorf("expected 1 heading for pattern, got %d", count)
	}
	if !strings.Contains(joined, "first-content") || !strings.Contains(joined, "second-content") {
		t.Error("both contents should be present")
	}
}

func TestRenderInstructionsBlock_NilVsEmpty(t *testing.T) {
	linesNil := templatebuilder.RenderInstructionsBlock(nil, "/base", identity)
	linesEmpty := templatebuilder.RenderInstructionsBlock([]templatebuilder.Instruction{}, "/base", identity)
	if len(linesNil) != 0 {
		t.Errorf("nil: expected 0 lines, got %d", len(linesNil))
	}
	if len(linesEmpty) != 0 {
		t.Errorf("empty: expected 0 lines, got %d", len(linesEmpty))
	}
}

func TestRenderInstructionsBlock_GlobalSortedByPath(t *testing.T) {
	instructions := []templatebuilder.Instruction{
		{Name: "z", FilePath: "/base/z.instructions.md", Content: "z-content"},
		{Name: "a", FilePath: "/base/a.instructions.md", Content: "a-content"},
	}
	lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", identity)
	joined := strings.Join(lines, "\n")
	aIdx := strings.Index(joined, "a-content")
	zIdx := strings.Index(joined, "z-content")
	if aIdx > zIdx {
		t.Error("global instructions should be sorted by relative path (a before z)")
	}
}
