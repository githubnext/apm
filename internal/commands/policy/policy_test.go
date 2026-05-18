package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPolicySourceFields(t *testing.T) {
	ps := PolicySource{
		Label:    "mypolicy",
		URL:      "https://example.com/policy.yml",
		FilePath: "/tmp/policy.yml",
		CacheAge: 120,
		Stale:    true,
	}
	if ps.Label != "mypolicy" {
		t.Errorf("Label = %q, want mypolicy", ps.Label)
	}
	if !ps.Stale {
		t.Error("Stale should be true")
	}
	if ps.CacheAge != 120 {
		t.Errorf("CacheAge = %d, want 120", ps.CacheAge)
	}
}

func TestPolicyStatusFields(t *testing.T) {
	s := &PolicyStatus{
		Discovered:  true,
		ProjectRoot: "/my/project",
		RuleCount:   map[string]int{"allow": 3, "deny": 1},
	}
	if !s.Discovered {
		t.Error("Discovered should be true")
	}
	if s.RuleCount["allow"] != 3 {
		t.Errorf("RuleCount allow = %d, want 3", s.RuleCount["allow"])
	}
}

func TestDiscoverPolicyFile_NotFound(t *testing.T) {
	dir := t.TempDir()
	path, err := discoverPolicyFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestDiscoverPolicyFile_Found(t *testing.T) {
	dir := t.TempDir()
	policyFile := filepath.Join(dir, "apm-policy.yml")
	if err := os.WriteFile(policyFile, []byte("allow:\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path, err := discoverPolicyFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != policyFile {
		t.Errorf("discoverPolicyFile = %q, want %q", path, policyFile)
	}
}

func TestCountRules(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "policy.yml")
	content := "allow:\ndeny:\n# comment\nrules:\n  - foo\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	counts, err := countRules(f)
	if err != nil {
		t.Fatalf("countRules: %v", err)
	}
	if counts["allow"] != 1 {
		t.Errorf("allow count = %d, want 1", counts["allow"])
	}
	if counts["deny"] != 1 {
		t.Errorf("deny count = %d, want 1", counts["deny"])
	}
}

func TestStatusOptionsFields(t *testing.T) {
	opts := StatusOptions{
		ProjectRoot: "/proj",
		Format:      "json",
		Verbose:     true,
		NoFetch:     false,
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want json", opts.Format)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestStripSourcePrefix(t *testing.T) {
	tests := []struct{ in, want string }{
		{"org:myorg", "myorg"},
		{"url:https://example.com", "https://example.com"},
		{"file:/tmp/policy.yml", "/tmp/policy.yml"},
		{"plain", "plain"},
		{"", ""},
	}
	for _, tc := range tests {
		got := stripSourcePrefix(tc.in)
		if got != tc.want {
			t.Errorf("stripSourcePrefix(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		secs int
		want string
	}{
		{-1, "n/a"},
		{0, "0s ago"},
		{30, "30s ago"},
		{59, "59s ago"},
		{60, "1m ago"},
		{3599, "59m ago"},
		{3600, "1h ago"},
		{86399, "23h ago"},
		{86400, "1d ago"},
	}
	for _, tc := range tests {
		got := formatAge(tc.secs)
		if got != tc.want {
			t.Errorf("formatAge(%d) = %q, want %q", tc.secs, got, tc.want)
		}
	}
}
