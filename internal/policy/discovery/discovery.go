// Package discovery implements auto-discovery and fetching of org-level apm-policy.yml files.
// Migrated from src/apm_cli/policy/discovery.py.
//
// Discovery flow:
//  1. Extract org from git remote (github.com/contoso/my-project -> "contoso")
//  2. Fetch <org>/.github/apm-policy.yml via GitHub API (Contents API)
//  3. Resolve inheritance chain via policy/inheritance package
//  4. Cache the merged effective policy with chain metadata
//  5. Parse and return the policy
//
// Supports:
//   - GitHub.com and GitHub Enterprise (*.ghe.com)
//   - Manual override via --policy <path|url>
//   - Cache with TTL (default 1 hour), stale fallback up to MAX_STALE_TTL
//   - Atomic cache writes (temp file + os.Rename)
//   - Hash-pin verification ("algo:hex" format) for supply-chain hardening
package discovery

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/githubnext/apm/internal/policy/schema"
	"github.com/githubnext/apm/internal/utils/pathsecurity"
)

const (
	policyCacheDir     = ".policy-cache"
	defaultCacheTTL    = 3600              // 1 hour (seconds)
	maxStaleTTL        = 7 * 24 * 3600    // 7 days
	cacheSchemaVersion = "3"
)

// scpLikeRE matches SCP-style SSH remote URLs: user@host:path
var scpLikeRE = regexp.MustCompile(`^(?:[^@:/?#]+@)(?P<host>[^:/?#]+):(?P<path>.+)$`)

// PolicyFetchResult is the outcome of a policy fetch attempt.
// The Outcome field discriminates discovery outcomes.
type PolicyFetchResult struct {
	Policy          *schema.ApmPolicy
	Source          string // "org:contoso/.github", "file:/path", "url:https://..."
	Cached          bool
	Err             string // error message if fetch failed
	CacheAgeSeconds int
	CacheStale      bool
	FetchErr        string
	Outcome         string
	RawBytesHash    string // "<algo>:<hex>" of leaf bytes off the wire
	ExpectedHash    string // pin that was checked, if any
}

// Found returns true when a policy was found.
func (r *PolicyFetchResult) Found() bool { return r.Policy != nil }

// cacheEntry is an internal representation of a cached policy read.
type cacheEntry struct {
	Policy       *schema.ApmPolicy
	Source       string
	AgeSeconds   int
	Stale        bool
	ChainRefs    []string
	Fingerprint  string
	RawBytesHash string
}

// ---------------------------------------------------------------------------
// Public entry points
// ---------------------------------------------------------------------------

// DiscoverPolicyWithChain discovers policy with full inheritance chain resolution.
// This is the shared entry point for all command sites that need chain-aware policy discovery.
func DiscoverPolicyWithChain(projectRoot string, expectedHash string) *PolicyFetchResult {
	if os.Getenv("APM_POLICY_DISABLE") == "1" {
		return &PolicyFetchResult{Outcome: "disabled"}
	}

	// If no explicit hash, read from project apm.yml (stub -- just pass through)
	if expectedHash == "" {
		if pin := readProjectHashPin(projectRoot); pin != "" {
			expectedHash = pin
		}
	}

	fetchResult := DiscoverPolicy(projectRoot, "", false, expectedHash)

	// Chain resolution if leaf has extends (stub -- not implemented in this iteration)
	_ = fetchResult
	return fetchResult
}

