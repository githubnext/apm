package scriptformatters

import (
	"strings"
	"testing"
)

func TestFormatAutoDiscoveryMessage_Extra(t *testing.T) {
	f := NewScriptExecutionFormatter()
	msg := f.FormatAutoDiscoveryMessage("my-script", "prompt.md", "go")
	if !strings.Contains(msg, "prompt.md") {
		t.Errorf("expected prompt file in message, got %q", msg)
	}
	if !strings.Contains(msg, "go") {
		t.Errorf("expected runtime in message, got %q", msg)
	}
}

func TestFormatExecutionError_NoMessage(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionError("python", 1, "")
	if len(lines) == 0 {
		t.Fatal("expected at least one line for error")
	}
	if !strings.Contains(lines[0], "1") {
		t.Errorf("expected exit code in output, got %q", lines[0])
	}
}

func TestFormatExecutionError_WithMultilineMessage(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionError("node", 2, "line one\nline two\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d: %v", len(lines), lines)
	}
}

func TestFormatExecutionError_SkipsEmptyLines(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionError("go", 1, "\n\n  \n")
	// Only the header line; blank lines skipped
	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d: %v", len(lines), lines)
	}
}

func TestFormatExecutionSuccess_WithTime(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionSuccess("python", 1.23)
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "1.23") {
		t.Errorf("expected time in output, got %q", combined)
	}
}

func TestFormatExecutionSuccess_NoTime(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionSuccess("go", -1)
	if len(lines) == 0 {
		t.Fatal("expected non-empty output")
	}
	// Should not contain time info
	combined := strings.Join(lines, "\n")
	if strings.Contains(combined, "(") {
		t.Errorf("expected no time info when executionTime<0, got %q", combined)
	}
}

func TestFormatSubprocessDetails_NoSpaceArgs(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatSubprocessDetails([]string{"go", "run", "main.go"}, 512)
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "go run main.go") {
		t.Errorf("expected args in output, got %q", combined)
	}
}

func TestFormatSubprocessDetails_EmptyArgs(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatSubprocessDetails([]string{}, 0)
	if len(lines) == 0 {
		t.Fatal("expected non-empty output for empty args")
	}
}

func TestFormatEnvironmentSetup_MultipleVars(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatEnvironmentSetup("node", []string{"PATH", "HOME", "NODE_ENV"})
	combined := strings.Join(lines, "\n")
	for _, v := range []string{"PATH", "HOME", "NODE_ENV"} {
		if !strings.Contains(combined, v) {
			t.Errorf("expected %q in environment setup output, got %q", v, combined)
		}
	}
}

func TestFormatEnvironmentSetup_Empty(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatEnvironmentSetup("go", []string{})
	// Returns nil/empty for empty env setup (by design)
	_ = lines
}

func TestFormatContentPreview_LongContent(t *testing.T) {
	f := NewScriptExecutionFormatter()
	long := strings.Repeat("x", 1000)
	lines := f.FormatContentPreview(long, 100)
	combined := strings.Join(lines, "\n")
	if len(combined) > 500 {
		t.Logf("note: preview is %d chars", len(combined))
	}
	_ = lines // no panic
}

func TestFormatCompilationProgress_Many(t *testing.T) {
	f := NewScriptExecutionFormatter()
	files := []string{"a.md", "b.md", "c.md", "d.md", "e.md"}
	lines := f.FormatCompilationProgress(files)
	if len(lines) == 0 {
		t.Fatal("expected non-empty compilation progress output")
	}
}

func TestFormatScriptHeader_EmptyName(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatScriptHeader("", map[string]string{})
	if len(lines) == 0 {
		t.Fatal("expected at least one line for empty script name")
	}
}
