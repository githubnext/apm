package deps

import (
	"strings"
	"testing"
)

// TestParityDependencyNodeGetID mirrors DependencyNode.get_id() parity.
func TestParityDependencyNodeGetID(t *testing.T) {
	n := &DependencyNode{RepoURL: "https://github.com/owner/repo"}
	if got := n.GetID(); got != "https://github.com/owner/repo" {
		t.Errorf("GetID no-ref: got %q", got)
	}

	n2 := &DependencyNode{RepoURL: "https://github.com/owner/repo", Reference: "main"}
	if got := n2.GetID(); got != "https://github.com/owner/repo#main" {
		t.Errorf("GetID with-ref: got %q", got)
	}
}

// TestParityDependencyNodeAncestorChain mirrors DependencyNode.get_ancestor_chain().
func TestParityDependencyNodeAncestorChain(t *testing.T) {
	root := &DependencyNode{RepoURL: "root"}
	mid := &DependencyNode{RepoURL: "mid", Parent: root}
	leaf := &DependencyNode{RepoURL: "leaf", Parent: mid}

	if got := leaf.GetAncestorChain(); got != "root > mid > leaf" {
		t.Errorf("ancestor chain: got %q", got)
	}
}

// TestParityCircularRefString mirrors CircularRef.__str__().
func TestParityCircularRefString(t *testing.T) {
	cr := &CircularRef{CyclePath: []string{"a", "b", "c"}, DetectedAtDepth: 2}
	got := cr.String()
	if !strings.Contains(got, "a -> b -> c -> a") {
		t.Errorf("circular ref string: got %q", got)
	}

	// Single element -- no arrow appended
	cr2 := &CircularRef{CyclePath: []string{"a"}}
	got2 := cr2.String()
	if !strings.Contains(got2, "a") {
		t.Errorf("single circular ref: got %q", got2)
	}

	// Empty
	cr3 := &CircularRef{}
	got3 := cr3.String()
	if !strings.Contains(got3, "empty path") {
		t.Errorf("empty circular ref: got %q", got3)
	}
}

// TestParityFlatDependencyMapFirstWins verifies first-wins conflict semantics.
func TestParityFlatDependencyMapFirstWins(t *testing.T) {
	f := NewFlatDependencyMap()
	f.AddDependency("https://github.com/owner/repo", "main")
	f.AddDependency("https://github.com/owner/repo", "v1.0") // should not overwrite

	if got := f.Dependencies["https://github.com/owner/repo"]; got != "main" {
		t.Errorf("first-wins: expected main, got %q", got)
	}
	if f.TotalDependencies() != 1 {
		t.Errorf("total: expected 1, got %d", f.TotalDependencies())
	}
}

// TestParityFlatDependencyMapInstallOrder mirrors install_order list.
func TestParityFlatDependencyMapInstallOrder(t *testing.T) {
	f := NewFlatDependencyMap()
	f.AddDependency("a", "main")
	f.AddDependency("b", "main")
	f.AddDependency("c", "main")

	order := f.GetInstallationList()
	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Errorf("install order: got %v", order)
	}
}

// TestParityDependencyTreeMaxDepth mirrors DependencyTree.max_depth.
func TestParityDependencyTreeMaxDepth(t *testing.T) {
	tree := NewDependencyTree("root")
	tree.AddNode(&DependencyNode{RepoURL: "a", Depth: 1})
	tree.AddNode(&DependencyNode{RepoURL: "b", Depth: 3})
	tree.AddNode(&DependencyNode{RepoURL: "c", Depth: 2})

	if tree.MaxDepth != 3 {
		t.Errorf("max_depth: expected 3, got %d", tree.MaxDepth)
	}
}

// TestParityDependencyTreeHasDependency mirrors DependencyTree.has_dependency().
func TestParityDependencyTreeHasDependency(t *testing.T) {
	tree := NewDependencyTree("root")
	tree.AddNode(&DependencyNode{RepoURL: "https://github.com/owner/repo", Depth: 1})

	if !tree.HasDependency("https://github.com/owner/repo") {
		t.Error("should find existing dep")
	}
	if tree.HasDependency("https://github.com/other/repo") {
		t.Error("should not find missing dep")
	}
}

// TestParityDependencyGraphIsValid mirrors DependencyGraph.is_valid().
func TestParityDependencyGraphIsValid(t *testing.T) {
	g := NewDependencyGraph("root")
	if !g.IsValid() {
		t.Error("empty graph should be valid")
	}

	g.AddError("some error")
	if g.IsValid() {
		t.Error("graph with errors should be invalid")
	}

	g2 := NewDependencyGraph("root2")
	g2.AddCircularDependency(CircularRef{CyclePath: []string{"a", "b"}})
	if g2.IsValid() {
		t.Error("graph with circular dep should be invalid")
	}
}

// TestParityDependencyGraphSummaryKeys mirrors DependencyGraph.get_summary() keys.
func TestParityDependencyGraphSummaryKeys(t *testing.T) {
	g := NewDependencyGraph("my-pkg")
	g.Flattened.AddDependency("dep1", "main")
	summary := g.GetSummary()

	expected := []string{
		"root_package", "total_dependencies", "max_depth",
		"has_circular_dependencies", "circular_count",
		"has_conflicts", "conflict_count",
		"has_errors", "error_count", "is_valid",
	}
	for _, k := range expected {
		if _, ok := summary[k]; !ok {
			t.Errorf("summary missing key %q", k)
		}
	}
	if summary["root_package"] != "my-pkg" {
		t.Errorf("root_package: got %v", summary["root_package"])
	}
	if summary["total_dependencies"].(int) != 1 {
		t.Errorf("total_dependencies: got %v", summary["total_dependencies"])
	}
}

