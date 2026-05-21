package errors_test

import (
"strings"
"testing"

apmerrors "github.com/githubnext/apm/internal/core/errors"
)

func TestRenderNoHarnessError_ContainsApmInstall(t *testing.T) {
msg := apmerrors.RenderNoHarnessError()
if !strings.Contains(msg, "apm install") {
t.Error("expected 'apm install' in no-harness error")
}
}

func TestRenderAmbiguousError_SingleElement(t *testing.T) {
msg := apmerrors.RenderAmbiguousError([]string{"claude"})
if !strings.Contains(msg, "claude") {
t.Error("expected 'claude' in ambiguous error for single element")
}
}

func TestRenderAmbiguousError_ThreeElements(t *testing.T) {
msg := apmerrors.RenderAmbiguousError([]string{"claude", "copilot", "gemini"})
if !strings.Contains(msg, "gemini") {
t.Error("expected 'gemini' in ambiguous error")
}
}

func TestRenderUnknownTargetError_StripsQuotes(t *testing.T) {
msg := apmerrors.RenderUnknownTargetError("'claude'", []string{"copilot"})
if !strings.Contains(msg, "claude") {
t.Error("expected stripped target name in error")
}
}

func TestRenderUnknownTargetError_ValidListShown(t *testing.T) {
msg := apmerrors.RenderUnknownTargetError("badval", []string{"claude", "copilot"})
if !strings.Contains(msg, "claude") {
t.Error("expected valid target list in error")
}
}

func TestRenderConflictingSchemaError_ContainsTargets(t *testing.T) {
msg := apmerrors.RenderConflictingSchemaError()
if !strings.Contains(msg, "targets:") {
t.Error("expected 'targets:' in conflicting schema error")
}
}

func TestNoHarnessError_ImplementsError(t *testing.T) {
var e apmerrors.NoHarnessError
e.Message = "test"
if e.Error() != "test" {
t.Errorf("expected 'test', got %q", e.Error())
}
}

func TestAmbiguousHarnessError_ImplementsError(t *testing.T) {
var e apmerrors.AmbiguousHarnessError
e.Message = "ambiguous"
if e.Error() != "ambiguous" {
t.Errorf("expected 'ambiguous', got %q", e.Error())
}
}
