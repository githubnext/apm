package urlnormalize_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/urlnormalize"
)

func TestNormalizeRepoURL_HTTPSNonDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com:8443/owner/repo")
	if !strings.Contains(got, "8443") {
		t.Errorf("non-default port should be preserved, got %q", got)
	}
}

func TestNormalizeRepoURL_SSHNonDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("ssh://git@github.com:2222/owner/repo")
	if !strings.Contains(got, "2222") {
		t.Errorf("non-default SSH port should be preserved, got %q", got)
	}
}

func TestNormalizeRepoURL_StripsUserPassword(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://user:secret@github.com/owner/repo")
	if strings.Contains(got, "secret") {
		t.Errorf("password should be stripped, got %q", got)
	}
	if !strings.Contains(got, "user") {
		t.Errorf("username should be preserved, got %q", got)
	}
}

func TestNormalizeRepoURL_GitHubHostLowercased(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://GitHub.COM/owner/repo")
	if !strings.HasPrefix(got, "https://github.com/") {
		t.Errorf("host should be lowercased, got %q", got)
	}
}

func TestNormalizeRepoURL_MultipleCallsIdempotent(t *testing.T) {
	input := "https://github.com/Owner/Repo.git"
	r1 := urlnormalize.NormalizeRepoURL(input)
	r2 := urlnormalize.NormalizeRepoURL(r1)
	if r1 != r2 {
		t.Errorf("normalization should be idempotent: %q vs %q", r1, r2)
	}
}

func TestNormalizeRepoURL_StripsDotGitAndLowercases(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com/OWNER/REPO.git")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCacheKey_SameNormalizedURLsSameKey(t *testing.T) {
	k1 := urlnormalize.CacheKey("https://github.com/owner/repo.git")
	k2 := urlnormalize.CacheKey("https://github.com/owner/repo")
	if k1 != k2 {
		t.Errorf("normalized URLs should produce same cache key: %q vs %q", k1, k2)
	}
}

func TestCacheKey_Exactly16Chars(t *testing.T) {
	k := urlnormalize.CacheKey("https://github.com/a/b")
	if len(k) != 16 {
		t.Errorf("expected 16-char cache key, got %d chars: %q", len(k), k)
	}
}

func TestCacheKey_OnlyHexChars(t *testing.T) {
	k := urlnormalize.CacheKey("https://github.com/a/b")
	for _, c := range k {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("cache key contains non-hex char %q: %q", string(c), k)
		}
	}
}

func TestNormalizeRepoURL_SCPSyntaxToSSH(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("git@github.com:owner/repo")
	if !strings.HasPrefix(got, "ssh://git@github.com/") {
		t.Errorf("SCP syntax should convert to ssh://, got %q", got)
	}
}

func TestNormalizeRepoURL_SCPWithUppercaseHost(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("git@GitHub.COM:owner/repo")
	if !strings.Contains(got, "github.com") {
		t.Errorf("SCP host should be lowercased, got %q", got)
	}
}

func TestNormalizeRepoURL_GitHubPathAlwaysLowercase(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com/MyOrg/MyRepo")
	if strings.Contains(got, "MyOrg") || strings.Contains(got, "MyRepo") {
		t.Errorf("github.com paths should be lowercased, got %q", got)
	}
}
