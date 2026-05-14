// Package agentintegrator handles integration of APM package agents into
// .github/agents/, .claude/agents/, .cursor/agents/ etc.
// Ported from src/apm_cli/integration/agent_integrator.py
package agentintegrator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/githubnext/apm/internal/integration/baseintegrator"
	"github.com/githubnext/apm/internal/integration/targets"
)

// AgentIntegrator handles agent file integration for a single package.
type AgentIntegrator struct{}

// FindAgentFiles returns all .agent.md and .chatmode.md files in a package.
// Searches package root, .apm/agents/ (with rglob), and .apm/chatmodes/ (legacy).
func FindAgentFiles(packagePath string) []string {
	var agentFiles []string
	seen := map[string]struct{}{}

	add := func(p string) {
		abs, _ := filepath.Abs(p)
		if _, ok := seen[abs]; !ok {
			seen[abs] = struct{}{}
			agentFiles = append(agentFiles, p)
		}
	}

	// Package root: *.agent.md and *.chatmode.md
	if entries, err := os.ReadDir(packagePath); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			n := e.Name()
			if strings.HasSuffix(n, ".agent.md") || strings.HasSuffix(n, ".chatmode.md") {
				add(filepath.Join(packagePath, n))
			}
		}
	}

	// .apm/agents/ -- rglob *.agent.md + plain .md files
	apmAgentsDir := filepath.Join(packagePath, ".apm", "agents")
	if _, err := os.Stat(apmAgentsDir); err == nil {
		filepath.WalkDir(apmAgentsDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			n := d.Name()
			if strings.HasSuffix(n, ".agent.md") {
				add(path)
			} else if strings.HasSuffix(n, ".md") {
				add(path)
			}
			return nil
		})
	}

	// .apm/chatmodes/ (legacy)
	apmChatmodesDir := filepath.Join(packagePath, ".apm", "chatmodes")
	if _, err := os.Stat(apmChatmodesDir); err == nil {
		if entries, err := os.ReadDir(apmChatmodesDir); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				if strings.HasSuffix(e.Name(), ".chatmode.md") {
					add(filepath.Join(apmChatmodesDir, e.Name()))
				}
			}
		}
	}

	return agentFiles
}

// GetTargetFilenameForTarget generates the target filename for an agent file
// using the extension from target's agents mapping.
func GetTargetFilenameForTarget(sourceFile string, target *targets.TargetProfile) string {
	mapping, ok := target.Primitives["agents"]
	ext := ".agent.md"
	if ok {
		ext = mapping.Extension
	}
	name := filepath.Base(sourceFile)
	var stem string
	if strings.HasSuffix(name, ".agent.md") {
		stem = name[:len(name)-9]
	} else if strings.HasSuffix(name, ".chatmode.md") {
		stem = name[:len(name)-12]
	} else {
		stem = strings.TrimSuffix(name, filepath.Ext(name))
	}
	return stem + ext
}

// PortableRelpath returns a relative path from base to target using forward slashes.
func PortableRelpath(targetPath, basePath string) string {
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return targetPath
	}
	return filepath.ToSlash(rel)
}

// CopyAgent copies a source agent file to target, returning links resolved count (stub: 0).
func CopyAgent(source, target string) (int, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return 0, err
	}
	if err := os.WriteFile(target, data, 0644); err != nil {
		return 0, err
	}
	return 0, nil
}

// IntegrateAgentsForTarget integrates agents from a package for a single target.
func IntegrateAgentsForTarget(
	target *targets.TargetProfile,
	installPath string,
	projectRoot string,
	force bool,
	managedFiles map[string]struct{},
	diag baseintegrator.Diagnostics,
) baseintegrator.IntegrationResult {
	mapping, ok := target.Primitives["agents"]
	if !ok {
		return baseintegrator.IntegrationResult{}
	}

	effectiveRoot := mapping.DeployRoot
	if effectiveRoot == "" {
		effectiveRoot = target.RootDir
	}
	targetRoot := filepath.Join(projectRoot, effectiveRoot)

	if !target.AutoCreate {
		if _, err := os.Stat(filepath.Join(projectRoot, target.RootDir)); os.IsNotExist(err) {
			return baseintegrator.IntegrationResult{}
		}
	}

	agentFiles := FindAgentFiles(installPath)
	if len(agentFiles) == 0 {
		return baseintegrator.IntegrationResult{}
	}

	agentsDir := targetRoot
	if mapping.Subdir != "" {
		agentsDir = filepath.Join(targetRoot, mapping.Subdir)
	}
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return baseintegrator.IntegrationResult{}
	}

	var result baseintegrator.IntegrationResult

	for _, sourceFile := range agentFiles {
		targetFilename := GetTargetFilenameForTarget(sourceFile, target)
		targetPath := filepath.Join(agentsDir, targetFilename)
		relPath := PortableRelpath(targetPath, projectRoot)

		if baseintegrator.CheckCollision(targetPath, relPath, managedFiles, force, diag) {
			result.FilesSkipped++
			continue
		}

		var linksResolved int
		var err error

		switch mapping.FormatID {
		case "codex_agent":
			err = writeCodexAgent(sourceFile, targetPath)
		case "windsurf_agent_skill":
			linksResolved, err = writeWindsurfAgentSkill(sourceFile, targetPath, diag)
		default:
			linksResolved, err = CopyAgent(sourceFile, targetPath)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to write agent %s: %v\n", targetFilename, err)
			continue
		}

		result.LinksResolved += linksResolved
		result.FilesIntegrated++
		result.TargetPaths = append(result.TargetPaths, targetPath)
	}

	return result
}

