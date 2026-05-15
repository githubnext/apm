// Package vscode implements the VS Code MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/vscode.py.
//
// VSCode uses .vscode/mcp.json at the project root with a "servers" key
// (plus an "inputs" section for ${input:VAR} variable definitions).
package vscode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/githubnext/apm/internal/adapters/client/copilot"
)

// inputVarRE matches ${input:NAME} placeholders.
var inputVarRE = regexp.MustCompile(`\$\{input:([^}]+)\}`)

// envVarRE matches ${VAR} and ${env:VAR}.
var envVarRE = regexp.MustCompile(`\$\{(?:env:)?([A-Za-z_][A-Za-z0-9_]*)\}`)

// legacyAngleVarRE matches <VARNAME> legacy placeholders.
var legacyAngleVarRE = regexp.MustCompile(`<([A-Z_][A-Z0-9_]*)>`)

// Adapter is the VS Code MCP client adapter.
type Adapter struct {
	*copilot.Adapter
}

// New creates a new VS Code adapter.
func New(projectRoot string, userScope bool) *Adapter {
	base := copilot.New(projectRoot, userScope)
	base.SupportsRuntimeEnvSubstitution = false
	return &Adapter{Adapter: base}
}

// TargetName returns "vscode".
func (a *Adapter) TargetName() string { return "vscode" }

// MCPServersKey returns "servers".
func (a *Adapter) MCPServersKey() string { return "servers" }

// SupportsUserScope returns false.
func (a *Adapter) SupportsUserScope() bool { return false }

// GetConfigPath returns the path to .vscode/mcp.json in the project root.
func (a *Adapter) GetConfigPath() string {
	root := a.ProjectRoot
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			root = "."
		}
	}
	return filepath.Join(root, ".vscode", "mcp.json")
}

// GetCurrentConfig reads the current .vscode/mcp.json.
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

// UpdateConfig writes the complete config to .vscode/mcp.json.
func (a *Adapter) UpdateConfig(configUpdates map[string]interface{}) error {
	configPath := a.GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(configUpdates, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0o644)
}

// InputVarDef is a VS Code input variable definition for the "inputs" array.
type InputVarDef struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Password    bool   `json:"password"`
}

// FormatServerConfig formats a server entry for VS Code's mcp.json schema.
// Returns the server config map and a list of input variable definitions.
func (a *Adapter) FormatServerConfig(
	serverInfo map[string]interface{},
	envOverrides map[string]interface{},
	runtimeVars map[string]string,
) (map[string]interface{}, []InputVarDef, error) {
	serverConfig := map[string]interface{}{}
	var inputVars []InputVarDef

	// Self-defined stdio deps.
	if raw, ok := serverInfo["_raw_stdio"].(map[string]interface{}); ok {
		serverConfig["type"] = "stdio"
		serverConfig["command"] = strField(raw, "command")
		serverConfig["args"] = raw["args"]
		if rawEnv, ok := raw["env"].(map[string]interface{}); ok && len(rawEnv) > 0 {
			translated := translateEnvVarsForVSCode(rawEnv)
			serverConfig["env"] = translated
			inputVars = append(inputVars, extractInputVariables(translated, strField(serverInfo, "name"))...)
		}
		return serverConfig, inputVars, nil
	}

	// Package-based servers.
	packages := toSliceOfMaps(serverInfo["packages"])
	if len(packages) > 0 {
		pkg := selectBestPackage(packages)
		if pkg == nil {
			return serverConfig, inputVars, nil
		}
		registryName := inferRegistryName(pkg)
		runtimeHint := strField(pkg, "runtime_hint")
		pkgArgs := extractPackageArgs(pkg)
		pkgName := strField(pkg, "name")

		switch {
		case runtimeHint == "npx" || registryName == "npm":
			extraArgs := filterOut(pkgArgs, pkgName)
			args := append([]interface{}{"-y", pkgName}, toInterfaceSlice(extraArgs)...)
			serverConfig = map[string]interface{}{
				"type":    "stdio",
				"command": "npx",
				"args":    args,
			}
		case runtimeHint == "docker" || registryName == "docker":
			args := pkgArgs
			if len(args) == 0 {
				args = []string{"run", "-i", "--rm", pkgName}
			}
			serverConfig = map[string]interface{}{
				"type":    "stdio",
				"command": "docker",
				"args":    toInterfaceSlice(args),
			}
		case registryName == "pypi" || runtimeHint == "uvx" || strings.Contains(runtimeHint, "python"):
			cmd := "uvx"
			if runtimeHint != "" && runtimeHint != "uvx" && runtimeHint != "pip" {
				cmd = runtimeHint
			}
			var args []string
			if len(pkgArgs) > 0 {
				args = pkgArgs
			} else {
				args = []string{pkgName}
			}
			serverConfig = map[string]interface{}{
				"type":    "stdio",
				"command": cmd,
				"args":    toInterfaceSlice(args),
			}
		case runtimeHint != "":
			args := pkgArgs
			if len(args) == 0 {
				args = []string{pkgName}
			}
			serverConfig = map[string]interface{}{
				"type":    "stdio",
				"command": runtimeHint,
				"args":    toInterfaceSlice(args),
			}
		}

		// Environment variables -> ${input:var-name} references.
		envVars := toSliceOfMaps(pkg["environment_variables"])
		if len(envVars) == 0 {
			envVars = toSliceOfMaps(pkg["environmentVariables"])
		}
		if len(envVars) > 0 {
			env := map[string]interface{}{}
			for _, ev := range envVars {
				name := strField(ev, "name")
				if name == "" {
					continue
				}
				inputVarName := strings.ReplaceAll(strings.ToLower(name), "_", "-")
				env[name] = "${input:" + inputVarName + "}"
				desc := strField(ev, "description")
				if desc == "" {
					desc = name + " for MCP server"
				}
				inputVars = append(inputVars, InputVarDef{
					Type:        "promptString",
					ID:          inputVarName,
					Description: desc,
					Password:    true,
				})
			}
			if len(env) > 0 {
				serverConfig["env"] = env
			}
		}

		return serverConfig, inputVars, nil
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
		}
		serverConfig = map[string]interface{}{
			"type": "sse",
			"url":  strings.TrimSpace(strField(remote, "url")),
		}
		if transport == "http" || transport == "streamable-http" {
			serverConfig["type"] = "http"
		}
		headers := toSliceOfMaps(remote["headers"])
		if len(headers) > 0 {
			hmap := map[string]interface{}{}
			for _, h := range headers {
				name := strField(h, "name")
				value := strField(h, "value")
				if name != "" {
					// Translate env-var placeholders to VS Code syntax.
					hmap[name] = translateEnvValueForVSCode(value)
				}
			}
			serverConfig["headers"] = hmap
		}
		return serverConfig, inputVars, nil
	}

	return serverConfig, inputVars, nil
}

