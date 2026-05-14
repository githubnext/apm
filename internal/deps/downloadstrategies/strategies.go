// Package downloadstrategies implements the DownloadDelegate -- the
// backend-specific HTTP download logic for APM packages.
//
// Encapsulates resilient HTTP GET, GitHub Contents API, Azure DevOps,
// GitLab, Artifactory archive, and generic-host file download logic.
// The owning GitHubPackageDownloader creates a single DownloadDelegate
// and delegates all download operations to it (Facade/Delegate pattern).
//
// Migrated from: src/apm_cli/deps/download_strategies.py
package downloadstrategies

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/githubnext/apm/internal/core/auth"
	"github.com/githubnext/apm/internal/models/depreference"
	"github.com/githubnext/apm/internal/utils/githubhost"
)

// HostProvider is the interface the DownloadDelegate requires from its owner.
// This avoids a circular package dependency on the github_downloader package.
type HostProvider interface {
	// GithubToken returns the GitHub personal access token (may be empty).
	GithubToken() string
	// AdoToken returns the Azure DevOps PAT (may be empty).
	AdoToken() string
	// ArtifactoryToken returns the Artifactory bearer token (may be empty).
	ArtifactoryToken() string
	// GithubHost returns the configured GitHub host (may be empty for default).
	GithubHost() string
	// AuthResolver returns the authentication resolver.
	AuthResolver() *auth.AuthResolver
	// ResilientGet performs an HTTP GET with retry/rate-limit handling.
	// Callers should treat a non-nil error as exhausted retries.
	ResilientGet(reqURL string, headers map[string]string, timeoutSecs int) (*http.Response, error)
}

// resolveToken extracts the token string from *string (nil -> "").
func resolveToken(t *string) string {
	if t == nil {
		return ""
	}
	return *t
}

// authResolve wraps AuthResolver.Resolve, handling the *int port parameter.
func authResolve(ar *auth.AuthResolver, host, org string, port int) (token, source string) {
	var portPtr *int
	if port != 0 {
		portPtr = &port
	}
	ctx := ar.Resolve(host, org, portPtr)
	if ctx == nil {
		return "", ""
	}
	return resolveToken(ctx.Token), ctx.Source
}

// DownloadDelegate encapsulates backend-specific download logic.
//
// Holds real implementations of HTTP resilient-get, URL building, and
// file download for GitHub, Azure DevOps, and Artifactory backends.
type DownloadDelegate struct {
	host HostProvider
}

// New creates a DownloadDelegate that delegates shared state to host.
func New(host HostProvider) *DownloadDelegate {
	return &DownloadDelegate{host: host}
}

// debug prints a message when APM_DEBUG is set.
func debug(msg string) {
	if os.Getenv("APM_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", msg)
	}
}

// ---------------------------------------------------------------------------
// HTTP resilient GET (standalone helper for callers without a HostProvider)
// ---------------------------------------------------------------------------

// ResilientGet performs an HTTP GET with exponential-backoff retry on 429/503
// and rate-limit header awareness.
//
// Returns the *http.Response and nil on success.  If all retries are
// exhausted it returns the last response (which may be rate-limited) plus a
// non-nil error.
func ResilientGet(reqURL string, headers map[string]string, timeoutSecs, maxRetries int) (*http.Response, error) {
	if timeoutSecs <= 0 {
		timeoutSecs = 30
	}
	if maxRetries <= 0 {
		maxRetries = 3
	}
	client := &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second}

	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries-1 {
				wait := jitter(math.Pow(2, float64(attempt)))
				debug(fmt.Sprintf("Connection error, retry in %.1fs (attempt %d/%d)", wait, attempt+1, maxRetries))
				time.Sleep(time.Duration(wait*float64(time.Second)))
			}
			continue
		}

		// Rate limiting: 429, 503, or 403 with X-RateLimit-Remaining: 0.
		isRateLimited := resp.StatusCode == 429 || resp.StatusCode == 503
		if !isRateLimited && resp.StatusCode == 403 {
			if rem := resp.Header.Get("X-RateLimit-Remaining"); rem != "" {
				if n, err := strconv.Atoi(rem); err == nil && n == 0 {
					isRateLimited = true
				}
			}
		}

		if isRateLimited {
			lastResp = resp
			wait := backoffFromRateLimitHeaders(resp, attempt)
			debug(fmt.Sprintf("Rate limited (%d), retry in %.1fs (attempt %d/%d)", resp.StatusCode, wait, attempt+1, maxRetries))
			time.Sleep(time.Duration(wait * float64(time.Second)))
			continue
		}

		// Log rate-limit proximity.
		if rem := resp.Header.Get("X-RateLimit-Remaining"); rem != "" {
			if n, err := strconv.Atoi(rem); err == nil && n < 10 {
				debug(fmt.Sprintf("GitHub API rate limit low: %d requests remaining", n))
			}
		}
		return resp, nil
	}

	if lastResp != nil {
		return lastResp, fmt.Errorf("rate limit retries exhausted for %s", reqURL)
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("all %d attempts failed for %s", maxRetries, reqURL)
}

