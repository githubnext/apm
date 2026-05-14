// Package commandintegrator provides command integration for APM packages.
// Deploys .prompt.md files as slash commands for Claude, Cursor, OpenCode, etc.
// Ported from src/apm_cli/integration/command_integrator.py
package commandintegrator

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/githubnext/apm/internal/integration/baseintegrator"
	"github.com/githubnext/apm/internal/integration/targets"
)

// IntegrationResult holds results of a command integration operation.
type IntegrationResult struct {
	FilesIntegrated int
	FilesUpdated    int
	FilesSkipped    int
	TargetPaths     []string
	LinksResolved   int
}

// inputNameRe matches valid command argument names.
var inputNameRe = regexp.MustCompile(`^[A-Za-z][\w-]{0,63}$`)

// inputRefRe matches ${{input:name}} and ${ input : name } references.
var inputRefRe = regexp.MustCompile(`\$\{\{?\s*input\s*:\s*([\w-]+)\s*\}?\}`)

// preservedCommandKeys is the set of frontmatter keys preserved by the command transformer.
var preservedCommandKeys = map[string]bool{
	"description":   true,
	"allowed-tools": true,
	"allowedTools":  true,
	"model":         true,
	"argument-hint": true,
	"argumentHint":  true,
	"input":         true,
}

// isValidInputName returns true if name is a safe argument identifier.
func isValidInputName(name string) bool {
	return inputNameRe.MatchString(name)
}

// extractInputNames extracts argument names from an APM 'input' frontmatter value.
// input may be a string (single name) or []interface{} (list of names or maps with name key).
func extractInputNames(input interface{}) (valid []string, rejected []string) {
	if input == nil {
		return nil, nil
	}
	switch v := input.(type) {
	case string:
		if isValidInputName(v) {
			valid = append(valid, v)
		} else {
			rejected = append(rejected, v)
		}
	case []interface{}:
		for _, item := range v {
			switch sv := item.(type) {
			case string:
				if isValidInputName(sv) {
					valid = append(valid, sv)
				} else {
					rejected = append(rejected, sv)
				}
			case map[string]interface{}:
				if name, ok := sv["name"].(string); ok {
					if isValidInputName(name) {
						valid = append(valid, name)
					} else {
						rejected = append(rejected, name)
					}
				}
			}
		}
	}
	return valid, rejected
}

// parseFrontmatter parses YAML-style frontmatter from markdown content.
// Returns (metadata map, body content). Simple implementation for the keys we care about.
func parseFrontmatter(content string) (map[string]interface{}, string) {
	meta := map[string]interface{}{}
	body := content

	if !strings.HasPrefix(content, "---") {
		return meta, body
	}
	// Find closing ---
	rest := content[3:]
	if rest != "" && rest[0] == '\n' {
		rest = rest[1:]
	}
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return meta, body
	}
	yamlPart := rest[:idx]
	body = rest[idx+4:]
	if strings.HasPrefix(body, "\n") {
		body = body[1:]
	}

	// Parse simple key: value lines
	for _, line := range strings.Split(yamlPart, "\n") {
		if colonIdx := strings.Index(line, ":"); colonIdx > 0 {
			key := strings.TrimSpace(line[:colonIdx])
			val := strings.TrimSpace(line[colonIdx+1:])
			// Remove surrounding quotes
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}
			meta[key] = val
		}
	}
	return meta, body
}

// buildCommandContent builds the command file content from metadata and body.
func buildCommandContent(meta map[string]interface{}, body string) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	orderedKeys := []string{"description", "allowed-tools", "model", "argument-hint", "arguments"}
	written := map[string]bool{}
	for _, k := range orderedKeys {
		if v, ok := meta[k]; ok {
			sb.WriteString(k)
			sb.WriteString(": ")
			switch sv := v.(type) {
			case string:
				sb.WriteString(sv)
			case []string:
				sb.WriteString("\n")
				for _, item := range sv {
					sb.WriteString("  - ")
					sb.WriteString(item)
					sb.WriteString("\n")
				}
				written[k] = true
				continue
			default:
				sb.WriteString("")
			}
			sb.WriteString("\n")
			written[k] = true
		}
	}
	sb.WriteString("---\n")
	sb.WriteString(body)
	return sb.String()
}

