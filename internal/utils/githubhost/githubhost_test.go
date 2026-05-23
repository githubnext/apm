// Package githubhost_test provides parity tests for githubhost utilities.
// Mirrors src/apm_cli/utils/github_host.py behaviour.
package githubhost_test

import (
	"os"
	"testing"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

func TestParityIsValidFQDN(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"github.com", true},
		{"dev.azure.com", true},
		{"my.custom.host.example.com", true},
		{"", false},
		{"localhost", false},
		{"invalid-", false},
		{"-invalid.com", false},
		{"has space.com", false},
	}
	for _, c := range cases {
		got := githubhost.IsValidFQDN(c.in)
		if got != c.want {
			t.Errorf("IsValidFQDN(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParityDefaultHost(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	if h := githubhost.DefaultHost(); h != "github.com" {
		t.Errorf("DefaultHost() = %q, want github.com", h)
	}
	t.Setenv("GITHUB_HOST", "myghe.com")
	if h := githubhost.DefaultHost(); h != "myghe.com" {
		t.Errorf("DefaultHost() with GITHUB_HOST = %q, want myghe.com", h)
	}
}

func TestParityIsAzureDevOpsHostname(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"dev.azure.com", true},
		{"myorg.visualstudio.com", true},
		{"github.com", false},
		{"", false},
	}
	for _, c := range cases {
		got := githubhost.IsAzureDevOpsHostname(c.in)
		if got != c.want {
			t.Errorf("IsAzureDevOpsHostname(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParityIsGitHubHostname(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"github.com", true},
		{"myenterprise.ghe.com", true},
		{"gitlab.com", false},
		{"dev.azure.com", false},
		{"", false},
	}
	for _, c := range cases {
		got := githubhost.IsGitHubHostname(c.in)
		if got != c.want {
			t.Errorf("IsGitHubHostname(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParityIsGitLabHostname(t *testing.T) {
	if got := githubhost.IsGitLabHostname("gitlab.com"); !got {
		t.Error("IsGitLabHostname(gitlab.com) should be true")
	}
	if got := githubhost.IsGitLabHostname("github.com"); got {
		t.Error("IsGitLabHostname(github.com) should be false")
	}
	if got := githubhost.IsGitLabHostname(""); got {
		t.Error("IsGitLabHostname('') should be false")
	}
	t.Setenv("GITLAB_HOST", "mygitlab.example.com")
	if got := githubhost.IsGitLabHostname("mygitlab.example.com"); !got {
		t.Error("IsGitLabHostname with GITLAB_HOST should be true")
	}
}

func TestParityIsADOAuthFailureSignal(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"HTTP 401 Unauthorized", true},
		{"authentication failed", true},
		{"could not read username", true},
		{"403 Forbidden", true},
		{"Unauthorized access", true},
		{"", false},
		{"everything is fine", false},
	}
	for _, c := range cases {
		got := githubhost.IsADOAuthFailureSignal(c.in)
		if got != c.want {
			t.Errorf("IsADOAuthFailureSignal(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParityBuildADOBearerGitEnv(t *testing.T) {
	env := githubhost.BuildADOBearerGitEnv("mytoken")
	if env["GIT_CONFIG_COUNT"] != "1" {
		t.Errorf("GIT_CONFIG_COUNT = %q", env["GIT_CONFIG_COUNT"])
	}
	if env["GIT_CONFIG_KEY_0"] != "http.extraheader" {
		t.Errorf("GIT_CONFIG_KEY_0 = %q", env["GIT_CONFIG_KEY_0"])
	}
	want := "Authorization: Bearer mytoken"
	if env["GIT_CONFIG_VALUE_0"] != want {
		t.Errorf("GIT_CONFIG_VALUE_0 = %q, want %q", env["GIT_CONFIG_VALUE_0"], want)
	}
}
