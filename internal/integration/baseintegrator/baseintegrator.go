// Package baseintegrator provides shared collision detection, sync removal,
// link resolution, and file-discovery helpers for file-level integrators.
// Ported from src/apm_cli/integration/base_integrator.py
package baseintegrator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/githubnext/apm/internal/integration/coworkpaths"
	"github.com/githubnext/apm/internal/integration/targets"
)

// IntegrationResult holds the outcome of a file-level integration operation.
type IntegrationResult struct {
	FilesIntegrated  int
	FilesUpdated     int // kept for CLI compat, always 0 today
	FilesSkipped     int
	TargetPaths      []string
	LinksResolved    int
	// hook-specific
	ScriptsCopied    int
	// skill-specific
	SubSkillsPromoted int
	SkillCreated      bool
}

// Diagnostics is a minimal interface for recording integration diagnostics.
type Diagnostics interface {
	Skip(relPath string)
	Warn(msg, detail string)
}

// CheckCollision returns true if targetPath is a user-authored collision.
// A collision exists when: managed set is non-nil, file exists, relPath is NOT
// in the managed set, and force is false.
func CheckCollision(
	targetPath string,
	relPath string,
	managedFiles map[string]struct{},
	force bool,
	diag Diagnostics,
) bool {
	if managedFiles == nil {
		return false
	}
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return false
	}
	norm := strings.ReplaceAll(relPath, "\\", "/")
	if _, ok := managedFiles[norm]; ok {
		return false
	}
	if force {
		return false
	}
	if diag != nil {
		diag.Skip(relPath)
	} else {
		fmt.Fprintf(os.Stderr, "[!] Skipping %s -- local file exists (not managed by APM). Use 'apm install --force' to overwrite.\n", relPath)
	}
	return true
}

// NormalizeManagedFiles normalizes path separators to forward slashes for O(1) lookups.
func NormalizeManagedFiles(managedFiles map[string]struct{}) map[string]struct{} {
	if managedFiles == nil {
		return nil
	}
	out := make(map[string]struct{}, len(managedFiles))
	for p := range managedFiles {
		out[strings.ReplaceAll(p, "\\", "/")] = struct{}{}
	}
	return out
}

// BucketAliases maps raw {prim}_{target} keys to canonical bucket names.
var BucketAliases = map[string]string{
	"prompts_copilot":      "prompts",
	"agents_copilot":       "agents_github",
	"commands_claude":      "commands",
	"commands_cursor":      "commands_cursor",
	"commands_opencode":    "commands_opencode",
	"instructions_copilot": "instructions",
	"instructions_cursor":  "rules_cursor",
	"instructions_claude":  "rules_claude",
}

// PartitionBucketKey returns the canonical bucket key for a (primitive, target) pair.
func PartitionBucketKey(primName, targetName string) string {
	raw := primName + "_" + targetName
	if alias, ok := BucketAliases[raw]; ok {
		return alias
	}
	return raw
}

// PartitionManagedFiles partitions managedFiles by integration prefix.
// When profiles is nil, falls back to targets.KnownTargets.
func PartitionManagedFiles(
	managedFiles map[string]struct{},
	profiles []*targets.TargetProfile,
) map[string]map[string]struct{} {
	source := profiles
	if source == nil {
		for _, p := range targets.KnownTargets {
			source = append(source, p)
		}
	}

	buckets := map[string]map[string]struct{}{
		"skills": {},
		"hooks":  {},
	}

	var skillPrefixes []string
	var hookPrefixes []string

	// prefix -> bucket key
	prefixMap := map[string]string{}

	for _, target := range source {
		for primName, mapping := range target.Primitives {
			if target.ResolvedDeployRoot != "" {
				if primName == "skills" {
					skillPrefixes = append(skillPrefixes, coworkpaths.CoworkLockfilePrefix)
				}
				continue
			}
			effectiveRoot := mapping.DeployRoot
			if effectiveRoot == "" {
				effectiveRoot = target.RootDir
			}
			var prefix string
			if mapping.Subdir != "" {
				prefix = effectiveRoot + "/" + mapping.Subdir + "/"
			} else {
				prefix = effectiveRoot + "/"
			}
			if primName == "skills" {
				skillPrefixes = append(skillPrefixes, prefix)
			} else if primName == "hooks" {
				hookPrefixes = append(hookPrefixes, prefix)
			} else {
				raw := primName + "_" + target.Name
				bucketKey, ok := BucketAliases[raw]
				if !ok {
					bucketKey = raw
				}
				if _, exists := buckets[bucketKey]; !exists {
					buckets[bucketKey] = map[string]struct{}{}
				}
				prefixMap[prefix] = bucketKey
			}
		}
	}

	// Build a trie for longest-prefix-match routing.
	type trieNode struct {
		children map[string]*trieNode
		bucket   string
	}
	root := &trieNode{children: map[string]*trieNode{}}
	for prefix, bucketKey := range prefixMap {
		segs := splitSegments(prefix)
		node := root
		for _, seg := range segs {
			child, ok := node.children[seg]
			if !ok {
				child = &trieNode{children: map[string]*trieNode{}}
				node.children[seg] = child
			}
			node = child
		}
		node.bucket = bucketKey
	}

	for p := range managedFiles {
		segs := splitSegments(p)
		node := root
		lastBucket := ""
		for _, seg := range segs {
			child, ok := node.children[seg]
			if !ok {
				break
			}
			node = child
			if node.bucket != "" {
				lastBucket = node.bucket
			}
		}
		if lastBucket != "" {
			buckets[lastBucket][p] = struct{}{}
			continue
		}
		// Fall back to cross-target buckets
		if hasAnyPrefix(p, skillPrefixes) {
			buckets["skills"][p] = struct{}{}
		} else if hasAnyPrefix(p, hookPrefixes) {
			buckets["hooks"][p] = struct{}{}
		}
	}

	return buckets
}

