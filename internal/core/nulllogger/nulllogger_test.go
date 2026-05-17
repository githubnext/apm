package nulllogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/nulllogger"
)

// NullCommandLogger methods must not panic.
func TestNullCommandLoggerNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}

	l.Start("msg", "")
	l.Start("msg", "running")
	l.Progress("msg", "")
	l.Success("msg", "")
	l.Warning("msg", "")
	l.Error("msg", "")
	l.VerboseDetail("msg")
	l.TreeItem("item")
	l.PackageInlineWarning("warn")
	l.MCPLookupHeartbeat(0)
	l.MCPLookupHeartbeat(1)
	l.MCPLookupHeartbeat(5)
}

func TestNullCommandLoggerVerboseDefault(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	if l.Verbose {
		t.Error("Verbose should default to false")
	}
}

func TestNullCommandLoggerMCPHeartbeatSingular(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(1)
}

func TestNullCommandLoggerMCPHeartbeatZero(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// Zero count should produce no output without panic
	l.MCPLookupHeartbeat(0)
}

func TestNullCommandLoggerMCPHeartbeatPlural(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// Plural form (count > 1)
	for _, n := range []int{2, 3, 5, 10, 100} {
		l.MCPLookupHeartbeat(n)
	}
}

func TestNullCommandLoggerStartWithSymbol(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// Start with explicit symbol should not panic
	l.Start("installing", "arrow")
	l.Start("done", "[+]")
}

func TestNullCommandLoggerAllMethods(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	msgs := []string{"", "hello", "multi word message", "with/slash"}
	for _, m := range msgs {
		l.Start(m, "")
		l.Progress(m, "")
		l.Success(m, "")
		l.Warning(m, "")
		l.Error(m, "")
		l.VerboseDetail(m)
		l.TreeItem(m)
		l.PackageInlineWarning(m)
	}
}

func TestNullCommandLoggerVerboseSet(t *testing.T) {
	l := &nulllogger.NullCommandLogger{Verbose: true}
	if !l.Verbose {
		t.Error("Verbose should be true when set")
	}
	// VerboseDetail should still not panic when Verbose=true
	l.VerboseDetail("verbose message")
}

func TestNullCommandLoggerLargeCount(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// Large counts should not panic
	l.MCPLookupHeartbeat(999)
	l.MCPLookupHeartbeat(1000)
}
