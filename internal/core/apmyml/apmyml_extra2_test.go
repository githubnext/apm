package apmyml_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core/apmyml"
)

func TestParseTargetsField_SingularTarget_Simple(t *testing.T) {
	data := map[string]interface{}{"target": "claude"}
	targets, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "claude" {
		t.Errorf("expected [claude], got %v", targets)
	}
}

func TestParseTargetsField_PluralTargets_TwoItems(t *testing.T) {
	data := map[string]interface{}{
		"targets": []interface{}{"claude", "copilot"},
	}
	targets, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %v", targets)
	}
}

func TestParseTargetsField_NeitherKey_ReturnsEmpty(t *testing.T) {
	data := map[string]interface{}{"other": "value"}
	targets, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 0 {
		t.Errorf("expected empty targets, got %v", targets)
	}
}

func TestParseTargetsField_BothKeys_ReturnsError(t *testing.T) {
	data := map[string]interface{}{
		"target":  "claude",
		"targets": []interface{}{"copilot"},
	}
	_, err := apmyml.ParseTargetsField(data)
	if err == nil {
		t.Fatal("expected error for both keys present")
	}
	if _, ok := err.(*apmyml.ConflictingTargetsError); !ok {
		t.Errorf("expected ConflictingTargetsError, got %T", err)
	}
}

func TestParseTargetsField_EmptyTargetsList_ReturnsError(t *testing.T) {
	data := map[string]interface{}{
		"targets": []interface{}{},
	}
	_, err := apmyml.ParseTargetsField(data)
	if err == nil {
		t.Fatal("expected error for empty targets list")
	}
	if _, ok := err.(*apmyml.EmptyTargetsListError); !ok {
		t.Errorf("expected EmptyTargetsListError, got %T", err)
	}
}

func TestParseTargetsField_UnknownTarget_ReturnsError(t *testing.T) {
	data := map[string]interface{}{"target": "unknown-editor"}
	_, err := apmyml.ParseTargetsField(data)
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
	if _, ok := err.(*apmyml.UnknownTargetError); !ok {
		t.Errorf("expected UnknownTargetError, got %T", err)
	}
}

func TestCanonicalTargets_ContainsClaude(t *testing.T) {
	if !apmyml.CanonicalTargets["claude"] {
		t.Error("CanonicalTargets should contain 'claude'")
	}
}

func TestCanonicalTargets_ContainsCopilot(t *testing.T) {
	if !apmyml.CanonicalTargets["copilot"] {
		t.Error("CanonicalTargets should contain 'copilot'")
	}
}

func TestCanonicalTargets_ContainsCursor(t *testing.T) {
	if !apmyml.CanonicalTargets["cursor"] {
		t.Error("CanonicalTargets should contain 'cursor'")
	}
}

func TestConflictingTargetsError_ErrorMessage(t *testing.T) {
	err := &apmyml.ConflictingTargetsError{Message: "conflict msg"}
	if err.Error() != "conflict msg" {
		t.Errorf("expected 'conflict msg', got %q", err.Error())
	}
}

func TestEmptyTargetsListError_ErrorMessage(t *testing.T) {
	err := &apmyml.EmptyTargetsListError{Message: "empty targets"}
	if err.Error() != "empty targets" {
		t.Errorf("expected 'empty targets', got %q", err.Error())
	}
}

func TestUnknownTargetError_ErrorMessage(t *testing.T) {
	err := &apmyml.UnknownTargetError{Token: "bad", Message: "bad target"}
	if err.Error() != "bad target" {
		t.Errorf("expected 'bad target', got %q", err.Error())
	}
	if err.Token != "bad" {
		t.Errorf("expected Token='bad', got %q", err.Token)
	}
}

func TestParseTargetsField_CSV_Trimmed(t *testing.T) {
	data := map[string]interface{}{"target": " claude , copilot "}
	targets, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, tgt := range targets {
		if strings.TrimSpace(tgt) != tgt {
			t.Errorf("target should be trimmed, got %q", tgt)
		}
	}
}

func TestParseTargetsField_Windsurf_IsCanonical(t *testing.T) {
	data := map[string]interface{}{"target": "windsurf"}
	targets, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "windsurf" {
		t.Errorf("expected [windsurf], got %v", targets)
	}
}

func TestParseTargetsField_Gemini_IsCanonical(t *testing.T) {
	data := map[string]interface{}{"target": "gemini"}
	targets, err := apmyml.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 || targets[0] != "gemini" {
		t.Errorf("expected [gemini], got %v", targets)
	}
}
