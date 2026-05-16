package registry_test

import (
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/marketplace/registry"
)

func TestFromDictValid(t *testing.T) {
m := map[string]interface{}{"name": "my-market", "url": "https://example.com"}
s, err := registry.FromDict(m)
if err != nil {
t.Fatalf("FromDict: %v", err)
}
if s.Name != "my-market" {
t.Errorf("Name: want my-market, got %s", s.Name)
}
if s.URL != "https://example.com" {
t.Errorf("URL: want https://example.com, got %s", s.URL)
}
}

func TestFromDictMissingName(t *testing.T) {
m := map[string]interface{}{"url": "https://example.com"}
_, err := registry.FromDict(m)
if err == nil {
t.Error("FromDict: expected error for missing name")
}
}

func TestFromDictEmptyName(t *testing.T) {
m := map[string]interface{}{"name": "", "url": "https://example.com"}
_, err := registry.FromDict(m)
if err == nil {
t.Error("FromDict: expected error for empty name")
}
}

func TestToDict(t *testing.T) {
s := registry.MarketplaceSource{Name: "foo", URL: "https://foo.com"}
d := s.ToDict()
if d["name"] != "foo" {
t.Errorf("ToDict name: want foo, got %v", d["name"])
}
if d["url"] != "https://foo.com" {
t.Errorf("ToDict url: want https://foo.com, got %v", d["url"])
}
}

func configDir(dir string) func() string { return func() string { return dir } }

func TestRegistryAddAndList(t *testing.T) {
dir := t.TempDir()
r := registry.New(configDir(dir))
sources, err := r.GetAll()
if err != nil {
t.Fatalf("List empty: %v", err)
}
if len(sources) != 0 {
t.Errorf("expected empty list, got %v", sources)
}

err = r.Add(registry.MarketplaceSource{Name: "alpha", URL: "https://alpha.com"})
if err != nil {
t.Fatalf("Add: %v", err)
}

sources, err = r.GetAll()
if err != nil {
t.Fatalf("List after Add: %v", err)
}
if len(sources) != 1 || sources[0].Name != "alpha" {
t.Errorf("List: want [{alpha ...}], got %v", sources)
}
}

func TestRegistryRemove(t *testing.T) {
dir := t.TempDir()
r := registry.New(configDir(dir))
_ = r.Add(registry.MarketplaceSource{Name: "beta", URL: "https://beta.com"})

err := r.Remove("beta")
if err != nil {
t.Fatalf("Remove: %v", err)
}
sources, _ := r.GetAll()
if len(sources) != 0 {
t.Errorf("expected empty after remove, got %v", sources)
}
}

func TestRegistryPersistence(t *testing.T) {
dir := t.TempDir()
r1 := registry.New(configDir(dir))
_ = r1.Add(registry.MarketplaceSource{Name: "gamma", URL: "https://gamma.com"})

r2 := registry.New(configDir(dir))
sources, err := r2.GetAll()
if err != nil {
t.Fatalf("List r2: %v", err)
}
if len(sources) != 1 || sources[0].Name != "gamma" {
t.Errorf("persistence: want gamma, got %v", sources)
}
}

func TestRegistryFileCreated(t *testing.T) {
dir := t.TempDir()
r := registry.New(configDir(dir))
_ = r.Add(registry.MarketplaceSource{Name: "delta", URL: "https://delta.com"})
if _, err := os.Stat(filepath.Join(dir, "marketplaces.json")); err != nil {
t.Errorf("marketplaces.json not created: %v", err)
}
}
