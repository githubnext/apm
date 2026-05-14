// Package mcpcommand orchestrates the apm install --mcp code path,
// composing the sibling MCP modules into the user-visible install flow.
// Mirrors src/apm_cli/install/mcp/command.py.
package mcpcommand

import (
	"strings"
)

// EnvPair parses a "KEY=VALUE" string into (key, value).
// Returns empty strings if the format is invalid.
func ParseEnvPair(pair string) (string, string, bool) {
	idx := strings.Index(pair, "=")
	if idx < 0 {
		return "", "", false
	}
	return pair[:idx], pair[idx+1:], true
}

// ParseEnvPairs converts a slice of "KEY=VALUE" strings to a map.
// Invalid pairs are skipped.
func ParseEnvPairs(pairs []string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := ParseEnvPair(p)
		if ok {
			out[k] = v
		}
	}
	return out
}

// ParseHeaderPair parses a "Name: Value" or "Name=Value" header string.
func ParseHeaderPair(pair string) (string, string, bool) {
	if idx := strings.Index(pair, ": "); idx >= 0 {
		return strings.TrimSpace(pair[:idx]), strings.TrimSpace(pair[idx+2:]), true
	}
	if idx := strings.Index(pair, "="); idx >= 0 {
		return strings.TrimSpace(pair[:idx]), strings.TrimSpace(pair[idx+1:]), true
	}
	return "", "", false
}

// ParseHeaderPairs converts a slice of header strings to a map.
func ParseHeaderPairs(pairs []string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := ParseHeaderPair(p)
		if ok {
			out[k] = v
		}
	}
	return out
}

// MCPInstallRequest holds all the parameters for the --mcp install path.
type MCPInstallRequest struct {
	MCPName     string
	Transport   string
	URL         string
	EnvPairs    []string
	HeaderPairs []string
	MCPVersion  string
	CommandArgv []string
	Dev         bool
	Force       bool
	Runtime     string
	Exclude     string
	Verbose     bool
	RegistryURL string
	Scope       string
}

// MCPInstallResult summarises what the --mcp install path did.
type MCPInstallResult struct {
	Outcome    string // "added", "replaced", "skipped"
	EntryKey   string
	Integrated bool
}

// TransportDefault returns the default transport for the given inputs,
// mirroring the Python entry builder routing logic.
func TransportDefault(url string, commandArgv []string, transport string) string {
	if transport != "" {
		return transport
	}
	if len(commandArgv) > 0 {
		return "stdio"
	}
	if url != "" {
		return "http"
	}
	return ""
}
