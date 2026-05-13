// Package shadowdetector detects cross-marketplace plugin name shadowing.
package shadowdetector

import "strings"

// ShadowMatch represents a plugin name found in a secondary marketplace.
type ShadowMatch struct {
MarketplaceName string
PluginName      string
}

// MarketplaceLister is an interface for listing plugins in a marketplace.
type MarketplaceLister interface {
ListPluginNames(marketplace string) ([]string, error)
ListRegisteredMarketplaces() []string
}

// DetectShadows checks registered marketplaces for duplicate plugin names.
func DetectShadows(pluginName, primaryMarketplace string, lister MarketplaceLister) []ShadowMatch {
var results []ShadowMatch
if lister == nil {
return results
}
for _, mp := range lister.ListRegisteredMarketplaces() {
if mp == primaryMarketplace {
continue
}
names, err := lister.ListPluginNames(mp)
if err != nil {
continue
}
lower := strings.ToLower(pluginName)
for _, n := range names {
if strings.ToLower(n) == lower {
results = append(results, ShadowMatch{MarketplaceName: mp, PluginName: n})
break
}
}
}
return results
}
