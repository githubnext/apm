// Package dependency defines dependency reference types for APM packages.
// Mirrors src/apm_cli/models/dependency/types.py.
package dependency

import "regexp"

// GitReferenceType classifies a git reference as branch, tag, or commit.
type GitReferenceType int

const (
	GitRefBranch GitReferenceType = iota
	GitRefTag
	GitRefCommit
)

func (t GitReferenceType) String() string {
	switch t {
	case GitRefBranch:
		return "branch"
	case GitRefTag:
		return "tag"
	case GitRefCommit:
		return "commit"
	default:
		return "unknown"
	}
}

// VirtualPackageType classifies a virtual package as a file or subdirectory.
type VirtualPackageType int

const (
	VirtualPackageFile VirtualPackageType = iota
	VirtualPackageSubdirectory
)

func (t VirtualPackageType) String() string {
	switch t {
	case VirtualPackageFile:
		return "file"
	case VirtualPackageSubdirectory:
		return "subdirectory"
	default:
		return "unknown"
	}
}

// RemoteRef represents a single remote git reference with its commit SHA.
type RemoteRef struct {
	Name      string
	RefType   GitReferenceType
	CommitSHA string
}

// ResolvedReference represents a resolved git reference.
type ResolvedReference struct {
	OriginalRef     string
	RefType         GitReferenceType
	ResolvedCommit  string
	RefName         string
}

func (r ResolvedReference) String() string {
	if r.ResolvedCommit == "" {
		return r.RefName
	}
	if r.RefType == GitRefCommit {
		if len(r.ResolvedCommit) > 8 {
			return r.ResolvedCommit[:8]
		}
		return r.ResolvedCommit
	}
	commit := r.ResolvedCommit
	if len(commit) > 8 {
		commit = commit[:8]
	}
	return r.RefName + " (" + commit + ")"
}

var (
	commitSHARE = regexp.MustCompile(`^[a-f0-9]{7,40}$`)
	semverRE    = regexp.MustCompile(`^v?\d+\.\d+\.\d+`)
)

// ParseGitReference parses a git reference string to determine its type.
// Mirrors src/apm_cli/models/dependency/types.py:parse_git_reference.
func ParseGitReference(ref string) (GitReferenceType, string) {
	if ref == "" {
		return GitRefBranch, "main"
	}

	// Check for commit SHA (7-40 hex chars)
	if commitSHARE.MatchString(ref) {
		return GitRefCommit, ref
	}

	// Check for semantic version tag
	if semverRE.MatchString(ref) {
		return GitRefTag, ref
	}

	return GitRefBranch, ref
}
