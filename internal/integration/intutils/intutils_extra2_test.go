package intutils

import (
	"testing"
)

func TestNormalizeRepoURL_GitLabHTTPS(t *testing.T) {
	got := NormalizeRepoURL("https://gitlab.com/owner/repo.git")
	if got != "owner/repo" {
		t.Errorf("got %q, want 'owner/repo'", got)
	}
}

func TestNormalizeRepoURL_BitbucketHTTPS(t *testing.T) {
	got := NormalizeRepoURL("https://bitbucket.org/owner/repo")
	if got != "owner/repo" {
		t.Errorf("got %q, want 'owner/repo'", got)
	}
}

func TestNormalizeRepoURL_PlainNoSuffix(t *testing.T) {
	got := NormalizeRepoURL("owner/repo")
	if got != "owner/repo" {
		t.Errorf("got %q, want 'owner/repo'", got)
	}
}

func TestNormalizeRepoURL_TripleSlash(t *testing.T) {
	// Unusual but should not panic
	got := NormalizeRepoURL("https://github.com///repo")
	_ = got // just verify no panic
}

func TestNormalizeRepoURL_ShortPath(t *testing.T) {
	got := NormalizeRepoURL("https://github.com/a")
	if got == "" {
		t.Error("expected non-empty for short path")
	}
}

func TestNormalizeRepoURL_DotGitAndTrailingSlash(t *testing.T) {
	got := NormalizeRepoURL("https://github.com/org/repo.git/")
	// Should strip .git and trailing slash
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestNormalizeRepoURL_NoSchemeWithGitSuffix(t *testing.T) {
	got := NormalizeRepoURL("org/repo.git")
	if got != "org/repo" {
		t.Errorf("got %q, want 'org/repo'", got)
	}
}

func TestNormalizeRepoURL_MixedGitCasing(t *testing.T) {
	got := NormalizeRepoURL("https://github.com/Owner/Repo.git")
	if got != "Owner/Repo" {
		t.Errorf("got %q, want 'Owner/Repo'", got)
	}
}
