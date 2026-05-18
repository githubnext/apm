// Package claude implements the Claude Code MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/claude.py.
//
// Claude Code uses .mcp.json at the project root (project scope) or
// ~/.claude.json (user scope) with top-level "mcpServers" key.
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

// Adapter is the Claude Code MCP client adapter.
//
// It inherits all helper methods from the Copilot adapter but:
//   - Does NOT support runtime env substitution (installs literal values)
//   - Normalises entries for Claude Code's on-disk shape (strips Copilot-only fields)
//   - Supports both project scope (.mcp.json) and user scope (~/.claude.json)
type Adapter struct {
	*copilot.Adapter
}

// New creates a new Claude adapter.
func New(projectRoot string, userScope bool) *Adapter {
	base := copilot.New(projectRoot, userScope)
	base.SupportsRuntimeEnvSubstitution = false
	return &Adapter{Adapter: base}
}

// TargetName returns "claude".
func (a *Adapter) TargetName() string { return "claude" }

// MCPServersKey returns "mcpServers".
func (a *Adapter) MCPServersKey() string { return "mcpServers" }

// SupportsUserScope returns true.
func (a *Adapter) SupportsUserScope() bool { return true }

// GetConfigPath returns the scope-resolved config file path.
//
// Project scope: <project_root>/.mcp.json
// User scope:    ~/.claude.json
func (a *Adapter) GetConfigPath() string {
	if a.UserScope {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".claude.json")
	}
	root := a.ProjectRoot
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			root = "."
		}
	}
	return filepath.Join(root, ".mcp.json")
}

// GetCurrentConfig reads the current config from the appropriate file.
func (a *Adapter) GetCurrentConfig() map[string]interface{} {
	data, err := os.ReadFile(a.GetConfigPath())
	if err != nil {
		return map[string]interface{}{}
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return map[string]interface{}{}
	}
	return cfg
}

// UpdateConfig merges configUpdates into the mcpServers section.
//
// For user scope, creates the file with 0o600 permissions on first write.
func (a *Adapter) UpdateConfig(configUpdates map[string]interface{}) error {
	configPath := a.GetConfigPath()
	current := a.GetCurrentConfig()

	if _, ok := current["mcpServers"]; !ok {
		current["mcpServers"] = map[string]interface{}{}
	}
	servers, _ := current["mcpServers"].(map[string]interface{})
	for k, v := range configUpdates {
		servers[k] = v
	}
	current["mcpServers"] = servers

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}

	perm := os.FileMode(0o644)
	if a.UserScope {
		perm = 0o600
	}
	return os.WriteFile(configPath, data, perm)
}

// FormatServerConfig wraps the Copilot formatter then normalises the entry
// for Claude Code's on-disk shape.
func (a *Adapter) FormatServerConfig(
	serverInfo map[string]interface{},
	envOverrides map[string]interface{},
	runtimeVars map[string]string,
) (map[string]interface{}, error) {
	raw, err := a.Adapter.FormatServerConfig(serverInfo, envOverrides, runtimeVars)
	if err != nil {
		return nil, err
	}
	return normalizeMCPEntryForClaudeCode(raw), nil
}

// normalizeMCPEntryForClaudeCode strips Copilot-only fields and emits
// the Claude Code on-disk shape.
//
// For remote servers: keeps type/url/headers per Claude Code docs.
// For stdio servers:  drops type:"local", tools, and empty id; emits type:"stdio".
func normalizeMCPEntryForClaudeCode(entry map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	entryType, _ := out["type"].(string)

	if entryType == "http" || entryType == "remote" {
		// Remote: keep as-is, delete Copilot-only fields.
		delete(out, "tools")
		delete(out, "id")
		return out
	}

	// stdio: normalise.
	delete(out, "tools")
	if id, _ := out["id"].(string); id == "" {
		delete(out, "id")
	}
	delete(out, "type")

	// Only emit type:stdio when command is present.
	if cmd, _ := out["command"].(string); cmd != "" {
		out["type"] = "stdio"
	}

	return out
}

// ConfigureMCPServer installs a single MCP server into the Claude config.
func (a *Adapter) ConfigureMCPServer(
	serverURL, serverName string,
	enabled bool,
	envOverrides map[string]interface{},
	serverInfoCache map[string]interface{},
	runtimeVars map[string]string,
) bool {
	if serverURL == "" {
		fmt.Fprintln(os.Stderr, "[x] server_url cannot be empty")
		return false
	}

	var serverInfo map[string]interface{}
	if serverInfoCache != nil {
		if v, ok := serverInfoCache[serverURL]; ok {
			serverInfo, _ = v.(map[string]interface{})
		}
	}
	if serverInfo == nil {
		fmt.Fprintf(os.Stderr, "[x] MCP server '%s' not found in registry\n", serverURL)
		return false
	}

	serverConfig, err := a.FormatServerConfig(serverInfo, envOverrides, runtimeVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error formatting server config: %s\n", err)
		return false
	}

	configKey := serverKeyFor(serverURL, serverName)
	if err := a.UpdateConfig(map[string]interface{}{configKey: serverConfig}); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error writing Claude config: %s\n", err)
		return false
	}

	fmt.Printf("[+] Configured MCP server '%s' for Claude Code\n", configKey)
	return true
}

// ---- helpers ----

func serverKeyFor(serverURL, serverName string) string {
	if serverName != "" {
		return serverName
	}
	if idx := strings.LastIndex(serverURL, "/"); idx >= 0 {
		return serverURL[idx+1:]
	}
	return serverURL
}
