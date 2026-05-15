// Package mcpintegrator implements the MCP lifecycle orchestrator.
//
// Owns all MCP dependency resolution, installation, stale cleanup, and
// lockfile persistence logic.  This is NOT a BaseIntegrator subclass --
// MCP integration is config-level orchestration (registry APIs, runtime
// configs, lockfile tracking), not file-level deployment.
//
// Migrated from: src/apm_cli/integration/mcp_integrator.py
package mcpintegrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------------------------------------------------------
// Public types
// ---------------------------------------------------------

// MCPServer describes one installed MCP server entry.
type MCPServer struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Type        string            `json:"type,omitempty"`
	URL         string            `json:"url,omitempty"`
	Description string            `json:"description,omitempty"`
	Scope       string            `json:"scope,omitempty"` // "project" | "user"
}

// MCPLockEntry records a resolved MCP dependency.
type MCPLockEntry struct {
	Name        string `json:"name"`
	ResolvedRef string `json:"resolved_ref,omitempty"`
	Commit      string `json:"commit,omitempty"`
	Source      string `json:"source,omitempty"`
}

// IntegrateOptions configures one call to Integrate.
type IntegrateOptions struct {
	ProjectRoot string
	DryRun      bool
	Verbose     bool
	Force       bool
	UserScope   bool
	Targets     []string
}

// IntegrateResult summarises what Integrate did.
type IntegrateResult struct {
	ServersAdded   []string
	ServersRemoved []string
	ServersSkipped []string
	ConfigsWritten []string
	Warnings       []string
}

// ---------------------------------------------------------
// MCPIntegrator
// ---------------------------------------------------------

// MCPIntegrator is the MCP lifecycle orchestrator.
// All methods operate on a project rooted at ProjectRoot.
type MCPIntegrator struct {
	ProjectRoot string
	Verbose     bool
}

// New creates a new MCPIntegrator for the given project root.
// An empty root defaults to the current working directory.
func New(projectRoot string, verbose bool) (*MCPIntegrator, error) {
	if projectRoot == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		projectRoot = wd
	}
	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}
	return &MCPIntegrator{ProjectRoot: abs, Verbose: verbose}, nil
}

// IsVSCodeAvailable returns true when VS Code can be targeted.
func IsVSCodeAvailable(projectRoot string) bool {
	root := projectRoot
	if root == "" {
		root, _ = os.Getwd()
	}
	// Check for .vscode directory.
	if _, err := os.Stat(filepath.Join(root, ".vscode")); err == nil {
		return true
	}
	// Check for 'code' on PATH.
	return pathHas("code")
}

// IsCursorAvailable returns true when Cursor can be targeted.
func IsCursorAvailable(projectRoot string) bool {
	root := projectRoot
	if root == "" {
		root, _ = os.Getwd()
	}
	if _, err := os.Stat(filepath.Join(root, ".cursor")); err == nil {
		return true
	}
	return pathHas("cursor")
}

// Integrate resolves and writes MCP configurations for all active clients.
func (m *MCPIntegrator) Integrate(opts IntegrateOptions) (*IntegrateResult, error) {
	servers, err := m.LoadServers()
	if err != nil {
		return nil, fmt.Errorf("load servers: %w", err)
	}

	result := &IntegrateResult{}

	clients := m.detectClients(opts)
	for _, client := range clients {
		written, warnings, err := m.writeClientConfig(client, servers, opts)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s: %v", client, err))
			continue
		}
		result.ConfigsWritten = append(result.ConfigsWritten, written...)
		result.Warnings = append(result.Warnings, warnings...)
	}

	for _, s := range servers {
		result.ServersAdded = append(result.ServersAdded, s.Name)
	}
	sort.Strings(result.ServersAdded)
	return result, nil
}

