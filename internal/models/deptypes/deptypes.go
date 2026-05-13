// Package deptypes defines dependency type enums and dataclasses.
package deptypes

import "regexp"

// GitReferenceType represents the type of a git reference.
type GitReferenceType int

const (
GitRefBranch GitReferenceType = iota
GitRefTag
GitRefCommit
)

// RemoteRef is a single remote git reference with its commit SHA.
type RemoteRef struct {
Name      string
RefType   GitReferenceType
CommitSHA string
}

// VirtualPackageType is the type of a virtual package.
type VirtualPackageType int

const (
VirtualPackageFile VirtualPackageType = iota
VirtualPackageSubdirectory
)

// ResolvedReference represents a resolved git reference.
type ResolvedReference struct {
OriginalRef     string
RefType         GitReferenceType
ResolvedCommit  string
RefName         string
}

var commitRe = regexp.MustCompile(`^[a-f0-9]{7,40}$`)
var semverRe = regexp.MustCompile(`^v?\d+\.\d+\.\d+`)

// ParseGitReference parses a git reference string to determine its type.
func ParseGitReference(ref string) (GitReferenceType, string) {
if ref == "" {
return GitRefBranch, "main"
}
r := ref
if commitRe.MatchString(r) {
return GitRefCommit, r
}
if semverRe.MatchString(r) {
return GitRefTag, r
}
return GitRefBranch, r
}
