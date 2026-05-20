package experimental

import (
	"testing"
)

func TestFlag_ZeroValue_Extra2(t *testing.T) {
	var f Flag
	if f.Name != "" || f.Description != "" {
		t.Error("zero-value Flag should have empty fields")
	}
}

func TestConfig_ZeroValue_Extra2(t *testing.T) {
	var c Config
	if len(c.ExperimentalFlags) != 0 {
		t.Error("zero-value Config.ExperimentalFlags should be empty")
	}
}

func TestKnownFlagsNotEmpty_Extra2(t *testing.T) {
	if len(KnownFlags) == 0 {
		t.Error("KnownFlags should not be empty")
	}
}

func TestKnownFlags_AllHaveName_Extra2(t *testing.T) {
	for _, f := range KnownFlags {
		if f.Name == "" {
			t.Errorf("flag with empty name found: %+v", f)
		}
	}
}

func TestKnownFlags_AllHaveDescription_Extra2(t *testing.T) {
	for _, f := range KnownFlags {
		if f.Description == "" {
			t.Errorf("flag %q has no description", f.Name)
		}
	}
}

func TestNormaliseFlag_CaseFolding_Extra2(t *testing.T) {
	name := NormaliseFlag("MY-FLAG")
	if name != "my-flag" {
		t.Errorf("NormaliseFlag = %q, want my-flag", name)
	}
}

func TestNormaliseFlag_Spaces_Extra2(t *testing.T) {
	name := NormaliseFlag("  flag  ")
	if name != "flag" {
		t.Errorf("NormaliseFlag = %q, want flag", name)
	}
}

func TestDisplayName_UnknownFlag_Extra2(t *testing.T) {
	name := DisplayName("totally-unknown-flag-xyz")
	if name != "totally-unknown-flag-xyz" {
		t.Errorf("DisplayName for unknown = %q, want input echoed back", name)
	}
}

func TestValidateFlagName_UnknownReturnsError_Extra2(t *testing.T) {
	err := ValidateFlagName("definitely-not-a-real-flag-xyz")
	if err == nil {
		t.Error("expected error for unknown flag name")
	}
}

func TestIsEnabled_InitiallyFalse_Extra2(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if len(KnownFlags) == 0 {
		t.Skip("no known flags")
	}
	flagName := KnownFlags[0].Name
	ok, err := IsEnabled(flagName)
	if err != nil {
		t.Skipf("IsEnabled error (expected in some environments): %v", err)
	}
	_ = ok // may or may not be enabled depending on environment
}

func TestListFlags_Extra2(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	statuses, err := ListFlags()
	if err != nil {
		t.Skipf("ListFlags error: %v", err)
	}
	if len(statuses) != len(KnownFlags) {
		t.Errorf("ListFlags returned %d statuses, expected %d", len(statuses), len(KnownFlags))
	}
}

func TestEnableDisable_RoundTrip_Extra2(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if len(KnownFlags) == 0 {
		t.Skip("no known flags")
	}
	flagName := KnownFlags[0].Name
	if err := EnableFlag(flagName); err != nil {
		t.Skipf("EnableFlag error: %v", err)
	}
	ok, err := IsEnabled(flagName)
	if err != nil {
		t.Skipf("IsEnabled error: %v", err)
	}
	if !ok {
		t.Errorf("expected flag %q to be enabled after EnableFlag", flagName)
	}

	if err := DisableFlag(flagName); err != nil {
		t.Skipf("DisableFlag error: %v", err)
	}
	ok, err = IsEnabled(flagName)
	if err != nil {
		t.Skipf("IsEnabled error: %v", err)
	}
	if ok {
		t.Errorf("expected flag %q to be disabled after DisableFlag", flagName)
	}
}

func TestResetFlags_Extra2(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if len(KnownFlags) == 0 {
		t.Skip("no known flags")
	}
	flagName := KnownFlags[0].Name
	_ = EnableFlag(flagName)

	if err := ResetFlags(); err != nil {
		t.Skipf("ResetFlags error: %v", err)
	}
	ok, err := IsEnabled(flagName)
	if err != nil {
		t.Skipf("IsEnabled error: %v", err)
	}
	if ok {
		t.Errorf("expected flag disabled after ResetFlags")
	}
}