func splitSegments(path string) []string {
	var segs []string
	for _, s := range strings.Split(path, "/") {
		if s != "" {
			segs = append(segs, s)
		}
	}
	return segs
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// ValidateDeployPath returns true if relPath is safe for APM to deploy or remove.
// Checks: no path traversal, starts with an allowed integration prefix, resolves within projectRoot.
func ValidateDeployPath(
	relPath string,
	projectRoot string,
	allowedPrefixes []string,
	profiles []*targets.TargetProfile,
) bool {
	if strings.Contains(relPath, "..") {
		return false
	}

	if allowedPrefixes == nil {
		allowedPrefixes = targets.GetIntegrationPrefixes(profiles)
	}

	if strings.HasPrefix(relPath, coworkpaths.CoworkURIScheme) {
		if !hasAnyPrefix(relPath, allowedPrefixes) {
			return false
		}
		coworkRoot, err := coworkpaths.ResolveCoworkSkillsDir()
		if err != nil || coworkRoot == "" {
			return false
		}
		_, err = coworkpaths.FromLockfilePath(relPath, coworkRoot)
		return err == nil
	}

	if !hasAnyPrefix(relPath, allowedPrefixes) {
		return false
	}

	target := filepath.Join(projectRoot, relPath)
	resolved, err := filepath.EvalSymlinks(target)
	if err != nil {
		// If path doesn't exist yet, check using Clean
		resolved = filepath.Clean(target)
	}
	projResolved, err := filepath.EvalSymlinks(projectRoot)
	if err != nil {
		projResolved = filepath.Clean(projectRoot)
	}
	return strings.HasPrefix(resolved, projResolved+string(os.PathSeparator)) || resolved == projResolved
}

// CleanupEmptyParents removes empty parent directories bottom-up.
// Stops at stopAt and does not remove stopAt itself.
func CleanupEmptyParents(deletedPaths []string, stopAt string) {
	if len(deletedPaths) == 0 {
		return
	}
	stopResolved, err := filepath.EvalSymlinks(stopAt)
	if err != nil {
		stopResolved = filepath.Clean(stopAt)
	}

	candidates := map[string]struct{}{}
	for _, p := range deletedPaths {
		parent := filepath.Dir(p)
		for parent != stopAt {
			parentResolved, _ := filepath.EvalSymlinks(parent)
			if parentResolved == stopResolved {
				break
			}
			candidates[parent] = struct{}{}
			next := filepath.Dir(parent)
			if next == parent {
				break
			}
			parent = next
		}
	}

	// Sort deepest-first
	sorted := make([]string, 0, len(candidates))
	for d := range candidates {
		sorted = append(sorted, d)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return strings.Count(sorted[i], string(os.PathSeparator)) > strings.Count(sorted[j], string(os.PathSeparator))
	})

	for _, d := range sorted {
		entries, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		if len(entries) == 0 {
			os.Remove(d) // ignore errors
		}
	}
}

// SyncRemoveResult holds the result of a sync removal operation.
type SyncRemoveResult struct {
	FilesRemoved int
	Errors       int
}

