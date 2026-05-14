// Package windsurf provides the Windsurf/Cascade MCP client adapter.
// Migrated from src/apm_cli/adapters/client/windsurf.py
//
// Windsurf uses the standard mcpServers JSON format at
// ~/.codeium/windsurf/mcp_config.json (global). The config schema is
// identical to GitHub Copilot CLI.
package windsurf

import (
	"os"
	"path/filepath"
)

// Adapter implements the Windsurf/Cascade MCP client adapter.
type Adapter struct {
	// SupportsUserScope indicates this adapter targets global user config.
	SupportsUserScope bool
	// ClientLabel is the user-facing label for this adapter.
	ClientLabel string
	// TargetName is the adapter identifier.
	TargetName string
	// MCPServersKey is the JSON key for MCP servers.
	MCPServersKey string
	// SupportsRuntimeEnvSubstitution mirrors the Python field.
	// Pinned to false until windsurf runtime-env audit is complete.
	SupportsRuntimeEnvSubstitution bool
}

// New returns a new Windsurf adapter with default settings.
func New() *Adapter {
	return &Adapter{
		SupportsUserScope:              true,
		ClientLabel:                    "Windsurf",
		TargetName:                     "windsurf",
		MCPServersKey:                  "mcpServers",
		SupportsRuntimeEnvSubstitution: false,
	}
}

// GetConfigPath returns the path to ~/.codeium/windsurf/mcp_config.json.
// This is a global config path -- Windsurf reads MCP server definitions
// from the user-level directory, not the workspace.
func (a *Adapter) GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	return filepath.Join(home, ".codeium", "windsurf", "mcp_config.json")
}

// GetRuntimeName returns the runtime name.
func (a *Adapter) GetRuntimeName() string { return a.TargetName }

// IsAvailable always returns true for Windsurf (file-based config, no binary check).
func (a *Adapter) IsAvailable() bool { return true }
