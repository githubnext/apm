package depgraph_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/deps/depgraph"
)

func TestCircularRef_String_Empty(t *testing.T) {
	cr := depgraph.CircularRef{CyclePath: []string{}, DetectedAtDepth: 0}
	s := cr.String()
	if !strings.Contains(s, "Circular dependency") {
		t.Errorf("expected circular dependency message, got %q", s)
	}
	if !strings.Contains(s, "empty path") {
		t.Errorf("expected empty path in message, got %q", s)
	}
}

func TestCircularRef_String_SingleNode(t *testing.T) {
	cr := depgraph.CircularRef{CyclePath: []string{"github.com/a/b"}, DetectedAtDepth: 1}
	s := cr.String()
	if !strings.Contains(s, "github.com/a/b") {
		t.Errorf("expected repo URL in message, got %q", s)
	}
}

func TestCircularRef_String_TwoNodes(t *testing.T) {
	cr := depgraph.CircularRef{CyclePath: []string{"github.com/a/b", "github.com/c/d"}, DetectedAtDepth: 2}
	s := cr.String()
	if !strings.Contains(s, "->") {
		t.Errorf("expected arrow separator in cycle string, got %q", s)
	}
	if !strings.Contains(s, "github.com/a/b") || !strings.Contains(s, "github.com/c/d") {
		t.Errorf("expected both nodes in cycle string, got %q", s)
	}
}

func TestCircularRef_String_ReturnsToStart(t *testing.T) {
	cr := depgraph.CircularRef{CyclePath: []string{"a", "b", "c"}, DetectedAtDepth: 3}
	s := cr.String()
	// should end with "-> a" to show the return to start
	if !strings.HasSuffix(s, "-> a") {
		t.Errorf("expected cycle to return to start, got %q", s)
	}
}

func TestCircularRef_DetectedAtDepth(t *testing.T) {
	cr := depgraph.CircularRef{CyclePath: []string{"x", "y"}, DetectedAtDepth: 5}
	if cr.DetectedAtDepth != 5 {
		t.Errorf("expected DetectedAtDepth=5, got %d", cr.DetectedAtDepth)
	}
}

func TestConflictInfo_String_Basic(t *testing.T) {
	ci := depgraph.ConflictInfo{
		RepoURL: "github.com/owner/repo",
		Winner:  depgraph.DependencyRef{UniqueKey: "owner/repo@v1.0.0"},
		Conflicts: []depgraph.DependencyRef{
			{UniqueKey: "owner/repo@v2.0.0"},
		},
		Reason: "first declared dependency wins",
	}
	s := ci.String()
	if !strings.Contains(s, "github.com/owner/repo") {
		t.Errorf("expected repo URL in string, got %q", s)
	}
	if !strings.Contains(s, "wins") {
		t.Errorf("expected 'wins' in string, got %q", s)
	}
}

func TestConflictInfo_String_MultipleConflicts(t *testing.T) {
	ci := depgraph.ConflictInfo{
		RepoURL: "github.com/a/b",
		Winner:  depgraph.DependencyRef{UniqueKey: "a/b@v1"},
		Conflicts: []depgraph.DependencyRef{
			{UniqueKey: "a/b@v2"},
			{UniqueKey: "a/b@v3"},
		},
		Reason: "first declared",
	}
	s := ci.String()
	if !strings.Contains(s, "a/b@v2") || !strings.Contains(s, "a/b@v3") {
		t.Errorf("expected both conflict keys, got %q", s)
	}
}

func TestConflictInfo_EmptyConflicts(t *testing.T) {
	ci := depgraph.ConflictInfo{
		RepoURL:   "github.com/x/y",
		Winner:    depgraph.DependencyRef{UniqueKey: "x/y@v1"},
		Conflicts: nil,
		Reason:    "only one",
	}
	s := ci.String()
	if s == "" {
		t.Error("expected non-empty string for ConflictInfo")
	}
}

func TestFlatDependencyMap_GetInstallationList_Empty(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	list := m.GetInstallationList()
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}

func TestFlatDependencyMap_GetInstallationList_SingleItem(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	ref := depgraph.DependencyRef{
		UniqueKey: "owner/repo@v1.0.0",
		RepoURL:   "github.com/owner/repo",
		Reference: "v1.0.0",
	}
	m.AddDependency(ref, false)
	list := m.GetInstallationList()
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].UniqueKey != ref.UniqueKey {
		t.Errorf("expected %q, got %q", ref.UniqueKey, list[0].UniqueKey)
	}
}

func TestFlatDependencyMap_GetInstallationList_ConflictNotAdded(t *testing.T) {
	m := depgraph.NewFlatDependencyMap()
	ref1 := depgraph.DependencyRef{UniqueKey: "a/b@v1", RepoURL: "github.com/a/b", Reference: "v1"}
	ref2 := depgraph.DependencyRef{UniqueKey: "a/b@v1", RepoURL: "github.com/a/b", Reference: "v2"}
	m.AddDependency(ref1, false)
	m.AddDependency(ref2, true) // same UniqueKey = conflict path
	list := m.GetInstallationList()
	if len(list) != 1 {
		t.Errorf("expected 1 item after conflict, got %d", len(list))
	}
}

func TestDependencyGraph_AddCircularDependency(t *testing.T) {
	g := depgraph.NewDependencyGraph("root")
	cr := depgraph.CircularRef{CyclePath: []string{"a", "b", "a"}, DetectedAtDepth: 2}
	g.AddCircularDependency(cr)
	if !g.HasCircularDependencies() {
		t.Error("expected HasCircularDependencies to be true")
	}
	if g.IsValid() {
		t.Error("expected IsValid to be false when circular dep present")
	}
}

func TestDependencyGraph_GetSummary_Keys(t *testing.T) {
	g := depgraph.NewDependencyGraph("myroot")
	summary := g.GetSummary()
	if _, ok := summary["root_package"]; !ok {
		t.Error("expected 'root_package' key in summary")
	}
	if _, ok := summary["is_valid"]; !ok {
		t.Error("expected 'is_valid' key in summary")
	}
}

func TestDependencyGraph_HasErrors_AfterAdd(t *testing.T) {
	g := depgraph.NewDependencyGraph("root")
	g.AddError("something went wrong")
	if !g.HasErrors() {
		t.Error("expected HasErrors to be true after adding error")
	}
	if g.IsValid() {
		t.Error("expected IsValid to be false after errors")
	}
}

func TestDependencyTree_MaxDepth_AfterMultipleAdds(t *testing.T) {
	tree := depgraph.NewDependencyTree()
	for _, depth := range []int{0, 1, 3, 2} {
		node := &depgraph.DependencyNode{
			Ref:   depgraph.DependencyRef{UniqueKey: strings.Repeat("x", depth+1), RepoURL: "url"},
			Depth: depth,
		}
		tree.AddNode(node)
	}
	if tree.MaxDepth != 3 {
		t.Errorf("expected MaxDepth=3, got %d", tree.MaxDepth)
	}
}
