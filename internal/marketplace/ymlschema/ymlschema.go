// Package ymlschema provides dataclasses, loader, and validation for
// marketplace authoring config.
// Migrated from src/apm_cli/marketplace/yml_schema.py.
package ymlschema

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Errors

// MarketplaceYmlError is raised on marketplace YAML validation failures.
type MarketplaceYmlError struct {
	Msg string
}

func (e *MarketplaceYmlError) Error() string { return e.Msg }

func mErr(format string, args ...interface{}) *MarketplaceYmlError {
	return &MarketplaceYmlError{Msg: fmt.Sprintf(format, args...)}
}

// Regex patterns
var (
	semverRE      = regexp.MustCompile(`^\d+\.\d+\.\d+(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$`)
	sourceRE      = regexp.MustCompile(`^(?:[^/]+/[^/]+|\./.*)$`)
	localSourceRE = regexp.MustCompile(`^\.\/`)
)

const (
	maxTagsCount  = 50
	maxTagLength  = 100
)

var tagPlaceholders = []string{"{version}", "{name}"}

// MarketplaceOwner is the owner block of marketplace.yml.
type MarketplaceOwner struct {
	Name  string
	Email string
	URL   string
}

// MarketplaceBuild is the APM-only build configuration block.
type MarketplaceBuild struct {
	TagPattern string
}

// PackageEntry is a single entry in the packages list.
type PackageEntry struct {
	Name              string
	Source            string
	Subdir            string
	Version           string
	Ref               string
	TagPattern        string
	IncludePrerelease bool
	Description       string
	Homepage          string
	Tags              []string
	Author            map[string]string // {name, email?, url?}
	License           string
	Repository        string
	IsLocal           bool
}

// MarketplaceConfig is the parsed marketplace configuration.
type MarketplaceConfig struct {
	Name                string
	Description         string
	Version             string
	Owner               MarketplaceOwner
	Output              string
	Metadata            map[string]interface{}
	Build               MarketplaceBuild
	Packages            []PackageEntry
	SourcePath          string
	IsLegacy            bool
	NameOverridden      bool
	DescriptionOverridden bool
	VersionOverridden   bool
}

// parseSimpleYAML is a minimal line-by-line YAML parser for flat string values.
// Returns top-level key->value pairs (no nesting). Values are trimmed and unquoted.
func parseSimpleYAML(content string) map[string]string {
	result := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, "\"'")
		if val != "" && !strings.HasPrefix(val, "{") && !strings.HasPrefix(val, "[") && !strings.HasPrefix(val, "-") {
			result[key] = val
		}
	}
	return result
}

func validateSemver(version, context string) error {
	if !semverRE.MatchString(version) {
		return mErr("'%s' value '%s' is not valid semver (expected x.y.z)", context, version)
	}
	return nil
}

func validateTagPattern(pattern, context string) error {
	for _, ph := range tagPlaceholders {
		if strings.Contains(pattern, ph) {
			return nil
		}
	}
	return mErr("'%s' must contain at least one of %s, got '%s'", context, strings.Join(tagPlaceholders, ", "), pattern)
}

func validatePathSegments(path string) error {
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if p == ".." {
			return fmt.Errorf("path traversal detected in: %s", path)
		}
	}
	return nil
}

func parseOwner(raw map[string]interface{}) (MarketplaceOwner, error) {
	name, ok := raw["name"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return MarketplaceOwner{}, mErr("'owner.name' is required and must be a non-empty string")
	}
	owner := MarketplaceOwner{Name: strings.TrimSpace(name)}
	if email, ok := raw["email"].(string); ok {
		owner.Email = strings.TrimSpace(email)
	}
	if url, ok := raw["url"].(string); ok {
		owner.URL = strings.TrimSpace(url)
	}
	return owner, nil
}

func parseBuild(raw interface{}) (MarketplaceBuild, error) {
	if raw == nil {
		return MarketplaceBuild{TagPattern: "v{version}"}, nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return MarketplaceBuild{}, mErr("'build' must be a mapping")
	}
	tagPattern := "v{version}"
	if tp, ok := m["tagPattern"].(string); ok && strings.TrimSpace(tp) != "" {
		tagPattern = strings.TrimSpace(tp)
	}
	if err := validateTagPattern(tagPattern, "build.tagPattern"); err != nil {
		return MarketplaceBuild{}, err
	}
	return MarketplaceBuild{TagPattern: tagPattern}, nil
}

