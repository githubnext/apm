// Package hostbackends provides vendor-specific URL/API construction for remote git hosts.
// Migrated from src/apm_cli/deps/host_backends.py.
//
// Each supported host kind is a concrete backend struct implementing the HostBackend interface.
// A dispatch function (BackendFor / BackendForHost) picks the right backend by consulting
// the auth package's ClassifyHost function.
package hostbackends

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/githubnext/apm/internal/core/auth"
	"github.com/githubnext/apm/internal/utils/githubhost"
)

var sha40RE = regexp.MustCompile(`^[a-f0-9]{40}$`)

// DepRef is the minimal interface expected of a dependency reference by backend URL builders.
type DepRef interface {
	// GetHost returns the host string for this dependency (may be "").
	GetHost() string
	// GetPort returns the non-standard port, or nil if default.
	GetPort() *int
	// GetRepoURL returns the "owner/repo" URL string.
	GetRepoURL() string
	// GetADOOrganization returns the ADO organisation name, or "".
	GetADOOrganization() string
	// GetADOProject returns the ADO project name, or "".
	GetADOProject() string
	// GetADORepo returns the ADO repo name, or "".
	GetADORepo() string
	// IsAzureDevOps returns true when this dep references Azure DevOps.
	IsAzureDevOps() bool
	// IsInsecure returns true when the dep was declared with a plain HTTP URL.
	IsInsecure() bool
}

// HostBackend exposes URL/API construction for one remote git host kind.
type HostBackend interface {
	// Kind returns a canonical host-kind string: "github", "ghe_cloud", "ghes", "ado", "gitlab", "generic".
	Kind() string
	// IsGitHubFamily returns true for github.com, *.ghe.com, and configured GHES hosts.
	IsGitHubFamily() bool
	// IsGeneric returns true for non-GitHub-family non-ADO hosts (GitLab, Bitbucket, Gitea, ...).
	IsGeneric() bool
	// GetHostInfo returns the HostInfo for this backend.
	GetHostInfo() auth.HostInfo

	// BuildCloneHTTPSURL builds the HTTPS clone URL.
	// token may be "" (anonymous / bearer), non-empty embeds credentials.
	// authScheme "bearer" suppresses embedding the token in the URL.
	BuildCloneHTTPSURL(dep DepRef, token string, authScheme string) string
	// BuildCloneSSHURL builds the SSH clone URL.
	BuildCloneSSHURL(dep DepRef) string
	// BuildCloneHTTPURL builds a plain HTTP clone URL (only for is_insecure deps).
	BuildCloneHTTPURL(dep DepRef) (string, error)
	// BuildCommitsAPIURL returns the cheap commit-resolution API URL, or "" when unavailable.
	BuildCommitsAPIURL(dep DepRef, ref string) string
	// BuildContentsAPIURLs returns ordered Contents-API URL candidates for fetching a file.
	BuildContentsAPIURLs(owner, repo, filePath, ref string) []string
}

// ---------------------------------------------------------------------------
// URL builder helpers (mirror Python's github_host.py helpers)
// ---------------------------------------------------------------------------

func buildHTTPSCloneURL(host, repoURL, token string, port *int) string {
	if token != "" {
		// embed as https://x-access-token:<token>@host/owner/repo.git
		netloc := netloc(host, port)
		return fmt.Sprintf("https://x-access-token:%s@%s/%s.git", url.PathEscape(token), netloc, repoURL)
	}
	netloc := netloc(host, port)
	return fmt.Sprintf("https://%s/%s.git", netloc, repoURL)
}

func buildSSHURL(host, repoURL string, port *int) string {
	if port != nil {
		return fmt.Sprintf("ssh://git@%s:%d/%s.git", host, *port, repoURL)
	}
	return fmt.Sprintf("git@%s:%s.git", host, repoURL)
}

func buildADOHTTPSCloneURL(org, project, repo, host, token string) string {
	if host == "" {
		host = "dev.azure.com"
	}
	base := fmt.Sprintf("https://%s/%s/%s/_git/%s", host, org, project, repo)
	if token != "" {
		base = fmt.Sprintf("https://%s@%s/%s/%s/_git/%s", token, host, org, project, repo)
	}
	return base
}

func buildADOSSHURL(org, project, repo string) string {
	return fmt.Sprintf("git@ssh.dev.azure.com:v3/%s/%s/%s", org, project, repo)
}

func buildGitLabHTTPSCloneURL(host, repoURL, token string, port *int) string {
	netloc := netloc(host, port)
	if token != "" {
		return fmt.Sprintf("https://oauth2:%s@%s/%s.git", url.PathEscape(token), netloc, repoURL)
	}
	return fmt.Sprintf("https://%s/%s.git", netloc, repoURL)
}

func netloc(host string, port *int) string {
	if port != nil {
		return fmt.Sprintf("%s:%d", host, *port)
	}
	return host
}

func urlHost(dep DepRef, fallback auth.HostInfo) string {
	h := dep.GetHost()
	if h != "" {
		return h
	}
	return fallback.Host
}

// ---------------------------------------------------------------------------
// GitHub-family shared base
// ---------------------------------------------------------------------------

