package scriptformatters_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/output/scriptformatters"
)

func TestFormatScriptHeader_SingleParam(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatScriptHeader("run.sh", map[string]string{"env": "prod"})
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "run.sh") {
		t.Errorf("first line should mention script name: %q", lines[0])
	}
}

func TestFormatScriptHeader_NoParams(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatScriptHeader("test.sh", nil)
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}

func TestFormatCompilationProgress_SingleFile(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatCompilationProgress([]string{"a.md"})
	if len(lines) < 1 {
		t.Fatal("expected lines")
	}
}

func TestFormatCompilationProgress_NilInput(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatCompilationProgress(nil)
	if lines != nil {
		t.Error("expected nil for empty input")
	}
}

func TestFormatRuntimeExecution_ContainsRuntime(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatRuntimeExecution("copilot", "cmd", 100)
	found := false
	for _, l := range lines {
		if strings.Contains(l, "copilot") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'copilot' in output")
	}
}

func TestFormatRuntimeExecution_ThreeLines(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatRuntimeExecution("llm", "run", 50)
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestFormatContentPreview_FourLines(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatContentPreview("hello world", 200)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestFormatContentPreview_ZeroMaxUsesDefault(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatContentPreview("short", 0)
	if len(lines) == 0 {
		t.Error("expected output")
	}
}

func TestFormatExecutionSuccess_NegativeTime(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatExecutionSuccess("copilot", -1)
	if len(lines) == 0 {
		t.Error("expected output")
	}
}

func TestFormatExecutionError_ZeroCode(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatExecutionError("copilot", 0, "")
	if len(lines) == 0 {
		t.Error("expected at least one line")
	}
}

func TestFormatSubprocessDetails_SingleArg(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	lines := f.FormatSubprocessDetails([]string{"cmd"}, 10)
	if len(lines) == 0 {
		t.Error("expected output")
	}
}

func TestFormatAutoDiscoveryMessage_NonEmpty(t *testing.T) {
	f := scriptformatters.NewScriptExecutionFormatter()
	msg := f.FormatAutoDiscoveryMessage("script.sh", "prompt.md", "copilot")
	if msg == "" {
		t.Error("expected non-empty message")
	}
}
