package installtui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_DefaultWriter(t *testing.T) {
	tui := New(nil, true)
	if tui == nil {
		t.Fatal("expected non-nil InstallTui")
	}
}

func TestNew_CustomWriterBuf(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	if tui == nil {
		t.Fatal("expected non-nil InstallTui")
	}
}

func TestTaskStarted_NoPanic(t *testing.T) {
	tui := New(nil, true)
	tui.TaskStarted("task-a")
	tui.TaskStarted("task-b")
}

func TestTaskCompleted_NoPanic(t *testing.T) {
	tui := New(nil, true)
	tui.TaskStarted("task-a")
	tui.TaskCompleted("task-a")
}

func TestTaskFailed_NoPanicVariant(t *testing.T) {
	tui := New(nil, true)
	tui.TaskStarted("task-x")
	tui.TaskFailed("task-x", "timeout")
}

func TestTaskCompleted_UnknownTask(t *testing.T) {
	tui := New(nil, true)
	// Completing a task never started should not panic
	tui.TaskCompleted("ghost-task")
}

func TestStartPhase_MultiplePhasesNoPanic(t *testing.T) {
	tui := New(nil, true)
	tui.StartPhase("resolve")
	tui.StartPhase("download")
	tui.StartPhase("integrate")
}

func TestBuildSpinnerLine_NoActive(t *testing.T) {
	line := buildSpinnerLine("|", "downloading", 0, 5, 0, "")
	if line == "" {
		t.Error("expected non-empty spinner line")
	}
}

func TestBuildSpinnerLine_WithActive(t *testing.T) {
	line := buildSpinnerLine("/", "installing", 3, 2, 1, "mypkg")
	if !strings.Contains(line, "mypkg") {
		t.Errorf("spinner line should contain first active task name, got %q", line)
	}
}

func TestBuildSpinnerLine_PhaseIncluded(t *testing.T) {
	line := buildSpinnerLine("-", "resolve", 1, 0, 0, "dep")
	if !strings.Contains(line, "resolve") {
		t.Errorf("expected phase in spinner line, got %q", line)
	}
}

func TestOpenClose_Quiet(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.Close()
}

func TestEnterExit_NilError(t *testing.T) {
	tui := New(nil, true)
	got := tui.Enter()
	if got != tui {
		t.Error("Enter should return self")
	}
	tui.Exit(nil)
}
