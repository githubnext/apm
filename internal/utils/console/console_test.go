package console_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/console"
)

func TestStatusSymbols(t *testing.T) {
	cases := map[string]string{
		"success": "[*]",
		"error":   "[x]",
		"warning": "[!]",
		"info":    "[i]",
		"check":   "[+]",
	}
	for k, want := range cases {
		if got := console.StatusSymbols[k]; got != want {
			t.Errorf("StatusSymbols[%q] = %q, want %q", k, got, want)
		}
	}
}

func TestEcho_noColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	console.Echo(&buf, "hello", "green", "", false)
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected 'hello' in output, got %q", buf.String())
	}
}

func TestEcho_withSymbol(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	console.Echo(&buf, "done", "", "check", false)
	if !strings.Contains(buf.String(), "[+]") {
		t.Errorf("expected symbol [+] in output, got %q", buf.String())
	}
}

func TestPrintFilesTable_smoke(t *testing.T) {
	// Just ensure no panic.
	console.PrintFilesTable([][]string{{"file.go", "main source"}}, "Files")
}
