package depgraph

import (
"strings"
"testing"
)

func TestDependencyRefID(t *testing.T) {
r := DependencyRef{UniqueKey: "owner/repo", Reference: "main"}
if r.ID() != "owner/repo#main" {
t.Errorf("unexpected ID: %s", r.ID())
}
r2 := DependencyRef{UniqueKey: "owner/repo"}
if r2.ID() != "owner/repo" {
t.Errorf("unexpected ID without ref: %s", r2.ID())
}
}

func TestDependencyNodeGetID(t *testing.T) {
n := &DependencyNode{Ref: DependencyRef{UniqueKey: "a/b", Reference: "v1"}}
if n.GetID() != "a/b#v1" {
t.Errorf("unexpected node ID: %s", n.GetID())
}
}

func TestDependencyNodeGetDisplayName(t *testing.T) {
n := &DependencyNode{Ref: DependencyRef{DisplayName: "my-pkg", UniqueKey: "a/b"}}
if n.GetDisplayName() != "my-pkg" {
t.Errorf("unexpected display name: %s", n.GetDisplayName())
}
}

func TestDependencyNodeGetAncestorChain(t *testing.T) {
root := &DependencyNode{Ref: DependencyRef{DisplayName: "root", UniqueKey: "root"}}
child := &DependencyNode{Ref: DependencyRef{DisplayName: "child", UniqueKey: "child"}, Parent: root}
chain := child.GetAncestorChain()
if !strings.Contains(chain, "root") || !strings.Contains(chain, "child") {
t.Errorf("ancestor chain missing expected names: %s", chain)
}
}

func TestDependencyTreeAddGetNode(t *testing.T) {
tree := NewDependencyTree()
ref := DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo", DisplayName: "repo"}
node := &DependencyNode{Ref: ref, Depth: 0}
tree.AddNode(node)

got := tree.GetNode("owner/repo")
if got == nil {
t.Fatal("expected node, got nil")
}
if got.GetDisplayName() != "repo" {
t.Errorf("unexpected display name: %s", got.GetDisplayName())
}
}

func TestDependencyTreeGetNodesAtDepth(t *testing.T) {
tree := NewDependencyTree()
n0 := &DependencyNode{Ref: DependencyRef{UniqueKey: "a"}, Depth: 0}
n1 := &DependencyNode{Ref: DependencyRef{UniqueKey: "b"}, Depth: 1}
n1b := &DependencyNode{Ref: DependencyRef{UniqueKey: "c"}, Depth: 1}
tree.AddNode(n0)
tree.AddNode(n1)
tree.AddNode(n1b)

depth1 := tree.GetNodesAtDepth(1)
if len(depth1) != 2 {
t.Errorf("expected 2 nodes at depth 1, got %d", len(depth1))
}
}

func TestDependencyTreeHasDependency(t *testing.T) {
tree := NewDependencyTree()
node := &DependencyNode{Ref: DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo"}}
tree.AddNode(node)
if !tree.HasDependency("owner/repo") {
t.Error("expected HasDependency to return true")
}
if tree.HasDependency("other/repo") {
t.Error("expected HasDependency to return false for missing repo")
}
}

func TestFlatDependencyMap(t *testing.T) {
m := NewFlatDependencyMap()
ref := DependencyRef{UniqueKey: "owner/pkg", RepoURL: "owner/pkg"}
m.AddDependency(ref, false)

got, ok := m.GetDependency("owner/pkg")
if !ok {
t.Fatal("expected to find dependency")
}
if got.RepoURL != "owner/pkg" {
t.Errorf("unexpected RepoURL: %s", got.RepoURL)
}
if m.HasConflicts() {
t.Error("expected no conflicts")
}
if m.TotalDependencies() != 1 {
t.Errorf("expected 1 dependency, got %d", m.TotalDependencies())
}
}

func TestFlatDependencyMapConflict(t *testing.T) {
m := NewFlatDependencyMap()
ref1 := DependencyRef{UniqueKey: "owner/pkg", RepoURL: "owner/pkg", Reference: "v1"}
ref2 := DependencyRef{UniqueKey: "owner/pkg", RepoURL: "owner/pkg", Reference: "v2"}
m.AddDependency(ref1, false)
m.AddDependency(ref2, true)
if !m.HasConflicts() {
t.Error("expected conflict")
}
}

func TestFlatDependencyMapGetInstallationList(t *testing.T) {
m := NewFlatDependencyMap()
m.AddDependency(DependencyRef{UniqueKey: "a/b"}, false)
m.AddDependency(DependencyRef{UniqueKey: "c/d"}, false)
list := m.GetInstallationList()
if len(list) != 2 {
t.Errorf("expected 2 items, got %d", len(list))
}
}

func TestDependencyGraph(t *testing.T) {
g := NewDependencyGraph("root")
if g.HasErrors() {
t.Error("new graph should have no errors")
}
if !g.IsValid() {
t.Error("new graph should be valid")
}
g.AddError("test error")
if !g.HasErrors() {
t.Error("graph should have error after AddError")
}
if g.IsValid() {
t.Error("graph with errors should not be valid")
}
}
