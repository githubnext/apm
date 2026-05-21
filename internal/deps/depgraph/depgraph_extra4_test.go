package depgraph

import "testing"

func TestDependencyRef_RepoURL_Extra4(t *testing.T) {
ref := DependencyRef{RepoURL: "org/repo"}
if ref.RepoURL != "org/repo" {
t.Errorf("unexpected RepoURL: %s", ref.RepoURL)
}
}

func TestDependencyRef_Reference_Extra4(t *testing.T) {
ref := DependencyRef{Reference: "v1.2.3"}
if ref.Reference != "v1.2.3" {
t.Errorf("unexpected Reference: %s", ref.Reference)
}
}

func TestDependencyRef_DisplayName_Extra4(t *testing.T) {
ref := DependencyRef{DisplayName: "mylib"}
if ref.DisplayName != "mylib" {
t.Errorf("unexpected DisplayName: %s", ref.DisplayName)
}
}

func TestDependencyRef_UniqueKey_Extra4(t *testing.T) {
ref := DependencyRef{UniqueKey: "org/repo/sub"}
if ref.UniqueKey != "org/repo/sub" {
t.Errorf("unexpected UniqueKey: %s", ref.UniqueKey)
}
}

func TestDependencyRef_VirtualPath_Extra4(t *testing.T) {
ref := DependencyRef{VirtualPath: "packages/core"}
if ref.VirtualPath != "packages/core" {
t.Errorf("unexpected VirtualPath: %s", ref.VirtualPath)
}
}

func TestDependencyTree_MaxDepthZero_Extra4(t *testing.T) {
tree := NewDependencyTree()
if tree.MaxDepth != 0 {
t.Errorf("expected 0 max depth for empty tree, got %d", tree.MaxDepth)
}
}

func TestDependencyTree_SingleNodeDepth0_Extra4(t *testing.T) {
tree := NewDependencyTree()
node := &DependencyNode{Ref: DependencyRef{RepoURL: "org/repo", Reference: "v1.0.0"}}
tree.AddNode(node)
nodes := tree.GetNodesAtDepth(0)
if len(nodes) != 1 {
t.Fatalf("expected 1 node at depth 0, got %d", len(nodes))
}
}

func TestFlatDependencyMap_ZeroTotal_Extra4(t *testing.T) {
m := NewFlatDependencyMap()
if m.TotalDependencies() != 0 {
t.Errorf("expected 0 deps, got %d", m.TotalDependencies())
}
}

func TestFlatDependencyMap_NoConflictsEmpty_Extra4(t *testing.T) {
m := NewFlatDependencyMap()
if m.HasConflicts() {
t.Error("expected no conflicts in empty map")
}
}

func TestDependencyGraph_ValidNew_Extra4(t *testing.T) {
g := NewDependencyGraph("root")
if !g.IsValid() {
t.Error("expected new graph to be valid")
}
}

func TestDependencyGraph_NoErrorsNew_Extra4(t *testing.T) {
g := NewDependencyGraph("root")
if g.HasErrors() {
t.Error("expected new graph to have no errors")
}
}

func TestCircularRef_StringNonEmpty_Extra4(t *testing.T) {
c := CircularRef{CyclePath: []string{"a", "b", "a"}}
s := c.String()
if len(s) == 0 {
t.Error("expected non-empty string")
}
}

func TestDependencyNode_RefField_Extra4(t *testing.T) {
node := &DependencyNode{
Ref: DependencyRef{RepoURL: "org/pkg", Reference: "v2.0"},
}
if node.Ref.RepoURL != "org/pkg" {
t.Errorf("unexpected RepoURL: %s", node.Ref.RepoURL)
}
}
