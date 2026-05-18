package experimental

import (
	"os"
	"testing"
)

func TestGetOverriddenFlags_Empty(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	flags, err := GetOverriddenFlags()
	if err != nil {
		t.Fatalf("GetOverriddenFlags error: %v", err)
	}
	if len(flags) != 0 {
		t.Errorf("expected 0 overridden flags on fresh config, got %d", len(flags))
	}
}

func TestGetOverriddenFlags_AfterEnable(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	if err := EnableFlag(KnownFlags[0].Name); err != nil {
		t.Fatalf("EnableFlag: %v", err)
	}
	flags, err := GetOverriddenFlags()
	if err != nil {
		t.Fatalf("GetOverriddenFlags: %v", err)
	}
	if len(flags) == 0 {
		t.Error("expected at least one overridden flag after enable")
	}
}

func TestGetMalformedFlagKeys_EmptyOnFreshConfig(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	keys, err := GetMalformedFlagKeys()
	if err != nil {
		t.Fatalf("GetMalformedFlagKeys: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 malformed keys on fresh config, got %d", len(keys))
	}
}

func TestNormaliseFlag_LowerCase(t *testing.T) {
	if got := NormaliseFlag("FOO-BAR"); got != "foo-bar" {
		t.Errorf("NormaliseFlag(FOO-BAR) = %q, want foo-bar", got)
	}
}

func TestNormaliseFlag_TrimSpaces(t *testing.T) {
	if got := NormaliseFlag("  foo  "); got != "foo" {
		t.Errorf("NormaliseFlag with spaces = %q, want foo", got)
	}
}

func TestNormaliseFlag_AlreadyNormal(t *testing.T) {
	if got := NormaliseFlag("parallel-install"); got != "parallel-install" {
		t.Errorf("NormaliseFlag(already-normal) = %q, want unchanged", got)
	}
}

func TestDisplayName_AllKnownFlags(t *testing.T) {
	for _, f := range KnownFlags {
		got := DisplayName(f.Name)
		if got != f.DisplayName {
			t.Errorf("DisplayName(%q) = %q, want %q", f.Name, got, f.DisplayName)
		}
	}
}

func TestDisplayName_Unknown_ReturnsInput(t *testing.T) {
	if got := DisplayName("no-such-flag-xyz"); got != "no-such-flag-xyz" {
		t.Errorf("DisplayName(unknown) = %q, want input unchanged", got)
	}
}

func TestValidateFlagName_AllKnown(t *testing.T) {
	for _, f := range KnownFlags {
		if err := ValidateFlagName(f.Name); err != nil {
			t.Errorf("ValidateFlagName(%q) unexpected error: %v", f.Name, err)
		}
	}
}

func TestValidateFlagName_Unknown(t *testing.T) {
	if err := ValidateFlagName("totally-unknown-flag"); err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestEnableDisable_Idempotent(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	flag := KnownFlags[0].Name
	if err := EnableFlag(flag); err != nil {
		t.Fatalf("first enable: %v", err)
	}
	if err := EnableFlag(flag); err != nil {
		t.Fatalf("second enable should not error: %v", err)
	}
	en, _ := IsEnabled(flag)
	if !en {
		t.Error("expected enabled after double enable")
	}
}

func TestListFlags_ContainsAllKnown(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	statuses, err := ListFlags()
	if err != nil {
		t.Fatalf("ListFlags: %v", err)
	}
	names := map[string]bool{}
	for _, s := range statuses {
		names[s.Name] = true
	}
	for _, f := range KnownFlags {
		if !names[f.Name] {
			t.Errorf("ListFlags missing known flag: %q", f.Name)
		}
	}
}

func TestResetFlags_ClearsEnabled(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	for _, f := range KnownFlags {
		_ = EnableFlag(f.Name)
	}
	if err := ResetFlags(); err != nil {
		t.Fatalf("ResetFlags: %v", err)
	}
	for _, f := range KnownFlags {
		en, _ := IsEnabled(f.Name)
		if en != f.Default {
			t.Errorf("after reset, IsEnabled(%q) = %v, want default %v", f.Name, en, f.Default)
		}
	}
}

func TestKnownFlags_AllHaveDescription(t *testing.T) {
	for _, f := range KnownFlags {
		if f.Description == "" {
			t.Errorf("flag %q has empty description", f.Name)
		}
	}
}
