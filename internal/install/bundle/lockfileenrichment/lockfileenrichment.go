// Package lockfileenrichment provides lockfile enrichment for pack-time metadata.
//
// Migrated from src/apm_cli/bundle/lockfile_enrichment.py
package lockfileenrichment

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"
)

// crossTargetMaps maps target names to source->destination prefix mappings
// for cross-target skills/ and agents/ path remapping.
var crossTargetMaps = map[string]map[string]string{
	"claude": {
		".github/skills/": ".claude/skills/",
		".github/agents/": ".claude/agents/",
	},
	"vscode": {
		".claude/skills/": ".github/skills/",
		".claude/agents/": ".github/agents/",
	},
	"copilot": {
		".claude/skills/": ".github/skills/",
		".claude/agents/": ".github/agents/",
	},
	"cursor": {
		".github/skills/": ".cursor/skills/",
		".github/agents/": ".cursor/agents/",
	},
	"opencode": {
		".github/skills/": ".opencode/skills/",
		".github/agents/": ".opencode/agents/",
	},
	"codex": {
		".github/skills/": ".agents/skills/",
		".github/agents/": ".codex/agents/",
	},
	"windsurf": {
		".github/skills/": ".windsurf/skills/",
		".github/agents/": ".windsurf/skills/",
	},
	"agent-skills": {
		".github/skills/": ".agents/skills/",
	},
}

// knownTargetPrefixes maps target names to their effective pack prefixes.
var knownTargetPrefixes = map[string][]string{
	"copilot":      {".github/"},
	"vscode":       {".github/"},
	"claude":       {".claude/"},
	"cursor":       {".cursor/"},
	"opencode":     {".opencode/"},
	"codex":        {".codex/", ".agents/"},
	"windsurf":     {".windsurf/"},
	"agent-skills": {".agents/"},
}

// allTargetPrefixes returns the union of pack prefixes for every deployable target.
func allTargetPrefixes() []string {
	seen := map[string]bool{}
	var prefixes []string
	// Stable order
	order := []string{"copilot", "vscode", "claude", "cursor", "opencode", "codex", "windsurf", "agent-skills"}
	for _, t := range order {
		for _, p := range knownTargetPrefixes[t] {
			if !seen[p] {
				seen[p] = true
				prefixes = append(prefixes, p)
			}
		}
	}
	return prefixes
}

// getTargetPrefixes resolves pack-prefixes for a single target name.
func getTargetPrefixes(target string) []string {
	if target == "all" {
		return allTargetPrefixes()
	}
	if target == "vscode" {
		return knownTargetPrefixes["copilot"]
	}
	if ps, ok := knownTargetPrefixes[target]; ok {
		return ps
	}
	return allTargetPrefixes()
}

// FilterFilesResult holds the result of FilterFilesByTarget.
type FilterFilesResult struct {
	Files        []string
	PathMappings map[string]string // bundle_path -> disk_path for cross-target remapped files
}

// FilterFilesByTarget filters deployed file paths by target prefix with cross-target mapping.
// target may be a single string or comma-separated list.
func FilterFilesByTarget(deployedFiles []string, target string) FilterFilesResult {
	targets := strings.Split(target, ",")
	for i := range targets {
		targets[i] = strings.TrimSpace(targets[i])
	}

	var prefixes []string
	seenPrefixes := map[string]bool{}
	crossMap := map[string]string{}

	for _, t := range targets {
		for _, p := range getTargetPrefixes(t) {
			if !seenPrefixes[p] {
				seenPrefixes[p] = true
				prefixes = append(prefixes, p)
			}
		}
		for k, v := range crossTargetMaps[t] {
			crossMap[k] = v
		}
	}

	var direct []string
	directSet := map[string]bool{}

	for _, f := range deployedFiles {
		for _, p := range prefixes {
			if strings.HasPrefix(f, p) {
				direct = append(direct, f)
				directSet[f] = true
				break
			}
		}
	}

	pathMappings := map[string]string{}
	if len(crossMap) > 0 {
		for _, f := range deployedFiles {
			if directSet[f] {
				continue
			}
			for srcPrefix, dstPrefix := range crossMap {
				if strings.HasPrefix(f, srcPrefix) {
					mapped := dstPrefix + f[len(srcPrefix):]
					// Path traversal guard
					normalised := path.Clean(mapped)
					if strings.Contains(normalised, "..") {
						continue
					}
					if !strings.HasPrefix(normalised, strings.TrimSuffix(dstPrefix, "/")) {
						continue
					}
					// Preserve trailing slash
					if strings.HasSuffix(mapped, "/") && !strings.HasSuffix(normalised, "/") {
						normalised += "/"
					}
					mapped = normalised
					if !directSet[mapped] {
						direct = append(direct, mapped)
						directSet[mapped] = true
						pathMappings[mapped] = f
					}
					break
				}
			}
		}
	}

	return FilterFilesResult{Files: direct, PathMappings: pathMappings}
}

// PackMeta holds pack section metadata for the enriched lockfile.
type PackMeta struct {
	Format      string
	Target      string
	PackedAt    string
	MappedFrom  []string
	BundleFiles map[string]string
}

// EnrichLockfileForPack generates a pack: metadata YAML block.
// It returns the pack section YAML string to prepend to the lockfile.
func EnrichLockfileForPack(meta PackMeta) string {
	if meta.PackedAt == "" {
		meta.PackedAt = time.Now().UTC().Format(time.RFC3339)
	}

	var sb strings.Builder
	sb.WriteString("pack:\n")
	sb.WriteString(fmt.Sprintf("  format: %s\n", yamlStr(meta.Format)))
	sb.WriteString(fmt.Sprintf("  target: %s\n", yamlStr(meta.Target)))
	sb.WriteString(fmt.Sprintf("  packed_at: %s\n", yamlStr(meta.PackedAt)))

	if len(meta.MappedFrom) > 0 {
		sb.WriteString("  mapped_from:\n")
		for _, m := range meta.MappedFrom {
			sb.WriteString(fmt.Sprintf("  - %s\n", yamlStr(m)))
		}
	}

	if len(meta.BundleFiles) > 0 {
		sb.WriteString("  bundle_files:\n")
		keys := make([]string, 0, len(meta.BundleFiles))
		for k := range meta.BundleFiles {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("    %s: %s\n", yamlStr(k), yamlStr(meta.BundleFiles[k])))
		}
	}

	return sb.String()
}

// CollectMappedFromPrefixes returns the source prefixes that were actually remapped,
// given all cross-target path mappings for the target and the original->mapped pairs.
func CollectMappedFromPrefixes(target string, originalPaths []string) []string {
	targets := strings.Split(target, ",")
	crossMap := map[string]string{}
	for _, t := range targets {
		for k, v := range crossTargetMaps[strings.TrimSpace(t)] {
			crossMap[k] = v
		}
	}

	used := map[string]bool{}
	for _, orig := range originalPaths {
		for srcPrefix := range crossMap {
			if strings.HasPrefix(orig, srcPrefix) {
				used[srcPrefix] = true
				break
			}
		}
	}

	result := make([]string, 0, len(used))
	for k := range used {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

// yamlStr returns a YAML-safe quoted string for simple values.
func yamlStr(s string) string {
	if s == "" || strings.ContainsAny(s, ":#{}[]|>&*!,") || strings.Contains(s, "  ") {
		return fmt.Sprintf("%q", s)
	}
	return s
}
