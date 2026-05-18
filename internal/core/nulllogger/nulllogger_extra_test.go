package nulllogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/nulllogger"
)

func TestNullCommandLoggerVerboseDetailNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// VerboseDetail when Verbose=false should not panic
	l.Verbose = false
	l.VerboseDetail("some detail message")
}

func TestNullCommandLoggerVerboseDetailVerbose(t *testing.T) {
	l := &nulllogger.NullCommandLogger{Verbose: true}
	l.VerboseDetail("detail with verbose=true")
}

func TestNullCommandLoggerStartEmptySymbol(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Start("message", "")
}

func TestNullCommandLoggerStartNonEmptySymbol(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Start("message", "[*]")
}

func TestNullCommandLoggerProgressVariants(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Progress("downloading", "")
	l.Progress("", "")
	l.Progress("long message with many words here", "sym")
}

func TestNullCommandLoggerSuccessVariants(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Success("done", "")
	l.Success("installed", "[+]")
}

func TestNullCommandLoggerWarningVariants(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Warning("warn msg", "")
	l.Warning("!", "[!]")
}

func TestNullCommandLoggerErrorVariants(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Error("error msg", "")
	l.Error("fail", "[x]")
}

func TestNullCommandLoggerTreeItemEmpty(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.TreeItem("")
}

func TestNullCommandLoggerTreeItemLong(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.TreeItem("very long tree item message with lots of text")
}

func TestNullCommandLoggerPackageInlineWarning(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.PackageInlineWarning("inline warning")
	l.PackageInlineWarning("")
}

func TestNullCommandLoggerMCPHeartbeatRange(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	for _, n := range []int{0, 1, 2, 5, 10, 50, 100, 1000} {
		l.MCPLookupHeartbeat(n)
	}
}

func TestNullCommandLoggerZeroValue(t *testing.T) {
	// Zero-value struct should work for all methods
	var l nulllogger.NullCommandLogger
	l.Start("s", "")
	l.Progress("p", "")
	l.Success("ok", "")
	l.Warning("w", "")
	l.Error("e", "")
	l.VerboseDetail("v")
	l.TreeItem("t")
	l.PackageInlineWarning("pw")
	l.MCPLookupHeartbeat(3)
}
