package apmyml_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/apmyml"
)

func TestParseTargetsField_BothKeys_Error(t *testing.T) {
	data := map[string]interface{}{
		"targets": []interface{}{"claude"},
		"target":  "copilot",
	}
	_, err := apmyml.ParseTargetsField(data)
	if err == nil {
		t.Error("expected error when both 'targets' and 'target' are present")
	}
}

func TestParseTargetsField_NeitherKey(t *testing.T) {
	data := map[string]interface{}{}
	got, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty list, got %v", got)
	}
}

func TestParseTargetsField_EmptyList_Error(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{}}
	_, err := apmyml.ParseTargetsField(data)
	if err == nil {
		t.Error("expected error for empty targets list")
	}
}

func TestParseTargetsField_CSVSingular(t *testing.T) {
	data := map[string]interface{}{"target": "claude,copilot"}
	got, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 targets from CSV, got %v", got)
	}
}

func TestParseTargetsField_ListUnderSingularKey(t *testing.T) {
	data := map[string]interface{}{"target": []interface{}{"claude", "cursor"}}
	got, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 from list-under-singular, got %v", got)
	}
}

func TestCanonicalTargets_ContainsExpected(t *testing.T) {
	expected := []string{"claude", "copilot", "cursor", "codex", "gemini", "windsurf"}
	for _, tgt := range expected {
		if !apmyml.CanonicalTargets[tgt] {
			t.Errorf("CanonicalTargets missing %q", tgt)
		}
	}
}

func TestParseTargetsField_UnknownTarget_Error(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{"notavalidtarget"}}
	_, err := apmyml.ParseTargetsField(data)
	if err == nil {
		t.Error("expected error for unknown target 'notavalidtarget'")
	}
}

func TestConflictingTargetsError_Message(t *testing.T) {
	var e error = &apmyml.ConflictingTargetsError{Message: "conflict"}
	if e.Error() != "conflict" {
		t.Errorf("Error() = %q, want 'conflict'", e.Error())
	}
}

func TestEmptyTargetsListError_Message(t *testing.T) {
	var e error = &apmyml.EmptyTargetsListError{Message: "empty"}
	if e.Error() != "empty" {
		t.Errorf("Error() = %q, want 'empty'", e.Error())
	}
}

func TestParseTargetsField_AgentSkills(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{"agent-skills"}}
	got, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "agent-skills" {
		t.Errorf("expected [agent-skills], got %v", got)
	}
}
