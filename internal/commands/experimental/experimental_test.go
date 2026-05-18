package experimental

import (
	"os"
	"testing"
)

func TestKnownFlagsNotEmpty(t *testing.T) {
	if len(KnownFlags) == 0 {
		t.Error("KnownFlags is empty")
	}
}

func TestKnownFlagsNames(t *testing.T) {
	seen := map[string]bool{}
	for _, f := range KnownFlags {
		if f.Name == "" {
			t.Errorf("Flag has empty name: %+v", f)
		}
		if seen[f.Name] {
			t.Errorf("Duplicate flag name: %q", f.Name)
		}
		seen[f.Name] = true
	}
}

func TestNormaliseFlag(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Parallel-Install", "parallel-install"},
		{"  telemetry  ", "telemetry"},
		{"STRICT-POLICY", "strict-policy"},
	}
	for _, c := range cases {
		if got := NormaliseFlag(c.in); got != c.want {
			t.Errorf("NormaliseFlag(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDisplayName(t *testing.T) {
	// Known flag returns DisplayName.
	for _, f := range KnownFlags {
		if got := DisplayName(f.Name); got != f.DisplayName {
			t.Errorf("DisplayName(%q) = %q, want %q", f.Name, got, f.DisplayName)
		}
	}
	// Unknown flag returns the input unchanged.
	if got := DisplayName("nonexistent-flag"); got != "nonexistent-flag" {
		t.Errorf("DisplayName(unknown) = %q, want input as-is", got)
	}
}

func TestValidateFlagName(t *testing.T) {
	// A known flag name is valid.
	if err := ValidateFlagName(KnownFlags[0].Name); err != nil {
		t.Errorf("ValidateFlagName(known) error: %v", err)
	}
	// An unknown flag name should error.
	if err := ValidateFlagName("not-a-real-flag-xyz"); err == nil {
		t.Error("ValidateFlagName(unknown) = nil, want error")
	}
}

func TestIsEnabledAndEnableDisable(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	flag := KnownFlags[0].Name
	enabled, err := IsEnabled(flag)
	if err != nil {
		t.Fatalf("IsEnabled error: %v", err)
	}
	if enabled != KnownFlags[0].Default {
		t.Errorf("IsEnabled(%q) before enable = %v, want default %v", flag, enabled, KnownFlags[0].Default)
	}

	if err := EnableFlag(flag); err != nil {
		t.Fatalf("EnableFlag error: %v", err)
	}
	enabled, _ = IsEnabled(flag)
	if !enabled {
		t.Errorf("IsEnabled(%q) after enable = false, want true", flag)
	}

	if err := DisableFlag(flag); err != nil {
		t.Fatalf("DisableFlag error: %v", err)
	}
	enabled, _ = IsEnabled(flag)
	if enabled {
		t.Errorf("IsEnabled(%q) after disable = true, want false", flag)
	}
}

func TestResetFlags(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	_ = EnableFlag(KnownFlags[0].Name)
	if err := ResetFlags(); err != nil {
		t.Fatalf("ResetFlags error: %v", err)
	}
}

func TestListFlags(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	statuses, err := ListFlags()
	if err != nil {
		t.Fatalf("ListFlags error: %v", err)
	}
	if len(statuses) != len(KnownFlags) {
		t.Errorf("ListFlags() len = %d, want %d", len(statuses), len(KnownFlags))
	}
}
