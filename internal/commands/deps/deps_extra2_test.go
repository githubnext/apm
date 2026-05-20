package deps

import (
	"strings"
	"testing"
)

func TestDepEntry_ZeroValue(t *testing.T) {
	var de DepEntry
	if de.Name != "" {
		t.Errorf("expected empty Name, got %q", de.Name)
	}
}

func TestDepEntry_SourceLabel(t *testing.T) {
	cases := []struct {
		dm   map[string]any
		want string
	}{
		{map[string]any{"source": "github"}, "github"},
		{map[string]any{}, ""},
		{nil, ""},
	}
	for _, c := range cases {
		got := sourceLabel(c.dm)
		if c.want != "" && !strings.Contains(got, c.want) {
			t.Errorf("sourceLabel(%v) = %q, want to contain %q", c.dm, got, c.want)
		}
	}
}

func TestSanitizeMermaid_EmptyString(t *testing.T) {
	got := sanitizeMermaid("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestSanitizeMermaid_NoSpecial(t *testing.T) {
	got := sanitizeMermaid("hello")
	if got != "hello" {
		t.Errorf("expected hello, got %q", got)
	}
}

func TestSanitizeMermaid_Slash(t *testing.T) {
	got := sanitizeMermaid("owner/repo")
	if strings.Contains(got, "/") {
		t.Errorf("expected / to be replaced, got %q", got)
	}
}

func TestListOptions_ZeroValue(t *testing.T) {
	var lo ListOptions
	_ = lo
}

func TestTreeNode_Children(t *testing.T) {
	child := TreeNode{Name: "child"}
	parent := TreeNode{
		Name:     "parent",
		Children: []TreeNode{child},
	}
	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0].Name != "child" {
		t.Errorf("unexpected child name: %q", parent.Children[0].Name)
	}
}

func TestTreeNode_ZeroValue(t *testing.T) {
	var tn TreeNode
	if tn.Name != "" || len(tn.Children) != 0 {
		t.Error("expected zero value")
	}
}

func TestTreeOptions_ZeroValue(t *testing.T) {
	var to TreeOptions
	_ = to
}

func TestSyncOptions_ZeroValue(t *testing.T) {
	var so SyncOptions
	_ = so
}

func TestOrphanOptions_ZeroValue(t *testing.T) {
	var oo OrphanOptions
	_ = oo
}

func TestCheckIssue_ZeroValue(t *testing.T) {
	var ci CheckIssue
	if ci.Name != "" || ci.Problem != "" {
		t.Error("expected zero value")
	}
}

func TestGraphOptions_ZeroValue(t *testing.T) {
	var go_ GraphOptions
	_ = go_
}

func TestListResult_NonNilOrphaned(t *testing.T) {
	lr := &ListResult{Orphaned: []string{"x"}}
	if len(lr.Orphaned) != 1 {
		t.Errorf("expected 1 orphan")
	}
}
