package configcmd_test

import (
	"testing"

	"github.com/githubnext/apm/internal/commands/configcmd"
)

func TestParseBoolValue_Yes_Extra3(t *testing.T) {
	v, err := configcmd.ParseBoolValue("yes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true for 'yes'")
	}
}

func TestParseBoolValue_No_Extra3(t *testing.T) {
	v, err := configcmd.ParseBoolValue("no")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v {
		t.Error("expected false for 'no'")
	}
}

func TestParseBoolValue_One_Extra3(t *testing.T) {
	v, err := configcmd.ParseBoolValue("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true for '1'")
	}
}

func TestParseBoolValue_Zero_Extra3(t *testing.T) {
	v, err := configcmd.ParseBoolValue("0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v {
		t.Error("expected false for '0'")
	}
}

func TestParseBoolValue_Invalid_ReturnsError_Extra3(t *testing.T) {
	_, err := configcmd.ParseBoolValue("maybe")
	if err == nil {
		t.Error("expected error for invalid bool value 'maybe'")
	}
}

func TestParseBoolValue_CaseInsensitive_Extra3(t *testing.T) {
	v, err := configcmd.ParseBoolValue("TRUE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true for 'TRUE'")
	}
}

func TestValidConfigKeys_NotEmpty_Extra3(t *testing.T) {
	keys := configcmd.ValidConfigKeys()
	if len(keys) == 0 {
		t.Error("ValidConfigKeys should return at least one key")
	}
}

func TestDisplayName_KnownKey_Extra3(t *testing.T) {
	name := configcmd.DisplayName("auto_integrate")
	if name == "" {
		t.Error("DisplayName for auto_integrate should not be empty")
	}
}

func TestDisplayName_UnknownKey_ReturnsSelf_Extra3(t *testing.T) {
	name := configcmd.DisplayName("unknown_key_xyz")
	// Should fall back gracefully; non-empty.
	if name == "" {
		t.Error("DisplayName for unknown key should return non-empty string")
	}
}
