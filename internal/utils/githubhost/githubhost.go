// Package githubhost provides utilities for handling GitHub, GitHub Enterprise,
// Azure DevOps, and other Git host hostnames and URLs.
package githubhost

import (
	"os"
	"regexp"
	"strings"
)

// validFQDNRe matches a valid fully-qualified domain name.
var validFQDNRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)

// DefaultHost returns the default Git host (can be overridden via GITHUB_HOST env var).
func DefaultHost() string {
	if h := os.Getenv("GITHUB_HOST"); h != "" {
		return h
	}
	return "github.com"
}

// IsAzureDevOpsHostname returns true if hostname is Azure DevOps (cloud or server).
func IsAzureDevOpsHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := strings.ToLower(hostname)
	return h == "dev.azure.com" || strings.HasSuffix(h, ".visualstudio.com")
}

// IsVisualStudioLegacyHostname returns true if hostname is a legacy *.visualstudio.com ADO host.
func IsVisualStudioLegacyHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	return strings.HasSuffix(strings.ToLower(hostname), ".visualstudio.com")
}

// IsGitLabHostname returns true if hostname is GitLab SaaS or a configured GitLab host.
func IsGitLabHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := normalizeHost(hostname)

	// GHES precedence: GITHUB_HOST match is enterprise GitHub, not GitLab
	ghesHost := normalizeHost(os.Getenv("GITHUB_HOST"))
	if ghesHost != "" && ghesHost == h &&
		ghesHost != "github.com" && ghesHost != "gitlab.com" &&
		!strings.HasSuffix(ghesHost, ".ghe.com") &&
		IsValidFQDN(ghesHost) {
		return false
	}

	if h == "gitlab.com" {
		return true
	}
	gitlabSingle := normalizeHost(os.Getenv("GITLAB_HOST"))
	if gitlabSingle != "" && gitlabSingle == h {
		return IsValidFQDN(h)
	}
	rawList := os.Getenv("APM_GITLAB_HOSTS")
	for _, part := range strings.Split(rawList, ",") {
		entry := normalizeHost(part)
		if entry != "" && entry == h && IsValidFQDN(entry) {
			return true
		}
	}
	return false
}

// HasGitHubGitLabHostEnvConflict returns true when hostname is claimed as both GHES and GitLab.
func HasGitHubGitLabHostEnvConflict(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := normalizeHost(hostname)
	if !IsValidFQDN(h) {
		return false
	}
	ghesHost := normalizeHost(os.Getenv("GITHUB_HOST"))
	if ghesHost == "" || ghesHost != h ||
		ghesHost == "github.com" || ghesHost == "gitlab.com" ||
		strings.HasSuffix(ghesHost, ".ghe.com") {
		return false
	}
	gitlabSingle := normalizeHost(os.Getenv("GITLAB_HOST"))
	if gitlabSingle != "" && gitlabSingle == h {
		return true
	}
	rawList := os.Getenv("APM_GITLAB_HOSTS")
	for _, part := range strings.Split(rawList, ",") {
		if normalizeHost(part) == h {
			return true
		}
	}
	return false
}

// IsGHEHostname returns true if hostname is GitHub Enterprise Server or GHE.com.
func IsGHEHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := normalizeHost(hostname)
	if h == "github.com" {
		return false
	}
	if strings.HasSuffix(h, ".ghe.com") {
		return true
	}
	ghesHost := normalizeHost(os.Getenv("GITHUB_HOST"))
	return ghesHost != "" && ghesHost == h && IsValidFQDN(h)
}

// IsGitHubHostname returns true if hostname is github.com or a GHES instance.
func IsGitHubHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := normalizeHost(hostname)
	return h == "github.com" || IsGHEHostname(h)
}

// IsArtifactoryHostname returns true if hostname is an Artifactory instance.
func IsArtifactoryHostname(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := normalizeHost(hostname)
	// Check APM_ARTIFACTORY_HOSTS env
	rawList := os.Getenv("APM_ARTIFACTORY_HOSTS")
	for _, part := range strings.Split(rawList, ",") {
		entry := normalizeHost(part)
		if entry != "" && entry == h {
			return true
		}
	}
	return false
}

// ClassifyHost returns the host type: "github", "ghes", "ghe_com", "gitlab",
// "azure_devops", "artifactory", or "unknown".
func ClassifyHost(hostname string) string {
	if hostname == "" {
		return "unknown"
	}
	h := normalizeHost(hostname)
	if h == "github.com" {
		return "github"
	}
	if strings.HasSuffix(h, ".ghe.com") {
		return "ghe_com"
	}
	ghesHost := normalizeHost(os.Getenv("GITHUB_HOST"))
	if ghesHost != "" && ghesHost == h && ghesHost != "github.com" && IsValidFQDN(h) {
		return "ghes"
	}
	if IsAzureDevOpsHostname(h) {
		return "azure_devops"
	}
	if IsGitLabHostname(h) {
		return "gitlab"
	}
	if IsArtifactoryHostname(h) {
		return "artifactory"
	}
	return "unknown"
}

// IsValidFQDN returns true if hostname is a syntactically valid fully-qualified domain name.
func IsValidFQDN(hostname string) bool {
	if hostname == "" || len(hostname) > 253 {
		return false
	}
	return validFQDNRe.MatchString(hostname)
}

// ParseHostFromURL extracts the hostname from a URL string.
func ParseHostFromURL(rawURL string) string {
	// Strip scheme
	s := rawURL
	if idx := strings.Index(s, "://"); idx >= 0 {
		s = s[idx+3:]
	}
	// Strip path
	if idx := strings.Index(s, "/"); idx >= 0 {
		s = s[:idx]
	}
	// Strip port
	if idx := strings.LastIndex(s, ":"); idx >= 0 {
		s = s[:idx]
	}
	// Strip user info
	if idx := strings.Index(s, "@"); idx >= 0 {
		s = s[idx+1:]
	}
	return strings.ToLower(strings.TrimSpace(s))
}

// AzureDevOpsOrgFromHostname extracts the org name from a legacy *.visualstudio.com host.
func AzureDevOpsOrgFromHostname(hostname string) string {
	h := strings.ToLower(hostname)
	if !strings.HasSuffix(h, ".visualstudio.com") {
		return ""
	}
	parts := strings.SplitN(h, ".", 2)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

// IsSupportedGitHost returns true for any hostname that APM recognises as a valid
// Git host: github.com, GHES, GHE.com, GitLab, Azure DevOps, or Artifactory.
// Any syntactically valid FQDN is accepted (self-hosted instances).
func IsSupportedGitHost(hostname string) bool {
	if hostname == "" {
		return false
	}
	h := normalizeHost(hostname)
	return IsValidFQDN(h)
}

// IsArtifactoryPath returns true when path segments start with "artifactory/".
func IsArtifactoryPath(segments []string) bool {
	return len(segments) >= 4 && strings.EqualFold(segments[0], "artifactory")
}

// ParseArtifactoryPath extracts (prefix, owner, repo) from Artifactory path segments.
// Segments are expected as ["artifactory", "<repo-key>", "<owner>", "<repo>", ...].
// Returns empty strings if the segments do not match.
func ParseArtifactoryPath(segments []string) (prefix, owner, repo string) {
	if !IsArtifactoryPath(segments) {
		return
	}
	prefix = strings.Join(segments[:2], "/")
	owner = segments[2]
	repo = segments[3]
	return
}

func normalizeHost(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	// Strip path component
	if idx := strings.Index(s, "/"); idx >= 0 {
		s = s[:idx]
	}
	return s
}
