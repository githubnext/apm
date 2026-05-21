package commandlogger

import (
	"testing"
)

func TestNewCommandLogger_VerboseFalse(t *testing.T) {
	l := NewCommandLogger("check", false, false)
	if l.Verbose {
		t.Error("expected Verbose=false")
	}
}

func TestNewCommandLogger_DryRunTrue(t *testing.T) {
	l := NewCommandLogger("install", false, true)
	if !l.DryRun {
		t.Error("expected DryRun=true")
	}
}

func TestNewCommandLogger_CommandField(t *testing.T) {
	l := NewCommandLogger("uninstall", false, false)
	if l.Command != "uninstall" {
		t.Errorf("expected Command=uninstall, got %q", l.Command)
	}
}

func TestNewCommandLogger_NotNil(t *testing.T) {
	l := NewCommandLogger("run", true, true)
	if l == nil {
		t.Error("expected non-nil logger")
	}
}

func TestStripSourcePrefix_OrgColon(t *testing.T) {
	result := StripSourcePrefix("org:acme")
	if result != "acme" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestStripSourcePrefix_URLColon(t *testing.T) {
	result := StripSourcePrefix("url:http://example.com")
	if result != "http://example.com" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestCommandLogger_MethodsNoopWhenNoWriter(t *testing.T) {
	l := NewCommandLogger("test", false, false)
	l.BlankLine()
	l.VerboseDetail("detail")
	l.TreeItem("item")
	l.MCPLookupHeartbeat(3)
}

func TestCommandLogger_PackageInlineWarning(t *testing.T) {
	l := NewCommandLogger("test", false, false)
	l.PackageInlineWarning("msg")
}

func TestCommandLogger_DryRunNotice(t *testing.T) {
	l := NewCommandLogger("test", false, true)
	l.DryRunNotice("would install foo")
}
