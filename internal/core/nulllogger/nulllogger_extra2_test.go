package nulllogger_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/nulllogger"
)

func TestNullCommandLogger_ProgressNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// should not panic
	l.Progress("progress message", "")
}

func TestNullCommandLogger_SuccessNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Success("success message", "[+]")
}

func TestNullCommandLogger_WarningNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Warning("warning message", "[!]")
}

func TestNullCommandLogger_ErrorNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.Error("error message", "[x]")
}

func TestNullCommandLogger_TreeItem_NonEmpty(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.TreeItem("tree item value")
}

func TestNullCommandLogger_PackageInlineWarning_Silent(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.PackageInlineWarning("suppressed")
}

func TestNullCommandLogger_MCPHeartbeat_ExactlyOne(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(1) // should say "server" (singular)
}

func TestNullCommandLogger_MCPHeartbeat_NegativeNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	l.MCPLookupHeartbeat(-1)
}

func TestNullCommandLogger_MultipleCallsNoOp(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	for i := 0; i < 5; i++ {
		l.Start("step", "")
		l.Progress("step", "")
		l.Success("step", "")
		l.Warning("step", "")
		l.Error("step", "")
	}
}

func TestNullCommandLogger_VerboseFlag_DoesNotAffectOtherMethods(t *testing.T) {
	l := &nulllogger.NullCommandLogger{Verbose: true}
	l.Start("start", "[*]")
	l.Progress("prog", "[>]")
	l.Success("ok", "[+]")
}

func TestNullCommandLogger_MCPHeartbeat_Zero_Stable(t *testing.T) {
	l := &nulllogger.NullCommandLogger{}
	// Should not panic or produce output for 0
	for i := 0; i < 3; i++ {
		l.MCPLookupHeartbeat(0)
	}
}
