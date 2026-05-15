// Package gitlabresolver resolves GitLab direct-shorthand package specs.
// Migrated from src/apm_cli/install/gitlab_resolver.py
package gitlabresolver

import "strings"

// GitLabDirectShorthandUnresolved is the error message when shorthand probing fails.
const GitLabDirectShorthandUnresolved = "Direct GitLab host/path did not resolve to a reachable " +
	"repository with an installable package path. Use an explicit 'git' URL with a 'path' field " +
	"for a deeper project or subdirectory."

// ShorthandParts holds the parsed pieces of a GitLab direct shorthand spec.
type ShorthandParts struct {
	Host     string
	Segments []string
	Ref      string
}

// ParseShorthand splits a GitLab host/path shorthand into its components.
// Returns nil when the input does not look like a GitLab shorthand.
func ParseShorthand(pkg string) *ShorthandParts {
	// Expected form: host/seg1/seg2[...][#ref]
	ref := ""
	if idx := strings.LastIndex(pkg, "#"); idx >= 0 {
		ref = pkg[idx+1:]
		pkg = pkg[:idx]
	}
	parts := strings.SplitN(pkg, "/", 2)
	if len(parts) < 2 {
		return nil
	}
	host := parts[0]
	// Must contain a dot to be a hostname
	if !strings.Contains(host, ".") {
		return nil
	}
	segments := strings.Split(parts[1], "/")
	if len(segments) == 0 {
		return nil
	}
	return &ShorthandParts{Host: host, Segments: segments, Ref: ref}
}

// BoundaryCandidates iterates candidate repo/virtualPath splits for a segment list.
// It yields candidates from longest to shortest repo path (greedy first).
type BoundaryCandidates struct {
	Host     string
	Segments []string
	Ref      string
	idx      int
}

// NewBoundaryCandidates creates an iterator over boundary candidates.
func NewBoundaryCandidates(parts *ShorthandParts) *BoundaryCandidates {
	return &BoundaryCandidates{
		Host:     parts.Host,
		Segments: parts.Segments,
		Ref:      parts.Ref,
		// Start from the longest possible repo path (need at least 2 segments for owner/repo)
		idx: len(parts.Segments),
	}
}

// BoundaryCandidate is one candidate repo/virtualPath split.
type BoundaryCandidate struct {
	RepoPath    string // "owner/repo"
	VirtualPath string // sub-path within the repo, or ""
}

// Next returns the next candidate, advancing the iterator.
// Returns (zero, false) when exhausted.
func (b *BoundaryCandidates) Next() (BoundaryCandidate, bool) {
	if b.idx < 2 {
		return BoundaryCandidate{}, false
	}
	repoSegs := b.Segments[:b.idx]
	virtualSegs := b.Segments[b.idx:]
	b.idx--
	return BoundaryCandidate{
		RepoPath:    strings.Join(repoSegs, "/"),
		VirtualPath: strings.Join(virtualSegs, "/"),
	}, true
}

// IsGitLabHost reports whether host looks like a GitLab instance (not GitHub/ADO).
func IsGitLabHost(host string) bool {
	h := strings.ToLower(host)
	if h == "github.com" || strings.HasSuffix(h, ".ghe.com") {
		return false
	}
	if h == "dev.azure.com" || strings.HasSuffix(h, ".visualstudio.com") {
		return false
	}
	return true
}
