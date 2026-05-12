// Package helpers provides miscellaneous utility functions for APM.
package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsToolAvailable reports whether a command-line tool can be found on PATH.
func IsToolAvailable(toolName string) bool {
	_, err := exec.LookPath(toolName)
	return err == nil
}

// GetAvailablePackageManagers returns a map of package manager name -> name
// for every package manager binary found on PATH.
func GetAvailablePackageManagers() map[string]string {
	candidates := []string{
		// Python
		"uv", "pip", "pipx",
		// JavaScript
		"npm", "yarn", "pnpm",
		// System
		"brew",   // macOS
		"apt",    // Debian/Ubuntu
		"yum",    // CentOS/RHEL
		"dnf",    // Fedora
		"apk",    // Alpine
		"pacman", // Arch
	}
	out := make(map[string]string)
	for _, name := range candidates {
		if IsToolAvailable(name) {
			out[name] = name
		}
	}
	return out
}

// DetectPlatform returns a normalised platform name: "macos", "linux",
// "windows", or "unknown".
func DetectPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}

// pluginJSONRelPaths is the ordered list of relative paths where plugin.json
// may live inside a plugin directory.
var pluginJSONRelPaths = []string{
	"plugin.json",
	filepath.Join(".github", "plugin", "plugin.json"),
	filepath.Join(".claude-plugin", "plugin.json"),
	filepath.Join(".cursor-plugin", "plugin.json"),
}

// FindPluginJSON searches for plugin.json in the well-known locations inside
// pluginPath and returns the first match. Returns an empty string when not found.
func FindPluginJSON(pluginPath string) string {
	for _, rel := range pluginJSONRelPaths {
		candidate := filepath.Join(pluginPath, rel)
		if fileExists(candidate) {
			return candidate
		}
	}
	return ""
}

// fileExists reports whether path refers to a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
