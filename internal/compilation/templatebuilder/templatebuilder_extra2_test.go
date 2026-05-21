package templatebuilder

import (
	"strings"
	"testing"
)

func TestInstruction_FieldZeroValues(t *testing.T) {
	var inst Instruction
	if inst.Name != "" || inst.FilePath != "" || inst.ApplyTo != "" || inst.Content != "" {
		t.Error("zero-value Instruction fields should be empty strings")
	}
}

func TestTemplateData_FieldZeroValues(t *testing.T) {
	var td TemplateData
	if td.InstructionsContent != "" || td.Version != "" || td.ChatmodeContent != "" {
		t.Error("zero-value TemplateData fields should be empty strings")
	}
}

func TestInstruction_AssignFields(t *testing.T) {
	inst := Instruction{
		Name:     "my-inst",
		FilePath: "/path/to/file.md",
		ApplyTo:  "**/*.go",
		Content:  "some content",
	}
	if inst.Name != "my-inst" {
		t.Errorf("expected Name=my-inst got %q", inst.Name)
	}
	if inst.ApplyTo != "**/*.go" {
		t.Errorf("expected ApplyTo=**/*.go got %q", inst.ApplyTo)
	}
}

func TestTemplateData_AssignFields(t *testing.T) {
	td := TemplateData{
		InstructionsContent: "ic",
		Version:             "v1.2.3",
		ChatmodeContent:     "cm",
	}
	if td.Version != "v1.2.3" {
		t.Errorf("expected Version=v1.2.3 got %q", td.Version)
	}
}

func TestRenderInstructionsBlock_SingleGlobalHeading(t *testing.T) {
	insts := []Instruction{
		{Name: "a", FilePath: "/base/a.md", Content: "line1"},
	}
	emit := func(i Instruction) []string { return []string{i.Content} }
	lines := RenderInstructionsBlock(insts, "/base", emit)
	found := false
	for _, l := range lines {
		if l == "## Global Instructions" {
			found = true
		}
	}
	if !found {
		t.Error("expected Global Instructions heading")
	}
}

func TestRenderInstructionsBlock_ScopedHeadingContainsPattern(t *testing.T) {
	insts := []Instruction{
		{Name: "b", FilePath: "/base/b.md", ApplyTo: "src/**", Content: "x"},
	}
	emit := func(i Instruction) []string { return []string{i.Content} }
	lines := RenderInstructionsBlock(insts, "/base", emit)
	found := false
	for _, l := range lines {
		if strings.Contains(l, "src/**") {
			found = true
		}
	}
	if !found {
		t.Error("expected scoped heading containing src/**")
	}
}

func TestRenderInstructionsBlock_MultipleEmitLines(t *testing.T) {
	insts := []Instruction{
		{Name: "c", FilePath: "/base/c.md", Content: "multi"},
	}
	emit := func(i Instruction) []string { return []string{"line1", "line2", "line3"} }
	lines := RenderInstructionsBlock(insts, "/base", emit)
	count := 0
	for _, l := range lines {
		if l == "line1" || l == "line2" || l == "line3" {
			count++
		}
	}
	if count != 3 {
		t.Errorf("expected 3 emitted lines, got %d", count)
	}
}

func TestRenderInstructionsBlock_PatternHeadingFormat(t *testing.T) {
	insts := []Instruction{
		{Name: "d", FilePath: "/base/d.md", ApplyTo: "*.py", Content: "py-content"},
	}
	emit := func(i Instruction) []string { return []string{i.Content} }
	lines := RenderInstructionsBlock(insts, "/base", emit)
	expected := "## Files matching `*.py`"
	found := false
	for _, l := range lines {
		if l == expected {
			found = true
		}
	}
	if !found {
		t.Errorf("expected heading %q", expected)
	}
}

func TestRenderInstructionsBlock_NoGlobalHeadingWhenNoGlobal(t *testing.T) {
	insts := []Instruction{
		{Name: "e", FilePath: "/base/e.md", ApplyTo: "*.go", Content: "go"},
	}
	emit := func(i Instruction) []string { return []string{i.Content} }
	lines := RenderInstructionsBlock(insts, "/base", emit)
	for _, l := range lines {
		if l == "## Global Instructions" {
			t.Error("should not emit Global Instructions heading when no global insts")
		}
	}
}

func TestRenderInstructionsBlock_BlankLineSeparator(t *testing.T) {
	insts := []Instruction{
		{Name: "f", FilePath: "/base/f.md", Content: "global"},
	}
	emit := func(i Instruction) []string { return []string{i.Content} }
	lines := RenderInstructionsBlock(insts, "/base", emit)
	// Should have at least one blank line after heading
	hasBlank := false
	for _, l := range lines {
		if l == "" {
			hasBlank = true
		}
	}
	if !hasBlank {
		t.Error("expected at least one blank separator line")
	}
}

func TestRenderInstructionsBlock_TwoScopedTwoDifferentPatterns(t *testing.T) {
	insts := []Instruction{
		{Name: "g1", FilePath: "/b/g1.md", ApplyTo: "a/**", Content: "ca"},
		{Name: "g2", FilePath: "/b/g2.md", ApplyTo: "z/**", Content: "cz"},
	}
	emit := func(i Instruction) []string { return []string{i.Content} }
	lines := RenderInstructionsBlock(insts, "/b", emit)
	hasA, hasZ := false, false
	for _, l := range lines {
		if strings.Contains(l, "a/**") {
			hasA = true
		}
		if strings.Contains(l, "z/**") {
			hasZ = true
		}
	}
	if !hasA || !hasZ {
		t.Error("expected headings for both patterns a/** and z/**")
	}
}
