// Package mcpdep implements the MCP dependency model.
// Ported from src/apm_cli/models/dependency/mcp.py
package mcpdep

import (
	"fmt"
	"net/url"
	"strings"
)

var validNameChars = func() [256]bool {
	var t [256]bool
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._@/:=-" {
		t[c] = true
	}
	return t
}()

var validNameStart = func() [256]bool {
	var t [256]bool
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@_" {
		t[c] = true
	}
	return t
}()

var validTransports = map[string]bool{
	"stdio":            true,
	"sse":              true,
	"http":             true,
	"streamable-http":  true,
}

var allowedURLSchemes = map[string]bool{
	"http":  true,
	"https": true,
}

// MCPDependency represents an MCP server dependency with optional overlay configuration.
// Supports three forms: string (registry reference), object with overlays, and self-defined.
type MCPDependency struct {
	Name      string
	Transport string      // "stdio" | "sse" | "streamable-http" | "http"
	Env       map[string]string
	Args      interface{} // map[string]interface{} for overlay, []string for self-defined
	Version   string
	// Registry: nil = default registry, false (RegistryFalse sentinel) = self-defined, string = custom URL
	Registry  interface{}
	Package   string
	Headers   map[string]string
	Tools     []string
	URL       string
	Command   string
}

// RegistryFalse is a sentinel value for Registry = false (self-defined dependency).
const RegistryFalse = registryFalseSentinel(0)

type registryFalseSentinel int

// IsRegistryResolved returns true when the dependency is resolved via a registry.
func (d *MCPDependency) IsRegistryResolved() bool {
	_, isFalse := d.Registry.(registryFalseSentinel)
	return !isFalse
}

// IsSelfDefined returns true when the dependency is self-defined (registry: false).
func (d *MCPDependency) IsSelfDefined() bool {
	_, isFalse := d.Registry.(registryFalseSentinel)
	return isFalse
}

// FromString creates an MCPDependency from a plain string (registry reference).
func FromString(s string) (*MCPDependency, error) {
	d := &MCPDependency{Name: s}
	if err := d.Validate(false); err != nil {
		return nil, err
	}
	return d, nil
}

// FromDict parses an MCPDependency from a map.
func FromDict(m map[string]interface{}) (*MCPDependency, error) {
	name, ok := m["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("MCP dependency dict must contain 'name'")
	}

	transport, _ := m["transport"].(string)
	if transport == "" {
		transport, _ = m["type"].(string) // legacy 'type' -> 'transport'
	}

	env, _ := m["env"].(map[string]interface{})
	var envMap map[string]string
	if env != nil {
		envMap = make(map[string]string, len(env))
		for k, v := range env {
			envMap[k] = fmt.Sprintf("%v", v)
		}
	}

	headers, _ := m["headers"].(map[string]interface{})
	var headersMap map[string]string
	if headers != nil {
		headersMap = make(map[string]string, len(headers))
		for k, v := range headers {
			headersMap[k] = fmt.Sprintf("%v", v)
		}
	}

	var tools []string
	if rawTools, ok := m["tools"].([]interface{}); ok {
		for _, t := range rawTools {
			if s, ok := t.(string); ok {
				tools = append(tools, s)
			}
		}
	}

	version, _ := m["version"].(string)
	pkg, _ := m["package"].(string)
	rawURL, _ := m["url"].(string)
	command, _ := m["command"].(string)

	var registry interface{}
	if regRaw, hasReg := m["registry"]; hasReg {
		if b, ok := regRaw.(bool); ok && !b {
			registry = RegistryFalse
		} else {
			registry = regRaw
		}
	}

	d := &MCPDependency{
		Name:      name,
		Transport: transport,
		Env:       envMap,
		Args:      m["args"],
		Version:   version,
		Registry:  registry,
		Package:   pkg,
		Headers:   headersMap,
		Tools:     tools,
		URL:       rawURL,
		Command:   command,
	}

	strict := d.IsSelfDefined()
	if err := d.Validate(strict); err != nil {
		return nil, err
	}
	return d, nil
}

// ToDict serializes to map, including only non-zero fields.
func (d *MCPDependency) ToDict() map[string]interface{} {
	result := map[string]interface{}{"name": d.Name}
	if d.Transport != "" {
		result["transport"] = d.Transport
	}
	if d.Env != nil {
		result["env"] = d.Env
	}
	if d.Args != nil {
		result["args"] = d.Args
	}
	if d.Version != "" {
		result["version"] = d.Version
	}
	if d.Registry != nil {
		if d.IsSelfDefined() {
			result["registry"] = false
		} else {
			result["registry"] = d.Registry
		}
	}
	if d.Package != "" {
		result["package"] = d.Package
	}
	if d.Headers != nil {
		result["headers"] = d.Headers
	}
	if d.Tools != nil {
		result["tools"] = d.Tools
	}
	if d.URL != "" {
		result["url"] = d.URL
	}
	if d.Command != "" {
		result["command"] = d.Command
	}
	return result
}

