package apmyml

import (
	"testing"
)

func TestCanonicalTargets_NotEmpty(t *testing.T) {
	if len(CanonicalTargets) == 0 {
		t.Fatal("CanonicalTargets should not be empty")
	}
}

func TestCanonicalTargets_ContainsGemini(t *testing.T) {
	if !CanonicalTargets["gemini"] {
		t.Fatal("expected 'gemini' in CanonicalTargets")
	}
}

func TestCanonicalTargets_ContainsCodex(t *testing.T) {
	if !CanonicalTargets["codex"] {
		t.Fatal("expected 'codex' in CanonicalTargets")
	}
}

func TestCanonicalTargets_ContainsWindsurf(t *testing.T) {
	if !CanonicalTargets["windsurf"] {
		t.Fatal("expected 'windsurf' in CanonicalTargets")
	}
}

func TestParseTargetsField_PluralList(t *testing.T) {
	data := map[string]interface{}{
		"targets": []interface{}{"claude", "copilot"},
	}
	targets, err := ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
}

func TestParseTargetsField_SingularValue(t *testing.T) {
	data := map[string]interface{}{
		"target": "cursor",
	}
	targets, err := ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "cursor" {
		t.Fatalf("expected [cursor], got %v", targets)
	}
}

func TestParseTargetsField_NeitherKeyReturnsEmpty(t *testing.T) {
	data := map[string]interface{}{}
	targets, err := ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 0 {
		t.Fatalf("expected empty, got %v", targets)
	}
}

func TestParseTargetsField_BothKeyConflict(t *testing.T) {
	data := map[string]interface{}{
		"target":  "claude",
		"targets": []interface{}{"copilot"},
	}
	_, err := ParseTargetsField(data)
	if err == nil {
		t.Fatal("expected error for both keys")
	}
}

func TestConflictingTargetsError_Message(t *testing.T) {
	e := &ConflictingTargetsError{Message: "conflicting targets"}
	if e.Error() != "conflicting targets" {
		t.Fatalf("unexpected error: %q", e.Error())
	}
}

func TestEmptyTargetsListError_Message(t *testing.T) {
	e := &EmptyTargetsListError{Message: "empty list"}
	if e.Error() != "empty list" {
		t.Fatalf("unexpected error: %q", e.Error())
	}
}

func TestUnknownTargetError_MessageContainsTarget(t *testing.T) {
	e := &UnknownTargetError{Token: "unknown-tool", Message: "unknown target: unknown-tool"}
	msg := e.Error()
	if msg == "" {
		t.Fatal("error message should not be empty")
	}
}

func TestParseTargetsField_AgentSkillsCanonical(t *testing.T) {
	data := map[string]interface{}{
		"target": "agent-skills",
	}
	targets, err := ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "agent-skills" {
		t.Fatalf("expected [agent-skills], got %v", targets)
	}
}

func TestParseTargetsField_UnknownTargetError(t *testing.T) {
	data := map[string]interface{}{
		"target": "nonexistent-tool",
	}
	_, err := ParseTargetsField(data)
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}
