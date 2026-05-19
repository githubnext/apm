package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMCPServerInfo_Fields(t *testing.T) {
	info := MCPServerInfo{
		ID:          "uuid-abc",
		Name:        "my-server",
		Description: "A test server",
		Publisher:   "acme",
		Homepage:    "https://example.com",
		Repository:  "https://github.com/acme/server",
		License:     "MIT",
		Tags:        []string{"mcp", "go"},
	}
	if info.ID != "uuid-abc" {
		t.Errorf("unexpected ID %q", info.ID)
	}
	if len(info.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(info.Tags))
	}
}

func TestMCPServerInfo_JSONRoundTrip(t *testing.T) {
	info := MCPServerInfo{
		ID:   "id-1",
		Name: "srv",
		Tags: []string{"tag-a"},
	}
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var out MCPServerInfo
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if out.ID != info.ID || out.Name != info.Name {
		t.Errorf("round-trip mismatch: got %+v", out)
	}
	if len(out.Tags) != 1 || out.Tags[0] != "tag-a" {
		t.Errorf("unexpected Tags after round-trip: %v", out.Tags)
	}
}

func TestVersionEntry_Fields(t *testing.T) {
	ve := VersionEntry{
		Version:   "1.2.3",
		CreatedAt: "2025-01-01T00:00:00Z",
		PackageID: "pkg-xyz",
	}
	if ve.Version != "1.2.3" {
		t.Errorf("unexpected Version %q", ve.Version)
	}
	if ve.PackageID != "pkg-xyz" {
		t.Errorf("unexpected PackageID %q", ve.PackageID)
	}
}

func TestSearchResult_Fields(t *testing.T) {
	sr := SearchResult{
		Items:      []MCPServerInfo{{ID: "a"}, {ID: "b"}},
		TotalCount: 10,
		Page:       2,
		PerPage:    5,
	}
	if len(sr.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(sr.Items))
	}
	if sr.TotalCount != 10 {
		t.Errorf("unexpected TotalCount %d", sr.TotalCount)
	}
}

func TestBaseURL_ReturnsCorrectURL(t *testing.T) {
	c, err := NewSimpleRegistryClient("https://registry.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL() != "https://registry.example.com" {
		t.Errorf("unexpected BaseURL: %q", c.BaseURL())
	}
}

func TestGetServer_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, err := NewSimpleRegistryClient(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = c.GetServer("nonexistent")
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestSearchServers_ValidResponse(t *testing.T) {
	result := SearchResult{
		Items:      []MCPServerInfo{{ID: "s1", Name: "server-one"}},
		TotalCount: 1,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(result)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, err := NewSimpleRegistryClient(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := c.SearchServers("server", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].ID != "s1" {
		t.Errorf("unexpected search result: %+v", got)
	}
}

func TestListServers_DefaultPagination(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		q := r.URL.Query()
		if q.Get("page") == "" || q.Get("per_page") == "" {
			t.Errorf("expected pagination params in URL: %q", r.URL.String())
		}
		json.NewEncoder(w).Encode(SearchResult{TotalCount: 0})
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, err := NewSimpleRegistryClient(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = c.ListServers(0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 HTTP call, got %d", callCount)
	}
}

func TestGetServerVersions_Success(t *testing.T) {
	versions := []VersionEntry{{Version: "1.0.0"}, {Version: "1.1.0"}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(versions)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, err := NewSimpleRegistryClient(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := c.GetServerVersions("my-server")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 versions, got %d", len(got))
	}
}

func TestTruncate_EdgeCases(t *testing.T) {
	// Already tested in base test; add boundary cases
	s := truncate("", 0)
	if s != "" {
		t.Errorf("expected empty for empty input, got %q", s)
	}
	s2 := truncate("a", 1)
	if s2 != "a" {
		t.Errorf("expected 'a', got %q", s2)
	}
}
