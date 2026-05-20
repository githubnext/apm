package console_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/console"
)

func TestStatusSymbols_CheckKeyBracket(t *testing.T) {
	sym, ok := console.StatusSymbols["check"]
	if !ok {
		t.Fatal("expected 'check' key in StatusSymbols")
	}
	if sym != "[+]" {
		t.Errorf("expected '[+]' for 'check', got %q", sym)
	}
}

func TestStatusSymbols_InfoKeyBracket(t *testing.T) {
	sym, ok := console.StatusSymbols["info"]
	if !ok {
		t.Fatal("expected 'info' key in StatusSymbols")
	}
	if sym != "[i]" {
		t.Errorf("expected '[i]' for 'info', got %q", sym)
	}
}

func TestStatusSymbols_WarningKeyBracket(t *testing.T) {
	sym, ok := console.StatusSymbols["warning"]
	if !ok {
		t.Fatal("expected 'warning' key in StatusSymbols")
	}
	if sym != "[!]" {
		t.Errorf("expected '[!]' for 'warning', got %q", sym)
	}
}

func TestStatusSymbols_DefaultKey(t *testing.T) {
	sym, ok := console.StatusSymbols["default"]
	if !ok {
		t.Fatal("expected 'default' key in StatusSymbols")
	}
	if sym == "" {
		t.Error("expected non-empty symbol for 'default'")
	}
}

func TestStatusSymbols_UpdateKeyTilde(t *testing.T) {
	sym, ok := console.StatusSymbols["update"]
	if !ok {
		t.Fatal("expected 'update' key in StatusSymbols")
	}
	if sym != "[~]" {
		t.Errorf("expected '[~]' for 'update', got %q", sym)
	}
}

func TestStatusSymbols_RemoveKeyDash(t *testing.T) {
	sym, ok := console.StatusSymbols["remove"]
	if !ok {
		t.Fatal("expected 'remove' key in StatusSymbols")
	}
	if sym != "[-]" {
		t.Errorf("expected '[-]' for 'remove', got %q", sym)
	}
}

func TestStatusSymbols_EqualKeyEquals(t *testing.T) {
	sym, ok := console.StatusSymbols["equal"]
	if !ok {
		t.Fatal("expected 'equal' key in StatusSymbols")
	}
	if sym != "[=]" {
		t.Errorf("expected '[=]' for 'equal', got %q", sym)
	}
}

func TestStatusSymbols_CrossKey(t *testing.T) {
	sym, ok := console.StatusSymbols["cross"]
	if !ok {
		t.Fatal("expected 'cross' key in StatusSymbols")
	}
	if sym != "[x]" {
		t.Errorf("expected '[x]' for 'cross', got %q", sym)
	}
}

func TestStatusSymbols_MissingKeyReturnsFalse(t *testing.T) {
	_, ok := console.StatusSymbols["nonexistent_key_xyz"]
	if ok {
		t.Error("expected false for nonexistent key")
	}
}

func TestStatusSymbols_RunningKey(t *testing.T) {
	sym, ok := console.StatusSymbols["running"]
	if !ok {
		t.Fatal("expected 'running' key in StatusSymbols")
	}
	if sym != "[>]" {
		t.Errorf("expected '[>]' for 'running', got %q", sym)
	}
}

func TestStatusSymbols_ListKey(t *testing.T) {
	sym, ok := console.StatusSymbols["list"]
	if !ok {
		t.Fatal("expected 'list' key in StatusSymbols")
	}
	if sym != "[#]" {
		t.Errorf("expected '[#]' for 'list', got %q", sym)
	}
}

func TestStatusSymbols_MetricsKey(t *testing.T) {
	sym, ok := console.StatusSymbols["metrics"]
	if !ok {
		t.Fatal("expected 'metrics' key in StatusSymbols")
	}
	if sym != "[#]" {
		t.Errorf("expected '[#]' for 'metrics', got %q", sym)
	}
}

func TestStatusSymbols_AllSymbolsNonEmpty(t *testing.T) {
	for k, v := range console.StatusSymbols {
		if v == "" {
			t.Errorf("StatusSymbols[%q] is empty", k)
		}
	}
}

func TestStatusSymbols_AllSymbolsASCII(t *testing.T) {
	for k, v := range console.StatusSymbols {
		for _, r := range v {
			if r > 0x7E || r < 0x20 {
				t.Errorf("StatusSymbols[%q] = %q contains non-printable-ASCII char %U", k, v, r)
			}
		}
	}
}