// DiscoverPolicy discovers and loads the applicable policy for a project.
//
// Resolution order:
//  1. If policyOverride is a local file path -- load from file
//  2. If policyOverride is an https:// URL -- fetch from URL
//  3. If policyOverride is "owner/repo" or "host/owner/repo" -- fetch from repo
//  4. If policyOverride is "" -- auto-discover from project's git remote
func DiscoverPolicy(projectRoot, policyOverride string, noCache bool, expectedHash string) *PolicyFetchResult {
	if policyOverride != "" {
		// Try as local file
		if info, err := os.Stat(policyOverride); err == nil && !info.IsDir() {
			return loadFromFile(policyOverride, expectedHash)
		}
		if strings.HasPrefix(policyOverride, "http://") {
			return &PolicyFetchResult{
				Err:     "Refusing plaintext http:// policy URL -- use https://",
				Source:  "url:" + policyOverride,
				Outcome: "cache_miss_fetch_fail",
			}
		}
		if strings.HasPrefix(policyOverride, "https://") {
			return fetchFromURL(policyOverride, projectRoot, noCache, expectedHash)
		}
		if policyOverride != "org" {
			return fetchFromRepo(policyOverride, projectRoot, noCache, expectedHash)
		}
	}
	return autoDiscover(projectRoot, noCache, expectedHash)
}

// ---------------------------------------------------------------------------
// File loading
// ---------------------------------------------------------------------------

func loadFromFile(path, expectedHash string) *PolicyFetchResult {
	content, err := os.ReadFile(path)
	if err != nil {
		return &PolicyFetchResult{
			Err:     fmt.Sprintf("Failed to read %s: %v", path, err),
			Outcome: "cache_miss_fetch_fail",
		}
	}
	sourceLabel := "file:" + path

	if mismatch := verifyHashPin(content, expectedHash, sourceLabel); mismatch != nil {
		return mismatch
	}

	policy, parseErr := parsePolicy(content)
	if parseErr != nil {
		return &PolicyFetchResult{
			Err:     fmt.Sprintf("Invalid policy file %s: %v", path, parseErr),
			Source:  sourceLabel,
			Outcome: "malformed",
		}
	}

	outcome := "found"
	if isPolicyEmpty(policy) {
		outcome = "empty"
	}
	var rawHash string
	if expectedHash != "" {
		rawHash = computeHashNormalized(content, expectedHash)
	}
	return &PolicyFetchResult{
		Policy:       policy,
		Source:       sourceLabel,
		Outcome:      outcome,
		RawBytesHash: rawHash,
		ExpectedHash: expectedHash,
	}
}

// ---------------------------------------------------------------------------
// Auto-discovery
// ---------------------------------------------------------------------------

func autoDiscover(projectRoot string, noCache bool, expectedHash string) *PolicyFetchResult {
	org, host, err := extractOrgFromGitRemote(projectRoot)
	if err != nil || org == "" {
		return &PolicyFetchResult{
			Err:     "Could not determine org from git remote",
			Outcome: "no_git_remote",
		}
	}
	repoRef := org + "/.github"
	if host != "" && host != "github.com" {
		repoRef = host + "/" + repoRef
	}
	return fetchFromRepo(repoRef, projectRoot, noCache, expectedHash)
}

// extractOrgFromGitRemote runs git remote get-url origin and parses the org and host.
func extractOrgFromGitRemote(projectRoot string) (org, host string, err error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = projectRoot
	out, execErr := cmd.Output()
	if execErr != nil {
		return "", "", execErr
	}
	remoteURL := strings.TrimSpace(string(out))
	return parseRemoteURL(remoteURL)
}

