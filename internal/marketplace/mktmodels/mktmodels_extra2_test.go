package mktmodels

import (
	"testing"
)

func TestMarketplaceSource_ZeroValue(t *testing.T) {
	var s MarketplaceSource
	if s.Name != "" || s.Owner != "" {
		t.Error("zero value should have empty fields")
	}
}

func TestNewMarketplaceSource_NonEmptyDefaults(t *testing.T) {
	s := NewMarketplaceSource("test", "owner", "repo", "", "", "")
	if s.Host != "github.com" {
		t.Errorf("Host = %q, want github.com", s.Host)
	}
	if s.Branch != "main" {
		t.Errorf("Branch = %q, want main", s.Branch)
	}
	if s.Path != "marketplace.json" {
		t.Errorf("Path = %q, want marketplace.json", s.Path)
	}
}

func TestMarketplaceSource_ToDict_OnlyRequiredKeys(t *testing.T) {
	s := MarketplaceSource{Name: "n", Owner: "o", Repo: "r", Host: "github.com", Branch: "main", Path: "marketplace.json"}
	d := s.ToDict()
	if d["name"] != "n" || d["owner"] != "o" || d["repo"] != "r" {
		t.Error("missing required keys")
	}
	if _, ok := d["host"]; ok {
		t.Error("default host should be omitted")
	}
}

func TestMarketplacePlugin_ZeroValue(t *testing.T) {
	var p MarketplacePlugin
	if p.Name != "" {
		t.Error("zero value Name should be empty")
	}
	if len(p.Tags) != 0 {
		t.Error("zero value Tags should be empty")
	}
}

func TestMarketplacePlugin_MatchesQuery_NoMatch(t *testing.T) {
	p := &MarketplacePlugin{Name: "foo", Description: "bar"}
	if p.MatchesQuery("xyz") {
		t.Error("should not match xyz")
	}
}

func TestMarketplacePlugin_MatchesQuery_Tag(t *testing.T) {
	p := &MarketplacePlugin{Name: "foo", Tags: []string{"editor"}}
	if !p.MatchesQuery("editor") {
		t.Error("should match by tag")
	}
}

func TestMarketplaceManifest_ZeroValue(t *testing.T) {
	var m MarketplaceManifest
	if m.FindPlugin("any") != nil {
		t.Error("empty manifest should return nil for any plugin")
	}
}

func TestMarketplaceManifest_Search_Empty(t *testing.T) {
	var m MarketplaceManifest
	results := m.Search("q")
	if len(results) != 0 {
		t.Error("empty manifest should return no results")
	}
}

func TestParseMarketplaceJSONBytes_InvalidJSON(t *testing.T) {
	_, err := ParseMarketplaceJSONBytes([]byte("{invalid"), "test")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
