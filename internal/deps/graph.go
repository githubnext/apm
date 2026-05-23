// Package deps implements dependency graph data structures for APM.
// Mirrors src/apm_cli/deps/dependency_graph.py.
package deps

// DependencyNode represents a single dependency node in the dependency graph.
// Mirrors src/apm_cli/deps/dependency_graph.py:DependencyNode.
type DependencyNode struct {
	RepoURL   string
	Reference string // git ref (branch/tag/commit), empty means default
	Depth     int
	Children  []*DependencyNode
	Parent    *DependencyNode
	IsDev     bool
}

// GetID returns a unique identifier for this node.
// Mirrors DependencyNode.get_id().
func (n *DependencyNode) GetID() string {
	if n.Reference != "" {
		return n.RepoURL + "#" + n.Reference
	}
	return n.RepoURL
}

// GetAncestorChain builds a breadcrumb from this node's ancestry.
// Mirrors DependencyNode.get_ancestor_chain().
func (n *DependencyNode) GetAncestorChain() string {
	var parts []string
	current := n
	for current != nil {
		parts = append(parts, current.RepoURL)
		current = current.Parent
	}
	// reverse
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
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

// CircularRef represents a circular dependency reference.
// Mirrors src/apm_cli/deps/dependency_graph.py:CircularRef.
type CircularRef struct {
	CyclePath        []string
	DetectedAtDepth  int
}

// String formats the circular dependency for display.
func (c *CircularRef) String() string {
	if len(c.CyclePath) == 0 {
		return "Circular dependency detected: (empty path)"
	}
	result := "Circular dependency detected: "
	for i, p := range c.CyclePath {
		if i > 0 {
			result += " -> "
		}
		result += p
	}
	if len(c.CyclePath) > 1 && c.CyclePath[0] != c.CyclePath[len(c.CyclePath)-1] {
		result += " -> " + c.CyclePath[0]
	}
	return result
}

// ConflictInfo describes a dependency version conflict.
// Mirrors src/apm_cli/deps/dependency_graph.py:ConflictInfo.
type ConflictInfo struct {
	RepoURL   string
	WinnerRef string // reference string of the winning dep
	Conflicts []string
	Reason    string
}

// FlatDependencyMap is the final flattened dependency mapping ready for install.
// Mirrors src/apm_cli/deps/dependency_graph.py:FlatDependencyMap.
type FlatDependencyMap struct {
	Dependencies  map[string]string // unique_key -> resolved ref
	Conflicts     []ConflictInfo
	InstallOrder  []string
}

// NewFlatDependencyMap creates an empty FlatDependencyMap.
func NewFlatDependencyMap() *FlatDependencyMap {
	return &FlatDependencyMap{
		Dependencies: make(map[string]string),
	}
}

// AddDependency adds a dependency to the flat map (first-wins on conflict).
func (f *FlatDependencyMap) AddDependency(uniqueKey, ref string) {
	if _, exists := f.Dependencies[uniqueKey]; !exists {
		f.Dependencies[uniqueKey] = ref
		f.InstallOrder = append(f.InstallOrder, uniqueKey)
	}
}

// HasConflicts reports whether any conflicts were recorded.
func (f *FlatDependencyMap) HasConflicts() bool {
	return len(f.Conflicts) > 0
}

// TotalDependencies returns the count of unique dependencies.
func (f *FlatDependencyMap) TotalDependencies() int {
	return len(f.Dependencies)
}

// GetInstallationList returns dependency keys in install order.
func (f *FlatDependencyMap) GetInstallationList() []string {
	result := make([]string, 0, len(f.InstallOrder))
	for _, key := range f.InstallOrder {
		if _, ok := f.Dependencies[key]; ok {
			result = append(result, key)
		}
	}
	return result
}

// DependencyTree is the hierarchical representation before flattening.
// Mirrors src/apm_cli/deps/dependency_graph.py:DependencyTree.
type DependencyTree struct {
	RootPackage string
	Nodes       map[string]*DependencyNode
	MaxDepth    int
}

// NewDependencyTree creates an empty DependencyTree for the given root.
func NewDependencyTree(rootPackage string) *DependencyTree {
	return &DependencyTree{
		RootPackage: rootPackage,
		Nodes:       make(map[string]*DependencyNode),
	}
}

// AddNode adds a node to the tree.
func (t *DependencyTree) AddNode(node *DependencyNode) {
	id := node.GetID()
	t.Nodes[id] = node
	if node.Depth > t.MaxDepth {
		t.MaxDepth = node.Depth
	}
}

// GetNode retrieves a node by its unique key.
func (t *DependencyTree) GetNode(id string) *DependencyNode {
	return t.Nodes[id]
}

// HasDependency checks if a repo URL is present in the tree.
func (t *DependencyTree) HasDependency(repoURL string) bool {
	for _, n := range t.Nodes {
		if n.RepoURL == repoURL {
			return true
		}
	}
	return false
}

// DependencyGraph is the complete resolved dependency information.
// Mirrors src/apm_cli/deps/dependency_graph.py:DependencyGraph.
type DependencyGraph struct {
	RootPackage          string
	Tree                 *DependencyTree
	Flattened            *FlatDependencyMap
	CircularDependencies []CircularRef
	ResolutionErrors     []string
}

// NewDependencyGraph creates a new empty DependencyGraph.
func NewDependencyGraph(rootPackage string) *DependencyGraph {
	return &DependencyGraph{
		RootPackage: rootPackage,
		Tree:        NewDependencyTree(rootPackage),
		Flattened:   NewFlatDependencyMap(),
	}
}

// HasCircularDependencies reports whether any circular deps were detected.
func (g *DependencyGraph) HasCircularDependencies() bool {
	return len(g.CircularDependencies) > 0
}

// HasConflicts reports whether any conflicts exist.
func (g *DependencyGraph) HasConflicts() bool {
	return g.Flattened.HasConflicts()
}

// HasErrors reports whether any resolution errors exist.
func (g *DependencyGraph) HasErrors() bool {
	return len(g.ResolutionErrors) > 0
}

// IsValid reports whether the graph is free of circular deps and errors.
func (g *DependencyGraph) IsValid() bool {
	return !g.HasCircularDependencies() && !g.HasErrors()
}

// AddError appends a resolution error.
func (g *DependencyGraph) AddError(err string) {
	g.ResolutionErrors = append(g.ResolutionErrors, err)
}

// AddCircularDependency records a circular dependency detection.
func (g *DependencyGraph) AddCircularDependency(ref CircularRef) {
	g.CircularDependencies = append(g.CircularDependencies, ref)
}

// GetSummary returns a summary map of the dependency graph.
func (g *DependencyGraph) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"root_package":             g.RootPackage,
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
