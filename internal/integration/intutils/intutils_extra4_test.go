package intutils_test

import (
	"testing"

	"github.com/githubnext/apm/internal/integration/intutils"
)

func TestNormalizeRepoURL_HTTPPathOnly_v4(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo")
	if got != "owner/repo" {
		t.Errorf("got %q, want %q", got, "owner/repo")
	}
}

func TestNormalizeRepoURL_WithDotGitHTTPS_v4(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/owner/repo.git")
	if got != "owner/repo" {
		t.Errorf("got %q, want %q", got, "owner/repo")
	}
}

func TestNormalizeRepoURL_EmptyString_v4(t *testing.T) {
	got := intutils.NormalizeRepoURL("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestNormalizeRepoURL_SlashTrailingPlain_v4(t *testing.T) {
	got := intutils.NormalizeRepoURL("owner/repo/")
	if got != "owner/repo" {
		t.Errorf("got %q, want %q", got, "owner/repo")
	}
}

func TestNormalizeRepoURL_DeepPath_v4(t *testing.T) {
	got := intutils.NormalizeRepoURL("https://github.com/org/repo/tree/main")
	// Should return path after the host
	if got == "" {
		t.Error("expected non-empty result for deep path")
	}
}

func TestNormalizeRepoURL_SSHLike_v4(t *testing.T) {
	// No scheme, treated as plain path
	got := intutils.NormalizeRepoURL("git@github.com:owner/repo.git")
	if got == "" {
		t.Error("expected non-empty result for SSH-like URL")
	}
}
