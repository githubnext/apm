package errors

import (
	"strings"
	"testing"
)

func TestRenderNoHarnessError_ContainsSignals(t *testing.T) {
	out := RenderNoHarnessError()
	if !strings.Contains(out, ".claude/") {
		t.Error("RenderNoHarnessError: missing .claude/ signal")
	}
	if !strings.Contains(out, ".cursor/") {
		t.Error("RenderNoHarnessError: missing .cursor/ signal")
	}
	if !strings.Contains(out, "copilot-instructions.md") {
		t.Error("RenderNoHarnessError: missing copilot signal")
	}
}

func TestRenderNoHarnessError_NoExtraXPrefix(t *testing.T) {
	out := RenderNoHarnessError()
	lines := strings.Split(out, "\n")
	if lines[0][:3] != "[x]" {
		t.Errorf("first line should start with [x], got %q", lines[0])
	}
}

func TestRenderAmbiguousError_ContainsDetected(t *testing.T) {
	out := RenderAmbiguousError([]string{"claude", "copilot"})
	if !strings.Contains(out, "claude") {
		t.Error("expected 'claude' in output")
	}
	if !strings.Contains(out, "copilot") {
		t.Error("expected 'copilot' in output")
	}
}

func TestRenderAmbiguousError_SingleDetected(t *testing.T) {
	out := RenderAmbiguousError([]string{"cursor"})
	if !strings.Contains(out, "cursor") {
		t.Errorf("expected 'cursor' in output, got %q", out)
	}
	if !strings.Contains(out, "[x]") {
		t.Errorf("expected [x] prefix, got %q", out)
	}
}

func TestRenderAmbiguousError_EmptySlice(t *testing.T) {
	out := RenderAmbiguousError([]string{})
	if !strings.Contains(out, "[x]") {
		t.Errorf("expected [x] prefix even with empty list")
	}
}

func TestRenderUnknownTargetError_ContainsValid(t *testing.T) {
	out := RenderUnknownTargetError("notreal", []string{"claude", "copilot", "cursor"})
	if !strings.Contains(out, "notreal") {
		t.Errorf("expected bad value in output, got %q", out)
	}
	if !strings.Contains(out, "claude") {
		t.Errorf("expected 'claude' in valid targets, got %q", out)
	}
}

func TestRenderUnknownTargetError_FiltersAgentSkills(t *testing.T) {
	out := RenderUnknownTargetError("bad", []string{"claude", "agent-skills"})
	if strings.Contains(out, "agent-skills") {
		t.Errorf("agent-skills should be filtered from visible targets, got %q", out)
	}
}

func TestRenderConflictingSchemaError_HasFix(t *testing.T) {
	out := RenderConflictingSchemaError()
	if !strings.Contains(out, "[x]") {
		t.Errorf("expected [x] prefix, got %q", out)
	}
	if !strings.Contains(out, "targets:") {
		t.Errorf("expected 'targets:' fix, got %q", out)
	}
}

func TestErrorTypes_Hierarchy(t *testing.T) {
	var e error = &NoHarnessError{}
	if e == nil {
		t.Error("NoHarnessError should implement error")
	}
	var e2 error = &AmbiguousHarnessError{}
	if e2 == nil {
		t.Error("AmbiguousHarnessError should implement error")
	}
	var e3 error = &UnknownTargetError{}
	if e3 == nil {
		t.Error("UnknownTargetError should implement error")
	}
}

func TestTargetResolutionError_Message(t *testing.T) {
	e := &TargetResolutionError{Message: "custom msg"}
	if e.Error() != "custom msg" {
		t.Errorf("Error() = %q, want 'custom msg'", e.Error())
	}
}

func TestRenderUnknownTargetError_EmptyValue(t *testing.T) {
	out := RenderUnknownTargetError("", []string{"claude"})
	if !strings.Contains(out, "[x]") {
		t.Errorf("expected [x] prefix for empty value, got %q", out)
	}
}
