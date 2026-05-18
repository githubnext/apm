package console_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/console"
)

func TestStatusSymbols_SuccessKey(t *testing.T) {
	sym, ok := console.StatusSymbols["success"]
	if !ok {
		t.Fatal("expected 'success' key in StatusSymbols")
	}
	if sym == "" {
		t.Error("expected non-empty symbol for 'success'")
	}
}

func TestStatusSymbols_ErrorKey(t *testing.T) {
	sym, ok := console.StatusSymbols["error"]
	if !ok {
		t.Fatal("expected 'error' key in StatusSymbols")
	}
	if sym != "[x]" {
		t.Errorf("expected '[x]' for 'error', got %q", sym)
	}
}

func TestStatusSymbols_WarningKey(t *testing.T) {
	sym, ok := console.StatusSymbols["warning"]
	if !ok {
		t.Fatal("expected 'warning' key in StatusSymbols")
	}
	if sym != "[!]" {
		t.Errorf("expected '[!]' for 'warning', got %q", sym)
	}
}

func TestStatusSymbols_InfoKey(t *testing.T) {
	sym, ok := console.StatusSymbols["info"]
	if !ok {
		t.Fatal("expected 'info' key in StatusSymbols")
	}
	if sym != "[i]" {
		t.Errorf("expected '[i]' for 'info', got %q", sym)
	}
}

func TestStatusSymbols_CheckKey(t *testing.T) {
	sym, ok := console.StatusSymbols["check"]
	if !ok {
		t.Fatal("expected 'check' key in StatusSymbols")
	}
	if sym != "[+]" {
		t.Errorf("expected '[+]' for 'check', got %q", sym)
	}
}

func TestStatusSymbols_AllASCII(t *testing.T) {
	for key, sym := range console.StatusSymbols {
		for _, c := range sym {
			if c > 0x7E || c < 0x20 {
				t.Errorf("symbol for %q contains non-ASCII char %q", key, c)
			}
		}
	}
}

func TestEcho_SymbolPrefix(t *testing.T) {
	var sb strings.Builder
	console.Echo(&sb, "hello", "", "success", false)
	got := sb.String()
	if !strings.Contains(got, "hello") {
		t.Errorf("expected 'hello' in output, got %q", got)
	}
	sym := console.StatusSymbols["success"]
	if !strings.Contains(got, sym) {
		t.Errorf("expected symbol %q in output, got %q", sym, got)
	}
}

func TestEcho_UnknownSymbol_NoPrefix(t *testing.T) {
	var sb strings.Builder
	console.Echo(&sb, "world", "", "no-such-symbol", false)
	got := sb.String()
	if !strings.Contains(got, "world") {
		t.Errorf("expected 'world' in output, got %q", got)
	}
}

func TestEcho_EmptyMessage(t *testing.T) {
	var sb strings.Builder
	console.Echo(&sb, "", "", "", false)
	got := sb.String()
	if got != "\n" {
		t.Errorf("expected single newline for empty message, got %q", got)
	}
}

func TestEcho_MultipleSymbols(t *testing.T) {
	syms := []string{"success", "error", "warning", "info", "check", "running"}
	for _, sym := range syms {
		var sb strings.Builder
		console.Echo(&sb, "msg", "", sym, false)
		got := sb.String()
		expected := console.StatusSymbols[sym]
		if !strings.Contains(got, expected) {
			t.Errorf("symbol %q: expected %q in output, got %q", sym, expected, got)
		}
	}
}

func TestPrintFilesTable_OneRow(t *testing.T) {
	// smoke test: does not panic
	console.PrintFilesTable([][]string{{"file.go", "description"}}, "Title")
}

func TestPrintFilesTable_ShortRow(t *testing.T) {
	// row with single element (no description)
	console.PrintFilesTable([][]string{{"file.go"}}, "")
}

func TestPrintFilesTable_EmptyRow(t *testing.T) {
	// row with no elements
	console.PrintFilesTable([][]string{{}}, "")
}

func TestPrintFilesTable_MultipleRows(t *testing.T) {
	rows := [][]string{
		{"alpha.go", "desc a"},
		{"beta.go", "desc b"},
		{"gamma.go", "desc c"},
	}
	console.PrintFilesTable(rows, "Files")
}

func TestDownloadSpinner_Invokes(t *testing.T) {
	called := false
	console.DownloadSpinner("my-repo", func() {
		called = true
	})
	if !called {
		t.Error("expected callback to be invoked")
	}
}

func TestDownloadSpinner_EmptyName(t *testing.T) {
	// should not panic
	console.DownloadSpinner("", func() {})
}

func TestPanel_LongTitle(t *testing.T) {
	// Should not panic on long title
	console.Panel("content", strings.Repeat("x", 200), "")
}

func TestPanel_NoTitle_NoSeparator(t *testing.T) {
	// should not panic and just print content
	console.Panel("only content", "", "")
}

func TestPanel_SpecialChars(t *testing.T) {
	// special characters in content are fine
	console.Panel("line1\nline2", "Title", "bold")
}

func TestEcho_InfoSymbol(t *testing.T) {
	var sb strings.Builder
	console.Echo(&sb, "test", "", "info", false)
	got := sb.String()
	if !strings.Contains(got, "[i]") {
		t.Errorf("expected '[i]' symbol for 'info', got %q", got)
	}
}

func TestEcho_GearSymbol(t *testing.T) {
	var sb strings.Builder
	console.Echo(&sb, "test", "", "gear", false)
	got := sb.String()
	if !strings.Contains(got, "test") {
		t.Errorf("expected 'test' in output, got %q", got)
	}
}

func TestEcho_DefaultSymbol(t *testing.T) {
	var sb strings.Builder
	console.Echo(&sb, "test", "", "default", false)
	got := sb.String()
	if !strings.Contains(got, console.StatusSymbols["default"]) {
		t.Errorf("expected default symbol in output, got %q", got)
	}
}

func TestStatusSymbols_UpdateKey(t *testing.T) {
	sym, ok := console.StatusSymbols["update"]
	if !ok {
		t.Fatal("expected 'update' key")
	}
	if sym != "[~]" {
		t.Errorf("expected '[~]' for 'update', got %q", sym)
	}
}

func TestStatusSymbols_RemoveKey(t *testing.T) {
	sym, ok := console.StatusSymbols["remove"]
	if !ok {
		t.Fatal("expected 'remove' key")
	}
	if sym != "[-]" {
		t.Errorf("expected '[-]' for 'remove', got %q", sym)
	}
}

func TestStatusSymbols_EqualKey(t *testing.T) {
	sym, ok := console.StatusSymbols["equal"]
	if !ok {
		t.Fatal("expected 'equal' key")
	}
	if sym != "[=]" {
		t.Errorf("expected '[=]' for 'equal', got %q", sym)
	}
}
