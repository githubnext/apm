package registry_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/registry"
)

func TestFromDictToDict_RoundTrip(t *testing.T) {
	m := map[string]interface{}{
		"name": "test-market",
		"url":  "https://github.com/org/marketplace",
	}
	src, err := registry.FromDict(m)
	if err != nil {
		t.Fatalf("FromDict error: %v", err)
	}
	out := src.ToDict()
	if out["name"] != "test-market" {
		t.Errorf("name mismatch: %v", out["name"])
	}
	if out["url"] != "https://github.com/org/marketplace" {
		t.Errorf("url mismatch: %v", out["url"])
	}
}

func TestFromDict_EmptyURL(t *testing.T) {
	m := map[string]interface{}{"name": "nourl"}
	src, err := registry.FromDict(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Name != "nourl" {
		t.Errorf("name mismatch: %q", src.Name)
	}
	if src.URL != "" {
		t.Errorf("expected empty URL, got %q", src.URL)
	}
}

func TestFromDict_ExtraFieldsPreserved(t *testing.T) {
	m := map[string]interface{}{
		"name":    "mkt",
		"url":     "https://example.com",
		"custom":  "val",
	}
	src, err := registry.FromDict(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Extra["custom"] != "val" {
		t.Errorf("expected extra field 'custom'='val', got %v", src.Extra["custom"])
	}
}

func TestFromDict_MissingNameError(t *testing.T) {
	m := map[string]interface{}{"url": "https://example.com"}
	_, err := registry.FromDict(m)
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestRegistry_AddAndGetByName(t *testing.T) {
	dir := t.TempDir()
	r := registry.New(func() string { return dir })
	src := registry.MarketplaceSource{Name: "mkt", URL: "https://example.com"}
	if err := r.Add(src); err != nil {
		t.Fatalf("Add error: %v", err)
	}
	got, err := r.GetByName("mkt")
	if err != nil {
		t.Fatalf("GetByName error: %v", err)
	}
	if got.Name != "mkt" {
		t.Errorf("got %q, want %q", got.Name, "mkt")
	}
}
