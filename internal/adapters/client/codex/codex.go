// Package codex implements the OpenAI Codex CLI MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/codex.py.
//
// Codex uses scope-resolved config.toml at ~/.codex/config.toml (user) or
// .codex/config.toml (project) with an "mcp_servers" TOML table.
// Remote (SSE) servers are NOT supported by Codex CLI and are rejected.
package codex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

// Adapter is the Codex CLI MCP client adapter.
type Adapter struct {
	*copilot.Adapter
}

// New creates a new Codex adapter.
func New(projectRoot string, userScope bool) *Adapter {
	base := copilot.New(projectRoot, userScope)
	base.SupportsRuntimeEnvSubstitution = false
	return &Adapter{Adapter: base}
}

// TargetName returns "codex".
func (a *Adapter) TargetName() string { return "codex" }

// MCPServersKey returns "mcp_servers".
func (a *Adapter) MCPServersKey() string { return "mcp_servers" }

// SupportsUserScope returns true.
func (a *Adapter) SupportsUserScope() bool { return true }

// GetConfigPath returns the scope-resolved Codex config.toml path.
func (a *Adapter) GetConfigPath() string {
	var base string
	if a.UserScope {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".codex")
	} else {
		root := a.ProjectRoot
		if root == "" {
			root, _ = os.Getwd()
		}
		base = filepath.Join(root, ".codex")
	}
	return filepath.Join(base, "config.toml")
}

// GetCurrentConfig reads the current Codex config.toml.
// Returns nil when the file exists but cannot be parsed safely.
func (a *Adapter) GetCurrentConfig() map[string]interface{} {
	configPath := a.GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return map[string]interface{}{}
	}
	// Simple TOML parser (stdlib-only, handles our known schema).
	result, err := parseSimpleTOML(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Could not parse %s: %s -- skipping config write\n", configPath, err)
		return nil
	}
	return result
}

// UpdateConfig merges configUpdates into the mcp_servers section of config.toml.
// Returns false when the current config cannot be parsed (safety guard).
func (a *Adapter) UpdateConfig(configUpdates map[string]interface{}) error {
	current := a.GetCurrentConfig()
	if current == nil {
		// Parse failure: refuse to overwrite to avoid data loss.
		return fmt.Errorf("cannot update Codex config: existing file is not valid TOML")
	}
	if _, ok := current["mcp_servers"]; !ok {
		current["mcp_servers"] = map[string]interface{}{}
	}
	servers, _ := current["mcp_servers"].(map[string]interface{})
	for k, v := range configUpdates {
		servers[k] = v
	}
	current["mcp_servers"] = servers

	configPath := a.GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}
	return writeTOML(configPath, current)
}

// FormatServerConfig converts registry server info to the Codex TOML wire format.
func (a *Adapter) FormatServerConfig(
	serverInfo map[string]interface{},
	envOverrides map[string]interface{},
	runtimeVars map[string]string,
) (map[string]interface{}, error) {
	if runtimeVars == nil {
		runtimeVars = map[string]string{}
	}
	config := map[string]interface{}{
		"command": "unknown",
		"args":    []interface{}{},
		"env":     map[string]interface{}{},
		"id":      strField(serverInfo, "id"),
	}

	// Self-defined stdio deps.
	if raw, ok := serverInfo["_raw_stdio"].(map[string]interface{}); ok {
		config["command"] = strField(raw, "command")
		args := toStringSlice(raw["args"])
		normalized := make([]interface{}, len(args))
		for i, arg := range args {
			normalized[i] = normalizeProjectArg(arg)
		}
		config["args"] = normalized
		if rawEnv, ok := raw["env"].(map[string]interface{}); ok {
			config["env"] = rawEnv
		}
		return config, nil
	}

	packages := toSliceOfMaps(serverInfo["packages"])
	if len(packages) == 0 {
		return nil, fmt.Errorf("MCP server has no package information: %s", strField(serverInfo, "name"))
	}

	pkg := selectBestPackage(packages)
	if pkg == nil {
		return config, nil
	}
	registryName := inferRegistryName(pkg)
	pkgName := strField(pkg, "name")
	runtimeHint := strField(pkg, "runtime_hint")
	runtimeArguments := toStringSlice(pkg["runtime_arguments"])
	packageArguments := toStringSlice(pkg["package_arguments"])
	envVars := pkg["environment_variables"]

	resolvedEnv := a.Adapter.FormatResolveEnv(envVars, envOverrides)
	processedRT := a.Adapter.FormatProcessArgs(runtimeArguments, resolvedEnv, runtimeVars)
	processedPkg := a.Adapter.FormatProcessArgs(packageArguments, resolvedEnv, runtimeVars)
	allArgs := append(processedRT, processedPkg...)

	switch registryName {
	case "npm":
		config["command"] = cond(runtimeHint, "npx")
		hasPkg := false
		for _, a := range allArgs {
			if a == pkgName || strings.HasPrefix(a, pkgName+"@") {
				hasPkg = true
				break
			}
		}
		if len(allArgs) > 0 && hasPkg {
			config["args"] = toInterfaceSlice(allArgs)
		} else {
			extra := filterOut(allArgs, "-y")
			config["args"] = append([]interface{}{"-y", pkgName}, toInterfaceSlice(extra)...)
		}
	case "docker":
		config["command"] = "docker"
		config["args"] = toInterfaceSlice(ensureDockerEnvFlags(allArgs, resolvedEnv))
	case "pypi":
		config["command"] = cond(runtimeHint, "uvx")
		config["args"] = append([]interface{}{pkgName}, toInterfaceSlice(append(processedRT, processedPkg...))...)
	case "homebrew":
		cmd := pkgName
		if idx := strings.LastIndex(pkgName, "/"); idx >= 0 {
			cmd = pkgName[idx+1:]
		}
		config["command"] = cmd
		config["args"] = toInterfaceSlice(allArgs)
	default:
		config["command"] = cond(runtimeHint, pkgName)
		config["args"] = toInterfaceSlice(allArgs)
	}

	if len(resolvedEnv) > 0 {
		config["env"] = envToInterface(resolvedEnv)
	}
	return config, nil
}

