package hostbackends_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/hostbackends"
)

// mockDep implements DepRef for testing.
type mockDep struct {
	host    string
	port    *int
	repoURL string
	ado     bool
	adoOrg  string
	adoProj string
	adoRepo string
	insecure bool
}

func (m *mockDep) GetHost() string              { return m.host }
func (m *mockDep) GetPort() *int                { return m.port }
func (m *mockDep) GetRepoURL() string           { return m.repoURL }
func (m *mockDep) GetADOOrganization() string   { return m.adoOrg }
func (m *mockDep) GetADOProject() string        { return m.adoProj }
func (m *mockDep) GetADORepo() string           { return m.adoRepo }
func (m *mockDep) IsAzureDevOps() bool          { return m.ado }
func (m *mockDep) IsInsecure() bool             { return m.insecure }

func TestBackendFor_GitHub(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	if b.Kind() != "github" {
		t.Errorf("expected github, got %s", b.Kind())
	}
	if !b.IsGitHubFamily() {
		t.Error("expected IsGitHubFamily true")
	}
	if b.IsGeneric() {
		t.Error("expected IsGeneric false")
	}
}

func TestBackendFor_ADO(t *testing.T) {
	dep := &mockDep{
		ado:     true,
		adoOrg:  "myorg",
		adoProj: "myproj",
		adoRepo: "myrepo",
		repoURL: "myorg/myrepo",
	}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	if b.Kind() != "ado" {
		t.Errorf("expected ado, got %s", b.Kind())
	}
}

func TestBackendFor_GitLab(t *testing.T) {
	dep := &mockDep{host: "gitlab.com", repoURL: "user/project"}
	b := hostbackends.BackendFor(dep, "gitlab.com")
	if b.Kind() != "gitlab" {
		t.Errorf("expected gitlab, got %s", b.Kind())
	}
	if !b.IsGeneric() {
		t.Error("expected IsGeneric true for gitlab")
	}
}

func TestBackendForHost_GitHub(t *testing.T) {
	b := hostbackends.BackendForHost("github.com", nil)
	if b.Kind() != "github" {
		t.Errorf("expected github, got %s", b.Kind())
	}
}

func TestGitHubBackend_BuildCloneHTTPSURL_NoToken(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneHTTPSURL(dep, "", "")
	if !strings.HasPrefix(url, "https://github.com/") {
		t.Errorf("unexpected URL: %s", url)
	}
	if !strings.Contains(url, "owner/repo") {
		t.Errorf("URL missing repo: %s", url)
	}
}

func TestGitHubBackend_BuildCloneHTTPSURL_WithToken(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneHTTPSURL(dep, "mytoken", "")
	if !strings.Contains(url, "x-access-token") {
		t.Errorf("expected token in URL, got: %s", url)
	}
}

func TestGitHubBackend_BuildCloneHTTPSURL_BearerSkipsToken(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneHTTPSURL(dep, "mytoken", "bearer")
	if strings.Contains(url, "mytoken") {
		t.Errorf("bearer scheme should suppress token, got: %s", url)
	}
}

func TestGitHubBackend_BuildCloneSSHURL(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCloneSSHURL(dep)
	if !strings.HasPrefix(url, "git@github.com:") {
		t.Errorf("unexpected SSH URL: %s", url)
	}
}

func TestGitHubBackend_BuildCloneHTTPURL(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url, err := b.BuildCloneHTTPURL(dep)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "http://") {
		t.Errorf("unexpected HTTP URL: %s", url)
	}
}

func TestGitHubBackend_BuildCommitsAPIURL_Branch(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCommitsAPIURL(dep, "main")
	if !strings.Contains(url, "owner") || !strings.Contains(url, "repo") {
		t.Errorf("unexpected API URL: %s", url)
	}
}

func TestGitHubBackend_BuildCommitsAPIURL_SHA40_Empty(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	sha := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	url := b.BuildCommitsAPIURL(dep, sha)
	if url != "" {
		t.Errorf("expected empty URL for full SHA, got: %s", url)
	}
}

func TestGitHubBackend_BuildContentsAPIURLs(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	urls := b.BuildContentsAPIURLs("owner", "repo", "path/file.txt", "main")
	if len(urls) == 0 {
		t.Error("expected at least one contents URL")
	}
	if !strings.Contains(urls[0], "contents") {
		t.Errorf("expected contents URL, got: %s", urls[0])
	}
}

func TestADOBackend_BuildCloneHTTPSURL(t *testing.T) {
	dep := &mockDep{
		ado:     true,
		adoOrg:  "myorg",
		adoProj: "myproj",
		adoRepo: "myrepo",
		repoURL: "myorg/myrepo",
	}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	url := b.BuildCloneHTTPSURL(dep, "", "")
	if !strings.Contains(url, "dev.azure.com") {
		t.Errorf("unexpected ADO URL: %s", url)
	}
	if !strings.Contains(url, "myorg") {
		t.Errorf("expected org in ADO URL: %s", url)
	}
}

func TestADOBackend_NoOrg_ErrorURL(t *testing.T) {
	dep := &mockDep{ado: true, repoURL: "x/y"}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	url := b.BuildCloneHTTPSURL(dep, "", "")
	if !strings.HasPrefix(url, "error://") {
		t.Errorf("expected error URL when org missing, got: %s", url)
	}
}

func TestADOBackend_BuildCloneSSHURL(t *testing.T) {
	dep := &mockDep{
		ado:     true,
		adoOrg:  "org",
		adoProj: "proj",
		adoRepo: "repo",
	}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	url := b.BuildCloneSSHURL(dep)
	if !strings.Contains(url, "ssh.dev.azure.com") {
		t.Errorf("unexpected SSH URL: %s", url)
	}
}

func TestADOBackend_BuildCloneHTTPURL_Error(t *testing.T) {
	dep := &mockDep{ado: true, adoOrg: "org", adoProj: "proj", adoRepo: "repo"}
	b := hostbackends.BackendFor(dep, "dev.azure.com")
	_, err := b.BuildCloneHTTPURL(dep)
	if err == nil {
		t.Error("expected error for ADO plain HTTP clone")
	}
}

func TestGitLabBackend_BuildCloneHTTPSURL_WithToken(t *testing.T) {
	dep := &mockDep{host: "gitlab.com", repoURL: "user/proj"}
	b := hostbackends.BackendFor(dep, "gitlab.com")
	url := b.BuildCloneHTTPSURL(dep, "mytoken", "")
	if !strings.Contains(url, "oauth2") {
		t.Errorf("expected oauth2 in GitLab URL, got: %s", url)
	}
}

func TestGitLabBackend_BuildCommitsAPIURL(t *testing.T) {
	dep := &mockDep{host: "gitlab.com", repoURL: "user/proj"}
	b := hostbackends.BackendFor(dep, "gitlab.com")
	url := b.BuildCommitsAPIURL(dep, "main")
	if !strings.Contains(url, "projects") {
		t.Errorf("expected projects in GitLab commits URL, got: %s", url)
	}
}

func TestGenericBackend_BuildContentsAPIURLs(t *testing.T) {
	b := hostbackends.BackendForHost("gitea.example.com", nil)
	urls := b.BuildContentsAPIURLs("user", "repo", "file.txt", "main")
	if len(urls) < 2 {
		t.Errorf("expected 2 generic contents URLs, got %d", len(urls))
	}
}

func TestBackendFor_FallbackToDefault(t *testing.T) {
	dep := &mockDep{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "")
	if b == nil {
		t.Error("expected non-nil backend with empty fallback")
	}
}
