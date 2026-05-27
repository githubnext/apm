// Package marketplace provides types and logic for APM marketplace plugin management.
// Mirrors Python apm_cli.marketplace.models and apm_cli.marketplace.client.
package marketplace

// MarketplaceSource represents a registered marketplace repository.
// Mirrors Python MarketplaceSource dataclass.
type MarketplaceSource struct {
	Name   string
	Owner  string
	Repo   string
	Host   string
	Branch string
	Path   string
}

// DefaultMarketplaceSource returns a MarketplaceSource with default field values.
func DefaultMarketplaceSource(name, owner, repo string) MarketplaceSource {
	return MarketplaceSource{
		Name:   name,
		Owner:  owner,
		Repo:   repo,
		Host:   "github.com",
		Branch: "main",
		Path:   "marketplace.json",
	}
}

// ToDict serializes a MarketplaceSource to a map, omitting default-valued fields.
// Mirrors Python MarketplaceSource.to_dict().
func (m MarketplaceSource) ToDict() map[string]string {
	result := map[string]string{
		"name":  m.Name,
		"owner": m.Owner,
		"repo":  m.Repo,
	}
	if m.Host != "github.com" {
		result["host"] = m.Host
	}
	if m.Branch != "main" {
		result["branch"] = m.Branch
	}
	if m.Path != "marketplace.json" {
		result["path"] = m.Path
	}
	return result
}

// MarketplacePlugin represents a single plugin entry inside a marketplace manifest.
// Mirrors Python MarketplacePlugin dataclass.
type MarketplacePlugin struct {
	Name                string
	Source              interface{} // string (relative) or map (github/url/git-subdir)
	Description         string
	Version             string
	Tags                []string
	SourceMarketplace   string
}

// MatchesQuery returns true if the plugin matches a search query (case-insensitive).
// Mirrors Python MarketplacePlugin.matches_query().
func (p MarketplacePlugin) MatchesQuery(query string) bool {
	q := toLower(query)
	if contains(toLower(p.Name), q) {
		return true
	}
	if contains(toLower(p.Description), q) {
		return true
	}
	for _, tag := range p.Tags {
		if contains(toLower(tag), q) {
			return true
		}
	}
	return false
}

// MarketplaceManifest is the parsed marketplace.json content.
// Mirrors Python MarketplaceManifest dataclass.
type MarketplaceManifest struct {
	Name        string
	Plugins     []MarketplacePlugin
	OwnerName   string
	Description string
	PluginRoot  string
}

// FindPlugin finds a plugin by exact name (case-insensitive).
// Mirrors Python MarketplaceManifest.find_plugin().
func (m *MarketplaceManifest) FindPlugin(name string) *MarketplacePlugin {
	lower := toLower(name)
	for i := range m.Plugins {
		if toLower(m.Plugins[i].Name) == lower {
			return &m.Plugins[i]
		}
	}
	return nil
}

// Search returns plugins matching a query.
// Mirrors Python MarketplaceManifest.search().
func (m *MarketplaceManifest) Search(query string) []MarketplacePlugin {
	var out []MarketplacePlugin
	for _, p := range m.Plugins {
		if p.MatchesQuery(query) {
			out = append(out, p)
		}
	}
	return out
}

// MarketplaceError is returned for marketplace client errors.
type MarketplaceError struct {
	Message string
}

func (e *MarketplaceError) Error() string { return e.Message }

// NotFoundError is returned when a plugin is not found.
type NotFoundError struct {
	PluginName string
}

func (e *NotFoundError) Error() string {
	return "plugin not found: " + e.PluginName
}

// toLower is an ASCII-safe lowercase helper.
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
