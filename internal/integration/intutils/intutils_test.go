package intutils_test

import (
	"testing"

	"github.com/githubnext/apm/internal/integration/intutils"
)

func TestNormalizeRepoURL_Plain(t *testing.T) {
	// No scheme: trim trailing slash and .git suffix
	got := intutils.NormalizeRepoURL("owner/repo")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_TrailingSlash(t *testing.T) {
	got := intutils.NormalizeRepoURL("owner/repo/")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_GitSuffix(t *testing.T) {
	got := intutils.NormalizeRepoURL("owner/repo.git")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_HTTPS(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_HTTPSWithGit(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo.git")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_HTTPSWithTrailingSlash(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo/")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_SSH(t *testing.T) {
	got := intutils.NormalizeRepoURL("git://github.com/owner/repo.git")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}
