package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatusOptions_ZeroValue(t *testing.T) {
	opts := StatusOptions{}
	if opts.ProjectRoot != "" || opts.Format != "" {
		t.Errorf("unexpected non-empty fields: %+v", opts)
	}
	if opts.Verbose || opts.NoFetch {
		t.Error("expected bool fields false by default")
	}
}

func TestDebugOptions_Fields(t *testing.T) {
	opts := DebugOptions{
		ProjectRoot: "/proj",
		Format:      "json",
		Source:      "remote",
	}
	if opts.ProjectRoot != "/proj" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q", opts.Format)
	}
	if opts.Source != "remote" {
		t.Errorf("Source = %q", opts.Source)
	}
}

func TestPolicySource_AllFields(t *testing.T) {
	ps := PolicySource{
		Label:    "corp-policy",
		URL:      "https://example.com/policy.yaml",
		FilePath: "/etc/apm-policy.yaml",
		CacheAge: 120,
		Stale:    true,
	}
	if ps.Label != "corp-policy" {
		t.Errorf("Label = %q", ps.Label)
	}
	if ps.CacheAge != 120 {
		t.Errorf("CacheAge = %d", ps.CacheAge)
	}
	if !ps.Stale {
		t.Error("expected Stale=true")
	}
}

func TestPolicyStatus_AllFields(t *testing.T) {
	s := &PolicyStatus{
		Discovered:  true,
		Error:       "not found",
		ProjectRoot: "/myproject",
		CheckedAt:   "2025-06-01T10:00:00Z",
		RuleCount: map[string]int{
			"deny": 5,
			"warn": 2,
		},
	}
	if !s.Discovered {
		t.Error("expected Discovered=true")
	}
	if s.Error != "not found" {
		t.Errorf("Error = %q", s.Error)
	}
	if s.RuleCount["deny"] != 5 {
		t.Errorf("RuleCount[deny] = %d", s.RuleCount["deny"])
	}
}

func TestDiscoverPolicyFile_ApmPolicyYML(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "apm-policy.yml")
	if err := os.WriteFile(f, []byte("deny:\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	found, err := discoverPolicyFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != f {
		t.Errorf("found = %q, want %q", found, f)
	}
}

func TestDiscoverPolicyFile_NoPolicyFile(t *testing.T) {
	dir := t.TempDir()
	found, err := discoverPolicyFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Returns empty string (not an error) when no policy file exists.
	if found != "" {
		t.Errorf("expected empty path when no policy file, got %q", found)
	}
}

func TestStripSourcePrefix_WithPrefix(t *testing.T) {
	// Strips "org:" prefix
	got := stripSourcePrefix("org:github.com/myorg/policy")
	if got != "github.com/myorg/policy" {
		t.Errorf("expected prefix stripped, got %q", got)
	}
	// Strips "url:" prefix
	got2 := stripSourcePrefix("url:https://example.com/policy.yaml")
	if got2 != "https://example.com/policy.yaml" {
		t.Errorf("expected url: prefix stripped, got %q", got2)
	}
}

func TestStripSourcePrefix_NoPrefix(t *testing.T) {
	got := stripSourcePrefix("plain text")
	if got != "plain text" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestFormatAge_Zero(t *testing.T) {
	s := formatAge(0)
	if s == "" {
		t.Error("expected non-empty formatAge(0)")
	}
}

func TestFormatAge_Large(t *testing.T) {
	s := formatAge(7200)
	if s == "" {
		t.Error("expected non-empty formatAge(7200)")
	}
}

func TestRunStatus_EmptyProjectRoot(t *testing.T) {
	opts := StatusOptions{ProjectRoot: t.TempDir()}
	err := RunStatus(opts)
	// may fail due to missing policy, but must not panic
	_ = err
}
