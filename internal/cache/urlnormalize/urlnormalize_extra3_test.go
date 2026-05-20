package urlnormalize_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/urlnormalize"
)

func TestNormalizeRepoURL_TrailingGitStripped_Extra3(t *testing.T) {
	a := urlnormalize.NormalizeRepoURL("https://github.com/owner/repo.git")
	b := urlnormalize.NormalizeRepoURL("https://github.com/owner/repo")
	if a != b {
		t.Errorf("trailing .git should be stripped: %q vs %q", a, b)
	}
}

func TestNormalizeRepoURL_HostLowercased_Extra3(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://GITHUB.COM/owner/repo")
	if strings.Contains(got, "GITHUB.COM") {
		t.Errorf("host should be lowercased, got %q", got)
	}
}

func TestNormalizeRepoURL_GitHubPathLowercased_Extra3(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com/Owner/Repo")
	if strings.Contains(got, "/Owner/") || strings.Contains(got, "/Repo") {
		t.Errorf("github.com path should be lowercased, got %q", got)
	}
}

func TestNormalizeRepoURL_PasswordStripped_Extra3(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://user:secret@github.com/owner/repo")
	if strings.Contains(got, "secret") {
		t.Errorf("password should be stripped from normalized URL, got %q", got)
	}
}

func TestNormalizeRepoURL_UserRetained_Extra3(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://user:secret@github.com/owner/repo")
	if !strings.Contains(got, "user@") {
		t.Errorf("username should be retained, got %q", got)
	}
}

func TestNormalizeRepoURL_DefaultHTTPSPortStripped_Extra3(t *testing.T) {
	a := urlnormalize.NormalizeRepoURL("https://github.com:443/owner/repo")
	b := urlnormalize.NormalizeRepoURL("https://github.com/owner/repo")
	if a != b {
		t.Errorf("default HTTPS port 443 should be stripped: %q vs %q", a, b)
	}
}

func TestNormalizeRepoURL_DefaultSSHPortStripped_Extra3(t *testing.T) {
	a := urlnormalize.NormalizeRepoURL("ssh://git@github.com:22/owner/repo")
	b := urlnormalize.NormalizeRepoURL("ssh://git@github.com/owner/repo")
	if a != b {
		t.Errorf("default SSH port 22 should be stripped: %q vs %q", a, b)
	}
}

func TestNormalizeRepoURL_SCPLikeConversion_Extra3(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("git@github.com:owner/repo")
	if !strings.HasPrefix(got, "ssh://") {
		t.Errorf("SCP-like URL should be converted to ssh://, got %q", got)
	}
}

func TestCacheKey_Length16_Extra3(t *testing.T) {
	key := urlnormalize.CacheKey("https://github.com/owner/repo")
	if len(key) != 16 {
		t.Errorf("CacheKey should be 16 hex chars, got %d: %q", len(key), key)
	}
}

func TestCacheKey_Deterministic_Extra3(t *testing.T) {
	a := urlnormalize.CacheKey("https://github.com/owner/repo")
	b := urlnormalize.CacheKey("https://github.com/owner/repo")
	if a != b {
		t.Errorf("CacheKey should be deterministic: %q vs %q", a, b)
	}
}

func TestCacheKey_DifferentURLsDifferentKeys_Extra3(t *testing.T) {
	a := urlnormalize.CacheKey("https://github.com/owner/repo1")
	b := urlnormalize.CacheKey("https://github.com/owner/repo2")
	if a == b {
		t.Error("different URLs should produce different cache keys")
	}
}
