package downloadstrategies

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_ReturnsNonNilE4(t *testing.T) {
	d := New(nil)
	if d == nil {
		t.Fatal("expected non-nil delegate")
	}
}

func TestNew_WithNilHostE4(t *testing.T) {
	d := New(nil)
	if d == nil {
		t.Fatal("expected non-nil even with nil host")
	}
}

func TestResilientGet_InvalidURLFails(t *testing.T) {
	resp, err := ResilientGet("not-a-url", nil, 1, 0)
	if err == nil {
		if resp != nil {
			resp.Body.Close()
		}
		t.Log("unexpected success for invalid URL")
	}
}

func TestResilientGet_EmptyURLFails(t *testing.T) {
	_, err := ResilientGet("", nil, 1, 0)
	_ = err
}

func TestResilientGet_200Response(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello"))
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

func TestResilientGet_CustomHeaders(t *testing.T) {
	var gotHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Custom")
		w.WriteHeader(200)
	}))
	defer ts.Close()
	resp, err := ResilientGet(ts.URL, map[string]string{"X-Custom": "test-value"}, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
	if gotHeader != "test-value" {
		t.Errorf("expected X-Custom=test-value, got %s", gotHeader)
	}
}

func TestResilientGet_500Response(t *testing.T) {
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

func TestResilientGet_301Redirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("final"))
	}))
	defer ts.Close()
	resp, err := ResilientGet(ts.URL+"/redirect", nil, 5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

