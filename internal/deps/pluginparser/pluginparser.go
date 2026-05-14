// Package pluginparser parses Claude Code plugin.json manifests and
// synthesises apm.yml files from plugin directory layouts.
//
// Migrated from: src/apm_cli/deps/plugin_parser.py
package pluginparser

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// PluginManifest holds the optional metadata from plugin.json.
type PluginManifest struct {
	Name       string          `json:"name"`
	MCPServers json.RawMessage `json:"mcpServers,omitempty"`
	Agents     []string        `json:"agents,omitempty"`
	Skills     []string        `json:"skills,omitempty"`
	Commands   []string        `json:"commands,omitempty"`
	Hooks      json.RawMessage `json:"hooks,omitempty"`
	Extra      map[string]json.RawMessage
}

// MCPServerConfig holds a single MCP server configuration.
type MCPServerConfig struct {
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	URL       string            `json:"url,omitempty"`
	Type      string            `json:"type,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Tools     []string          `json:"tools,omitempty"`
}

// MCPDepEntry is a dependency entry generated from an MCP server config.
type MCPDepEntry struct {
	Name          string
	Transport     string
	Command       string
	Args          []string
	URL           string
	Headers       map[string]string
	Env           map[string]string
	Tools         []string
	Registry      bool
}

// ParsePluginManifest parses a plugin.json file at the given path.
// Returns the parsed manifest or an error.
func ParsePluginManifest(pluginJSONPath string) (*PluginManifest, error) {
	if _, err := os.Stat(pluginJSONPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin.json not found: %s", pluginJSONPath)
	}
	data, err := os.ReadFile(pluginJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin.json: %w", err)
	}
	var manifest PluginManifest
	if err2 := json.Unmarshal(data, &manifest); err2 != nil {
		return nil, fmt.Errorf("invalid JSON in plugin.json: %w", err2)
	}
	if manifest.Name == "" {
		log.Printf("plugin.json at %s is missing 'name' field; falling back to directory name", pluginJSONPath)
	}
	return &manifest, nil
}

// NormalizePluginDirectory normalises a Claude plugin directory into an APM package.
//
// Works with or without plugin.json. Returns the path to the generated apm.yml.
func NormalizePluginDirectory(pluginPath string, pluginJSONPath string) (string, error) {
	var manifest *PluginManifest

	if pluginJSONPath != "" {
		if _, err := os.Stat(pluginJSONPath); err == nil {
			m, err2 := ParsePluginManifest(pluginJSONPath)
			if err2 != nil {
				// Treat as empty manifest; fall back to dir-name defaults
				m = &PluginManifest{}
			}
			manifest = m
		}
	}

	if manifest == nil {
		manifest = &PluginManifest{}
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(pluginPath)
	}

	return SynthesizeApmYMLFromPlugin(pluginPath, manifest)
}

// SynthesizeApmYMLFromPlugin synthesises apm.yml from plugin metadata.
func SynthesizeApmYMLFromPlugin(pluginPath string, manifest *PluginManifest) (string, error) {
	if manifest.Name == "" {
		manifest.Name = filepath.Base(pluginPath)
	}

	// Create .apm directory structure
	apmDir := filepath.Join(pluginPath, ".apm")
	if err := os.MkdirAll(apmDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create .apm directory: %w", err)
	}

	// Map plugin structure into .apm/ subdirectories
	if err := mapPluginArtifacts(pluginPath, apmDir, manifest); err != nil {
		return "", err
	}

	// Extract MCP servers
	mcpServers, err := extractMCPServers(pluginPath, manifest)
	if err != nil {
		log.Printf("failed to extract MCP servers from plugin %s: %v", pluginPath, err)
	}

	var mcpDeps []MCPDepEntry
	if len(mcpServers) > 0 {
		mcpDeps = mcpServersToDeps(mcpServers, pluginPath)
	}

	// Generate apm.yml
	content := generateApmYML(manifest, mcpDeps)
	apmYMLPath := filepath.Join(pluginPath, "apm.yml")
	if err2 := os.WriteFile(apmYMLPath, []byte(content), 0o644); err2 != nil {
		return "", fmt.Errorf("failed to write apm.yml: %w", err2)
	}

	return apmYMLPath, nil
}

// extractMCPServers reads MCP server definitions from the plugin manifest.
func extractMCPServers(pluginPath string, manifest *PluginManifest) (map[string]MCPServerConfig, error) {
	logger := log.Default()

	if manifest.MCPServers == nil {
		// Fall back to auto-discovery
		servers := map[string]MCPServerConfig{}
		for _, candidate := range []string{".mcp.json", filepath.Join(".github", ".mcp.json")} {
			fullPath := filepath.Join(pluginPath, candidate)
			info, err := os.Lstat(fullPath)
			if err == nil && info.Mode()&fs.ModeSymlink == 0 && info.Mode().IsRegular() {
				s, err2 := readMCPJSON(fullPath)
				if err2 == nil && len(s) > 0 {
					servers = s
					break
				}
			}
		}
		if len(servers) > 0 {
			return substitutePlaceholder(servers, pluginPath, logger), nil
		}
		return servers, nil
	}

	// Determine type of mcpServers value
	raw := manifest.MCPServers
	var servers map[string]MCPServerConfig

	// Try dict
	if err := json.Unmarshal(raw, &servers); err == nil {
		return substitutePlaceholder(servers, pluginPath, logger), nil
	}

	// Try string (file path)
	var strVal string
	if err := json.Unmarshal(raw, &strVal); err == nil {
		s, err2 := readMCPFile(pluginPath, strVal)
		if err2 != nil {
			logger.Printf("MCP file read failed: %v", err2)
			return map[string]MCPServerConfig{}, nil
		}
		return substitutePlaceholder(s, pluginPath, logger), nil
	}

	// Try array of string paths
	var arrVal []string
	if err := json.Unmarshal(raw, &arrVal); err == nil {
		result := map[string]MCPServerConfig{}
		for _, entry := range arrVal {
			s, err2 := readMCPFile(pluginPath, entry)
			if err2 != nil {
				logger.Printf("MCP file read failed: %v", err2)
				continue
			}
			for k, v := range s {
				result[k] = v
			}
		}
		return substitutePlaceholder(result, pluginPath, logger), nil
	}

	logger.Printf("unsupported mcpServers type in plugin %s", pluginPath)
	return map[string]MCPServerConfig{}, nil
}

// readMCPFile reads a JSON file at relPath relative to pluginPath and returns its mcpServers dict.
func readMCPFile(pluginPath, relPath string) (map[string]MCPServerConfig, error) {
	absPlug, _ := filepath.Abs(pluginPath)
	target := filepath.Join(absPlug, relPath)
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", relPath)
	}
	// Security: must stay inside pluginPath
	if !strings.HasPrefix(absTarget, absPlug+string(os.PathSeparator)) {
		return nil, fmt.Errorf("MCP file path escapes plugin root: %s", relPath)
	}
	info, err := os.Lstat(absTarget)
	if err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("MCP file not found or invalid: %s", absTarget)
	}
	return readMCPJSON(absTarget)
}

// readMCPJSON parses a JSON file and returns the mcpServers dict.
func readMCPJSON(path string) (map[string]MCPServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		MCPServers map[string]MCPServerConfig `json:"mcpServers"`
	}
	if err2 := json.Unmarshal(data, &wrapper); err2 != nil {
		return nil, err2
	}
	if wrapper.MCPServers == nil {
		return map[string]MCPServerConfig{}, nil
	}
	return wrapper.MCPServers, nil
}

// substitutePlaceholder replaces ${CLAUDE_PLUGIN_ROOT} in string values.
func substitutePlaceholder(servers map[string]MCPServerConfig, pluginPath string, _ *log.Logger) map[string]MCPServerConfig {
	absRoot, _ := filepath.Abs(pluginPath)
	placeholder := "${CLAUDE_PLUGIN_ROOT}"

	replaceStr := func(s string) string {
		return strings.ReplaceAll(s, placeholder, absRoot)
	}

	result := make(map[string]MCPServerConfig, len(servers))
	for name, cfg := range servers {
		cfg.Command = replaceStr(cfg.Command)
		cfg.URL = replaceStr(cfg.URL)
		newArgs := make([]string, len(cfg.Args))
		for i, a := range cfg.Args {
			newArgs[i] = replaceStr(a)
		}
		cfg.Args = newArgs
		if cfg.Env != nil {
			newEnv := make(map[string]string, len(cfg.Env))
			for k, v := range cfg.Env {
				newEnv[k] = replaceStr(v)
			}
			cfg.Env = newEnv
		}
		result[name] = cfg
	}
	return result
}

// mcpServersToDeps converts raw MCP server configs to dependency dicts.
func mcpServersToDeps(servers map[string]MCPServerConfig, pluginPath string) []MCPDepEntry {
	var deps []MCPDepEntry
	for name, cfg := range servers {
		dep := MCPDepEntry{Name: name, Registry: false}
		if cfg.Command != "" {
			dep.Transport = "stdio"
			dep.Command = cfg.Command
			dep.Args = cfg.Args
		} else if cfg.URL != "" {
			transport := cfg.Type
			validTransports := map[string]bool{"http": true, "sse": true, "streamable-http": true}
			if !validTransports[transport] {
				transport = "http"
			}
			dep.Transport = transport
			dep.URL = cfg.URL
			dep.Headers = cfg.Headers
		} else {
			log.Printf("skipping MCP server %q from plugin %q: no 'command' or 'url'", name, filepath.Base(pluginPath))
			continue
		}
		dep.Env = cfg.Env
		dep.Tools = cfg.Tools
		deps = append(deps, dep)
	}
	return deps
}

// mapPluginArtifacts copies plugin components to .apm/ subdirectories.
func mapPluginArtifacts(pluginPath, apmDir string, manifest *PluginManifest) error {
	type mapping struct {
		src  string
		dst  string
		isDir bool
	}

	// Standard component mappings
	componentMappings := []mapping{
		{"agents", filepath.Join(apmDir, "agents"), true},
		{"skills", filepath.Join(apmDir, "skills"), true},
		{"commands", filepath.Join(apmDir, "prompts"), true},
		{"hooks", filepath.Join(apmDir, "hooks"), true},
	}

	for _, m := range componentMappings {
		srcPath := filepath.Join(pluginPath, m.src)
		info, err := os.Lstat(srcPath)
		if err != nil || info.Mode()&fs.ModeSymlink != 0 {
			continue
		}
		if !info.IsDir() {
			continue
		}
		// Verify path is within plugin root
		abs, _ := filepath.Abs(srcPath)
		absPlugin, _ := filepath.Abs(pluginPath)
		if !strings.HasPrefix(abs, absPlugin+string(os.PathSeparator)) {
			continue
		}
		if err2 := copyDir(srcPath, m.dst); err2 != nil {
			log.Printf("warning: failed to copy %s to %s: %v", srcPath, m.dst, err2)
		}
	}

	// Pass-through files
	passthroughs := []string{".mcp.json", ".lsp.json", "settings.json"}
	for _, fname := range passthroughs {
		src := filepath.Join(pluginPath, fname)
		info, err := os.Lstat(src)
		if err != nil || info.Mode()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() {
			continue
		}
		dst := filepath.Join(apmDir, fname)
		if err2 := copyFile(src, dst); err2 != nil {
			log.Printf("warning: failed to copy %s: %v", fname, err2)
		}
	}

	return nil
}

// copyFile copies a single regular file.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		// Skip symlinks
		info, err2 := os.Lstat(srcPath)
		if err2 != nil || info.Mode()&fs.ModeSymlink != 0 {
			continue
		}
		if entry.IsDir() {
			if err3 := copyDir(srcPath, dstPath); err3 != nil {
				log.Printf("warning: copyDir %s: %v", srcPath, err3)
			}
		} else {
			if err3 := copyFile(srcPath, dstPath); err3 != nil {
				log.Printf("warning: copyFile %s: %v", srcPath, err3)
			}
		}
	}
	return nil
}

// generateApmYML generates the apm.yml content from plugin metadata.
func generateApmYML(manifest *PluginManifest, mcpDeps []MCPDepEntry) string {
	var sb strings.Builder
	sb.WriteString("# Generated by APM from Claude plugin\n")
	sb.WriteString("name: ")
	sb.WriteString(yamlString(manifest.Name))
	sb.WriteString("\n\n")

	if len(mcpDeps) > 0 {
		sb.WriteString("dependencies:\n  mcp:\n")
		for _, dep := range mcpDeps {
			sb.WriteString("  - name: ")
			sb.WriteString(yamlString(dep.Name))
			sb.WriteString("\n    registry: false\n")
			sb.WriteString("    transport: ")
			sb.WriteString(dep.Transport)
			sb.WriteString("\n")
			if dep.Command != "" {
				sb.WriteString("    command: ")
				sb.WriteString(yamlString(dep.Command))
				sb.WriteString("\n")
				if len(dep.Args) > 0 {
					sb.WriteString("    args:\n")
					for _, a := range dep.Args {
						sb.WriteString("    - ")
						sb.WriteString(yamlString(a))
						sb.WriteString("\n")
					}
				}
			}
			if dep.URL != "" {
				sb.WriteString("    url: ")
				sb.WriteString(dep.URL)
				sb.WriteString("\n")
			}
			if len(dep.Env) > 0 {
				sb.WriteString("    env:\n")
				for k, v := range dep.Env {
					sb.WriteString("      ")
					sb.WriteString(k)
					sb.WriteString(": ")
					sb.WriteString(yamlString(v))
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String()
}

// yamlString wraps a string in quotes if needed.
func yamlString(s string) string {
	if strings.ContainsAny(s, ":{}[]|>&*!,#?@`\"'\\") ||
		strings.Contains(s, " ") ||
		strings.Contains(s, "\n") {
		escaped := strings.ReplaceAll(s, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return s
}
