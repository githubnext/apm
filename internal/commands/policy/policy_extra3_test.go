package policy

import (
	"strings"
	"testing"
)

func TestPolicySource_Stale_Extra3(t *testing.T) {
	s := PolicySource{Stale: true, Label: "local"}
	if !s.Stale {
		t.Error("Stale should be true")
	}
	if s.Label != "local" {
		t.Errorf("Label = %q, want local", s.Label)
	}
}

func TestPolicySource_CacheAge_Extra3(t *testing.T) {
	s := PolicySource{CacheAge: 3600}
	if s.CacheAge != 3600 {
		t.Errorf("CacheAge = %d, want 3600", s.CacheAge)
	}
}

func TestPolicySource_FilePath_Extra3(t *testing.T) {
	s := PolicySource{FilePath: "/etc/policy.yml"}
	if s.FilePath != "/etc/policy.yml" {
		t.Errorf("FilePath = %q", s.FilePath)
	}
}

func TestPolicyStatus_InheritanceChain_Extra3(t *testing.T) {
	status := PolicyStatus{
		Discovered: true,
		InheritanceChain: []PolicySource{
			{Label: "parent"},
			{Label: "child"},
		},
	}
	if len(status.InheritanceChain) != 2 {
		t.Errorf("InheritanceChain len = %d, want 2", len(status.InheritanceChain))
	}
	if status.InheritanceChain[0].Label != "parent" {
		t.Errorf("chain[0] = %q, want parent", status.InheritanceChain[0].Label)
	}
}

func TestPolicyStatus_RuleCount_Extra3(t *testing.T) {
	status := PolicyStatus{
		RuleCount: map[string]int{"allow": 5, "deny": 2},
	}
	if status.RuleCount["allow"] != 5 {
		t.Errorf("allow count = %d, want 5", status.RuleCount["allow"])
	}
}

func TestStatusOptions_NoFetch_Extra3(t *testing.T) {
	opts := StatusOptions{NoFetch: true, Format: "json"}
	if !opts.NoFetch {
		t.Error("NoFetch should be true")
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want json", opts.Format)
	}
}

func TestStatusOptions_Verbose_Extra3(t *testing.T) {
	opts := StatusOptions{Verbose: true}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestFormatAge_Seconds_Extra3(t *testing.T) {
	s := formatAge(45)
	if s == "" {
		t.Error("formatAge(45) should not be empty")
	}
}

func TestFormatAge_Minutes_Extra3(t *testing.T) {
	s := formatAge(120)
	if !strings.Contains(s, "min") && !strings.Contains(s, "2") {
		t.Errorf("formatAge(120) = %q should indicate minutes", s)
	}
}

func TestFormatAge_Hours_Extra3(t *testing.T) {
	s := formatAge(7200)
	if !strings.Contains(s, "h") && !strings.Contains(s, "2") {
		t.Errorf("formatAge(7200) = %q should indicate hours", s)
	}
}

func TestStripSourcePrefix_Variants_Extra3(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"org:myorg", "myorg"},
		{"url:https://example.com", "https://example.com"},
		{"file:/etc/policy.yml", "/etc/policy.yml"},
		{"nosource", "nosource"},
	}
	for _, c := range cases {
		got := stripSourcePrefix(c.in)
		if got != c.out {
			t.Errorf("stripSourcePrefix(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}