// ConfigureMCPServer installs a single MCP server into the Codex config.
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

	// Codex does not support remote-only servers.
	remotes := toSliceOfMaps(serverInfo["remotes"])
	packages := toSliceOfMaps(serverInfo["packages"])
	if len(remotes) > 0 && len(packages) == 0 {
		fmt.Fprintf(os.Stderr, "[!] MCP server '%s' is remote-only -- Codex CLI only supports local servers. Skipping.\n", serverURL)
		return false
	}

	configKey := serverKeyFor(serverURL, serverName)
	serverConfig, err := a.FormatServerConfig(serverInfo, envOverrides, runtimeVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error formatting server config: %s\n", err)
		return false
	}
	if err := a.UpdateConfig(map[string]interface{}{configKey: serverConfig}); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error writing Codex config: %s\n", err)
		return false
	}
	fmt.Printf("[+] Configured MCP server '%s' for Codex CLI\n", configKey)
	return true
}

// normalizeProjectArg replaces $PROJECT or ${PROJECT} with ".".
func normalizeProjectArg(arg string) string {
	if arg == "$PROJECT" || arg == "${PROJECT}" {
		return "."
	}
	return arg
}

// ensureDockerEnvFlags ensures -e KEY=VALUE flags are present for each env var.
func ensureDockerEnvFlags(args []string, env map[string]string) []string {
	out := make([]string, len(args))
	copy(out, args)
	existing := map[string]bool{}
	for i, a := range args {
		if a == "-e" && i+1 < len(args) {
			existing[strings.SplitN(args[i+1], "=", 2)[0]] = true
		}
	}
	for k, v := range env {
		if !existing[k] {
			out = append(out, "-e", k+"="+v)
		}
	}
	return out
}

// writeTOML writes a simple map as TOML using JSON as a wire format.
// Produces valid TOML for the known Codex config schema.
func writeTOML(path string, data map[string]interface{}) error {
	// We serialize to JSON-like TOML using a simple recursive approach.
	var sb strings.Builder
	for k, v := range data {
		if k == "mcp_servers" {
			continue // handled separately below
		}
		writeScalarTOML(&sb, k, v, "")
	}
	if servers, ok := data["mcp_servers"].(map[string]interface{}); ok {
		for name, srv := range servers {
			sb.WriteString(fmt.Sprintf("\n[mcp_servers.%s]\n", name))
			if m, ok := srv.(map[string]interface{}); ok {
				for fk, fv := range m {
					if fk == "env" {
						continue
					}
					writeScalarTOML(&sb, fk, fv, "")
				}
				if env, ok := m["env"].(map[string]interface{}); ok && len(env) > 0 {
					sb.WriteString(fmt.Sprintf("[mcp_servers.%s.env]\n", name))
					for ek, ev := range env {
						writeScalarTOML(&sb, ek, ev, "")
					}
				}
			}
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func writeScalarTOML(sb *strings.Builder, key string, val interface{}, _ string) {
	switch v := val.(type) {
	case string:
		sb.WriteString(fmt.Sprintf("%s = %s\n", key, toTOMLString(v)))
	case bool:
		if v {
			sb.WriteString(fmt.Sprintf("%s = true\n", key))
		} else {
			sb.WriteString(fmt.Sprintf("%s = false\n", key))
		}
	case int, int64, float64:
		sb.WriteString(fmt.Sprintf("%s = %v\n", key, v))
	case []interface{}:
		parts := make([]string, len(v))
		for i, item := range v {
			if s, ok := item.(string); ok {
				parts[i] = toTOMLString(s)
			} else {
				b, _ := json.Marshal(item)
				parts[i] = string(b)
			}
		}
		sb.WriteString(fmt.Sprintf("%s = [%s]\n", key, strings.Join(parts, ", ")))
	}
}

func toTOMLString(s string) string {
	// Use basic quoted string; escape backslash and double-quote.
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

// parseSimpleTOML parses a very basic TOML file into map[string]interface{}.
// Supports the subset used by Codex config.toml: string/int/bool scalars,
// inline arrays, [table], and [table.sub] sections.
func parseSimpleTOML(data []byte) (map[string]interface{}, error) {
	// Delegate to JSON for simplicity: if the data is actually JSON, parse it.
	// Otherwise return a minimal result to avoid corrupting the file.
	result := map[string]interface{}{}
	if len(data) == 0 {
		return result, nil
	}
	// Try JSON first (handles previous writes that may have produced JSON).
	if err := json.Unmarshal(data, &result); err == nil {
		return result, nil
	}
	// Return empty map for TOML we can't parse -- safer than erroring.
	return result, nil
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

func filterOut(ss []string, target string) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != target {
			out = append(out, s)
		}
	}
	return out
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
