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

func TestTargetResolutionError_Types(t *testing.T) {
	var e error = &NoHarnessError{TargetResolutionError{Message: "no harness"}}
	if e.Error() != "no harness" {
		t.Errorf("NoHarnessError: got %q", e.Error())
	}
	var e2 error = &AmbiguousHarnessError{TargetResolutionError{Message: "ambiguous"}}
	if e2.Error() != "ambiguous" {
		t.Errorf("AmbiguousHarnessError: got %q", e2.Error())
	}
	var e3 error = &UnknownTargetError{TargetResolutionError{Message: "unknown target"}}
	if e3.Error() != "unknown target" {
		t.Errorf("UnknownTargetError: got %q", e3.Error())
	}
	var e4 error = &ConflictingTargetsError{TargetResolutionError{Message: "conflict"}}
	if e4.Error() != "conflict" {
		t.Errorf("ConflictingTargetsError: got %q", e4.Error())
	}
	var e5 error = &EmptyTargetsListError{TargetResolutionError{Message: "empty"}}
	if e5.Error() != "empty" {
		t.Errorf("EmptyTargetsListError: got %q", e5.Error())
	}
}

func TestRenderAmbiguousError_Suggestion(t *testing.T) {
	out := RenderAmbiguousError([]string{"cursor"})
	if !strings.Contains(out, "cursor") {
		t.Error("RenderAmbiguousError: missing cursor in output")
	}
	if !strings.Contains(out, "--target cursor") {
		t.Error("RenderAmbiguousError: missing suggestion")
	}
}

func TestRenderNoHarnessError_ContainsMarkers(t *testing.T) {
	out := RenderNoHarnessError()
	if !strings.Contains(out, ".claude/") {
		t.Error("RenderNoHarnessError: missing .claude/ marker")
	}
	if !strings.Contains(out, "--target") {
		t.Error("RenderNoHarnessError: missing --target suggestion")
	}
	if !strings.Contains(out, "apm install") {
		t.Error("RenderNoHarnessError: missing apm install command")
	}
}

func TestRenderUnknownTargetError_ShowsValid(t *testing.T) {
	valid := []string{"claude", "cursor", "gemini"}
	out := RenderUnknownTargetError("bogus", valid)
	if !strings.Contains(out, "bogus") {
		t.Error("RenderUnknownTargetError: missing unknown value in output")
	}
	if !strings.Contains(out, "claude") {
		t.Error("RenderUnknownTargetError: missing valid target claude")
	}
	if !strings.Contains(out, "cursor") {
		t.Error("RenderUnknownTargetError: missing valid target cursor")
	}
}

func TestRenderUnknownTargetError_BracketInput(t *testing.T) {
	out := RenderUnknownTargetError("['badval']", []string{"claude"})
	if !strings.Contains(out, "badval") {
		t.Errorf("expected cleaned-up value in output, got: %s", out)
	}
}