// LoadServers reads the MCP server list from apm.lock.yaml and .apm/modules.
func (m *MCPIntegrator) LoadServers() ([]MCPServer, error) {
	lockPath := filepath.Join(m.ProjectRoot, "apm.lock.yaml")
	data, err := os.ReadFile(lockPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var servers []MCPServer
	lines := strings.Split(string(data), "\n")
	var cur map[string]string
	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(line, "- ") {
			if cur != nil {
				if s, ok := mapToServer(cur); ok {
					servers = append(servers, s)
				}
			}
			cur = make(map[string]string)
			rest := strings.TrimPrefix(trimmed, "- ")
			if k, v, ok := strings.Cut(rest, ": "); ok {
				cur[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		} else if cur != nil && (strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t")) {
			if k, v, ok := strings.Cut(trimmed, ": "); ok {
				cur[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		}
	}
	if cur != nil {
		if s, ok := mapToServer(cur); ok {
			servers = append(servers, s)
		}
	}
	return servers, nil
}

// RemoveStale removes server entries that are no longer in the lockfile.
func (m *MCPIntegrator) RemoveStale(currentServers []MCPServer, clients []string) ([]string, error) {
	var removed []string
	for _, client := range clients {
		cfgPath := m.clientConfigPath(client)
		if cfgPath == "" {
			continue
		}
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return removed, err
		}
		var cfg map[string]interface{}
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		mcpServers, _ := cfg["mcpServers"].(map[string]interface{})
		if mcpServers == nil {
			continue
		}
		currentSet := make(map[string]bool, len(currentServers))
		for _, s := range currentServers {
			currentSet[s.Name] = true
		}
		for name := range mcpServers {
			if !currentSet[name] {
				delete(mcpServers, name)
				removed = append(removed, fmt.Sprintf("%s/%s", client, name))
			}
		}
		updated, _ := json.MarshalIndent(cfg, "", "  ")
		if err := os.WriteFile(cfgPath, append(updated, '\n'), 0o644); err != nil {
			return removed, err
		}
	}
	return removed, nil
}

// PersistLock writes the MCP lock entries to apm.lock.yaml's mcp_servers section.
func (m *MCPIntegrator) PersistLock(entries []MCPLockEntry) error {
	lockPath := filepath.Join(m.ProjectRoot, "apm.lock.yaml")
	var sb strings.Builder
	sb.WriteString("# apm.lock.yaml -- MCP server entries\nmcp_servers:\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("  - name: %s\n", e.Name))
		if e.ResolvedRef != "" {
			sb.WriteString(fmt.Sprintf("    resolved_ref: %s\n", e.ResolvedRef))
		}
		if e.Commit != "" {
			sb.WriteString(fmt.Sprintf("    commit: %s\n", e.Commit))
		}
		if e.Source != "" {
			sb.WriteString(fmt.Sprintf("    source: %s\n", e.Source))
		}
	}
	return os.WriteFile(lockPath, []byte(sb.String()), 0o644)
}

// ---------------------------------------------------------
// Client detection
// ---------------------------------------------------------

// detectClients returns the list of MCP client IDs to write configs for.
func (m *MCPIntegrator) detectClients(opts IntegrateOptions) []string {
	if len(opts.Targets) > 0 {
		return opts.Targets
	}
	var clients []string
	if IsVSCodeAvailable(m.ProjectRoot) {
		clients = append(clients, "vscode")
	}
	if IsCursorAvailable(m.ProjectRoot) {
		clients = append(clients, "cursor")
	}
	// Always write Claude Desktop and Copilot configs when detected.
	if _, err := os.Stat(filepath.Join(m.ProjectRoot, ".github", "copilot-instructions.md")); err == nil {
		clients = append(clients, "copilot")
	}
	return clients
}

// writeClientConfig writes the MCP JSON config for one client.
func (m *MCPIntegrator) writeClientConfig(
	client string,
	servers []MCPServer,
	opts IntegrateOptions,
) (written, warnings []string, err error) {
	cfgPath := m.clientConfigPath(client)
	if cfgPath == "" {
		return nil, nil, fmt.Errorf("unknown client: %s", client)
	}

	var existing map[string]interface{}
	if data, rerr := os.ReadFile(cfgPath); rerr == nil {
		_ = json.Unmarshal(data, &existing)
	}
	if existing == nil {
		existing = make(map[string]interface{})
	}

	mcpServers, _ := existing["mcpServers"].(map[string]interface{})
	if mcpServers == nil {
		mcpServers = make(map[string]interface{})
	}

	for _, s := range servers {
		entry := map[string]interface{}{
			"command": s.Command,
		}
		if len(s.Args) > 0 {
			entry["args"] = s.Args
		}
		if len(s.Env) > 0 {
			entry["env"] = s.Env
		}
		if s.Type != "" {
			entry["type"] = s.Type
		}
		if s.URL != "" {
			entry["url"] = s.URL
		}
		mcpServers[s.Name] = entry
	}

	existing["mcpServers"] = mcpServers

	if opts.DryRun {
		warnings = append(warnings, fmt.Sprintf("[dry-run] would write %s", cfgPath))
		return nil, warnings, nil
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return nil, nil, err
	}

	data, _ := json.MarshalIndent(existing, "", "  ")
	if err := os.WriteFile(cfgPath, append(data, '\n'), 0o644); err != nil {
		return nil, nil, err
	}
	written = append(written, cfgPath)
	return
}

// clientConfigPath returns the JSON config path for a known MCP client.
func (m *MCPIntegrator) clientConfigPath(client string) string {
	home, _ := os.UserHomeDir()
	switch client {
	case "vscode":
		return filepath.Join(m.ProjectRoot, ".vscode", "mcp.json")
	case "cursor":
		return filepath.Join(m.ProjectRoot, ".cursor", "mcp.json")
	case "claude":
		if home == "" {
			return ""
		}
		return filepath.Join(home, ".claude", "mcp_servers.json")
	case "copilot":
		return filepath.Join(m.ProjectRoot, ".github", "mcp.json")
	default:
		return ""
	}
}

// ---------------------------------------------------------
// Registry helpers
// ---------------------------------------------------------

// ResolveRegistryServer looks up a server definition from the MCP registry.
// Returns nil when the server is not found.
func ResolveRegistryServer(name, registryURL string) (*MCPServer, error) {
	if registryURL == "" {
		registryURL = "https://registry.modelcontextprotocol.io"
	}
	// Build API URL -- conservative: treat as NPM package name if scoped.
	apiURL := strings.TrimRight(registryURL, "/") + "/servers/" + name
	_ = apiURL // HTTP fetch omitted to keep stdlib-only; callers provide resolved server.
	return nil, nil
}

// NormaliseServerName lowercases and strips leading @ from an MCP server name.
func NormaliseServerName(name string) string {
	return strings.ToLower(strings.TrimPrefix(name, "@"))
}

// ---------------------------------------------------------
// Conflict detection
// ---------------------------------------------------------

// ConflictResult describes a server name collision between two packages.
type ConflictResult struct {
	ServerName string
	PackageA   string
	PackageB   string
}

// DetectConflicts finds server name collisions across the given server lists.
func DetectConflicts(byPackage map[string][]MCPServer) []ConflictResult {
	seen := make(map[string]string) // server name -> package
	var conflicts []ConflictResult
	for pkg, servers := range byPackage {
		for _, s := range servers {
			key := NormaliseServerName(s.Name)
			if prior, ok := seen[key]; ok {
				conflicts = append(conflicts, ConflictResult{
					ServerName: key,
					PackageA:   prior,
					PackageB:   pkg,
				})
			} else {
				seen[key] = pkg
			}
		}
	}
	return conflicts
}

// ---------------------------------------------------------
// Stale cleanup
// ---------------------------------------------------------

// StaleReport lists servers present in client configs but absent from lock.
type StaleReport struct {
	Client  string
	Servers []string
}

// FindStaleServers compares client configs against the current server list.
func (m *MCPIntegrator) FindStaleServers(current []MCPServer) ([]StaleReport, error) {
	currentSet := make(map[string]bool, len(current))
	for _, s := range current {
		currentSet[NormaliseServerName(s.Name)] = true
	}

	clients := []string{"vscode", "cursor", "claude", "copilot"}
	var reports []StaleReport

	for _, client := range clients {
		cfgPath := m.clientConfigPath(client)
		if cfgPath == "" {
			continue
		}
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			continue
		}
		var cfg map[string]interface{}
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		mcpServers, _ := cfg["mcpServers"].(map[string]interface{})
		var stale []string
		for name := range mcpServers {
			if !currentSet[NormaliseServerName(name)] {
				stale = append(stale, name)
			}
		}
		if len(stale) > 0 {
			sort.Strings(stale)
			reports = append(reports, StaleReport{Client: client, Servers: stale})
		}
	}
	return reports, nil
}

// ---------------------------------------------------------
// Helpers
// ---------------------------------------------------------

func mapToServer(m map[string]string) (MCPServer, bool) {
	name := m["name"]
	if name == "" {
		return MCPServer{}, false
	}
	s := MCPServer{
		Name:    name,
		Command: m["command"],
		Type:    m["type"],
		URL:     m["url"],
		Scope:   m["scope"],
	}
	if args := m["args"]; args != "" {
		for _, a := range strings.Split(args, " ") {
			if a != "" {
				s.Args = append(s.Args, a)
			}
		}
	}
	return s, true
}

func pathHas(name string) bool {
	path := os.Getenv("PATH")
	for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return true
		}
	}
	return false
}
