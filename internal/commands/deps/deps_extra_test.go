package deps

import (
	"testing"
)

func TestListOptions_Fields(t *testing.T) {
	opts := ListOptions{
		ProjectRoot:  "/my/project",
		Scope:        "user",
		JSON:         true,
		InsecureOnly: false,
		NoColor:      true,
	}
	if opts.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot mismatch: %q", opts.ProjectRoot)
	}
	if opts.Scope != "user" {
		t.Errorf("Scope mismatch: %q", opts.Scope)
	}
	if !opts.JSON {
		t.Error("JSON should be true")
	}
	if opts.InsecureOnly {
		t.Error("InsecureOnly should be false")
	}
	if !opts.NoColor {
		t.Error("NoColor should be true")
	}
}

func TestCheckIssue_Fields(t *testing.T) {
	ci := CheckIssue{
		Name:    "owner/pkg",
		Problem: "outdated version",
	}
	if ci.Name != "owner/pkg" {
		t.Errorf("Name mismatch: %q", ci.Name)
	}
	if ci.Problem != "outdated version" {
		t.Errorf("Problem mismatch: %q", ci.Problem)
	}
}

func TestSyncResult_Fields(t *testing.T) {
	sr := SyncResult{
		Removed: []string{"old-pkg", "stale-pkg"},
		Added:   []string{},
	}
	if len(sr.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(sr.Removed))
	}
	if sr.Removed[0] != "old-pkg" {
		t.Errorf("Removed[0] = %q", sr.Removed[0])
	}
	if len(sr.Added) != 0 {
		t.Error("expected no added")
	}
}

func TestOrphanResult_Fields(t *testing.T) {
	or_ := OrphanResult{
		Orphaned: []string{"orphan-a", "orphan-b", "orphan-c"},
	}
	if len(or_.Orphaned) != 3 {
		t.Errorf("expected 3 orphaned, got %d", len(or_.Orphaned))
	}
	if or_.Orphaned[2] != "orphan-c" {
		t.Errorf("Orphaned[2] = %q", or_.Orphaned[2])
	}
}

func TestCheckResult_Fields(t *testing.T) {
	cr := CheckResult{
		Issues: []CheckIssue{
			{Name: "pkg-a", Problem: "missing"},
			{Name: "pkg-b", Problem: "version mismatch"},
		},
	}
	if len(cr.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(cr.Issues))
	}
	if cr.Issues[1].Name != "pkg-b" {
		t.Errorf("Issues[1].Name = %q", cr.Issues[1].Name)
	}
}

func TestListResult_EmptyOrphans(t *testing.T) {
	r := ListResult{
		Deps:     []DepEntry{{Name: "pkg-a", Source: "github"}},
		Orphaned: nil,
	}
	if len(r.Deps) != 1 {
		t.Errorf("expected 1 dep, got %d", len(r.Deps))
	}
	if len(r.Orphaned) != 0 {
		t.Errorf("expected no orphaned, got %d", len(r.Orphaned))
	}
}

func TestGraphOptions_Fields(t *testing.T) {
	opts := GraphOptions{
		ProjectRoot: "/root",
		Format:      "mermaid",
	}
	if opts.ProjectRoot != "/root" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Format != "mermaid" {
		t.Errorf("Format = %q", opts.Format)
	}
}

func TestDepEntry_CommitAndRef(t *testing.T) {
	e := DepEntry{
		Name:    "my-dep",
		Commit:  "abc123def",
		Ref:     "v2.0.1",
		RepoURL: "https://github.com/owner/my-dep",
	}
	if e.Commit != "abc123def" {
		t.Errorf("Commit = %q", e.Commit)
	}
	if e.Ref != "v2.0.1" {
		t.Errorf("Ref = %q", e.Ref)
	}
	if e.RepoURL != "https://github.com/owner/my-dep" {
		t.Errorf("RepoURL = %q", e.RepoURL)
	}
}

func TestSanitizeMermaid_AllSpecialChars(t *testing.T) {
	// Verify all non-alphanum chars become underscores
	cases := []struct{ in, want string }{
		{"a/b/c", "a_b_c"},
		{"@org/pkg@1.2.3", "_org_pkg_1_2_3"},
		{"no-change", "no_change"},
		{"abc", "abc"},
	}
	for _, tc := range cases {
		got := sanitizeMermaid(tc.in)
		if got != tc.want {
			t.Errorf("sanitizeMermaid(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestTreeNode_DeepNesting(t *testing.T) {
	root := TreeNode{
		Name:    "root",
		Version: "v1.0.0",
		Children: []TreeNode{
			{
				Name:    "child",
				Version: "v2.0.0",
				Children: []TreeNode{
					{Name: "grandchild", Version: "v3.0.0"},
				},
			},
		},
	}
	if len(root.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(root.Children))
	}
	if len(root.Children[0].Children) != 1 {
		t.Errorf("expected 1 grandchild, got %d", len(root.Children[0].Children))
	}
	if root.Children[0].Children[0].Name != "grandchild" {
		t.Errorf("grandchild name = %q", root.Children[0].Children[0].Name)
	}
}
