// Package copilot implements the GitHub Copilot CLI MCP client adapter.
//
// Mirrors src/apm_cli/adapters/client/copilot.py.
//
// The adapter writes MCP server configuration to ~/.copilot/mcp-config.json.
// Unlike legacy adapters, it emits runtime-substitution placeholders (${VAR})
// rather than resolving secrets at install time (see issue #1152).
package copilot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Legacy angle-bracket placeholder pattern: <VARNAME>
var legacyAngleVarRE = regexp.MustCompile(`<([A-Z_][A-Z0-9_]*)>`)

// Combined env-var placeholder regex covering all three syntaxes Copilot accepts.
var copilotEnvRE = regexp.MustCompile(`<([A-Z_][A-Z0-9_]*)>|\$\{(?:env:)?([A-Za-z_][A-Za-z0-9_]*)\}`)

// envVarRE matches ${VAR} and ${env:VAR}.
var envVarRE = regexp.MustCompile(`\$\{(?:env:)?([A-Za-z_][A-Za-z0-9_]*)\}`)

// defaultGitHubEnv holds non-secret literal defaults that stay literal in translate mode.
var defaultGitHubEnv = map[string]string{
	"GITHUB_TOOLSETS":          "context",
	"GITHUB_DYNAMIC_TOOLSETS":  "1",
}

// process-wide aggregation state (mirrors class-level Python ClassVar fields).
var (
	globalMu                  sync.Mutex
	legacyAngleOffenders      = map[string][]string{}
	securityUpgradedKeys      = map[string]bool{}
	unsetEnvKeysByServer      = map[string][]string{}
	installRunSummaryEmitted  bool
)

// TranslateEnvPlaceholder converts env-var placeholders to ${VAR} form.
//
// Translations:
//
//	${env:VAR}  -> ${VAR}
//	${VAR}      -> ${VAR} (no-op)
//	<VAR>       -> ${VAR} (legacy migration)
//	non-string  -> passthrough
func TranslateEnvPlaceholder(value string) string {
	return copilotEnvRE.ReplaceAllStringFunc(value, func(m string) string {
		sub := copilotEnvRE.FindStringSubmatch(m)
		if sub[1] != "" {
			return "${" + sub[1] + "}"
		}
		return "${" + sub[2] + "}"
	})
}

// ExtractLegacyAngleVars returns the set of <VAR> names in value.
func ExtractLegacyAngleVars(value string) []string {
	matches := legacyAngleVarRE.FindAllStringSubmatch(value, -1)
	seen := map[string]bool{}
	out := []string{}
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

// HasEnvPlaceholder returns true if value contains any recognised env-var
// placeholder syntax.
func HasEnvPlaceholder(value string) bool {
	return copilotEnvRE.MatchString(value)
}

// Adapter is the Copilot CLI MCP client adapter.
//
// It targets ~/.copilot/mcp-config.json and emits ${VAR} runtime-substitution
// placeholders (SupportsRuntimeEnvSubstitution = true).
type Adapter struct {
	ProjectRoot                    string
	UserScope                      bool
	SupportsRuntimeEnvSubstitution bool

	// per-server tracking populated during FormatServerConfig
	lastEnvPlaceholderKeys []string
	lastLegacyAngleVars    []string
}

// New creates a new Copilot adapter.
func New(projectRoot string, userScope bool) *Adapter {
	return &Adapter{
		ProjectRoot:                    projectRoot,
		UserScope:                      userScope,
		SupportsRuntimeEnvSubstitution: true,
	}
}

// TargetName returns "copilot".
func (a *Adapter) TargetName() string { return "copilot" }

// MCPServersKey returns "mcpServers".
func (a *Adapter) MCPServersKey() string { return "mcpServers" }

// SupportsUserScope returns true.
func (a *Adapter) SupportsUserScope() bool { return true }

// GetConfigPath returns the path to ~/.copilot/mcp-config.json.
func (a *Adapter) GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".copilot", "mcp-config.json")
}

// GetCurrentConfig reads and returns the current Copilot config, or {} on error.
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

