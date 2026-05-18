package downloadstrategies

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildSSHURL_ZeroPort(t *testing.T) {
	got := buildSSHURL("gitlab.com", "group/project", 0)
	want := "git@gitlab.com:group/project.git"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildHTTPSCloneURL_TokenWithPort(t *testing.T) {
	got := buildHTTPSCloneURL("ghe.corp.com", "org/repo", "token123", 8443)
	want := "https://x-access-token:token123@ghe.corp.com:8443/org/repo.git"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildHTTPSCloneURL_NoTokenNoPort(t *testing.T) {
	got := buildHTTPSCloneURL("github.com", "a/b", "", 0)
	if got != "https://github.com/a/b.git" {
		t.Errorf("got %q", got)
	}
}

func TestBuildADOAPIURL_EmptyHost(t *testing.T) {
	got := buildADOAPIURL("myorg", "myproj", "myrepo", "/readme.md", "main", "")
	if got == "" {
		t.Error("expected non-empty URL")
	}
}

func TestResilientGet_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	resp, err := ResilientGet(srv.URL, nil, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestResilientGet_WithHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom") != "value" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	resp, err := ResilientGet(srv.URL, map[string]string{"X-Custom": "value"}, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestNew_WithNilHost(t *testing.T) {
	d := New(nil)
	if d == nil {
		t.Error("New(nil) returned nil")
	}
}

func TestBuildSSHURL_NonStandardHost(t *testing.T) {
	got := buildSSHURL("git.internal.corp", "team/service", 22)
	if got == "" {
		t.Error("expected non-empty SSH URL")
	}
}

func TestBuildADOAPIURL_ReturnsValidURL(t *testing.T) {
	got := buildADOAPIURL("org1", "proj1", "repo1", "/path/to/file.yaml", "feature/branch", "")
	if len(got) < 20 {
		t.Errorf("URL too short: %q", got)
	}
}
