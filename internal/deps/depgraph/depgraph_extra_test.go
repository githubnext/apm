package depgraph_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/depgraph"
)

func TestDependencyRef_ID_NoReference(t *testing.T) {
	ref := depgraph.DependencyRef{UniqueKey: "owner/repo"}
	if ref.ID() != "owner/repo" {
		t.Errorf("expected owner/repo, got %q", ref.ID())
	}
}

func TestDependencyRef_ID_WithReference(t *testing.T) {
	ref := depgraph.DependencyRef{UniqueKey: "owner/repo", Reference: "v1.0.0"}
	if ref.ID() != "owner/repo#v1.0.0" {
		t.Errorf("expected owner/repo#v1.0.0, got %q", ref.ID())
	}
}

func TestDependencyRef_AllFields(t *testing.T) {
	ref := depgraph.DependencyRef{
		RepoURL:     "https://github.com/owner/repo",
		Reference:   "main",
		UniqueKey:   "owner/repo",
		VirtualPath: "subdir",
		DisplayName: "owner/repo@main",
	}
	if ref.RepoURL == "" || ref.Reference == "" || ref.UniqueKey == "" {
		t.Error("fields not set correctly")
	}
}

func TestDependencyNode_GetID(t *testing.T) {
	node := &depgraph.DependencyNode{
		Ref: depgraph.DependencyRef{UniqueKey: "a/b", Reference: "main"},
	}
	if node.GetID() != "a/b#main" {
		t.Errorf("unexpected ID: %q", node.GetID())
	}
}

func TestDependencyNode_GetDisplayName(t *testing.T) {
	node := &depgraph.DependencyNode{
		Ref: depgraph.DependencyRef{DisplayName: "my-package"},
	}
	if node.GetDisplayName() != "my-package" {
		t.Errorf("unexpected display name: %q", node.GetDisplayName())
	}
}

func TestDependencyNode_GetAncestorChain_Single(t *testing.T) {
	node := &depgraph.DependencyNode{
		Ref: depgraph.DependencyRef{DisplayName: "root"},
	}
	chain := node.GetAncestorChain()
	if chain != "root" {
		t.Errorf("single node chain: expected root, got %q", chain)
	}
}

func TestDependencyNode_GetAncestorChain_ThreeLevels(t *testing.T) {
	root := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{DisplayName: "root"}}
	mid := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{DisplayName: "mid"}, Parent: root}
	leaf := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{DisplayName: "leaf"}, Parent: mid}

	chain := leaf.GetAncestorChain()
	if !strings.Contains(chain, "root") || !strings.Contains(chain, "mid") || !strings.Contains(chain, "leaf") {
		t.Errorf("expected root>mid>leaf chain, got %q", chain)
	}
}

func TestNewDependencyTree(t *testing.T) {
	tree := depgraph.NewDependencyTree()
	if tree == nil {
		t.Fatal("expected non-nil tree")
	}
}

func TestDependencyTree_AddAndGetNode(t *testing.T) {
	tree := depgraph.NewDependencyTree()
	node := &depgraph.DependencyNode{
		Ref: depgraph.DependencyRef{UniqueKey: "a/b"},
	}
	tree.AddNode(node)
	got := tree.GetNode("a/b")
	if got == nil {
		t.Fatal("expected to find node a/b")
	}
	if got.Ref.UniqueKey != "a/b" {
		t.Errorf("unexpected node: %v", got)
	}
}

func TestDependencyTree_GetNode_Missing(t *testing.T) {
	tree := depgraph.NewDependencyTree()
	if tree.GetNode("nonexistent") != nil {
		t.Error("expected nil for missing node")
	}
}

func TestDependencyTree_HasDependency(t *testing.T) {
	tree := depgraph.NewDependencyTree()
	node := &depgraph.DependencyNode{
		Ref: depgraph.DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo"},
	}
	tree.AddNode(node)
	if !tree.HasDependency("owner/repo") {
		t.Error("expected HasDependency=true")
	}
	if tree.HasDependency("other/repo") {
		t.Error("expected HasDependency=false for absent repo")
	}
}

func TestDependencyTree_GetNodesAtDepth(t *testing.T) {
	tree := depgraph.NewDependencyTree()
	n0 := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{UniqueKey: "root"}, Depth: 0}
	n1a := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{UniqueKey: "a"}, Depth: 1}
	n1b := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{UniqueKey: "b"}, Depth: 1}
	n2 := &depgraph.DependencyNode{Ref: depgraph.DependencyRef{UniqueKey: "c"}, Depth: 2}
	for _, n := range []*depgraph.DependencyNode{n0, n1a, n1b, n2} {
		tree.AddNode(n)
	}
	depth1 := tree.GetNodesAtDepth(1)
	if len(depth1) != 2 {
		t.Errorf("expected 2 nodes at depth 1, got %d", len(depth1))
	}
}

func TestNewFlatDependencyMap(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	if m == nil {
		t.Fatal("expected non-nil map")
	}
	if m.TotalDependencies() != 0 {
		t.Error("new map should have 0 dependencies")
	}
}

func TestFlatDependencyMap_AddAndGet(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	ref := depgraph.DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo"}
	m.AddDependency(ref, false)
	got, ok := m.GetDependency("owner/repo")
	if !ok {
		t.Fatal("expected to find dependency")
	}
	if got.UniqueKey != "owner/repo" {
		t.Errorf("unexpected dependency: %v", got)
	}
}

func TestFlatDependencyMap_HasConflicts(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	if m.HasConflicts() {
		t.Error("new map should have no conflicts")
	}
	// Add the same key twice with isConflict=true on the second add to trigger conflict recording
	ref1 := depgraph.DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo"}
	ref2 := depgraph.DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo", Reference: "v2.0.0"}
	m.AddDependency(ref1, false) // first add: stored normally
	m.AddDependency(ref2, true)  // second add with same key and isConflict=true
	if !m.HasConflicts() {
		t.Error("map should report HasConflicts=true after conflicting add")
	}
}

func TestFlatDependencyMap_TotalDependencies(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	for i := 0; i < 5; i++ {
		key := "owner/repo" + string(rune('a'+i))
		m.AddDependency(depgraph.DependencyRef{UniqueKey: key}, false)
	}
	if m.TotalDependencies() != 5 {
		t.Errorf("expected 5, got %d", m.TotalDependencies())
	}
}

func TestNewDependencyGraph(t *testing.T) {
	g := depgraph.NewDependencyGraph("my-root")
	if g == nil {
		t.Fatal("expected non-nil graph")
	}
	if g.IsValid() == false {
		// Fresh graph with no errors should be valid
	}
}

func TestDependencyGraph_AddError(t *testing.T) {
	g := depgraph.NewDependencyGraph("root")
	if g.HasErrors() {
		t.Error("new graph should have no errors")
	}
	g.AddError("something went wrong")
	if !g.HasErrors() {
		t.Error("graph should have errors after AddError")
	}
}

func TestDependencyGraph_IsValid(t *testing.T) {
	g := depgraph.NewDependencyGraph("root")
	if !g.IsValid() {
		t.Error("empty graph should be valid")
	}
	g.AddError("error")
	if g.IsValid() {
		t.Error("graph with errors should not be valid")
	}
}

func TestDependencyGraph_GetSummary(t *testing.T) {
	g := depgraph.NewDependencyGraph("root")
	s := g.GetSummary()
	if s == nil {
		t.Fatal("expected non-nil summary")
	}
}
