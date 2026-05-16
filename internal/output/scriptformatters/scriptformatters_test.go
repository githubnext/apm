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
