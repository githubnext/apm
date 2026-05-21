package shadowdetector

import "testing"

type mockListerExtra2 struct {
	marketplaces []string
	plugins      map[string][]string
	errors       map[string]error
}

func (m *mockListerExtra2) ListRegisteredMarketplaces() []string { return m.marketplaces }
func (m *mockListerExtra2) ListPluginNames(mp string) ([]string, error) {
	if err, ok := m.errors[mp]; ok {
		return nil, err
	}
	return m.plugins[mp], nil
}

func TestShadowMatch_Fields(t *testing.T) {
	sm := ShadowMatch{MarketplaceName: "npm", PluginName: "foo"}
	if sm.MarketplaceName != "npm" {
		t.Errorf("expected npm, got %q", sm.MarketplaceName)
	}
	if sm.PluginName != "foo" {
		t.Errorf("expected foo, got %q", sm.PluginName)
	}
}

func TestDetectShadows_NilLister(t *testing.T) {
	result := DetectShadows("plugin", "primary", nil)
	if len(result) != 0 {
		t.Errorf("expected empty result with nil lister, got %d", len(result))
	}
}

func TestDetectShadows_PrimaryExcluded(t *testing.T) {
	lister := &mockListerExtra2{
		marketplaces: []string{"primary"},
		plugins:      map[string][]string{"primary": {"myplugin"}},
	}
	result := DetectShadows("myplugin", "primary", lister)
	if len(result) != 0 {
		t.Errorf("primary marketplace should be excluded, got %d results", len(result))
	}
}

func TestDetectShadows_CaseInsensitiveMatch(t *testing.T) {
	lister := &mockListerExtra2{
		marketplaces: []string{"primary", "secondary"},
		plugins: map[string][]string{
			"secondary": {"MYPLUGIN"},
		},
	}
	result := DetectShadows("myplugin", "primary", lister)
	if len(result) != 1 {
		t.Errorf("expected 1 shadow match, got %d", len(result))
	}
}

func TestDetectShadows_NoSecondaryMarketplaces(t *testing.T) {
	lister := &mockListerExtra2{
		marketplaces: []string{"primary"},
		plugins:      map[string][]string{},
	}
	result := DetectShadows("plugin", "primary", lister)
	if len(result) != 0 {
		t.Errorf("expected no results, got %d", len(result))
	}
}

func TestDetectShadows_MultipleMatchesAcrossMarketplaces(t *testing.T) {
	lister := &mockListerExtra2{
		marketplaces: []string{"primary", "mp1", "mp2"},
		plugins: map[string][]string{
			"mp1": {"plugin-a", "plugin-b"},
			"mp2": {"plugin-b", "plugin-c"},
		},
	}
	result := DetectShadows("plugin-b", "primary", lister)
	if len(result) != 2 {
		t.Errorf("expected 2 matches, got %d", len(result))
	}
}
