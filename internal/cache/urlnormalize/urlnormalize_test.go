package urlnormalize_test

import (
	"testing"

	"github.com/githubnext/apm/internal/cache/urlnormalize"
)

func TestNormalizeRepoURL_StripsDotGit(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://github.com/owner/repo.git", "https://github.com/owner/repo"},
		{"https://github.com/owner/repo", "https://github.com/owner/repo"},
		{"ssh://git@github.com/owner/repo.git", "ssh://git@github.com/owner/repo"},
	}
	for _, tc := range tests {
		got := urlnormalize.NormalizeRepoURL(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeRepoURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestNormalizeRepoURL_LowercasesGitHubPath(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com/Owner/Repo")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_SCPLike(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("git@github.com:owner/repo.git")
	// SCP is converted to ssh://
	if got != "ssh://git@github.com/owner/repo" {
		t.Errorf("got %q", got)
	}
}

func TestNormalizeRepoURL_StripsDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://github.com:443/owner/repo")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_KeepsNonDefaultPort(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://example.com:8080/owner/repo")
	want := "https://example.com:8080/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_LowercasesHost(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://GITHUB.COM/owner/repo")
	want := "https://github.com/owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeRepoURL_StripsPassword(t *testing.T) {
	got := urlnormalize.NormalizeRepoURL("https://user:secret@example.com/org/repo")
	want := "https://user@example.com/org/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCacheKey_Length16(t *testing.T) {
	key := urlnormalize.CacheKey("https://github.com/owner/repo.git")
	if len(key) != 16 {
		t.Errorf("expected 16 char key, got %d: %q", len(key), key)
	}
}

func TestCacheKey_Deterministic(t *testing.T) {
	url := "https://github.com/owner/repo"
	k1 := urlnormalize.CacheKey(url)
	k2 := urlnormalize.CacheKey(url)
	if k1 != k2 {
		t.Errorf("non-deterministic: %q vs %q", k1, k2)
	}
}

func TestCacheKey_NormalizesBeforeHashing(t *testing.T) {
	k1 := urlnormalize.CacheKey("https://github.com/Owner/Repo.git")
	k2 := urlnormalize.CacheKey("https://github.com/owner/repo")
	if k1 != k2 {
		t.Errorf("normalization not applied: %q vs %q", k1, k2)
	}
}