type gitHubFamilyBase struct {
	hostInfo auth.HostInfo
	kind     string
}

func (b *gitHubFamilyBase) Kind() string              { return b.kind }
func (b *gitHubFamilyBase) IsGitHubFamily() bool      { return true }
func (b *gitHubFamilyBase) IsGeneric() bool            { return false }
func (b *gitHubFamilyBase) GetHostInfo() auth.HostInfo { return b.hostInfo }

func (b *gitHubFamilyBase) BuildCloneHTTPSURL(dep DepRef, token string, authScheme string) string {
	host := urlHost(dep, b.hostInfo)
	port := dep.GetPort()
	if authScheme == "bearer" {
		token = ""
	}
	return buildHTTPSCloneURL(host, dep.GetRepoURL(), token, port)
}

func (b *gitHubFamilyBase) BuildCloneSSHURL(dep DepRef) string {
	host := urlHost(dep, b.hostInfo)
	return buildSSHURL(host, dep.GetRepoURL(), dep.GetPort())
}

func (b *gitHubFamilyBase) BuildCloneHTTPURL(dep DepRef) (string, error) {
	host := urlHost(dep, b.hostInfo)
	port := dep.GetPort()
	n := netloc(host, port)
	return fmt.Sprintf("http://%s/%s.git", n, dep.GetRepoURL()), nil
}

func (b *gitHubFamilyBase) BuildCommitsAPIURL(dep DepRef, ref string) string {
	if sha40RE.MatchString(strings.ToLower(ref)) {
		return ""
	}
	parts := strings.SplitN(dep.GetRepoURL(), "/", 2)
	if len(parts) != 2 {
		return ""
	}
	return fmt.Sprintf("%s/repos/%s/%s/commits/%s", b.hostInfo.APIBase, parts[0], parts[1], ref)
}

func (b *gitHubFamilyBase) BuildContentsAPIURLs(owner, repo, filePath, ref string) []string {
	return []string{
		fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", b.hostInfo.APIBase, owner, repo, filePath, ref),
	}
}

// ---------------------------------------------------------------------------
// Concrete backends
// ---------------------------------------------------------------------------

// GitHubBackend is the backend for github.com.
type GitHubBackend struct{ gitHubFamilyBase }

// GHECloudBackend is the backend for *.ghe.com (GitHub Enterprise Cloud -- Data Residency).
type GHECloudBackend struct{ gitHubFamilyBase }

// GHESBackend is the backend for self-hosted GitHub Enterprise Server.
type GHESBackend struct{ gitHubFamilyBase }

// ADOBackend is the backend for Azure DevOps.
type ADOBackend struct {
	hostInfo auth.HostInfo
}

func (b *ADOBackend) Kind() string              { return "ado" }
func (b *ADOBackend) IsGitHubFamily() bool      { return false }
func (b *ADOBackend) IsGeneric() bool            { return false }
func (b *ADOBackend) GetHostInfo() auth.HostInfo { return b.hostInfo }

func (b *ADOBackend) BuildCloneHTTPSURL(dep DepRef, token string, authScheme string) string {
	if dep.GetADOOrganization() == "" {
		// Missing org -- return a diagnostic URL so callers can surface the error.
		return "error://ado-missing-org"
	}
	host := urlHost(dep, b.hostInfo)
	if host == "" {
		host = "dev.azure.com"
	}
	return buildADOHTTPSCloneURL(dep.GetADOOrganization(), dep.GetADOProject(), dep.GetADORepo(), host, token)
}

func (b *ADOBackend) BuildCloneSSHURL(dep DepRef) string {
	return buildADOSSHURL(dep.GetADOOrganization(), dep.GetADOProject(), dep.GetADORepo())
}

func (b *ADOBackend) BuildCloneHTTPURL(_ DepRef) (string, error) {
	return "", fmt.Errorf("Azure DevOps does not support plain HTTP cloning; use HTTPS or SSH")
}

func (b *ADOBackend) BuildCommitsAPIURL(_ DepRef, _ string) string { return "" }

func (b *ADOBackend) BuildContentsAPIURLs(_, _, _, _ string) []string { return nil }

// GitLabBackend is the backend for GitLab (gitlab.com and self-managed instances).
type GitLabBackend struct {
	hostInfo auth.HostInfo
}

func (b *GitLabBackend) Kind() string              { return "gitlab" }
func (b *GitLabBackend) IsGitHubFamily() bool      { return false }
func (b *GitLabBackend) IsGeneric() bool            { return true }
func (b *GitLabBackend) GetHostInfo() auth.HostInfo { return b.hostInfo }

func (b *GitLabBackend) BuildCloneHTTPSURL(dep DepRef, token string, authScheme string) string {
	host := urlHost(dep, b.hostInfo)
	port := dep.GetPort()
	if token != "" && authScheme != "bearer" {
		return buildGitLabHTTPSCloneURL(host, dep.GetRepoURL(), token, port)
	}
	return buildHTTPSCloneURL(host, dep.GetRepoURL(), "", port)
}