// UpdateConfig merges configUpdates into the mcpServers section of mcp-config.json.
func (a *Adapter) UpdateConfig(configUpdates map[string]interface{}) error {
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

// ConfigureMCPServer installs a single MCP server into the Copilot config.
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

	a.lastEnvPlaceholderKeys = nil
	a.lastLegacyAngleVars = nil

	// Snapshot previously baked keys for security-upgrade detection.
	prevBakedKeys, prevBakedHeaders := a.collectPreviouslyBakedKeys(serverURL, serverName)

	serverConfig, err := a.FormatServerConfig(serverInfo, envOverrides, runtimeVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error configuring MCP server: %s\n", err)
		return false
	}

	configKey := serverKeyFor(serverURL, serverName)

	if err := a.UpdateConfig(map[string]interface{}{configKey: serverConfig}); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Error writing config: %s\n", err)
		return false
	}

	// Aggregate diagnostics.
	if a.SupportsRuntimeEnvSubstitution {
		if len(a.lastLegacyAngleVars) > 0 {
			globalMu.Lock()
			legacyAngleOffenders[configKey] = a.lastLegacyAngleVars
			globalMu.Unlock()
		}
		upgradedKeys := intersect(prevBakedKeys, a.lastEnvPlaceholderKeys)
		if prevBakedHeaders && len(a.lastEnvPlaceholderKeys) > 0 {
			upgradedKeys = union(upgradedKeys, a.lastEnvPlaceholderKeys)
		}
		if len(upgradedKeys) > 0 {
			globalMu.Lock()
			for _, k := range upgradedKeys {
				securityUpgradedKeys[k] = true
			}
			globalMu.Unlock()
		}
	}

	a.emitInstallSummary(configKey, serverConfig)
	return true
}

// collectPreviouslyBakedKeys returns the env keys and headers baked status
// for the current on-disk config of the given server.
func (a *Adapter) collectPreviouslyBakedKeys(serverURL, serverName string) ([]string, bool) {
	current := a.GetCurrentConfig()
	servers, _ := current["mcpServers"].(map[string]interface{})
	key := serverKeyFor(serverURL, serverName)
	existing, _ := servers[key].(map[string]interface{})
	if existing == nil {
		return nil, false
	}
	var bakedEnvKeys []string
	if envBlock, ok := existing["env"].(map[string]interface{}); ok {
		for k, v := range envBlock {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" && !HasEnvPlaceholder(s) {
				bakedEnvKeys = append(bakedEnvKeys, k)
			}
		}
	}
	headersBaked := false
	if hBlock, ok := existing["headers"].(map[string]interface{}); ok {
		for _, v := range hBlock {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" && !HasEnvPlaceholder(s) {
				headersBaked = true
				break
			}
		}
	}
	return bakedEnvKeys, headersBaked
}

// emitInstallSummary records unset env vars for the post-install summary.
func (a *Adapter) emitInstallSummary(configKey string, serverConfig map[string]interface{}) {
	if !a.SupportsRuntimeEnvSubstitution {
		return
	}
	keys := map[string]bool{}
	for _, k := range a.lastEnvPlaceholderKeys {
		keys[k] = true
	}
	for _, blockKey := range []string{"env", "headers"} {
		if block, ok := serverConfig[blockKey].(map[string]interface{}); ok {
			for _, v := range block {
				if s, ok := v.(string); ok {
					for _, m := range envVarRE.FindAllStringSubmatch(s, -1) {
						keys[m[1]] = true
					}
				}
			}
		}
	}
	var unset []string
	for name := range keys {
		if os.Getenv(name) == "" {
			unset = append(unset, name)
		}
	}
	if len(unset) > 0 {
		globalMu.Lock()
		existing := unsetEnvKeysByServer[configKey]
		seen := map[string]bool{}
		for _, u := range existing {
			seen[u] = true
		}
		for _, u := range unset {
			if !seen[u] {
				existing = append(existing, u)
			}
		}
		unsetEnvKeysByServer[configKey] = existing
		globalMu.Unlock()
	}
}

// ResetInstallRunState resets process-wide aggregation buckets (for tests).
func ResetInstallRunState() {
	globalMu.Lock()
	defer globalMu.Unlock()
	legacyAngleOffenders = map[string][]string{}
	securityUpgradedKeys = map[string]bool{}
	unsetEnvKeysByServer = map[string][]string{}
	installRunSummaryEmitted = false
}

