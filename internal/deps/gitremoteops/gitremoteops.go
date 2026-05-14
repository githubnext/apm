// Package gitremoteops provides helpers for parsing git remote references.
package gitremoteops

import (
"regexp"
"sort"
"strings"
)

// GitReferenceType identifies the kind of a git reference.
type GitReferenceType int

const (
GitRefBranch GitReferenceType = iota
GitRefTag
)

// RemoteRef is a single remote git reference with its commit SHA.
type RemoteRef struct {
Name      string
RefType   GitReferenceType
CommitSHA string
}

var semverTagRe = regexp.MustCompile(`^v?\d+\.\d+\.\d+`)

// ParseLsRemoteOutput parses "git ls-remote --tags --heads" output.
func ParseLsRemoteOutput(output string) []RemoteRef {
tags := map[string]string{} // name -> commit sha
var branches []RemoteRef

for _, line := range strings.Split(output, "\n") {
line = strings.TrimSpace(line)
if line == "" {
continue
}
parts := strings.SplitN(line, "\t", 2)
if len(parts) != 2 {
continue
}
sha := strings.TrimSpace(parts[0])
refname := strings.TrimSpace(parts[1])

switch {
case strings.HasPrefix(refname, "refs/tags/"):
tagName := refname[len("refs/tags/"):]
if strings.HasSuffix(tagName, "^{}") {
tags[tagName[:len(tagName)-3]] = sha
} else {
if _, ok := tags[tagName]; !ok {
tags[tagName] = sha
}
}
case strings.HasPrefix(refname, "refs/heads/"):
branchName := refname[len("refs/heads/"):]
branches = append(branches, RemoteRef{Name: branchName, RefType: GitRefBranch, CommitSHA: sha})
}
}

var refs []RemoteRef
for name, sha := range tags {
refs = append(refs, RemoteRef{Name: name, RefType: GitRefTag, CommitSHA: sha})
}
refs = append(refs, branches...)
return refs
}

// SortRefsBySemver sorts tag refs by semantic version (descending), non-semver tags last.
func SortRefsBySemver(refs []RemoteRef) []RemoteRef {
sorted := make([]RemoteRef, len(refs))
copy(sorted, refs)
sort.Slice(sorted, func(i, j int) bool {
ai := semverTagRe.MatchString(sorted[i].Name)
aj := semverTagRe.MatchString(sorted[j].Name)
if ai != aj {
return ai
}
return sorted[i].Name > sorted[j].Name
})
return sorted
}
