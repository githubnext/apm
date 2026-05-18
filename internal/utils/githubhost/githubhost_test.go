package githubhost_test

import (
	"os"
	"testing"

	"github.com/githubnext/apm/internal/utils/githubhost"
)

func TestDefaultHost(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	if got := githubhost.DefaultHost(); got != "github.com" {
		t.Errorf("want github.com got %s", got)
	}
	os.Setenv("GITHUB_HOST", "myghe.example.com")
	if got := githubhost.DefaultHost(); got != "myghe.example.com" {
		t.Errorf("want myghe.example.com got %s", got)
	}
	os.Unsetenv("GITHUB_HOST")
}

func TestIsAzureDevOpsHostname(t *testing.T) {
	tests := []struct{ h string; want bool }{
		{"dev.azure.com", true},
		{"myorg.visualstudio.com", true},
		{"github.com", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := githubhost.IsAzureDevOpsHostname(tt.h); got != tt.want {
			t.Errorf("IsAzureDevOpsHostname(%q)=%v want %v", tt.h, got, tt.want)
		}
	}
}

func TestIsValidFQDN(t *testing.T) {
	tests := []struct{ h string; want bool }{
		{"github.com", true},
		{"myghe.example.com", true},
		{"localhost", false},
		{"", false},
		{"not valid!", false},
	}
	for _, tt := range tests {
		if got := githubhost.IsValidFQDN(tt.h); got != tt.want {
			t.Errorf("IsValidFQDN(%q)=%v want %v", tt.h, got, tt.want)
		}
	}
}

func TestClassifyHost(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	os.Unsetenv("GITLAB_HOST")
	os.Unsetenv("APM_GITLAB_HOSTS")

	if got := githubhost.ClassifyHost("github.com"); got != "github" {
		t.Errorf("want github got %s", got)
	}
	if got := githubhost.ClassifyHost("myorg.ghe.com"); got != "ghe_com" {
		t.Errorf("want ghe_com got %s", got)
	}
	if got := githubhost.ClassifyHost("dev.azure.com"); got != "azure_devops" {
		t.Errorf("want azure_devops got %s", got)
	}

	os.Setenv("GITLAB_HOST", "gitlab.mycompany.com")
	if got := githubhost.ClassifyHost("gitlab.mycompany.com"); got != "gitlab" {
		t.Errorf("want gitlab got %s", got)
	}
	os.Unsetenv("GITLAB_HOST")
}

func TestIsGHEHostname(t *testing.T) {
os.Unsetenv("GITHUB_HOST")
tests := []struct {
h    string
want bool
}{
{"myorg.ghe.com", true},
{"github.com", false},
{"", false},
{"dev.azure.com", false},
}
for _, tt := range tests {
if got := githubhost.IsGHEHostname(tt.h); got != tt.want {
t.Errorf("IsGHEHostname(%q)=%v want %v", tt.h, got, tt.want)
}
}
}

func TestIsGitHubHostname(t *testing.T) {
os.Unsetenv("GITHUB_HOST")
if !githubhost.IsGitHubHostname("github.com") {
t.Error("github.com should be a GitHub hostname")
}
if githubhost.IsGitHubHostname("") {
t.Error("empty string should not be a GitHub hostname")
}
if githubhost.IsGitHubHostname("dev.azure.com") {
t.Error("azure devops should not be a GitHub hostname")
}
}

func TestAzureDevOpsOrgFromHostname(t *testing.T) {
tests := []struct {
h    string
want string
}{
{"myorg.visualstudio.com", "myorg"},
{"ACME.visualstudio.com", "acme"},
{"github.com", ""},
{"dev.azure.com", ""},
{"", ""},
}
for _, tt := range tests {
got := githubhost.AzureDevOpsOrgFromHostname(tt.h)
if got != tt.want {
t.Errorf("AzureDevOpsOrgFromHostname(%q)=%q want %q", tt.h, got, tt.want)
}
}
}

func TestParseHostFromURL(t *testing.T) {
tests := []struct {
url  string
want string
}{
{"https://github.com/owner/repo", "github.com"},
{"http://dev.azure.com/org/proj", "dev.azure.com"},
{"github.com/owner/repo", "github.com"},
{"https://myhost.com:8080/path", "myhost.com"},
}
for _, tt := range tests {
got := githubhost.ParseHostFromURL(tt.url)
if got != tt.want {
t.Errorf("ParseHostFromURL(%q)=%q want %q", tt.url, got, tt.want)
}
}
}

func TestIsVisualStudioLegacyHostname(t *testing.T) {
if !githubhost.IsVisualStudioLegacyHostname("myorg.visualstudio.com") {
t.Error("expected true for *.visualstudio.com")
}
if githubhost.IsVisualStudioLegacyHostname("github.com") {
t.Error("expected false for github.com")
}
if githubhost.IsVisualStudioLegacyHostname("") {
t.Error("expected false for empty")
}
}