func getStr(m map[string]interface{}, key string) (string, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func requireStr(m map[string]interface{}, key, context string) (string, error) {
	s, ok := getStr(m, key)
	if !ok || strings.TrimSpace(s) == "" {
		path := key
		if context != "" {
			path = context + "." + key
		}
		return "", mErr("'%s' is required", path)
	}
	return strings.TrimSpace(s), nil
}

func checkUnknownKeys(data map[string]interface{}, permitted map[string]bool, context string) error {
	var unknown []string
	for k := range data {
		if !permitted[k] {
			unknown = append(unknown, k)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		var perm []string
		for k := range permitted {
			perm = append(perm, k)
		}
		sort.Strings(perm)
		return mErr("Unknown key(s) in %s: %s. Permitted keys: %s", context, strings.Join(unknown, ", "), strings.Join(perm, ", "))
	}
	return nil
}

var packageEntryKeys = map[string]bool{
	"name": true, "source": true, "subdir": true, "version": true, "ref": true,
	"tag_pattern": true, "include_prerelease": true, "description": true,
	"homepage": true, "tags": true, "author": true, "license": true,
	"repository": true, "keywords": true,
}

var apmMarketplaceKeys = map[string]bool{
	"name": true, "description": true, "version": true, "owner": true,
	"output": true, "metadata": true, "build": true, "packages": true,
}

func parsePackageEntry(raw interface{}, index int) (PackageEntry, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		// Try map[interface{}]interface{} (some YAML parsers)
		if mi, ok2 := raw.(map[interface{}]interface{}); ok2 {
			m = make(map[string]interface{})
			for k, v := range mi {
				m[fmt.Sprint(k)] = v
			}
		} else {
			return PackageEntry{}, mErr("packages[%d] must be a mapping", index)
		}
	}
	if err := checkUnknownKeys(m, packageEntryKeys, fmt.Sprintf("packages[%d]", index)); err != nil {
		return PackageEntry{}, err
	}

	name, err := requireStr(m, "name", fmt.Sprintf("packages[%d]", index))
	if err != nil {
		return PackageEntry{}, err
	}
	source, err := requireStr(m, "source", fmt.Sprintf("packages[%d]", index))
	if err != nil {
		return PackageEntry{}, err
	}
	if !sourceRE.MatchString(source) {
		return PackageEntry{}, mErr("'packages[%d].source' must match '<owner>/<repo>' or './<path>' shape, got '%s'", index, source)
	}
	isLocal := localSourceRE.MatchString(source)

	entry := PackageEntry{Name: name, Source: source, IsLocal: isLocal}

	if v, ok := getStr(m, "subdir"); ok && strings.TrimSpace(v) != "" {
		entry.Subdir = strings.TrimSpace(v)
	}
	if v, ok := getStr(m, "version"); ok && strings.TrimSpace(v) != "" {
		entry.Version = strings.TrimSpace(v)
	}
	if v, ok := getStr(m, "ref"); ok && strings.TrimSpace(v) != "" {
		entry.Ref = strings.TrimSpace(v)
	}
	if !isLocal && entry.Version == "" && entry.Ref == "" {
		return PackageEntry{}, mErr("packages[%d] ('%s'): remote packages require at least one of 'version' or 'ref'", index, name)
	}
	if v, ok := getStr(m, "tag_pattern"); ok && strings.TrimSpace(v) != "" {
		tp := strings.TrimSpace(v)
		if err := validateTagPattern(tp, fmt.Sprintf("packages[%d].tag_pattern", index)); err != nil {
			return PackageEntry{}, err
		}
		entry.TagPattern = tp
	}
	if v, ok := m["include_prerelease"].(bool); ok {
		entry.IncludePrerelease = v
	}
	if v, ok := getStr(m, "description"); ok {
		entry.Description = strings.TrimSpace(v)
	}
	if v, ok := getStr(m, "homepage"); ok {
		entry.Homepage = strings.TrimSpace(v)
	}
	if v, ok := getStr(m, "license"); ok {
		entry.License = strings.TrimSpace(v)
	}
	if v, ok := getStr(m, "repository"); ok {
		entry.Repository = strings.TrimSpace(v)
	}

	// Tags + keywords merge
	var tags []string
	if rawTags, ok := m["tags"].([]interface{}); ok {
		for _, t := range rawTags {
			if s, ok := t.(string); ok {
				tags = append(tags, s)
			}
		}
	}
	if rawKW, ok := m["keywords"].([]interface{}); ok {
		seen := map[string]bool{}
		for _, t := range tags {
			seen[t] = true
		}
		for _, t := range rawKW {
			if s, ok := t.(string); ok && !seen[s] {
				tags = append(tags, s)
				seen[s] = true
			}
		}
	}
	if len(tags) > maxTagsCount {
		tags = tags[:maxTagsCount]
	}
	for i, t := range tags {
		if len(t) > maxTagLength {
			tags[i] = t[:maxTagLength]
		}
	}
	entry.Tags = tags

	// Author
	if rawAuthor, ok := m["author"]; ok && rawAuthor != nil {
		switch a := rawAuthor.(type) {
		case string:
			n := strings.TrimSpace(a)
			if n == "" {
				return PackageEntry{}, mErr("'packages[%d].author' must be a non-empty string or object with 'name'", index)
			}
			entry.Author = map[string]string{"name": n}
		case map[string]interface{}:
			n, ok := getStr(a, "name")
			if !ok || strings.TrimSpace(n) == "" {
				return PackageEntry{}, mErr("'packages[%d].author.name' is required", index)
			}
			auth := map[string]string{"name": strings.TrimSpace(n)}
			for _, k := range []string{"email", "url"} {
				if v, ok := getStr(a, k); ok && strings.TrimSpace(v) != "" {
					auth[k] = strings.TrimSpace(v)
				}
			}
			entry.Author = auth
		}
	}

	return entry, nil
}

// LoadFromFile loads a MarketplaceConfig from a file path.
// It reads the file as raw text and uses a minimal parser.
func LoadFromFile(path string, isLegacy bool) (*MarketplaceConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, mErr("Cannot read '%s': %v", path, err)
	}

	// Use simple key-value extraction for top-level scalars
	flat := parseSimpleYAML(string(content))

	cfg := &MarketplaceConfig{
		SourcePath: path,
		IsLegacy:   isLegacy,
		Build:      MarketplaceBuild{TagPattern: "v{version}"},
		Output:     ".claude-plugin/marketplace.json",
	}
	if isLegacy {
		cfg.Output = "marketplace.json"
	}

	if v := flat["name"]; v != "" {
		cfg.Name = v
		cfg.NameOverridden = isLegacy
	}
	if v := flat["description"]; v != "" {
		cfg.Description = v
		cfg.DescriptionOverridden = isLegacy
	}
	if v := flat["version"]; v != "" {
		cfg.Version = v
		cfg.VersionOverridden = isLegacy
		if cfg.Version != "" {
			if err := validateSemver(cfg.Version, "version"); err != nil {
				return nil, err
			}
		}
	}
	if v := flat["output"]; v != "" {
		cfg.Output = v
		if err := validatePathSegments(cfg.Output); err != nil {
			return nil, mErr("invalid output path: %v", err)
		}
	}

	// Owner (required)
	ownerName := flat["owner.name"]
	if ownerName == "" {
		// Try to extract owner.name from nested YAML manually
		ownerName = extractNestedValue(string(content), "owner", "name")
	}
	if ownerName == "" {
		return nil, mErr("'owner' is required")
	}
	cfg.Owner = MarketplaceOwner{
		Name:  ownerName,
		Email: extractNestedValue(string(content), "owner", "email"),
		URL:   extractNestedValue(string(content), "owner", "url"),
	}

	return cfg, nil
}

// extractNestedValue extracts a value from a 2-level YAML structure without
// a full YAML parser. Used for simple cases like owner.name.
func extractNestedValue(content, parent, key string) string {
	inParent := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			// Top-level line
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, parent+":") {
				inParent = true
			} else if trimmed != "" {
				inParent = false
			}
			continue
		}
		if inParent {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, key+":") {
				val := strings.TrimSpace(trimmed[len(key)+1:])
				return strings.Trim(val, "\"'")
			}
		}
	}
	return ""
}
