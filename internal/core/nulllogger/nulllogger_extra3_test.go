package nulllogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/nulllogger"
)

func TestNullCommandLogger_ZeroValue(t *testing.T) {
	var l nulllogger.NullCommandLogger
	if l.Verbose {
		t.Error("expected Verbose=false for zero value")
	}
}

func TestNullCommandLogger_StartNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Start("starting", "")
}

func TestNullCommandLogger_StartWithSymbol(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Start("starting", "spin")
}

func TestNullCommandLogger_ErrorNoOp_Extra3(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Error("something failed", "")
}

func TestNullCommandLogger_MCPLookupHeartbeat_Zero(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(0)
}

func TestNullCommandLogger_MCPLookupHeartbeat_Positive(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(5)
}

func TestNullCommandLogger_MCPLookupHeartbeat_One(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(1)
}

func TestNullCommandLogger_TreeItem(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.TreeItem("leaf node")
}

func TestNullCommandLogger_VerboseDetailNoop(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.VerboseDetail("should be discarded")
}
