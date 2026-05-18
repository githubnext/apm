package installtui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_NotNilOut(t *testing.T) {
	tui := New(nil, true)
	if tui == nil {
		t.Fatal("expected non-nil InstallTui")
	}
}

func TestNew_QuietDisablesAnimation(t *testing.T) {
	tui := New(nil, true)
	if tui.animate {
		t.Error("quiet mode should disable animation")
	}
}

func TestNew_CustomWriter(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	if tui == nil {
		t.Fatal("expected non-nil with custom writer")
	}
}

func TestOpenClose_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.Close()
}

func TestStartPhase_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.StartPhase("resolve")
	tui.StartPhase("download")
	tui.Close()
}

func TestTaskStartedCompleted_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.StartPhase("resolve")
	tui.TaskStarted("dep1")
	tui.TaskCompleted("dep1")
	tui.Close()
}

func TestTaskFailed_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Open()
	tui.StartPhase("download")
	tui.TaskStarted("dep2")
	tui.TaskFailed("dep2", "network timeout")
	tui.Close()
}

func TestBuildSpinnerLine_Basic(t *testing.T) {
	line := buildSpinnerLine("|", "resolve", 3, 5, 1, "mypkg")
	if !strings.Contains(line, "resolve") {
		t.Errorf("expected phase name in spinner line, got %q", line)
	}
}

func TestBuildSpinnerLine_AllZero(t *testing.T) {
	line := buildSpinnerLine("/", "download", 0, 0, 0, "")
	if line == "" {
		t.Error("expected non-empty spinner line")
	}
}

func TestBuildSpinnerLine_NoUnicode(t *testing.T) {
	line := buildSpinnerLine("-", "finalize", 1, 10, 0, "pkg")
	for _, r := range line {
		if r > 127 {
			t.Errorf("spinner line contains non-ASCII character %q", string(r))
		}
	}
}

func TestEnterExit_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui2 := tui.Enter()
	if tui2 == nil {
		t.Fatal("Enter() should return self")
	}
	tui.Exit(nil)
}

func TestEnterExit_WithError(t *testing.T) {
	var buf bytes.Buffer
	tui := New(&buf, true)
	tui.Enter()
	// Exit with an error should not panic
	tui.Exit(nil)
}
