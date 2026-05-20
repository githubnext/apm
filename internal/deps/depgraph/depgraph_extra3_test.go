package depgraph

import (
	"strings"
	"testing"
)

func TestDependencyRef_ZeroValue_Extra3(t *testing.T) {
	var r DependencyRef
	if r.RepoURL != "" || r.Reference != "" || r.UniqueKey != "" || r.VirtualPath != "" || r.DisplayName != "" {
		t.Error("zero value DependencyRef should have empty fields")
	}
	if r.ID() != "" {
		t.Errorf("zero value ID() = %q, want empty", r.ID())
	}
}

func TestDependencyRef_ID_VirtualPath_Extra3(t *testing.T) {
	r := DependencyRef{
		UniqueKey:   "owner/repo/pkg",
		Reference:   "v1.0",
		VirtualPath: "pkg",
	}
	id := r.ID()
	if !strings.Contains(id, "v1.0") {
		t.Errorf("ID() = %q should contain reference", id)
	}
	if !strings.Contains(id, "owner/repo/pkg") {
		t.Errorf("ID() = %q should contain unique key", id)
	}
}

func TestDependencyNode_IsDev_Extra3(t *testing.T) {
	n := &DependencyNode{
		Ref:   DependencyRef{DisplayName: "devpkg", UniqueKey: "owner/devpkg"},
		Depth: 1,
		IsDev: true,
	}
	if !n.IsDev {
		t.Error("IsDev should be true")
	}
	if n.GetDisplayName() != "devpkg" {
		t.Errorf("GetDisplayName() = %q, want devpkg", n.GetDisplayName())
	}
}

func TestDependencyNode_NoParent_Extra3(t *testing.T) {
	n := &DependencyNode{
		Ref: DependencyRef{DisplayName: "root"},
	}
	chain := n.GetAncestorChain()
	if chain != "root" {
		t.Errorf("single node chain = %q, want root", chain)
	}
}

func TestDependencyNode_DeepChain_Extra3(t *testing.T) {
	a := &DependencyNode{Ref: DependencyRef{DisplayName: "a"}}
	b := &DependencyNode{Ref: DependencyRef{DisplayName: "b"}, Parent: a}
	c := &DependencyNode{Ref: DependencyRef{DisplayName: "c"}, Parent: b}
	d := &DependencyNode{Ref: DependencyRef{DisplayName: "d"}, Parent: c}
	chain := d.GetAncestorChain()
	if chain != "a > b > c > d" {
		t.Errorf("deep chain = %q, want a > b > c > d", chain)
	}
}

func TestDependencyTree_AddDuplicate_Extra3(t *testing.T) {
	tree := NewDependencyTree()
	n1 := &DependencyNode{Ref: DependencyRef{UniqueKey: "k1", DisplayName: "p1"}, Depth: 0}
	n2 := &DependencyNode{Ref: DependencyRef{UniqueKey: "k1", DisplayName: "p1v2"}, Depth: 0}
	tree.AddNode(n1)
	tree.AddNode(n2)
	got := tree.GetNode("k1")
	if got == nil {
		t.Fatal("expected node to exist")
	}
}

func TestDependencyTree_MultipleDepths_Extra3(t *testing.T) {
	tree := NewDependencyTree()
	for i := 0; i < 5; i++ {
		key := strings.Repeat("x", i+1)
		n := &DependencyNode{Ref: DependencyRef{UniqueKey: key, DisplayName: key}, Depth: i}
		tree.AddNode(n)
	}
	if tree.MaxDepth != 4 {
		t.Errorf("MaxDepth = %d, want 4", tree.MaxDepth)
	}
	nodes := tree.GetNodesAtDepth(2)
	if len(nodes) != 1 {
		t.Errorf("nodes at depth 2 = %d, want 1", len(nodes))
	}
}

func TestFlatDependencyMap_ConflictRecord_Extra3(t *testing.T) {
	m := NewFlatDependencyMap()
	r1 := DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo"}
	r2 := DependencyRef{UniqueKey: "owner/repo", RepoURL: "owner/repo"}
	m.AddDependency(r1, false)
	m.AddDependency(r2, true)
	if !m.HasConflicts() {
		t.Error("expected conflict to be recorded")
	}
}

func TestFlatDependencyMap_InstallOrder_Extra3(t *testing.T) {
	m := NewFlatDependencyMap()
	keys := []string{"c", "a", "b"}
	for _, k := range keys {
		m.AddDependency(DependencyRef{UniqueKey: k, RepoURL: k}, false)
	}
	list := m.GetInstallationList()
	if len(list) != 3 {
		t.Fatalf("list len = %d, want 3", len(list))
	}
	if list[0].UniqueKey != "c" || list[1].UniqueKey != "a" || list[2].UniqueKey != "b" {
		t.Error("install order not preserved")
	}
}

func TestDependencyGraph_Summary_Extra3(t *testing.T) {
	g := NewDependencyGraph("mypkg")
	s := g.GetSummary()
	if s["root_package"] != "mypkg" {
		t.Errorf("root_package = %v, want mypkg", s["root_package"])
	}
	if s["is_valid"] != true {
		t.Errorf("is_valid = %v, want true", s["is_valid"])
	}
	if s["total_dependencies"] != 0 {
		t.Errorf("total_dependencies = %v, want 0", s["total_dependencies"])
	}
}

func TestDependencyGraph_InvalidAfterCycle_Extra3(t *testing.T) {
	g := NewDependencyGraph("root")
	g.AddCircularDependency(CircularRef{CyclePath: []string{"a", "b", "a"}})
	if g.IsValid() {
		t.Error("graph with circular dependency should not be valid")
	}
	if !g.HasCircularDependencies() {
		t.Error("expected HasCircularDependencies to be true")
	}
}

func TestDependencyGraph_MultipleErrors_Extra3(t *testing.T) {
	g := NewDependencyGraph("root")
	g.AddError("error1")
	g.AddError("error2")
	if len(g.ResolutionErrors) != 2 {
		t.Errorf("errors count = %d, want 2", len(g.ResolutionErrors))
	}
	if !g.HasErrors() {
		t.Error("expected HasErrors true")
	}
	s := g.GetSummary()
	if s["error_count"] != 2 {
		t.Errorf("error_count = %v, want 2", s["error_count"])
	}
}

func TestCircularRef_LongCycle_Extra3(t *testing.T) {
	cr := CircularRef{CyclePath: []string{"a", "b", "c", "d", "a"}, DetectedAtDepth: 4}
	s := cr.String()
	if !strings.Contains(s, "a") || !strings.Contains(s, "d") {
		t.Errorf("CircularRef.String() = %q missing expected nodes", s)
	}
	if cr.DetectedAtDepth != 4 {
		t.Errorf("DetectedAtDepth = %d, want 4", cr.DetectedAtDepth)
	}
}

func TestConflictInfo_MultipleLoserRefs_Extra3(t *testing.T) {
	ci := ConflictInfo{
		RepoURL: "owner/repo",
		Winner:  DependencyRef{UniqueKey: "owner/repo@v1"},
		Conflicts: []DependencyRef{
			{UniqueKey: "owner/repo@v2"},
			{UniqueKey: "owner/repo@v3"},
		},
		Reason: "first declared wins",
	}
	s := ci.String()
	if !strings.Contains(s, "owner/repo@v1") {
		t.Errorf("conflict string missing winner: %q", s)
	}
	if !strings.Contains(s, "owner/repo@v2") {
		t.Errorf("conflict string missing loser v2: %q", s)
	}
}
