// Package client implements a simple MCP registry HTTP client for server discovery.
//
// Migrated from: src/apm_cli/registry/client.py
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRegistryURL    = "https://api.mcp.github.com"
	defaultConnectTimeout = 10.0
	defaultReadTimeout    = 30.0
)

// MCPServerInfo holds metadata for a single MCP server entry.
type MCPServerInfo struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Publisher   string         `json:"publisher"`
	Homepage    string         `json:"homepage"`
	Repository  string         `json:"repository"`
	License     string         `json:"license"`
	Tags        []string       `json:"tags"`
	Versions    []VersionEntry `json:"versions"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// VersionEntry holds version metadata for an MCP server.
type VersionEntry struct {
	Version   string `json:"version"`
	CreatedAt string `json:"created_at"`
	PackageID string `json:"package_id"`
}

// SearchResult is the response envelope for registry searches.
type SearchResult struct {
	Items      []MCPServerInfo `json:"items"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
}

// resolveTimeout returns (connect, read) timeouts from env or defaults.
func resolveTimeout() (float64, float64) {
	readFloat := func(key string, def float64) float64 {
		raw := os.Getenv(key)
		if raw == "" {
			return def
		}
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil || v <= 0 {
			return def
		}
		return v
	}
	return readFloat("MCP_REGISTRY_CONNECT_TIMEOUT", defaultConnectTimeout),
		readFloat("MCP_REGISTRY_READ_TIMEOUT", defaultReadTimeout)
}

// SimpleRegistryClient is a lightweight HTTP client for MCP server discovery.
type SimpleRegistryClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewSimpleRegistryClient creates a registry client targeting the given URL.
// Passing an empty string uses MCP_REGISTRY_URL env var or the default public registry.
func NewSimpleRegistryClient(registryURL string) (*SimpleRegistryClient, error) {
	envOverride := strings.TrimSpace(os.Getenv("MCP_REGISTRY_URL"))
	resolved := registryURL
	if resolved == "" {
		resolved = envOverride
	}
	if resolved == "" {
		resolved = defaultRegistryURL
	}
	resolved = strings.TrimRight(strings.TrimSpace(resolved), "/")

	parsed, err := url.Parse(resolved)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("registry URL %q is not a valid absolute URL", resolved)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, fmt.Errorf("registry URL scheme %q is not supported; use https", parsed.Scheme)
	}
	if parsed.Scheme == "http" && os.Getenv("MCP_REGISTRY_ALLOW_HTTP") != "1" {
		return nil, fmt.Errorf("http:// registry URL rejected; set MCP_REGISTRY_ALLOW_HTTP=1 to allow")
	}

	_, readTO := resolveTimeout()
	return &SimpleRegistryClient{
		baseURL: resolved,
		httpClient: &http.Client{
			Timeout: time.Duration(readTO * float64(time.Second)),
		},
	}, nil
}

// get performs an authenticated GET to path and decodes the JSON response.
func (c *SimpleRegistryClient) get(path string, out interface{}) error {
	u := c.baseURL + path
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s: %w", u, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d from %s: %s", resp.StatusCode, u, truncate(string(body), 200))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode JSON from %s: %w", u, err)
	}
	return nil
}

// SearchServers searches the registry for servers matching query.
func (c *SimpleRegistryClient) SearchServers(query string, page, perPage int) (*SearchResult, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	path := fmt.Sprintf("/v0/servers?q=%s&page=%d&per_page=%d",
		url.QueryEscape(query), page, perPage)
	var result SearchResult
	if err := c.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServer retrieves a single server by its ID or qualified name.
func (c *SimpleRegistryClient) GetServer(serverID string) (*MCPServerInfo, error) {
	path := "/v0/servers/" + url.PathEscape(serverID)
	var info MCPServerInfo
	if err := c.get(path, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetServerVersions returns the available versions for a server.
func (c *SimpleRegistryClient) GetServerVersions(serverID string) ([]VersionEntry, error) {
	path := "/v0/servers/" + url.PathEscape(serverID) + "/versions"
	var versions []VersionEntry
	if err := c.get(path, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// ListServers retrieves a page of servers from the registry index.
func (c *SimpleRegistryClient) ListServers(page, perPage int) (*SearchResult, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	path := fmt.Sprintf("/v0/servers?page=%d&per_page=%d", page, perPage)
	var result SearchResult
	if err := c.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BaseURL returns the base URL this client is targeting.
func (c *SimpleRegistryClient) BaseURL() string { return c.baseURL }

// truncate caps s to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
