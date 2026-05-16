package experimental_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/experimental"
)

func TestFlags(t *testing.T) {
	flags := experimental.Flags()
	if len(flags) == 0 {
		t.Error("expected at least one registered flag")
	}
	for name, flag := range flags {
		if flag.Name != name {
			t.Errorf("flag key %q but flag.Name %q mismatch", name, flag.Name)
		}
		if flag.Description == "" {
			t.Errorf("flag %q has empty description", name)
		}
		if flag.Default {
			t.Errorf("flag %q default should be false", name)
		}
	}
}

func TestFlagsContainsKnownFlags(t *testing.T) {
	flags := experimental.Flags()
	expected := []string{"verbose_version", "copilot_cowork"}
	for _, name := range expected {
		if _, ok := flags[name]; !ok {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestDisplayName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"verbose_version", "verbose-version"},
		{"copilot_cowork", "copilot-cowork"},
		{"no_underscores", "no-underscores"},
		{"simple", "simple"},
		{"", ""},
	}
	for _, tc := range cases {
		got := experimental.DisplayName(tc.in)
		if got != tc.want {
			t.Errorf("DisplayName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
