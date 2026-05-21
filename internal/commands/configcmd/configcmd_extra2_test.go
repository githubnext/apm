package configcmd_test

import (
	"testing"

	"github.com/githubnext/apm/internal/commands/configcmd"
)

func TestParseBoolValue_TrueLower(t *testing.T) {
	v, err := configcmd.ParseBoolValue("true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true")
	}
}

func TestParseBoolValue_FalseLower(t *testing.T) {
	v, err := configcmd.ParseBoolValue("false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v {
		t.Error("expected false")
	}
}

func TestParseBoolValue_One(t *testing.T) {
	v, err := configcmd.ParseBoolValue("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true for '1'")
	}
}

func TestParseBoolValue_Zero(t *testing.T) {
	v, err := configcmd.ParseBoolValue("0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v {
		t.Error("expected false for '0'")
	}
}

func TestParseBoolValue_Yes(t *testing.T) {
	v, err := configcmd.ParseBoolValue("yes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true for 'yes'")
	}
}

func TestParseBoolValue_No(t *testing.T) {
	v, err := configcmd.ParseBoolValue("no")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v {
		t.Error("expected false for 'no'")
	}
}

func TestParseBoolValue_InvalidOn(t *testing.T) {
	_, err := configcmd.ParseBoolValue("on")
	if err == nil {
		t.Error("expected error for unsupported value 'on'")
	}
}

func TestParseBoolValue_InvalidOff(t *testing.T) {
	_, err := configcmd.ParseBoolValue("off")
	if err == nil {
		t.Error("expected error for unsupported value 'off'")
	}
}

func TestValidConfigKeys_NoDuplicates(t *testing.T) {
	keys := configcmd.ValidConfigKeys()
	seen := map[string]bool{}
	for _, k := range keys {
		if seen[k] {
			t.Errorf("duplicate key: %q", k)
		}
		seen[k] = true
	}
}

func TestValidConfigKeys_AllNonEmpty(t *testing.T) {
	for _, k := range configcmd.ValidConfigKeys() {
		if k == "" {
			t.Error("found empty key in ValidConfigKeys")
		}
	}
}

func TestDisplayName_AutoIntegrate(t *testing.T) {
	got := configcmd.DisplayName("auto_integrate")
	if got == "" {
		t.Error("expected non-empty display name for auto_integrate")
	}
}

func TestDisplayName_TempDir(t *testing.T) {
	got := configcmd.DisplayName("temp_dir")
	if got == "" {
		t.Error("expected non-empty display name for temp_dir")
	}
}

func TestDisplayName_Unknown(t *testing.T) {
	got := configcmd.DisplayName("nonexistent_key_xyz")
	if got == "" {
		t.Error("expected fallback display name for unknown key")
	}
}

func TestDisplayName_AllValidKeys(t *testing.T) {
	for _, k := range configcmd.ValidConfigKeys() {
		got := configcmd.DisplayName(k)
		if got == "" {
			t.Errorf("DisplayName(%q) returned empty string", k)
		}
	}
}
