// Package scope defines installation scope resolution for APM packages.
// Ported from src/apm_cli/core/scope.py
package scope

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/constants"
)

// InstallScope controls where packages are deployed.
type InstallScope int

const (
	// ScopeProject deploys to the current working directory.
	ScopeProject InstallScope = iota
	// ScopeUser deploys to user-level directories (~/.apm/).
	ScopeUser
)

// UserAPMDir is the directory under $HOME for user-scope metadata.
const UserAPMDir = ".apm"

// String returns the string representation of the scope.
func (s InstallScope) String() string {
	if s == ScopeUser {
		return "user"
	}
	return "project"
}

// ParseScope parses a scope string into an InstallScope.
func ParseScope(s string) (InstallScope, bool) {
	switch strings.ToLower(s) {
	case "user":
		return ScopeUser, true
	case "project":
		return ScopeProject, true
	default:
		return ScopeProject, false
	}
}

// GetDeployRoot returns the root used to construct deployment paths.
// For project scope this is cwd; for user scope this is $HOME.
func GetDeployRoot(s InstallScope) (string, error) {
	if s == ScopeUser {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return home, nil
	}
	return os.Getwd()
}

// GetAPMDir returns the directory that holds APM metadata (manifest, lockfile, modules).
// Project scope: cwd. User scope: ~/.apm/.
func GetAPMDir(s InstallScope) (string, error) {
	if s == ScopeUser {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, UserAPMDir), nil
	}
	return os.Getwd()
}

// GetModulesDir returns the apm_modules directory for scope.
func GetModulesDir(s InstallScope) (string, error) {
	apmDir, err := GetAPMDir(s)
	if err != nil {
		return "", err
	}
	return filepath.Join(apmDir, constants.APMModulesDir), nil
}

// GetManifestPath returns the apm.yml path for scope.
func GetManifestPath(s InstallScope) (string, error) {
	apmDir, err := GetAPMDir(s)
	if err != nil {
		return "", err
	}
	return filepath.Join(apmDir, constants.APMYMLFilename), nil
}

// GetLockfileDir returns the directory containing the lockfile for scope.
func GetLockfileDir(s InstallScope) (string, error) {
	return GetAPMDir(s)
}

// EnsureUserDirs creates ~/.apm/ and ~/.apm/apm_modules/ if they do not exist.
// Returns the user APM root (~/.apm/).
func EnsureUserDirs() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	userRoot := filepath.Join(home, UserAPMDir)
	if err := os.MkdirAll(userRoot, 0o755); err != nil {
		return "", err
	}
	modsDir := filepath.Join(userRoot, constants.APMModulesDir)
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		return "", err
	}
	return userRoot, nil
}
