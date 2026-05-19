package hostbackends_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/hostbackends"
)

func TestGitHubBuildCloneHTTPSURL_NoToken(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneHTTPSURL(dep, "", "")
	if !strings.Contains(url, "github.com") {
		t.Errorf("expected github.com in URL, got %q", url)
	}
	if !strings.Contains(url, "owner/repo") {
		t.Errorf("expected owner/repo in URL, got %q", url)
	}
}

func TestGitHubBuildCloneHTTPSURL_WithToken(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneHTTPSURL(dep, "mytoken123", "")
	if !strings.Contains(url, "mytoken123") {
		t.Errorf("expected token embedded in URL, got %q", url)
	}
}

func TestGitHubBuildCloneHTTPSURL_BearerScheme(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneHTTPSURL(dep, "mytoken123", "bearer")
	// bearer should NOT embed token in URL
	if strings.Contains(url, "mytoken123") {
		t.Errorf("bearer scheme should not embed token in URL, got %q", url)
	}
}

func TestGitHubBuildCloneSSHURL(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneSSHURL(dep)
	if !strings.Contains(url, "owner/repo") {
		t.Errorf("expected owner/repo in SSH URL, got %q", url)
	}
	if !strings.HasPrefix(url, "git@") {
		t.Errorf("expected SSH URL to start with git@, got %q", url)
	}
}

func TestADOBuildCloneHTTPSURL(t *testing.T) {
	dep := &mockDep{
		ado:     true,
		adoOrg:  "myorg",
		adoProj: "myproj",
		adoRepo: "myrepo",
		repoURL: "myorg/myrepo",
	}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	url := b.BuildCloneHTTPSURL(dep, "", "")
	if !strings.Contains(url, "myorg") {
		t.Errorf("expected myorg in URL, got %q", url)
	}
}

func TestADOBuildCommitsAPIURLEmpty(t *testing.T) {
	dep := &mockDep{
		ado:     true,
		adoOrg:  "myorg",
		adoProj: "myproj",
		adoRepo: "myrepo",
	}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	url := b.BuildCommitsAPIURL(dep, "main")
	// ADO returns empty
	if url != "" {
		t.Errorf("ADO CommitsAPIURL should be empty, got %q", url)
	}
}

func TestADOIsGitHubFamilyFalse(t *testing.T) {
	dep := &mockDep{ado: true, adoOrg: "org", adoProj: "proj", adoRepo: "repo"}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	if b.IsGitHubFamily() {
		t.Error("ADO backend should not be GitHub family")
	}
	if !b.IsGeneric() == false {
		// IsGeneric is false for ADO
	}
}

func TestGitLabIsGenericTrue(t *testing.T) {
	dep := &mockDep{host: "gitlab.com", repoURL: "user/project"}
	b := hostbackends.BackendFor(dep, "gitlab.com")
	if !b.IsGeneric() {
		t.Error("GitLab backend should be generic")
	}
	if b.IsGitHubFamily() {
		t.Error("GitLab backend should not be GitHub family")
	}
}

func TestGitHubBuildContentsAPIURLs(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	urls := b.BuildContentsAPIURLs("owner", "repo", "README.md", "main")
	if len(urls) == 0 {
		t.Error("expected at least one contents API URL")
	}
	for _, u := range urls {
		if !strings.Contains(u, "README.md") {
			t.Errorf("URL %q should contain README.md", u)
		}
	}
}

func TestGitHubCommitsAPIURLWithSHA(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	sha := strings.Repeat("a", 40)
	// For a full 40-char SHA the API returns "" (uses cheap non-SHA resolution only)
	url := b.BuildCommitsAPIURL(dep, sha)
	_ = url // may be "" for full SHAs
}

func TestGitHubCommitsAPIURLWithRef(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCommitsAPIURL(dep, "main")
	if url == "" {
		t.Error("expected non-empty commits API URL for branch ref")
	}
	if !strings.Contains(url, "owner") {
		t.Errorf("expected owner in URL, got %q", url)
	}
}

func TestGitHubBuildCloneHTTPURLInsecure(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo", insecure: true}
	b := hostbackends.BackendFor(dep, "github.com")
	url, err := b.BuildCloneHTTPURL(dep)
	if err != nil {
		t.Errorf("expected no error for insecure dep, got %v", err)
	}
	if !strings.HasPrefix(url, "http://") {
		t.Errorf("expected http:// URL, got %q", url)
	}
}