// FormatServerConfig converts registry server info to Copilot CLI's wire format.
func (a *Adapter) FormatServerConfig(
	serverInfo map[string]interface{},
	envOverrides map[string]interface{},
	runtimeVars map[string]string,
) (map[string]interface{}, error) {
	if runtimeVars == nil {
		runtimeVars = map[string]string{}
	}

	config := map[string]interface{}{
		"type":  "local",
		"tools": []interface{}{"*"},
		"id":    strField(serverInfo, "id"),
	}

	// Self-defined stdio deps carry raw command/args.
	if raw, ok := serverInfo["_raw_stdio"].(map[string]interface{}); ok {
		config["command"] = strField(raw, "command")
		resolvedEnv := map[string]string{}
		if rawEnv, ok := raw["env"].(map[string]interface{}); ok {
			resolvedEnv = a.resolveEnvVarsDict(rawEnv, envOverrides)
			config["env"] = envToInterface(resolvedEnv)
		}
		args := toStringSlice(raw["args"])
		resolved := make([]interface{}, len(args))
		for i, arg := range args {
			resolved[i] = a.resolveVariablePlaceholders(arg, resolvedEnv, runtimeVars)
		}
		config["args"] = resolved
		if toolsOverride := serverInfo["_apm_tools_override"]; toolsOverride != nil {
			config["tools"] = toolsOverride
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
			return nil, fmt.Errorf("unsupported remote transport %q for Copilot (server %s)", transport, strField(serverInfo, "name"))
		}
		remoteConfig := map[string]interface{}{
			"type":  "http",
			"url":   strings.TrimSpace(strField(remote, "url")),
			"tools": []interface{}{"*"},
			"id":    strField(serverInfo, "id"),
		}
		serverName := strField(serverInfo, "name")
		if a.isGitHubServer(serverName, strField(remote, "url")) {
			if token := a.getGitHubToken(); token != "" {
				remoteConfig["headers"] = map[string]interface{}{
					"Authorization": "Bearer " + token,
				}
			}
		}
		headers := toSliceOfMaps(remote["headers"])
		for _, header := range headers {
			name := strField(header, "name")
			value := strField(header, "value")
			if name != "" && value != "" {
				resolved := a.resolveEnvVariable(name, value, envOverrides)
				if _, ok := remoteConfig["headers"]; !ok {
					remoteConfig["headers"] = map[string]interface{}{}
				}
				remoteConfig["headers"].(map[string]interface{})[name] = resolved
			}
		}
		if toolsOverride := serverInfo["_apm_tools_override"]; toolsOverride != nil {
			remoteConfig["tools"] = toolsOverride
		}
		return remoteConfig, nil
	}

	// Local packages.
	packages := toSliceOfMaps(serverInfo["packages"])
	if len(packages) == 0 {
		return nil, fmt.Errorf("MCP server has incomplete configuration (no packages or remotes): %s", strField(serverInfo, "name"))
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

	resolvedEnv := a.resolveEnvironmentVariables(envVars, envOverrides)
	processedRT := a.processArguments(runtimeArguments, resolvedEnv, runtimeVars)
	processedPkg := a.processArguments(packageArguments, resolvedEnv, runtimeVars)

	switch registryName {
	case "npm":
		config["command"] = cond(runtimeHint, "npx")
		args := append([]interface{}{"-y", packageName}, toInterfaceSlice(processedRT)...)
		config["args"] = append(args, toInterfaceSlice(processedPkg)...)
		if len(resolvedEnv) > 0 {
			config["env"] = envToInterface(resolvedEnv)
		}
	case "docker":
		config["command"] = "docker"
		if len(processedRT) > 0 {
			config["args"] = toInterfaceSlice(injectEnvVarsIntoDockerArgs(processedRT, resolvedEnv))
		} else {
			config["args"] = toInterfaceSlice(processDockerArgs([]string{"run", "-i", "--rm", packageName}, resolvedEnv))
		}
	case "pypi":
		config["command"] = cond(runtimeHint, "uvx")
		args := append([]interface{}{packageName}, toInterfaceSlice(processedRT)...)
		config["args"] = append(args, toInterfaceSlice(processedPkg)...)
		if len(resolvedEnv) > 0 {
			config["env"] = envToInterface(resolvedEnv)
		}
	case "homebrew":
		cmd := packageName
		if idx := strings.LastIndex(packageName, "/"); idx >= 0 {
			cmd = packageName[idx+1:]
		}
		config["command"] = cmd
		args := append(toInterfaceSlice(processedRT), toInterfaceSlice(processedPkg)...)
		config["args"] = args
		if len(resolvedEnv) > 0 {
			config["env"] = envToInterface(resolvedEnv)
		}
	default:
		config["command"] = cond(runtimeHint, packageName)
		config["args"] = append(toInterfaceSlice(processedRT), toInterfaceSlice(processedPkg)...)
		if len(resolvedEnv) > 0 {
			config["env"] = envToInterface(resolvedEnv)
		}
	}

	if toolsOverride := serverInfo["_apm_tools_override"]; toolsOverride != nil {
		config["tools"] = toolsOverride
	}
	return config, nil
}

// resolveEnvironmentVariables resolves a list or dict of env var definitions.
//
// In translate mode (SupportsRuntimeEnvSubstitution=true): emits ${NAME} placeholders.
// In legacy mode: resolves from envOverrides or os.Getenv.
func (a *Adapter) resolveEnvironmentVariables(envVars interface{}, envOverrides map[string]interface{}) map[string]string {
	if a.SupportsRuntimeEnvSubstitution {
		return a.translateEnvVars(envVars)
	}
	return a.resolveEnvVarsLegacy(envVars, envOverrides)
}

// translateEnvVars emits ${NAME} placeholders for registry-defined env vars.
func (a *Adapter) translateEnvVars(envVars interface{}) map[string]string {
	result := map[string]string{}
	var placeholderKeys []string

	switch ev := envVars.(type) {
	case map[string]interface{}:
		// Self-defined stdio shape: {NAME: value-or-placeholder}
		for name, rawValue := range ev {
			if name == "" {
				continue
			}
			s, ok := rawValue.(string)
			if !ok {
				result[name] = fmt.Sprintf("%v", rawValue)
				continue
			}
			if HasEnvPlaceholder(s) {
				a.lastLegacyAngleVars = append(a.lastLegacyAngleVars, ExtractLegacyAngleVars(s)...)
				translated := TranslateEnvPlaceholder(s)
				result[name] = translated
				for _, m := range envVarRE.FindAllStringSubmatch(translated, -1) {
					placeholderKeys = append(placeholderKeys, m[1])
				}
			} else if def, ok := defaultGitHubEnv[name]; ok && s == def {
				result[name] = s
			} else {
				result[name] = "${" + name + "}"
				placeholderKeys = append(placeholderKeys, name)
			}
		}
	case []interface{}:
		// Registry-sourced shape: [{name, description, required}, ...]
		for _, item := range ev {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name := strField(m, "name")
			if name == "" {
				continue
			}
			if _, isDefault := defaultGitHubEnv[name]; isDefault {
				result[name] = defaultGitHubEnv[name]
			} else {
				result[name] = "${" + name + "}"
				placeholderKeys = append(placeholderKeys, name)
			}
		}
	}

	a.lastEnvPlaceholderKeys = append(a.lastEnvPlaceholderKeys, placeholderKeys...)
	return result
}

// resolveEnvVarsLegacy resolves env vars from overrides or os.Getenv (legacy mode).
func (a *Adapter) resolveEnvVarsLegacy(envVars interface{}, envOverrides map[string]interface{}) map[string]string {
	result := map[string]string{}
	if envOverrides == nil {
		envOverrides = map[string]interface{}{}
	}
	switch ev := envVars.(type) {
	case map[string]interface{}:
		for name, rawValue := range ev {
			s, _ := rawValue.(string)
			if ov, ok := envOverrides[name]; ok {
				result[name] = fmt.Sprintf("%v", ov)
			} else if val := os.Getenv(name); val != "" {
				result[name] = val
			} else if s != "" {
				result[name] = s
			}
		}
	case []interface{}:
		for _, item := range ev {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name := strField(m, "name")
			if name == "" {
				continue
			}
			if ov, ok := envOverrides[name]; ok {
				result[name] = fmt.Sprintf("%v", ov)
			} else if val := os.Getenv(name); val != "" {
				result[name] = val
			}
		}
	}
	return result
}

// resolveEnvVarsDict translates a dict-shaped env block.
func (a *Adapter) resolveEnvVarsDict(env map[string]interface{}, envOverrides map[string]interface{}) map[string]string {
	return a.resolveEnvironmentVariables(env, envOverrides)
}

// resolveEnvVariable resolves a single header/env value.
func (a *Adapter) resolveEnvVariable(name, value string, envOverrides map[string]interface{}) string {
	if a.SupportsRuntimeEnvSubstitution && HasEnvPlaceholder(value) {
		a.lastLegacyAngleVars = append(a.lastLegacyAngleVars, ExtractLegacyAngleVars(value)...)
		translated := TranslateEnvPlaceholder(value)
		for _, m := range envVarRE.FindAllStringSubmatch(translated, -1) {
			a.lastEnvPlaceholderKeys = append(a.lastEnvPlaceholderKeys, m[1])
		}
		return translated
	}
	if envOverrides != nil {
		if ov, ok := envOverrides[name]; ok {
			return fmt.Sprintf("%v", ov)
		}
	}
	if val := os.Getenv(name); val != "" {
		return val
	}
	return value
}

// resolveVariablePlaceholders resolves ${input:VAR} and env-var references in a single arg.
func (a *Adapter) resolveVariablePlaceholders(arg string, resolvedEnv map[string]string, runtimeVars map[string]string) string {
	// Replace ${input:KEY} from runtimeVars.
	inputRE := regexp.MustCompile(`\$\{input:([^}]+)\}`)
	arg = inputRE.ReplaceAllStringFunc(arg, func(m string) string {
		sub := inputRE.FindStringSubmatch(m)
		if v, ok := runtimeVars[sub[1]]; ok {
			return v
		}
		return m
	})
	// Replace ${VAR} / ${env:VAR} from resolvedEnv.
	arg = envVarRE.ReplaceAllStringFunc(arg, func(m string) string {
		sub := envVarRE.FindStringSubmatch(m)
		if v, ok := resolvedEnv[sub[1]]; ok {
			return v
		}
		return m
	})
	return arg
}

// processArguments resolves placeholders in a list of argument strings.
func (a *Adapter) processArguments(args []string, resolvedEnv map[string]string, runtimeVars map[string]string) []string {
	out := make([]string, len(args))
	for i, arg := range args {
		out[i] = a.resolveVariablePlaceholders(arg, resolvedEnv, runtimeVars)
	}
	return out
}

// isGitHubServer returns true when the server or URL is hosted on github.com.
func (a *Adapter) isGitHubServer(name, url string) bool {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "github") {
		return true
	}
	lurl := strings.ToLower(url)
	return strings.Contains(lurl, "github.com") || strings.Contains(lurl, "api.github.com")
}

