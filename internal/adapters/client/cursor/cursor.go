// Package cursor implements the Cursor IDE MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/cursor.py.
//
// Cursor uses .cursor/mcp.json at the project root with "mcpServers" key.
// APM only writes when .cursor/ already exists (opt-in).
// Emits Cursor-native transport discriminators (type: stdio / type: http).
package cursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

// Adapter is the Cursor IDE MCP client adapter.
type Adapter struct {
	*copilot.Adapter
}

// New creates a new Cursor adapter.
func New(projectRoot string, userScope bool) *Adapter {
	base := copilot.New(projectRoot, userScope)
	base.SupportsRuntimeEnvSubstitution = false
	return &Adapter{Adapter: base}
}

// TargetName returns "cursor".
func (a *Adapter) TargetName() string { return "cursor" }

// MCPServersKey returns "mcpServers".
func (a *Adapter) MCPServersKey() string { return "mcpServers" }

// SupportsUserScope returns false.
func (a *Adapter) SupportsUserScope() bool { return false }

// GetConfigPath returns the path to .cursor/mcp.json in the project root.
func (a *Adapter) GetConfigPath() string {
	root := a.ProjectRoot
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			root = "."
		}
	}
	return filepath.Join(root, ".cursor", "mcp.json")
}

// GetCurrentConfig reads the current .cursor/mcp.json.
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

// UpdateConfig merges configUpdates only when .cursor/ already exists.
func (a *Adapter) UpdateConfig(configUpdates map[string]interface{}) error {
	root := a.ProjectRoot
	if root == "" {
		root, _ = os.Getwd()
	}
	cursorDir := filepath.Join(root, ".cursor")
	info, err := os.Stat(cursorDir)
	if err != nil || !info.IsDir() {
		// Opt-in: silently skip when .cursor/ doesn't exist.
		return nil
	}

	current := a.GetCurrentConfig()
	if _, ok := current["mcpServers"]; !ok {
		current["mcpServers"] = map[string]interface{}{}
	}
	servers, _ := current["mcpServers"].(map[string]interface{})
	for k, v := range configUpdates {
		servers[k] = v
	}
	current["mcpServers"] = servers

	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.GetConfigPath(), data, 0o644)
}

// FormatServerConfig formats a server entry in Cursor's native schema.
//
// Differences from Copilot:
//   - No "type":"local", no "tools", no "id" fields
//   - Stdio: emits explicit type:"stdio"
//   - HTTP:  emits type:"http"
func (a *Adapter) FormatServerConfig(
	serverInfo map[string]interface{},
	envOverrides map[string]interface{},
	runtimeVars map[string]string,
) (map[string]interface{}, error) {
	raw, err := a.Adapter.FormatServerConfig(serverInfo, envOverrides, runtimeVars)
	if err != nil {
		return nil, err
	}
	return normalizeMCPEntryForCursor(raw), nil
}

// normalizeMCPEntryForCursor strips Copilot-only fields and emits Cursor's wire format.
func normalizeMCPEntryForCursor(entry map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	delete(out, "tools")
	delete(out, "id")

	entryType, _ := out["type"].(string)
	if entryType == "local" {
		out["type"] = "stdio"
	} else if entryType == "http" || entryType == "remote" {
		out["type"] = "http"
	}
	return out
}

// ConfigureMCPServer installs a single MCP server into the Cursor config.
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
		fmt.Fprintf(os.Stderr, "[x] Error writing Cursor config: %s\n", err)
		return false
	}
	fmt.Printf("[+] Configured MCP server '%s' for Cursor\n", configKey)
	return true
}

func serverKeyFor(serverURL, serverName string) string {
	if serverName != "" {
		return serverName
	}
	if idx := strings.LastIndex(serverURL, "/"); idx >= 0 {
		return serverURL[idx+1:]
	}
	return serverURL
}
