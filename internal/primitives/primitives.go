// Package primitives defines data structures for APM primitive files
// (chatmodes, instructions, contexts, skills).
// Mirrors src/apm_cli/primitives/models.py.
package primitives

import "path/filepath"

// PrimitiveType classifies a primitive by its kind.
type PrimitiveType string

const (
	PrimitiveTypeChatmode    PrimitiveType = "chatmode"
	PrimitiveTypeInstruction PrimitiveType = "instruction"
	PrimitiveTypeContext     PrimitiveType = "context"
	PrimitiveTypeSkill       PrimitiveType = "skill"
)

// Chatmode represents a chatmode primitive.
// Mirrors src/apm_cli/primitives/models.py:Chatmode.
type Chatmode struct {
	Name        string
	FilePath    string
	Description string
	ApplyTo     string // Glob pattern for file targeting (empty if not set)
	Content     string
	Author      string
	Version     string
	Source      string
}

// Instruction represents an instruction primitive (.instructions.md).
// Mirrors src/apm_cli/primitives/models.py:Instruction.
type Instruction struct {
	Name        string
	FilePath    string
	Description string
	ApplyTo     string
	Content     string
	Author      string
	Version     string
	Source      string
}

// Context represents a context primitive (.context.md).
// Mirrors src/apm_cli/primitives/models.py:Context.
type Context struct {
	Name        string
	FilePath    string
	Description string
	Scope       string
	Content     string
	Author      string
	Version     string
	Source      string
}

// Skill represents a skill primitive (SKILL.md).
// Mirrors src/apm_cli/primitives/models.py:Skill.
type Skill struct {
	Name        string
	FilePath    string
	Description string
	Content     string
	Author      string
	Version     string
	Source      string
}

// PrimitiveConflict records a conflict between two primitives.
// Mirrors src/apm_cli/primitives/models.py:PrimitiveConflict.
type PrimitiveConflict struct {
	Type     PrimitiveType
	Name     string
	Path1    string
	Path2    string
	Reason   string
}

// PrimitiveCollection holds all discovered primitives for a package.
// Mirrors src/apm_cli/primitives/models.py:PrimitiveCollection.
type PrimitiveCollection struct {
	Chatmodes    []Chatmode
	Instructions []Instruction
	Contexts     []Context
	Skills       []Skill
	Conflicts    []PrimitiveConflict
}

// FileNameWithoutExt returns the base filename without extension.
func FileNameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return base[:len(base)-len(ext)]
}

// TotalCount returns total number of primitives in the collection.
func (pc *PrimitiveCollection) TotalCount() int {
	return len(pc.Chatmodes) + len(pc.Instructions) + len(pc.Contexts) + len(pc.Skills)
}

// HasConflicts returns true if there are any conflicts.
func (pc *PrimitiveCollection) HasConflicts() bool {
	return len(pc.Conflicts) > 0
}
