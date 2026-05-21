package githubhost_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

func TestIsGitHubHostname_GitHubCom(t *testing.T) {
	if !githubhost.IsGitHubHostname("github.com") {
		t.Error("expected github.com to be a GitHub hostname")
	}
}

func TestIsGitHubHostname_Empty(t *testing.T) {
	if githubhost.IsGitHubHostname("") {
		t.Error("expected empty string to not be a GitHub hostname")
	}
}

func TestIsGitHubHostname_Other(t *testing.T) {
	if githubhost.IsGitHubHostname("gitlab.com") {
		t.Error("expected gitlab.com to not be a GitHub hostname")
	}
}

func TestIsValidFQDN_Valid(t *testing.T) {
	cases := []string{"github.com", "api.github.com", "example.co.uk"}
	for _, c := range cases {
		if !githubhost.IsValidFQDN(c) {
			t.Errorf("expected %q to be a valid FQDN", c)
		}
	}
}

func TestIsValidFQDN_Invalid(t *testing.T) {
	cases := []string{"", "localhost", "nodot", "a..b.com"}
	for _, c := range cases {
		if githubhost.IsValidFQDN(c) {
			t.Errorf("expected %q to be invalid FQDN", c)
		}
	}
}

func TestParseHostFromURL_Standard(t *testing.T) {
	got := githubhost.ParseHostFromURL("https://github.com/owner/repo")
	if got != "github.com" {
		t.Errorf("got %q, want %q", got, "github.com")
	}
}

func TestParseHostFromURL_Empty_v4(t *testing.T) {
	got := githubhost.ParseHostFromURL("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestClassifyHost_GitHub_v4(t *testing.T) {
	got := githubhost.ClassifyHost("github.com")
	if got == "" {
		t.Error("expected non-empty classification for github.com")
	}
}

func TestIsSupportedGitHost_GitHub_v4(t *testing.T) {
	if !githubhost.IsSupportedGitHost("github.com") {
		t.Error("expected github.com to be a supported git host")
	}
}
