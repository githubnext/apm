package mktmodels

import (
"strings"
"testing"
)

func TestNewMarketplaceSourceDefaults(t *testing.T) {
s := NewMarketplaceSource("test", "owner", "repo", "", "", "")
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

func TestNewMarketplaceSourceCustom(t *testing.T) {
s := NewMarketplaceSource("test", "owner", "repo", "github.example.com", "develop", "custom.json")
if s.Host != "github.example.com" {
t.Errorf("unexpected host: %s", s.Host)
}
if s.Branch != "develop" {
t.Errorf("unexpected branch: %s", s.Branch)
}
if s.Path != "custom.json" {
t.Errorf("unexpected path: %s", s.Path)
}
}

func TestMarketplaceSourceToDict(t *testing.T) {
s := NewMarketplaceSource("test", "owner", "repo", "", "", "")
d := s.ToDict()
if d["name"] != "test" || d["owner"] != "owner" || d["repo"] != "repo" {
t.Errorf("unexpected dict: %v", d)
}
// Default values should not be included
if _, ok := d["host"]; ok {
t.Error("default host should be omitted from dict")
}
if _, ok := d["branch"]; ok {
t.Error("default branch should be omitted from dict")
}
}

func TestMarketplaceSourceToDictNonDefault(t *testing.T) {
s := NewMarketplaceSource("test", "owner", "repo", "enterprise.com", "dev", "custom.json")
d := s.ToDict()
if d["host"] != "enterprise.com" {
t.Errorf("expected host in dict: %v", d)
}
}

func TestParseMarketplaceJSONBytes(t *testing.T) {
data := []byte(`{"plugins":[{"name":"my-pkg","repository":"acme/pkg","description":"A package","tags":["ai"]}]}`)
manifest, err := ParseMarketplaceJSONBytes(data, "test-source")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(manifest.Plugins) == 0 {
t.Error("expected at least one package")
}
}

func TestMarketplaceManifestFindPlugin(t *testing.T) {
data := []byte(`{"plugins":[{"name":"my-pkg","repository":"acme/pkg","description":"A package"}]}`)
manifest, err := ParseMarketplaceJSONBytes(data, "src")
if err != nil {
t.Fatalf("parse error: %v", err)
}
p := manifest.FindPlugin("my-pkg")
if p == nil {
t.Fatal("expected to find plugin my-pkg")
}
if p.Name != "my-pkg" {
t.Errorf("unexpected name: %s", p.Name)
}
}

func TestMarketplaceManifestFindPluginMissing(t *testing.T) {
manifest := MarketplaceManifest{}
if manifest.FindPlugin("nonexistent") != nil {
t.Error("expected nil for missing plugin")
}
}

func TestMarketplaceManifestSearch(t *testing.T) {
data := []byte(`{"plugins":[
{"name":"alpha","repository":"o/r","description":"alpha tool"},
{"name":"beta","repository":"o/r2","description":"beta tool"},
{"name":"gamma","repository":"o/r3","description":"something else"}
]}`)
manifest, err := ParseMarketplaceJSONBytes(data, "src")
if err != nil {
t.Fatalf("parse error: %v", err)
}
results := manifest.Search("alpha")
if len(results) == 0 {
t.Error("expected at least one result for 'alpha'")
}
for _, r := range results {
if !strings.Contains(strings.ToLower(r.Name), "alpha") &&
!strings.Contains(strings.ToLower(r.Description), "alpha") {
t.Errorf("result %s doesn't match query 'alpha'", r.Name)
}
}
}
