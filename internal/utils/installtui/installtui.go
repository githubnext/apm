// Package installtui provides a shared Live-region TUI controller for the install pipeline.
//
// A single InstallTui instance is opened by apm install and is re-used
// across the resolve, download, integrate, and MCP-registry phases.
// Per-phase code calls StartPhase() once when the phase boundary is crossed,
// then TaskStarted() / TaskCompleted() / TaskFailed() for every dep / server /
// artifact in flight.
//
// When ShouldAnimate() is false (CI, dumb terminal, APM_PROGRESS=never,
// --quiet), every method on this struct is a cheap no-op. Callers do NOT
// need to gate their calls.
//
// This module uses a single ASCII spinner (| / - \) and never emits emoji
// or Unicode box-drawing, to stay safe under Windows cp1252.
package installtui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// DeferShowDuration is how long after Open() before the spinner is shown.
// Installs that finish under this threshold never paint a spinner.
const DeferShowDuration = 250 * time.Millisecond

// refreshInterval is the spinner update interval (8 Hz).
const refreshInterval = 125 * time.Millisecond

var spinnerFrames = []string{"|", "/", "-", "\\"}

// InstallTui is the TUI controller.
type InstallTui struct {
	out     io.Writer
	animate bool
	quiet   bool

	mu            sync.Mutex
	phase         string
	activeTasks   map[string]bool
	failedTasks   map[string]string
	completedCount int
	failedCount    int

	// spinner state
	stopCh    chan struct{}
	stoppedCh chan struct{}
	started   bool
}

// New creates a new InstallTui. quiet disables animation regardless of TTY.
func New(out io.Writer, quiet bool) *InstallTui {
	if out == nil {
		out = os.Stdout
	}
	animate := ShouldAnimate() && !quiet
	return &InstallTui{
		out:         out,
		animate:     animate,
		quiet:       quiet,
		activeTasks: make(map[string]bool),
		failedTasks: make(map[string]string),
	}
}

// ShouldAnimate returns true if the TUI should animate.
// Respects NO_COLOR, TERM=dumb, APM_PROGRESS env, and TTY detection.
func ShouldAnimate() bool {
	prog := os.Getenv("APM_PROGRESS")
	if prog == "never" || prog == "0" || prog == "false" {
		return false
	}
	if prog == "always" || prog == "1" || prog == "true" {
		return true
	}
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	// Check if stdout is a TTY
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Open begins the TUI session (deferred by DeferShowDuration).
func (t *InstallTui) Open() {
	if !t.animate {
		return
	}
	t.mu.Lock()
	if t.started {
		t.mu.Unlock()
		return
	}
	t.started = true
	t.stopCh = make(chan struct{})
	t.stoppedCh = make(chan struct{})
	t.mu.Unlock()

	go t.spinLoop()
}

// Close tears down the TUI session.
func (t *InstallTui) Close() {
	t.mu.Lock()
	if !t.started || t.stopCh == nil {
		t.mu.Unlock()
		return
	}
	stopCh := t.stopCh
	t.mu.Unlock()

	close(stopCh)
	<-t.stoppedCh
	// Clear the spinner line
	fmt.Fprint(t.out, "\r\033[K")
}

// StartPhase signals a new install phase.
func (t *InstallTui) StartPhase(phase string) {
	if !t.animate {
		return
	}
	t.mu.Lock()
	t.phase = phase
	t.activeTasks = make(map[string]bool)
	t.completedCount = 0
	t.failedCount = 0
	t.mu.Unlock()
}

// TaskStarted records that a task has started.
func (t *InstallTui) TaskStarted(name string) {
	if !t.animate {
		return
	}
	t.mu.Lock()
	t.activeTasks[name] = true
	t.mu.Unlock()
}

// TaskCompleted records that a task completed successfully.
func (t *InstallTui) TaskCompleted(name string) {
	if !t.animate {
		return
	}
	t.mu.Lock()
	delete(t.activeTasks, name)
	t.completedCount++
	t.mu.Unlock()
}

// TaskFailed records that a task failed with the given reason.
func (t *InstallTui) TaskFailed(name, reason string) {
	if !t.animate {
		return
	}
	t.mu.Lock()
	delete(t.activeTasks, name)
	t.failedTasks[name] = reason
	t.failedCount++
	t.mu.Unlock()
}

func (t *InstallTui) spinLoop() {
	defer close(t.stoppedCh)
	// Defer showing the spinner
	select {
	case <-t.stopCh:
		return
	case <-time.After(DeferShowDuration):
	}
	frame := 0
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			t.mu.Lock()
			phase := t.phase
			active := len(t.activeTasks)
			completed := t.completedCount
			failed := t.failedCount
			// Pick first active task name
			var firstName string
			for k := range t.activeTasks {
				firstName = k
				break
			}
			t.mu.Unlock()

			spinner := spinnerFrames[frame%len(spinnerFrames)]
			frame++
			line := buildSpinnerLine(spinner, phase, active, completed, failed, firstName)
			fmt.Fprintf(t.out, "\r\033[K%s", line)
		}
	}
}

func buildSpinnerLine(spinner, phase string, active, completed, failed int, firstName string) string {
	var sb strings.Builder
	sb.WriteString(spinner)
	sb.WriteString(" ")
	if phase != "" {
		sb.WriteString("[")
		sb.WriteString(phase)
		sb.WriteString("] ")
	}
	if firstName != "" {
		name := firstName
		if len(name) > 40 {
			name = name[:37] + "..."
		}
		sb.WriteString(name)
		sb.WriteString(" ")
	}
	if active > 0 || completed > 0 || failed > 0 {
		sb.WriteString(fmt.Sprintf("(%d active, %d done", active, completed))
		if failed > 0 {
			sb.WriteString(fmt.Sprintf(", %d failed", failed))
		}
		sb.WriteString(")")
	}
	return sb.String()
}

// Enter implements a context-manager-style entry. Returns the tui itself.
func (t *InstallTui) Enter() *InstallTui {
	t.Open()
	return t
}

// Exit implements a context-manager-style exit.
func (t *InstallTui) Exit(origErr error) {
	t.Close()
}
