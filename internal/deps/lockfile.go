// Package deps -- LockedDependency and LockFile data structures.
// Mirrors src/apm_cli/deps/lockfile.py (core types only).
package deps

import (
	"sort"
)

// LockedDependency is a resolved dependency with exact commit/version info.
// Mirrors src/apm_cli/deps/lockfile.py:LockedDependency.
type LockedDependency struct {
	RepoURL                string
	Host                   string
	Port                   int    // 0 means not set
	RegistryPrefix         string
	ResolvedCommit         string
	ResolvedRef            string
	Version                string
	VirtualPath            string
	IsVirtual              bool
	Depth                  int
	ResolvedBy             string
	PackageType            string
	DeployedFiles          []string
	DeployedFileHashes     map[string]string
	Source                 string // "local" for local deps
	LocalPath              string
	ContentHash            string
	IsDev                  bool
	DiscoveredVia          string
	MarketplacePluginName  string
	IsInsecure             bool
	AllowInsecure          bool
	SkillSubset            []string
}

// GetUniqueKey returns a stable key for this locked dependency.
// Mirrors LockedDependency.get_unique_key().
func (d *LockedDependency) GetUniqueKey() string {
	if d.Source == "local" && d.LocalPath != "" {
		return d.LocalPath
	}
	if d.IsVirtual && d.VirtualPath != "" {
		return d.RepoURL + "/" + d.VirtualPath
	}
	return d.RepoURL
}

// ToMap serializes the locked dependency to a map (for YAML output).
// Mirrors LockedDependency.to_dict().
func (d *LockedDependency) ToMap() map[string]interface{} {
	m := map[string]interface{}{"repo_url": d.RepoURL}
	if d.Host != "" {
		m["host"] = d.Host
	}
	if d.Port != 0 {
		m["port"] = d.Port
	}
	if d.RegistryPrefix != "" {
		m["registry_prefix"] = d.RegistryPrefix
	}
	if d.ResolvedCommit != "" {
		m["resolved_commit"] = d.ResolvedCommit
	}
	if d.ResolvedRef != "" {
		m["resolved_ref"] = d.ResolvedRef
	}
	if d.Version != "" {
		m["version"] = d.Version
	}
	if d.VirtualPath != "" {
		m["virtual_path"] = d.VirtualPath
	}
	if d.IsVirtual {
		m["is_virtual"] = true
	}
	if d.Depth != 1 {
		m["depth"] = d.Depth
	}
	if d.ResolvedBy != "" {
		m["resolved_by"] = d.ResolvedBy
	}
	if d.PackageType != "" {
		m["package_type"] = d.PackageType
	}
	if len(d.DeployedFiles) > 0 {
		files := make([]string, len(d.DeployedFiles))
		copy(files, d.DeployedFiles)
		sort.Strings(files)
		m["deployed_files"] = files
	}
	if len(d.DeployedFileHashes) > 0 {
		m["deployed_file_hashes"] = d.DeployedFileHashes
	}
	if d.Source != "" {
		m["source"] = d.Source
	}
	if d.LocalPath != "" {
		m["local_path"] = d.LocalPath
	}
	if d.ContentHash != "" {
		m["content_hash"] = d.ContentHash
	}
	if d.IsDev {
		m["is_dev"] = true
	}
	if d.DiscoveredVia != "" {
		m["discovered_via"] = d.DiscoveredVia
	}
	if d.MarketplacePluginName != "" {
		m["marketplace_plugin_name"] = d.MarketplacePluginName
	}
	if d.IsInsecure {
		m["is_insecure"] = true
	}
	if d.AllowInsecure {
		m["allow_insecure"] = true
	}
	if len(d.SkillSubset) > 0 {
		ss := make([]string, len(d.SkillSubset))
		copy(ss, d.SkillSubset)
		sort.Strings(ss)
		m["skill_subset"] = ss
	}
	return m
}

// LockedDependencyFromMap deserializes a LockedDependency from a map.
// Mirrors LockedDependency.from_dict().
func LockedDependencyFromMap(data map[string]interface{}) *LockedDependency {
	ld := &LockedDependency{
		RepoURL:               stringField(data, "repo_url"),
		Host:                  stringField(data, "host"),
		RegistryPrefix:        stringField(data, "registry_prefix"),
		ResolvedCommit:        stringField(data, "resolved_commit"),
		ResolvedRef:           stringField(data, "resolved_ref"),
		Version:               stringField(data, "version"),
		VirtualPath:           stringField(data, "virtual_path"),
		IsVirtual:             boolField(data, "is_virtual"),
		Depth:                 intFieldDefault(data, "depth", 1),
		ResolvedBy:            stringField(data, "resolved_by"),
		PackageType:           stringField(data, "package_type"),
		Source:                stringField(data, "source"),
		LocalPath:             stringField(data, "local_path"),
		ContentHash:           stringField(data, "content_hash"),
		IsDev:                 boolField(data, "is_dev"),
		DiscoveredVia:         stringField(data, "discovered_via"),
		MarketplacePluginName: stringField(data, "marketplace_plugin_name"),
		IsInsecure:            boolField(data, "is_insecure"),
		AllowInsecure:         boolField(data, "allow_insecure"),
	}

	// Port with validation (1-65535).
	if pRaw, ok := data["port"]; ok && pRaw != nil {
		switch v := pRaw.(type) {
		case int:
			if v >= 1 && v <= 65535 {
				ld.Port = v
			}
		case float64:
			iv := int(v)
			if iv >= 1 && iv <= 65535 {
				ld.Port = iv
			}
		}
	}

	// deployed_files
	if raw, ok := data["deployed_files"]; ok {
		if sl, ok := raw.([]interface{}); ok {
			for _, v := range sl {
				if s, ok := v.(string); ok {
					ld.DeployedFiles = append(ld.DeployedFiles, s)
				}
			}
		}
	}
	// deployed_file_hashes
	if raw, ok := data["deployed_file_hashes"]; ok {
		if m, ok := raw.(map[string]interface{}); ok {
			ld.DeployedFileHashes = make(map[string]string, len(m))
			for k, v := range m {
				if s, ok := v.(string); ok {
					ld.DeployedFileHashes[k] = s
				}
			}
		}
	}
	// skill_subset
	if raw, ok := data["skill_subset"]; ok {
		if sl, ok := raw.([]interface{}); ok {
			for _, v := range sl {
				if s, ok := v.(string); ok {
					ld.SkillSubset = append(ld.SkillSubset, s)
				}
			}
		}
	}

	return ld
}

// InstalledPackage records a successfully-installed dependency.
// Mirrors src/apm_cli/deps/installed_package.py:InstalledPackage.
type InstalledPackage struct {
	RepoURL        string
	Reference      string
	ResolvedCommit string
	Depth          int
	ResolvedBy     string
	IsDev          bool
	RegistryHost   string // from RegistryConfig.host if set
	RegistryPrefix string // from RegistryConfig.prefix if set
}

// helpers

func stringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func boolField(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func intFieldDefault(m map[string]interface{}, key string, def int) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case float64:
			return int(n)
		}
	}
	return def
}
