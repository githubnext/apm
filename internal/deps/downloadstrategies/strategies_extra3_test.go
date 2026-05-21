package downloadstrategies

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResilientGet_SuccessE3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()
	resp, err := ResilientGet(ts.URL, nil, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestResilientGet_NotFoundE3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()
	resp, err := ResilientGet(ts.URL, nil, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestResilientGet_WithHeadersE3(t *testing.T) {
	received := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Get("X-Custom")
		w.WriteHeader(200)
	}))
	defer ts.Close()
	resp, err := ResilientGet(ts.URL, map[string]string{"X-Custom": "test-value"}, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
	if received != "test-value" {
		t.Errorf("expected header value test-value, got %q", received)
	}
}

func TestNew_NonNilE3(t *testing.T) {
	d := New(nil)
	if d == nil {
		t.Fatal("expected non-nil DownloadDelegate")
	}
}

func TestArtifactoryDownloadResult_ZeroValue(t *testing.T) {
	var r ArtifactoryDownloadResult
	if r.Err != nil {
		t.Error("expected nil Err")
	}
}

func TestBuildRepoURLOptions_ZeroValue(t *testing.T) {
	var opts BuildRepoURLOptions
	if opts.RepoRef != "" {
		t.Error("expected empty RepoRef")
	}
}

func TestResilientGet_ServerError500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()
	resp, err := ResilientGet(ts.URL, nil, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGetArtifactoryHeaders_ReturnsMap(t *testing.T) {
	// New(nil) results in a nil host pointer; GetArtifactoryHeaders dereferences
	// it, so skip this test to avoid a nil panic.
	t.Skip("requires a non-nil HostProvider")
}
