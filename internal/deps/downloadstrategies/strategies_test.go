package downloadstrategies

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildSSHURL_NoPort(t *testing.T) {
	got := buildSSHURL("github.com", "owner/repo", 0)
	want := "git@github.com:owner/repo.git"
	if got != want {
		t.Errorf("buildSSHURL no-port: got %q want %q", got, want)
	}
}

func TestBuildSSHURL_WithPort(t *testing.T) {
	got := buildSSHURL("ghe.example.com", "org/project", 2222)
	want := "ssh://git@ghe.example.com:2222/org/project.git"
	if got != want {
		t.Errorf("buildSSHURL with-port: got %q want %q", got, want)
	}
}

func TestBuildHTTPSCloneURL_NoToken(t *testing.T) {
	got := buildHTTPSCloneURL("github.com", "owner/repo", "", 0)
	want := "https://github.com/owner/repo.git"
	if got != want {
		t.Errorf("buildHTTPSCloneURL no-token: got %q want %q", got, want)
	}
}

func TestBuildHTTPSCloneURL_WithToken(t *testing.T) {
	got := buildHTTPSCloneURL("github.com", "owner/repo", "mytoken", 0)
	want := "https://x-access-token:mytoken@github.com/owner/repo.git"
	if got != want {
		t.Errorf("buildHTTPSCloneURL with-token: got %q want %q", got, want)
	}
}

func TestBuildHTTPSCloneURL_WithPort(t *testing.T) {
	got := buildHTTPSCloneURL("ghe.corp.com", "org/repo", "", 8443)
	want := "https://ghe.corp.com:8443/org/repo.git"
	if got != want {
		t.Errorf("buildHTTPSCloneURL with-port: got %q want %q", got, want)
	}
}

func TestBuildADOAPIURL_DefaultHost(t *testing.T) {
	got := buildADOAPIURL("myorg", "myproject", "myrepo", "/path/file.txt", "main", "")
	if got == "" {
		t.Error("expected non-empty ADO API URL")
	}
	if len(got) < 10 {
		t.Errorf("URL too short: %s", got)
	}
}

func TestBuildADOAPIURL_CustomHost(t *testing.T) {
	got := buildADOAPIURL("org", "proj", "repo", "/file.txt", "main", "ado.corp.com")
	if got == "" {
		t.Error("expected non-empty ADO API URL with custom host")
	}
}

func TestResilientGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	resp, err := ResilientGet(srv.URL, nil, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestResilientGet_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	resp, err := ResilientGet(srv.URL, nil, 5, 1)
	// 404 is not retried; should return response with no error
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestNew_NotNil(t *testing.T) {
	d := New(nil)
	if d == nil {
		t.Error("New returned nil")
	}
}
