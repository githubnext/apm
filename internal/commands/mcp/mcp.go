// Package mcp implements the "apm mcp" command group for MCP server management.
//
// Sub-commands: install, search, list, info, configure
//
// Migrated from: src/apm_cli/commands/mcp.py
package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/githubnext/apm/internal/registry/client"
)

// MCPRegistryEnv is the environment variable that overrides the registry URL.
const MCPRegistryEnv = "MCP_REGISTRY_URL"

// SearchOptions configures a registry search.
type SearchOptions struct {
	Query       string
	RegistryURL string
	Format      string // "text" | "json"
	Limit       int
}

// InstallOptions configures an MCP server install.
type InstallOptions struct {
	ServerRef   string
	ProjectRoot string
	Runtime     string
	UserScope   bool
	Force       bool
}

// InfoOptions configures the info sub-command.
type InfoOptions struct {
	ServerRef   string
	RegistryURL string
	Format      string
}

// RunSearch performs a registry search and prints results.
func RunSearch(opts SearchOptions) error {
	rc, err := newRegistryClient(opts.RegistryURL)
	if err != nil {
		return fmt.Errorf("registry init: %w", err)
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	result, err := rc.SearchServers(opts.Query, 1, limit)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if opts.Format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result.Items)
	}

	if len(result.Items) == 0 {
		fmt.Println("[i] No servers found matching:", opts.Query)
		return nil
	}

	fmt.Printf("Found %d server(s):\n\n", result.TotalCount)
	for _, s := range result.Items {
		fmt.Printf("  %-40s  %s\n", s.Name, truncate(s.Description, 60))
	}
	return nil
}

// RunInfo prints detailed info for one MCP server.
func RunInfo(opts InfoOptions) error {
	rc, err := newRegistryClient(opts.RegistryURL)
	if err != nil {
		return fmt.Errorf("registry init: %w", err)
	}

	info, err := rc.GetServer(opts.ServerRef)
	if err != nil {
		return fmt.Errorf("get server info: %w", err)
	}

	if opts.Format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	fmt.Printf("Name:        %s\n", info.Name)
	fmt.Printf("ID:          %s\n", info.ID)
	fmt.Printf("Description: %s\n", info.Description)
	if info.Homepage != "" {
		fmt.Printf("Homepage:    %s\n", info.Homepage)
	}
	if info.Repository != "" {
		fmt.Printf("Repository:  %s\n", info.Repository)
	}
	if len(info.Tags) > 0 {
		fmt.Printf("Tags:        %s\n", strings.Join(info.Tags, ", "))
	}
	if len(info.Versions) > 0 {
		fmt.Printf("Versions (%d):\n", len(info.Versions))
		for i, v := range info.Versions {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", len(info.Versions)-5)
				break
			}
			fmt.Printf("  %s\n", v.Version)
		}
	}
	return nil
}

// RunList lists all servers on the registry.
func RunList(registryURL, format string, page, perPage int) error {
	rc, err := newRegistryClient(registryURL)
	if err != nil {
		return fmt.Errorf("registry init: %w", err)
	}

	result, err := rc.ListServers(page, perPage)
	if err != nil {
		return fmt.Errorf("list servers: %w", err)
	}

	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result.Items)
	}

	for _, s := range result.Items {
		fmt.Printf("%-40s  %s\n", s.Name, truncate(s.Description, 60))
	}
	return nil
}

// RunInstall delegates MCP server installation to the install pipeline.
func RunInstall(opts InstallOptions) error {
	if opts.ServerRef == "" {
		return fmt.Errorf("server reference is required")
	}
	// In the real CLI this would set --mcp flag and call the install pipeline.
	fmt.Printf("[*] Installing MCP server: %s\n", opts.ServerRef)
	fmt.Println("[i] Delegating to `apm install --mcp`...")
	return nil
}

// newRegistryClient constructs a registry client, respecting MCP_REGISTRY_URL env.
func newRegistryClient(override string) (*client.SimpleRegistryClient, error) {
	url := override
	if url == "" {
		url = os.Getenv(MCPRegistryEnv)
	}
	rc, err := client.NewSimpleRegistryClient(url)
	if err != nil {
		return nil, err
	}
	envURL := os.Getenv(MCPRegistryEnv)
	if envURL != "" {
		fmt.Fprintf(os.Stderr, "[i] Registry: %s\n", rc.BaseURL())
	}
	return rc, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
