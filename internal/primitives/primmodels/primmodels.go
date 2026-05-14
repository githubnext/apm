// Package primmodels defines data models for APM primitives.
package primmodels

// Primitive is the common interface for all APM primitive types.
type Primitive interface {
	Validate() []string
}

// Chatmode represents a chatmode primitive.
type Chatmode struct {
Name        string
FilePath    string
Description string
ApplyTo     string
Content     string
Author      string
Version     string
Source      string
}

// Validate returns a list of validation errors for a Chatmode.
func (c *Chatmode) Validate() []string {
var errs []string
if c.Description == "" {
errs = append(errs, "Missing 'description' in frontmatter")
}
if c.Content == "" {
errs = append(errs, "Empty content")
}
return errs
}

// Instruction represents an instruction primitive.
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

// Validate returns a list of validation errors for an Instruction.
func (i *Instruction) Validate() []string {
var errs []string
if i.Description == "" {
errs = append(errs, "Missing 'description' in frontmatter")
}
if i.Content == "" {
errs = append(errs, "Empty content")
}
return errs
}

// Context represents a context primitive.
type Context struct {
Name        string
FilePath    string
Content     string
Description string
Author      string
Version     string
Source      string
}

// Validate returns validation errors for a Context.
func (c *Context) Validate() []string {
if c.Content == "" {
return []string{"Empty content"}
}
return nil
}

// Skill represents a skill primitive.
type Skill struct {
Name        string
FilePath    string
Description string
ApplyTo     string
Content     string
Author      string
Version     string
Source      string
}

// Validate returns validation errors for a Skill.
func (s *Skill) Validate() []string {
return nil
}

// Agent represents an agent primitive.
type Agent struct {
Name        string
FilePath    string
Description string
Content     string
Author      string
Version     string
Source      string
}

// Hook represents a hook primitive.
type Hook struct {
Name        string
FilePath    string
Description string
Content     string
Author      string
Version     string
Source      string
}

// ConflictIndex tracks primitives by name to detect conflicts.
type ConflictIndex struct {
Chatmodes    map[string]*Chatmode
Instructions map[string]*Instruction
Skills       map[string]*Skill
Agents       map[string]*Agent
}

// NewConflictIndex creates an initialized ConflictIndex.
func NewConflictIndex() *ConflictIndex {
return &ConflictIndex{
Chatmodes:    map[string]*Chatmode{},
Instructions: map[string]*Instruction{},
Skills:       map[string]*Skill{},
Agents:       map[string]*Agent{},
}
}
