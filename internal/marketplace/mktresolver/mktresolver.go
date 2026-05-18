// Package mktresolver resolves marketplace plugin specifiers to canonical refs.
// Migrated from src/apm_cli/marketplace/resolver.py
package mktresolver

import (
	"fmt"
	"regexp"
	"strings"
)

// marketplaceRE matches NAME@MARKETPLACE[#ref] specifiers.
var marketplaceRE = regexp.MustCompile(`^([a-zA-Z0-9._-]+)@([a-zA-Z0-9._-]+)(?:#(.+))?$`)

// semverRangeCharsRE matches characters that indicate a semver range.
var semverRangeCharsRE = regexp.MustCompile(`[~^<>=!]`)

// MarketplacePlugin represents a plugin entry from a marketplace registry.
type MarketplacePlugin struct {
	Name   string
	Repo   string
	Source map[string]interface{}
}

// MarketplaceSource represents a marketplace source configuration.
type MarketplaceSource struct {
	Name string
	Host string
	Repo string
}

// MarketplacePluginResolution is the outcome of ResolveMarketplacePlugin.
type MarketplacePluginResolution struct {
	// Canonical is the resolved owner/repo#ref string.
	Canonical string
	// Plugin is the matched marketplace plugin entry.
	Plugin MarketplacePlugin
	// DependencyReference is set when a structured ref is needed
	// (e.g. GitLab in-marketplace subdirectory plugins).
	DependencyReference interface{}
}

// ParsedMarketplaceRef holds the parsed parts of a NAME@MARKETPLACE[#ref] string.
type ParsedMarketplaceRef struct {
	Name        string
	Marketplace string
	Ref         string
}

// ParseMarketplaceRef parses a NAME@MARKETPLACE[#ref] specifier.
// Returns nil when the input does not match the pattern.
func ParseMarketplaceRef(spec string) *ParsedMarketplaceRef {
	m := marketplaceRE.FindStringSubmatch(spec)
	if m == nil {
		return nil
	}
	return &ParsedMarketplaceRef{Name: m[1], Marketplace: m[2], Ref: m[3]}
}

// IsMarketplaceRef reports whether spec looks like a marketplace ref.
func IsMarketplaceRef(spec string) bool {
	return marketplaceRE.MatchString(spec)
}

// IsSemverRange reports whether ref contains semver range characters.
func IsSemverRange(ref string) bool {
	return semverRangeCharsRE.MatchString(ref)
}

// NormalizeOwnerRepoSlug lowercases an owner/repo slug and strips .git.
func NormalizeOwnerRepoSlug(repo string) string {
	r := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(repo), "/"))
	r = strings.TrimSuffix(r, ".git")
	return strings.ToLower(r)
}

// MarketplaceProjectSlug returns the normalized owner/repo slug.
func MarketplaceProjectSlug(owner, repo string) string {
	return NormalizeOwnerRepoSlug(owner + "/" + repo)
}

// NormalizeRepoFieldForMatch normalizes a repo field to a logical project path
// for marketplace matching. Returns "" when the field names a different host.
func NormalizeRepoFieldForMatch(repoField, marketplaceHost string) string {
	raw := strings.TrimSuffix(strings.TrimRight(strings.TrimSpace(repoField), "/"), ".git")
	hostL := strings.ToLower(strings.TrimSpace(marketplaceHost))

	// Handle full URLs
	for _, prefix := range []string{"https://", "http://", "ssh://"} {
		if strings.HasPrefix(raw, prefix) {
			rest := strings.TrimPrefix(raw, prefix)
			// Strip scheme-specific prefix for ssh://
			slash := strings.Index(rest, "/")
			if slash < 0 {
				return ""
			}
			host := strings.ToLower(rest[:slash])
			if host != hostL {
				return ""
			}
			return strings.ToLower(strings.TrimLeft(rest[slash:], "/"))
		}
	}

	// git@ SSH shorthand
	if strings.Contains(raw, "@") && strings.Contains(raw, ":") {
		atIdx := strings.Index(raw, "@")
		colonIdx := strings.Index(raw, ":")
		if atIdx < colonIdx {
			host := strings.ToLower(raw[atIdx+1 : colonIdx])
			if host != hostL {
				return ""
			}
			return strings.ToLower(raw[colonIdx+1:])
		}
	}

	// Bare host/owner/repo
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, hostL+"/") {
		return strings.TrimPrefix(lower, hostL+"/")
	}
	return lower
}

