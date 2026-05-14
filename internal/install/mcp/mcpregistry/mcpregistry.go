// Package mcpregistry validates and resolves MCP registry URLs.
// Mirrors src/apm_cli/install/mcp/registry.py.
package mcpregistry

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const maxRegistryURLLength = 2048

// AllowedSchemes are the URL schemes accepted for registry URLs.
var AllowedSchemes = map[string]bool{
	"https": true,
	"http":  true,
}

// ValidationError is returned for invalid registry URLs.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// RedactURLCredentials strips user:password@ from a URL before logging.
func RedactURLCredentials(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	if u.User == nil {
		return rawURL
	}
	// Rebuild without userinfo.
	clean := *u
	clean.User = nil
	return clean.String()
}

// isLocalOrMetadataHost returns true for loopback, link-local, RFC1918, or
// cloud metadata hosts.
func isLocalOrMetadataHost(host string) bool {
	if host == "" {
		return false
	}
	lower := strings.ToLower(host)
	if lower == "localhost" || lower == "ip6-localhost" || lower == "ip6-loopback" {
		return true
	}
	// Try as IP address.
	ip := net.ParseIP(lower)
	if ip == nil {
		// Try as decimal integer (obfuscated form like 2130706433 == 127.0.0.1).
		if n, err := strconv.ParseInt(lower, 10, 64); err == nil {
			b := [4]byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
			ip = net.IP(b[:])
		}
	}
	if ip == nil {
		return false
	}
	cloudMetadata := map[string]bool{
		"169.254.169.254": true,
		"100.100.100.200": true,
	}
	if cloudMetadata[ip.String()] {
		return true
	}
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsPrivate()
}

// ValidateRegistryURL validates the --registry URL value.
// Returns (normalizedURL, localWarning, error).
// localWarning is non-empty for local/metadata hosts (soft warning only).
func ValidateRegistryURL(rawURL string) (string, string, error) {
	if len(rawURL) > maxRegistryURLLength {
		return "", "", &ValidationError{
			Message: fmt.Sprintf("--registry URL too long (%d > %d chars)", len(rawURL), maxRegistryURLLength),
		}
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", &ValidationError{Message: fmt.Sprintf("invalid --registry URL: %v", err)}
	}
	scheme := strings.ToLower(u.Scheme)
	if !AllowedSchemes[scheme] {
		return "", "", &ValidationError{
			Message: fmt.Sprintf("--registry URL scheme %q is not allowed; use https:// (or http:// for local mirrors)", scheme),
		}
	}
	if u.Host == "" {
		return "", "", &ValidationError{Message: "--registry URL must have a host"}
	}
	normalized := u.String()
	var localWarn string
	if isLocalOrMetadataHost(u.Hostname()) {
		localWarn = fmt.Sprintf("--registry URL '%s' points to a local or metadata host; verify intent.", RedactURLCredentials(rawURL))
	}
	return normalized, localWarn, nil
}

// ResolveRegistryURL determines the effective registry URL from the CLI flag
// and the MCP_REGISTRY_URL environment variable. The CLI flag takes precedence.
func ResolveRegistryURL(flagValue, envValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return envValue
}

// RegistryEnvOverride returns the environment additions needed to expose the
// registry URL to the MCPIntegrator subprocess.
// Returns (envKey->value map, allowHTTP bool).
func RegistryEnvOverride(registryURL string) (map[string]string, bool) {
	if registryURL == "" {
		return nil, false
	}
	env := map[string]string{
		"MCP_REGISTRY_URL": registryURL,
	}
	u, err := url.Parse(registryURL)
	allowHTTP := err == nil && strings.ToLower(u.Scheme) == "http"
	if allowHTTP {
		env["MCP_REGISTRY_ALLOW_HTTP"] = "1"
	}
	return env, allowHTTP
}
