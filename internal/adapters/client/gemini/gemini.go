// Package gemini implements the Gemini CLI MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/gemini.py.
//
// Gemini CLI uses .gemini/settings.json at the project root with an "mcpServers" key.
// Transport is inferred from key presence (command=stdio, url=SSE, httpUrl=HTTP).
// APM only writes when .gemini/ already exists (opt-in).
package gemini

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

// Adapter is the Gemini CLI MCP client adapter.
type Adapter struct {
	*copilot.Adapter
}

// New creates a new Gemini adapter.
func New(projectRoot string, userScope bool) *Adapter {
	base := copilot.New(projectRoot, userScope)
	base.SupportsRuntimeEnvSubstitution = false
	return &Adapter{Adapter: base}
}

// TargetName returns "gemini".
func (a *Adapter) TargetName() string { return "gemini" }

// MCPServersKey returns "mcpServers".
func (a *Adapter) MCPServersKey() string { return "mcpServers" }

// SupportsUserScope returns true (Gemini has a global settings path).
func (a *Adapter) SupportsUserScope() bool { return true }

// GetConfigPath returns the path to .gemini/settings.json in the project root.
func (a *Adapter) GetConfigPath() string {
	root := a.ProjectRoot
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			root = "."
		}
	}
	return filepath.Join(root, ".gemini", "settings.json")
}

// GetCurrentConfig reads the current .gemini/settings.json.
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

// UpdateConfig merges configUpdates into mcpServers only when .gemini/ exists.
func (a *Adapter) UpdateConfig(configUpdates map[string]interface{}) error {
	root := a.ProjectRoot
	if root == "" {
		root, _ = os.Getwd()
	}
	geminiDir := filepath.Join(root, ".gemini")
	info, err := os.Stat(geminiDir)
	if err != nil || !info.IsDir() {
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

	configPath := a.GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0o644)
}

// FormatServerConfig formats a server entry for Gemini CLI's schema.
//
// Gemini's schema differs from Copilot:
//   - No type, tools, or id fields
//   - Transport inferred from key: command (stdio), url (SSE), httpUrl (streamable HTTP)
//   - Tool filtering via includeTools/excludeTools
func (a *Adapter) FormatServerConfig(
	serverInfo map[string]interface{},
	envOverrides map[string]interface{},
	runtimeVars map[string]string,
) (map[string]interface{}, error) {
	if runtimeVars == nil {
		runtimeVars = map[string]string{}
	}
	config := map[string]interface{}{}

	// Self-defined stdio deps.
	if raw, ok := serverInfo["_raw_stdio"].(map[string]interface{}); ok {
		config["command"] = strField(raw, "command")
		args := toStringSlice(raw["args"])
		if len(args) > 0 {
			config["args"] = toInterfaceSlice(args)
		}
		if rawEnv, ok := raw["env"].(map[string]interface{}); ok && len(rawEnv) > 0 {
			config["env"] = rawEnv
		}
		return config, nil
	}

	// Remote endpoints.
	remotes := toSliceOfMaps(serverInfo["remotes"])
	if len(remotes) > 0 {
		remote := selectRemoteWithURL(remotes)
		if remote == nil {
			remote = remotes[0]
		}
		transport := strings.TrimSpace(strField(remote, "transport_type"))
		if transport == "" {
			transport = "http"
		} else if transport != "sse" && transport != "http" && transport != "streamable-http" {
			return nil, fmt.Errorf("unsupported remote transport %q for Gemini (server %s)", transport, strField(serverInfo, "name"))
		}
		url := strings.TrimSpace(strField(remote, "url"))
		if transport == "sse" {
			config["url"] = url
		} else {
			config["httpUrl"] = url
		}
		headers := toSliceOfMaps(remote["headers"])
		for _, header := range headers {
			name := strField(header, "name")
			value := strField(header, "value")
			if name != "" && value != "" {
				if _, ok := config["headers"]; !ok {
					config["headers"] = map[string]interface{}{}
				}
				config["headers"].(map[string]interface{})[name] = value
			}
		}
		return config, nil
	}

	// Local packages.
	packages := toSliceOfMaps(serverInfo["packages"])
	if len(packages) == 0 {
		return nil, fmt.Errorf("MCP server has no package information or remote endpoints: %s", strField(serverInfo, "name"))
	}

	pkg := selectBestPackage(packages)
	if pkg == nil {
		return config, nil
	}
	registryName := inferRegistryName(pkg)
	packageName := strField(pkg, "name")
	runtimeHint := strField(pkg, "runtime_hint")
	runtimeArguments := toStringSlice(pkg["runtime_arguments"])
	packageArguments := toStringSlice(pkg["package_arguments"])
	envVars := pkg["environment_variables"]

	resolvedEnv := a.Adapter.FormatResolveEnv(envVars, envOverrides)
	processedRT := a.Adapter.FormatProcessArgs(runtimeArguments, resolvedEnv, runtimeVars)
	processedPkg := a.Adapter.FormatProcessArgs(packageArguments, resolvedEnv, runtimeVars)

	switch registryName {
	case "npm":
		config["command"] = cond(runtimeHint, "npx")
		args := append([]interface{}{"-y", packageName}, toInterfaceSlice(processedRT)...)
		config["args"] = append(args, toInterfaceSlice(processedPkg)...)
	case "docker":
		config["command"] = "docker"
		if len(processedRT) > 0 {
			config["args"] = toInterfaceSlice(processedRT)
		} else {
			config["args"] = toInterfaceSlice([]string{"run", "-i", "--rm", packageName})
		}
	case "pypi":
		config["command"] = cond(runtimeHint, "uvx")
		args := append([]interface{}{packageName}, toInterfaceSlice(processedRT)...)
		config["args"] = append(args, toInterfaceSlice(processedPkg)...)
	case "homebrew":
		cmd := packageName
		if idx := strings.LastIndex(packageName, "/"); idx >= 0 {
			cmd = packageName[idx+1:]
		}
		config["command"] = cmd
		config["args"] = append(toInterfaceSlice(processedRT), toInterfaceSlice(processedPkg)...)
	default:
		config["command"] = cond(runtimeHint, packageName)
		config["args"] = append(toInterfaceSlice(processedRT), toInterfaceSlice(processedPkg)...)
	}

	if len(resolvedEnv) > 0 {
		config["env"] = envToInterface(resolvedEnv)
	}
	return config, nil
}