func (b *GitLabBackend) BuildCloneSSHURL(dep DepRef) string {
	host := urlHost(dep, b.hostInfo)
	return buildSSHURL(host, dep.GetRepoURL(), dep.GetPort())
}

func (b *GitLabBackend) BuildCloneHTTPURL(dep DepRef) (string, error) {
	host := urlHost(dep, b.hostInfo)
	n := netloc(host, dep.GetPort())
	return fmt.Sprintf("http://%s/%s.git", n, dep.GetRepoURL()), nil
}

func (b *GitLabBackend) BuildCommitsAPIURL(dep DepRef, ref string) string {
	if sha40RE.MatchString(strings.ToLower(ref)) {
		return ""
	}
	proj := url.PathEscape(dep.GetRepoURL())
	return fmt.Sprintf("%s/projects/%s/repository/commits/%s", b.hostInfo.APIBase, proj, ref)
}

func (b *GitLabBackend) BuildContentsAPIURLs(owner, repo, filePath, ref string) []string {
	proj := url.PathEscape(owner + "/" + repo)
	f := url.PathEscape(filePath)
	return []string{
		fmt.Sprintf("%s/projects/%s/repository/files/%s/raw?ref=%s", b.hostInfo.APIBase, proj, f, ref),
	}
}

// GenericGitBackend is the backend for non-GitHub/non-ADO/non-GitLab hosts (Gitea/Gogs/Bitbucket, ...).
type GenericGitBackend struct {
	hostInfo auth.HostInfo
}

func (b *GenericGitBackend) Kind() string              { return "generic" }
func (b *GenericGitBackend) IsGitHubFamily() bool      { return false }
func (b *GenericGitBackend) IsGeneric() bool            { return true }
func (b *GenericGitBackend) GetHostInfo() auth.HostInfo { return b.hostInfo }

func (b *GenericGitBackend) BuildCloneHTTPSURL(dep DepRef, token string, authScheme string) string {
	host := urlHost(dep, b.hostInfo)
	port := dep.GetPort()
	if authScheme == "bearer" {
		token = ""
	}
	return buildHTTPSCloneURL(host, dep.GetRepoURL(), token, port)
}

func (b *GenericGitBackend) BuildCloneSSHURL(dep DepRef) string {
	host := urlHost(dep, b.hostInfo)
	return buildSSHURL(host, dep.GetRepoURL(), dep.GetPort())
}

func (b *GenericGitBackend) BuildCloneHTTPURL(dep DepRef) (string, error) {
	host := urlHost(dep, b.hostInfo)
	n := netloc(host, dep.GetPort())
	return fmt.Sprintf("http://%s/%s.git", n, dep.GetRepoURL()), nil
}

func (b *GenericGitBackend) BuildCommitsAPIURL(_ DepRef, _ string) string { return "" }

func (b *GenericGitBackend) BuildContentsAPIURLs(owner, repo, filePath, ref string) []string {
	host := b.hostInfo.Host
	return []string{
		fmt.Sprintf("https://%s/api/v1/repos/%s/%s/contents/%s?ref=%s", host, owner, repo, filePath, ref),
		fmt.Sprintf("https://%s/api/v3/repos/%s/%s/contents/%s?ref=%s", host, owner, repo, filePath, ref),
	}
}

// ---------------------------------------------------------------------------
// Dispatch
// ---------------------------------------------------------------------------

// BackendFor picks the right HostBackend for a DepRef.
// Falls back to GenericGitBackend when the host kind cannot be classified.
func BackendFor(dep DepRef, fallbackHost string) HostBackend {
	var host string
	var port *int
	if dep != nil && dep.GetHost() != "" {
		host = dep.GetHost()
		port = dep.GetPort()
	} else {
		if fallbackHost != "" {
			host = fallbackHost
		} else {
			host = githubhost.DefaultHost()
		}
	}

	// ADO short-circuit
	if dep != nil && dep.IsAzureDevOps() {
		info := auth.ClassifyHost(host, port)
		return &ADOBackend{hostInfo: info}
	}

	info := auth.ClassifyHost(host, port)
	return backendFromInfo(info)
}

// BackendForHost picks the right HostBackend for a bare hostname.
func BackendForHost(host string, port *int) HostBackend {
	info := auth.ClassifyHost(host, port)
	return backendFromInfo(info)
}

func backendFromInfo(info auth.HostInfo) HostBackend {
	base := gitHubFamilyBase{hostInfo: info}
	switch info.Kind {
	case "github":
		base.kind = "github"
		return &GitHubBackend{base}
	case "ghe_cloud":
		base.kind = "ghe_cloud"
		return &GHECloudBackend{base}
	case "ghes":
		base.kind = "ghes"
		return &GHESBackend{base}
	case "ado":
		return &ADOBackend{hostInfo: info}
	case "gitlab":
		return &GitLabBackend{hostInfo: info}
	default:
		return &GenericGitBackend{hostInfo: info}
	}
}

// Ensure ADOBackend satisfies a narrower interface for compile-time check.
var _ interface {
	BuildCloneSSHURL(dep DepRef) string
	GetHostInfo() auth.HostInfo
} = (*ADOBackend)(nil)
