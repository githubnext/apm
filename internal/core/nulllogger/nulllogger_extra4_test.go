package nulllogger

import (
"testing"
)

func TestNullCommandLogger_ZeroValue(t *testing.T) {
var l NullCommandLogger
// methods should not panic
l.Start("starting", "")
l.Progress("progress", "")
l.Success("done", "")
l.Warning("warn", "")
l.Error("err", "")
l.VerboseDetail("verbose")
l.TreeItem("item")
l.PackageInlineWarning("warning")
l.MCPLookupHeartbeat(0)
}

func TestNullCommandLogger_MCPHeartbeat_Single(t *testing.T) {
var l NullCommandLogger
// count=1 uses "server" (singular), should not panic
l.MCPLookupHeartbeat(1)
}

func TestNullCommandLogger_MCPHeartbeat_Multiple(t *testing.T) {
var l NullCommandLogger
// count>1 uses "servers" (plural), should not panic
l.MCPLookupHeartbeat(5)
}

func TestNullCommandLogger_MCPHeartbeat_Zero(t *testing.T) {
var l NullCommandLogger
// count <= 0 is a no-op; should not panic
l.MCPLookupHeartbeat(0)
l.MCPLookupHeartbeat(-1)
}

func TestNullCommandLogger_Verbose_Field(t *testing.T) {
l := &NullCommandLogger{Verbose: true}
if !l.Verbose {
t.Error("expected Verbose=true")
}
}

func TestNullCommandLogger_StartNonEmptySymbol(t *testing.T) {
var l NullCommandLogger
// symbol non-empty path; should not panic
l.Start("message", "custom-symbol")
}
