// Package templatebuilder provides template building utilities for AGENTS.md compilation.
package templatebuilder

import (
"path/filepath"
"sort"
"strings"
)

// Instruction represents an instruction primitive for template rendering.
type Instruction struct {
Name     string
FilePath string
ApplyTo  string
Content  string
}

// TemplateData holds data for template generation.
type TemplateData struct {
InstructionsContent string
Version             string
ChatmodeContent     string
}

const globalInstructionsHeading = "## Global Instructions"

// RenderInstructionsBlock renders the body lines of an instructions section.
// Global instructions (no ApplyTo) go under globalInstructionsHeading.
// Pattern-scoped instructions are grouped under "## Files matching `<pattern>`" headings.
func RenderInstructionsBlock(instructions []Instruction, baseDir string, emitInstruction func(Instruction) []string) []string {
var global []Instruction
scoped := map[string][]Instruction{}

for _, inst := range instructions {
if inst.Content == "" {
continue
}
if inst.ApplyTo == "" {
global = append(global, inst)
} else {
scoped[inst.ApplyTo] = append(scoped[inst.ApplyTo], inst)
}
}

// Sort global instructions by relative path
sort.Slice(global, func(i, j int) bool {
return relKey(baseDir, global[i].FilePath) < relKey(baseDir, global[j].FilePath)
})

var lines []string

if len(global) > 0 {
lines = append(lines, globalInstructionsHeading)
lines = append(lines, "")
for _, inst := range global {
lines = append(lines, emitInstruction(inst)...)
}
}

// Sort patterns for deterministic output
var patterns []string
for p := range scoped {
patterns = append(patterns, p)
}
sort.Strings(patterns)

for _, pattern := range patterns {
insts := scoped[pattern]
sort.Slice(insts, func(i, j int) bool {
return relKey(baseDir, insts[i].FilePath) < relKey(baseDir, insts[j].FilePath)
})
lines = append(lines, "## Files matching `"+pattern+"`")
lines = append(lines, "")
for _, inst := range insts {
lines = append(lines, emitInstruction(inst)...)
}
}

return lines
}

func relKey(base, path string) string {
rel, err := filepath.Rel(base, path)
if err != nil {
return path
}
return strings.ToLower(rel)
}