// String returns a human-friendly identifier.
func (d *MCPDependency) String() string {
	if d.Transport != "" {
		return fmt.Sprintf("%s (%s)", d.Name, d.Transport)
	}
	return d.Name
}

// Validate validates the dependency. Returns error on invalid state.
func (d *MCPDependency) Validate(strict bool) error {
	if d.Name == "" {
		return fmt.Errorf("MCP dependency 'name' must not be empty")
	}
	if !isValidName(d.Name) {
		return fmt.Errorf(
			"Invalid MCP dependency name %q: must start with a letter, digit, '@', or '_' "+
				"and contain only [a-zA-Z0-9._@/:=-] (max 128 chars). "+
				"Example: 'io.github.acme/cool-server' or 'my-server'.",
			d.Name,
		)
	}
	for _, seg := range strings.Split(d.Name, "/") {
		if seg == ".." {
			return fmt.Errorf(
				"Invalid MCP dependency name %q: must not contain '..' path segments. "+
					"Example: 'io.github.acme/cool-server' or 'my-server'.",
				d.Name,
			)
		}
	}
	if d.URL != "" {
		u, err := url.Parse(d.URL)
		if err != nil || !allowedURLSchemes[strings.ToLower(u.Scheme)] {
			scheme := ""
			if err == nil {
				scheme = strings.ToLower(u.Scheme)
			}
			return fmt.Errorf(
				"Invalid MCP url %q: scheme %q is not supported; use http:// or https://. "+
					"WebSocket URLs (ws/wss) are not supported for MCP transports.",
				d.URL, scheme,
			)
		}
	}
	if d.Headers != nil {
		for k, v := range d.Headers {
			if strings.ContainsAny(k, "\r\n") || strings.ContainsAny(v, "\r\n") {
				return fmt.Errorf(
					"Invalid header '%s=%s': control characters (CR/LF) not allowed in keys or values",
					k, v,
				)
			}
		}
	}
	if d.Command != "" {
		for _, seg := range strings.Split(d.Command, "/") {
			if seg == ".." {
				return fmt.Errorf(
					"Invalid MCP command %q: must not contain '..' path segments. "+
						"Use an absolute path or a command name on PATH instead.",
					d.Command,
				)
			}
		}
	}
	if !strict {
		return nil
	}
	if d.Transport != "" && !validTransports[d.Transport] {
		var sortedTransports []string
		for t := range validTransports {
			sortedTransports = append(sortedTransports, t)
		}
		return fmt.Errorf(
			"MCP dependency %q has unsupported transport %q. Valid values: %s",
			d.Name, d.Transport, strings.Join(sortedTransports, ", "),
		)
	}
	if d.IsSelfDefined() {
		if d.Transport == "" {
			return fmt.Errorf("Self-defined MCP dependency %q requires 'transport'", d.Name)
		}
		if (d.Transport == "http" || d.Transport == "sse" || d.Transport == "streamable-http") && d.URL == "" {
			return fmt.Errorf(
				"Self-defined MCP dependency %q with transport %q requires 'url'",
				d.Name, d.Transport,
			)
		}
		if d.Transport == "stdio" && d.Command == "" {
			return fmt.Errorf(
				"Self-defined MCP dependency %q with transport 'stdio' requires 'command'",
				d.Name,
			)
		}
		if d.Transport == "stdio" && d.Command != "" && d.Args == nil {
			if strings.ContainsAny(d.Command, " \t") {
				parts := strings.Fields(d.Command)
				if len(parts) == 0 {
					return fmt.Errorf(
						"Self-defined MCP dependency %q: 'command' is empty or whitespace-only. "+
							"Set 'command' to a binary path, e.g. command: npx",
						d.Name,
					)
				}
				first := parts[0]
				rest := parts[1:]
				var quotedArgs []string
				for _, tok := range rest {
					quotedArgs = append(quotedArgs, fmt.Sprintf("%q", tok))
				}
				suggestedArgs := "[" + strings.Join(quotedArgs, ", ") + "]"
				return fmt.Errorf(
					"'command' contains whitespace in MCP dependency %q.\n"+
						"  Rule: 'command' must be a single binary path -- APM does not split on whitespace. Use 'args' for additional arguments.\n"+
						"  Got:  command=%q (%d additional args)\n"+
						"  Fix:  command: %s\n"+
						"        args: %s\n"+
						"  See:  https://microsoft.github.io/apm/guides/mcp-servers/",
					d.Name, first, len(rest), first, suggestedArgs,
				)
			}
		}
	}
	return nil
}

func isValidName(name string) bool {
	if len(name) == 0 || len(name) > 128 {
		return false
	}
	if !validNameStart[name[0]] {
		return false
	}
	for i := 1; i < len(name); i++ {
		if !validNameChars[name[i]] {
			return false
		}
	}
	return true
}
