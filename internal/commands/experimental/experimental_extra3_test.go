package experimental

import (
	"testing"
)

func TestFlag_Fields_Extra3(t *testing.T) {
	f := Flag{Name: "my-flag", Description: "desc"}
	if f.Name != "my-flag" {
		t.Errorf("Name = %q, want my-flag", f.Name)
	}
	if f.Description != "desc" {
		t.Errorf("Description = %q, want desc", f.Description)
	}
}

func TestFlagStatus_Struct_Extra3(t *testing.T) {
	fs := FlagStatus{Enabled: true}
	if !fs.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestNormaliseFlag_AlreadyLower_Extra3(t *testing.T) {
	got := NormaliseFlag("my-flag")
	if got != "my-flag" {
		t.Errorf("NormaliseFlag returned %q, want my-flag", got)
	}
}

func TestNormaliseFlag_Empty_Extra3(t *testing.T) {
	got := NormaliseFlag("")
	if got != "" {
		t.Errorf("NormaliseFlag(\"\") = %q, want empty", got)
	}
}

func TestDisplayName_KnownFlag_Extra3(t *testing.T) {
	if len(KnownFlags) == 0 {
		t.Skip("no known flags")
	}
	name := KnownFlags[0].Name
	dn := DisplayName(name)
	if dn == "" {
		t.Errorf("DisplayName(%q) returned empty string", name)
	}
}

func TestValidateFlagName_EmptyName_Extra3(t *testing.T) {
	err := ValidateFlagName("")
	if err == nil {
		t.Error("expected error for empty flag name")
	}
}

func TestKnownFlags_NoDuplicateNames_Extra3(t *testing.T) {
	seen := map[string]bool{}
	for _, f := range KnownFlags {
		if seen[f.Name] {
			t.Errorf("duplicate flag name: %q", f.Name)
		}
		seen[f.Name] = true
	}
}

func TestGetOverriddenFlags_Type_Extra3(t *testing.T) {
	flags, err := GetOverriddenFlags()
	if err != nil {
		t.Fatalf("GetOverriddenFlags error: %v", err)
	}
	for _, f := range flags {
		if f.Name == "" {
			t.Error("overridden flag with empty name")
		}
	}
}

func TestGetStaleConfigKeys_Type_Extra3(t *testing.T) {
	keys, err := GetStaleConfigKeys()
	if err != nil {
		t.Fatalf("GetStaleConfigKeys error: %v", err)
	}
	// May be empty, just check no panics
	_ = keys
}

func TestGetMalformedFlagKeys_Same_Extra3(t *testing.T) {
	a, _ := GetMalformedFlagKeys()
	b, _ := GetStaleConfigKeys()
	if len(a) != len(b) {
		t.Errorf("GetMalformedFlagKeys and GetStaleConfigKeys should return same length: %d vs %d", len(a), len(b))
	}
}
