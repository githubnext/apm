package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello..."},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "abc..."},
	}
	for _, tc := range tests {
		got := truncate(tc.input, tc.maxLen)
		if got != tc.want {
			t.Errorf("truncate(%q, %d) = %q; want %q", tc.input, tc.maxLen, got, tc.want)
		}
	}
}

func TestNewSimpleRegistryClient_InvalidURL(t *testing.T) {
	tests := []struct {
		url string
	}{
		{"not-a-url"},
		{"ftp://example.com"},
		{"://bad"},
	}
	for _, tc := range tests {
		_, err := NewSimpleRegistryClient(tc.url)
		if err == nil {
			t.Errorf("NewSimpleRegistryClient(%q): expected error, got nil", tc.url)
		}
	}
}

func TestNewSimpleRegistryClient_HTTPRejectedWithoutFlag(t *testing.T) {
	os.Unsetenv("MCP_REGISTRY_ALLOW_HTTP")
	_, err := NewSimpleRegistryClient("http://example.com")
	if err == nil {
		t.Error("expected error for http:// without MCP_REGISTRY_ALLOW_HTTP=1")
	}
}

func TestNewSimpleRegistryClient_HTTPAllowedWithFlag(t *testing.T) {
	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, err := NewSimpleRegistryClient("http://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL() != "http://example.com" {
		t.Errorf("BaseURL() = %q; want %q", c.BaseURL(), "http://example.com")
	}
}

func TestNewSimpleRegistryClient_TrailingSlashStripped(t *testing.T) {
	c, err := NewSimpleRegistryClient("https://example.com/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL() != "https://example.com" {
		t.Errorf("BaseURL() = %q; want no trailing slash", c.BaseURL())
	}
}

func TestNewSimpleRegistryClient_EnvOverride(t *testing.T) {
	t.Setenv("MCP_REGISTRY_URL", "https://custom.registry.example.com")
	c, err := NewSimpleRegistryClient("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL() != "https://custom.registry.example.com" {
		t.Errorf("BaseURL() = %q; want custom URL", c.BaseURL())
	}
}

func TestSearchServers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := SearchResult{
			Items:      []MCPServerInfo{{ID: "srv1", Name: "MyServer"}},
			TotalCount: 1,
			Page:       1,
			PerPage:    20,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, err := NewSimpleRegistryClient(srv.URL)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	res, err := c.SearchServers("myserver", 0, 0)
	if err != nil {
		t.Fatalf("SearchServers: %v", err)
	}
	if len(res.Items) != 1 || res.Items[0].ID != "srv1" {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestGetServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := MCPServerInfo{ID: "abc-123", Name: "TestServer"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, _ := NewSimpleRegistryClient(srv.URL)

	info, err := c.GetServer("abc-123")
	if err != nil {
		t.Fatalf("GetServer: %v", err)
	}
	if info.ID != "abc-123" {
		t.Errorf("ID = %q; want abc-123", info.ID)
	}
}

func TestGetServer_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, _ := NewSimpleRegistryClient(srv.URL)

	_, err := c.GetServer("missing")
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestGetServerVersions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []VersionEntry{{Version: "1.0.0", PackageID: "pkg1"}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versions)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, _ := NewSimpleRegistryClient(srv.URL)

	versions, err := c.GetServerVersions("srv1")
	if err != nil {
		t.Fatalf("GetServerVersions: %v", err)
	}
	if len(versions) != 1 || versions[0].Version != "1.0.0" {
		t.Errorf("unexpected versions: %+v", versions)
	}
}

func TestListServers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := SearchResult{Items: []MCPServerInfo{{ID: "s1"}, {ID: "s2"}}, TotalCount: 2}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))
	defer srv.Close()

	t.Setenv("MCP_REGISTRY_ALLOW_HTTP", "1")
	c, _ := NewSimpleRegistryClient(srv.URL)

	res, err := c.ListServers(0, 0)
	if err != nil {
		t.Fatalf("ListServers: %v", err)
	}
	if len(res.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(res.Items))
	}
}

func TestResolveTimeout_Defaults(t *testing.T) {
	os.Unsetenv("MCP_REGISTRY_CONNECT_TIMEOUT")
	os.Unsetenv("MCP_REGISTRY_READ_TIMEOUT")
	conn, read := resolveTimeout()
	if conn != defaultConnectTimeout {
		t.Errorf("connect timeout = %v; want %v", conn, defaultConnectTimeout)
	}
	if read != defaultReadTimeout {
		t.Errorf("read timeout = %v; want %v", read, defaultReadTimeout)
	}
}

func TestResolveTimeout_EnvOverride(t *testing.T) {
	t.Setenv("MCP_REGISTRY_CONNECT_TIMEOUT", "5.0")
	t.Setenv("MCP_REGISTRY_READ_TIMEOUT", "60.0")
	conn, read := resolveTimeout()
	if conn != 5.0 {
		t.Errorf("connect timeout = %v; want 5.0", conn)
	}
	if read != 60.0 {
		t.Errorf("read timeout = %v; want 60.0", read)
	}
}
