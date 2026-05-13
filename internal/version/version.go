// Package version provides version resolution for APM CLI.
// Migrated from src/apm_cli/version.py
package version

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// BuildVersion is optionally injected at build time via -ldflags.
var BuildVersion string

// BuildSHA is optionally injected at build time via -ldflags.
var BuildSHA string

var versionRe = regexp.MustCompile(`version\s*=\s*["']([^"']+)["']`)
var pep440Re = regexp.MustCompile(`^\d+\.\d+\.\d+(a\d+|b\d+|rc\d+)?$`)

// GetVersion returns the current version string.
// Priority: build-time constant > pyproject.toml parse > "unknown".
func GetVersion() string {
	if BuildVersion != "" {
		return BuildVersion
	}
	// Locate pyproject.toml relative to this source file (dev mode).
	_, file, _, ok := runtime.Caller(0)
	if ok {
		repoRoot := filepath.Join(filepath.Dir(file), "..", "..", "..")
		pyproject := filepath.Join(repoRoot, "pyproject.toml")
		if v := versionFromPyproject(pyproject); v != "" {
			return v
		}
	}
	return "unknown"
}

func versionFromPyproject(path string) string {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return ""
	}
	m := versionRe.FindStringSubmatch(string(data))
	if m == nil {
		return ""
	}
	v := m[1]
	if !pep440Re.MatchString(v) {
		return ""
	}
	return v
}

// GetBuildSHA returns the short git commit SHA.
func GetBuildSHA() string {
	if BuildSHA != "" {
		return BuildSHA
	}
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
