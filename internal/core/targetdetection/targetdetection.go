// Package targetdetection implements target auto-detection for APM CLI.
// Migrated from src/apm_cli/core/target_detection.py.
package targetdetection

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ValidTargets is the set of canonical target names.
var ValidTargets = map[string]bool{
	"vscode":       true,
	"claude":       true,
	"cursor":       true,
	"opencode":     true,
	"codex":        true,
	"gemini":       true,
	"windsurf":     true,
	"agent-skills": true,
	"all":          true,
	"minimal":      true,
	"copilot":      true, // alias
	"agents":       true, // alias (deprecated)
}

// NormalizeTarget resolves user-facing aliases to canonical internal names.
func NormalizeTarget(t string) string {
	switch t {
	case "copilot", "vscode", "agents":
		return "vscode"
	default:
		return t
	}
}

// CANONICAL_TARGETS_ORDERED lists display-ordered canonical target names.
var CanonicalTargetsOrdered = []string{
	"claude",
	"copilot",
	"cursor",
	"codex",
	"gemini",
	"opencode",
	"windsurf",
}

// CanonicalDeployDirs maps canonical target names to their deploy directories.
var CanonicalDeployDirs = map[string]string{
	"claude":   ".claude/",
	"copilot":  ".github/",
	"cursor":   ".cursor/",
	"codex":    ".codex/",
	"gemini":   ".gemini/",
	"opencode": ".opencode/",
	"windsurf": ".windsurf/",
}

// CanonicalSignal maps canonical target names to their primary detection signal.
var CanonicalSignal = map[string]string{
	"claude":   "CLAUDE.md",
	"copilot":  ".github/copilot-instructions.md",
	"cursor":   ".cursor/",
	"codex":    ".codex/",
	"gemini":   "GEMINI.md",
	"opencode": ".opencode/",
	"windsurf": ".windsurf/",
}

// signalEntry is one row in the whitelist.
type signalEntry struct {
	target    string
	checkType string // "dir" or "file"
	path      string
}

// signalWhitelist is the ordered list of filesystem markers.
var signalWhitelist = []signalEntry{
	{"claude", "dir", ".claude"},
	{"claude", "file", "CLAUDE.md"},
	{"cursor", "dir", ".cursor"},
	{"cursor", "file", ".cursorrules"},
	{"copilot", "file", ".github/copilot-instructions.md"},
	{"codex", "dir", ".codex"},
	{"gemini", "dir", ".gemini"},
	{"gemini", "file", "GEMINI.md"},
	{"opencode", "dir", ".opencode"},
	{"windsurf", "dir", ".windsurf"},
}

// Signal represents a detected filesystem marker.
type Signal struct {
	Target string
	Source string
}

// ResolvedTargets is the result of target resolution.
type ResolvedTargets struct {
	Targets    []string // sorted canonical target names
	Source     string   // human-readable source description
	AutoCreate bool
}

// DetectSignals scans projectRoot for harness markers.
func DetectSignals(projectRoot string) []Signal {
	var found []Signal
	for _, entry := range signalWhitelist {
		full := filepath.Join(projectRoot, entry.path)
		switch entry.checkType {
		case "dir":
			if info, err := os.Stat(full); err == nil && info.IsDir() {
				found = append(found, Signal{Target: entry.target, Source: entry.path + "/"})
			}
		case "file":
			if info, err := os.Stat(full); err == nil && !info.IsDir() {
				found = append(found, Signal{Target: entry.target, Source: entry.path})
			}
		}
	}
	return found
}

