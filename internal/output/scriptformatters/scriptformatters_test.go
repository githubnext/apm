package scriptformatters

import (
"strings"
"testing"
)

func TestNewScriptExecutionFormatter(t *testing.T) {
f := NewScriptExecutionFormatter()
if f == nil {
t.Fatal("expected non-nil formatter")
}
}

func TestFormatScriptHeader(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatScriptHeader("my-script", map[string]string{"env": "prod"})
combined := strings.Join(lines, "\n")
if !strings.Contains(combined, "my-script") {
t.Errorf("header should contain script name, got: %s", combined)
}
}

func TestFormatScriptHeaderEmpty(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatScriptHeader("test", nil)
if len(lines) == 0 {
t.Error("expected at least one header line")
}
}

func TestFormatCompilationProgress(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatCompilationProgress([]string{"prompt1.md", "prompt2.md"})
combined := strings.Join(lines, "\n")
_ = combined // may or may not contain filenames
if len(lines) == 0 {
t.Error("expected at least one progress line")
}
}

func TestFormatRuntimeExecution(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatRuntimeExecution("python", "run.py", 512)
combined := strings.Join(lines, "\n")
if !strings.Contains(combined, "python") {
t.Errorf("expected runtime name in output: %s", combined)
}
}

func TestFormatContentPreview(t *testing.T) {
f := NewScriptExecutionFormatter()
content := "line1\nline2\nline3\nline4\nline5"
lines := f.FormatContentPreview(content, 3)
_ = lines // validates no panic
}

func TestFormatEnvironmentSetup(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatEnvironmentSetup("node", []string{"API_KEY", "DEBUG"})
combined := strings.Join(lines, "\n")
_ = combined
if len(lines) == 0 {
t.Error("expected at least one env setup line")
}
}

func TestFormatExecutionSuccess(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatExecutionSuccess("go", 1.23)
if len(lines) == 0 {
t.Error("expected at least one success line")
}
}

func TestFormatExecutionError(t *testing.T) {
f := NewScriptExecutionFormatter()
lines := f.FormatExecutionError("python", 1, "module not found")
combined := strings.Join(lines, "\n")
if !strings.Contains(combined, "python") && !strings.Contains(combined, "not found") {
t.Errorf("expected error info in output: %s", combined)
}
}

func TestFormatAutoDiscoveryMessage(t *testing.T) {
	f := NewScriptExecutionFormatter()
	msg := f.FormatAutoDiscoveryMessage("my-script", "prompt.md", "go")
	if msg == "" {
		t.Error("expected non-empty auto discovery message")
	}
}

func TestFormatCompilationProgressSingle(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatCompilationProgress([]string{"only.md"})
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "Compiling prompt") {
		t.Errorf("single prompt: expected 'Compiling prompt', got: %s", combined)
	}
}

func TestFormatCompilationProgressNone(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatCompilationProgress(nil)
	if lines != nil {
		t.Errorf("expected nil for empty prompt list, got: %v", lines)
	}
}

func TestFormatCompilationProgressLastLineReplaced(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatCompilationProgress([]string{"a.md", "b.md", "c.md"})
	if len(lines) == 0 {
		t.Fatal("expected non-empty lines")
	}
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, "+-") {
		t.Errorf("last line should start with '+-', got: %q", last)
	}
}

func TestFormatEnvironmentSetupEmpty(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatEnvironmentSetup("node", nil)
	if lines != nil {
		t.Errorf("expected nil for empty env vars, got: %v", lines)
	}
}

func TestFormatEnvironmentSetupLastLinePlusMinus(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatEnvironmentSetup("go", []string{"TOKEN", "SECRET"})
	if len(lines) == 0 {
		t.Fatal("expected lines")
	}
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, "+-") {
		t.Errorf("last var line should start with '+-', got: %q", last)
	}
}

func TestFormatExecutionSuccessNoTime(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionSuccess("node", -1)
	if len(lines) == 0 {
		t.Error("expected at least one line")
	}
	combined := strings.Join(lines, "\n")
	if strings.Contains(combined, "s)") {
		t.Errorf("should not show time when executionTime < 0, got: %s", combined)
	}
}

func TestFormatExecutionErrorMultilineMsg(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionError("ruby", 2, "line1\nline2\nline3")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines for multiline error, got %d", len(lines))
	}
}

func TestFormatExecutionErrorEmptyMsg(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatExecutionError("go", 1, "")
	if len(lines) != 1 {
		t.Errorf("expected exactly 1 line for empty error msg, got %d", len(lines))
	}
}

func TestFormatContentPreviewTruncates(t *testing.T) {
	f := NewScriptExecutionFormatter()
	long := strings.Repeat("x", 300)
	lines := f.FormatContentPreview(long, 100)
	for _, l := range lines {
		if strings.Contains(l, "...") {
			return
		}
	}
	t.Error("expected truncation ellipsis in preview")
}

func TestFormatContentPreviewDefaultMaxPreview(t *testing.T) {
	f := NewScriptExecutionFormatter()
	short := "short content"
	lines := f.FormatContentPreview(short, 0)
	found := false
	for _, l := range lines {
		if l == short {
			found = true
		}
	}
	if !found {
		t.Errorf("expected full content in preview lines: %v", lines)
	}
}

func TestFormatSubprocessDetails(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatSubprocessDetails([]string{"python", "-c", "print('hi')"}, 42)
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "python") {
		t.Errorf("expected python in subprocess details, got: %s", combined)
	}
	if !strings.Contains(combined, "42") {
		t.Errorf("expected content length in subprocess details, got: %s", combined)
	}
}

func TestFormatSubprocessDetailsSpacedArg(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatSubprocessDetails([]string{"my script", "arg"}, 0)
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, `"my script"`) {
		t.Errorf("expected quoted spaced arg in output, got: %s", combined)
	}
}

func TestFormatRuntimeExecutionContentLength(t *testing.T) {
	f := NewScriptExecutionFormatter()
	lines := f.FormatRuntimeExecution("go", "main", 1024)
	combined := strings.Join(lines, "\n")
	if !strings.Contains(combined, "1024") {
		t.Errorf("expected content length in output, got: %s", combined)
	}
}

func TestFormatScriptHeaderMultipleParams(t *testing.T) {
	f := NewScriptExecutionFormatter()
	params := map[string]string{"a": "1", "b": "2", "c": "3"}
	lines := f.FormatScriptHeader("batch", params)
	// header line + 3 param lines
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (1 header + 3 params), got %d", len(lines))
	}
}
