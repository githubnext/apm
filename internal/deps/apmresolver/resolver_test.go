package apmresolver

import (
	"testing"
)

func TestNew(t *testing.T) {
	r := New(Options{MaxDepth: 10})
	if r == nil {
		t.Fatal("New returned nil")
	}
	if r.maxDepth != 10 {
		t.Errorf("maxDepth = %d, want 10", r.maxDepth)
	}
}

func TestNewDefaults(t *testing.T) {
	r := New(Options{})
	if r.maxDepth != 50 {
		t.Errorf("default maxDepth = %d, want 50", r.maxDepth)
	}
}

func TestParseApmYMLDepsEmpty(t *testing.T) {
	refs := parseApmYMLDeps("")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}

func TestParseApmYMLDepsNoDeps(t *testing.T) {
	content := "name: my-package\nversion: 1.0.0\n"
	refs := parseApmYMLDeps(content)
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}

func TestParseApmYMLDepsSingleDep(t *testing.T) {
	content := `name: my-package
dependencies:
  - owner/repo
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RepoURL != "owner/repo" {
		t.Errorf("unexpected ref RepoURL: %s", refs[0].RepoURL)
	}
}

func TestParseApmYMLDepsMultiple(t *testing.T) {
	content := `name: pkg
dependencies:
  - owner1/repo1
  - owner2/repo2
  - owner3/repo3
other:
  - ignored
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 3 {
		t.Fatalf("expected 3 refs, got %d", len(refs))
	}
}

func TestParseApmYMLDepsWithComments(t *testing.T) {
	content := `dependencies:
  - owner/repo1 # this is a comment
  - owner/repo2
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestParseApmYMLDepsWithQuotes(t *testing.T) {
	content := `dependencies:
  - "owner/repo1"
  - 'owner/repo2'
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestParseApmYMLDepsDevDeps(t *testing.T) {
	content := `devDependencies:
  - owner/devrepo
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RepoURL != "owner/devrepo" {
		t.Errorf("unexpected RepoURL: %s", refs[0].RepoURL)
	}
}

func TestResolveMaxParallel(t *testing.T) {
	// With explicit value
	got := resolveMaxParallel(8)
	if got != 8 {
		t.Errorf("resolveMaxParallel(8) = %d, want 8", got)
	}
	// Zero falls back to default
	got = resolveMaxParallel(0)
	if got <= 0 {
		t.Errorf("resolveMaxParallel(0) = %d, want > 0", got)
	}
}
