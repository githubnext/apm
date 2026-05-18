package shadowdetector_test

import (
	"fmt"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/shadowdetector"
)

type fakeLister struct {
	marketplaces []string
	plugins      map[string][]string
	errs         map[string]error
}

func (f *fakeLister) ListRegisteredMarketplaces() []string { return f.marketplaces }
func (f *fakeLister) ListPluginNames(mp string) ([]string, error) {
	if e, ok := f.errs[mp]; ok {
		return nil, e
	}
	return f.plugins[mp], nil
}

func TestDetectShadows_ReturnsShadowMatchFields(t *testing.T) {
	lister := &fakeLister{
		marketplaces: []string{"primary", "secondary"},
		plugins: map[string][]string{
			"primary":   {"myplugin"},
			"secondary": {"MyPlugin"},
		},
	}
	matches := shadowdetector.DetectShadows("myplugin", "primary", lister)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].MarketplaceName != "secondary" {
		t.Errorf("MarketplaceName = %q, want %q", matches[0].MarketplaceName, "secondary")
	}
	if matches[0].PluginName != "MyPlugin" {
		t.Errorf("PluginName = %q, want %q", matches[0].PluginName, "MyPlugin")
	}
}

func TestDetectShadows_ErrorMarketplaceSkipped(t *testing.T) {
	lister := &fakeLister{
		marketplaces: []string{"primary", "bad", "other"},
		plugins: map[string][]string{
			"primary": {"plugin"},
			"other":   {"plugin"},
		},
		errs: map[string]error{"bad": fmt.Errorf("fetch error")},
	}
	matches := shadowdetector.DetectShadows("plugin", "primary", lister)
	// "bad" is skipped, "other" should match
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].MarketplaceName != "other" {
		t.Errorf("unexpected marketplace: %q", matches[0].MarketplaceName)
	}
}

func TestDetectShadows_NoMatch(t *testing.T) {
	lister := &fakeLister{
		marketplaces: []string{"primary", "secondary"},
		plugins: map[string][]string{
			"primary":   {"plugin-a"},
			"secondary": {"plugin-b"},
		},
	}
	matches := shadowdetector.DetectShadows("plugin-a", "primary", lister)
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestDetectShadows_MultipleSecondaryMarketplaces(t *testing.T) {
	lister := &fakeLister{
		marketplaces: []string{"primary", "mp1", "mp2", "mp3"},
		plugins: map[string][]string{
			"primary": {"tool"},
			"mp1":     {"other"},
			"mp2":     {"TOOL"},
			"mp3":     {"tool"},
		},
	}
	matches := shadowdetector.DetectShadows("tool", "primary", lister)
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestDetectShadows_OnlyOneMatchPerMarketplace(t *testing.T) {
	// Even if the plugin name appears twice in one marketplace, we only break on first
	lister := &fakeLister{
		marketplaces: []string{"primary", "secondary"},
		plugins: map[string][]string{
			"primary":   {"tool"},
			"secondary": {"TOOL", "tool"},
		},
	}
	matches := shadowdetector.DetectShadows("tool", "primary", lister)
	if len(matches) != 1 {
		t.Errorf("expected 1 match (break on first), got %d", len(matches))
	}
}

func TestDetectShadows_EmptyPluginName(t *testing.T) {
	lister := &fakeLister{
		marketplaces: []string{"primary", "secondary"},
		plugins: map[string][]string{
			"primary":   {""},
			"secondary": {""},
		},
	}
	matches := shadowdetector.DetectShadows("", "primary", lister)
	if len(matches) != 1 {
		t.Errorf("expected 1 match for empty plugin name, got %d", len(matches))
	}
}

func TestShadowMatch_ZeroValue(t *testing.T) {
	var m shadowdetector.ShadowMatch
	if m.MarketplaceName != "" || m.PluginName != "" {
		t.Error("zero value should have empty fields")
	}
}
