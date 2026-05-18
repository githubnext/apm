package mktmodels

import (
	"testing"
)

func TestNewMarketplaceSource_defaults(t *testing.T) {
	s := NewMarketplaceSource("my-mkt", "acme", "registry", "", "", "")
	if s.Host != "github.com" {
		t.Errorf("expected default host github.com, got %s", s.Host)
	}
	if s.Branch != "main" {
		t.Errorf("expected default branch main, got %s", s.Branch)
	}
	if s.Path != "marketplace.json" {
		t.Errorf("expected default path marketplace.json, got %s", s.Path)
	}
}

func TestNewMarketplaceSource_custom(t *testing.T) {
	s := NewMarketplaceSource("mkt", "org", "repo", "ghe.company.com", "release", "index.json")
	if s.Host != "ghe.company.com" {
		t.Errorf("expected custom host, got %s", s.Host)
	}
	if s.Branch != "release" {
		t.Errorf("expected release branch, got %s", s.Branch)
	}
	if s.Path != "index.json" {
		t.Errorf("expected index.json path, got %s", s.Path)
	}
}

func TestMarketplaceSource_ToDict_defaults(t *testing.T) {
	s := NewMarketplaceSource("n", "o", "r", "", "", "")
	d := s.ToDict()
	if _, ok := d["host"]; ok {
		t.Error("default host should be omitted from ToDict")
	}
	if _, ok := d["branch"]; ok {
		t.Error("default branch should be omitted from ToDict")
	}
	if _, ok := d["path"]; ok {
		t.Error("default path should be omitted from ToDict")
	}
	if d["name"] != "n" || d["owner"] != "o" || d["repo"] != "r" {
		t.Errorf("unexpected ToDict values: %v", d)
	}
}

func TestMarketplaceSource_ToDict_custom(t *testing.T) {
	s := NewMarketplaceSource("n", "o", "r", "ghe.example.com", "dev", "custom/path.json")
	d := s.ToDict()
	if d["host"] != "ghe.example.com" {
		t.Errorf("expected ghe.example.com, got %s", d["host"])
	}
	if d["branch"] != "dev" {
		t.Errorf("expected dev branch, got %s", d["branch"])
	}
	if d["path"] != "custom/path.json" {
		t.Errorf("expected custom path, got %s", d["path"])
	}
}

func TestMarketplacePlugin_MatchesQuery_name(t *testing.T) {
	p := MarketplacePlugin{Name: "MyPlugin", Description: "Some desc", Tags: []string{"ai"}}
	if !p.MatchesQuery("myplugin") {
		t.Error("should match name case-insensitively")
	}
	if !p.MatchesQuery("PLUGIN") {
		t.Error("should match partial name case-insensitively")
	}
}

func TestMarketplacePlugin_MatchesQuery_description(t *testing.T) {
	p := MarketplacePlugin{Name: "foo", Description: "Helps you write Go code faster"}
	if !p.MatchesQuery("go code") {
		t.Error("should match description")
	}
}

func TestMarketplacePlugin_MatchesQuery_tag(t *testing.T) {
	p := MarketplacePlugin{Name: "foo", Tags: []string{"automation", "testing"}}
	if !p.MatchesQuery("testing") {
		t.Error("should match tag")
	}
}

func TestMarketplacePlugin_MatchesQuery_nomatch(t *testing.T) {
	p := MarketplacePlugin{Name: "alpha", Description: "beta gamma", Tags: []string{"delta"}}
	if p.MatchesQuery("zeta") {
		t.Error("should not match unrelated query")
	}
}

func TestMarketplaceManifest_FindPlugin_caseless(t *testing.T) {
	m := MarketplaceManifest{
		Plugins: []MarketplacePlugin{
			{Name: "MyTool"},
			{Name: "OtherTool"},
		},
	}
	p := m.FindPlugin("mytool")
	if p == nil {
		t.Fatal("expected to find plugin case-insensitively")
	}
	if p.Name != "MyTool" {
		t.Errorf("expected MyTool, got %s", p.Name)
	}
}

func TestMarketplaceManifest_FindPlugin_missing(t *testing.T) {
	m := MarketplaceManifest{Plugins: []MarketplacePlugin{{Name: "Alpha"}}}
	p := m.FindPlugin("Beta")
	if p != nil {
		t.Error("expected nil for missing plugin")
	}
}

func TestMarketplaceManifest_Search(t *testing.T) {
	m := MarketplaceManifest{
		Plugins: []MarketplacePlugin{
			{Name: "go-helper", Tags: []string{"golang"}},
			{Name: "python-helper", Tags: []string{"python"}},
			{Name: "general", Description: "supports golang and python"},
		},
	}
	results := m.Search("golang")
	if len(results) < 1 {
		t.Errorf("expected at least 1 result, got %d", len(results))
	}
}

func TestMarketplaceManifest_Search_empty(t *testing.T) {
	m := MarketplaceManifest{}
	results := m.Search("anything")
	if len(results) != 0 {
		t.Errorf("expected no results from empty manifest")
	}
}

func TestParseMarketplaceJSONBytes_invalid(t *testing.T) {
	_, err := ParseMarketplaceJSONBytes([]byte("not json"), "test")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseMarketplaceJSONBytes_empty(t *testing.T) {
	m, err := ParseMarketplaceJSONBytes([]byte(`{}`), "test-source")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Plugins) != 0 {
		t.Errorf("expected no plugins, got %d", len(m.Plugins))
	}
}
