// Package instructionintegrator deploys .instructions.md files from APM packages
// to the appropriate target directory with format-specific transforms.
//
// Supported format transforms:
//   - cursor_rules:   applyTo: -> globs: (.mdc extension)
//   - claude_rules:   applyTo: -> paths: list
//   - windsurf_rules: applyTo: -> trigger: glob + globs:
//   - default:        verbatim copy
package instructionintegrator

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IntegrationResult holds the result of an instruction integration operation.
type IntegrationResult struct {
	FilesIntegrated int
	FilesUpdated    int
	FilesSkipped    int
	TargetPaths     []string
	LinksResolved   int
}

// FormatID identifies the content transform to apply.
type FormatID string

const (
	FormatVerbatim      FormatID = ""
	FormatCursorRules   FormatID = "cursor_rules"
	FormatClaudeRules   FormatID = "claude_rules"
	FormatWindsurfRules FormatID = "windsurf_rules"
)

// TargetConfig holds deploy configuration for an integration target.
type TargetConfig struct {
	// RootDir is the target root (e.g. ".github").
	RootDir string
	// Subdir is the subdirectory under RootDir for the primitive.
	Subdir string
	// Extension is the file extension for renamed files (e.g. ".mdc").
	Extension string
	// FormatID selects the content transform.
	FormatID FormatID
	// DeployRoot overrides RootDir when set.
	DeployRoot string
	// AutoCreate creates the target directory even if RootDir doesn't exist.
	AutoCreate bool
}

