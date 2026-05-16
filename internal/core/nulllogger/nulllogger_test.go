package nulllogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/nulllogger"
)

// NullCommandLogger methods must not panic.
func TestNullCommandLoggerNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}

	// These should not panic
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
	// Should not panic for single server
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(1)
}
