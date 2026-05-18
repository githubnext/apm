// Package githubdownloader provides GitHub package downloading for APM dependencies.
//
// Implements GitHubPackageDownloader: git clone/fetch over HTTPS or SSH with
// auth resolution, bare-cache support, remote ref listing, raw-file download
// from GitHub/ADO/GitLab, and a resilient HTTP GET delegate.
//
// Migrated from: src/apm_cli/deps/github_downloader.py
package githubdownloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/githubnext/apm/internal/core/auth"
	"github.com/githubnext/apm/internal/models/depreference"
	"github.com/githubnext/apm/internal/utils/githubhost"
)

// ProtocolPreference controls which git transport is attempted first.
type ProtocolPreference int

const (
	ProtocolPreferHTTPS ProtocolPreference = iota
	ProtocolPreferSSH
	ProtocolHTTPSOnly
	ProtocolSSHOnly
)

// RemoteRef represents a reference returned by git ls-remote.
type RemoteRef struct {
	Name string
	SHA  string
}

// DownloadResult summarises the outcome of a package download.
type DownloadResult struct {
	DestDir   string
	SHA       string
	Ref       string
	Transport string // "https" | "ssh"
}

// RawFileResult is a raw file fetched from a remote host.
type RawFileResult struct {
	Content     []byte
	ContentType string
	ETag        string
}

// ProgressReporter receives download progress callbacks.
type ProgressReporter interface {
	Update(op string, cur, max int64, message string)
}

// GitHubPackageDownloader downloads APM packages from git hosts.
type GitHubPackageDownloader struct {
	authResolver  *auth.AuthResolver
	cacheDir      string
	concurrency   int
	timeoutSecs   float64
	allowFallback bool
	protoPref     ProtocolPreference
	httpClient    *http.Client
	mu            sync.Mutex
}

// Options controls downloader construction.
type Options struct {
	CacheDir      string
	Concurrency   int
	TimeoutSecs   float64
	AllowFallback bool
	ProtoPref     ProtocolPreference
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Concurrency:   4,
		TimeoutSecs:   300,
		AllowFallback: true,
		ProtoPref:     ProtocolPreferHTTPS,
	}
}

// New constructs a GitHubPackageDownloader.
func New(resolver *auth.AuthResolver, opts Options) *GitHubPackageDownloader {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 4
	}
	if opts.TimeoutSecs <= 0 {
		opts.TimeoutSecs = 300
	}
	if resolver == nil {
		resolver = auth.NewAuthResolver(nil)
	}
	return &GitHubPackageDownloader{
		authResolver:  resolver,
		cacheDir:      opts.CacheDir,
		concurrency:   opts.Concurrency,
		timeoutSecs:   opts.TimeoutSecs,
		allowFallback: opts.AllowFallback,
		protoPref:     opts.ProtoPref,
		httpClient: &http.Client{
			Timeout: time.Duration(opts.TimeoutSecs) * time.Second,
		},
	}
}

// -------------------------------------------------------------------
// Remote ref listing
// -------------------------------------------------------------------

// ParseLsRemoteOutput parses the output of `git ls-remote`.
func ParseLsRemoteOutput(output string) []RemoteRef {
	var refs []RemoteRef
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		refs = append(refs, RemoteRef{SHA: parts[0], Name: parts[1]})
	}
	return refs
}