// Logger is a minimal interface for sync-remove diagnostic output.
type Logger interface {
	Warning(msg string, symbol string)
}

// SyncRemoveFiles removes APM-managed files matching prefix from managedFiles.
// Falls back to a legacy glob when managedFiles is nil.
func SyncRemoveFiles(
	projectRoot string,
	managedFiles map[string]struct{},
	prefix string,
	legacyGlobDir string,
	legacyGlobPattern string,
	profiles []*targets.TargetProfile,
	logger Logger,
) SyncRemoveResult {
	stats := SyncRemoveResult{}

	if managedFiles != nil {
		coworkRootResolved := false
		coworkRootCached := ""
		coworkOrphansSkipped := 0

		for relPath := range managedFiles {
			if !strings.HasPrefix(relPath, prefix) {
				continue
			}
			if !ValidateDeployPath(relPath, projectRoot, nil, profiles) {
				continue
			}

			var targetPath string
			if strings.HasPrefix(relPath, coworkpaths.CoworkURIScheme) {
				if !coworkRootResolved {
					coworkRootCached, _ = coworkpaths.ResolveCoworkSkillsDir()
					coworkRootResolved = true
				}
				if coworkRootCached == "" {
					coworkOrphansSkipped++
					continue
				}
				resolved, err := coworkpaths.FromLockfilePath(relPath, coworkRootCached)
				if err != nil {
					continue
				}
				targetPath = resolved
			} else {
				targetPath = filepath.Join(projectRoot, relPath)
			}

			if _, err := os.Stat(targetPath); err == nil {
				if err := os.Remove(targetPath); err != nil {
					stats.Errors++
				} else {
					stats.FilesRemoved++
				}
			}
		}

		if coworkOrphansSkipped > 0 {
			word := "entry"
			if coworkOrphansSkipped != 1 {
				word = "entries"
			}
			msg := fmt.Sprintf(
				"Cowork: skipping %d orphaned lockfile %s -- OneDrive path not detected.\n"+
					"Run: apm config set copilot-cowork-skills-dir <path>  "+
					"(or set APM_COPILOT_COWORK_SKILLS_DIR)\n"+
					"to clean up these entries on the next install/uninstall.",
				coworkOrphansSkipped, word,
			)
			if logger != nil {
				logger.Warning(msg, "warning")
			} else {
				fmt.Fprintf(os.Stderr, "[!] %s\n", msg)
			}
		}
	} else if legacyGlobDir != "" && legacyGlobPattern != "" {
		if _, err := os.Stat(legacyGlobDir); err == nil {
			matches, err := filepath.Glob(filepath.Join(legacyGlobDir, legacyGlobPattern))
			if err == nil {
				for _, f := range matches {
					if err := os.Remove(f); err != nil {
						stats.Errors++
					} else {
						stats.FilesRemoved++
					}
				}
			}
		}
	}

	return stats
}

// FindFilesByGlob searches packagePath (and optional subdirs) for pattern.
// Symlinks and hardlinks are rejected.
func FindFilesByGlob(packagePath string, pattern string, subdirs []string) []string {
	var results []string
	seen := map[uint64]struct{}{}

	dirs := []string{packagePath}
	for _, s := range subdirs {
		dirs = append(dirs, filepath.Join(packagePath, s))
	}

	for _, d := range dirs {
		if _, err := os.Stat(d); err != nil {
			continue
		}
		matches, err := filepath.Glob(filepath.Join(d, pattern))
		if err != nil {
			continue
		}
		sort.Strings(matches)
		for _, f := range matches {
			info, err := os.Lstat(f)
			if err != nil {
				continue
			}
			// Reject symlinks
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			// Reject hardlinks (nlink > 1)
			if sys, ok := info.Sys().(*syscall.Stat_t); ok {
				if sys.Nlink > 1 {
					continue
				}
			}
			resolved, err := filepath.EvalSymlinks(f)
			if err != nil {
				resolved = filepath.Clean(f)
			}
			pkgResolved, err := filepath.EvalSymlinks(packagePath)
			if err != nil {
				pkgResolved = filepath.Clean(packagePath)
			}
			if !strings.HasPrefix(resolved, pkgResolved+string(os.PathSeparator)) && resolved != pkgResolved {
				continue
			}
			// Use inode as unique key
			if sys, ok := info.Sys().(*syscall.Stat_t); ok {
				inode := sys.Ino
				if _, exists := seen[inode]; exists {
					continue
				}
				seen[inode] = struct{}{}
			}
			results = append(results, f)
		}
	}
	return results
}
