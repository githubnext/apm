// Package depgraph provides data structures for dependency graph representation
// and resolution.
//
// Mirrors src/apm_cli/deps/dependency_graph.py.
package depgraph

import "fmt"

// DependencyRef captures the key information from a resolved dependency
// reference needed for graph operations.
type DependencyRef struct {
	// RepoURL is the canonical repository URL (e.g. "owner/repo").
	RepoURL string
	// Reference is the git reference (branch/tag/commit), may be empty.
	Reference string
	// UniqueKey is the deduplication key (repo_url or repo_url/virtual_path).
	UniqueKey string
	// VirtualPath is the optional virtual package path suffix.
	VirtualPath string
	// DisplayName is a human-readable short name for diagnostics.
	DisplayName string
}

// ID returns a unique identifier that includes the reference when set.
func (r *DependencyRef) ID() string {
	if r.Reference != "" {
		return r.UniqueKey + "#" + r.Reference
	}
	return r.UniqueKey
}

// DependencyNode represents a single node in the dependency graph.
type DependencyNode struct {
	Ref      DependencyRef
	Depth    int
	Children []*DependencyNode
	Parent   *DependencyNode
	IsDev    bool // reached exclusively via devDependencies
}

// GetID returns the unique identifier for this node.
func (n *DependencyNode) GetID() string {
	return n.Ref.ID()
}

// GetDisplayName returns the display name for this dependency.
func (n *DependencyNode) GetDisplayName() string {
	return n.Ref.DisplayName
}

// GetAncestorChain builds a human-readable breadcrumb from this node's ancestry.
// Example: "root-pkg > mid-pkg > this-pkg"
func (n *DependencyNode) GetAncestorChain() string {
	var parts []string
	cur := n
	for cur != nil {
		parts = append([]string{cur.GetDisplayName()}, parts...)
		cur = cur.Parent
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " > "
		}
		result += p
	}
	return result
}

// CircularRef describes a detected circular dependency.
type CircularRef struct {
	// CyclePath is the ordered list of repo URLs forming the cycle.
	CyclePath []string
	// DetectedAtDepth is the depth at which the cycle was detected.
	DetectedAtDepth int
}

// formatCompleteCycle returns a string showing the full cycle visually.
func (cr *CircularRef) formatCompleteCycle() string {
	if len(cr.CyclePath) == 0 {
		return "(empty path)"
	}
	result := ""
	for i, p := range cr.CyclePath {
		if i > 0 {
			result += " -> "
		}
		result += p
	}
	// Ensure visual return to start.
	if len(cr.CyclePath) > 1 && cr.CyclePath[0] != cr.CyclePath[len(cr.CyclePath)-1] {
		result += " -> " + cr.CyclePath[0]
	}
	return result
}

func (cr *CircularRef) String() string {
	return "Circular dependency detected: " + cr.formatCompleteCycle()
}

// DependencyTree is the hierarchical representation of dependencies before
// flattening.
type DependencyTree struct {
	nodes       map[string]*DependencyNode
	nodesByDepth map[int][]*DependencyNode
	MaxDepth    int
}

// NewDependencyTree creates an empty DependencyTree.
func NewDependencyTree() *DependencyTree {
	return &DependencyTree{
		nodes:        make(map[string]*DependencyNode),
		nodesByDepth: make(map[int][]*DependencyNode),
	}
}

// AddNode inserts a node into the tree.
func (t *DependencyTree) AddNode(node *DependencyNode) {
	id := node.GetID()
	if _, exists := t.nodes[id]; !exists {
		t.nodesByDepth[node.Depth] = append(t.nodesByDepth[node.Depth], node)
	}
	t.nodes[id] = node
	if node.Depth > t.MaxDepth {
		t.MaxDepth = node.Depth
	}
}

// GetNode returns the node for the given unique key, or nil.
func (t *DependencyTree) GetNode(uniqueKey string) *DependencyNode {
	return t.nodes[uniqueKey]
}

// GetNodesAtDepth returns all nodes at a given depth.
func (t *DependencyTree) GetNodesAtDepth(depth int) []*DependencyNode {
	nodes := t.nodesByDepth[depth]
	out := make([]*DependencyNode, len(nodes))
	copy(out, nodes)
	return out
}

// HasDependency reports whether any node has the given repo URL.
func (t *DependencyTree) HasDependency(repoURL string) bool {
	for _, node := range t.nodes {
		if node.Ref.RepoURL == repoURL {
			return true
		}
	}
	return false
}

// ConflictInfo describes a dependency conflict.
type ConflictInfo struct {
	RepoURL   string
	Winner    DependencyRef
	Conflicts []DependencyRef
	Reason    string
}