// parseRemoteURL parses a git remote URL into (org, host, error).
func parseRemoteURL(rawURL string) (org, host string, err error) {
	if rawURL == "" {
		return "", "", fmt.Errorf("empty URL")
	}

	// SCP-style SSH: user@host:path
	if m := scpLikeRE.FindStringSubmatch(rawURL); len(m) > 0 {
		var hostPart, pathPart string
		for i, name := range scpLikeRE.SubexpNames() {
			switch name {
			case "host":
				hostPart = m[i]
			case "path":
				pathPart = m[i]
			}
		}
		pathPart = strings.TrimSuffix(strings.TrimRight(pathPart, "/"), ".git")
		parts := strings.Split(pathPart, "/")
		var cleaned []string
		for _, p := range parts {
			if p != "" {
				cleaned = append(cleaned, p)
			}
		}
		if len(cleaned) == 0 {
			return "", "", fmt.Errorf("cannot parse path from SCP URL")
		}
		// Azure DevOps SSH has v3/ prefix
		if hostPart == "ssh.dev.azure.com" && len(cleaned) >= 2 && cleaned[0] == "v3" {
			return cleaned[1], hostPart, nil
		}
		return cleaned[0], hostPart, nil
	}

	// HTTPS
	if strings.Contains(rawURL, "://") {
		u, parseErr := url.Parse(rawURL)
		if parseErr != nil {
			return "", "", parseErr
		}
		h := u.Hostname()
		pathPart := strings.TrimSuffix(strings.Trim(u.Path, "/"), ".git")
		parts := strings.Split(pathPart, "/")
		var cleaned []string
		for _, p := range parts {
			if p != "" {
				cleaned = append(cleaned, p)
			}
		}
		if h != "" && len(cleaned) > 0 {
			return cleaned[0], h, nil
		}
	}
	return "", "", fmt.Errorf("could not parse remote URL: %s", rawURL)
}

// ---------------------------------------------------------------------------
// URL fetch
// ---------------------------------------------------------------------------

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// Refuse redirects (security: prevent SSRF via redirect)
		return http.ErrUseLastResponse
	},
}

func fetchFromURL(rawURL, projectRoot string, noCache bool, expectedHash string) *PolicyFetchResult {
	sourceLabel := "url:" + rawURL
	var ce *cacheEntry

	if !noCache {
		ce = readCacheEntry(rawURL, projectRoot, defaultCacheTTL, expectedHash)
		if ce != nil && !ce.Stale {
			outcome := "found"
			if isPolicyEmpty(ce.Policy) {
				outcome = "empty"
			}
			return &PolicyFetchResult{
				Policy:          ce.Policy,
				Source:          ce.Source,
				Cached:          true,
				CacheAgeSeconds: ce.AgeSeconds,
				Outcome:         outcome,
				RawBytesHash:    ce.RawBytesHash,
				ExpectedHash:    expectedHash,
			}
		}
	}

	resp, err := httpClient.Get(rawURL)
	var content []byte
	var fetchErrStr string
	if err != nil {
		fetchErrStr = fmt.Sprintf("Error fetching %s: %v", rawURL, err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return &PolicyFetchResult{Source: sourceLabel, Err: "404: Policy file not found", Outcome: "absent"}
		}
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			loc := resp.Header.Get("Location")
			fetchErrStr = fmt.Sprintf("Refusing HTTP redirect (%d) from %s to %s", resp.StatusCode, rawURL, loc)
		} else if resp.StatusCode != 200 {
			fetchErrStr = fmt.Sprintf("HTTP %d fetching %s", resp.StatusCode, rawURL)
		} else {
			content, err = io.ReadAll(resp.Body)
			if err != nil {
				fetchErrStr = fmt.Sprintf("Error reading response from %s: %v", rawURL, err)
			}
		}
	}

	if fetchErrStr != "" {
		return staleOrError(ce, fetchErrStr, sourceLabel, "cache_miss_fetch_fail")
	}

	if gr := detectGarbage(content, rawURL, sourceLabel, ce); gr != nil {
		return gr
	}

	if mismatch := verifyHashPin(content, expectedHash, sourceLabel); mismatch != nil {
		return mismatch
	}

	policy, parseErr := parsePolicy(content)
	if parseErr != nil {
		return &PolicyFetchResult{
			Err:     fmt.Sprintf("Invalid policy from %s: %v", rawURL, parseErr),
			Source:  sourceLabel,
			Outcome: "malformed",
		}
	}

	actualHash := computeHashNormalized(content, expectedHash)
	writeCache(rawURL, policy, projectRoot, []string{rawURL}, actualHash)
	outcome := "found"
	if isPolicyEmpty(policy) {
		outcome = "empty"
	}
	return &PolicyFetchResult{
		Policy:       policy,
		Source:       sourceLabel,
		Outcome:      outcome,
		RawBytesHash: actualHash,
		ExpectedHash: expectedHash,
	}
}

