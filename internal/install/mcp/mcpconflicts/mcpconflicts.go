// Package mcpconflicts validates MCP CLI flag-conflict matrix (E1-E15).
// Mirrors src/apm_cli/install/mcp/conflicts.py.
package mcpconflicts

import "fmt"

// ValidationError is returned when a flag conflict is detected.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

func conflict(msg string) *ValidationError { return &ValidationError{Message: msg} }

// ConflictConfig holds all the flag values passed to ValidateMCPConflicts.
type ConflictConfig struct {
	MCPName              string
	HasMCPName           bool
	Packages             []string
	PreDashPackages      []string
	Transport            string
	URL                  string
	Env                  map[string]string
	Headers              map[string]string
	MCPVersion           string
	CommandArgv          []string
	Global               bool
	Only                 string
	Update               bool
	UseSSH               bool
	UseHTTPS             bool
	AllowProtocolFallback bool
	RegistryURL          string
}

// ValidateMCPConflicts applies the E1-E15 conflict matrix.
// Returns nil on success or a *ValidationError on conflict.
func ValidateMCPConflicts(cfg ConflictConfig) error {
	// E10: flags require --mcp
	if !cfg.HasMCPName {
		requiresMCPFlags := []struct {
			Value interface{}
			Label string
		}{
			{cfg.Transport, "--transport"},
			{cfg.URL, "--url"},
			{cfg.Env, "--env"},
			{cfg.Headers, "--header"},
			{cfg.MCPVersion, "--mcp-version"},
			{cfg.RegistryURL, "--registry"},
		}
		for _, f := range requiresMCPFlags {
			switch v := f.Value.(type) {
			case string:
				if v != "" {
					return conflict(fmt.Sprintf("%s requires --mcp", f.Label))
				}
			case map[string]string:
				if len(v) > 0 {
					return conflict(fmt.Sprintf("%s requires --mcp", f.Label))
				}
			}
		}
		return nil
	}

	// E7/E8: NAME shape
	if cfg.MCPName == "" {
		return conflict("MCP name cannot be empty")
	}
	if len(cfg.MCPName) > 0 && cfg.MCPName[0] == '-' {
		return conflict("MCP name cannot start with '-'; did you forget a value for --mcp?")
	}

	// E1: positional packages mixed with --mcp
	if len(cfg.PreDashPackages) > 0 {
		return conflict("cannot mix --mcp with positional packages")
	}

	// E2: --global not supported for MCP
	if cfg.Global {
		return conflict("MCP servers are project-scoped; --global is not supported for MCP entries")
	}

	// E3: --only apm conflicts with --mcp
	if cfg.Only == "apm" {
		return conflict("cannot use --only apm with --mcp")
	}

	// E4: transport selection flags
	if cfg.UseSSH || cfg.UseHTTPS || cfg.AllowProtocolFallback {
		return conflict("transport selection flags (--ssh/--https/--allow-protocol-fallback) don't apply to MCP entries")
	}

	// E5: --update
	if cfg.Update {
		return conflict("use 'apm update' instead to update MCP entries")
	}

	// E9: --header without --url
	if len(cfg.Headers) > 0 && cfg.URL == "" {
		return conflict("--header requires --url")
	}

	// E11: --url with stdio command
	if cfg.URL != "" && len(cfg.CommandArgv) > 0 {
		return conflict("cannot specify both --url and a stdio command")
	}

	// E12: --transport stdio with --url
	if cfg.Transport == "stdio" && cfg.URL != "" {
		return conflict("stdio transport doesn't accept --url")
	}

	// E13: remote transports with stdio command
	switch cfg.Transport {
	case "http", "sse", "streamable-http":
		if len(cfg.CommandArgv) > 0 {
			return conflict("remote transports don't accept stdio command")
		}
	}

	// E14: --env with --url and no command
	if len(cfg.Env) > 0 && cfg.URL != "" && len(cfg.CommandArgv) == 0 {
		return conflict("--env applies to stdio MCPs; use --header for remote")
	}

	// E15: --registry only applies to registry-resolved entries
	if cfg.RegistryURL != "" && (cfg.URL != "" || len(cfg.CommandArgv) > 0) {
		return conflict("--registry only applies to registry-resolved MCP servers; remove --url or the post-`--` stdio command, or drop --registry")
	}

	return nil
}
