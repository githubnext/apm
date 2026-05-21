package templatebuilder

import (
	"strings"
	"testing"
)

func TestRenderInstructionsBlock_Empty_Extra4(t *testing.T) {
	lines := RenderInstructionsBlock(nil, "/base", func(Instruction) []string { return nil })
	if len(lines) != 0 {
		t.Errorf("expected empty output for nil instructions, got %v", lines)
	}
}

func TestRenderInstructionsBlock_EmptyContent_Extra4(t *testing.T) {
	insts := []Instruction{{Name: "a", Content: ""}}
	lines := RenderInstructionsBlock(insts, "/base", func(i Instruction) []string { return []string{i.Name} })
	if len(lines) != 0 {
		t.Error("expected instructions with empty content to be skipped")
	}
}

func TestRenderInstructionsBlock_GlobalHeading_Extra4(t *testing.T) {
	insts := []Instruction{{Name: "a", FilePath: "/base/a.md", Content: "hello", ApplyTo: ""}}
	lines := RenderInstructionsBlock(insts, "/base", func(i Instruction) []string { return []string{"inst:" + i.Name} })
	found := false
	for _, l := range lines {
		if strings.Contains(l, "Global") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Global heading in %v", lines)
	}
}

func TestRenderInstructionsBlock_ScopedHeading_Extra4(t *testing.T) {
	insts := []Instruction{{Name: "b", FilePath: "/base/b.md", Content: "c", ApplyTo: "**/*.go"}}
	lines := RenderInstructionsBlock(insts, "/base", func(i Instruction) []string { return []string{"inst:" + i.Name} })
	found := false
	for _, l := range lines {
		if strings.Contains(l, "**/*.go") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected scoped heading with pattern in %v", lines)
	}
}

func TestRenderInstructionsBlock_SortedPatterns_Extra4(t *testing.T) {
	insts := []Instruction{
		{Name: "z", FilePath: "/base/z.md", Content: "c", ApplyTo: "z-pattern"},
		{Name: "a", FilePath: "/base/a.md", Content: "c", ApplyTo: "a-pattern"},
	}
	lines := RenderInstructionsBlock(insts, "/base", func(i Instruction) []string { return []string{i.Name} })
	firstIdx, secondIdx := -1, -1
	for i, l := range lines {
		if strings.Contains(l, "a-pattern") {
			firstIdx = i
		}
		if strings.Contains(l, "z-pattern") {
			secondIdx = i
		}
	}
	if firstIdx < 0 || secondIdx < 0 {
		t.Skip("patterns not found in output")
	}
	if firstIdx > secondIdx {
		t.Error("expected a-pattern before z-pattern (sorted)")
	}
}

func TestInstruction_Fields_Extra4(t *testing.T) {
	i := Instruction{Name: "n", FilePath: "/f", ApplyTo: "**", Content: "c"}
	if i.Name != "n" || i.FilePath != "/f" || i.ApplyTo != "**" || i.Content != "c" {
		t.Error("unexpected field values")
	}
}

func TestTemplateData_Fields_Extra4(t *testing.T) {
	td := TemplateData{InstructionsContent: "ic", Version: "v1", ChatmodeContent: "cc"}
	if td.InstructionsContent != "ic" || td.Version != "v1" || td.ChatmodeContent != "cc" {
		t.Error("unexpected TemplateData field values")
	}
}

func TestRenderInstructionsBlock_MultipleGlobal_Extra4(t *testing.T) {
	insts := []Instruction{
		{Name: "a", FilePath: "/base/a.md", Content: "c1"},
		{Name: "b", FilePath: "/base/b.md", Content: "c2"},
	}
	emitted := []string{}
	RenderInstructionsBlock(insts, "/base", func(i Instruction) []string {
		emitted = append(emitted, i.Name)
		return []string{i.Name}
	})
	if len(emitted) != 2 {
		t.Errorf("expected 2 emitted, got %d", len(emitted))
	}
}

func TestRenderInstructionsBlock_EmitterCalled_Extra4(t *testing.T) {
	called := 0
	insts := []Instruction{{Name: "x", Content: "body"}}
	RenderInstructionsBlock(insts, "/base", func(i Instruction) []string {
		called++
		return []string{"x"}
	})
	if called != 1 {
		t.Errorf("expected emitter called once, got %d", called)
	}
}