func jitter(base float64) float64 {
	if base > 30 {
		base = 30
	}
	return base * (0.5 + rand.Float64())
}

func backoffFromRateLimitHeaders(resp *http.Response, attempt int) float64 {
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		if v, err := strconv.ParseFloat(ra, 64); err == nil {
			if v < 60 {
				return v
			}
			return 60
		}
	}
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
			wait := float64(ts) - float64(time.Now().Unix())
			if wait > 0 && wait < 60 {
				return wait
			}
		}
	}
	return jitter(math.Pow(2, float64(attempt)))
}

// ---------------------------------------------------------------------------
// Repository URL building
// ---------------------------------------------------------------------------

// BuildRepoURLOptions controls how BuildRepoURL constructs its result.
type BuildRepoURLOptions struct {
	RepoRef    string
	UseSSH     bool
	DepRef     *depreference.DependencyReference
	Token      string
	AuthScheme string // "basic" | "bearer"  (default: "basic")
}

// BuildRepoURL constructs the repository URL for git clone operations.
// Supports GitHub, Azure DevOps, GitLab, and generic hosts.
func (d *DownloadDelegate) BuildRepoURL(opts BuildRepoURLOptions) string {
	var host string
	if opts.DepRef != nil && opts.DepRef.Host != "" {
		host = opts.DepRef.Host
	} else if h := d.host.GithubHost(); h != "" {
		host = h
	} else {
		host = githubhost.DefaultHost()
	}

	token := opts.Token
	if token == "" {
		token = d.host.GithubToken()
	}

	repoRef := opts.RepoRef
	if opts.DepRef != nil && repoRef == "" {
		repoRef = opts.DepRef.RepoURL
	}

	var port int
	if opts.DepRef != nil {
		port = opts.DepRef.Port
	}

	if opts.UseSSH {
		return buildSSHURL(host, repoRef, port)
	}
	if token != "" {
		return buildHTTPSCloneURL(host, repoRef, token, port)
	}
	return buildHTTPSCloneURL(host, repoRef, "", port)
}

func buildSSHURL(host, repoRef string, port int) string {
	if port != 0 {
		return fmt.Sprintf("ssh://git@%s:%d/%s.git", host, port, repoRef)
	}
	return fmt.Sprintf("git@%s:%s.git", host, repoRef)
}

func buildHTTPSCloneURL(host, repoRef, token string, port int) string {
	var netloc string
	if port != 0 {
		netloc = fmt.Sprintf("%s:%d", host, port)
	} else {
		netloc = host
	}
	if token != "" {
		return fmt.Sprintf("https://x-access-token:%s@%s/%s.git", token, netloc, repoRef)
	}
	return fmt.Sprintf("https://%s/%s.git", netloc, repoRef)
}

// ---------------------------------------------------------------------------
// Artifactory helpers
// ---------------------------------------------------------------------------

// GetArtifactoryHeaders returns HTTP headers for Artifactory requests.
func (d *DownloadDelegate) GetArtifactoryHeaders() map[string]string {
	headers := make(map[string]string)
	if tok := d.host.ArtifactoryToken(); tok != "" {
		headers["Authorization"] = "Bearer " + tok
	}
	return headers
}

// ArtifactoryDownloadResult holds the result of an Artifactory archive download.
type ArtifactoryDownloadResult struct {
	Data []byte
	Err  error
}