// frontmatterRe matches a YAML frontmatter block at the top of a file.
var frontmatterRe = regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---\s*\n?`)

// parseFrontmatter extracts applyTo and description from YAML frontmatter.
func parseFrontmatter(content string) (applyTo, description, body string) {
	m := frontmatterRe.FindStringSubmatchIndex(content)
	if m == nil {
		return "", "", content
	}
	fmBlock := content[m[2]:m[3]]
	body = content[m[1]:]
	for _, line := range strings.Split(fmBlock, "\n") {
		stripped := strings.TrimSpace(line)
		if strings.HasPrefix(stripped, "applyTo:") {
			applyTo = strings.Trim(strings.TrimPrefix(stripped, "applyTo:"), " '\"")
		} else if strings.HasPrefix(stripped, "description:") {
			description = strings.Trim(strings.TrimPrefix(stripped, "description:"), " '\"")
		}
	}
	return applyTo, description, body
}

// ConvertToCursorRules converts APM instruction content to Cursor Rules .mdc format.
// Maps applyTo: -> globs: and extracts or generates description.
func ConvertToCursorRules(content string) string {
	applyTo, description, body := parseFrontmatter(content)

	if description == "" {
		for _, line := range strings.Split(body, "\n") {
			stripped := strings.TrimLeft(strings.TrimSpace(line), "#")
			stripped = strings.TrimSpace(stripped)
			if stripped != "" {
				parts := strings.SplitN(stripped, ".", 2)
				description = strings.TrimSpace(parts[0])
				break
			}
		}
	}

	var parts []string
	parts = append(parts, "---")
	if description != "" {
		parts = append(parts, "description: "+description)
	}
	if applyTo != "" {
		parts = append(parts, `globs: "`+applyTo+`"`)
	}
	parts = append(parts, "---")

	return strings.Join(parts, "\n") + "\n\n" + strings.TrimLeft(body, "\n")
}

// ConvertToClaudeRules converts APM instruction content to Claude Code rules .md format.
// Maps applyTo: -> paths: list. Instructions without applyTo become unconditional rules.
func ConvertToClaudeRules(content string) string {
	applyTo, _, body := parseFrontmatter(content)

	if applyTo != "" {
		fm := "---\npaths:\n  - \"" + applyTo + "\"\n---"
		return fm + "\n\n" + strings.TrimLeft(body, "\n")
	}
	return strings.TrimLeft(body, "\n")
}

// ConvertToWindsurfRules converts APM instruction content to Windsurf rules .md format.
// Maps applyTo: -> trigger: glob + globs:. Instructions without applyTo use trigger: always_on.
func ConvertToWindsurfRules(content string) string {
	applyTo, _, body := parseFrontmatter(content)

	var parts []string
	parts = append(parts, "---")
	if applyTo != "" {
		safeApplyTo := strings.ReplaceAll(strings.ReplaceAll(applyTo, "\n", " "), "\r", " ")
		safeApplyTo = strings.TrimSpace(safeApplyTo)
		parts = append(parts, "trigger: glob")
		parts = append(parts, `globs: "`+safeApplyTo+`"`)
	} else {
		parts = append(parts, "trigger: always_on")
	}
	parts = append(parts, "---")

	return strings.Join(parts, "\n") + "\n\n" + strings.TrimLeft(body, "\n")
}

// FindInstructionFiles returns all .instructions.md files in a package's .apm/instructions/ dir.
func FindInstructionFiles(packagePath string) ([]string, error) {
	var files []string
	instructionsDir := filepath.Join(packagePath, ".apm", "instructions")
	_ = filepath.WalkDir(instructionsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".instructions.md") {
			files = append(files, path)
		}
		return nil
	})
	return files, nil
}

// CopyInstruction copies an instruction file to target, applying the given format transform.
// Returns number of links resolved (always 0 in this stdlib implementation).
func CopyInstruction(source, target string, format FormatID) (int, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return 0, err
	}
	content := string(data)

	switch format {
	case FormatCursorRules:
		content = ConvertToCursorRules(content)
	case FormatClaudeRules:
		content = ConvertToClaudeRules(content)
	case FormatWindsurfRules:
		content = ConvertToWindsurfRules(content)
	}

	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return 0, err
	}
	return 0, nil
}

// IntegrateInstructionsForTarget deploys instruction files to the given target directory.
func IntegrateInstructionsForTarget(
	installPath string,
	projectRoot string,
	cfg TargetConfig,
	force bool,
	managedFiles map[string]bool,
) (IntegrationResult, error) {
	result := IntegrationResult{}

	effectiveRoot := cfg.DeployRoot
	if effectiveRoot == "" {
		effectiveRoot = cfg.RootDir
	}

	if !cfg.AutoCreate {
		if _, err := os.Stat(filepath.Join(projectRoot, cfg.RootDir)); os.IsNotExist(err) {
			return result, nil
		}
	}

	instructionFiles, err := FindInstructionFiles(installPath)
	if err != nil {
		return result, err
	}
	if len(instructionFiles) == 0 {
		return result, nil
	}

	deployDir := filepath.Join(projectRoot, effectiveRoot, cfg.Subdir)
	if err := os.MkdirAll(deployDir, 0o755); err != nil {
		return result, err
	}

	needsRename := cfg.FormatID == FormatCursorRules ||
		cfg.FormatID == FormatClaudeRules ||
		cfg.FormatID == FormatWindsurfRules

	for _, src := range instructionFiles {
		var targetName string
		if needsRename {
			stem := filepath.Base(src)
			if strings.HasSuffix(stem, ".instructions.md") {
				stem = stem[:len(stem)-len(".instructions.md")]
			}
			ext := cfg.Extension
			if ext == "" {
				ext = ".md"
			}
			targetName = stem + ext
		} else {
			targetName = filepath.Base(src)
		}

		targetPath := filepath.Join(deployDir, targetName)
		relPath := filepath.ToSlash(strings.TrimPrefix(targetPath, projectRoot+string(filepath.Separator)))

		if checkCollision(targetPath, relPath, managedFiles, force) {
			result.FilesSkipped++
			continue
		}

		links, err := CopyInstruction(src, targetPath, cfg.FormatID)
		if err != nil {
			return result, err
		}
		result.FilesIntegrated++
		result.LinksResolved += links
		result.TargetPaths = append(result.TargetPaths, targetPath)
	}

	return result, nil
}

// SyncForTarget removes APM-managed instruction files for a given target.
func SyncForTarget(
	projectRoot string,
	cfg TargetConfig,
	managedFiles map[string]bool,
) (filesRemoved int, errors int) {
	effectiveRoot := cfg.DeployRoot
	if effectiveRoot == "" {
		effectiveRoot = cfg.RootDir
	}
	prefix := effectiveRoot + "/" + cfg.Subdir + "/"

	if managedFiles != nil {
		for rel := range managedFiles {
			if strings.HasPrefix(rel, prefix) {
				abs := filepath.Join(projectRoot, filepath.FromSlash(rel))
				if rmErr := os.Remove(abs); rmErr == nil {
					filesRemoved++
				}
			}
		}
		return filesRemoved, errors
	}

	// Legacy glob removal
	var legacyPattern string
	switch cfg.FormatID {
	case FormatCursorRules:
		legacyPattern = "*.mdc"
	case FormatWindsurfRules, FormatClaudeRules:
		// Avoid broad deletion of user-authored .md files
		return 0, 0
	default:
		legacyPattern = "*.instructions.md"
	}

	legacyDir := filepath.Join(projectRoot, effectiveRoot, cfg.Subdir)
	entries, err := os.ReadDir(legacyDir)
	if err != nil {
		return 0, 0
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		matched, _ := filepath.Match(legacyPattern, e.Name())
		if matched {
			if rmErr := os.Remove(filepath.Join(legacyDir, e.Name())); rmErr == nil {
				filesRemoved++
			}
		}
	}
	return filesRemoved, errors
}

// checkCollision returns true if the target is a user-authored file that should not be overwritten.
func checkCollision(targetPath, relPath string, managedFiles map[string]bool, force bool) bool {
	if managedFiles == nil {
		return false
	}
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return false
	}
	normalized := strings.ReplaceAll(relPath, "\\", "/")
	if managedFiles[normalized] {
		return false
	}
	if force {
		return false
	}
	return true
}
