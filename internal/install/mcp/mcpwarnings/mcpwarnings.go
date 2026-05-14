// Package mcpwarnings provides MCP install-time non-blocking safety warnings.
// F5 (SSRF) and F7 (shell metacharacters) -- mirroring
// src/apm_cli/install/mcp/warnings.py.
package mcpwarnings

import (
	"net"
	"net/url"
	"strings"
)

// shellMetacharTokens are the shell constructs that would be evaluated by a
// real shell but are NOT evaluated when an MCP stdio server runs via execve.
var shellMetacharTokens = []string{"$(", "`", ";", "&&", "||", "|", ">>", ">", "<"}

// metadataHosts are well-known cloud IMDS addresses.
var metadataHosts = map[string]bool{
	"169.254.169.254": true, // AWS / Azure / GCP
	"100.100.100.200": true, // Alibaba Cloud
	"fd00:ec2::254":   true, // AWS IPv6
}

// IsInternalOrMetadataHost returns true when host resolves or parses to an
// internal IP (loopback, link-local, RFC1918) or a cloud metadata endpoint.
func IsInternalOrMetadataHost(host string) bool {
	if host == "" {
		return false
	}
	bare := strings.Trim(host, "[]")
	if metadataHosts[bare] || metadataHosts[host] {
		return true
	}
	candidates := []string{bare}
	if bare != host {
		candidates = append(candidates, host)
	}
	// Attempt DNS resolution for non-literal hostnames.
	if net.ParseIP(bare) == nil {
		addrs, err := net.LookupHost(bare)
		if err == nil {
			candidates = append(candidates, addrs...)
		}
	}
	for _, c := range candidates {
		ip := net.ParseIP(c)
		if ip == nil {
			continue
		}
		if metadataHosts[ip.String()] {
			return true
		}
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsPrivate() {
			return true
		}
	}
	return false
}

// WarnSSRFURL returns a non-empty warning string when the URL points at an
// internal or cloud metadata address. Returns "" when safe.
func WarnSSRFURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if IsInternalOrMetadataHost(host) {
		return "URL '" + rawURL + "' points to an internal or metadata address; verify intent before installing."
	}
	return ""
}

// WarnShellMetachars returns warning strings for any shell metacharacter
// found in env values or the stdio command field.
func WarnShellMetachars(env map[string]string, command string) []string {
	var warnings []string
	for key, value := range env {
		sval := value
		for _, tok := range shellMetacharTokens {
			if strings.Contains(sval, tok) {
				warnings = append(warnings, "Env value for '"+key+"' contains shell metacharacter '"+tok+"'; reminder these are NOT shell-evaluated.")
				break
			}
		}
	}
	if command != "" {
		for _, tok := range shellMetacharTokens {
			if strings.Contains(command, tok) {
				warnings = append(warnings, "'command' contains shell metacharacter '"+tok+"'; reminder MCP stdio servers run via execve (no shell). This will be passed literally.")
				break
			}
		}
	}
	return warnings
}
