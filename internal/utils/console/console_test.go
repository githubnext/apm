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

func TestPrintFilesTable_noTitle(t *testing.T) {
	// Empty title should not panic.
	console.PrintFilesTable([][]string{{"a.go", "pkg a"}, {"b.go", "pkg b"}}, "")
}

func TestPrintFilesTable_emptyRows(t *testing.T) {
	console.PrintFilesTable([][]string{}, "No Files")
}

func TestEcho_nilWriter(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	// nil writer falls back to os.Stdout -- just ensure no panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	console.Echo(nil, "msg", "", "", false)
}

func TestEcho_unknownSymbol(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	// Unknown symbol keys should not appear as prefix.
	console.Echo(&buf, "msg", "", "notasymbol", false)
	if !strings.Contains(buf.String(), "msg") {
		t.Errorf("expected 'msg' in output, got %q", buf.String())
	}
}

func TestEcho_boldFlag(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	console.Echo(&buf, "bold text", "green", "", true)
	if !strings.Contains(buf.String(), "bold text") {
		t.Errorf("expected 'bold text' in output, got %q", buf.String())
	}
}

func TestStatusSymbols_extraKeys(t *testing.T) {
	extras := []string{"running", "gear", "cross", "list", "preview", "download", "update", "remove"}
	for _, k := range extras {
		if v := console.StatusSymbols[k]; v == "" {
			t.Errorf("StatusSymbols[%q] is empty", k)
		}
	}
}

func TestDownloadSpinner_smoke(t *testing.T) {
	called := false
	console.DownloadSpinner("test-repo", func() { called = true })
	if !called {
		t.Error("expected callback to be called")
	}
}

func TestPanel_noTitle(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panel panicked: %v", r)
		}
	}()
	console.Panel("content here", "", "default")
}

func TestPanel_withTitle(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panel panicked: %v", r)
		}
	}()
	console.Panel("content here", "Section Title", "default")
}
