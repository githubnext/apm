package core

import (
	"os"
	"path/filepath"
)

// TargetType is the canonical internal target name.
type TargetType = string

// ReasonNoTargetFolder is returned by DetectTarget when no integration folder
// is present.
const ReasonNoTargetFolder = "no target folder found"

// AllCanonicalTargets is the complete set of real (non-pseudo) canonical
// targets. "minimal" is intentionally excluded.
var AllCanonicalTargets = map[string]bool{
	"vscode":   true,
	"claude":   true,
	"cursor":   true,
	"opencode": true,
	"codex":    true,
	"gemini":   true,
	"windsurf": true,
}

// TargetAliases maps user-facing names to canonical internal names.
var TargetAliases = map[string]string{
	"copilot": "vscode",
	"agents":  "vscode",
	"vscode":  "vscode",
}

// DetectTarget detects the appropriate target for compilation and integration.
// It returns (target, reason) following the priority rules:
//  1. Explicit --target flag
//  2. apm.yml target setting
//  3. Auto-detect from existing folders
func DetectTarget(projectRoot, explicitTarget, configTarget string) (string, string) {
	// Priority 1: explicit --target flag
	if explicitTarget != "" {
		return resolveAlias(explicitTarget), "explicit --target flag"
	}
	// Priority 2: apm.yml target
	if configTarget != "" {
		return resolveAlias(configTarget), "apm.yml target"
	}
	// Priority 3: auto-detect from folders
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
		return "all", "detected " + joinStrings(detected, " and ") + " folders"
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
	return "minimal", ReasonNoTargetFolder
}

// ShouldCompileAgentsMD reports whether AGENTS.md should be compiled for the
// given target. AGENTS.md is generated for vscode, opencode, codex, gemini,
// windsurf, all, and minimal.
func ShouldCompileAgentsMD(target string) bool {
	switch target {
	case "vscode", "opencode", "codex", "gemini", "windsurf", "all", "minimal":
		return true
	}
	return false
}

// ShouldCompileClaudeMD reports whether CLAUDE.md should be compiled.
func ShouldCompileClaudeMD(target string) bool {
	return target == "claude" || target == "all"
}

// ShouldCompileGeminiMD reports whether GEMINI.md should be compiled.
func ShouldCompileGeminiMD(target string) bool {
	return target == "gemini" || target == "all"
}

// ShouldCompileCopilotInstructionsMD reports whether
// .github/copilot-instructions.md should be compiled.
func ShouldCompileCopilotInstructionsMD(target string) bool {
	return target == "vscode" || target == "all"
}

// GetTargetDescription returns a human-readable description of what will be
// generated for a target (accepts both internal types and user-facing aliases).
func GetTargetDescription(target string) string {
	normalized := target
	if target == "copilot" || target == "agents" {
		normalized = "vscode"
	}
	descriptions := map[string]string{
		"vscode":       "AGENTS.md + .github/copilot-instructions.md + .github/prompts/ + .github/agents/",
		"claude":       "CLAUDE.md + .claude/commands/ + .claude/agents/ + .claude/skills/",
		"cursor":       ".cursor/agents/ + .cursor/skills/ + .cursor/rules/",
		"opencode":     "AGENTS.md + .opencode/agents/ + .opencode/commands/ + .opencode/skills/",
		"codex":        "AGENTS.md + .agents/skills/ + .codex/agents/ + .codex/hooks.json",
		"gemini":       "GEMINI.md + .gemini/commands/ + .gemini/skills/ + .gemini/settings.json (MCP/hooks)",
		"windsurf":     "AGENTS.md + .windsurf/rules/ + .windsurf/skills/ + .windsurf/workflows/ + .windsurf/hooks.json",
		"agent-skills": ".agents/skills/ only (cross-client shared skills -- no agents, hooks, or commands)",
		"all":          "AGENTS.md + CLAUDE.md + GEMINI.md + .github/copilot-instructions.md + .github/ + .claude/ + .cursor/ + .opencode/ + .codex/ + .gemini/ + .windsurf/ + .agents/",
		"minimal":      "AGENTS.md only (create .github/, .claude/, or .gemini/ for full integration)",
	}
	if d, ok := descriptions[normalized]; ok {
		return d
	}
	return "unknown target"
}

// NormalizeTargetList normalizes a target value to a list of canonical names.
// Returns nil for nil input (meaning "auto-detect").
func NormalizeTargetList(targets []string) []string {
	if targets == nil {
		return nil
	}
	for _, t := range targets {
		if t == "all" {
			return sortedKeys(AllCanonicalTargets)
		}
	}
	seen := map[string]bool{}
	var result []string
	for _, item := range targets {
		canonical := item
		if a, ok := TargetAliases[item]; ok {
			canonical = a
		}
		if !seen[canonical] {
			seen[canonical] = true
			result = append(result, canonical)
		}
	}
	return result
}

// resolveAlias converts a user-facing target name to its canonical form.
func resolveAlias(t string) string {
	if a, ok := TargetAliases[t]; ok {
		return a
	}
	return t
}

// dirExists reports whether path is an existing directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