func (ci *ConflictInfo) String() string {
	var conflictStrs []string
	for _, c := range ci.Conflicts {
		conflictStrs = append(conflictStrs, c.UniqueKey)
	}
	result := fmt.Sprintf("Conflict for %s: %s wins", ci.RepoURL, ci.Winner.UniqueKey)
	if len(conflictStrs) > 0 {
		result += " over "
		for i, s := range conflictStrs {
			if i > 0 {
				result += ", "
			}
			result += s
		}
	}
	result += " (" + ci.Reason + ")"
	return result
}

// FlatDependencyMap is the final flattened dependency mapping ready for
// installation.
type FlatDependencyMap struct {
	Dependencies map[string]DependencyRef
	Conflicts    []ConflictInfo
	InstallOrder []string
}

// NewFlatDependencyMap creates an empty FlatDependencyMap.
func NewFlatDependencyMap() *FlatDependencyMap {
	return &FlatDependencyMap{
		Dependencies: make(map[string]DependencyRef),
	}
}

// AddDependency adds a dependency to the flat map, recording conflicts when
// isConflict is true.
func (m *FlatDependencyMap) AddDependency(ref DependencyRef, isConflict bool) {
	key := ref.UniqueKey
	if _, exists := m.Dependencies[key]; !exists {
		m.Dependencies[key] = ref
		m.InstallOrder = append(m.InstallOrder, key)
		return
	}
	if !isConflict {
		return
	}
	// Record conflict; first-declared wins.
	existing := m.Dependencies[key]
	for i := range m.Conflicts {
		if m.Conflicts[i].RepoURL == ref.RepoURL {
			m.Conflicts[i].Conflicts = append(m.Conflicts[i].Conflicts, ref)
			return
		}
	}
	m.Conflicts = append(m.Conflicts, ConflictInfo{
		RepoURL:   ref.RepoURL,
		Winner:    existing,
		Conflicts: []DependencyRef{ref},
		Reason:    "first declared dependency wins",
	})
}

// GetDependency returns the dependency for the given unique key or the zero
// value with ok == false.
func (m *FlatDependencyMap) GetDependency(uniqueKey string) (DependencyRef, bool) {
	ref, ok := m.Dependencies[uniqueKey]
	return ref, ok
}

// HasConflicts reports whether any conflicts were recorded.
func (m *FlatDependencyMap) HasConflicts() bool {
	return len(m.Conflicts) > 0
}

// TotalDependencies returns the count of unique dependencies.
func (m *FlatDependencyMap) TotalDependencies() int {
	return len(m.Dependencies)
}

// GetInstallationList returns dependencies in install order.
func (m *FlatDependencyMap) GetInstallationList() []DependencyRef {
	out := make([]DependencyRef, 0, len(m.InstallOrder))
	for _, key := range m.InstallOrder {
		if ref, ok := m.Dependencies[key]; ok {
			out = append(out, ref)
		}
	}
	return out
}

// DependencyGraph is the complete resolved dependency information.
type DependencyGraph struct {
	RootPackageName     string
	Tree                *DependencyTree
	Flattened           *FlatDependencyMap
	CircularDependencies []CircularRef
	ResolutionErrors    []string
}

// NewDependencyGraph creates an empty DependencyGraph.
func NewDependencyGraph(rootPackageName string) *DependencyGraph {
	return &DependencyGraph{
		RootPackageName: rootPackageName,
		Tree:            NewDependencyTree(),
		Flattened:       NewFlatDependencyMap(),
	}
}

// HasCircularDependencies reports whether any cycles were detected.
func (g *DependencyGraph) HasCircularDependencies() bool {
	return len(g.CircularDependencies) > 0
}

// HasConflicts reports whether any dependency conflicts were found.
func (g *DependencyGraph) HasConflicts() bool {
	return g.Flattened.HasConflicts()
}

// HasErrors reports whether any resolution errors occurred.
func (g *DependencyGraph) HasErrors() bool {
	return len(g.ResolutionErrors) > 0
}

// IsValid reports whether the graph has no circular dependencies and no errors.
func (g *DependencyGraph) IsValid() bool {
	return !g.HasCircularDependencies() && !g.HasErrors()
}

// GetSummary returns a summary map of the dependency resolution.
func (g *DependencyGraph) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"root_package":             g.RootPackageName,
		"total_dependencies":       g.Flattened.TotalDependencies(),
		"max_depth":                g.Tree.MaxDepth,
		"has_circular_dependencies": g.HasCircularDependencies(),
		"circular_count":           len(g.CircularDependencies),
		"has_conflicts":            g.HasConflicts(),
		"conflict_count":           len(g.Flattened.Conflicts),
		"has_errors":               g.HasErrors(),
		"error_count":              len(g.ResolutionErrors),
		"is_valid":                 g.IsValid(),
	}
}

// AddError appends a resolution error.
func (g *DependencyGraph) AddError(err string) {
	g.ResolutionErrors = append(g.ResolutionErrors, err)
}

// AddCircularDependency records a circular dependency detection.
func (g *DependencyGraph) AddCircularDependency(cr CircularRef) {
	g.CircularDependencies = append(g.CircularDependencies, cr)
}
