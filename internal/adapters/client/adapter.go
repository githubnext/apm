// Package client provides MCP client adapter interfaces and types.
package client

import "errors"

// MCPClientAdapter is the base interface for MCP client adapters.
type MCPClientAdapter interface {
	// GetCurrentConfig returns the current MCP configuration.
	GetCurrentConfig() (map[string]interface{}, error)
	// UpdateConfig updates the MCP configuration.
	UpdateConfig(config map[string]interface{}) error
	// ConfigureMCPServer adds or updates a single MCP server entry.
	ConfigureMCPServer(serverName, packageName string, enabled bool) error
	// RemoveMCPServer removes a server entry from the configuration.
	RemoveMCPServer(serverName string) error
	// GetTargetName returns the adapter's target identifier.
	GetTargetName() string
}

// ErrServerNotFound is returned when a server is not in the config.
var ErrServerNotFound = errors.New("MCP server not found in config")

// ErrConfigInvalid is returned for malformed configurations.
var ErrConfigInvalid = errors.New("invalid MCP configuration")

// MCPServerEntry represents a single MCP server configuration entry.
type MCPServerEntry struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	Type    string            `json:"type,omitempty"`
}
