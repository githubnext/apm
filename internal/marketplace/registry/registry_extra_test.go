package registry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/registry"
)

func TestFromDictWithExtra(t *testing.T) {
	m := map[string]interface{}{
		"name": "my-src",
		"url":  "https://example.com",
		"note": "extra field",
	}
	src, err := registry.FromDict(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := src.ToDict()
	if d["note"] != "extra field" {
		t.Errorf("extra field not preserved: %v", d)
	}
}

func TestRegistryGetByName(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	r.Add(registry.MarketplaceSource{Name: "alpha", URL: "https://alpha.com"})
	r.Add(registry.MarketplaceSource{Name: "beta", URL: "https://beta.com"})

	src, err := r.GetByName("alpha")
	if err != nil {
		t.Fatalf("GetByName error: %v", err)
	}
	if src.Name != "alpha" {
		t.Errorf("expected alpha, got %q", src.Name)
	}
}

func TestRegistryGetByNameCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	r.Add(registry.MarketplaceSource{Name: "Alpha", URL: "https://alpha.com"})

	src, err := r.GetByName("ALPHA")
	if err != nil {
		t.Fatalf("GetByName case-insensitive error: %v", err)
	}
	if src.Name != "Alpha" {
		t.Errorf("expected Alpha, got %q", src.Name)
	}
}

func TestRegistryGetByNameNotFound(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	_, err := r.GetByName("noexist")
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestRegistryNamesOrder(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	r.Add(registry.MarketplaceSource{Name: "z-last"})
	r.Add(registry.MarketplaceSource{Name: "a-first"})

	names, err := r.Names()
	if err != nil {
		t.Fatalf("Names error: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestRegistryCount(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	count, _ := r.Count()
	if count != 0 {
		t.Errorf("empty registry count = %d, want 0", count)
	}
	r.Add(registry.MarketplaceSource{Name: "x"})
	r.Add(registry.MarketplaceSource{Name: "y"})
	count, _ = r.Count()
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestRegistryAddDuplicate(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	r.Add(registry.MarketplaceSource{Name: "dup", URL: "https://first.com"})
	// Add replaces duplicates (upsert), so no error expected
	err := r.Add(registry.MarketplaceSource{Name: "dup", URL: "https://second.com"})
	if err != nil {
		t.Errorf("unexpected error on duplicate add (should upsert): %v", err)
	}
	// Verify the entry was replaced
	src, _ := r.GetByName("dup")
	if src.URL != "https://second.com" {
		t.Errorf("expected second URL, got %q", src.URL)
	}
}

func TestRegistryRemoveNonExistent(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	err := r.Remove("nonexistent")
	if err == nil {
		t.Error("expected error removing nonexistent source")
	}
}

func TestRegistryGetAll(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	r.Add(registry.MarketplaceSource{Name: "a"})
	r.Add(registry.MarketplaceSource{Name: "b"})
	all, err := r.GetAll()
	if err != nil {
		t.Fatalf("GetAll error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2, got %d", len(all))
	}
}

func TestRegistryFileCreatedOnWrite(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	r.Add(registry.MarketplaceSource{Name: "src"})

	// The file should exist in the config dir
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	var found bool
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			found = true
		}
	}
	if !found {
		t.Error("expected a JSON registry file to be created")
	}
}

func TestFromDictURLOptional(t *testing.T) {
	m := map[string]interface{}{
		"name": "no-url",
	}
	src, err := registry.FromDict(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Name != "no-url" {
		t.Errorf("expected no-url, got %q", src.Name)
	}
	if src.URL != "" {
		t.Errorf("URL should be empty, got %q", src.URL)
	}
}
