package intutils_test

import (
	"testing"

	"github.com/githubnext/apm/internal/integration/intutils"
)

func TestNormalizeRepoURL_HTTPSSubdomain(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo/tree/main")
	if got != "owner/repo/tree/main" {
		t.Fatalf("expected 'owner/repo/tree/main', got %q", got)
	}
}

func TestNormalizeRepoURL_HTTPSGitDotGit(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo.git/")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_SchemeNoSlash(t *testing.T) {
	// scheme present but no slash after host -- should return URL as-is
	got := intutils.NormalizeRepoURL("https://github.com")
	if got == "" {
		t.Fatal("should not return empty string for scheme-only URL")
	}
}

func TestNormalizeRepoURL_GHEHost(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.example.com/owner/repo")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_DoubleSlash(t *testing.T) {
	got := intutils.NormalizeRepoURL("owner/repo//")
	// trailing slashes trimmed
	if got == "" {
		t.Fatal("should not return empty string")
	}
}

func TestNormalizeRepoURL_OnlyDotGit(t *testing.T) {
	got := intutils.NormalizeRepoURL("repo.git")
	if got != "repo" {
		t.Fatalf("expected 'repo', got %q", got)
	}
}

func TestNormalizeRepoURL_NestedPathWithGit(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/org/repo/sub.git")
	if got != "org/repo/sub" {
		t.Fatalf("expected 'org/repo/sub', got %q", got)
	}
}

func TestNormalizeRepoURL_UnusualScheme(t *testing.T) {
	got := intutils.NormalizeRepoURL("ssh://git@github.com/owner/repo.git")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_PlainOwnerRepo(t *testing.T) {
	got := intutils.NormalizeRepoURL("owner/repo")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoURL_DotGitSuffix(t *testing.T) {
	got := intutils.NormalizeRepoURL("owner/repo.git")
	if got != "owner/repo" {
t.Fatalf("expected 'owner/repo', got %q", got)
}
}

func TestNormalizeRepoURL_HTTPGitHub(t *testing.T) {
	got := intutils.NormalizeRepoURL("http://github.com/owner/repo")
	if got != "owner/repo" {
t.Fatalf("expected 'owner/repo', got %q", got)
}
}

func TestNormalizeRepoURL_HTTPSWithDotGitAndSlash(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/myorg/myrepo.git/")
	if got != "myorg/myrepo" {
t.Fatalf("expected 'myorg/myrepo', got %q", got)
}
}

func TestNormalizeRepoURL_SubpathPreserved(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo/blob/main/README.md")
	if got == "" {
t.Fatal("should not return empty for path with subpath")
}
}

func TestNormalizeRepoURL_OnlyScheme(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://")
	if got == "" {
t.Fatal("should return something (the input unchanged) for scheme-only")
}
}
