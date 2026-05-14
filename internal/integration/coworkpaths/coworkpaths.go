// Package coworkpaths handles OneDrive-backed Cowork skills directory resolution
// and lockfile path translation.
// Ported from src/apm_cli/integration/copilot_cowork_paths.py
package coworkpaths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CoworkURIScheme is the synthetic URI prefix used in lockfile entries.
const CoworkURIScheme = "cowork://"

// CoworkLockfilePrefix is the full prefix for skill entries in the lockfile.
const CoworkLockfilePrefix = "cowork://skills/"

const oneDriveGlob = "OneDrive*"
const coworkSubdir = "Documents/Cowork"
const coworkSkillsSubdir = "Documents/Cowork/skills"

// CoworkResolutionError is raised when OneDrive resolution fails.
type CoworkResolutionError struct {
	Msg string
}

func (e *CoworkResolutionError) Error() string { return e.Msg }

// ResolveCoworkSkillsDir locates the Cowork skills directory on the current machine.
// Resolution order:
//  1. APM_COPILOT_COWORK_SKILLS_DIR env var
//  2. Platform auto-detection (macOS, Windows)
//
// Returns empty string when no OneDrive mount is found.
func ResolveCoworkSkillsDir() (string, error) {
	if override := os.Getenv("APM_COPILOT_COWORK_SKILLS_DIR"); override != "" {
		if err := validatePathSegments(override, "APM_COPILOT_COWORK_SKILLS_DIR"); err != nil {
			return "", &CoworkResolutionError{
				Msg: fmt.Sprintf("APM_COPILOT_COWORK_SKILLS_DIR contains a traversal sequence: %v", err),
			}
		}
		abs, err := filepath.Abs(override)
		if err != nil {
			return "", err
		}
		return abs, nil
	}

	switch runtime.GOOS {
	case "windows":
		for _, envName := range []string{"ONEDRIVECOMMERCIAL", "ONEDRIVE"} {
			winRoot := os.Getenv(envName)
			if winRoot != "" {
				winSkills := filepath.Join(winRoot, filepath.FromSlash(coworkSkillsSubdir))
				if err := validatePathSegments(winSkills, envName+" env var"); err != nil {
					return "", &CoworkResolutionError{
						Msg: fmt.Sprintf("%s contains a traversal sequence: %v", envName, err),
					}
				}
				abs, err := filepath.Abs(winSkills)
				if err != nil {
					return "", err
				}
				return abs, nil
			}
		}
		return "", nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cloudStorage := filepath.Join(home, "Library", "CloudStorage")
		info, err := os.Stat(cloudStorage)
		if err != nil || !info.IsDir() {
			return "", nil
		}
		entries, err := os.ReadDir(cloudStorage)
		if err != nil {
			return "", nil
		}
		var candidates []string
		for _, e := range entries {
			if matchOneDriveGlob(e.Name()) {
				candidates = append(candidates, filepath.Join(cloudStorage, e.Name()))
			}
		}
		if len(candidates) == 0 {
			return "", nil
		}
		if len(candidates) > 1 {
			listing := ""
			for _, c := range candidates {
				listing += fmt.Sprintf("  - %s\n", c)
			}
			suggestion := filepath.Join(candidates[0], filepath.FromSlash(coworkSkillsSubdir))
			return "", &CoworkResolutionError{
				Msg: fmt.Sprintf("Multiple OneDrive mounts detected:\n%s"+
					"Set APM_COPILOT_COWORK_SKILLS_DIR to the desired skills directory, e.g.:\n"+
					"  export APM_COPILOT_COWORK_SKILLS_DIR=\"%s\"",
					listing, suggestion),
			}
		}
		return filepath.Join(candidates[0], filepath.FromSlash(coworkSkillsSubdir)), nil
	default:
		return "", nil
	}
}

// matchOneDriveGlob returns true if name matches "OneDrive*".
func matchOneDriveGlob(name string) bool {
	return strings.HasPrefix(name, "OneDrive")
}

// validatePathSegments rejects traversal sequences in a path string.
func validatePathSegments(p string, context string) error {
	parts := strings.Split(filepath.ToSlash(p), "/")
	for _, part := range parts {
		if part == ".." {
			return fmt.Errorf("%s: path contains '..' segment", context)
		}
	}
	return nil
}

// ToLockfilePath encodes an absolute cowork path as a cowork:// lockfile entry.
func ToLockfilePath(absolute string, coworkRoot string) (string, error) {
	absResolved, err := filepath.Abs(absolute)
	if err != nil {
		return "", err
	}
	rootResolved, err := filepath.Abs(coworkRoot)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absResolved, rootResolved+string(filepath.Separator)) &&
		absResolved != rootResolved {
		return "", errors.New("path escapes cowork root")
	}
	rel, err := filepath.Rel(rootResolved, absResolved)
	if err != nil {
		return "", err
	}
	return CoworkURIScheme + "skills/" + filepath.ToSlash(rel), nil
}

// FromLockfilePath decodes a cowork:// lockfile entry to an absolute path.
func FromLockfilePath(lockfilePath string, coworkRoot string) (string, error) {
	if !strings.HasPrefix(lockfilePath, CoworkURIScheme) {
		return "", fmt.Errorf("not a cowork lockfile path: %q", lockfilePath)
	}
	relPosix := lockfilePath[len(CoworkURIScheme):]
	if err := validatePathSegments(relPosix, "cowork lockfile path"); err != nil {
		return "", err
	}
	skillsPrefix := "skills/"
	if strings.HasPrefix(relPosix, skillsPrefix) {
		relPosix = relPosix[len(skillsPrefix):]
	}
	candidate := filepath.Join(coworkRoot, filepath.FromSlash(relPosix))
	rootResolved, err := filepath.Abs(coworkRoot)
	if err != nil {
		return "", err
	}
	candidateResolved, err := filepath.Abs(candidate)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(candidateResolved, rootResolved+string(filepath.Separator)) &&
		candidateResolved != rootResolved {
		return "", errors.New("decoded path escapes cowork root")
	}
	return candidateResolved, nil
}

// IsCoworkPath returns true if lockfilePath uses the cowork:// scheme.
func IsCoworkPath(lockfilePath string) bool {
	return strings.HasPrefix(lockfilePath, CoworkURIScheme)
}
