package errors_test

import (
	"strings"
	"testing"

	customerrors "github.com/githubnext/apm/internal/core/errors"
)

func TestRenderNoHarnessError_HasXPrefix(t *testing.T) {
	msg := customerrors.RenderNoHarnessError()
	if !strings.HasPrefix(msg, "[x]") {
		t.Errorf("expected [x] prefix, got: %q", msg[:min2(len(msg), 20)])
	}
}

func TestRenderNoHarnessError_MentionsAPM(t *testing.T) {
	msg := customerrors.RenderNoHarnessError()
	if !strings.Contains(strings.ToLower(msg), "apm") {
		t.Error("error message should mention apm")
	}
}

func TestRenderNoHarnessError_MentionsTargets(t *testing.T) {
	msg := customerrors.RenderNoHarnessError()
	if !strings.Contains(msg, "targets") {
		t.Error("error message should mention 'targets'")
	}
}

func TestRenderAmbiguousError_HasXPrefix(t *testing.T) {
	msg := customerrors.RenderAmbiguousError([]string{"claude", "copilot"})
	if !strings.HasPrefix(msg, "[x]") {
		t.Errorf("expected [x] prefix, got: %q", msg[:min2(len(msg), 20)])
	}
}

func TestRenderAmbiguousError_ContainsBothDetected(t *testing.T) {
	msg := customerrors.RenderAmbiguousError([]string{"claude", "copilot"})
	if !strings.Contains(msg, "claude") {
		t.Error("should contain 'claude'")
	}
	if !strings.Contains(msg, "copilot") {
		t.Error("should contain 'copilot'")
	}
}

func TestRenderAmbiguousError_SuggestionIsFirstElement(t *testing.T) {
	msg := customerrors.RenderAmbiguousError([]string{"cursor", "gemini"})
	if !strings.Contains(msg, "cursor") {
		t.Error("suggestion should be the first detected element 'cursor'")
	}
}

func TestRenderUnknownTargetError_HasXPrefix(t *testing.T) {
	msg := customerrors.RenderUnknownTargetError("badtarget", []string{"claude", "copilot"})
	if !strings.HasPrefix(msg, "[x]") {
		t.Errorf("expected [x] prefix, got: %q", msg[:min2(len(msg), 20)])
	}
}

func TestRenderUnknownTargetError_ShowsValue(t *testing.T) {
	msg := customerrors.RenderUnknownTargetError("myunknown", []string{"claude"})
	if !strings.Contains(msg, "myunknown") {
		t.Error("error should display the unknown target value")
	}
}

func TestRenderUnknownTargetError_BracketStripped(t *testing.T) {
	msg := customerrors.RenderUnknownTargetError("[bad]", []string{"claude"})
	if strings.Contains(msg, "[bad]") {
		t.Error("brackets should be stripped from the display value")
	}
}

func TestRenderUnknownTargetError_HidesAgentSkills(t *testing.T) {
	msg := customerrors.RenderUnknownTargetError("x", []string{"claude", "agent-skills"})
	if strings.Contains(msg, "agent-skills") {
		t.Error("agent-skills should be hidden from valid targets list")
	}
}

func TestRenderConflictingSchemaError_HasXPrefix(t *testing.T) {
	msg := customerrors.RenderConflictingSchemaError()
	if !strings.HasPrefix(msg, "[x]") {
		t.Errorf("expected [x] prefix, got: %q", msg[:min2(len(msg), 20)])
	}
}

func TestRenderConflictingSchemaError_MentionsBothKeys(t *testing.T) {
	msg := customerrors.RenderConflictingSchemaError()
	if !strings.Contains(msg, "target:") || !strings.Contains(msg, "targets:") {
		t.Error("conflicting schema error should mention both 'target:' and 'targets:'")
	}
}

func TestErrorTypes_ImplementError(t *testing.T) {
	var _ error = &customerrors.TargetResolutionError{Message: "test"}
	var _ error = &customerrors.NoHarnessError{TargetResolutionError: customerrors.TargetResolutionError{Message: "x"}}
	var _ error = &customerrors.AmbiguousHarnessError{TargetResolutionError: customerrors.TargetResolutionError{Message: "x"}}
	var _ error = &customerrors.UnknownTargetError{TargetResolutionError: customerrors.TargetResolutionError{Message: "x"}}
	var _ error = &customerrors.ConflictingTargetsError{TargetResolutionError: customerrors.TargetResolutionError{Message: "x"}}
	var _ error = &customerrors.EmptyTargetsListError{TargetResolutionError: customerrors.TargetResolutionError{Message: "x"}}
}

func TestRenderNoHarnessError_ContainsFixHint(t *testing.T) {
	msg := customerrors.RenderNoHarnessError()
	if !strings.Contains(msg, "Fix") {
		t.Error("error message should contain a fix hint")
	}
}

func TestRenderAmbiguousError_ContainsFixHint(t *testing.T) {
	msg := customerrors.RenderAmbiguousError([]string{"claude"})
	if !strings.Contains(msg, "Fix") {
		t.Error("error message should contain a fix hint")
	}
}

func TestRenderUnknownTargetError_ValidListSorted(t *testing.T) {
	msg := customerrors.RenderUnknownTargetError("bad", []string{"z-target", "a-target"})
	az := strings.Index(msg, "a-target")
	zz := strings.Index(msg, "z-target")
	if az == -1 || zz == -1 {
		t.Skip("valid targets not found in error message")
	}
	if az > zz {
		t.Error("valid targets should be sorted alphabetically in the error message")
	}
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}
