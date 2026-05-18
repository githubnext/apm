package experimental_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/experimental"
)

func TestDisplayName_SingleWord(t *testing.T) {
	got := experimental.DisplayName("feature")
	if got != "feature" {
		t.Errorf("DisplayName(no underscores) = %q, want feature", got)
	}
}

func TestDisplayName_EmptyString(t *testing.T) {
	got := experimental.DisplayName("")
	if got != "" {
		t.Errorf("DisplayName('') = %q, want empty", got)
	}
}

func TestFlags_NoDuplicateNames(t *testing.T) {
	flags := experimental.Flags()
	seen := map[string]bool{}
	for k, f := range flags {
		if seen[f.Name] {
			t.Errorf("duplicate flag.Name %q in registry", f.Name)
		}
		seen[f.Name] = true
		if k != f.Name {
			t.Errorf("key %q does not match flag.Name %q", k, f.Name)
		}
	}
}

func TestFlags_AllDefaultFalse(t *testing.T) {
	flags := experimental.Flags()
	for name, f := range flags {
		if f.Default {
			t.Errorf("flag %q has Default=true, expected false", name)
		}
	}
}

func TestFlags_AllHaveDescription(t *testing.T) {
	flags := experimental.Flags()
	for name, f := range flags {
		if f.Description == "" {
			t.Errorf("flag %q has empty Description", name)
		}
	}
}

func TestFlagsConsistentLength(t *testing.T) {
	// Two successive calls should return the same number of flags.
	flags1 := experimental.Flags()
	flags2 := experimental.Flags()
	if len(flags1) != len(flags2) {
		t.Errorf("Flags() returned different lengths: %d vs %d", len(flags1), len(flags2))
	}
}

func TestDisplayName_TrailingUnderscore(t *testing.T) {
	got := experimental.DisplayName("flag_")
	if got != "flag-" {
		t.Errorf("DisplayName('flag_') = %q, want 'flag-'", got)
	}
}

func TestDisplayName_MultiUnderscores(t *testing.T) {
	got := experimental.DisplayName("a__b")
	if got != "a--b" {
		t.Errorf("DisplayName('a__b') = %q, want 'a--b'", got)
	}
}

func TestFlagHintNotEmpty(t *testing.T) {
	flags := experimental.Flags()
	for name, f := range flags {
		if f.Hint == "" {
			t.Logf("flag %q has empty Hint (informational)", name)
		}
	}
	// At least one flag should have a hint
	hasHint := false
	for _, f := range flags {
		if f.Hint != "" {
			hasHint = true
		}
	}
	if !hasHint {
		t.Error("expected at least one flag with a non-empty Hint")
	}
}

func TestFlagsNamesAreSnakeCase(t *testing.T) {
	flags := experimental.Flags()
	for name := range flags {
		for _, ch := range name {
			if ch == '-' {
				t.Errorf("flag key %q contains hyphen; should use snake_case", name)
			}
		}
	}
}
