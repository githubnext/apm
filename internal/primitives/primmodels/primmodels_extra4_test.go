package primmodels_test

import (
"testing"

"github.com/githubnext/apm/internal/primitives/primmodels"
)

func TestChatmode_SourceField_Extra4(t *testing.T) {
c := primmodels.Chatmode{Source: "mysource"}
if c.Source != "mysource" {
t.Errorf("expected 'mysource', got %q", c.Source)
}
}

func TestChatmode_FilePath_Extra4(t *testing.T) {
c := primmodels.Chatmode{FilePath: "/a/b/c.md"}
if c.FilePath != "/a/b/c.md" {
t.Errorf("expected '/a/b/c.md', got %q", c.FilePath)
}
}

func TestInstruction_FilePath_Extra4(t *testing.T) {
i := primmodels.Instruction{FilePath: "/x/y.md"}
if i.FilePath != "/x/y.md" {
t.Errorf("expected '/x/y.md', got %q", i.FilePath)
}
}

func TestInstruction_ApplyTo_Extra4(t *testing.T) {
i := primmodels.Instruction{ApplyTo: "**/*.go"}
if i.ApplyTo != "**/*.go" {
t.Errorf("expected '**/*.go', got %q", i.ApplyTo)
}
}

func TestContext_NameField_Extra4(t *testing.T) {
c := primmodels.Context{Name: "ctx-name"}
if c.Name != "ctx-name" {
t.Errorf("expected 'ctx-name', got %q", c.Name)
}
}

func TestSkill_NameAndDescription_Extra4(t *testing.T) {
s := primmodels.Skill{Name: "mskill", Description: "desc text"}
if s.Name != "mskill" {
t.Errorf("expected 'mskill', got %q", s.Name)
}
if s.Description != "desc text" {
t.Errorf("expected 'desc text', got %q", s.Description)
}
}

func TestNewConflictIndex_Empty_Extra4(t *testing.T) {
ci := primmodels.NewConflictIndex()
if ci == nil {
t.Error("expected non-nil ConflictIndex")
}
}

func TestConflictIndex_SkillsMapInitialized_Extra4(t *testing.T) {
ci := primmodels.NewConflictIndex()
if ci.Skills == nil {
t.Error("expected Skills map to be initialized")
}
}

func TestAgent_NameField_Extra4(t *testing.T) {
a := primmodels.Agent{Name: "agent-x"}
if a.Name != "agent-x" {
t.Errorf("expected 'agent-x', got %q", a.Name)
}
}

func TestHook_NameField_Extra4(t *testing.T) {
h := primmodels.Hook{Name: "pre-push"}
if h.Name != "pre-push" {
t.Errorf("expected 'pre-push', got %q", h.Name)
}
}