// ---------------------------------------------------------------------------
// Repo fetch (GitHub Contents API)
// ---------------------------------------------------------------------------

func fetchFromRepo(repoRef, projectRoot string, noCache bool, expectedHash string) *PolicyFetchResult {
	sourceLabel := "org:" + repoRef
	var ce *cacheEntry

	if !noCache {
		ce = readCacheEntry(repoRef, projectRoot, defaultCacheTTL, expectedHash)
		if ce != nil && !ce.Stale {
			outcome := "found"
			if isPolicyEmpty(ce.Policy) {
				outcome = "empty"
			}
			return &PolicyFetchResult{
				Policy:          ce.Policy,
				Source:          ce.Source,
				Cached:          true,
				CacheAgeSeconds: ce.AgeSeconds,
				Outcome:         outcome,
				RawBytesHash:    ce.RawBytesHash,
				ExpectedHash:    expectedHash,
			}
		}
	}

	content, fetchErr := fetchGithubContents(repoRef, "apm-policy.yml")
	if fetchErr != "" {
		if strings.Contains(fetchErr, "404") {
			return &PolicyFetchResult{Source: sourceLabel, Outcome: "absent"}
		}
		return staleOrError(ce, fetchErr, sourceLabel, "cache_miss_fetch_fail")
	}
	if content == nil {
		return &PolicyFetchResult{Source: sourceLabel, Outcome: "absent"}
	}

	if gr := detectGarbage(content, repoRef, sourceLabel, ce); gr != nil {
		return gr
	}

	if mismatch := verifyHashPin(content, expectedHash, sourceLabel); mismatch != nil {
		return mismatch
	}

	policy, parseErr := parsePolicy(content)
	if parseErr != nil {
		return &PolicyFetchResult{
			Err:     fmt.Sprintf("Invalid policy in %s: %v", repoRef, parseErr),
			Source:  sourceLabel,
			Outcome: "malformed",
		}
	}

	actualHash := computeHashNormalized(content, expectedHash)
	writeCache(repoRef, policy, projectRoot, []string{repoRef}, actualHash)
	outcome := "found"
	if isPolicyEmpty(policy) {
		outcome = "empty"
	}
	return &PolicyFetchResult{
		Policy:       policy,
		Source:       sourceLabel,
		Outcome:      outcome,
		RawBytesHash: actualHash,
		ExpectedHash: expectedHash,
	}
}

