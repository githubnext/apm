// Package pm provides MCP package manager adapter interfaces.
package pm

import "errors"

// MCPPackageManagerAdapter is the base interface for MCP package managers.
type MCPPackageManagerAdapter interface {
	// Install installs an MCP package.
	Install(packageName string, version string) error
	// Uninstall removes an installed MCP package.
	Uninstall(packageName string) error
	// ListInstalled lists all installed MCP packages.
	ListInstalled() ([]string, error)
	// Search queries available MCP packages.
	Search(query string) ([]string, error)
}

// ErrPackageNotFound is returned when a package is not installed.
var ErrPackageNotFound = errors.New("MCP package not found")

// ErrInstallFailed is returned when package installation fails.
var ErrInstallFailed = errors.New("MCP package install failed")
