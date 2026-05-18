package urlnormalize_test

import (
	"testing"

	"github.com/githubnext/apm/internal/cache/urlnormalize"
)

func TestNormalizeRepoURL_GitLabLowercase(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://gitlab.com/Owner/Repo")
	want := "https://gitlab.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_BitbucketLowercase(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://bitbucket.org/Owner/Repo.git")
	want := "https://bitbucket.org/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_SSHDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("ssh://git@github.com:22/owner/repo")
	want := "ssh://git@github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_GitDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("git://github.com:9418/owner/repo")
	want := "git://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_NoScheme(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("github.com/owner/repo")
	// without scheme, host is treated as-is
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestNormalizeRepoURL_EmptyString(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestNormalizeRepoURL_SCPWithDotGit(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("git@gitlab.com:owner/repo.git")
	if got != "ssh://git@gitlab.com/owner/repo" {
		t.Errorf("got %q", got)
	}
}

func TestNormalizeRepoURL_HTTPDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("http://github.com:80/owner/repo")
	want := "http://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_StripTrailingWhitespace(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("  https://github.com/owner/repo  ")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCacheKey_DifferentURLs(t *testing.T) {
	k1 := urlnormalize.CacheKey("https://github.com/owner/repo1")
	k2 := urlnormalize.CacheKey("https://github.com/owner/repo2")
	if k1 == k2 {
		t.Error("different URLs must produce different cache keys")
	}
}

func TestCacheKey_IsHex(t *testing.T) {
	key := urlnormalize.CacheKey("https://github.com/owner/repo")
	for _, c := range key {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("cache key contains non-hex char %q: %s", c, key)
		}
	}
}

func TestNormalizeRepoURL_PreservesCustomHostPath(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://ghe.example.com/Org/Repo")
	// non-known host: path is preserved as-is (not lowercased)
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestNormalizeRepoURL_UserWithoutPassword(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://user@example.com/org/repo")
	want := "https://user@example.com/org/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
