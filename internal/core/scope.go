package core

import "path/filepath"

// InstallScope controls where packages are deployed.
type InstallScope int

const (
	// ScopeProject deploys to the current working directory (default).
	ScopeProject InstallScope = iota
	// ScopeUser deploys to user-level directories (~/.apm/).
	ScopeUser
)

// userAPMDir is the directory under $HOME for user-scope metadata.
const userAPMDir = ".apm"

// GetDeployRoot returns the root directory used to construct deployment paths.
func GetDeployRoot(scope InstallScope, cwd, home string) string {
	if scope == ScopeUser {
		return home
	}
	return cwd
}

// GetAPMDir returns the directory that holds APM metadata.
func GetAPMDir(scope InstallScope, cwd, home string) string {
	if scope == ScopeUser {
		return filepath.Join(home, userAPMDir)
	}
	return cwd
}

// GetModulesDir returns the apm_modules directory for scope.
func GetModulesDir(scope InstallScope, cwd, home, apmModulesDir string) string {
	return filepath.Join(GetAPMDir(scope, cwd, home), apmModulesDir)
}

// GetManifestPath returns the apm.yml path for scope.
func GetManifestPath(scope InstallScope, cwd, home, apmYMLFilename string) string {
	return filepath.Join(GetAPMDir(scope, cwd, home), apmYMLFilename)
}

// GetLockfileDir returns the directory containing the lockfile for scope.
func GetLockfileDir(scope InstallScope, cwd, home string) string {
	return GetAPMDir(scope, cwd, home)
}

// ParseScope parses a scope string ("project" or "user") into an InstallScope.
// Returns ScopeProject and false for unknown values.
func ParseScope(s string) (InstallScope, bool) {
	switch s {
	case "user":
		return ScopeUser, true
	case "project", "":
		return ScopeProject, true
	}
	return ScopeProject, false
}

// String returns the string representation of the scope.
func (s InstallScope) String() string {
	if s == ScopeUser {
		return "user"
	}
	return "project"
}
