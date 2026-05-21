package apmyml_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core/apmyml"
)

func TestParseTargetsField_NilMap_Extra4(t *testing.T) {
	_, err := apmyml.ParseTargetsField(nil)
	if err != nil {
		t.Errorf("expected no error for nil map, got %v", err)
	}
}

func TestParseTargetsField_EmptyMap_Extra4(t *testing.T) {
	_, err := apmyml.ParseTargetsField(map[string]interface{}{})
	if err != nil {
		t.Errorf("expected no error for empty map, got %v", err)
	}
}

func TestParseTargetsField_ValidTarget_Extra4(t *testing.T) {
	m := map[string]interface{}{"targets": []interface{}{"copilot"}}
	targets, err := apmyml.ParseTargetsField(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) == 0 {
		t.Error("expected at least one target")
	}
}

func TestParseTargetsField_UnknownTarget_Error_Extra4(t *testing.T) {
	m := map[string]interface{}{"targets": []interface{}{"definitely_unknown_target_xyz"}}
	_, err := apmyml.ParseTargetsField(m)
	if err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestConflictingTargetsError_Message_Extra4(t *testing.T) {
	err := &apmyml.ConflictingTargetsError{Message: "targets and target conflict"}
	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestEmptyTargetsListError_Message_Extra4(t *testing.T) {
	err := &apmyml.EmptyTargetsListError{Message: "targets list is empty"}
	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestUnknownTargetError_Message_Extra4(t *testing.T) {
	err := &apmyml.UnknownTargetError{Token: "badtarget", Message: "unknown target: badtarget"}
	msg := err.Error()
	if !strings.Contains(msg, "badtarget") {
		t.Errorf("expected target name in error, got %q", msg)
	}
}

func TestCanonicalTargets_NotEmpty_Extra4(t *testing.T) {
	if len(apmyml.CanonicalTargets) == 0 {
		t.Error("expected non-empty CanonicalTargets")
	}
}

func TestCanonicalTargets_CopilotPresent_Extra4(t *testing.T) {
	if !apmyml.CanonicalTargets["copilot"] {
		t.Error("expected copilot in CanonicalTargets")
	}
}
