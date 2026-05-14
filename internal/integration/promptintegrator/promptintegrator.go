// Package promptintegrator provides prompt file integration for APM packages.
// Deploys .prompt.md files into .github/prompts/.
package promptintegrator

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// IntegrationResult holds the result of a prompt integration operation.
type IntegrationResult struct {
	FilesIntegrated int
	FilesUpdated    int
	FilesSkipped    int
	TargetPaths     []string
	LinksResolved   int
}

// FindPromptFiles returns all .prompt.md files found in a package directory.
// Searches in package root and .apm/prompts/ subdirectory.
func FindPromptFiles(packagePath string) ([]string, error) {
	var files []string

	// Search in package root
	entries, err := os.ReadDir(packagePath)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".prompt.md") {
				files = append(files, filepath.Join(packagePath, e.Name()))
			}
		}
	}

	// Search in .apm/prompts/
	apmPrompts := filepath.Join(packagePath, ".apm", "prompts")
	_ = filepath.WalkDir(apmPrompts, func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".prompt.md") {
			files = append(files, path)
		}
		return nil
	})

	return files, nil
}

// GetTargetFilename returns the target filename for a prompt file (no suffix change).
func GetTargetFilename(sourceFile string) string {
	return filepath.Base(sourceFile)
}

// CopyPrompt copies a prompt file verbatim to the target path.
// Returns number of links resolved (always 0 in this implementation).
func CopyPrompt(source, target string) (int, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return 0, err
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return 0, err
	}
	return 0, nil
}

// IntegratePackagePrompts integrates all prompt files from a package into .github/prompts/.
// managedFiles is the set of relative paths known to be APM-managed (nil = legacy mode).
// force overrides collision checks.
func IntegratePackagePrompts(
	installPath string,
	projectRoot string,
	force bool,
	managedFiles map[string]bool,
) (IntegrationResult, error) {
	result := IntegrationResult{}

	promptFiles, err := FindPromptFiles(installPath)
	if err != nil {
		return result, err
	}
	if len(promptFiles) == 0 {
		return result, nil
	}

	promptsDir := filepath.Join(projectRoot, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		return result, err
	}

	for _, src := range promptFiles {
		targetName := GetTargetFilename(src)
		targetPath := filepath.Join(promptsDir, targetName)
		relPath := filepath.ToSlash(strings.TrimPrefix(targetPath, projectRoot+string(filepath.Separator)))

		if checkCollision(targetPath, relPath, managedFiles, force) {
			result.FilesSkipped++
			continue
		}

		links, err := CopyPrompt(src, targetPath)
		if err != nil {
			return result, err
		}
		result.FilesIntegrated++
		result.LinksResolved += links
		result.TargetPaths = append(result.TargetPaths, targetPath)
	}

	return result, nil
}

// SyncIntegration removes APM-managed prompt files.
// managedFiles nil => legacy glob removal of *-apm.prompt.md.
func SyncIntegration(
	projectRoot string,
	managedFiles map[string]bool,
) (filesRemoved int, errors int) {
	promptsDir := filepath.Join(projectRoot, ".github", "prompts")

	if managedFiles != nil {
		for rel := range managedFiles {
			if strings.HasPrefix(rel, ".github/prompts/") {
				abs := filepath.Join(projectRoot, filepath.FromSlash(rel))
				if rmErr := os.Remove(abs); rmErr == nil {
					filesRemoved++
				}
			}
		}
		return filesRemoved, errors
	}

	// Legacy: remove *-apm.prompt.md
	entries, err := os.ReadDir(promptsDir)
	if err != nil {
		return 0, 0
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), "-apm.prompt.md") {
			if rmErr := os.Remove(filepath.Join(promptsDir, e.Name())); rmErr == nil {
				filesRemoved++
			}
		}
	}
	return filesRemoved, errors
}

// checkCollision returns true if target_path is a user-authored file that should not be overwritten.
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
