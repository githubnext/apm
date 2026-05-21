package experimental_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/experimental"
)

func TestFlags_AllHaveDescriptions(t *testing.T) {
	flags := experimental.Flags()
	for name, f := range flags {
		if f.Description == "" {
			t.Errorf("flag %q has empty description", name)
		}
	}
}

func TestDisplayName_MultiWord(t *testing.T) {
	result := experimental.DisplayName("my-feature-flag")
	if result == "" {
		t.Error("expected non-empty display name")
	}
}

func TestDisplayName_WithHyphen(t *testing.T) {
	result := experimental.DisplayName("foo-bar")
	if result == "" {
		t.Error("expected non-empty display name for hyphenated name")
	}
}

func TestFlags_CountAtLeastOne_Extra3(t *testing.T) {
	flags := experimental.Flags()
	if len(flags) == 0 {
		t.Error("expected at least one experimental flag")
	}
}

func TestValidateFlagName_ValidFlag(t *testing.T) {
	flags := experimental.Flags()
	for name := range flags {
		canonical, err := experimental.ValidateFlagName(name)
		if err != nil {
			t.Errorf("valid flag %q failed validation: %v", name, err)
		}
		if canonical == "" {
			t.Errorf("canonical name for %q is empty", name)
		}
		break
	}
}

func TestValidateFlagName_InvalidFlag(t *testing.T) {
	_, err := experimental.ValidateFlagName("nonexistent-flag-xyz")
	if err == nil {
		t.Error("expected error for unknown flag name")
	}
}
