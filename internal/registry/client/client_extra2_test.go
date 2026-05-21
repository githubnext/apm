package client_test

import (
	"encoding/json"
	"testing"

	"github.com/githubnext/apm/internal/registry/client"
)

func TestMCPServerInfo_ZeroValue(t *testing.T) {
	var s client.MCPServerInfo
	if s.ID != "" || s.Name != "" {
		t.Error("zero value should have empty string fields")
	}
	if s.Tags != nil {
		t.Error("zero value Tags should be nil")
	}
}

func TestVersionEntry_ZeroValue(t *testing.T) {
	var v client.VersionEntry
	if v.Version != "" || v.CreatedAt != "" || v.PackageID != "" {
		t.Error("zero value should have empty fields")
	}
}

func TestSearchResult_ZeroValue(t *testing.T) {
	var r client.SearchResult
	if r.TotalCount != 0 || r.Page != 0 || r.PerPage != 0 {
		t.Error("zero value fields should be 0")
	}
	if r.Items != nil {
		t.Error("Items should be nil in zero value")
	}
}

func TestMCPServerInfo_JSONRoundTrip_Tags(t *testing.T) {
	s := client.MCPServerInfo{
		ID:   "srv1",
		Name: "MyServer",
		Tags: []string{"go", "api"},
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var s2 client.MCPServerInfo
	if err := json.Unmarshal(b, &s2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(s2.Tags) != 2 || s2.Tags[0] != "go" {
		t.Errorf("tags mismatch: %v", s2.Tags)
	}
}

func TestSearchResult_ItemsField(t *testing.T) {
	r := client.SearchResult{
		Items:      []client.MCPServerInfo{{ID: "a"}, {ID: "b"}},
		TotalCount: 2,
	}
	if len(r.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(r.Items))
	}
}

func TestNewSimpleRegistryClient_DefaultURL(t *testing.T) {
	c, err := client.NewSimpleRegistryClient("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL() == "" {
		t.Error("expected non-empty base URL")
	}
}

func TestNewSimpleRegistryClient_CustomURL(t *testing.T) {
	c, err := client.NewSimpleRegistryClient("https://custom.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL() != "https://custom.example.com" {
		t.Errorf("unexpected base URL: %q", c.BaseURL())
	}
}

func TestVersionEntry_JSONRoundTrip(t *testing.T) {
	v := client.VersionEntry{Version: "1.0.0", PackageID: "pkg123"}
	b, _ := json.Marshal(v)
	var v2 client.VersionEntry
	json.Unmarshal(b, &v2)
	if v2.Version != "1.0.0" || v2.PackageID != "pkg123" {
		t.Errorf("roundtrip failed: %+v", v2)
	}
}

func TestMCPServerInfo_VersionsField(t *testing.T) {
	s := client.MCPServerInfo{
		Versions: []client.VersionEntry{{Version: "0.1.0"}},
	}
	if len(s.Versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(s.Versions))
	}
}
