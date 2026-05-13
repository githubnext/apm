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
