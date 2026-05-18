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

func TestFlagsAreImmutable(t *testing.T) {
	flags1 := experimental.Flags()
	flags2 := experimental.Flags()
	if len(flags1) != len(flags2) {
		t.Errorf("Flags() returned different lengths on repeated calls: %d vs %d", len(flags1), len(flags2))
	}
}

func TestFlagHasHint(t *testing.T) {
	flags := experimental.Flags()
	vv, ok := flags["verbose_version"]
	if !ok {
		t.Fatal("verbose_version not found")
	}
	if vv.Hint == "" {
		t.Error("verbose_version should have a non-empty Hint")
	}
}

func TestFlagCopilotCoworkHasHint(t *testing.T) {
	flags := experimental.Flags()
	cc, ok := flags["copilot_cowork"]
	if !ok {
		t.Fatal("copilot_cowork not found")
	}
	if cc.Hint == "" {
		t.Error("copilot_cowork should have a non-empty Hint")
	}
}

func TestDisplayNameMultipleUnderscores(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"a_b_c_d", "a-b-c-d"},
		{"_leading", "-leading"},
		{"trailing_", "trailing-"},
		{"__double__", "--double--"},
	}
	for _, tc := range cases {
		got := experimental.DisplayName(tc.in)
		if got != tc.want {
			t.Errorf("DisplayName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
