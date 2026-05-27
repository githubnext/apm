package marketplace_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace"
)

// TestParityMarketplaceSourceDefaults verifies default field values mirror Python.
func TestParityMarketplaceSourceDefaults(t *testing.T) {
	s := marketplace.DefaultMarketplaceSource("acme", "acme-org", "tools")
	if s.Host != "github.com" {
		t.Errorf("expected host=github.com, got %s", s.Host)
	}
	if s.Branch != "main" {
		t.Errorf("expected branch=main, got %s", s.Branch)
	}
	if s.Path != "marketplace.json" {
		t.Errorf("expected path=marketplace.json, got %s", s.Path)
	}
}

// TestParityMarketplaceSourceToDictOmitsDefaults verifies to_dict omits default fields.
func TestParityMarketplaceSourceToDictOmitsDefaults(t *testing.T) {
	s := marketplace.DefaultMarketplaceSource("acme", "acme-org", "tools")
	d := s.ToDict()
	if _, ok := d["host"]; ok {
		t.Error("expected host to be omitted when default")
	}
	if _, ok := d["branch"]; ok {
		t.Error("expected branch to be omitted when default")
	}
	if _, ok := d["path"]; ok {
		t.Error("expected path to be omitted when default")
	}
	if d["name"] != "acme" {
		t.Errorf("expected name=acme, got %s", d["name"])
	}
}

// TestParityMarketplaceSourceToDictIncludesNonDefaults verifies non-default fields appear.
func TestParityMarketplaceSourceToDictIncludesNonDefaults(t *testing.T) {
	s := marketplace.MarketplaceSource{
		Name:   "internal",
		Owner:  "my-org",
		Repo:   "my-repo",
		Host:   "github.enterprise.com",
		Branch: "develop",
		Path:   "custom.json",
	}
	d := s.ToDict()
	if d["host"] != "github.enterprise.com" {
		t.Errorf("expected host in dict, got %v", d["host"])
	}
	if d["branch"] != "develop" {
		t.Errorf("expected branch in dict")
	}
	if d["path"] != "custom.json" {
		t.Errorf("expected path in dict")
	}
}

// TestParityPluginMatchesQueryName verifies query matching on plugin name.
func TestParityPluginMatchesQueryName(t *testing.T) {
	p := marketplace.MarketplacePlugin{
		Name:        "my-tool",
		Description: "A useful tool",
		Tags:        []string{"cli", "automation"},
	}
	if !p.MatchesQuery("my-tool") {
		t.Error("should match exact name")
	}
	if !p.MatchesQuery("MY-TOOL") {
		t.Error("should match case-insensitive name")
	}
}

// TestParityPluginMatchesQueryDescription verifies matching on description.
func TestParityPluginMatchesQueryDescription(t *testing.T) {
	p := marketplace.MarketplacePlugin{
		Name:        "tool",
		Description: "A useful database inspector",
	}
	if !p.MatchesQuery("database") {
		t.Error("should match description substring")
	}
}

// TestParityPluginMatchesQueryTags verifies matching on tags.
func TestParityPluginMatchesQueryTags(t *testing.T) {
	p := marketplace.MarketplacePlugin{
		Name: "tool",
		Tags: []string{"cli", "automation"},
	}
	if !p.MatchesQuery("automation") {
		t.Error("should match tag")
	}
	if !p.MatchesQuery("AUTOMATION") {
		t.Error("should match tag case-insensitive")
	}
}

// TestParityPluginNoMatch verifies non-matching queries return false.
func TestParityPluginNoMatch(t *testing.T) {
	p := marketplace.MarketplacePlugin{
		Name:        "tool",
		Description: "A useful tool",
		Tags:        []string{"cli"},
	}
	if p.MatchesQuery("database") {
		t.Error("should not match unrelated query")
	}
}

// TestParityManifestFindPlugin verifies case-insensitive plugin lookup.
func TestParityManifestFindPlugin(t *testing.T) {
	m := &marketplace.MarketplaceManifest{
		Name: "test",
		Plugins: []marketplace.MarketplacePlugin{
			{Name: "MyPlugin"},
			{Name: "OtherPlugin"},
		},
	}
	p := m.FindPlugin("myplugin")
	if p == nil {
		t.Fatal("expected to find plugin case-insensitively")
	}
	if p.Name != "MyPlugin" {
		t.Errorf("expected MyPlugin, got %s", p.Name)
	}
}

// TestParityManifestFindPluginNotFound verifies nil returned when missing.
func TestParityManifestFindPluginNotFound(t *testing.T) {
	m := &marketplace.MarketplaceManifest{Name: "test"}
	if m.FindPlugin("missing") != nil {
		t.Error("expected nil for missing plugin")
	}
}

// TestParityManifestSearch verifies search returns matching plugins.
func TestParityManifestSearch(t *testing.T) {
	m := &marketplace.MarketplaceManifest{
		Name: "test",
		Plugins: []marketplace.MarketplacePlugin{
			{Name: "dbinspector", Description: "Inspect databases"},
			{Name: "filewatcher", Description: "Watch files"},
			{Name: "dbmigrator", Tags: []string{"database"}},
		},
	}
	results := m.Search("database")
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

// TestParityManifestSearchEmpty verifies empty manifest returns nothing.
func TestParityManifestSearchEmpty(t *testing.T) {
	m := &marketplace.MarketplaceManifest{Name: "empty"}
	results := m.Search("anything")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// TestParityNotFoundError verifies error message format.
func TestParityNotFoundError(t *testing.T) {
	err := &marketplace.NotFoundError{PluginName: "my-plugin"}
	if err.Error() != "plugin not found: my-plugin" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

// TestParityMarketplaceError verifies generic marketplace error.
func TestParityMarketplaceError(t *testing.T) {
	err := &marketplace.MarketplaceError{Message: "connection failed"}
	if err.Error() != "connection failed" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
