// Package registry provides types for MCP server registry interactions.
// Mirrors Python apm_cli.registry.client and apm_cli.registry.operations.
package registry

import "fmt"

// DefaultRegistryURL is the default MCP registry API base URL.
const DefaultRegistryURL = "https://api.mcp.github.com"

// V01Prefix is the API version path prefix for MCP Registry v0.1.
const V01Prefix = "/v0.1"

// DefaultConnectTimeout is the default TCP connect timeout in seconds.
const DefaultConnectTimeout = 10.0

// DefaultReadTimeout is the default HTTP read timeout in seconds.
const DefaultReadTimeout = 30.0

// ServerNotFoundError is raised when a registry lookup returns 404.
// Mirrors Python ServerNotFoundError.
type ServerNotFoundError struct {
	ServerName  string
	RegistryURL string
}

func (e *ServerNotFoundError) Error() string {
	return fmt.Sprintf(
		"Server '%s' not found in registry %s. "+
			"If this is a self-hosted registry, verify it implements the "+
			"MCP Registry v0.1 API (apm uses /v0.1/servers/...).",
		e.ServerName, e.RegistryURL,
	)
}

// RegistryError wraps generic registry HTTP errors.
type RegistryError struct {
	StatusCode int
	Message    string
}

func (e *RegistryError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("registry error %d: %s", e.StatusCode, e.Message)
	}
	return e.Message
}

// ServerEntry represents a single MCP server entry returned from the registry.
// Mirrors Python MCPServerEntry in registry client.
type ServerEntry struct {
	Name        string
	Description string
	Source      string // e.g. "github", "url"
	Repository  string // owner/repo for github sources
	Version     string
	Tags        []string
}

// SearchResult wraps a list of server entries from a registry search.
type SearchResult struct {
	Servers []ServerEntry
	Total   int
}

// InstallStatus describes the current install state of an MCP server.
// Mirrors Python MCPServerOperations install status checks.
type InstallStatus int

const (
	// StatusNotInstalled means the server is not installed.
	StatusNotInstalled InstallStatus = iota
	// StatusInstalled means the server is installed and up to date.
	StatusInstalled
	// StatusConflict means multiple integrations define the same server.
	StatusConflict
	// StatusOutdated means a newer version is available.
	StatusOutdated
)

// String returns a human-readable install status label.
func (s InstallStatus) String() string {
	switch s {
	case StatusInstalled:
		return "installed"
	case StatusConflict:
		return "conflict"
	case StatusOutdated:
		return "outdated"
	default:
		return "not-installed"
	}
}

// ConflictEntry records a detected server name conflict across integrations.
type ConflictEntry struct {
	ServerName   string
	Integrations []string
}

// ServerReference holds a parsed server reference from user input.
// Format: [registry/]owner/repo[@version] or server-name.
type ServerReference struct {
	Raw         string
	Owner       string
	Repo        string
	Version     string
	RegistryURL string
}

// IsVersioned returns true if a version constraint was specified.
func (r ServerReference) IsVersioned() bool { return r.Version != "" }

// ParseServerReference parses a user-supplied server reference string.
// Mirrors Python ref_resolver.parse_server_reference().
// Supports formats: name, owner/repo, owner/repo@version.
func ParseServerReference(raw string) ServerReference {
	ref := ServerReference{Raw: raw}
	// strip optional @version suffix
	s := raw
	if idx := lastIndex(s, "@"); idx >= 0 {
		ref.Version = s[idx+1:]
		s = s[:idx]
	}
	if idx := indexOf(s, "/"); idx >= 0 {
		ref.Owner = s[:idx]
		ref.Repo = s[idx+1:]
	} else {
		ref.Repo = s
	}
	return ref
}

// SemVer holds a parsed semantic version triple.
type SemVer struct {
	Major int
	Minor int
	Patch int
	Pre   string
}

// String formats the semantic version.
func (v SemVer) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	return s
}

// Compare returns -1, 0, or 1 for less-than, equal, greater-than.
// Mirrors Python semver comparison. Pre-release versions sort lower.
func (v SemVer) Compare(other SemVer) int {
	if v.Major != other.Major {
		return cmpInt(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return cmpInt(v.Minor, other.Minor)
	}
	if v.Patch != other.Patch {
		return cmpInt(v.Patch, other.Patch)
	}
	// pre-release is lower than release
	if v.Pre == "" && other.Pre != "" {
		return 1
	}
	if v.Pre != "" && other.Pre == "" {
		return -1
	}
	if v.Pre < other.Pre {
		return -1
	}
	if v.Pre > other.Pre {
		return 1
	}
	return 0
}

func cmpInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func lastIndex(s, sub string) int {
	last := -1
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			last = i
		}
	}
	return last
}
