package installtui_test

import (
	"bytes"
	"testing"

	"github.com/githubnext/apm/internal/utils/installtui"
)

func TestNew_QuietMode(t *testing.T) {
	tui := installtui.New(nil, true)
	if tui == nil {
		t.Fatal("expected non-nil TUI")
	}
}

func TestOpen_Close_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := installtui.New(&buf, true)
	tui.Open()
	tui.Close()
}

func TestStartPhase_AfterOpen_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := installtui.New(&buf, true)
	tui.Open()
	tui.StartPhase("resolve")
	tui.Close()
}

func TestTaskFailed_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := installtui.New(&buf, true)
	tui.Open()
	tui.TaskStarted("pkg1")
	tui.TaskFailed("pkg1", "network error")
	tui.Close()
}

func TestEnterExit_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	tui := installtui.New(&buf, true)
	tui.Open()
	tui2 := tui.Enter()
	if tui2 == nil {
		t.Fatal("Enter returned nil")
	}
	tui.Exit(nil)
	tui.Close()
}

func TestMultipleTasksCompleted(t *testing.T) {
	var buf bytes.Buffer
	tui := installtui.New(&buf, true)
	tui.Open()
	tui.StartPhase("download")
	for i := 0; i < 5; i++ {
		tui.TaskStarted("task")
		tui.TaskCompleted("task")
	}
	tui.Close()
}