// ConfigureMCPServer installs a single MCP server into .vscode/mcp.json.
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

	serverConfig, inputVars, err := a.FormatServerConfig(serverInfo, envOverrides, runtimeVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error formatting server config: %s\n", err)
		return false
	}
	if len(serverConfig) == 0 {
		fmt.Fprintf(os.Stderr, "[x] Unable to configure server: %s\n", serverURL)
		return false
	}

	configKey := serverURL
	if serverName != "" {
		configKey = serverName
	}

	current := a.GetCurrentConfig()
	if _, ok := current["servers"]; !ok {
		current["servers"] = map[string]interface{}{}
	}
	if _, ok := current["inputs"]; !ok {
		current["inputs"] = []interface{}{}
	}
	servers, _ := current["servers"].(map[string]interface{})
	servers[configKey] = serverConfig
	current["servers"] = servers

	// Merge input vars (avoid duplicates by ID).
	existingInputs, _ := current["inputs"].([]interface{})
	existingIDs := map[string]bool{}
	for _, inp := range existingInputs {
		if m, ok := inp.(map[string]interface{}); ok {
			existingIDs[strField(m, "id")] = true
		}
	}
	for _, iv := range inputVars {
		if !existingIDs[iv.ID] {
			existingInputs = append(existingInputs, map[string]interface{}{
				"type":        iv.Type,
				"id":          iv.ID,
				"description": iv.Description,
				"password":    iv.Password,
			})
			existingIDs[iv.ID] = true
		}
	}
	current["inputs"] = existingInputs

	if err := a.UpdateConfig(current); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error writing VS Code config: %s\n", err)
		return false
	}
	fmt.Printf("[+] Configured MCP server '%s' for VS Code\n", configKey)
	return true
}

// translateEnvVarsForVSCode converts env dict values from ${VAR}/${env:VAR}/<VAR>
// to VS Code's ${env:VAR} syntax. ${input:...} references are preserved.
func translateEnvVarsForVSCode(env map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(env))
	for k, v := range env {
		if s, ok := v.(string); ok {
			out[k] = translateEnvValueForVSCode(s)
		} else {
			out[k] = v
		}
	}
	return out
}

// translateEnvValueForVSCode converts a single value to VS Code env-var syntax.
func translateEnvValueForVSCode(s string) string {
	// Legacy <VAR> -> ${env:VAR}
	s = legacyAngleVarRE.ReplaceAllString(s, "${env:$1}")
	// ${VAR} -> ${env:VAR} (only when not already ${env:...} or ${input:...})
	s = envVarRE.ReplaceAllStringFunc(s, func(m string) string {
		if strings.HasPrefix(m, "${env:") || strings.HasPrefix(m, "${input:") {
			return m
		}
		sub := envVarRE.FindStringSubmatch(m)
		return "${env:" + sub[1] + "}"
	})
	return s
}

// extractInputVariables scans a translated env map for ${input:VAR} references
// and returns InputVarDef entries.
func extractInputVariables(env map[string]interface{}, serverName string) []InputVarDef {
	seen := map[string]bool{}
	var out []InputVarDef
	for _, v := range env {
		s, ok := v.(string)
		if !ok {
			continue
		}
		for _, m := range inputVarRE.FindAllStringSubmatch(s, -1) {
			id := m[1]
			if seen[id] {
				continue
			}
			seen[id] = true
			out = append(out, InputVarDef{
				Type:        "promptString",
				ID:          id,
				Description: id + " for " + serverName,
				Password:    true,
			})
		}
	}
	return out
}

// extractPackageArgs returns the combined runtime+package arguments for a package entry.
func extractPackageArgs(pkg map[string]interface{}) []string {
	rt := toStringSlice(pkg["runtime_arguments"])
	pk := toStringSlice(pkg["package_arguments"])
	return append(rt, pk...)
}

// filterOut removes occurrences of target from ss.
func filterOut(ss []string, target string) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != target {
			out = append(out, s)
		}
	}
	return out
}

// ---- helpers shared with other adapter packages ----

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