// SyncForTarget removes APM-managed agent files for a single target.
func SyncForTarget(
	target *targets.TargetProfile,
	projectRoot string,
	managedFiles map[string]struct{},
) baseintegrator.SyncRemoveResult {
	mapping, ok := target.Primitives["agents"]
	if !ok {
		return baseintegrator.SyncRemoveResult{}
	}
	effectiveRoot := mapping.DeployRoot
	if effectiveRoot == "" {
		effectiveRoot = target.RootDir
	}
	prefix := effectiveRoot + "/" + mapping.Subdir + "/"
	legacyDir := filepath.Join(projectRoot, effectiveRoot, mapping.Subdir)
	legacyPattern := "*-apm.md"
	if mapping.Extension == ".agent.md" {
		legacyPattern = "*-apm.agent.md"
	}
	return baseintegrator.SyncRemoveFiles(
		projectRoot,
		managedFiles,
		prefix,
		legacyDir,
		legacyPattern,
		[]*targets.TargetProfile{target},
		nil,
	)
}

// frontmatterRE matches YAML frontmatter in markdown.
var frontmatterRE = regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---\s*\n?`)

// writeCodexAgent transforms an .agent.md file to Codex .toml format.
// Produces a minimal TOML output without an external dependency.
func writeCodexAgent(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	content := string(data)

	name := filepath.Base(source)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	if strings.HasSuffix(name, ".agent") {
		name = name[:len(name)-6]
	}
	description := ""
	body := content

	if m := frontmatterRE.FindStringSubmatchIndex(content); m != nil {
		fmStr := content[m[2]:m[3]]
		body = content[m[1]:]
		fm := parseSimpleYAML(fmStr)
		if v, ok := fm["name"]; ok {
			name = v
		}
		if v, ok := fm["description"]; ok {
			description = v
		}
	}

	body = strings.TrimSpace(body)

	// Produce minimal TOML
	var sb strings.Builder
	sb.WriteString("name = ")
	sb.WriteString(tomlQuote(name))
	sb.WriteString("\ndescription = ")
	sb.WriteString(tomlQuote(description))
	sb.WriteString("\ndeveloper_instructions = ")
	sb.WriteString(tomlMultilineQuote(body))
	sb.WriteString("\n")

	return os.WriteFile(target, []byte(sb.String()), 0644)
}

// writeWindsurfAgentSkill transforms an .agent.md file to a Windsurf Skill (SKILL.md).
func writeWindsurfAgentSkill(source, target string, diag baseintegrator.Diagnostics) (int, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return 0, err
	}
	content := string(data)

	name := filepath.Base(source)
	if strings.HasSuffix(name, ".agent.md") {
		name = name[:len(name)-9]
	} else if strings.HasSuffix(name, ".chatmode.md") {
		name = name[:len(name)-12]
	} else {
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}

	description := ""
	body := content
	var fmMap map[string]string

	if m := frontmatterRE.FindStringSubmatchIndex(content); m != nil {
		fmMap = parseSimpleYAML(content[m[2]:m[3]])
		body = content[m[1]:]
	} else {
		fmMap = map[string]string{}
	}

	if diag != nil {
		var dropped []string
		for _, k := range []string{"tools", "model"} {
			if v, ok := fmMap[k]; ok && v != "" {
				dropped = append(dropped, k)
			}
		}
		if len(dropped) > 0 {
			diag.Warn(
				fmt.Sprintf("Windsurf skill conversion dropped frontmatter field(s) %s from %s",
					strings.Join(dropped, ", "), filepath.Base(source)),
				"Windsurf Skills do not support agent-only fields; only name, description, and body are preserved.",
			)
		}
	}

	if v, ok := fmMap["name"]; ok {
		name = v
	}
	if v, ok := fmMap["description"]; ok {
		description = v
	}

	var fm strings.Builder
	fm.WriteString("name: ")
	fm.WriteString(yamlQuoteIfNeeded(name))
	if description != "" {
		fm.WriteString("\ndescription: ")
		fm.WriteString(yamlQuoteIfNeeded(description))
	}

	result := "---\n" + fm.String() + "\n---\n" + body

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return 0, err
	}
	if err := os.WriteFile(target, []byte(result), 0644); err != nil {
		return 0, err
	}
	return 0, nil
}

// parseSimpleYAML parses simple key: value YAML lines (no nesting, no lists).
func parseSimpleYAML(s string) map[string]string {
	result := map[string]string{}
	for _, line := range strings.Split(s, "\n") {
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		val := strings.TrimSpace(line[colon+1:])
		// Strip surrounding quotes
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		result[key] = val
	}
	return result
}

// tomlQuote wraps a string in TOML basic string quotes.
func tomlQuote(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return `"` + s + `"`
}

// tomlMultilineQuote wraps a string in TOML multi-line basic string quotes.
func tomlMultilineQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"""`, `\"\"\"`)
	return `"""` + "\n" + s + "\n" + `"""`
}

// yamlQuoteIfNeeded wraps a value in double quotes if it contains special chars.
func yamlQuoteIfNeeded(s string) string {
	specials := []string{":", "#", "[", "]", "{", "}", ",", "&", "*", "!", "|", ">", "'", "\"", "%", "@", "`"}
	needs := false
	for _, sp := range specials {
		if strings.Contains(s, sp) {
			needs = true
			break
		}
	}
	if needs {
		s = strings.ReplaceAll(s, `"`, `\"`)
		return `"` + s + `"`
	}
	return s
}
