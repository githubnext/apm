package apmresolver

import (
	"testing"
)

func TestParseApmYMLDeps_MultipleDevDeps(t *testing.T) {
	content := "devDependencies:\n  - owner/pkg1\n  - owner/pkg2\n"
	deps := parseApmYMLDeps(content)
	if len(deps) != 2 {
		t.Errorf("expected 2 dev deps, got %d", len(deps))
	}
}

func TestParseApmYMLDeps_MixedSection(t *testing.T) {
	content := "dependencies:\n  - owner/a\ndependencies:\n  - owner/b\n"
	deps := parseApmYMLDeps(content)
	if len(deps) < 1 {
		t.Error("expected at least 1 dep from mixed content")
	}
}

func TestParseApmYMLDeps_EmptyContent(t *testing.T) {
	deps := parseApmYMLDeps("")
	if deps != nil && len(deps) != 0 {
		t.Errorf("expected no deps from empty content, got %d", len(deps))
	}
}

func TestNew_ZeroMaxDepth(t *testing.T) {
	r := New(Options{MaxDepth: 0})
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNew_PositiveMaxDepth(t *testing.T) {
	r := New(Options{MaxDepth: 10})
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNew_NilDownloadFn(t *testing.T) {
	r := New(Options{DownloadFn: nil})
	if r == nil {
		t.Fatal("expected non-nil resolver even with nil DownloadFn")
	}
}

func TestParseApmYMLDeps_WithGitRef(t *testing.T) {
	content := "dependencies:\n  - owner/repo-ref\n"
	deps := parseApmYMLDeps(content)
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	if deps[0].RepoURL == "" {
		t.Error("expected non-empty RepoURL")
	}
}

func TestParseApmYMLDeps_WithVersion(t *testing.T) {
	content := "dependencies:\n  - owner/repo-ver\n"
	deps := parseApmYMLDeps(content)
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
}

func TestResolveMaxParallel_Zero(t *testing.T) {
	p := resolveMaxParallel(0)
	if p <= 0 {
		t.Errorf("expected positive parallelism, got %d", p)
	}
}
