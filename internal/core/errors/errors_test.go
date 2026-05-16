package errors

import (
	"strings"
	"testing"
)

func TestTargetResolutionError(t *testing.T) {
	e := &TargetResolutionError{Message: "test error"}
	if e.Error() != "test error" {
		t.Errorf("Error() = %q, want %q", e.Error(), "test error")
	}
}

func TestRenderNoHarnessError(t *testing.T) {
	out := RenderNoHarnessError()
	if !strings.Contains(out, "[x]") {
		t.Error("RenderNoHarnessError: missing [x] prefix")
	}
	if !strings.Contains(out, "apm install") {
		t.Error("RenderNoHarnessError: missing apm install suggestion")
	}
	if !strings.Contains(out, "targets:") {
		t.Error("RenderNoHarnessError: missing targets: yaml example")
	}
}

func TestRenderAmbiguousError(t *testing.T) {
	out := RenderAmbiguousError([]string{"claude", "copilot"})
	if !strings.Contains(out, "[x]") {
		t.Error("RenderAmbiguousError: missing [x] prefix")
	}
	if !strings.Contains(out, "claude") {
		t.Error("RenderAmbiguousError: missing detected harnesses")
	}
	if !strings.Contains(out, "copilot") {
		t.Error("RenderAmbiguousError: missing detected harnesses")
	}
}

func TestRenderAmbiguousError_Empty(t *testing.T) {
	out := RenderAmbiguousError([]string{})
	if !strings.Contains(out, "[x]") {
		t.Error("RenderAmbiguousError empty: missing [x] prefix")
	}
}

func TestRenderUnknownTargetError(t *testing.T) {
	valid := []string{"claude", "copilot", "cursor", "agent-skills"}
	out := RenderUnknownTargetError("notreal", valid)
	if !strings.Contains(out, "[x]") {
		t.Error("RenderUnknownTargetError: missing [x] prefix")
	}
	if !strings.Contains(out, "notreal") {
		t.Error("RenderUnknownTargetError: missing unknown value")
	}
	// agent-skills should be hidden
	if strings.Contains(out, "agent-skills") {
		t.Error("RenderUnknownTargetError: agent-skills should not appear in valid list")
	}
}

func TestRenderUnknownTargetError_Empty(t *testing.T) {
	out := RenderUnknownTargetError("", []string{})
	if !strings.Contains(out, "[x]") {
		t.Error("missing [x] prefix for empty target")
	}
}

func TestRenderConflictingSchemaError(t *testing.T) {
	out := RenderConflictingSchemaError()
	if !strings.Contains(out, "[x]") {
		t.Error("RenderConflictingSchemaError: missing [x] prefix")
	}
	if !strings.Contains(out, "targets:") {
		t.Error("RenderConflictingSchemaError: missing targets: reference")
	}
}
