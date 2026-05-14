// Package refresolver provides concurrent git ls-remote with in-memory ref caching.
// Migrated from src/apm_cli/marketplace/ref_resolver.py.
package refresolver

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/githubnext/apm/internal/marketplace/gitstderr"
	"github.com/githubnext/apm/internal/marketplace/gitutils"
	"github.com/githubnext/apm/internal/utils/githubhost"
)

// RemoteRef is a single ref returned by git ls-remote.
type RemoteRef struct {
	Name string // e.g. "refs/tags/v1.2.0" or "refs/heads/main"
	SHA  string // 40-char hex SHA
}

var shaRE = regexp.MustCompile(`^[0-9a-f]{40}$`)

// DefaultTTL is the default cache TTL (5 minutes).
const DefaultTTL = 5 * time.Minute

type cacheEntry struct {
	refs      []RemoteRef
	timestamp time.Time
}

// RefCache is an in-memory cache keyed on "owner/repo".
type RefCache struct {
	mu    sync.Mutex
	store map[string]*cacheEntry
	ttl   time.Duration
}

// NewRefCache creates a RefCache with the given TTL.
func NewRefCache(ttl time.Duration) *RefCache {
	return &RefCache{store: make(map[string]*cacheEntry), ttl: ttl}
}

// Get returns cached refs or nil on miss/expiry.
func (c *RefCache) Get(ownerRepo string) []RemoteRef {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.store[ownerRepo]
	if e == nil {
		return nil
	}
	if time.Since(e.timestamp) > c.ttl {
		delete(c.store, ownerRepo)
		return nil
	}
	out := make([]RemoteRef, len(e.refs))
	copy(out, e.refs)
	return out
}

// Put stores refs for ownerRepo.
func (c *RefCache) Put(ownerRepo string, refs []RemoteRef) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := make([]RemoteRef, len(refs))
	copy(cp, refs)
	c.store[ownerRepo] = &cacheEntry{refs: cp, timestamp: time.Now()}
}

// Clear drops all entries.
func (c *RefCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]*cacheEntry)
}

// Len returns the number of cached entries.
func (c *RefCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.store)
}

// GitLsRemoteError is raised when git ls-remote fails.
type GitLsRemoteError struct {
	Package string
	Summary string
	Hint    string
}

func (e *GitLsRemoteError) Error() string {
	if e.Hint != "" {
		return e.Summary + " " + e.Hint
	}
	return e.Summary
}

// OfflineMissError is raised in offline mode when the cache has no entry.
type OfflineMissError struct {
	Package string
	Remote  string
}

func (e *OfflineMissError) Error() string {
	return fmt.Sprintf("offline mode: no cached refs for remote '%s'", e.Remote)
}

// RefResolver runs git ls-remote and caches the results.
type RefResolver struct {
	timeoutSeconds float64
	offline        bool
	host           string
	token          string
	cache          *RefCache
	mu             sync.Mutex
	remoteLocks    map[string]*sync.Mutex
}

// New creates a RefResolver.
func New(timeoutSeconds float64, offline bool, host, token string) *RefResolver {
	if host == "" {
		host = githubhost.DefaultHost()
	}
	if host == "" {
		host = "github.com"
	}
	return &RefResolver{
		timeoutSeconds: timeoutSeconds,
		offline:        offline,
		host:           host,
		token:          token,
		cache:          NewRefCache(DefaultTTL),
		remoteLocks:    make(map[string]*sync.Mutex),
	}
}

func (r *RefResolver) remoteLock(ownerRepo string) *sync.Mutex {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.remoteLocks[ownerRepo]; !ok {
		r.remoteLocks[ownerRepo] = &sync.Mutex{}
	}
	return r.remoteLocks[ownerRepo]
}

// buildHTTPSCloneURL constructs an authenticated HTTPS clone URL.
func buildHTTPSCloneURL(host, ownerRepo, token string) string {
	base := fmt.Sprintf("https://%s/%s.git", host, ownerRepo)
	if token != "" {
		base = fmt.Sprintf("https://x-access-token:%s@%s/%s.git", token, host, ownerRepo)
	}
	return base
}

// parseLsRemoteOutput parses git ls-remote stdout into RemoteRefs.
func parseLsRemoteOutput(output string) []RemoteRef {
	var refs []RemoteRef
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		sha, refname := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if !shaRE.MatchString(sha) {
			continue
		}
		if strings.HasSuffix(refname, "^{}") {
			continue
		}
		refs = append(refs, RemoteRef{Name: refname, SHA: sha})
	}
	return refs
}

// ListRemoteRefs fetches all tags and heads from the configured Git host.
func (r *RefResolver) ListRemoteRefs(ownerRepo string) ([]RemoteRef, error) {
	lock := r.remoteLock(ownerRepo)
	lock.Lock()
	defer lock.Unlock()

	if cached := r.cache.Get(ownerRepo); cached != nil {
		return cached, nil
	}

	if r.offline {
		return nil, &OfflineMissError{Remote: ownerRepo}
	}

	url := buildHTTPSCloneURL(r.host, ownerRepo, r.token)
	timeout := time.Duration(r.timeoutSeconds * float64(time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--tags", "--heads", url)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_ASKPASS=echo")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, &GitLsRemoteError{
			Summary: fmt.Sprintf("git ls-remote timed out after %.0fs for '%s'.", r.timeoutSeconds, ownerRepo),
			Hint:    "Increase --timeout or check your network connection.",
		}
	}
	if runErr != nil {
		stderrStr := gitutils.RedactToken(stderr.String())
		exitCode := -1
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		translated := gitstderr.Translate(stderrStr, gitstderr.Options{
			ExitCode:  &exitCode,
			Operation: "ls-remote",
			Remote:    ownerRepo,
		})
		return nil, &GitLsRemoteError{
			Summary: translated.Summary,
			Hint:    translated.Hint,
		}
	}

	refs := parseLsRemoteOutput(stdout.String())
	r.cache.Put(ownerRepo, refs)
	return refs, nil
}

// ResolveRefSHA resolves a single ref to its concrete SHA.
func (r *RefResolver) ResolveRefSHA(ownerRepo, ref string) (string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	url := buildHTTPSCloneURL(r.host, ownerRepo, r.token)
	timeout := time.Duration(r.timeoutSeconds * float64(time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "ls-remote", url, ref)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_ASKPASS=echo")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", &GitLsRemoteError{
			Summary: fmt.Sprintf("git ls-remote timed out after %.0fs for '%s'.", r.timeoutSeconds, ownerRepo),
			Hint:    "Increase --timeout or check your network connection.",
		}
	}
	if runErr != nil {
		stderrStr := gitutils.RedactToken(stderr.String())
		exitCode := -1
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		translated := gitstderr.Translate(stderrStr, gitstderr.Options{
			ExitCode:  &exitCode,
			Operation: "ls-remote",
			Remote:    ownerRepo,
		})
		return "", &GitLsRemoteError{
			Summary: translated.Summary,
			Hint:    translated.Hint,
		}
	}

	refs := parseLsRemoteOutput(stdout.String())
	if len(refs) == 0 {
		return "", &GitLsRemoteError{
			Summary: fmt.Sprintf("Ref '%s' not found on remote '%s'.", ref, ownerRepo),
			Hint:    "Check that the ref exists and you have access to the repository.",
		}
	}
	return refs[0].SHA, nil
}

// Close releases resources.
func (r *RefResolver) Close() {
	r.cache.Clear()
	r.mu.Lock()
	r.remoteLocks = make(map[string]*sync.Mutex)
	r.mu.Unlock()
}
