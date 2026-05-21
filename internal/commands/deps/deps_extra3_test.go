package deps

import (
	"testing"
)

func TestListResult_ZeroValue_Extra3(t *testing.T) {
	var r ListResult
	if len(r.Deps) != 0 || len(r.Orphaned) != 0 {
		t.Error("zero ListResult should have empty slices")
	}
}

func TestTreeNode_ZeroValue_Extra3(t *testing.T) {
	var n TreeNode
	if n.Name != "" || len(n.Children) != 0 {
		t.Error("zero TreeNode should have empty fields")
	}
}

func TestTreeNode_AssignChildren_Extra3(t *testing.T) {
	child := TreeNode{Name: "child"}
	parent := TreeNode{Name: "parent", Children: []TreeNode{child}}
	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0].Name != "child" {
		t.Errorf("expected child.Name=child, got %q", parent.Children[0].Name)
	}
}

func TestCheckIssue_ZeroValue_Extra3(t *testing.T) {
	var issue CheckIssue
	if issue.Name != "" || issue.Problem != "" {
		t.Error("zero CheckIssue should have empty fields")
	}
}

func TestCheckResult_ZeroValue_Extra3(t *testing.T) {
	var r CheckResult
	if len(r.Issues) != 0 || r.OK {
		t.Error("zero CheckResult should have empty issues and OK=false")
	}
}

func TestOrphanResult_ZeroValue_Extra3(t *testing.T) {
	var r OrphanResult
	if len(r.Orphaned) != 0 {
		t.Error("zero OrphanResult should have empty Orphaned slice")
	}
}

func TestDepEntry_IsOrphanedField_Extra3(t *testing.T) {
	e := DepEntry{Name: "foo", IsOrphaned: true}
	if !e.IsOrphaned {
		t.Error("expected IsOrphaned=true")
	}
}

func TestDepEntry_IsInsecureField_Extra3(t *testing.T) {
	e := DepEntry{Name: "foo", IsInsecure: true}
	if !e.IsInsecure {
		t.Error("expected IsInsecure=true")
	}
}

func TestDepEntry_PrimitivesSlice_Extra3(t *testing.T) {
	e := DepEntry{
		Name:       "mypkg",
		Primitives: []string{"instructions", "contexts"},
	}
	if len(e.Primitives) != 2 {
		t.Errorf("expected 2 primitives, got %d", len(e.Primitives))
	}
}

func TestSyncResult_ZeroValue_Extra3(t *testing.T) {
	var r SyncResult
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Updated) != 0 {
		t.Error("zero SyncResult should have all-empty slices")
	}
}