// getGitHubToken retrieves a GitHub token from the environment.
func (a *Adapter) getGitHubToken() string {
	for _, k := range []string{
		"GITHUB_COPILOT_PAT",
		"GITHUB_TOKEN",
		"GITHUB_APM_PAT",
		"GITHUB_PERSONAL_ACCESS_TOKEN",
	} {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// FormatResolveEnv is an exported wrapper for resolveEnvironmentVariables,
// used by sibling adapter packages (gemini, vscode, etc.) that embed Adapter.
func (a *Adapter) FormatResolveEnv(envVars interface{}, envOverrides map[string]interface{}) map[string]string {
	return a.resolveEnvironmentVariables(envVars, envOverrides)
}

// FormatProcessArgs is an exported wrapper for processArguments,
// used by sibling adapter packages.
func (a *Adapter) FormatProcessArgs(args []string, resolvedEnv map[string]string, runtimeVars map[string]string) []string {
	return a.processArguments(args, resolvedEnv, runtimeVars)
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

func strField(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
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

// selectBestPackage prefers npm, then docker, then others.
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

// inferRegistryName returns the registry type for a package entry.
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
	name := strField(pkg, "name")
	if strings.HasPrefix(name, "@") || strings.Contains(name, "/") {
		return "npm"
	}
	return "npm"
}

// processDockerArgs builds docker args injecting env vars as -e KEY=VALUE.
func processDockerArgs(base []string, env map[string]string) []string {
	out := make([]string, len(base))
	copy(out, base)
	for k, v := range env {
		out = append(out, "-e", k+"="+v)
	}
	return out
}

// injectEnvVarsIntoDockerArgs injects env vars into an existing docker arg list.
func injectEnvVarsIntoDockerArgs(args []string, env map[string]string) []string {
	if len(env) == 0 {
		return args
	}
	out := make([]string, len(args))
	copy(out, args)
	for k, v := range env {
		out = append(out, "-e", k+"="+v)
	}
	return out
}

func cond(preferred, fallback string) string {
	if preferred != "" {
		return preferred
	}
	return fallback
}

func intersect(a, b []string) []string {
	mb := map[string]bool{}
	for _, s := range b {
		mb[s] = true
	}
	var out []string
	for _, s := range a {
		if mb[s] {
			out = append(out, s)
		}
	}
	return out
}

func union(a, b []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range a {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for _, s := range b {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
