// Package mktmodels defines frozen dataclasses and JSON parser for marketplace manifests.
// Ported from src/apm_cli/marketplace/models.py
package mktmodels

import (
	"encoding/json"
	"strings"
)

// MarketplaceSource is a registered marketplace repository.
// Stored in ~/.apm/marketplaces.json.
type MarketplaceSource struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Host   string `json:"host,omitempty"`
	Branch string `json:"branch,omitempty"`
	Path   string `json:"path,omitempty"`
}

// ToDict serializes to a map for JSON storage (omits defaults).
func (m *MarketplaceSource) ToDict() map[string]string {
	result := map[string]string{
		"name":  m.Name,
		"owner": m.Owner,
		"repo":  m.Repo,
	}
	if m.Host != "" && m.Host != "github.com" {
		result["host"] = m.Host
	}
	if m.Branch != "" && m.Branch != "main" {
		result["branch"] = m.Branch
	}
	if m.Path != "" && m.Path != "marketplace.json" {
		result["path"] = m.Path
	}
	return result
}

// NewMarketplaceSource creates a MarketplaceSource with defaults applied.
func NewMarketplaceSource(name, owner, repo, host, branch, path string) MarketplaceSource {
	if host == "" {
		host = "github.com"
	}
	if branch == "" {
		branch = "main"
	}
	if path == "" {
		path = "marketplace.json"
	}
	return MarketplaceSource{Name: name, Owner: owner, Repo: repo, Host: host, Branch: branch, Path: path}
}

// MarketplacePlugin is a single plugin entry inside a marketplace manifest.
type MarketplacePlugin struct {
	Name              string
	Source            interface{} // string or map[string]interface{}
	Description       string
	Version           string
	Tags              []string
	SourceMarketplace string
}

// MatchesQuery returns true if the plugin matches a search query (case-insensitive).
func (p *MarketplacePlugin) MatchesQuery(query string) bool {
	q := strings.ToLower(query)
	if strings.Contains(strings.ToLower(p.Name), q) {
		return true
	}
	if strings.Contains(strings.ToLower(p.Description), q) {
		return true
	}
	for _, tag := range p.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return true
		}
	}
	return false
}

// MarketplaceManifest holds parsed marketplace.json content.
type MarketplaceManifest struct {
	Name        string
	Plugins     []MarketplacePlugin
	OwnerName   string
	Description string
	PluginRoot  string
}

// FindPlugin finds a plugin by exact name (case-insensitive).
func (m *MarketplaceManifest) FindPlugin(name string) *MarketplacePlugin {
	lower := strings.ToLower(name)
	for i := range m.Plugins {
		if strings.ToLower(m.Plugins[i].Name) == lower {
			return &m.Plugins[i]
		}
	}
	return nil
}

// Search returns plugins matching a query.
func (m *MarketplaceManifest) Search(query string) []MarketplacePlugin {
	var result []MarketplacePlugin
	for _, p := range m.Plugins {
		if p.MatchesQuery(query) {
			result = append(result, p)
		}
	}
	return result
}

// parsePluginEntry parses a single plugin entry from either Copilot CLI or Claude Code format.
func parsePluginEntry(entry map[string]interface{}, sourceName string) *MarketplacePlugin {
	name, _ := entry["name"].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	description, _ := entry["description"].(string)
	version, _ := entry["version"].(string)
	var tags []string
	if rawTags, ok := entry["tags"].([]interface{}); ok {
		for _, t := range rawTags {
			if s, ok := t.(string); ok {
				tags = append(tags, s)
			}
		}
	}

	var source interface{}

	if rawSource, ok := entry["source"]; ok {
		switch s := rawSource.(type) {
		case string:
			source = s
		case map[string]interface{}:
			sourceType, _ := s["type"].(string)
			if sourceType == "" {
				sourceType, _ = s["source"].(string)
			}
			if sourceType == "npm" {
				return nil
			}
			if sourceType != "" {
				if _, hasType := s["type"]; !hasType {
					newS := make(map[string]interface{}, len(s)+1)
					for k, v := range s {
						newS[k] = v
					}
					newS["type"] = sourceType
					s = newS
				}
			}
			source = s
		default:
			return nil
		}
	} else if rawRepo, ok := entry["repository"].(string); ok {
		if strings.Contains(rawRepo, "/") {
			src := map[string]interface{}{"type": "github", "repo": rawRepo}
			if ref, ok := entry["ref"].(string); ok && ref != "" {
				src["ref"] = ref
			}
			source = src
		} else {
			return nil
		}
	} else {
		return nil
	}

	return &MarketplacePlugin{
		Name:              name,
		Source:            source,
		Description:       description,
		Version:           version,
		Tags:              tags,
		SourceMarketplace: sourceName,
	}
}

// ParseMarketplaceJSON parses a marketplace.json dict into a MarketplaceManifest.
// Accepts both Copilot CLI and Claude Code marketplace formats.
func ParseMarketplaceJSON(data map[string]interface{}, sourceName string) MarketplaceManifest {
	manifestName, _ := data["name"].(string)
	if manifestName == "" {
		manifestName = sourceName
		if manifestName == "" {
			manifestName = "unknown"
		}
	}
	description, _ := data["description"].(string)

	var ownerName string
	if ownerMap, ok := data["owner"].(map[string]interface{}); ok {
		ownerName, _ = ownerMap["name"].(string)
	} else if ownerStr, ok := data["owner"].(string); ok {
		ownerName = ownerStr
	}

	var pluginRoot string
	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		if pr, ok := metadata["pluginRoot"].(string); ok {
			pluginRoot = strings.TrimSpace(pr)
		}
	}

	var plugins []MarketplacePlugin
	if rawPlugins, ok := data["plugins"].([]interface{}); ok {
		for _, rawEntry := range rawPlugins {
			entry, ok := rawEntry.(map[string]interface{})
			if !ok {
				continue
			}
			p := parsePluginEntry(entry, sourceName)
			if p != nil {
				plugins = append(plugins, *p)
			}
		}
	}

	return MarketplaceManifest{
		Name:        manifestName,
		Plugins:     plugins,
		OwnerName:   ownerName,
		Description: description,
		PluginRoot:  pluginRoot,
	}
}

// ParseMarketplaceJSONBytes parses a marketplace.json byte slice.
func ParseMarketplaceJSONBytes(b []byte, sourceName string) (MarketplaceManifest, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return MarketplaceManifest{}, err
	}
	return ParseMarketplaceJSON(data, sourceName), nil
}
