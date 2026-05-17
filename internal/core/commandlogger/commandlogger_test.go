package commandlogger_test

import (
"testing"

"github.com/githubnext/apm/internal/core/commandlogger"
)

func TestStripSourcePrefix(t *testing.T) {
cases := []struct {
input, want string
}{
{"org:myorg", "myorg"},
{"url:https://example.com", "https://example.com"},
{"bare", "bare"},
{"", ""},
{"org:", "org:"},
}
for _, c := range cases {
got := commandlogger.StripSourcePrefix(c.input)
if got != c.want {
t.Errorf("StripSourcePrefix(%q) = %q, want %q", c.input, got, c.want)
}
}
}

func TestNewCommandLogger(t *testing.T) {
l := commandlogger.NewCommandLogger("install", true, false)
if l == nil {
t.Fatal("NewCommandLogger returned nil")
}
if l.Command != "install" {
t.Errorf("Command = %q, want install", l.Command)
}
if !l.Verbose {
t.Error("expected Verbose=true")
}
if l.DryRun {
t.Error("expected DryRun=false")
}
}

func TestCommandLogger_ShouldExecute(t *testing.T) {
live := commandlogger.NewCommandLogger("test", false, false)
if !live.ShouldExecute() {
t.Error("expected ShouldExecute=true for non-dry-run")
}
dry := commandlogger.NewCommandLogger("test", false, true)
if dry.ShouldExecute() {
t.Error("expected ShouldExecute=false for dry-run")
}
}

func TestNewInstallLogger(t *testing.T) {
l := commandlogger.NewInstallLogger(false, false, false)
if l == nil {
t.Fatal("NewInstallLogger returned nil")
}
}

func TestCommandLogger_DryRunNotice(t *testing.T) {
l := commandlogger.NewCommandLogger("install", false, true)
// DryRunNotice should not panic.
l.DryRunNotice("would install 3 packages")
}

func TestCommandLogger_MCPLookupHeartbeat_zero(t *testing.T) {
l := commandlogger.NewCommandLogger("install", false, false)
// count=0 should be a no-op (no panic).
l.MCPLookupHeartbeat(0)
}

func TestCommandLogger_MCPLookupHeartbeat_one(t *testing.T) {
l := commandlogger.NewCommandLogger("install", false, false)
l.MCPLookupHeartbeat(1)
}

func TestCommandLogger_MCPLookupHeartbeat_many(t *testing.T) {
l := commandlogger.NewCommandLogger("install", false, false)
l.MCPLookupHeartbeat(5)
}

func TestCommandLogger_BlankLine(t *testing.T) {
l := commandlogger.NewCommandLogger("test", false, false)
l.BlankLine()
}

func TestCommandLogger_TreeItem(t *testing.T) {
l := commandlogger.NewCommandLogger("test", false, false)
l.TreeItem("some-package@1.0.0")
}

func TestCommandLogger_VerboseDetail_notVerbose(t *testing.T) {
l := commandlogger.NewCommandLogger("test", false, false)
l.VerboseDetail("hidden detail")
}

func TestCommandLogger_VerboseDetail_verbose(t *testing.T) {
l := commandlogger.NewCommandLogger("test", true, false)
l.VerboseDetail("visible detail")
}

func TestStripSourcePrefix_urlWithColon(t *testing.T) {
got := commandlogger.StripSourcePrefix("url:https://host:8080/path")
want := "https://host:8080/path"
if got != want {
t.Errorf("StripSourcePrefix(%q) = %q, want %q", "url:https://host:8080/path", got, want)
}
}