// RepoFieldMatchesMarketplace reports whether repoField belongs to the marketplace source.
func RepoFieldMatchesMarketplace(repoField string, source MarketplaceSource) bool {
	slug := MarketplaceProjectSlug(source.Repo, "")
	normalized := NormalizeRepoFieldForMatch(repoField, source.Host)
	if normalized == "" {
		return false
	}
	return strings.HasSuffix(strings.TrimRight(normalized, "/"), strings.TrimRight(slug, "/"))
}

// GitSourceToCanonical converts a GitHub source dict to a canonical ref string.
func GitSourceToCanonical(source map[string]interface{}) string {
	repo, _ := source["repo"].(string)
	ref, _ := source["ref"].(string)
	if ref == "" {
		ref, _ = source["version"].(string)
	}
	if ref != "" {
		return fmt.Sprintf("%s#%s", NormalizeOwnerRepoSlug(repo), ref)
	}
	return NormalizeOwnerRepoSlug(repo)
}

// URLSourceToCanonical converts a URL source dict to a canonical URL ref.
func URLSourceToCanonical(source map[string]interface{}) string {
	url, _ := source["url"].(string)
	return strings.TrimSpace(url)
}

// PluginSourceType identifies the type of a plugin source.
type PluginSourceType int

const (
	PluginSourceGitHub PluginSourceType = iota
	PluginSourceURL
	PluginSourceGitSubdir
	PluginSourceRelative
	PluginSourceUnknown
)

// ClassifyPluginSource determines the type of a plugin source map.
func ClassifyPluginSource(source map[string]interface{}) PluginSourceType {
	if _, ok := source["github"]; ok {
		return PluginSourceGitHub
	}
	if _, ok := source["url"]; ok {
		return PluginSourceURL
	}
	if git, ok := source["git"].(string); ok && git != "" {
		if _, ok2 := source["path"]; ok2 {
			return PluginSourceGitSubdir
		}
	}
	if _, ok := source["relative"]; ok {
		return PluginSourceRelative
	}
	return PluginSourceUnknown
}

// ResolvePluginSource resolves a plugin source map to a canonical string.
// Returns ("", ErrUnknownSourceType) for unrecognised source formats.
func ResolvePluginSource(source map[string]interface{}, marketplaceHost string) (string, error) {
	switch ClassifyPluginSource(source) {
	case PluginSourceGitHub:
		gh, _ := source["github"].(map[string]interface{})
		if gh == nil {
			return "", fmt.Errorf("malformed github source")
		}
		return GitSourceToCanonical(gh), nil
	case PluginSourceURL:
		return URLSourceToCanonical(source), nil
	case PluginSourceGitSubdir:
		git, _ := source["git"].(string)
		path, _ := source["path"].(string)
		ref, _ := source["ref"].(string)
		canonical := NormalizeOwnerRepoSlug(git)
		if path != "" {
			canonical += " path:" + path
		}
		if ref != "" {
			canonical += "#" + ref
		}
		return canonical, nil
	case PluginSourceRelative:
		rel, _ := source["relative"].(string)
		return "relative:" + rel, nil
	default:
		return "", fmt.Errorf("unknown plugin source type")
	}
}

// MarketplaceHostNeedsExplicitGitPath reports whether the marketplace host
// (typically a GitLab instance) requires an explicit git URL + path for
// in-marketplace subdirectory plugins.
func MarketplaceHostNeedsExplicitGitPath(host string) bool {
	h := strings.ToLower(host)
	if h == "github.com" {
		return false
	}
	if strings.HasSuffix(h, ".ghe.com") {
		return false
	}
	return true
}
