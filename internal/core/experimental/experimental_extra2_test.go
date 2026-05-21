package experimental_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core/experimental"
)

func TestFlags_MapNotNil(t *testing.T) {
	flags := experimental.Flags()
	if flags == nil {
		t.Error("Flags() returned nil map")
	}
}

func TestFlags_AtLeastOneFlag(t *testing.T) {
	flags := experimental.Flags()
	if len(flags) == 0 {
		t.Error("expected at least one experimental flag")
	}
}

func TestFlags_EachFlagHasName(t *testing.T) {
	for k, f := range experimental.Flags() {
		if f.Name == "" {
			t.Errorf("flag key=%q has empty Name field", k)
		}
	}
}

func TestFlags_EachFlagDefaultFalse(t *testing.T) {
	for k, f := range experimental.Flags() {
		if f.Default {
			t.Errorf("flag key=%q has Default=true, expected false", k)
		}
	}
}

func TestDisplayName_ConvertsUnderscoreToSpace(t *testing.T) {
	got := experimental.DisplayName("some_flag_name")
	if strings.Contains(got, "_") {
		t.Errorf("expected no underscores in display name, got %q", got)
	}
}

func TestDisplayName_NonEmpty(t *testing.T) {
	got := experimental.DisplayName("my_flag")
	if got == "" {
		t.Error("expected non-empty display name")
	}
}

func TestDisplayName_Empty(t *testing.T) {
	got := experimental.DisplayName("")
	_ = got
}

func TestFlags_NamesAreSnakeCase(t *testing.T) {
	for k := range experimental.Flags() {
		if strings.Contains(k, " ") {
			t.Errorf("flag key %q contains space (expected snake_case)", k)
		}
		if strings.ToLower(k) != k {
			t.Errorf("flag key %q is not lowercase", k)
		}
	}
}

func TestFlags_DescriptionsAreASCII(t *testing.T) {
	for k, f := range experimental.Flags() {
		for _, r := range f.Description {
			if r > 127 {
				t.Errorf("flag %q description contains non-ASCII character %q", k, string(r))
			}
		}
	}
}

func TestFlags_DescriptionsNotTooLong(t *testing.T) {
	for k, f := range experimental.Flags() {
		if len(f.Description) > 120 {
			t.Errorf("flag %q description is longer than 120 chars: %d", k, len(f.Description))
		}
	}
}

func TestDisplayName_KnownFlag(t *testing.T) {
	for k := range experimental.Flags() {
		got := experimental.DisplayName(k)
		if got == "" {
			t.Errorf("DisplayName(%q) returned empty string", k)
		}
		break
	}
}

func TestFlags_NameMatchesKey(t *testing.T) {
	for k, f := range experimental.Flags() {
		if f.Name != k {
			t.Errorf("key %q != flag.Name %q", k, f.Name)
		}
	}
}