// ConfigureMCPServer installs a single MCP server into the Gemini config.
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
	root := a.ProjectRoot
	if root == "" {
		root, _ = os.Getwd()
	}
	geminiDir := filepath.Join(root, ".gemini")
	if info, err := os.Stat(geminiDir); err != nil || !info.IsDir() {
		return true // opt-in: silently succeed when .gemini/ absent
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
		fmt.Fprintf(os.Stderr, "[x] Error writing Gemini config: %s\n", err)
		return false
	}
	fmt.Printf("[+] Configured MCP server '%s' for Gemini CLI\n", configKey)
	return true
}

// ---- helpers ----

func strField(m map[string]interface{}, key string) string {
	v, _ := m[key].(string)
	return v
}

func toStringSlice(v interface{}) []string {
	switch s := v.(type) {
	case []string:
		return s
	case []interface{}:
		out := make([]string, 0, len(s))
		for _, item := range s {
			out = append(out, fmt.Sprintf("%v", item))
		}
		return out
	}
	return nil
}

func toSliceOfMaps(v interface{}) []map[string]interface{} {
	sl, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(sl))
	for _, item := range sl {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func toInterfaceSlice(ss []string) []interface{} {
	out := make([]interface{}, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func envToInterface(m map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func selectRemoteWithURL(remotes []map[string]interface{}) map[string]interface{} {
	for _, r := range remotes {
		if strings.TrimSpace(strField(r, "url")) != "" {
			return r
		}
	}
	return nil
}

func selectBestPackage(packages []map[string]interface{}) map[string]interface{} {
	priority := map[string]int{"npm": 0, "docker": 1, "pypi": 2, "homebrew": 3}
	best := packages[0]
	bestScore := 9999
	for _, p := range packages {
		score, ok := priority[inferRegistryName(p)]
		if !ok {
			score = 4
		}
		if score < bestScore {
			bestScore = score
			best = p
		}
	}
	return best
}

func inferRegistryName(pkg map[string]interface{}) string {
	if r := strField(pkg, "registry"); r != "" {
		lower := strings.ToLower(r)
		switch {
		case strings.Contains(lower, "npm"):
			return "npm"
		case strings.Contains(lower, "docker"):
			return "docker"
		case strings.Contains(lower, "pypi"):
			return "pypi"
		case strings.Contains(lower, "homebrew"):
			return "homebrew"
		}
		return lower
	}
	return "npm"
}

func cond(preferred, fallback string) string {
	if preferred != "" {
		return preferred
	}
	return fallback
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