// SemverSortKey returns a tuple-like sort key for a semver tag name.
// Non-semver names sort last.
var semverRe = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(.*)$`)

func SemverSortKey(name string) [4]int {
	m := semverRe.FindStringSubmatch(name)
	if m == nil {
		return [4]int{-1, 0, 0, 0}
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	pre := 0
	if m[4] != "" {
		pre = -1 // pre-release sorts before release
	}
	return [4]int{major, minor, patch, pre}
}

// SortRemoteRefs returns refs sorted newest semver first, then alphabetically.
func SortRemoteRefs(refs []RemoteRef) []RemoteRef {
	sorted := make([]RemoteRef, len(refs))
	copy(sorted, refs)
	sort.Slice(sorted, func(i, j int) bool {
		ki := SemverSortKey(sorted[i].Name)
		kj := SemverSortKey(sorted[j].Name)
		for idx := 0; idx < 4; idx++ {
			if ki[idx] != kj[idx] {
				return ki[idx] > kj[idx]
			}
		}
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}

// ListRemoteRefs runs `git ls-remote` and returns parsed refs.
func (d *GitHubPackageDownloader) ListRemoteRefs(dep *depreference.DependencyReference) ([]RemoteRef, error) {
	repoURL, err := d.buildRepoURL(dep, "https")
	if err != nil {
		return nil, err
	}
	env := d.gitEnv(dep)
	cmd := exec.Command("git", "ls-remote", "--tags", "--heads", repoURL)
	cmd.Env = append(os.Environ(), mapToEnv(env)...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ls-remote %s: %w", repoURL, err)
	}
	refs := ParseLsRemoteOutput(string(out))
	return SortRemoteRefs(refs), nil
}

// ResolveGitReference resolves a version range to a concrete ref/SHA pair.
func (d *GitHubPackageDownloader) ResolveGitReference(dep *depreference.DependencyReference) (string, string, error) {
	refs, err := d.ListRemoteRefs(dep)
	if err != nil {
		return "", "", err
	}
	want := dep.Reference
	for _, r := range refs {
		if r.Name == want || r.Name == "refs/tags/"+want || r.Name == "refs/heads/"+want {
			return r.Name, r.SHA, nil
		}
	}
	// SHA pinned?
	if len(want) >= 7 {
		for _, r := range refs {
			if strings.HasPrefix(r.SHA, want) {
				return r.Name, r.SHA, nil
			}
		}
	}
	return "", "", fmt.Errorf("ref %q not found in remote %s", want, dep.RepoURL)
}

// -------------------------------------------------------------------
// Clone / download
// -------------------------------------------------------------------

// Download clones or updates the package into destDir.
func (d *GitHubPackageDownloader) Download(dep *depreference.DependencyReference, destDir string, progress ProgressReporter) (*DownloadResult, error) {
	repoURL, err := d.buildRepoURL(dep, "https")
	if err != nil {
		return nil, err
	}
	env := d.gitEnv(dep)

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}

	ref := dep.Reference
	if ref == "" {
		ref = "HEAD"
	}

	isSubdir := dep.VirtualPath != "" || dep.ADOProject != ""

	args := []string{"clone", "--depth=1", "--branch", ref, repoURL, destDir}
	if isSubdir {
		args = []string{"clone", "--depth=1", "--filter=blob:none", "--sparse", "--branch", ref, repoURL, destDir}
	}

	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(), mapToEnv(env)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		// Fallback: try without --branch for bare SHA refs
		if d.allowFallback {
			return d.cloneFallback(dep, repoURL, destDir, env)
		}
		return nil, fmt.Errorf("git clone failed: %w\n%s", err, out)
	}

	sha, _ := d.resolveHEAD(destDir)
	return &DownloadResult{DestDir: destDir, SHA: sha, Ref: ref, Transport: "https"}, nil
}

func (d *GitHubPackageDownloader) cloneFallback(dep *depreference.DependencyReference, repoURL, destDir string, env map[string]string) (*DownloadResult, error) {
	_ = os.RemoveAll(destDir)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}
	cmd := exec.Command("git", "clone", "--depth=1", repoURL, destDir)
	cmd.Env = append(os.Environ(), mapToEnv(env)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git clone fallback failed: %w\n%s", err, out)
	}
	sha, _ := d.resolveHEAD(destDir)
	return &DownloadResult{DestDir: destDir, SHA: sha, Ref: dep.Reference, Transport: "https"}, nil
}

func (d *GitHubPackageDownloader) resolveHEAD(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// -------------------------------------------------------------------
// Raw file download
// -------------------------------------------------------------------

// DownloadRawFile fetches a single file from a remote host.
func (d *GitHubPackageDownloader) DownloadRawFile(dep *depreference.DependencyReference, filePath string) (*RawFileResult, error) {
	token := d.resolveToken(dep)
	rawURL := d.buildRawFileURL(dep, filePath)
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	if token != nil {
		req.Header.Set("Authorization", "token "+*token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, rawURL)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &RawFileResult{
		Content:     body,
		ContentType: resp.Header.Get("Content-Type"),
		ETag:        resp.Header.Get("ETag"),
	}, nil
}

// -------------------------------------------------------------------
// Artifactory support
// -------------------------------------------------------------------

// IsArtifactoryOnly reports whether the dep host is Artifactory-only.
func (d *GitHubPackageDownloader) IsArtifactoryOnly() bool {
	return os.Getenv("APM_ARTIFACTORY_ONLY") == "1"
}

// DownloadArtifactoryArchive fetches a tarball from Artifactory.
func (d *GitHubPackageDownloader) DownloadArtifactoryArchive(dep *depreference.DependencyReference, destDir string) error {
	baseURL := os.Getenv("APM_ARTIFACTORY_BASE_URL")
	if baseURL == "" {
		return errors.New("APM_ARTIFACTORY_BASE_URL not set")
	}
	token := os.Getenv("APM_ARTIFACTORY_TOKEN")
	archiveURL := strings.TrimRight(baseURL, "/") + "/" + dep.RepoURL + "/" + dep.Reference + ".tar.gz"

	req, err := http.NewRequest("GET", archiveURL, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Artifactory HTTP %d", resp.StatusCode)
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	tmp := filepath.Join(destDir, "archive.tar.gz")
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func (d *GitHubPackageDownloader) buildRepoURL(dep *depreference.DependencyReference, proto string) (string, error) {
	host := dep.Host
	if host == "" {
		host = githubhost.DefaultHost()
	}
	if proto == "ssh" {
		return fmt.Sprintf("git@%s:%s.git", host, dep.RepoURL), nil
	}
	token := d.resolveToken(dep)
	if token != nil {
		u := &url.URL{
			Scheme: "https",
			User:   url.UserPassword("x-access-token", *token),
			Host:   host,
			Path:   "/" + dep.RepoURL + ".git",
		}
		return u.String(), nil
	}
	return fmt.Sprintf("https://%s/%s.git", host, dep.RepoURL), nil
}

func (d *GitHubPackageDownloader) buildRawFileURL(dep *depreference.DependencyReference, filePath string) string {
	host := dep.Host
	if host == "" {
		host = "raw.githubusercontent.com"
	} else {
		host = "raw." + host
	}
	ref := dep.Reference
	if ref == "" {
		ref = "HEAD"
	}
	return fmt.Sprintf("https://%s/%s/%s/%s", host, dep.RepoURL, ref, filePath)
}

func (d *GitHubPackageDownloader) resolveToken(dep *depreference.DependencyReference) *string {
	host := dep.Host
	if host == "" {
		host = githubhost.DefaultHost()
	}
	ctx := d.authResolver.Resolve(host, "", nil)
	if ctx == nil {
		return nil
	}
	return ctx.Token
}

func (d *GitHubPackageDownloader) gitEnv(dep *depreference.DependencyReference) map[string]string {
	host := dep.Host
	if host == "" {
		host = githubhost.DefaultHost()
	}
	ctx := d.authResolver.Resolve(host, "", nil)
	if ctx == nil {
		return map[string]string{"GIT_TERMINAL_PROMPT": "0"}
	}
	env := make(map[string]string)
	for k, v := range ctx.GitEnv {
		env[k] = v
	}
	env["GIT_TERMINAL_PROMPT"] = "0"
	return env
}

func mapToEnv(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

// -------------------------------------------------------------------
// Transport plan / protocol selection
// -------------------------------------------------------------------

// TransportPlan describes which transports to attempt, in order.
type TransportPlan struct {
	Primary   string // "https" | "ssh"
	Fallbacks []string
}

// BuildTransportPlan returns the ordered list of transports for a given preference.
func BuildTransportPlan(pref ProtocolPreference, allowFallback bool) TransportPlan {
	switch pref {
	case ProtocolSSHOnly:
		return TransportPlan{Primary: "ssh"}
	case ProtocolHTTPSOnly:
		return TransportPlan{Primary: "https"}
	case ProtocolPreferSSH:
		if allowFallback {
			return TransportPlan{Primary: "ssh", Fallbacks: []string{"https"}}
		}
		return TransportPlan{Primary: "ssh"}
	default:
		if allowFallback {
			return TransportPlan{Primary: "https", Fallbacks: []string{"ssh"}}
		}
		return TransportPlan{Primary: "https"}
	}
}

// -------------------------------------------------------------------
// Validation helper
// -------------------------------------------------------------------

// ValidateAPMPackage checks that a downloaded directory contains a valid apm.yml.
func ValidateAPMPackage(dir string) error {
	candidates := []string{
		filepath.Join(dir, "apm.yml"),
		filepath.Join(dir, ".apm", "apm.yml"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no apm.yml found in %s", dir)
}

// -------------------------------------------------------------------
// Bare cache helpers (exported for testing)
// -------------------------------------------------------------------

// BareCloneURL builds the bare-cache path for a repo URL.
func BareCloneURL(cacheDir, repoURL string) string {
	safe := regexp.MustCompile(`[^a-zA-Z0-9_.-]`).ReplaceAllString(repoURL, "_")
	return filepath.Join(cacheDir, safe+".git")
}

// -------------------------------------------------------------------
// ADO (Azure DevOps) raw file download
// -------------------------------------------------------------------

// DownloadADOFile fetches a file from Azure DevOps REST API.
func (d *GitHubPackageDownloader) DownloadADOFile(org, project, repo, ref, filePath string) ([]byte, error) {
	token := os.Getenv("ADO_APM_PAT")
	if token == "" {
		token = os.Getenv("ADO_TOKEN")
	}
	apiURL := fmt.Sprintf(
		"https://dev.azure.com/%s/%s/_apis/git/repositories/%s/items?path=%s&versionDescriptor.version=%s&api-version=7.0",
		url.PathEscape(org), url.PathEscape(project), url.PathEscape(repo),
		url.QueryEscape(filePath), url.QueryEscape(ref),
	)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.SetBasicAuth("", token)
	}
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ADO HTTP %d for %s", resp.StatusCode, apiURL)
	}
	return io.ReadAll(resp.Body)
}

// -------------------------------------------------------------------
// Progress
// -------------------------------------------------------------------

// LogProgressReporter writes progress to stderr when APM_DEBUG is set.
type LogProgressReporter struct{}

func (l *LogProgressReporter) Update(op string, cur, max int64, message string) {
	if os.Getenv("APM_DEBUG") == "" {
		return
	}
	pct := ""
	if max > 0 {
		pct = fmt.Sprintf(" %.0f%%", float64(cur)/float64(max)*100)
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] %s%s %s\n", op, pct, message)
}

// -------------------------------------------------------------------
// Registry config
// -------------------------------------------------------------------

// RegistryConfig holds per-registry authentication settings parsed from environment.
type RegistryConfig struct {
	ArtifactoryBaseURL string
	ArtifactoryToken   string
	NpmRegistry        string
	NpmToken           string
}

// LoadRegistryConfig reads registry config from environment variables.
func LoadRegistryConfig() RegistryConfig {
	return RegistryConfig{
		ArtifactoryBaseURL: os.Getenv("APM_ARTIFACTORY_BASE_URL"),
		ArtifactoryToken:   os.Getenv("APM_ARTIFACTORY_TOKEN"),
		NpmRegistry:        os.Getenv("APM_NPM_REGISTRY"),
		NpmToken:           os.Getenv("APM_NPM_TOKEN"),
	}
}

// -------------------------------------------------------------------
// Sanitize git errors (remove tokens from messages)
// -------------------------------------------------------------------

// SanitizeGitError redacts bearer tokens and credentials from git error messages.
func SanitizeGitError(msg string) string {
	// Redact https://x-access-token:TOKEN@...
	re := regexp.MustCompile(`(https?://[^:@/]+:)[^@]+(@)`)
	msg = re.ReplaceAllString(msg, "${1}[REDACTED]${2}")
	// Redact Authorization: token XYZ
	re2 := regexp.MustCompile(`(?i)(Authorization:\s*(?:token|bearer)\s+)\S+`)
	return re2.ReplaceAllString(msg, "${1}[REDACTED]")
}

// -------------------------------------------------------------------
// JSON serialisation helpers (used by benchmarks)
// -------------------------------------------------------------------

// DownloadResultJSON marshals a DownloadResult to JSON.
func (r *DownloadResult) JSON() ([]byte, error) {
	return json.Marshal(r)
}
