package templatebuilder_test

import (
"strings"
"testing"

"github.com/githubnext/apm/internal/compilation/templatebuilder"
)

func emitLines(inst templatebuilder.Instruction) []string {
return []string{inst.Content, ""}
}

func TestRenderInstructionsBlock_GlobalOnly(t *testing.T) {
instructions := []templatebuilder.Instruction{
{Name: "a", Content: "line-a"},
{Name: "b", Content: "line-b"},
}
lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", emitLines)
joined := strings.Join(lines, "\n")
if !strings.Contains(joined, "## Global Instructions") {
t.Error("expected global heading")
}
if !strings.Contains(joined, "line-a") || !strings.Contains(joined, "line-b") {
t.Error("expected both instruction contents")
}
}

func TestRenderInstructionsBlock_ScopedOnly(t *testing.T) {
instructions := []templatebuilder.Instruction{
{Name: "c", Content: "scoped-content", ApplyTo: "**/*.ts"},
}
lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", emitLines)
joined := strings.Join(lines, "\n")
if !strings.Contains(joined, "## Files matching") {
t.Error("expected scoped heading")
}
if !strings.Contains(joined, "**/*.ts") {
t.Error("expected pattern in heading")
}
if strings.Contains(joined, "## Global Instructions") {
t.Error("should not have global heading when no global instructions")
}
}

func TestRenderInstructionsBlock_Empty(t *testing.T) {
lines := templatebuilder.RenderInstructionsBlock(nil, "/base", emitLines)
if len(lines) != 0 {
t.Errorf("expected empty output, got %v", lines)
}
}

func TestRenderInstructionsBlock_EmptyContentSkipped(t *testing.T) {
instructions := []templatebuilder.Instruction{
{Name: "empty", Content: ""},
{Name: "valid", Content: "hello"},
}
lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", emitLines)
count := 0
for _, l := range lines {
if l == "hello" {
count++
}
}
if count != 1 {
t.Errorf("expected 1 valid line, got %d", count)
}
}

func TestRenderInstructionsBlock_SortedPatterns(t *testing.T) {
instructions := []templatebuilder.Instruction{
{Name: "z", Content: "z-content", ApplyTo: "z-pattern"},
{Name: "a", Content: "a-content", ApplyTo: "a-pattern"},
}
lines := templatebuilder.RenderInstructionsBlock(instructions, "/base", emitLines)
joined := strings.Join(lines, "\n")
aIdx := strings.Index(joined, "a-pattern")
zIdx := strings.Index(joined, "z-pattern")
if aIdx > zIdx {
t.Error("patterns should be sorted alphabetically")
}
}