// transformPromptToCommand transforms a .prompt.md file into Claude command format.
// Returns (commandName, fileContent, droppedKeys bool).
func transformPromptToCommand(sourceFile string) (string, string, bool, error) {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", "", false, err
	}
	content := string(data)
	meta, body := parseFrontmatter(content)

	filename := filepath.Base(sourceFile)
	commandName := strings.TrimSuffix(filename, ".prompt.md")
	if commandName == filename {
		commandName = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	claudeMeta := map[string]interface{}{}

	if v, ok := meta["description"]; ok {
		claudeMeta["description"] = v
	}
	if v, ok := meta["allowed-tools"]; ok {
		claudeMeta["allowed-tools"] = v
	} else if v, ok := meta["allowedTools"]; ok {
		claudeMeta["allowed-tools"] = v
	}
	if v, ok := meta["model"]; ok {
		claudeMeta["model"] = v
	}
	if v, ok := meta["argument-hint"]; ok {
		claudeMeta["argument-hint"] = v
	} else if v, ok := meta["argumentHint"]; ok {
		claudeMeta["argument-hint"] = v
	}

	// Map 'input' to 'arguments' and 'argument-hint'
	inputNames, _ := extractInputNames(meta["input"])
	if len(inputNames) > 0 {
		claudeMeta["arguments"] = inputNames
		if _, ok := claudeMeta["argument-hint"]; !ok {
			hints := make([]string, len(inputNames))
			for i, n := range inputNames {
				hints[i] = "<" + n + ">"
			}
			claudeMeta["argument-hint"] = strings.Join(hints, " ")
		}
		// Replace ${{input:name}} with $name
		body = inputRefRe.ReplaceAllStringFunc(body, func(m string) string {
			sub := inputRefRe.FindStringSubmatch(m)
			if len(sub) > 1 {
				return "$" + sub[1]
			}
			return m
		})
	}

	// Compute dropped keys
	droppedKeys := false
	for k := range meta {
		if !preservedCommandKeys[k] {
			droppedKeys = true
			break
		}
	}

	fileContent := buildCommandContent(claudeMeta, body)
	return commandName, fileContent, droppedKeys, nil
}