// ResolveTargets resolves effective targets. Returns error on ambiguity or missing harness.
// Priority: flag > yamlTargets > auto-detect signals.
func ResolveTargets(projectRoot string, flag []string, yamlTargets []string) (ResolvedTargets, error) {
	// Priority 1: --target flag
	if len(flag) > 0 {
		for _, t := range flag {
			if !ValidTargets[t] {
				return ResolvedTargets{}, fmt.Errorf("unknown target: %s", t)
			}
		}
		sorted := sortedUnique(flag)
		return ResolvedTargets{Targets: sorted, Source: "--target flag", AutoCreate: true}, nil
	}

	// Priority 2: apm.yml targets
	if len(yamlTargets) > 0 {
		sorted := sortedUnique(yamlTargets)
		return ResolvedTargets{Targets: sorted, Source: "apm.yml", AutoCreate: true}, nil
	}

	// Priority 3: auto-detect
	signals := DetectSignals(projectRoot)
	targetSet := map[string]bool{}
	var sources []string
	for _, s := range signals {
		if !targetSet[s.Target] {
			targetSet[s.Target] = true
		}
		sources = append(sources, s.Source)
	}
	sort.Strings(sources)

	targetList := sortedKeys(targetSet)

	if len(targetList) == 0 {
		return ResolvedTargets{}, fmt.Errorf("no harness found in %s", projectRoot)
	}
	if len(targetList) >= 2 {
		return ResolvedTargets{}, fmt.Errorf("ambiguous harness: multiple targets detected: %s", strings.Join(targetList, ", "))
	}

	return ResolvedTargets{
		Targets:    targetList,
		Source:     "auto-detect from " + strings.Join(sources, ", "),
		AutoCreate: true,
	}, nil
}

// ExpandAllTargets expands 'all' to (signals union yamlTargets).
func ExpandAllTargets(projectRoot string, yamlTargets []string) ([]string, error) {
	signals := DetectSignals(projectRoot)
	combined := map[string]bool{}
	for _, s := range signals {
		combined[s.Target] = true
	}
	for _, t := range yamlTargets {
		combined[t] = true
	}
	result := sortedKeys(combined)
	if len(result) == 0 {
		return nil, fmt.Errorf("no harness found in %s", projectRoot)
	}
	return result, nil
}

// FormatProvenance formats a provenance line for CLI output.
func FormatProvenance(resolved ResolvedTargets) string {
	targets := strings.Join(resolved.Targets, ", ")
	return fmt.Sprintf("Targets: %s  (source: %s)", targets, resolved.Source)
}

// DetectTarget implements the legacy v1 detection API.
// Returns (target, reason).
func DetectTarget(projectRoot string, explicitTarget, configTarget string) (string, string) {
	if explicitTarget != "" {
		return NormalizeTarget(explicitTarget), "explicit --target flag"
	}
	if configTarget != "" {
		return NormalizeTarget(configTarget), "apm.yml target"
	}

	githubExists := dirExists(filepath.Join(projectRoot, ".github"))
	claudeExists := dirExists(filepath.Join(projectRoot, ".claude"))
	cursorExists := dirExists(filepath.Join(projectRoot, ".cursor"))
	opencodeExists := dirExists(filepath.Join(projectRoot, ".opencode"))
	codexExists := dirExists(filepath.Join(projectRoot, ".codex"))
	geminiExists := dirExists(filepath.Join(projectRoot, ".gemini"))
	windsurfExists := dirExists(filepath.Join(projectRoot, ".windsurf"))

	var detected []string
	if githubExists {
		detected = append(detected, ".github/")
	}
	if claudeExists {
		detected = append(detected, ".claude/")
	}
	if cursorExists {
		detected = append(detected, ".cursor/")
	}
	if opencodeExists {
		detected = append(detected, ".opencode/")
	}
	if codexExists {
		detected = append(detected, ".codex/")
	}
	if geminiExists {
		detected = append(detected, ".gemini/")
	}
	if windsurfExists {
		detected = append(detected, ".windsurf/")
	}

	if len(detected) >= 2 {
		return "all", fmt.Sprintf("detected %s folders", strings.Join(detected, " and "))
	}
	if githubExists {
		return "vscode", "detected .github/ folder"
	}
	if claudeExists {
		return "claude", "detected .claude/ folder"
	}
	if cursorExists {
		return "cursor", "detected .cursor/ folder"
	}
	if opencodeExists {
		return "opencode", "detected .opencode/ folder"
	}
	if codexExists {
		return "codex", "detected .codex/ folder"
	}
	if geminiExists {
		return "gemini", "detected .gemini/ folder"
	}
	if windsurfExists {
		return "windsurf", "detected .windsurf/ folder"
	}
	return "minimal", "no target folder found"
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func sortedUnique(items []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, s := range items {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	sort.Strings(result)
	return result
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
