// Package mcpentry builds MCP apm.yml entries from CLI parameters.
// Mirrors src/apm_cli/install/mcp/entry.py.
package mcpentry

// EntryKind distinguishes how the MCP entry was constructed.
type EntryKind int

const (
	EntryKindRegistryShorthand EntryKind = iota
	EntryKindRegistryDict
	EntryKindSelfDefinedStdio
	EntryKindSelfDefinedRemote
)

// MCPEntry represents an MCP dependency entry as it will appear in apm.yml.
// A nil map for Env/Headers means the field is absent.
type MCPEntry struct {
	// Name is the MCP server name.
	Name string
	// Kind indicates which routing path was taken.
	Kind EntryKind
	// Registry is false (bool) for self-defined, a URL string for custom
	// registries, and true (bool) for bare registry shorthand.
	Registry interface{}
	// Transport is the chosen transport ("stdio", "http", "sse", etc.).
	Transport string
	// URL is the remote endpoint URL (remote entries only).
	URL string
	// Command is the stdio executable (stdio entries only).
	Command string
	// Args are the extra argv for stdio servers.
	Args []string
	// Env maps environment variable names to values (stdio entries).
	Env map[string]string
	// Headers maps HTTP header names to values (remote entries).
	Headers map[string]string
	// Version is the optional semver constraint (registry entries).
	Version string
}

// IsSelfDefined returns true when the entry represents a self-defined MCP
// (i.e. not resolved from a registry).
func (e MCPEntry) IsSelfDefined() bool {
	return e.Kind == EntryKindSelfDefinedStdio || e.Kind == EntryKindSelfDefinedRemote
}

// BuildMCPEntry constructs an MCPEntry from the CLI inputs, mirroring the
// routing logic in the Python build_mcp_entry function.
// Returns (entry, isSelfDefined).
func BuildMCPEntry(
	name string,
	transport string,
	rawURL string,
	env map[string]string,
	headers map[string]string,
	version string,
	commandArgv []string,
	registryURL string,
) (MCPEntry, bool) {
	if len(commandArgv) > 0 {
		// Self-defined stdio
		e := MCPEntry{
			Name:      name,
			Kind:      EntryKindSelfDefinedStdio,
			Registry:  false,
			Transport: "stdio",
			Command:   commandArgv[0],
		}
		if len(commandArgv) > 1 {
			e.Args = commandArgv[1:]
		}
		if len(env) > 0 {
			e.Env = copyStringMap(env)
		}
		return e, true
	}

	if rawURL != "" {
		// Self-defined remote
		chosen := transport
		if chosen == "" {
			chosen = "http"
		}
		e := MCPEntry{
			Name:      name,
			Kind:      EntryKindSelfDefinedRemote,
			Registry:  false,
			Transport: chosen,
			URL:       rawURL,
		}
		if len(headers) > 0 {
			e.Headers = copyStringMap(headers)
		}
		return e, true
	}

	// Registry entry
	if version != "" || transport != "" || registryURL != "" {
		e := MCPEntry{
			Name:      name,
			Kind:      EntryKindRegistryDict,
			Transport: transport,
			Version:   version,
		}
		if registryURL != "" {
			e.Registry = registryURL
		} else {
			e.Registry = true
		}
		return e, false
	}

	// Bare registry shorthand
	return MCPEntry{
		Name:     name,
		Kind:     EntryKindRegistryShorthand,
		Registry: true,
	}, false
}

func copyStringMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
