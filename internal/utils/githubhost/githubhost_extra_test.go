package githubhost_test

import (
	"os"
	"testing"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

func TestIsVisualStudioLegacyHostname_Variants(t *testing.T) {
	tests := []struct {
		h    string
		want bool
	}{
		{"myorg.visualstudio.com", true},
		{"MYORG.VISUALSTUDIO.COM", true},
		{"github.com", false},
		{"dev.azure.com", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := githubhost.IsVisualStudioLegacyHostname(tt.h); got != tt.want {
			t.Errorf("IsVisualStudioLegacyHostname(%q) = %v, want %v", tt.h, got, tt.want)
		}
	}
}

func TestIsGitLabHostname_DefaultFalse(t *testing.T) {
	os.Unsetenv("GITLAB_HOST")
	os.Unsetenv("APM_GITLAB_HOSTS")
	os.Unsetenv("GITHUB_HOST")
	if githubhost.IsGitLabHostname("github.com") {
		t.Error("github.com should not be gitlab")
	}
	if githubhost.IsGitLabHostname("") {
		t.Error("empty should not be gitlab")
	}
}

func TestIsGitLabHostname_GitLabSaaS(t *testing.T) {
	os.Unsetenv("GITLAB_HOST")
	os.Unsetenv("APM_GITLAB_HOSTS")
	os.Unsetenv("GITHUB_HOST")
	if !githubhost.IsGitLabHostname("gitlab.com") {
		t.Error("gitlab.com should be recognized as gitlab")
	}
}

func TestIsGitLabHostname_ViaEnv(t *testing.T) {
	os.Setenv("GITLAB_HOST", "mygitlab.example.com")
	defer os.Unsetenv("GITLAB_HOST")
	if !githubhost.IsGitLabHostname("mygitlab.example.com") {
		t.Error("host matching GITLAB_HOST should be recognized as gitlab")
	}
}

func TestIsSupportedGitHost_GitHub(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	if !githubhost.IsSupportedGitHost("github.com") {
		t.Error("github.com should be a supported git host")
	}
}

func TestIsSupportedGitHost_ADO(t *testing.T) {
	if !githubhost.IsSupportedGitHost("dev.azure.com") {
		t.Error("dev.azure.com should be supported")
	}
}

func TestIsSupportedGitHost_Unknown(t *testing.T) {
	// IsSupportedGitHost returns true for any valid FQDN
	if !githubhost.IsSupportedGitHost("example.com") {
		t.Error("valid FQDN should be supported")
	}
	if githubhost.IsSupportedGitHost("") {
		t.Error("empty host should not be supported")
	}
	if githubhost.IsSupportedGitHost("localhost") {
		t.Error("localhost (no TLD) should not be supported")
	}
}

func TestIsArtifactoryHostname_False(t *testing.T) {
	if githubhost.IsArtifactoryHostname("github.com") {
		t.Error("github.com should not be artifactory")
	}
}

func TestParseHostFromURL_HTTPS(t *testing.T) {
	got := githubhost.ParseHostFromURL("https://github.com/owner/repo.git")
	if got != "github.com" {
		t.Errorf("ParseHostFromURL HTTPS: got %q, want github.com", got)
	}
}

func TestParseHostFromURL_SSH(t *testing.T) {
	got := githubhost.ParseHostFromURL("git@github.com:owner/repo.git")
	if got != "github.com" {
		t.Errorf("ParseHostFromURL SSH: got %q, want github.com", got)
	}
}

func TestParseHostFromURL_Empty(t *testing.T) {
	got := githubhost.ParseHostFromURL("")
	if got != "" {
		t.Errorf("ParseHostFromURL empty: got %q, want empty", got)
	}
}

func TestAzureDevOpsOrgFromHostname_VisualStudio(t *testing.T) {
	got := githubhost.AzureDevOpsOrgFromHostname("myorg.visualstudio.com")
	if got != "myorg" {
		t.Errorf("got %q, want myorg", got)
	}
}

func TestAzureDevOpsOrgFromHostname_NonADO(t *testing.T) {
	got := githubhost.AzureDevOpsOrgFromHostname("github.com")
	if got != "" {
		t.Errorf("non-ADO host should return empty org, got %q", got)
	}
}

func TestClassifyHost_GitHub(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	got := githubhost.ClassifyHost("github.com")
	if got != "github" {
		t.Errorf("ClassifyHost(github.com) = %q, want github", got)
	}
}

func TestClassifyHost_ADO(t *testing.T) {
	got := githubhost.ClassifyHost("dev.azure.com")
	if got != "azure_devops" {
		t.Errorf("ClassifyHost(dev.azure.com) = %q, want azure_devops", got)
	}
}

func TestHasGitHubGitLabHostEnvConflict_NoConflict(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	os.Unsetenv("GITLAB_HOST")
	if githubhost.HasGitHubGitLabHostEnvConflict("github.com") {
		t.Error("no conflict expected for github.com with no env vars")
	}
}
