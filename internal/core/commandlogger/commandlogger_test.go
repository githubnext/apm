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