// fetchGithubContents fetches apm-policy.yml from a GitHub/GHE repo via the Contents API.
// Returns (content, errString). One will be nil/"".
func fetchGithubContents(repoRef, filePath string) ([]byte, string) {
	parts := strings.Split(repoRef, "/")
	var host, owner, repo string
	switch len(parts) {
	case 2:
		host, owner, repo = "github.com", parts[0], parts[1]
	case 3:
		host, owner, repo = parts[0], parts[1], parts[2]
	default:
		if len(parts) >= 3 {
			host, owner, repo = parts[0], parts[1], strings.Join(parts[2:], "/")
		} else {
			return nil, fmt.Sprintf("Invalid repo reference: %s", repoRef)
		}
	}

	var apiURL string
	if host == "github.com" {
		apiURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, filePath)
	} else {
		apiURL = fmt.Sprintf("https://%s/api/v3/repos/%s/%s/contents/%s", host, owner, repo, filePath)
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Sprintf("Error building request for %s: %v", repoRef, err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := getTokenForHost(host); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Sprintf("Error fetching policy from %s: %v", repoRef, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return nil, "404: Policy file not found"
	case 403:
		return nil, fmt.Sprintf("403: Access denied to %s", repoRef)
	case 200:
		// continue
	default:
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			loc := resp.Header.Get("Location")
			return nil, fmt.Sprintf("Refusing HTTP redirect (%d) from %s to %s", resp.StatusCode, apiURL, loc)
		}
		return nil, fmt.Sprintf("HTTP %d fetching policy from %s", resp.StatusCode, repoRef)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Sprintf("Error reading response from %s: %v", repoRef, err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Sprintf("Error parsing response from %s: %v", repoRef, err)
	}

	if enc, ok := data["encoding"].(string); ok && enc == "base64" {
		if rawContent, ok := data["content"].(string); ok && rawContent != "" {
			cleaned := strings.ReplaceAll(rawContent, "\n", "")
			decoded, err := base64.StdEncoding.DecodeString(cleaned)
			if err != nil {
				return nil, fmt.Sprintf("Error decoding base64 content from %s: %v", repoRef, err)
			}
			return decoded, ""
		}
	}
	if rawContent, ok := data["content"].(string); ok && rawContent != "" {
		return []byte(rawContent), ""
	}
	return nil, fmt.Sprintf("Unexpected response format from %s", repoRef)
}

// getTokenForHost returns a GitHub/GHE token for the given host.
func getTokenForHost(host string) string {
	hostLower := strings.ToLower(host)
	isGitHub := hostLower == "github.com" || strings.HasSuffix(hostLower, ".ghe.com") ||
		(os.Getenv("GITHUB_HOST") != "" && hostLower == strings.ToLower(os.Getenv("GITHUB_HOST")))
	if !isGitHub {
		return ""
	}
	for _, env := range []string{"GITHUB_TOKEN", "GITHUB_APM_PAT", "GH_TOKEN"} {
		if t := os.Getenv(env); t != "" {
			return t
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Hash pin verification
// ---------------------------------------------------------------------------

// verifyHashPin verifies content against an expected hash pin.
// Returns nil when verification passes or there is no pin.
// Returns a PolicyFetchResult with outcome "hash_mismatch" on failure.
func verifyHashPin(content []byte, expectedHash, sourceLabel string) *PolicyFetchResult {
	if expectedHash == "" {
		return nil
	}
	algo, expectedHex, err := splitHashPin(expectedHash)
	if err != nil {
		return &PolicyFetchResult{
			Outcome:      "hash_mismatch",
			Source:       sourceLabel,
			Err:          fmt.Sprintf("Policy hash mismatch from %s: invalid pin (%v)", sourceLabel, err),
			ExpectedHash: expectedHash,
		}
	}

	var actualHex string
	switch algo {
	case "sha256":
		h := sha256.Sum256(content)
		actualHex = fmt.Sprintf("%x", h)
	default:
		return &PolicyFetchResult{
			Outcome: "hash_mismatch",
			Source:  sourceLabel,
			Err:     fmt.Sprintf("Unsupported hash algorithm: %s", algo),
		}
	}

	if actualHex != expectedHex {
		return &PolicyFetchResult{
			Outcome:      "hash_mismatch",
			Source:       sourceLabel,
			Err:          fmt.Sprintf("Policy hash mismatch from %s: expected %s:%s, got %s:%s", sourceLabel, algo, expectedHex, algo, actualHex),
			ExpectedHash: fmt.Sprintf("%s:%s", algo, expectedHex),
			RawBytesHash: fmt.Sprintf("%s:%s", algo, actualHex),
		}
	}
	return nil
}

// splitHashPin splits "<algo>:<hex>" into (algo, hex).
// Bare hex without prefix is treated as sha256 for backward compatibility.
func splitHashPin(pin string) (algo, hex string, err error) {
	raw := strings.TrimSpace(pin)
	if strings.Contains(raw, ":") {
		idx := strings.Index(raw, ":")
		algo = strings.ToLower(strings.TrimSpace(raw[:idx]))
		hex = strings.ToLower(strings.TrimSpace(raw[idx+1:]))
	} else {
		algo = "sha256"
		hex = strings.ToLower(raw)
	}
	if algo != "sha256" {
		return "", "", fmt.Errorf("unsupported algorithm %q", algo)
	}
	if len(hex) != 64 {
		return "", "", fmt.Errorf("invalid sha256 hex (length %d)", len(hex))
	}
	return algo, hex, nil
}

func computeHashNormalized(content []byte, expectedHash string) string {
	algo := "sha256"
	if expectedHash != "" {
		if a, _, err := splitHashPin(expectedHash); err == nil {
			algo = a
		}
	}
	switch algo {
	case "sha256":
		h := sha256.Sum256(content)
		return fmt.Sprintf("sha256:%x", h)
	}
	return ""
}

// ---------------------------------------------------------------------------
// Policy parsing
// ---------------------------------------------------------------------------

// parsePolicy parses raw YAML bytes into an ApmPolicy.
// Uses a minimal line-by-line scanner tracking current section context.
func parsePolicy(data []byte) (*schema.ApmPolicy, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return &schema.ApmPolicy{}, nil
	}

	p := &schema.ApmPolicy{}
	lines := strings.Split(string(data), "\n")

	// Track section by top-level key and sub-key
	var section, subSection, listKey string
	var listTarget *[]string

	setListTarget := func(key string) {
		switch {
		case section == "dependencies" && key == "allow":
			listTarget = &p.Deps.Allow
		case section == "dependencies" && key == "deny":
			listTarget = &p.Deps.Deny
		case section == "dependencies" && key == "require":
			listTarget = &p.Deps.Require
		case section == "mcp" && key == "allow":
			listTarget = &p.MCP.Allow
		case section == "mcp" && key == "deny":
			listTarget = &p.MCP.Deny
		case section == "mcp" && subSection == "transport" && key == "allow":
			listTarget = &p.MCP.Transport.Allow
		case section == "compilation" && subSection == "target" && key == "allow":
			listTarget = &p.Compilation.Targets.Allow
		default:
			listTarget = nil
		}
		listKey = key
		_ = listKey
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := 0
		for _, ch := range line {
			if ch == ' ' {
				indent++
			} else {
				break
			}
		}

		if strings.HasPrefix(trimmed, "- ") {
			val := strings.TrimPrefix(trimmed, "- ")
			val = strings.Trim(val, "\"'")
			if listTarget != nil {
				*listTarget = append(*listTarget, val)
			}
			continue
		}

		if idx := strings.Index(trimmed, ":"); idx >= 0 {
			key := strings.TrimSpace(trimmed[:idx])
			val := strings.TrimSpace(trimmed[idx+1:])
			val = strings.Trim(val, "\"'")

			if indent == 0 {
				// Top-level key
				section = key
				subSection = ""
				listTarget = nil
				switch key {
				case "version":
					p.Version = val
				case "enforcement":
					p.Enforcement = val
				case "fetch_failure":
					p.FetchFailure = val
				}
			} else if indent == 2 {
				// Section key
				subSection = ""
				listTarget = nil
				if val == "" {
					subSection = key
				} else {
					switch {
					case section == "dependencies" && key == "require_resolution":
						p.Deps.RequireResolution = val
					case section == "mcp" && key == "self_defined":
						p.MCP.SelfDefined = val
					case section == "compilation" && key == "source_attribution":
						// ignore
					}
					setListTarget(key)
					// If val is empty this is a list parent -- handled above
					// If non-empty, clear listTarget (it's a scalar, not list)
					if val != "" {
						listTarget = nil
					}
				}
			} else if indent == 4 {
				// Sub-section key
				listTarget = nil
				if val == "" {
					subSection = key
				} else {
					switch {
					case section == "mcp" && subSection == "transport" && key == "allow":
						// scalar allow -- no-op
					case section == "compilation" && subSection == "target" && key == "enforce":
						p.Compilation.Targets.Enforce = val
					case section == "compilation" && subSection == "strategy" && key == "enforce":
						p.Compilation.Strategy.Enforce = val
					}
					setListTarget(key)
					if val != "" {
						listTarget = nil
					}
				}
			}
		}
	}

	return p, nil
}

// isPolicyEmpty returns true when a policy has no actionable restrictions.
func isPolicyEmpty(p *schema.ApmPolicy) bool {
	if p == nil {
		return true
	}
	return len(p.Deps.Deny) == 0 &&
		p.Deps.Allow == nil &&
		len(p.Deps.Require) == 0 &&
		len(p.MCP.Deny) == 0 &&
		p.MCP.Allow == nil &&
		p.MCP.Transport.Allow == nil &&
		p.Compilation.Targets.Allow == nil
}


// ---------------------------------------------------------------------------
// Cache
// ---------------------------------------------------------------------------

type cacheMeta struct {
	RepoRef       string   `json:"repo_ref"`
	CachedAt      float64  `json:"cached_at"`
	ChainRefs     []string `json:"chain_refs"`
	SchemaVersion string   `json:"schema_version"`
	Fingerprint   string   `json:"fingerprint"`
	RawBytesHash  string   `json:"raw_bytes_hash"`
}

func cacheKey(repoRef string) string {
	h := sha256.Sum256([]byte(repoRef))
	return fmt.Sprintf("%x", h)[:16]
}

func getCacheDir(projectRoot string) (string, error) {
	resolved, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", err
	}
	base := filepath.Join(resolved, "apm_modules")
	candidate := filepath.Join(base, policyCacheDir)
	if _, err := pathsecurity.EnsurePathWithin(candidate, resolved); err != nil {
		return "", fmt.Errorf("policy cache path %q resolves outside project root %q", candidate, resolved)
	}
	return candidate, nil
}

func readCacheEntry(repoRef, projectRoot string, ttl int, expectedHash string) *cacheEntry {
	cacheDir, err := getCacheDir(projectRoot)
	if err != nil {
		return nil
	}
	key := cacheKey(repoRef)
	policyFile := filepath.Join(cacheDir, key+".yml")
	metaFile := filepath.Join(cacheDir, key+".meta.json")

	if _, err := os.Stat(policyFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(metaFile); os.IsNotExist(err) {
		return nil
	}

	metaBytes, err := os.ReadFile(metaFile)
	if err != nil {
		return nil
	}
	var meta cacheMeta
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return nil
	}
	if meta.SchemaVersion != cacheSchemaVersion {
		return nil
	}

	age := int(time.Now().Unix() - int64(meta.CachedAt))
	if age > maxStaleTTL {
		return nil
	}

	// Pin verification
	if expectedHash != "" {
		ea, eh, err := splitHashPin(expectedHash)
		if err != nil {
			return nil
		}
		expectedNorm := fmt.Sprintf("%s:%s", ea, eh)
		if strings.ToLower(meta.RawBytesHash) != expectedNorm {
			return nil
		}
	}

	policyContent, err := os.ReadFile(policyFile)
	if err != nil {
		return nil
	}
	policy, err := parsePolicy(policyContent)
	if err != nil {
		return nil
	}

	source := "org:" + repoRef
	if strings.HasPrefix(repoRef, "http://") || strings.HasPrefix(repoRef, "https://") {
		source = "url:" + repoRef
	}

	return &cacheEntry{
		Policy:       policy,
		Source:       source,
		AgeSeconds:   age,
		Stale:        age > ttl,
		ChainRefs:    meta.ChainRefs,
		Fingerprint:  meta.Fingerprint,
		RawBytesHash: meta.RawBytesHash,
	}
}

var writeMu sync.Mutex

func writeCache(repoRef string, policy *schema.ApmPolicy, projectRoot string, chainRefs []string, rawBytesHash string) {
	cacheDir, err := getCacheDir(projectRoot)
	if err != nil {
		return
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return
	}

	key := cacheKey(repoRef)
	policyFile := filepath.Join(cacheDir, key+".yml")
	metaFile := filepath.Join(cacheDir, key+".meta.json")

	serialized := serializePolicy(policy)
	fingerprint := fmt.Sprintf("%x", sha256.Sum256([]byte(serialized)))[:32]

	meta := cacheMeta{
		RepoRef:       repoRef,
		CachedAt:      float64(time.Now().UnixNano()) / 1e9,
		ChainRefs:     chainRefs,
		SchemaVersion: cacheSchemaVersion,
		Fingerprint:   fingerprint,
		RawBytesHash:  rawBytesHash,
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return
	}

	writeMu.Lock()
	defer writeMu.Unlock()

	uid := fmt.Sprintf("%d", time.Now().UnixNano())
	tmpPolicy := policyFile + "." + uid + ".tmp"
	if err := os.WriteFile(tmpPolicy, []byte(serialized), 0o644); err == nil {
		_ = os.Rename(tmpPolicy, policyFile)
	}
	tmpMeta := metaFile + "." + uid + ".tmp"
	if err := os.WriteFile(tmpMeta, metaBytes, 0o644); err == nil {
		_ = os.Rename(tmpMeta, metaFile)
	}
}

// serializePolicy serializes an ApmPolicy to a simple YAML-like string for caching.
func serializePolicy(p *schema.ApmPolicy) string {
	if p == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("version: %s\n", p.Version))
	sb.WriteString(fmt.Sprintf("enforcement: %s\n", p.Enforcement))
	sb.WriteString(fmt.Sprintf("fetch_failure: %s\n", p.FetchFailure))
	if len(p.Deps.Deny) > 0 {
		sb.WriteString("dependencies:\n")
		sb.WriteString("  deny:\n")
		for _, d := range p.Deps.Deny {
			sb.WriteString("  - " + d + "\n")
		}
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// Garbage detection
// ---------------------------------------------------------------------------

func detectGarbage(content []byte, identifier, sourceLabel string, ce *cacheEntry) *PolicyFetchResult {
	if content == nil {
		return nil
	}
	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		return nil
	}
	// Very basic check: a valid YAML policy starts with a known key or is a mapping
	// For garbage detection: if it starts with "<" (HTML) it's a captive portal
	if strings.HasPrefix(trimmed, "<") {
		msg := fmt.Sprintf("Response from %s is not valid YAML (possible captive portal or redirect)", identifier)
		if ce != nil {
			return &PolicyFetchResult{
				Policy:          ce.Policy,
				Source:          ce.Source,
				Cached:          true,
				CacheStale:      true,
				CacheAgeSeconds: ce.AgeSeconds,
				FetchErr:        msg,
				Outcome:         "cached_stale",
			}
		}
		return &PolicyFetchResult{
			Err:     msg,
			Source:  sourceLabel,
			FetchErr: msg,
			Outcome:  "garbage_response",
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Stale or error fallback
// ---------------------------------------------------------------------------

func staleOrError(ce *cacheEntry, fetchErrMsg, sourceLabel, outcomeOnMiss string) *PolicyFetchResult {
	if ce != nil {
		return &PolicyFetchResult{
			Policy:          ce.Policy,
			Source:          ce.Source,
			Cached:          true,
			CacheStale:      true,
			CacheAgeSeconds: ce.AgeSeconds,
			FetchErr:        fetchErrMsg,
			Outcome:         "cached_stale",
		}
	}
	return &PolicyFetchResult{
		Err:     fetchErrMsg,
		Source:  sourceLabel,
		FetchErr: fetchErrMsg,
		Outcome:  outcomeOnMiss,
	}
}

// readProjectHashPin is a stub -- returns "" if no apm.yml hash pin found.
func readProjectHashPin(projectRoot string) string {
	// Full implementation would parse apm.yml policy.hash field.
	// Returning "" for now -- callers pass the pin explicitly when available.
	return ""
}
