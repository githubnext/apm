package scriptrunner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PromptCompiler compiles .prompt.md files with parameter substitution.
type PromptCompiler struct {
	CompiledDir string
}

const defaultCompiledDir = ".apm/compiled"

// NewPromptCompiler returns a PromptCompiler with default settings.
func NewPromptCompiler() *PromptCompiler {
	return &PromptCompiler{CompiledDir: defaultCompiledDir}
}

// Compile compiles a .prompt.md file with parameter substitution.
// Returns the path to the compiled .txt file.
func (c *PromptCompiler) Compile(promptFile string, params map[string]string) (string, error) {
	promptPath, err := c.resolvePromptFile(promptFile)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(c.CompiledDir, 0o755); err != nil {
		return "", fmt.Errorf("creating compiled dir: %w", err)
	}

	data, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("reading prompt file: %w", err)
	}

	content := string(data)

	// Strip YAML frontmatter if present.
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			content = strings.TrimSpace(parts[2])
		}
	}

	compiled := substituteParameters(content, params)

	// Build output file name: strip .prompt from stem, add .txt.
	base := filepath.Base(promptPath)
	stem := strings.TrimSuffix(base, filepath.Ext(base)) // removes .md
	stem = strings.TrimSuffix(stem, ".prompt")           // removes .prompt
	outputName := stem + ".txt"
	outputPath := filepath.Join(c.CompiledDir, outputName)

	if err := os.WriteFile(outputPath, []byte(compiled), 0o644); err != nil {
		return "", fmt.Errorf("writing compiled file: %w", err)
	}

	return outputPath, nil
}

// resolvePromptFile locates the .prompt.md file checking local dirs then dependencies.
func (c *PromptCompiler) resolvePromptFile(promptFile string) (string, error) {
	promptPath := promptFile

	// Reject symlinks.
	if fi, err := os.Lstat(promptPath); err == nil {
		if fi.Mode()&fs.ModeSymlink != 0 {
			return "", fmt.Errorf("prompt file '%s' is a symlink; symlinks are not allowed for security reasons", promptFile)
		}
		return promptPath, nil
	}

	// Common project directories.
	for _, dir := range []string{".github/prompts", ".apm/prompts"} {
		candidate := filepath.Join(dir, promptFile)
		fi, err := os.Lstat(candidate)
		if err == nil && fi.Mode()&fs.ModeSymlink == 0 {
			return candidate, nil
		}
	}

	// Search in apm_modules (two-level walk).
	apmModulesDir := "apm_modules"
	depDirs := collectDependencyDirs(apmModulesDir)

	for _, dep := range depDirs {
		for _, subdir := range []string{".", "prompts", "workflows"} {
			var candidate string
			if subdir == "." {
				candidate = filepath.Join(dep.repoDir, promptFile)
			} else {
				candidate = filepath.Join(dep.repoDir, subdir, promptFile)
			}
			fi, err := os.Lstat(candidate)
			if err == nil && fi.Mode()&fs.ModeSymlink == 0 {
				return candidate, nil
			}
		}
	}

	// Build error message.
	return "", c.buildNotFoundError(promptFile, depDirs)
}

type depDir struct {
	orgName  string
	repoName string
	repoDir  string
}

func collectDependencyDirs(apmModulesDir string) []depDir {
	if _, err := os.Stat(apmModulesDir); err != nil {
		return nil
	}
	var result []depDir
	orgEntries, err := os.ReadDir(apmModulesDir)
	if err != nil {
		return nil
	}
	for _, orgEntry := range orgEntries {
		if !orgEntry.IsDir() || strings.HasPrefix(orgEntry.Name(), ".") {
			continue
		}
		orgDir := filepath.Join(apmModulesDir, orgEntry.Name())
		repoEntries, err := os.ReadDir(orgDir)
		if err != nil {
			continue
		}
		for _, repoEntry := range repoEntries {
			if !repoEntry.IsDir() || strings.HasPrefix(repoEntry.Name(), ".") {
				continue
			}
			result = append(result, depDir{
				orgName:  orgEntry.Name(),
				repoName: repoEntry.Name(),
				repoDir:  filepath.Join(orgDir, repoEntry.Name()),
			})
		}
	}
	return result
}

func (c *PromptCompiler) buildNotFoundError(promptFile string, deps []depDir) error {
	locations := []string{
		"Local: " + promptFile,
		"GitHub prompts: .github/prompts/" + promptFile,
		"APM prompts: .apm/prompts/" + promptFile,
	}
	if len(deps) > 0 {
		locations = append(locations, "Dependencies:")
		for _, d := range deps {
			locations = append(locations, fmt.Sprintf("  - %s/%s/%s", d.orgName, d.repoName, promptFile))
		}
	}
	return fmt.Errorf(
		"Prompt file '%s' not found.\nSearched in:\n%s\n\nTip: Run 'apm install' to ensure dependencies are installed.",
		promptFile,
		strings.Join(locations, "\n"),
	)
}

// substituteParameters replaces ${input:key} placeholders in content.
func substituteParameters(content string, params map[string]string) string {
	result := content
	for key, value := range params {
		placeholder := "${input:" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