// DownloadArtifactoryArchive downloads an archive from Artifactory.
func (d *DownloadDelegate) DownloadArtifactoryArchive(archiveURL string) ArtifactoryDownloadResult {
	headers := d.GetArtifactoryHeaders()

	resp, err := ResilientGet(archiveURL, headers, 120, 3)
	if err != nil {
		return ArtifactoryDownloadResult{Err: fmt.Errorf("artifactory archive download: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ArtifactoryDownloadResult{
			Err: fmt.Errorf("artifactory archive HTTP %d for %s", resp.StatusCode, archiveURL),
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ArtifactoryDownloadResult{Err: fmt.Errorf("reading artifactory archive: %w", err)}
	}
	return ArtifactoryDownloadResult{Data: data}
}

// DownloadFileFromArtifactory downloads a single file from Artifactory.
func (d *DownloadDelegate) DownloadFileFromArtifactory(fileURL string) ([]byte, error) {
	headers := d.GetArtifactoryHeaders()
	resp, err := ResilientGet(fileURL, headers, 30, 3)
	if err != nil {
		return nil, fmt.Errorf("artifactory file download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, fileURL)
	}
	return io.ReadAll(resp.Body)
}

// ---------------------------------------------------------------------------
// Raw download (CDN fast-path for github.com)
// ---------------------------------------------------------------------------

// TryRawDownload attempts to fetch a file via raw.githubusercontent.com.
// Returns nil if the file was not found or the request failed.
func (d *DownloadDelegate) TryRawDownload(owner, repo, ref, filePath string) []byte {
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, filePath)
	resp, err := ResilientGet(rawURL, nil, 30, 2)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return data
}

// ---------------------------------------------------------------------------
// Azure DevOps file download
// ---------------------------------------------------------------------------

// buildADOAPIURL constructs the Azure DevOps Items API URL for a file.
func buildADOAPIURL(org, project, repo, filePath, ref, host string) string {
	if host == "" {
		host = "dev.azure.com"
	}
	return fmt.Sprintf(
		"https://%s/%s/%s/_apis/git/repositories/%s/items?path=%s&versionType=branch&version=%s&api-version=6.0",
		host, url.PathEscape(org), url.PathEscape(project), url.PathEscape(repo),
		url.QueryEscape(filePath), url.QueryEscape(ref),
	)
}

func (d *DownloadDelegate) DownloadADOFile(depRef *depreference.DependencyReference, filePath, ref string) ([]byte, error) {
	if depRef == nil {
		return nil, fmt.Errorf("nil dep_ref for ADO download")
	}
	if depRef.ADOOrganization == "" || depRef.ADOProject == "" || depRef.ADORepo == "" {
		return nil, fmt.Errorf(
			"invalid ADO dep_ref: missing org/project/repo (got org=%q project=%q repo=%q)",
			depRef.ADOOrganization, depRef.ADOProject, depRef.ADORepo,
		)
	}

	host := depRef.Host
	if host == "" {
		host = "dev.azure.com"
	}
	apiURL := buildADOAPIURL(depRef.ADOOrganization, depRef.ADOProject, depRef.ADORepo, filePath, ref, host)

	headers := make(map[string]string)
	if tok := d.host.AdoToken(); tok != "" {
		authBytes := []byte(":" + tok)
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString(authBytes)
	}

	resp, err := d.host.ResilientGet(apiURL, headers, 30)
	if err != nil {
		return nil, fmt.Errorf("ADO download network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode == http.StatusNotFound {
		if ref == "main" || ref == "master" {
			fallbackRef := "master"
			if ref == "master" {
				fallbackRef = "main"
			}
			fallbackURL := buildADOAPIURL(depRef.ADOOrganization, depRef.ADOProject, depRef.ADORepo, filePath, fallbackRef, host)
			resp2, err2 := d.host.ResilientGet(fallbackURL, headers, 30)
			if err2 != nil {
				return nil, fmt.Errorf("ADO fallback download failed: %w", err2)
			}
			defer resp2.Body.Close()
			if resp2.StatusCode == http.StatusOK {
				return io.ReadAll(resp2.Body)
			}
			return nil, fmt.Errorf("file not found: %s in %s (tried refs: %s, %s)", filePath, depRef.RepoURL, ref, fallbackRef)
		}
		return nil, fmt.Errorf("file not found: %s at ref %q in %s", filePath, ref, depRef.RepoURL)
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("authentication failed for Azure DevOps %s", depRef.RepoURL)
	}
	return nil, fmt.Errorf("ADO download HTTP %d for %s", resp.StatusCode, apiURL)
}

// ---------------------------------------------------------------------------
// GitLab file download
// ---------------------------------------------------------------------------

// DownloadGitLabFile downloads a file via the GitLab REST v4 API.
func (d *DownloadDelegate) DownloadGitLabFile(depRef *depreference.DependencyReference, filePath, ref string) ([]byte, error) {
	if depRef == nil {
		return nil, fmt.Errorf("nil dep_ref for GitLab download")
	}
	host := depRef.Host
	if host == "" {
		host = githubhost.DefaultHost()
	}
	projectPath := depRef.RepoURL
	if projectPath == "" {
		return nil, fmt.Errorf("missing repository path for GitLab file download")
	}

	ar := d.host.AuthResolver()
	var token string
	if ar != nil {
		org := ""
		parts := strings.SplitN(projectPath, "/", 2)
		if len(parts) > 0 {
			org = parts[0]
		}
		t, _ := authResolve(ar, host, org, depRef.Port)
		token = t
	}

	headers := map[string]string{}
	if token != "" {
		headers["PRIVATE-TOKEN"] = token
	}

	enc := url.PathEscape(projectPath)
	encFile := url.PathEscape(filePath)
	encRef := url.QueryEscape(ref)
	apiURL := fmt.Sprintf("https://%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s", host, enc, encFile, encRef)

	resp, err := d.host.ResilientGet(apiURL, headers, 30)
	if err != nil {
		return nil, fmt.Errorf("GitLab download error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}
	if resp.StatusCode == http.StatusNotFound {
		// Try the other default branch.
		if ref == "main" || ref == "master" {
			fallbackRef := "master"
			if ref == "master" {
				fallbackRef = "main"
			}
			encFallback := url.QueryEscape(fallbackRef)
			fallbackURL := fmt.Sprintf("https://%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s", host, enc, encFile, encFallback)
			resp2, err2 := d.host.ResilientGet(fallbackURL, headers, 30)
			if err2 == nil {
				defer resp2.Body.Close()
				if resp2.StatusCode == http.StatusOK {
					return io.ReadAll(resp2.Body)
				}
			}
		}
		return nil, fmt.Errorf("file not found: %s at ref %q in %s", filePath, ref, projectPath)
	}
	return nil, fmt.Errorf("GitLab download HTTP %d", resp.StatusCode)
}

// ---------------------------------------------------------------------------
// GitHub file download (Contents API)
// ---------------------------------------------------------------------------

// DownloadGitHubFile downloads a file from a GitHub (or GHES/generic) repository.
func (d *DownloadDelegate) DownloadGitHubFile(depRef *depreference.DependencyReference, filePath, ref string) ([]byte, error) {
	if depRef == nil {
		return nil, fmt.Errorf("nil dep_ref for GitHub download")
	}
	host := depRef.Host
	if host == "" {
		host = githubhost.DefaultHost()
	}

	parts := strings.SplitN(depRef.RepoURL, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repo_url %q: expected owner/repo", depRef.RepoURL)
	}
	owner, repo := parts[0], parts[1]

	ar := d.host.AuthResolver()
	var token string
	if ar != nil {
		t, _ := authResolve(ar, host, owner, depRef.Port)
		token = t
	}

	isGitHubHost := githubhost.IsGitHubHostname(host) || d.isConfiguredGHES(host)

	// CDN fast-path for github.com without a token.
	if strings.EqualFold(host, "github.com") && token == "" {
		if data := d.TryRawDownload(owner, repo, ref, filePath); data != nil {
			return data, nil
		}
		// Try alternate default branch.
		if ref == "main" || ref == "master" {
			alt := "master"
			if ref == "master" {
				alt = "main"
			}
			if data := d.TryRawDownload(owner, repo, alt, filePath); data != nil {
				return data, nil
			}
		}
		// Fall through to Contents API.
	}

	// For non-GitHub generic hosts: try raw URL first.
	if !isGitHubHost {
		rawURL := fmt.Sprintf("https://%s/%s/%s/raw/%s/%s", host, owner, repo, ref, filePath)
		rawHeaders := d.buildGenericHostAuthHeaders(host, depRef, nil)
		if resp, err := d.host.ResilientGet(rawURL, rawHeaders, 30); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return io.ReadAll(resp.Body)
			}
		}
	}

	// Contents API path.
	apiURLs := d.buildContentsAPIURLs(host, owner, repo, filePath, ref, isGitHubHost)
	if len(apiURLs) == 0 {
		return nil, fmt.Errorf("could not build Contents API URL for %s", depRef.RepoURL)
	}

	var apiHeaders map[string]string
	if isGitHubHost {
		apiHeaders = map[string]string{"Accept": "application/vnd.github.v3.raw"}
		if token != "" {
			apiHeaders["Authorization"] = "token " + token
		}
	} else {
		apiHeaders = d.buildGenericHostAuthHeaders(host, depRef, nil)
		apiHeaders["Accept"] = "application/json"
	}

	for _, apiURL := range apiURLs {
		resp, err := d.host.ResilientGet(apiURL, apiHeaders, 30)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return extractContentsAPIPayload(resp, isGitHubHost)
		}
		if resp.StatusCode == http.StatusNotFound {
			continue
		}
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("authentication failed for %s/%s on %s", owner, repo, host)
		}
	}

	// Try alternate default branch as final fallback.
	if ref == "main" || ref == "master" {
		alt := "master"
		if ref == "master" {
			alt = "main"
		}
		altURLs := d.buildContentsAPIURLs(host, owner, repo, filePath, alt, isGitHubHost)
		for _, apiURL := range altURLs {
			resp, err := d.host.ResilientGet(apiURL, apiHeaders, 30)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return extractContentsAPIPayload(resp, isGitHubHost)
			}
		}
	}

	return nil, fmt.Errorf("file not found: %s at ref %q in %s", filePath, ref, depRef.RepoURL)
}

