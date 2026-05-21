package urlnormalize_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/urlnormalize"
)

func TestNormalizeRepoURL_TrailingDotGit_Stripped_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com/org/repo.git")
	if strings.HasSuffix(got, ".git") {
		t.Errorf("expected .git stripped, got %q", got)
	}
}

func TestNormalizeRepoURL_NoScheme_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("github.com/org/repo")
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestNormalizeRepoURL_HTTPSSchemePreserved_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com/org/repo")
	if !strings.HasPrefix(got, "https://") {
		t.Errorf("expected https:// prefix, got %q", got)
	}
}

func TestNormalizeRepoURL_SSHSchemePreserved_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("ssh://git@github.com/org/repo")
	if !strings.HasPrefix(got, "ssh://") {
		t.Errorf("expected ssh:// prefix, got %q", got)
	}
}

func TestNormalizeRepoURL_GitLabPathLowercased_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://gitlab.com/Org/Repo")
	if strings.Contains(got, "Org") || strings.Contains(got, "Repo") {
		t.Errorf("expected path lowercased, got %q", got)
	}
}

func TestNormalizeRepoURL_BitbucketPathLowercased_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://bitbucket.org/Org/Repo")
	if strings.Contains(got, "Org") || strings.Contains(got, "Repo") {
		t.Errorf("expected path lowercased, got %q", got)
	}
}

func TestNormalizeRepoURL_NonGitHubPathCasePreserved_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://example.com/MyOrg/MyRepo")
	if !strings.Contains(got, "MyOrg") {
		t.Errorf("expected case preserved for non-github host, got %q", got)
	}
}

func TestCacheKey_DifferentURLsDifferentKeys_Extra4(t *testing.T) {
	k1 := urlnormalize.CacheKey("https://github.com/org/repo1")
	k2 := urlnormalize.CacheKey("https://github.com/org/repo2")
	if k1 == k2 {
		t.Error("expected different cache keys for different repos")
	}
}

func TestCacheKey_DotGitAndWithoutEqual_Extra4(t *testing.T) {
	k1 := urlnormalize.CacheKey("https://github.com/org/repo.git")
	k2 := urlnormalize.CacheKey("https://github.com/org/repo")
	if k1 != k2 {
		t.Errorf("expected same key for .git vs no .git: %q vs %q", k1, k2)
	}
}

func TestNormalizeRepoURL_EmptyString_Extra4(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("")
	_ = got // no panic expected
}

func TestCacheKey_LengthIs16_Extra4(t *testing.T) {
	k := urlnormalize.CacheKey("https://github.com/test/test")
	if len(k) != 16 {
		t.Errorf("expected 16 chars, got %d: %q", len(k), k)
	}
}