// writeGeminiCommand transforms a .prompt.md to Gemini CLI TOML format.
func writeGeminiCommand(sourceFile, targetFile string) error {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	meta, body := parseFrontmatter(string(data))
	description, _ := meta["description"].(string)
	promptText := strings.TrimSpace(body)
	promptText = strings.ReplaceAll(promptText, "$ARGUMENTS", "{{args}}")

	var sb strings.Builder
	if description != "" {
		sb.WriteString("description = ")
		sb.WriteString(`"`)
		sb.WriteString(strings.ReplaceAll(description, `"`, `\"`))
		sb.WriteString(`"`)
		sb.WriteString("\n")
	}
	sb.WriteString("prompt = ")
	sb.WriteString(`"""`)
	sb.WriteString("\n")
	sb.WriteString(promptText)
	sb.WriteString("\n")
	sb.WriteString(`"""`)
	sb.WriteString("\n")

	if err := os.MkdirAll(filepath.Dir(targetFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(targetFile, []byte(sb.String()), 0o644)
}

// CommandIntegrator handles integration of .prompt.md files as slash commands.
type CommandIntegrator struct {
	passthroughNotified map[string]bool
}

// New returns a new CommandIntegrator.
func New() *CommandIntegrator {
	return &CommandIntegrator{
		passthroughNotified: map[string]bool{},
	}
}

// FindPromptFiles returns all .prompt.md files in a package.
func FindPromptFiles(packagePath string) []string {
	return baseintegrator.FindFilesByGlob(packagePath, "*.prompt.md", []string{".apm/prompts"})
}

// IntegrateCommandsForTarget integrates prompt files as commands for a single target.
func (ci *CommandIntegrator) IntegrateCommandsForTarget(
	tgt *targets.TargetProfile,
	packageInstallPath, projectRoot string,
	force bool,
	managedFiles map[string]struct{},
	diag baseintegrator.Diagnostics,
) IntegrationResult {
	mapping, ok := tgt.Primitives["commands"]
	if !ok {
		return IntegrationResult{}
	}

	effectiveRoot := mapping.DeployRoot
	if effectiveRoot == "" {
		effectiveRoot = tgt.RootDir
	}
	if !tgt.AutoCreate {
		if _, err := os.Stat(filepath.Join(projectRoot, tgt.RootDir)); err != nil {
			return IntegrationResult{}
		}
	}

	promptFiles := FindPromptFiles(packageInstallPath)
	if len(promptFiles) == 0 {
		return IntegrationResult{}
	}

	commandsDir := filepath.Join(projectRoot, effectiveRoot, mapping.Subdir)
	var result IntegrationResult
	anyDroppedKeys := false

	for _, promptFile := range promptFiles {
		filename := filepath.Base(promptFile)
		baseName := strings.TrimSuffix(filename, ".prompt.md")
		if baseName == filename {
			baseName = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		// Path security check
		if strings.Contains(baseName, "..") || strings.ContainsAny(baseName, "/\\") {
			result.FilesSkipped++
			continue
		}

		ext := mapping.Extension
		if ext == "" {
			ext = ".md"
		}
		targetPath := filepath.Join(commandsDir, baseName+ext)
		relPath := strings.ReplaceAll(func() string {
			rel, _ := filepath.Rel(projectRoot, targetPath)
			return rel
		}(), "\\", "/")

		if baseintegrator.CheckCollision(targetPath, relPath, managedFiles, force, diag) {
			result.FilesSkipped++
			continue
		}

		var written bool
		var hadDropped bool
		if mapping.FormatID == "gemini_command" {
			if err := writeGeminiCommand(promptFile, targetPath); err == nil {
				written = true
			}
			hadDropped = false
		} else {
			commandName, fileContent, dropped, err := transformPromptToCommand(promptFile)
			_ = commandName
			if err != nil {
				result.FilesSkipped++
				continue
			}
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				result.FilesSkipped++
				continue
			}
			if err := os.WriteFile(targetPath, []byte(fileContent), 0o644); err != nil {
				result.FilesSkipped++
				continue
			}
			written = true
			hadDropped = dropped
		}

		if !written {
			result.FilesSkipped++
			continue
		}
		if hadDropped {
			anyDroppedKeys = true
		}
		result.FilesIntegrated++
		result.TargetPaths = append(result.TargetPaths, targetPath)
	}
	_ = anyDroppedKeys
	return result
}

// SyncForTarget removes APM-managed command files for a single target.
func (ci *CommandIntegrator) SyncForTarget(
	tgt *targets.TargetProfile,
	projectRoot string,
	managedFiles map[string]struct{},
) map[string]int {
	mapping, ok := tgt.Primitives["commands"]
	if !ok {
		return map[string]int{"files_removed": 0, "errors": 0}
	}
	effectiveRoot := mapping.DeployRoot
	if effectiveRoot == "" {
		effectiveRoot = tgt.RootDir
	}
	prefix := effectiveRoot + "/" + mapping.Subdir + "/"
	legacyDir := filepath.Join(projectRoot, effectiveRoot, mapping.Subdir)

	res := baseintegrator.SyncRemoveFiles(
		projectRoot,
		managedFiles,
		prefix,
		legacyDir,
		"*-apm.md",
		nil,
		nil,
	)
	return map[string]int{"files_removed": res.FilesRemoved, "errors": res.Errors}
}

// IntegratePackageCommands integrates prompt files as Claude commands (legacy API).
func (ci *CommandIntegrator) IntegratePackageCommands(
	packageInstallPath, projectRoot string,
	force bool,
	managedFiles map[string]struct{},
	diag baseintegrator.Diagnostics,
) IntegrationResult {
	tgt, ok := targets.KnownTargets["claude"]
	if !ok {
		return IntegrationResult{}
	}
	_ = os.MkdirAll(filepath.Join(projectRoot, ".claude"), 0o755)
	return ci.IntegrateCommandsForTarget(tgt, packageInstallPath, projectRoot, force, managedFiles, diag)
}

// SyncIntegration removes APM-managed command files from .claude/commands/ (legacy).
func (ci *CommandIntegrator) SyncIntegration(
	projectRoot string,
	managedFiles map[string]struct{},
) map[string]int {
	tgt, ok := targets.KnownTargets["claude"]
	if !ok {
		return map[string]int{"files_removed": 0, "errors": 0}
	}
	return ci.SyncForTarget(tgt, projectRoot, managedFiles)
}
