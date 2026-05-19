package downloadstrategies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/githubnext/apm/internal/core/auth"
)

// stubHost is a minimal HostProvider for tests.
type stubHost struct {
	githubToken       string
	adoToken          string
	artifactoryToken  string
	githubHost        string
}

func (h *stubHost) GithubToken() string         { return h.githubToken }
func (h *stubHost) AdoToken() string             { return h.adoToken }
func (h *stubHost) ArtifactoryToken() string     { return h.artifactoryToken }
func (h *stubHost) GithubHost() string           { return h.githubHost }
func (h *stubHost) AuthResolver() *auth.AuthResolver { return nil }
func (h *stubHost) ResilientGet(reqURL string, headers map[string]string, timeoutSecs int) (*http.Response, error) {
	return ResilientGet(reqURL, headers, timeoutSecs, 0)
}

func newStub() *DownloadDelegate {
	return New(&stubHost{})
}

func TestBuildRepoURL_SSHNoPort(t *testing.T) {
	d := newStub()
	url := d.BuildRepoURL(BuildRepoURLOptions{
		RepoRef: "owner/repo",
		UseSSH:  true,
	})
	if url == "" {
		t.Error("expected non-empty SSH URL")
	}
}

func TestBuildRepoURL_HTTPSNoToken(t *testing.T) {
	d := newStub()
	url := d.BuildRepoURL(BuildRepoURLOptions{
		RepoRef: "owner/repo",
		UseSSH:  false,
	})
	if url == "" {
		t.Error("expected non-empty HTTPS URL")
	}
}

func TestBuildRepoURL_HTTPSWithToken(t *testing.T) {
	d := newStub()
	url := d.BuildRepoURL(BuildRepoURLOptions{
		RepoRef: "owner/repo",
		UseSSH:  false,
		Token:   "ghp_mytoken",
	})
	if url == "" {
		t.Error("expected non-empty HTTPS URL with token")
	}
}

func TestBuildRepoURL_HTTPSBearer(t *testing.T) {
	d := newStub()
	url := d.BuildRepoURL(BuildRepoURLOptions{
		RepoRef:    "owner/repo",
		UseSSH:     false,
		Token:      "mytoken",
		AuthScheme: "bearer",
	})
	if url == "" {
		t.Error("expected non-empty bearer HTTPS URL")
	}
}

func TestGetArtifactoryHeaders_NoToken(t *testing.T) {
	d := newStub()
	h := d.GetArtifactoryHeaders()
	if h == nil {
		t.Error("expected non-nil headers map")
	}
}

func TestResilientGet_Redirect(t *testing.T) {
	final := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer final.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, final.URL, http.StatusFound)
	}))
	defer srv.Close()

	resp, err := ResilientGet(srv.URL, nil, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 after redirect, got %d", resp.StatusCode)
	}
}

func TestResilientGet_Timeout(t *testing.T) {
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-done
	}))
	defer srv.Close()
	defer close(done)

	_, err := ResilientGet(srv.URL, nil, 1, 0)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestBuildSSHURL_EmptyHost(t *testing.T) {
	url := buildSSHURL("", "owner/repo", 0)
	if url == "" {
		t.Error("expected non-empty SSH URL for empty host (defaults to github.com)")
	}
}

func TestBuildHTTPSCloneURL_EmptyHostDefaultsGitHub(t *testing.T) {
	url := buildHTTPSCloneURL("", "owner/repo", "", 0)
	if url == "" {
		t.Error("expected non-empty HTTPS URL")
	}
}

func TestBuildHTTPSCloneURL_BearerToken(t *testing.T) {
	url := buildHTTPSCloneURL("github.com", "owner/repo", "tok", 0)
	if url == "" {
		t.Error("expected non-empty HTTPS URL with token")
	}
}

func TestBuildADOAPIURL_WithRef(t *testing.T) {
	got := buildADOAPIURL("myorg", "myproj", "myrepo", "path/file.yaml", "main", "")
	if got == "" {
		t.Error("expected non-empty ADO API URL")
	}
}

func TestBuildADOAPIURL_CustomHostAndRef(t *testing.T) {
	got := buildADOAPIURL("org2", "proj2", "repo2", "dir/a.yaml", "feature/x", "ado.example.com")
	if got == "" {
		t.Error("expected non-empty custom-host ADO URL")
	}
}

func TestNew_ReturnsNonNilDelegate(t *testing.T) {
	d := newStub()
	if d == nil {
		t.Error("New(stub) should return non-nil DownloadDelegate")
	}
}

func TestArtifactoryDownloadResult_ErrField(t *testing.T) {
	var r ArtifactoryDownloadResult
	if r.Err != nil {
		t.Error("zero value should have nil Err")
	}
	if r.Data != nil {
		t.Error("zero value should have nil Data")
	}
}

func TestResilientGet_MultipleHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") == "" || r.Header.Get("Accept") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	headers := map[string]string{
		"X-Token": "abc",
		"Accept":  "application/json",
	}
	resp, err := ResilientGet(srv.URL, headers, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
