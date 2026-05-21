package scriptrunner

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// formatCompilationProgress
// ---------------------------------------------------------------------------

func TestFormatCompilationProgress_Empty(t *testing.T) {
	lines := formatCompilationProgress([]string{})
	if len(lines) == 0 {
		t.Error("expected at least one line for empty input")
	}
}

func TestFormatCompilationProgress_SingleFile(t *testing.T) {
	lines := formatCompilationProgress([]string{"foo.prompt.md"})
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "foo.prompt.md") {
		t.Errorf("expected file name in output, got %q", joined)
	}
}

func TestFormatCompilationProgress_MultipleFiles(t *testing.T) {
	files := []string{"a.prompt.md", "b.prompt.md", "c.prompt.md"}
	lines := formatCompilationProgress(files)
	if len(lines) == 0 {
		t.Error("expected non-empty output for multiple files")
	}
}

// ---------------------------------------------------------------------------
// formatRuntimeExecution
// ---------------------------------------------------------------------------

func TestFormatRuntimeExecution_Copilot(t *testing.T) {
	lines := formatRuntimeExecution(RuntimeCopilot, "gh copilot suggest", 512)
	if len(lines) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestFormatRuntimeExecution_Unknown(t *testing.T) {
	lines := formatRuntimeExecution(RuntimeUnknown, "some-cmd", 0)
	if len(lines) == 0 {
		t.Error("expected non-empty output for unknown runtime")
	}
}

// ---------------------------------------------------------------------------
// formatContentPreview
// ---------------------------------------------------------------------------

func TestFormatContentPreview_ShortContent(t *testing.T) {
	lines := formatContentPreview("short content")
	if len(lines) == 0 {
		t.Error("expected non-empty preview")
	}
}

func TestFormatContentPreview_LongContent(t *testing.T) {
	content := strings.Repeat("x", 2000)
	lines := formatContentPreview(content)
	if len(lines) == 0 {
		t.Error("expected non-empty preview for long content")
	}
}

func TestFormatContentPreview_Empty(t *testing.T) {
	lines := formatContentPreview("")
	if lines == nil {
		t.Error("expected non-nil result for empty content")
	}
}

// ---------------------------------------------------------------------------
// formatEnvironmentSetup
// ---------------------------------------------------------------------------

func TestFormatEnvironmentSetup_NoVars(t *testing.T) {
	lines := formatEnvironmentSetup(RuntimeCopilot, []string{})
	if len(lines) == 0 {
		t.Error("expected at least one line")
	}
}

func TestFormatEnvironmentSetup_WithVars(t *testing.T) {
	lines := formatEnvironmentSetup(RuntimeCodex, []string{"FOO=bar", "BAZ=qux"})
	if len(lines) == 0 {
		t.Error("expected non-empty output")
	}
}

// ---------------------------------------------------------------------------
// formatExecutionSuccess / formatExecutionError
// ---------------------------------------------------------------------------

func TestFormatExecutionSuccess_Copilot(t *testing.T) {
	lines := formatExecutionSuccess(RuntimeCopilot)
	if len(lines) == 0 {
		t.Error("expected non-empty success lines")
	}
}

func TestFormatExecutionError_LLM(t *testing.T) {
	lines := formatExecutionError(RuntimeLLM)
	if len(lines) == 0 {
		t.Error("expected non-empty error lines")
	}
}

// ---------------------------------------------------------------------------
// isValidEnvVarName edge cases
// ---------------------------------------------------------------------------

func TestIsValidEnvVarName_Underscore(t *testing.T) {
	if !isValidEnvVarName("MY_VAR") {
		t.Error("expected MY_VAR to be valid")
	}
}

func TestIsValidEnvVarName_StartsWithDigit(t *testing.T) {
	if isValidEnvVarName("1BAD") {
		t.Error("expected 1BAD to be invalid")
	}
}

func TestIsValidEnvVarName_Empty(t *testing.T) {
	if isValidEnvVarName("") {
		t.Error("expected empty string to be invalid")
	}
}

func TestIsValidEnvVarName_Lowercase(t *testing.T) {
	// lowercase letters are allowed
	if !isValidEnvVarName("my_var") {
		t.Error("expected my_var to be valid")
	}
}

// ---------------------------------------------------------------------------
// envMapToSlice ordering stability
// ---------------------------------------------------------------------------

func TestEnvMapToSlice_Single(t *testing.T) {
	s := envMapToSlice(map[string]string{"K": "V"})
	if len(s) != 1 {
		t.Errorf("expected 1 entry, got %d", len(s))
	}
	if s[0] != "K=V" {
		t.Errorf("expected K=V, got %q", s[0])
	}
}

func TestEnvMapToSlice_Empty(t *testing.T) {
	s := envMapToSlice(map[string]string{})
	if len(s) != 0 {
		t.Errorf("expected 0 entries, got %d", len(s))
	}
}

// ---------------------------------------------------------------------------
// ScriptRunner fields
// ---------------------------------------------------------------------------

func TestScriptRunnerFields(t *testing.T) {
	sr := New(true)
	if sr.Compiler == nil {
		t.Error("Compiler should not be nil")
	}
	if !sr.UseColor {
		t.Error("UseColor should be true")
	}
}

func TestScriptRunnerFieldsFalse(t *testing.T) {
	sr := New(false)
	if sr.UseColor {
		t.Error("UseColor should be false")
	}
}
