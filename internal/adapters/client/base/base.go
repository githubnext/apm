// Package base defines the MCPClientAdapter interface and shared regex helpers.
//
// Mirrors src/apm_cli/adapters/client/base.py.
package base

import (
	"regexp"
)

// InputVarRE matches ${input:NAME} placeholders that adapters must warn about.
var InputVarRE = regexp.MustCompile(`\$\{input:([^}]+)\}`)

// EnvVarRE matches ${VAR} and ${env:VAR}, capturing the variable name.
// Does NOT match ${input:VAR} or GitHub Actions ${{ ... }}.
var EnvVarRE = regexp.MustCompile(`\$\{(?:env:)?([A-Za-z_][A-Za-z0-9_]*)\}`)

// MCPClientAdapter is the interface all MCP client adapters must satisfy.
type MCPClientAdapter interface {
	// GetConfigPath returns the path to this adapter's config file.
	GetConfigPath() string

	// UpdateConfig merges config_updates into the adapter's config file.
	UpdateConfig(configUpdates map[string]interface{}) error

	// GetCurrentConfig reads and returns the current config, or empty map on error.
	GetCurrentConfig() map[string]interface{}

	// ConfigureMCPServer installs a single MCP server into the adapter config.
	// Returns true on success.
	ConfigureMCPServer(serverURL, serverName string, enabled bool,
		envOverrides, serverInfoCache map[string]interface{},
		runtimeVars map[string]string) bool

	// FormatServerConfig converts registry server info to the adapter's wire format.
	FormatServerConfig(serverInfo map[string]interface{},
		envOverrides map[string]interface{},
		runtimeVars map[string]string) (map[string]interface{}, error)

	// TargetName returns the canonical adapter target name (e.g. "copilot", "vscode").
	TargetName() string

	// MCPServersKey returns the top-level JSON key for server entries.
	MCPServersKey() string

	// SupportsUserScope reports whether this adapter has a user/global config scope.
	SupportsUserScope() bool
}
