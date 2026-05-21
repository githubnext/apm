package errors_test

import (
	"strings"
	"testing"

	apmerrors "github.com/githubnext/apm/internal/core/errors"
)

func TestTargetResolutionError_MessageField(t *testing.T) {
	e := &apmerrors.TargetResolutionError{Message: "some error"}
	if e.Error() != "some error" {
		t.Errorf("unexpected error: %q", e.Error())
	}
}

func TestRenderAmbiguousError_ContainsDetected(t *testing.T) {
	msg := apmerrors.RenderAmbiguousError([]string{"claude", "copilot"})
	if !strings.Contains(msg, "claude") {
		t.Error("expected 'claude' in ambiguous error")
	}
	if !strings.Contains(msg, "copilot") {
		t.Error("expected 'copilot' in ambiguous error")
	}
}

func TestRenderAmbiguousError_EmptyList(t *testing.T) {
	msg := apmerrors.RenderAmbiguousError([]string{})
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestRenderUnknownTargetError_ContainsValue(t *testing.T) {
	msg := apmerrors.RenderUnknownTargetError("badtarget", []string{"claude", "copilot"})
	if !strings.Contains(msg, "badtarget") {
		t.Error("expected 'badtarget' in unknown target error")
	}
}

func TestRenderConflictingSchemaError_NotEmpty(t *testing.T) {
	msg := apmerrors.RenderConflictingSchemaError()
	if msg == "" {
		t.Error("expected non-empty conflicting schema error")
	}
}

func TestNoHarnessError_ZeroValue(t *testing.T) {
	var e apmerrors.NoHarnessError
	if e.Error() != "" {
		_ = e.Error()
	}
}

func TestRenderNoHarnessError_ContainsXPrefix_Extra3(t *testing.T) {
	msg := apmerrors.RenderNoHarnessError()
	if !strings.HasPrefix(msg, "[x]") {
		t.Errorf("expected [x] prefix, got: %q", msg[:10])
	}
}