// TestParityLockedDependencyGetUniqueKey mirrors LockedDependency.get_unique_key().
func TestParityLockedDependencyGetUniqueKey(t *testing.T) {
	// Normal dep
	ld := &LockedDependency{RepoURL: "https://github.com/owner/repo", Depth: 1}
	if got := ld.GetUniqueKey(); got != "https://github.com/owner/repo" {
		t.Errorf("normal dep key: got %q", got)
	}

	// Local dep
	ld2 := &LockedDependency{RepoURL: "https://github.com/owner/repo", Source: "local", LocalPath: "./local/pkg"}
	if got := ld2.GetUniqueKey(); got != "./local/pkg" {
		t.Errorf("local dep key: got %q", got)
	}

	// Virtual dep
	ld3 := &LockedDependency{
		RepoURL:     "https://github.com/owner/mono",
		IsVirtual:   true,
		VirtualPath: "packages/sub",
	}
	if got := ld3.GetUniqueKey(); got != "https://github.com/owner/mono/packages/sub" {
		t.Errorf("virtual dep key: got %q", got)
	}
}

// TestParityLockedDependencyToMap mirrors LockedDependency.to_dict() behavior.
func TestParityLockedDependencyToMap(t *testing.T) {
	ld := &LockedDependency{
		RepoURL:        "https://github.com/owner/repo",
		ResolvedCommit: "abc1234",
		Depth:          2,
		IsDev:          true,
	}
	m := ld.ToMap()
	if m["repo_url"] != "https://github.com/owner/repo" {
		t.Errorf("repo_url: got %v", m["repo_url"])
	}
	if m["resolved_commit"] != "abc1234" {
		t.Errorf("resolved_commit: got %v", m["resolved_commit"])
	}
	if m["depth"] != 2 {
		t.Errorf("depth: got %v", m["depth"])
	}
	if m["is_dev"] != true {
		t.Errorf("is_dev: got %v", m["is_dev"])
	}
	// depth==1 should be omitted (default)
	ld2 := &LockedDependency{RepoURL: "x", Depth: 1}
	m2 := ld2.ToMap()
	if _, ok := m2["depth"]; ok {
		t.Error("depth==1 should be omitted from map")
	}
}

// TestParityLockedDependencyDeployedFilesSorted mirrors sorted deployed_files.
func TestParityLockedDependencyDeployedFilesSorted(t *testing.T) {
	ld := &LockedDependency{
		RepoURL:       "x",
		DeployedFiles: []string{"z.md", "a.md", "m.md"},
	}
	m := ld.ToMap()
	files := m["deployed_files"].([]string)
	if files[0] != "a.md" || files[1] != "m.md" || files[2] != "z.md" {
		t.Errorf("deployed_files not sorted: %v", files)
	}
}

// TestParityLockedDependencyFromMap mirrors LockedDependency.from_dict().
func TestParityLockedDependencyFromMap(t *testing.T) {
	data := map[string]interface{}{
		"repo_url":        "https://github.com/owner/repo",
		"resolved_commit": "deadbeef",
		"depth":           float64(2),
		"is_dev":          true,
		"port":            float64(7999),
	}
	ld := LockedDependencyFromMap(data)
	if ld.RepoURL != "https://github.com/owner/repo" {
		t.Errorf("repo_url: got %q", ld.RepoURL)
	}
	if ld.ResolvedCommit != "deadbeef" {
		t.Errorf("resolved_commit: got %q", ld.ResolvedCommit)
	}
	if ld.Depth != 2 {
		t.Errorf("depth: got %d", ld.Depth)
	}
	if !ld.IsDev {
		t.Error("is_dev should be true")
	}
	if ld.Port != 7999 {
		t.Errorf("port: got %d", ld.Port)
	}
}

// TestParityLockedDependencyPortValidation mirrors port range validation in from_dict().
func TestParityLockedDependencyPortValidation(t *testing.T) {
	// Invalid port (out of range) should be ignored
	data := map[string]interface{}{
		"repo_url": "x",
		"port":     float64(99999),
	}
	ld := LockedDependencyFromMap(data)
	if ld.Port != 0 {
		t.Errorf("out-of-range port should be 0, got %d", ld.Port)
	}

	// Valid port
	data2 := map[string]interface{}{
		"repo_url": "x",
		"port":     float64(443),
	}
	ld2 := LockedDependencyFromMap(data2)
	if ld2.Port != 443 {
		t.Errorf("valid port: got %d", ld2.Port)
	}
}

// TestParityInstalledPackageFields checks InstalledPackage fields exist.
func TestParityInstalledPackageFields(t *testing.T) {
	ip := InstalledPackage{
		RepoURL:        "https://github.com/owner/repo",
		Reference:      "main",
		ResolvedCommit: "abc1234",
		Depth:          1,
		IsDev:          false,
	}
	if ip.RepoURL == "" {
		t.Error("RepoURL should be set")
	}
	if ip.Depth != 1 {
		t.Errorf("Depth: got %d", ip.Depth)
	}
}
