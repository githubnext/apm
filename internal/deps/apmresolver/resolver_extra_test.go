package apmresolver

import (
	"os"
	"testing"

	"github.com/githubnext/apm/internal/models/depreference"
)

func TestParseApmYMLDeps_InlineComment(t *testing.T) {
	content := `dependencies:
  - owner/repo # this is a comment
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RepoURL != "owner/repo" {
		t.Errorf("expected 'owner/repo', got %q", refs[0].RepoURL)
	}
}

func TestParseApmYMLDeps_SectionEndsAtNonIndented(t *testing.T) {
	content := `dependencies:
  - owner/repo
name: other
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 1 {
		t.Errorf("expected 1 ref (section stops at 'name:'), got %d", len(refs))
	}
}

func TestParseApmYMLDeps_DevDepsOnly(t *testing.T) {
	content := `devDependencies:
  - devowner/devrepo
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RepoURL != "devowner/devrepo" {
		t.Errorf("expected 'devowner/devrepo', got %q", refs[0].RepoURL)
	}
}

func TestParseApmYMLDeps_BothSections(t *testing.T) {
	content := `dependencies:
  - owner/repo-a
devDependencies:
  - owner/repo-b
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestParseApmYMLDeps_StripsQuotes(t *testing.T) {
	content := `dependencies:
  - "owner/quoted"
  - 'owner/single-quoted'
`
	refs := parseApmYMLDeps(content)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
	for _, ref := range refs {
		if ref.RepoURL == "" {
			t.Error("expected non-empty RepoURL after quote stripping")
		}
	}
}

func TestParseApmYMLDeps_EmptyLine(t *testing.T) {
	content := `dependencies:
  - owner/repo

  - owner/repo2
`
	refs := parseApmYMLDeps(content)
	// Empty lines within deps section are fine; both refs should be found
	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d", len(refs))
	}
}

func TestNew_WithApmModulesDir(t *testing.T) {
	r := New(Options{ApmModulesDir: "/custom/modules"})
	if r == nil {
		t.Fatal("New returned nil")
	}
}

func TestNew_WithDownloadFn(t *testing.T) {
	called := false
	r := New(Options{
		DownloadFn: func(ref *depreference.DependencyReference, apmModulesDir, parentChain, parentPkg string) string {
			called = true
			return ""
		},
	})
	if r == nil {
		t.Fatal("New returned nil")
	}
	_ = called
}

func TestResolveMaxParallel_EnvVar(t *testing.T) {
	orig := os.Getenv("APM_RESOLVE_PARALLEL")
	os.Setenv("APM_RESOLVE_PARALLEL", "7")
	defer os.Setenv("APM_RESOLVE_PARALLEL", orig)

	n := resolveMaxParallel(0)
	if n != 7 {
		t.Errorf("expected 7 from env var, got %d", n)
	}
}

func TestResolveMaxParallel_ExplicitOverridesEnv(t *testing.T) {
	orig := os.Getenv("APM_RESOLVE_PARALLEL")
	os.Setenv("APM_RESOLVE_PARALLEL", "7")
	defer os.Setenv("APM_RESOLVE_PARALLEL", orig)

	n := resolveMaxParallel(3)
	if n != 3 {
		t.Errorf("expected explicit 3 to override env, got %d", n)
	}
}

func TestResolveMaxParallel_InvalidEnv(t *testing.T) {
	orig := os.Getenv("APM_RESOLVE_PARALLEL")
	os.Setenv("APM_RESOLVE_PARALLEL", "not-a-number")
	defer os.Setenv("APM_RESOLVE_PARALLEL", orig)

	n := resolveMaxParallel(0)
	if n <= 0 {
		t.Errorf("expected positive default parallel, got %d", n)
	}
}

func TestNew_MaxParallel(t *testing.T) {
	r := New(Options{MaxParallel: 5})
	if r == nil {
		t.Fatal("New returned nil")
	}
}

func TestNew_DefaultMaxDepth(t *testing.T) {
	r := New(Options{MaxDepth: 0})
	if r.maxDepth != 50 {
		t.Errorf("expected default maxDepth=50, got %d", r.maxDepth)
	}
}

func TestNew_NegativeMaxDepth(t *testing.T) {
	r := New(Options{MaxDepth: -5})
	if r.maxDepth != 50 {
		t.Errorf("expected default maxDepth=50 for negative, got %d", r.maxDepth)
	}
}
