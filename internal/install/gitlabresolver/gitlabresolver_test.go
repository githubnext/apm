package gitlabresolver_test

import (
"testing"

"github.com/githubnext/apm/internal/install/gitlabresolver"
)

func TestParseShorthand(t *testing.T) {
cases := []struct {
input    string
wantNil  bool
wantHost string
wantSegs []string
wantRef  string
}{
{"gitlab.com/owner/repo", false, "gitlab.com", []string{"owner", "repo"}, ""},
{"gitlab.com/owner/repo#v1.0", false, "gitlab.com", []string{"owner", "repo"}, "v1.0"},
{"gitlab.com/owner/repo/subdir", false, "gitlab.com", []string{"owner", "repo", "subdir"}, ""},
{"notahost/foo", true, "", nil, ""},
{"singlepart", true, "", nil, ""},
{"", true, "", nil, ""},
}
for _, c := range cases {
got := gitlabresolver.ParseShorthand(c.input)
if c.wantNil {
if got != nil {
t.Errorf("ParseShorthand(%q): want nil, got %+v", c.input, got)
}
continue
}
if got == nil {
t.Errorf("ParseShorthand(%q): want non-nil", c.input)
continue
}
if got.Host != c.wantHost {
t.Errorf("ParseShorthand(%q) Host: want %q, got %q", c.input, c.wantHost, got.Host)
}
if len(got.Segments) != len(c.wantSegs) {
t.Errorf("ParseShorthand(%q) Segments: want %v, got %v", c.input, c.wantSegs, got.Segments)
}
if got.Ref != c.wantRef {
t.Errorf("ParseShorthand(%q) Ref: want %q, got %q", c.input, c.wantRef, got.Ref)
}
}
}

func TestBoundaryCandidates(t *testing.T) {
parts := gitlabresolver.ParseShorthand("gitlab.com/owner/repo/subdir")
if parts == nil {
t.Fatal("ParseShorthand returned nil")
}
bc := gitlabresolver.NewBoundaryCandidates(parts)
var results []gitlabresolver.BoundaryCandidate
for {
c, ok := bc.Next()
if !ok {
break
}
results = append(results, c)
}
if len(results) == 0 {
t.Fatal("expected candidates, got none")
}
// First candidate: owner/repo/subdir, virtual=""
if results[0].RepoPath != "owner/repo/subdir" {
t.Errorf("first RepoPath: want owner/repo/subdir, got %s", results[0].RepoPath)
}
// Should eventually get owner/repo with virtual subdir
found := false
for _, r := range results {
if r.RepoPath == "owner/repo" && r.VirtualPath == "subdir" {
found = true
}
}
if !found {
t.Errorf("expected owner/repo + subdir candidate among %v", results)
}
}

func TestBoundaryCandidatesMinTwo(t *testing.T) {
// A spec with only 2 segments should produce exactly one candidate: owner/repo with no virtual
parts := gitlabresolver.ParseShorthand("gitlab.com/owner/repo")
if parts == nil {
t.Fatal("ParseShorthand returned nil")
}
bc := gitlabresolver.NewBoundaryCandidates(parts)
cand, ok := bc.Next()
if !ok {
t.Fatal("expected at least one candidate")
}
if cand.RepoPath != "owner/repo" {
t.Errorf("RepoPath: want owner/repo, got %s", cand.RepoPath)
}
if cand.VirtualPath != "" {
t.Errorf("VirtualPath: want empty, got %s", cand.VirtualPath)
}
_, ok = bc.Next()
if ok {
t.Error("expected no more candidates for 2-segment spec")
}
}
