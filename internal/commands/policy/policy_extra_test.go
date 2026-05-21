package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPolicySource_FetchError(t *testing.T) {
	ps := PolicySource{
		Label:      "remote",
		FetchError: "timeout",
	}
	if ps.FetchError != "timeout" {
		t.Errorf("FetchError = %q, want 'timeout'", ps.FetchError)
	}
}

func TestPolicySource_EmptyURL(t *testing.T) {
	ps := PolicySource{Label: "local"}
	if ps.URL != "" {
		t.Errorf("expected empty URL, got %q", ps.URL)
	}
}

func TestPolicyStatus_NoSource(t *testing.T) {
	s := &PolicyStatus{Discovered: false}
	if s.Source != nil {
		t.Error("expected nil Source when not discovered")
	}
}

func TestPolicyStatus_InheritanceChain(t *testing.T) {
	s := &PolicyStatus{
		InheritanceChain: []PolicySource{
			{Label: "parent"},
			{Label: "child"},
		},
	}
	if len(s.InheritanceChain) != 2 {
		t.Errorf("expected 2 chain entries, got %d", len(s.InheritanceChain))
	}
	if s.InheritanceChain[0].Label != "parent" {
		t.Errorf("first entry should be 'parent', got %q", s.InheritanceChain[0].Label)
	}
}

func TestPolicyStatus_CheckedAt(t *testing.T) {
	s := &PolicyStatus{CheckedAt: "2025-01-01T00:00:00Z"}
	if s.CheckedAt != "2025-01-01T00:00:00Z" {
		t.Errorf("CheckedAt = %q", s.CheckedAt)
	}
}

func TestDiscoverPolicyFile_ApmPolicyYAML(t *testing.T) {
	dir := t.TempDir()
	yamlFile := filepath.Join(dir, "apm-policy.yaml")
	if err := os.WriteFile(yamlFile, []byte("deny:\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path, err := discoverPolicyFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Error("expected to find apm-policy.yaml")
	}
}

func TestDiscoverPolicyFile_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	path, err := discoverPolicyFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path for dir without policy, got %q", path)
	}
}

func TestCountRules_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "empty.yml")
	if err := os.WriteFile(f, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	counts, err := countRules(f)
	if err != nil {
		t.Fatalf("countRules: %v", err)
	}
	if len(counts) != 0 {
		t.Errorf("expected empty counts, got %v", counts)
	}
}

func TestCountRules_MultipleKeys(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "policy.yml")
	content := "allow:\n  - github.com/a\ndeny:\n  - github.com/b\n  - github.com/c\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	counts, err := countRules(f)
	if err != nil {
		t.Fatalf("countRules: %v", err)
	}
	if counts["allow"] < 1 {
		t.Errorf("expected allow >= 1, got %d", counts["allow"])
	}
}

func TestFormatAge_BoundaryValues(t *testing.T) {
	// Edge cases around minute and hour boundaries
	if got := formatAge(59); got != "59s ago" {
		t.Errorf("59s: got %q", got)
	}
	if got := formatAge(60); got != "1m ago" {
		t.Errorf("60s: got %q", got)
	}
	if got := formatAge(3600); got != "1h ago" {
		t.Errorf("3600s: got %q", got)
	}
	if got := formatAge(86400); got != "1d ago" {
		t.Errorf("86400s: got %q", got)
	}
}

func TestStatusOptions_Defaults(t *testing.T) {
	opts := StatusOptions{}
	if opts.ProjectRoot != "" {
		t.Errorf("expected empty ProjectRoot, got %q", opts.ProjectRoot)
	}
	if opts.Format != "" {
		t.Errorf("expected empty Format, got %q", opts.Format)
	}
	if opts.Verbose {
		t.Error("Verbose should default to false")
	}
	if opts.NoFetch {
		t.Error("NoFetch should default to false")
	}
}
