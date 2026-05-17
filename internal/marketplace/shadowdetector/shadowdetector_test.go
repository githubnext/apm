package shadowdetector_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/shadowdetector"
)

type mockLister struct {
	plugins       map[string][]string
	marketplaces  []string
}

func (m *mockLister) ListPluginNames(marketplace string) ([]string, error) {
	return m.plugins[marketplace], nil
}

func (m *mockLister) ListRegisteredMarketplaces() []string {
	return m.marketplaces
}

func TestDetectShadows_NoConflict(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{"secondary": {"other-plugin"}},
		marketplaces: []string{"primary", "secondary"},
	}
	results := shadowdetector.DetectShadows("my-plugin", "primary", lister)
	if len(results) != 0 {
		t.Errorf("expected no shadows, got %v", results)
	}
}

func TestDetectShadows_Conflict(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{"secondary": {"my-plugin", "other"}},
		marketplaces: []string{"primary", "secondary"},
	}
	results := shadowdetector.DetectShadows("my-plugin", "primary", lister)
	if len(results) != 1 {
		t.Fatalf("expected 1 shadow, got %d", len(results))
	}
	if results[0].MarketplaceName != "secondary" {
		t.Errorf("expected secondary, got %q", results[0].MarketplaceName)
	}
}

func TestDetectShadows_CaseInsensitive(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{"other": {"MY-PLUGIN"}},
		marketplaces: []string{"primary", "other"},
	}
	results := shadowdetector.DetectShadows("my-plugin", "primary", lister)
	if len(results) != 1 {
		t.Fatalf("expected 1 shadow, got %d", len(results))
	}
	if results[0].PluginName != "MY-PLUGIN" {
		t.Errorf("expected 'MY-PLUGIN', got %q", results[0].PluginName)
	}
}

func TestDetectShadows_SkipsPrimary(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{"primary": {"my-plugin"}},
		marketplaces: []string{"primary"},
	}
	results := shadowdetector.DetectShadows("my-plugin", "primary", lister)
	if len(results) != 0 {
		t.Errorf("should not detect shadow in primary marketplace itself")
	}
}

func TestDetectShadows_NilLister(t *testing.T) {
	results := shadowdetector.DetectShadows("x", "y", nil)
	if len(results) != 0 {
		t.Error("nil lister should return empty slice")
	}
}

func TestDetectShadows_MultipleConflicts(t *testing.T) {
	lister := &mockLister{
		plugins: map[string][]string{
			"mp-a": {"my-plugin"},
			"mp-b": {"MY-PLUGIN"},
		},
		marketplaces: []string{"primary", "mp-a", "mp-b"},
	}
	results := shadowdetector.DetectShadows("my-plugin", "primary", lister)
	if len(results) != 2 {
		t.Fatalf("expected 2 shadows, got %d", len(results))
	}
}

func TestDetectShadows_EmptyMarketplaces(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{},
		marketplaces: []string{},
	}
	results := shadowdetector.DetectShadows("any-plugin", "primary", lister)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestDetectShadows_OnlyPrimary(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{"primary": {"my-plugin"}},
		marketplaces: []string{"primary"},
	}
	results := shadowdetector.DetectShadows("my-plugin", "primary", lister)
	if len(results) != 0 {
		t.Error("primary marketplace should not be checked for shadows")
	}
}

func TestShadowMatchFields(t *testing.T) {
	lister := &mockLister{
		plugins:      map[string][]string{"other": {"TargetPlugin"}},
		marketplaces: []string{"main", "other"},
	}
	results := shadowdetector.DetectShadows("targetplugin", "main", lister)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].MarketplaceName != "other" {
		t.Errorf("MarketplaceName: got %q, want %q", results[0].MarketplaceName, "other")
	}
	if results[0].PluginName != "TargetPlugin" {
		t.Errorf("PluginName: got %q, want %q", results[0].PluginName, "TargetPlugin")
	}
}
