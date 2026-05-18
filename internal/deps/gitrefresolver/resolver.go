// Package gitrefresolver resolves git references to concrete SHAs.
// Migrated from src/apm_cli/deps/git_reference_resolver.py
package gitrefresolver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// GitReferenceType indicates the kind of a resolved git reference.
type GitReferenceType int

const (
	ReferenceTypeBranch GitReferenceType = iota
	ReferenceTypeTag
	ReferenceTypeCommit
	ReferenceTypeUnknown
)

// RemoteRef represents a single git ref returned by ls-remote.
type RemoteRef struct {
	Name   string
	SHA    string
	IsTag  bool
	IsBranch bool
}

// ResolvedReference is the output of reference resolution.
type ResolvedReference struct {
	SHA     string
	RefType GitReferenceType
	Ref     string
}

// GitHubAPIResult holds a resolved SHA from the GitHub commits API.
type GitHubAPIResult struct {
	SHA string
}

// fullSHARe matches a 40-hex-char full SHA.
var fullSHARe = regexp.MustCompile(`^[0-9a-f]{40}$`)

// shortSHARe matches a 7-40-hex-char short SHA.
var shortSHARe = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

// GitReferenceResolver resolves user-supplied refs to concrete SHAs.
type GitReferenceResolver struct {
	AuthToken string
	Host      string
	Timeout   time.Duration
}

// New creates a GitReferenceResolver.
func New(host, authToken string) *GitReferenceResolver {
	return &GitReferenceResolver{
		Host:      host,
		AuthToken: authToken,
		Timeout:   15 * time.Second,
	}
}

// IsFullSHA reports whether ref looks like a 40-hex-char commit SHA.
func IsFullSHA(ref string) bool {
	return fullSHARe.MatchString(ref)
}

// IsShortSHA reports whether ref could be a short SHA (7-40 hex chars).
func IsShortSHA(ref string) bool {
	return shortSHARe.MatchString(ref)
}

// ResolveViaGitHubAPI attempts a cheap SHA lookup via the GitHub commits API.
// Returns ("", false, nil) when the fast path is not applicable.
func (r *GitReferenceResolver) ResolveViaGitHubAPI(owner, repo, ref string) (string, bool, error) {
	if r.Host != "github.com" && !strings.HasSuffix(r.Host, ".ghe.com") {
		return "", false, nil
	}
	url := fmt.Sprintf("https://api.%s/repos/%s/%s/commits/%s", r.Host, owner, repo, ref)
	if r.Host == "github.com" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, ref)
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Accept", "application/vnd.github.sha")
	if r.AuthToken != "" {
		req.Header.Set("Authorization", "token "+r.AuthToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, nil
	}

	buf := make([]byte, 41)
	n, _ := resp.Body.Read(buf)
	sha := strings.TrimSpace(string(buf[:n]))
	if IsFullSHA(sha) {
		return sha, true, nil
	}
	return "", false, nil
}

// ParseLsRemoteOutput parses the output of git ls-remote into RemoteRef slices.
func ParseLsRemoteOutput(output string) []RemoteRef {
	var refs []RemoteRef
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		sha := parts[0]
		refName := parts[1]
		// Skip peeled tags
		if strings.HasSuffix(refName, "^{}") {
			continue
		}
		rr := RemoteRef{SHA: sha, Name: refName}
		if strings.HasPrefix(refName, "refs/tags/") {
			rr.IsTag = true
			rr.Name = strings.TrimPrefix(refName, "refs/tags/")
		} else if strings.HasPrefix(refName, "refs/heads/") {
			rr.IsBranch = true
			rr.Name = strings.TrimPrefix(refName, "refs/heads/")
		}
		refs = append(refs, rr)
	}
	return refs
}

// ListRemoteRefs runs git ls-remote against a remote URL.
func ListRemoteRefs(repoURL string, extraEnv []string) ([]RemoteRef, error) {
	cmd := exec.Command("git", "ls-remote", "--tags", "--heads", repoURL)
	cmd.Env = append(os.Environ(), extraEnv...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ls-remote failed: %w", err)
	}
	return ParseLsRemoteOutput(string(out)), nil
}

// FindRef searches a list of RemoteRefs for an exact match by short name.
func FindRef(refs []RemoteRef, name string) (RemoteRef, bool) {
	for _, r := range refs {
		if r.Name == name {
			return r, true
		}
	}
	return RemoteRef{}, false
}

// ClassifyRef determines the GitReferenceType for a raw ref string.
func ClassifyRef(refs []RemoteRef, rawRef string) GitReferenceType {
	if IsFullSHA(rawRef) {
		return ReferenceTypeCommit
	}
	for _, r := range refs {
		if r.Name == rawRef {
			if r.IsTag {
				return ReferenceTypeTag
			}
			if r.IsBranch {
				return ReferenceTypeBranch
			}
		}
	}
	if IsShortSHA(rawRef) {
		return ReferenceTypeCommit
	}
	return ReferenceTypeUnknown
}

// Resolve attempts to resolve a ref for owner/repo, trying the GitHub API first.
func (r *GitReferenceResolver) Resolve(owner, repo, ref string) (*ResolvedReference, error) {
	if IsFullSHA(ref) {
		return &ResolvedReference{SHA: ref, Ref: ref, RefType: ReferenceTypeCommit}, nil
	}

	// Try GitHub API fast path
	if sha, ok, err := r.ResolveViaGitHubAPI(owner, repo, ref); err == nil && ok {
		refType := ReferenceTypeBranch
		if IsFullSHA(sha) && sha == ref {
			refType = ReferenceTypeCommit
		}
		return &ResolvedReference{SHA: sha, Ref: ref, RefType: refType}, nil
	}

	return nil, fmt.Errorf("could not resolve ref %q for %s/%s", ref, owner, repo)
}
