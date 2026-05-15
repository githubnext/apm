// Package opencode provides the OpenCode MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/opencode.py.
//
// OpenCode uses opencode.json at the project root with an "mcp" key.
// Schema:
//
//	{
//	  "mcp": {
//	    "server-name": {
//	      "type": "local",
//	      "command": ["npx", "-y", "@modelcontextprotocol/server-foo"],
//	      "environment": { "KEY": "value" },
//	      "enabled": true
//	    }
//	  }
//	}
//
// APM only writes to opencode.json when .opencode/ already exists (opt-in).
package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ServerEntry represents an OpenCode MCP server config entry.
type ServerEntry struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Enabled     bool              `json:"enabled"`
}

// CopilotEntry represents a Copilot-format server config entry (input format).
type CopilotEntry struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// ToOpenCodeFormat converts a Copilot-format entry to OpenCode format.
//
// Copilot: {"command": "npx", "args": ["-y", "pkg"], "env": {...}}
// OpenCode: {"type": "local", "command": ["npx", "-y", "pkg"], "environment": {...}, "enabled": true}
func ToOpenCodeFormat(entry CopilotEntry, enabled bool) ServerEntry {
	result := ServerEntry{
		Type:    "local",
		Enabled: enabled,
	}

	if entry.Command != "" {
		result.Command = append([]string{entry.Command}, entry.Args...)
	} else if entry.URL != "" {
		result.Type = "remote"
		result.URL = entry.URL
		if len(entry.Headers) > 0 {
			result.Headers = entry.Headers
		}
	}

	if len(entry.Env) > 0 {
		result.Environment = entry.Env
	}

	return result
}

// Adapter manages the OpenCode MCP configuration.
type Adapter struct {
	ProjectRoot string
}

// New creates a new OpenCode adapter for the given project root.
func New(projectRoot string) *Adapter {
	return &Adapter{ProjectRoot: projectRoot}
}

// ConfigPath returns the path to opencode.json in the project root.
func (a *Adapter) ConfigPath() string {
	return filepath.Join(a.ProjectRoot, "opencode.json")
}

// IsOptedIn returns true if the .opencode/ directory exists.
func (a *Adapter) IsOptedIn() bool {
	info, err := os.Stat(filepath.Join(a.ProjectRoot, ".opencode"))
	return err == nil && info.IsDir()
}

// GetCurrentConfig reads the current opencode.json contents.
func (a *Adapter) GetCurrentConfig() map[string]interface{} {
	data, err := os.ReadFile(a.ConfigPath())
	if err != nil {
		return map[string]interface{}{}
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}

// UpdateConfig merges configUpdates (Copilot-format) into the mcp section of opencode.json.
// Returns silently if .opencode/ does not exist.
func (a *Adapter) UpdateConfig(configUpdates map[string]CopilotEntry, enabled bool) error {
	if !a.IsOptedIn() {
		return nil
	}

	current := a.GetCurrentConfig()
	mcpSection, ok := current["mcp"].(map[string]interface{})
	if !ok {
		mcpSection = map[string]interface{}{}
	}

	for name, entry := range configUpdates {
		oc := ToOpenCodeFormat(entry, enabled)
		ocMap := map[string]interface{}{
			"type":    oc.Type,
			"enabled": oc.Enabled,
		}
		if len(oc.Command) > 0 {
			ocMap["command"] = oc.Command
		}
		if oc.URL != "" {
			ocMap["url"] = oc.URL
		}
		if len(oc.Headers) > 0 {
			ocMap["headers"] = oc.Headers
		}
		if len(oc.Environment) > 0 {
			ocMap["environment"] = oc.Environment
		}
		mcpSection[name] = ocMap
	}

	current["mcp"] = mcpSection

	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return fmt.Errorf("opencode: marshal config: %w", err)
	}

	return os.WriteFile(a.ConfigPath(), data, 0o644)
}
