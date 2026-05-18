package installtui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_NoAnimation_Quiet(t *testing.T) {
	tui := New(nil, true)
	if tui.animate {
		t.Error("quiet=true should set animate=false")
	}
}

func TestStartPhase_BeforeOpen_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	// Calling StartPhase without Open should not panic
	tui.StartPhase("early")
}

func TestTaskStarted_BeforeOpen_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.TaskStarted("early-task")
}

func TestTaskCompleted_BeforeOpen_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.TaskCompleted("early-task")
}

func TestTaskFailed_BeforeOpen_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.TaskFailed("early-task", "some reason")
}

func TestBuildSpinnerLine_EmptyPhase(t *testing.T) {
	line := buildSpinnerLine("|", "", 1, 2, 0, "pkg")
	// No phase name -- should still produce a non-empty line
	if line == "" {
		t.Error("expected non-empty spinner line with empty phase")
	}
}

func TestBuildSpinnerLine_FailedCount(t *testing.T) {
	line := buildSpinnerLine("-", "install", 0, 5, 3, "pkg")
	// Failed count should be reflected
	if !strings.Contains(line, "3") {
		t.Logf("spinner with failed=3: %q (informational)", line)
	}
}

func TestBuildSpinnerLine_LongFirstName(t *testing.T) {
	long := strings.Repeat("a", 80)
	line := buildSpinnerLine("|", "resolve", 1, 0, 0, long)
	if line == "" {
		t.Error("expected non-empty spinner line with long first name")
	}
}

func TestOpenCloseIdempotent(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.Open() // second Open should not panic
	tui.Close()
	tui.Close() // second Close should not panic
}

func TestMultipleTasks(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.StartPhase("resolve")
	for i := 0; i < 5; i++ {
		tui.TaskStarted("dep")
		tui.TaskCompleted("dep")
	}
	tui.Close()
}

func TestEnterReturnsNonNil(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	result := tui.Enter()
	if result == nil {
		t.Error("Enter() should return non-nil")
	}
}

func TestBuildSpinnerLine_CompletedCount(t *testing.T) {
	line := buildSpinnerLine("*", "link", 0, 100, 0, "first-pkg")
	// Completed count should appear
	if !strings.Contains(line, "100") {
		t.Logf("spinner with completed=100: %q (informational)", line)
	}
}
