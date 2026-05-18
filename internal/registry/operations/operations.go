// Package operations implements MCP server install and conflict-detection logic.
//
// Migrated from: src/apm_cli/registry/operations.py
package operations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/githubnext/apm/internal/registry/client"
)

// ServerNeed records whether a server reference needs installation.
type ServerNeed struct {
	Reference  string
	NeedsInstall bool
	Reason     string
}

// InstallStatus summarises per-runtime installation state for one server.
type InstallStatus struct {
	Runtime   string
	Installed bool
	ServerID  string
}

// MCPServerOperations handles conflict detection and installation status for MCP servers.
type MCPServerOperations struct {
	registryClient *client.SimpleRegistryClient
}

// NewMCPServerOperations creates an MCPServerOperations for the given registry URL.
// Passing an empty string uses the default public registry.
func NewMCPServerOperations(registryURL string) (*MCPServerOperations, error) {
	rc, err := client.NewSimpleRegistryClient(registryURL)
	if err != nil {
		return nil, fmt.Errorf("create registry client: %w", err)
	}
	return &MCPServerOperations{registryClient: rc}, nil
}

// CheckServersNeedingInstallation returns the subset of serverRefs that are not yet
// installed in at least one of targetRuntimes.
// maxWorkers bounds the concurrency for registry lookups (default 4).
func (o *MCPServerOperations) CheckServersNeedingInstallation(
	targetRuntimes, serverRefs []string,
	projectRoot string,
	userScope bool,
	maxWorkers int,
) ([]string, error) {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}

	// Pre-load installed IDs per runtime.
	installedByRuntime := make(map[string]map[string]struct{}, len(targetRuntimes))
	for _, rt := range targetRuntimes {
		ids, err := o.getInstalledServerIDs([]string{rt}, projectRoot, userScope)
		if err != nil {
			return nil, fmt.Errorf("get installed IDs for %s: %w", rt, err)
		}
		installedByRuntime[rt] = ids
	}

	type result struct {
		ref    string
		needed bool
	}

	sem := make(chan struct{}, maxWorkers)
	results := make(chan result, len(serverRefs))
	var wg sync.WaitGroup

	for _, ref := range serverRefs {
		wg.Add(1)
		sem <- struct{}{}
		go func(serverRef string) {
			defer wg.Done()
			defer func() { <-sem }()

			needed := o.serverNeedsInstall(serverRef, targetRuntimes, installedByRuntime)
			results <- result{ref: serverRef, needed: needed}
		}(ref)
	}

	wg.Wait()
	close(results)

	var needing []string
	for r := range results {
		if r.needed {
			needing = append(needing, r.ref)
		}
	}
	return needing, nil
}

// serverNeedsInstall checks whether serverRef is installed in all target runtimes.
func (o *MCPServerOperations) serverNeedsInstall(
	serverRef string,
	targetRuntimes []string,
	installedByRuntime map[string]map[string]struct{},
) bool {
	info, err := o.registryClient.GetServer(serverRef)
	if err != nil || info == nil {
		return true
	}
	for _, rt := range targetRuntimes {
		ids, ok := installedByRuntime[rt]
		if !ok {
			return true
		}
		if _, found := ids[info.ID]; !found {
			return true
		}
	}
	return false
}

// getInstalledServerIDs reads the MCP config files for the given runtimes and returns
// the set of installed server IDs (UUIDs).
func (o *MCPServerOperations) getInstalledServerIDs(
	runtimes []string,
	projectRoot string,
	userScope bool,
) (map[string]struct{}, error) {
	ids := make(map[string]struct{})
	for _, rt := range runtimes {
		paths := mcpConfigPaths(rt, projectRoot, userScope)
		for _, p := range paths {
			data, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			extracted, err := extractServerIDs(data)
			if err != nil {
				continue
			}
			for _, id := range extracted {
				ids[id] = struct{}{}
			}
		}
	}
	return ids, nil
}

// mcpConfigPaths returns candidate MCP config file paths for a runtime.
func mcpConfigPaths(runtime, projectRoot string, userScope bool) []string {
	var paths []string
	home, _ := os.UserHomeDir()

	switch strings.ToLower(runtime) {
	case "claude":
		if userScope {
			if home != "" {
				paths = append(paths,
					filepath.Join(home, ".claude", "claude_desktop_config.json"),
					filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"),
				)
			}
		} else if projectRoot != "" {
			paths = append(paths, filepath.Join(projectRoot, ".claude", "claude_mcp_config.json"))
		}
	case "copilot", "vscode":
		if projectRoot != "" {
			paths = append(paths, filepath.Join(projectRoot, ".vscode", "mcp.json"))
		}
		if userScope && home != "" {
			paths = append(paths, filepath.Join(home, ".vscode", "mcp.json"))
		}
	case "cursor":
		if projectRoot != "" {
			paths = append(paths, filepath.Join(projectRoot, ".cursor", "mcp.json"))
		}
	}
	return paths
}

// extractServerIDs parses an MCP config JSON blob and returns all server IDs found.
func extractServerIDs(data []byte) ([]string, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	var ids []string

	// Look for mcpServers or servers key (varies by client).
	for _, key := range []string{"mcpServers", "servers"} {
		serversRaw, ok := raw[key]
		if !ok {
			continue
		}
		var servers map[string]json.RawMessage
		if err := json.Unmarshal(serversRaw, &servers); err != nil {
			continue
		}
		for _, v := range servers {
			var entry map[string]json.RawMessage
			if err := json.Unmarshal(v, &entry); err != nil {
				continue
			}
			if idRaw, ok := entry["id"]; ok {
				var id string
				if err := json.Unmarshal(idRaw, &id); err == nil && id != "" {
					ids = append(ids, id)
				}
			}
		}
	}
	return ids, nil
}

// GetInstallStatus returns per-runtime installation status for each serverRef.
func (o *MCPServerOperations) GetInstallStatus(
	serverRefs, targetRuntimes []string,
	projectRoot string,
	userScope bool,
) (map[string][]InstallStatus, error) {
	installedByRuntime := make(map[string]map[string]struct{}, len(targetRuntimes))
	for _, rt := range targetRuntimes {
		ids, err := o.getInstalledServerIDs([]string{rt}, projectRoot, userScope)
		if err != nil {
			return nil, err
		}
		installedByRuntime[rt] = ids
	}

	out := make(map[string][]InstallStatus)
	for _, ref := range serverRefs {
		info, err := o.registryClient.GetServer(ref)
		if err != nil || info == nil {
			for _, rt := range targetRuntimes {
				out[ref] = append(out[ref], InstallStatus{Runtime: rt, Installed: false})
			}
			continue
		}
		for _, rt := range targetRuntimes {
			ids := installedByRuntime[rt]
			_, installed := ids[info.ID]
			out[ref] = append(out[ref], InstallStatus{
				Runtime:   rt,
				Installed: installed,
				ServerID:  info.ID,
			})
		}
	}
	return out, nil
}
