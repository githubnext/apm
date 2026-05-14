package lockfile

import (
"testing"
)

const sampleYAML = `lockfile_version: "1"
generated_at: 2026-01-01T00:00:00Z
apm_version: "1.0.0"
dependencies:
  - repo_url: https://github.com/owner/repo
    resolved_commit: abc123
    depth: 1
    is_dev: false
  - repo_url: https://github.com/owner/repo2
    resolved_commit: def456
    depth: 2
    is_dev: true
mcp_servers:
  - my-server
local_deployed_files:
  - .github/copilot-instructions.md
`

func TestFromYAMLBasic(t *testing.T) {
lf, err := FromYAML(sampleYAML)
if err != nil {
t.Fatalf("FromYAML error: %v", err)
}
if lf.LockfileVersion != "1" {
t.Errorf("expected version 1, got %s", lf.LockfileVersion)
}
if lf.APMVersion != "1.0.0" {
t.Errorf("expected APMVersion 1.0.0, got %s", lf.APMVersion)
}
// 2 real deps + 1 self-entry (local_deployed_files not empty)
if !lf.HasDependency("https://github.com/owner/repo") {
t.Error("expected dep1")
}
if !lf.HasDependency("https://github.com/owner/repo2") {
t.Error("expected dep2")
}
if !lf.HasDependency(".") {
t.Error("expected self entry from local_deployed_files")
}
if len(lf.MCPServers) != 1 || lf.MCPServers[0] != "my-server" {
t.Errorf("unexpected mcp_servers: %v", lf.MCPServers)
}
}

func TestNewLockFile(t *testing.T) {
lf := NewLockFile()
if lf.LockfileVersion != "1" {
t.Errorf("expected version 1")
}
if lf.GeneratedAt == "" {
t.Error("expected non-empty generated_at")
}
}

func TestAddGetDependency(t *testing.T) {
lf := NewLockFile()
dep := &LockedDependency{
RepoURL: "https://github.com/foo/bar",
Depth:   1,
}
lf.AddDependency(dep)
got := lf.GetDependency("https://github.com/foo/bar")
if got == nil {
t.Error("expected dependency")
}
if got.RepoURL != dep.RepoURL {
t.Errorf("repo_url mismatch")
}
}

func TestGetAllDependenciesSorted(t *testing.T) {
lf := NewLockFile()
lf.AddDependency(&LockedDependency{RepoURL: "b", Depth: 2})
lf.AddDependency(&LockedDependency{RepoURL: "a", Depth: 1})
lf.AddDependency(&LockedDependency{RepoURL: "c", Depth: 1})
deps := lf.GetAllDependencies()
if deps[0].RepoURL != "a" || deps[1].RepoURL != "c" || deps[2].RepoURL != "b" {
t.Errorf("unexpected order: %v", func() []string {
var s []string
for _, d := range deps {
s = append(s, d.RepoURL)
}
return s
}())
}
}

func TestGetLockfilePath(t *testing.T) {
p := GetLockfilePath("/project")
if p == "" {
t.Error("expected non-empty path")
}
}

func TestLockedDepToDict(t *testing.T) {
dep := &LockedDependency{
RepoURL:        "https://example.com/repo",
ResolvedCommit: "abc",
Depth:          1,
IsDev:          true,
}
d := dep.ToDict()
if d["repo_url"] != "https://example.com/repo" {
t.Error("repo_url mismatch")
}
if d["is_dev"] != true {
t.Error("is_dev should be true")
}
// depth == 1 should not be emitted
if _, ok := d["depth"]; ok {
t.Error("depth=1 should be omitted")
}
}
