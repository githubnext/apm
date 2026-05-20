package hostbackends_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/hostbackends"
)

// mockDep2 is a separate mock to avoid redeclaration errors.
type mockDep2 struct {
	host     string
	port     *int
	repoURL  string
	ado      bool
	adoOrg   string
	adoProj  string
	adoRepo  string
	insecure bool
}

func (m *mockDep2) GetHost() string            { return m.host }
func (m *mockDep2) GetPort() *int              { return m.port }
func (m *mockDep2) GetRepoURL() string         { return m.repoURL }
func (m *mockDep2) GetADOOrganization() string { return m.adoOrg }
func (m *mockDep2) GetADOProject() string      { return m.adoProj }
func (m *mockDep2) GetADORepo() string         { return m.adoRepo }
func (m *mockDep2) IsAzureDevOps() bool        { return m.ado }
func (m *mockDep2) IsInsecure() bool           { return m.insecure }

func TestBackendForHost_GitHub_extra2(t *testing.T) {
	b := hostbackends.BackendForHost("github.com", nil)
	if !b.IsGitHubFamily() {
		t.Error("expected IsGitHubFamily=true for github.com")
	}
	if b.IsGeneric() {
		t.Error("expected IsGeneric=false for github.com")
	}
}

func TestBackendForHost_GitLab(t *testing.T) {
	b := hostbackends.BackendForHost("gitlab.com", nil)
	if b.IsGitHubFamily() {
		t.Error("expected IsGitHubFamily=false for gitlab.com")
	}
}

func TestBackendForHost_Generic(t *testing.T) {
	b := hostbackends.BackendForHost("bitbucket.org", nil)
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestBackendFor_GitLabSSH(t *testing.T) {
	dep := &mockDep2{host: "gitlab.com", repoURL: "user/project"}
	b := hostbackends.BackendFor(dep, "gitlab.com")
	url := b.BuildCloneSSHURL(dep)
	if !strings.Contains(url, "user/project") {
		t.Errorf("expected user/project in SSH URL, got %q", url)
	}
}

func TestBackendFor_ADO_Kind(t *testing.T) {
	dep := &mockDep2{
		ado:     true,
		adoOrg:  "myorg",
		adoProj: "myproj",
		adoRepo: "myrepo",
		repoURL: "myorg/myrepo",
	}
	b := hostbackends.BackendFor(dep, "")
	if b.Kind() != "ado" {
		t.Errorf("expected kind='ado', got %q", b.Kind())
	}
}

func TestBackendFor_ADO_HTTPSURLFormat(t *testing.T) {
	dep := &mockDep2{
		ado:     true,
		adoOrg:  "corp",
		adoProj: "proj",
		adoRepo: "repo",
		repoURL: "corp/repo",
	}
	b := hostbackends.BackendFor(dep, "")
	url := b.BuildCloneHTTPSURL(dep, "", "")
	if !strings.Contains(url, "corp") {
		t.Errorf("expected org in ADO URL, got %q", url)
	}
}

func TestBackendFor_GitHub_CommitsAPIURL(t *testing.T) {
	dep := &mockDep2{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	url := b.BuildCommitsAPIURL(dep, "main")
	if !strings.Contains(url, "owner") || !strings.Contains(url, "repo") {
		t.Errorf("expected owner/repo in commits API URL, got %q", url)
	}
}

func TestBackendFor_GitHub_ContentsAPIURLs(t *testing.T) {
	dep := &mockDep2{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	urls := b.BuildContentsAPIURLs("owner", "repo", "README.md", "main")
	if len(urls) == 0 {
		t.Error("expected at least one contents API URL")
	}
	for _, u := range urls {
		if !strings.Contains(u, "README.md") {
			t.Errorf("expected README.md in URL, got %q", u)
		}
	}
}

func TestBackendFor_GitHub_HostInfo_NonNil(t *testing.T) {
	dep := &mockDep2{repoURL: "owner/repo"}
	b := hostbackends.BackendFor(dep, "github.com")
	hi := b.GetHostInfo()
	_ = hi // just ensure it doesn't panic
}