// buildContentsAPIURLs returns ordered API URL candidates for the given file.
func (d *DownloadDelegate) buildContentsAPIURLs(host, owner, repo, filePath, ref string, isGitHubHost bool) []string {
	if isGitHubHost {
		apiBase := "api.github.com"
		if !strings.EqualFold(host, "github.com") {
			apiBase = host + "/api/v3"
		}
		return []string{fmt.Sprintf("https://%s/repos/%s/%s/contents/%s?ref=%s", apiBase, owner, repo, filePath, url.QueryEscape(ref))}
	}
	// Generic host: try multiple API version paths.
	return []string{
		fmt.Sprintf("https://%s/api/v1/repos/%s/%s/contents/%s?ref=%s", host, owner, repo, filePath, url.QueryEscape(ref)),
		fmt.Sprintf("https://%s/api/v3/repos/%s/%s/contents/%s?ref=%s", host, owner, repo, filePath, url.QueryEscape(ref)),
	}
}

// buildGenericHostAuthHeaders builds auth headers for non-GitHub hosts.
func (d *DownloadDelegate) buildGenericHostAuthHeaders(host string, depRef *depreference.DependencyReference, accept *string) map[string]string {
	headers := make(map[string]string)
	if accept != nil {
		headers["Accept"] = *accept
	}
	ar := d.host.AuthResolver()
	if ar == nil {
		return headers
	}
	var port int
	org := ""
	if depRef != nil {
		port = depRef.Port
		if parts := strings.SplitN(depRef.RepoURL, "/", 2); len(parts) > 0 {
			org = parts[0]
		}
	}
	token, src := authResolve(ar, host, org, port)
	if token == "" {
		return headers
	}
	// Only forward tokens for credential-helper-sourced or org-scoped sources,
	// or explicitly configured GHES.
	if src == "git-credential-fill" || strings.HasPrefix(src, "GITHUB_APM_PAT_") || d.isConfiguredGHES(host) {
		headers["Authorization"] = "token " + token
	}
	return headers
}

// isConfiguredGHES reports whether host is set as the configured GHES via GITHUB_HOST.
func (d *DownloadDelegate) isConfiguredGHES(host string) bool {
	ghHost := strings.TrimSpace(os.Getenv("GITHUB_HOST"))
	if ghHost == "" {
		return false
	}
	return strings.EqualFold(ghHost, host)
}

// extractContentsAPIPayload decodes a Contents-API response into raw bytes.
//
// GitHub family: returns response.Body bytes directly (vnd.github.v3.raw).
// Generic hosts (Gitea/Gogs): the server returns a JSON envelope
// {"content": "<base64>", "encoding": "base64"}.
func extractContentsAPIPayload(resp *http.Response, isGitHubHost bool) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if isGitHubHost {
		return body, nil
	}
	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(ct, "json") && (len(body) == 0 || body[0] != '{') {
		return body, nil
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return body, nil
	}
	contentField, ok := payload["content"]
	if !ok {
		return body, nil
	}
	encoding, _ := payload["encoding"].(string)
	contentStr, _ := contentField.(string)
	if strings.ToLower(encoding) == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(contentStr, "\n", ""))
		if err != nil {
			return body, nil
		}
		return decoded, nil
	}
	return []byte(contentStr), nil
}
