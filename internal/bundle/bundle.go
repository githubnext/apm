// Package bundle provides data types and helpers for APM bundle pack/unpack
// operations. It mirrors the Python apm_cli.bundle module.
package bundle

// PackResult describes the result of a pack operation.
type PackResult struct {
	// BundlePath is the filesystem path to the produced bundle.
	BundlePath string
	// Files is the list of relative file paths included in the bundle.
	Files []string
	// LockfileEnriched reports whether the lockfile was enriched during packing.
	LockfileEnriched bool
	// MappedCount is the number of path-mapping entries applied.
	MappedCount int
	// PathMappings holds source->dest path remappings applied during packing.
	PathMappings map[string]string
}

// UnpackResult describes the result of an unpack operation.
type UnpackResult struct {
	// ExtractedDir is the directory where bundle contents were placed.
	ExtractedDir string
	// Files is the deduplicated list of relative file paths extracted.
	Files []string
	// Verified reports whether completeness verification was performed and passed.
	Verified bool
	// DependencyFiles maps dependency keys to their lists of deployed files.
	DependencyFiles map[string][]string
	// SkippedCount is the number of files skipped (symlinks, missing, etc.).
	SkippedCount int
	// SecurityWarnings is the count of files with non-critical hidden characters.
	SecurityWarnings int
	// SecurityCritical is the count of files with critical hidden characters.
	SecurityCritical int
	// PackMeta holds arbitrary pack-section metadata from the embedded lockfile.
	PackMeta map[string]interface{}
}

// LocalBundleInfo is a frozen descriptor for a detected local bundle.
type LocalBundleInfo struct {
	// SourceDir is the filesystem path to the bundle root.
	SourceDir string
	// PluginJSON is the parsed plugin.json content (empty map when absent).
	PluginJSON map[string]interface{}
	// PackageID is derived from plugin.json["id"], falling back to the bundle directory name.
	PackageID string
	// Lockfile is the parsed apm.lock.yaml content, or nil for older bundles.
	Lockfile map[string]interface{}
	// PackTargets lists the targets the bundle was packed for.
	PackTargets []string
	// IsArchive is true when the source path was a .tar.gz.
	IsArchive bool
	// TempDir is the extraction directory for tarballs (caller must clean up).
	// Empty string when not applicable.
	TempDir string
}

// ExtractPackTargets returns the list of pack targets from a parsed bundle lockfile.
// Returns an empty slice when the lockfile is nil or carries no target.
func ExtractPackTargets(lockfile map[string]interface{}) []string {
	if lockfile == nil {
		return []string{}
	}
	pack, _ := lockfile["pack"].(map[string]interface{})
	if pack == nil {
		return []string{}
	}
	raw := pack["target"]
	if raw == nil {
		return []string{}
	}
	switch v := raw.(type) {
	case string:
		if v == "" {
			return []string{}
		}
		return []string{v}
	case []interface{}:
		targets := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				targets = append(targets, s)
			}
		}
		return targets
	}
	return []string{}
}

// CheckTargetMismatch returns a warning string when the bundle targets are not
// covered by the install targets. Returns an empty string when:
//   - bundleTargets is empty (pre-constraint bundle, no metadata), OR
//   - bundleTargets contains "all" (target-agnostic bundle), OR
//   - installTargets is a superset of bundleTargets.
func CheckTargetMismatch(bundleTargets, installTargets []string) string {
	if len(bundleTargets) == 0 {
		return ""
	}
	bundleSet := make(map[string]bool, len(bundleTargets))
	for _, t := range bundleTargets {
		if t != "" {
			bundleSet[t] = true
		}
	}
	if bundleSet["all"] {
		return ""
	}
	installSet := make(map[string]bool, len(installTargets))
	for _, t := range installTargets {
		if t != "" {
			installSet[t] = true
		}
	}
	var missing []string
	for t := range bundleSet {
		if !installSet[t] {
			missing = append(missing, t)
		}
	}
	if len(missing) == 0 {
		return ""
	}
	// Sort for determinism
	sortStrings(missing)
	packed := sortedKeys(bundleSet)
	active := sortedKeys(installSet)
	activeStr := joinStrings(active)
	if activeStr == "" {
		activeStr = "<none>"
	}
	return "Bundle was packed for targets [" + joinStrings(packed) + "] but install resolved to [" +
		activeStr + "]. The following packed targets will not receive files: " +
		joinStrings(missing)
}

// IsSafeRelPath returns true when rel is safe to write inside an output directory.
// It rejects absolute paths and paths containing ".." components.
func IsSafeRelPath(rel string) bool {
	if rel == "" {
		return false
	}
	if rel[0] == '/' || rel[0] == '\\' {
		return false
	}
	// Check for Windows absolute (e.g. C:\)
	if len(rel) >= 2 && rel[1] == ':' {
		return false
	}
	// Walk path components
	parts := splitPathParts(rel)
	for _, p := range parts {
		if p == ".." {
			return false
		}
	}
	return true
}

// --- string/path helpers (no external deps) ---

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func joinStrings(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}

func splitPathParts(p string) []string {
	var parts []string
	cur := ""
	for _, c := range p {
		if c == '/' || c == '\\' {
			if cur != "" {
				parts = append(parts, cur)
				cur = ""
			}
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		parts = append(parts, cur)
	}
	return parts
}
