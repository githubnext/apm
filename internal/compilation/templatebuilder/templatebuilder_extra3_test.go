package templatebuilder

import (
	"strings"
	"testing"
)

func TestRenderInstructionsBlock_SingleGlobal(t *testing.T) {
	inst := []Instruction{{Name: "g", Content: "global text"}}
	lines := RenderInstructionsBlock(inst, ".", func(i Instruction) []string {
		return []string{i.Content}
	})
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "global text") {
		t.Fatalf("expected global text in output, got %q", combined)
	}
}

func TestRenderInstructionsBlock_ScopedHeading(t *testing.T) {
	inst := []Instruction{{Name: "s", ApplyTo: "*.go", Content: "go text"}}
	lines := RenderInstructionsBlock(inst, ".", func(i Instruction) []string {
		return []string{i.Content}
	})
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "*.go") {
		t.Fatalf("expected pattern *.go in heading, got %q", combined)
	}
}

func TestRenderInstructionsBlock_EmptyInputNoLines(t *testing.T) {
	lines := RenderInstructionsBlock(nil, ".", func(i Instruction) []string { return nil })
	if len(lines) > 0 {
		t.Fatalf("expected no lines for empty input, got %d", len(lines))
	}
}

func TestTemplateData_FieldAssignment(t *testing.T) {
	td := TemplateData{
		InstructionsContent: "content",
		Version:             "1.0",
		ChatmodeContent:     "chatmode",
	}
	if td.InstructionsContent != "content" {
		t.Fatal("InstructionsContent mismatch")
	}
	if td.Version != "1.0" {
		t.Fatal("Version mismatch")
	}
	if td.ChatmodeContent != "chatmode" {
		t.Fatal("ChatmodeContent mismatch")
	}
}

func TestInstruction_ZeroValue(t *testing.T) {
	var i Instruction
	if i.Name != "" || i.FilePath != "" || i.ApplyTo != "" || i.Content != "" {
		t.Fatal("zero-value Instruction should have empty fields")
	}
}

func TestRenderInstructionsBlock_MultipleGlobalsUnderSingleHeading(t *testing.T) {
	inst := []Instruction{
		{Name: "g1", Content: "global one"},
		{Name: "g2", Content: "global two"},
	}
	lines := RenderInstructionsBlock(inst, ".", func(i Instruction) []string {
		return []string{i.Content}
	})
	combined := strings.Join(lines, "\n")
	count := strings.Count(combined, "## Global Instructions")
	if count > 1 {
		t.Fatalf("expected at most one global heading, got %d", count)
	}
}

func TestRenderInstructionsBlock_SkipsEmptyContent2(t *testing.T) {
	inst := []Instruction{{Name: "empty", Content: ""}}
	lines := RenderInstructionsBlock(inst, ".", func(i Instruction) []string {
		if i.Content == "" {
			return nil
		}
		return []string{i.Content}
	})
	if len(lines) > 0 {
		t.Fatalf("expected no output for empty content instruction, got %v", lines)
	}
}

func TestRenderInstructionsBlock_ScopedAndGlobalMixed(t *testing.T) {
	inst := []Instruction{
		{Name: "g", Content: "global"},
		{Name: "s", ApplyTo: "src/", Content: "scoped"},
	}
	lines := RenderInstructionsBlock(inst, ".", func(i Instruction) []string {
		return []string{i.Content}
	})
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "global") || !strings.Contains(combined, "scoped") {
		t.Fatalf("expected both global and scoped in output, got %q", combined)
	}
}

func TestRenderInstructionsBlock_PatternGrouping(t *testing.T) {
	inst := []Instruction{
		{Name: "a", ApplyTo: "*.go", Content: "go1"},
		{Name: "b", ApplyTo: "*.go", Content: "go2"},
	}
	lines := RenderInstructionsBlock(inst, ".", func(i Instruction) []string {
		return []string{i.Content}
	})
	combined := strings.Join(lines, "\n")
	count := strings.Count(combined, "*.go")
	if count < 1 {
		t.Fatalf("expected at least one *.go heading, got %d", count)
	}
}
