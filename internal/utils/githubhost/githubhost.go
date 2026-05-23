// Package githubhost provides host classification and URL utilities for
// GitHub, GitHub Enterprise, Azure DevOps, and GitLab hostnames.
// Mirrors src/apm_cli/utils/github_host.py.
package githubhost

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// fqdnPattern validates a Fully Qualified Domain Name per the Python implementation.
var fqdnPattern = regexp.MustCompile(
	`^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?)+$`,
)

// IsValidFQDN reports whether hostname is a valid FQDN (at least two labels).
func IsValidFQDN(hostname string) bool {
	if hostname == "" {
		return false
	}
	// Strip path components.
	if idx := strings.Index(hostname, "/"); idx >= 0 {
		hostname = hostname[:idx]
	}
	return fqdnPattern.MatchString(hostname)
}

// DefaultHost returns the configured default git host (GITHUB_HOST env or "github.com").
func DefaultHost() string {
	if h := os.Getenv("GITHUB_HOST"); h != "" {
		return h
	}
	return "github.com"
}

// IsAzureDevOpsHostname reports whether hostname is an ADO host
// (dev.azure.com or *.visualstudio.com).
func IsAzureDevOpsHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := strings.ToLower(hostname)
	return h == "dev.azure.com" || strings.HasSuffix(h, ".visualstudio.com")
}

// IsGitHubHostname reports whether hostname is GitHub SaaS or GHE Cloud (*.ghe.com).
func IsGitHubHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := strings.ToLower(hostname)
	return h == "github.com" || strings.HasSuffix(h, ".ghe.com")
}

// IsGitLabHostname reports whether hostname is GitLab SaaS or a configured
// self-managed GitLab host (GITLAB_HOST or APM_GITLAB_HOSTS env vars).
// GHES host takes precedence -- if GITHUB_HOST matches, this returns false.
func IsGitLabHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := strings.ToLower(strings.SplitN(hostname, "/", 2)[0])

	// GHES precedence check.
	ghesHost := strings.ToLower(strings.SplitN(os.Getenv("GITHUB_HOST"), "/", 2)[0])
	if ghesHost != "" && ghesHost == h &&
		ghesHost != "github.com" && ghesHost != "gitlab.com" &&
		!strings.HasSuffix(ghesHost, ".ghe.com") &&
		IsValidFQDN(ghesHost) {
		return false
	}

	if h == "gitlab.com" {
		return true
	}
	if single := strings.ToLower(strings.SplitN(os.Getenv("GITLAB_HOST"), "/", 2)[0]); single != "" && single == h {
		return IsValidFQDN(h)
	}
	for _, part := range strings.Split(os.Getenv("APM_GITLAB_HOSTS"), ",") {
		entry := strings.ToLower(strings.SplitN(strings.TrimSpace(part), "/", 2)[0])
		if entry != "" && entry == h && IsValidFQDN(entry) {
			return true
		}
	}
	return false
}

// SupportGHCLIHost reports whether host should use gh CLI token fallback.
func SupportGHCLIHost(host string) bool {
	if host == "" {
		return false
	}
	if IsGitHubHostname(host) {
		return true
	}
	configured := strings.ToLower(DefaultHost())
	hostLower := strings.ToLower(host)
	if hostLower != configured {
		return false
	}
	if configured == "github.com" || strings.HasSuffix(configured, ".ghe.com") {
		return false
	}
	if IsAzureDevOpsHostname(configured) {
		return false
	}
	return IsValidFQDN(configured)
}

// adoAuthFailureSignals are the case-insensitive signals for ADO auth failures.
var adoAuthFailureSignals = []string{
	"401",
	"403",
	"authentication failed",
	"unauthorized",
	"could not read username",
}

// IsADOAuthFailureSignal reports whether text contains an ADO auth-failure signal.
func IsADOAuthFailureSignal(text string) bool {
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)
	for _, sig := range adoAuthFailureSignals {
		if strings.Contains(lower, sig) {
			return true
		}
	}
	return false
}

// BuildAuthorizationHeaderGitEnv builds env vars to inject an HTTP Authorization
// header into git operations via GIT_CONFIG_COUNT/KEY_N/VALUE_N.
func BuildAuthorizationHeaderGitEnv(scheme, credential string) map[string]string {
	return map[string]string{
		"GIT_CONFIG_COUNT": "1",
		"GIT_CONFIG_KEY_0": "http.extraheader",
		"GIT_CONFIG_VALUE_0": fmt.Sprintf("Authorization: %s %s", scheme, credential),
	}
}

// BuildADOBearerGitEnv builds env vars to authenticate to Azure DevOps
// with an Entra ID bearer token.
func BuildADOBearerGitEnv(bearerToken string) map[string]string {
	return BuildAuthorizationHeaderGitEnv("Bearer", bearerToken)
}
