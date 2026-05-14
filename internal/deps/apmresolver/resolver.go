// Package apmresolver implements the APM dependency resolution engine.
//
// Provides BFS-based dependency resolution, circular dependency detection,
// and dependency flattening following an NPM-hoisting "first-wins" strategy.
//
// Migrated from: src/apm_cli/deps/apm_resolver.py
package apmresolver

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/githubnext/apm/internal/deps/depgraph"
	"github.com/githubnext/apm/internal/models/depreference"
)

const defaultResolveParallel = 4

// DownloadFunc is a callback invoked to download a missing dependency.
// It mirrors the Python DownloadCallback protocol.
// Parameters:
//   - ref: the dependency reference to download
//   - apmModulesDir: the apm_modules directory path
//   - parentChain: breadcrumb string (e.g. "root > mid > dep")
//   - parentPkg: the package that declared this dependency, or ""
//
// Returns the install path on success, or "" on failure.
type DownloadFunc func(ref *depreference.DependencyReference, apmModulesDir, parentChain, parentPkg string) string

// workItem is the unit of work dispatched during the BFS download phase.
type workItem struct {
	node       *depgraph.DependencyNode
	depRef     *depreference.DependencyReference
	parentNode *depgraph.DependencyNode
	isDev      bool
}

// workResult is returned by the worker goroutine.
type workResult struct {
	item      workItem
	installed bool
	err       string
}

// Resolver resolves APM dependencies recursively.
type Resolver struct {
	maxDepth         int
	apmModulesDir    string
	projectRoot      string
	downloadFn       DownloadFunc
	maxParallel      int

	mu                      sync.Mutex
	downloadedPackages      map[string]bool
	rejectedRemoteLocalKeys map[string]bool
	callbackFailures        map[string]string
}

// Options for constructing a Resolver.
type Options struct {
	// MaxDepth is the maximum resolution depth (default: 50).
	MaxDepth int
	// ApmModulesDir is an explicit apm_modules directory path (optional).
	ApmModulesDir string
	// DownloadFn is invoked when a transitive dep is not installed.
	DownloadFn DownloadFunc
	// MaxParallel controls the worker pool size for the BFS level batches.
	// 0 or negative falls back to the APM_RESOLVE_PARALLEL env var, then
	// to defaultResolveParallel.
	MaxParallel int
}
// New creates a Resolver with the given options.
func New(opts Options) *Resolver {
	maxDepth := opts.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 50
	}
	return &Resolver{
		maxDepth:                maxDepth,
		apmModulesDir:           opts.ApmModulesDir,
		downloadFn:              opts.DownloadFn,
		maxParallel:             resolveMaxParallel(opts.MaxParallel),
		downloadedPackages:      make(map[string]bool),
		rejectedRemoteLocalKeys: make(map[string]bool),
		callbackFailures:        make(map[string]string),
	}
}

func resolveMaxParallel(explicit int) int {
	if explicit > 0 {
		return explicit
	}
	if env := strings.TrimSpace(os.Getenv("APM_RESOLVE_PARALLEL")); env != "" {
		if n, err := strconv.Atoi(env); err == nil && n > 0 {
			return n
		}
	}
	return defaultResolveParallel
}

// ResolveDependencies performs a full BFS dependency resolution starting from
// the apm.yml in projectRoot.
func (r *Resolver) ResolveDependencies(projectRoot string) *depgraph.DependencyGraph {
	r.projectRoot = projectRoot
	if r.apmModulesDir == "" {
		r.apmModulesDir = filepath.Join(projectRoot, "apm_modules")
	}

	apmYMLPath := filepath.Join(projectRoot, "apm.yml")
	if _, err := os.Stat(apmYMLPath); os.IsNotExist(err) {
		g := depgraph.NewDependencyGraph("unknown")
		return g
	}

	tree := r.buildDependencyTree(apmYMLPath)
	circularDeps := r.detectCircularDependencies(tree)
	flattened := r.flattenDependencies(tree)

	g := depgraph.NewDependencyGraph(filepath.Base(projectRoot))
	g.Tree = tree
	g.Flattened = flattened
	for _, c := range circularDeps {
		g.AddCircularDependency(c)
	}
	return g
}

// buildDependencyTree performs BFS expansion of the dependency tree.
func (r *Resolver) buildDependencyTree(rootApmYML string) *depgraph.DependencyTree {
	tree := depgraph.NewDependencyTree()

	// Read root package dependencies from apm.yml using a simple line scanner.
	deps := r.readApmYMLDeps(rootApmYML)

	// BFS queue: (depRef, parentNode, depth, isDev)
	type queueItem struct {
		ref    *depreference.DependencyReference
		parent *depgraph.DependencyNode
		depth  int
		isDev  bool
	}

	var queue []queueItem
	for _, d := range deps {
		dCopy := d
		queue = append(queue, queueItem{ref: dCopy, parent: nil, depth: 1, isDev: false})
	}

	visited := make(map[string]bool)

	for len(queue) > 0 {
		// Collect all items at the current depth level for parallel dispatch.
		currentDepth := queue[0].depth
		var level []queueItem
		remaining := queue[:0]
		for _, qi := range queue {
			if qi.depth == currentDepth {
				level = append(level, qi)
			} else {
				remaining = append(remaining, qi)
			}
		}
		queue = remaining

		// Deduplicate within the level and filter already-visited.
		var work []workItem
		for _, qi := range level {
			key := qi.ref.GetUniqueKey()
			if visited[key] {
				continue
			}
			if qi.depth > r.maxDepth {
				continue
			}
			node := &depgraph.DependencyNode{
				Ref: depgraph.DependencyRef{
					RepoURL:     qi.ref.RepoURL,
					Reference:   qi.ref.Reference,
					UniqueKey:   key,
					VirtualPath: qi.ref.VirtualPath,
					DisplayName: qi.ref.GetDisplayName(),
				},
				Depth:  qi.depth,
				Parent: qi.parent,
				IsDev:  qi.isDev,
			}
			if qi.parent != nil {
				qi.parent.Children = append(qi.parent.Children, node)
			}
			tree.AddNode(node)
			visited[key] = true
			work = append(work, workItem{
				node:       node,
				depRef:     qi.ref,
				parentNode: qi.parent,
				isDev:      qi.isDev,
			})
		}

		if len(work) == 0 {
			continue
		}

		// Dispatch work items (potentially in parallel).
		results := r.dispatchLevel(work)

		// For each successfully loaded package, enqueue its transitive deps.
		for _, res := range results {
			if !res.installed {
				if res.err != "" {
					r.mu.Lock()
					r.callbackFailures[res.item.depRef.GetUniqueKey()] = res.err
					r.mu.Unlock()
				}
				continue
			}
			// Load transitive deps from the installed package.
			installPath := r.resolveInstallPath(res.item.depRef)
			if installPath == "" {
				continue
			}
			transApmYML := filepath.Join(installPath, "apm.yml")
			if _, err := os.Stat(transApmYML); err != nil {
				continue
			}
			transDeps := r.readApmYMLDeps(transApmYML)
			for _, td := range transDeps {
				tdCopy := td
				queue = append(queue, queueItem{
					ref:    tdCopy,
					parent: res.item.node,
					depth:  res.item.node.Depth + 1,
					isDev:  res.item.isDev,
				})
			}
		}
	}

	return tree
}

// dispatchLevel runs workItems, using a goroutine pool if maxParallel > 1.
func (r *Resolver) dispatchLevel(items []workItem) []workResult {
	results := make([]workResult, len(items))

	if r.maxParallel <= 1 || r.downloadFn == nil {
		for i, item := range items {
			results[i] = r.processWorkItem(item)
		}
		return results
	}

	sem := make(chan struct{}, r.maxParallel)
	var wg sync.WaitGroup
	for i, item := range items {
		wg.Add(1)
		go func(idx int, wi workItem) {
			defer wg.Done()
			sem <- struct{}{}
			results[idx] = r.processWorkItem(wi)
			<-sem
		}(i, item)
	}
	wg.Wait()
	return results
}

func (r *Resolver) processWorkItem(item workItem) workResult {
	if r.downloadFn == nil {
		// No downloader -- check if already installed.
		installPath := r.resolveInstallPath(item.depRef)
		installed := installPath != ""
		return workResult{item: item, installed: installed}
	}

	key := item.depRef.GetUniqueKey()
	r.mu.Lock()
	alreadyDownloaded := r.downloadedPackages[key]
	r.mu.Unlock()
	if alreadyDownloaded {
		return workResult{item: item, installed: true}
	}

	parentChain := ""
	if item.node != nil {
		parentChain = item.node.GetAncestorChain()
	}
	parentPkg := ""
	if item.parentNode != nil {
		parentPkg = item.parentNode.Ref.UniqueKey
	}

	result := r.downloadFn(item.depRef, r.apmModulesDir, parentChain, parentPkg)
	if result == "" {
		return workResult{item: item, installed: false, err: "download returned empty path"}
	}

	r.mu.Lock()
	r.downloadedPackages[key] = true
	r.mu.Unlock()
	return workResult{item: item, installed: true}
}

// resolveInstallPath returns the installation path for a dependency, or "".
func (r *Resolver) resolveInstallPath(ref *depreference.DependencyReference) string {
	key := ref.GetUniqueKey()
	// Normalize: use last path segment as dir name.
	parts := strings.Split(key, "/")
	name := parts[len(parts)-1]
	candidate := filepath.Join(r.apmModulesDir, name)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

// readApmYMLDeps reads dependency references from an apm.yml file using a
// minimal line-scanner (no external YAML library required).
func (r *Resolver) readApmYMLDeps(apmYMLPath string) []*depreference.DependencyReference {
	data, err := os.ReadFile(apmYMLPath)
	if err != nil {
		return nil
	}
	return parseApmYMLDeps(string(data))
}

// parseApmYMLDeps extracts dependency strings from apm.yml content and parses
// each into a DependencyReference.
func parseApmYMLDeps(content string) []*depreference.DependencyReference {
	var refs []*depreference.DependencyReference
	inDeps := false
	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimRight(rawLine, "\r")
		trimmed := strings.TrimSpace(line)

		if trimmed == "dependencies:" || trimmed == "devDependencies:" {
			inDeps = true
			continue
		}
		if inDeps {
			// End of the deps section: a non-indented, non-list line.
			if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && trimmed != "" && !strings.HasPrefix(trimmed, "-") {
				inDeps = false
				continue
			}
			if strings.HasPrefix(trimmed, "-") {
				raw := strings.TrimPrefix(trimmed, "-")
				raw = strings.TrimSpace(raw)
				// Strip inline comments.
				if idx := strings.Index(raw, " #"); idx >= 0 {
					raw = strings.TrimSpace(raw[:idx])
				}
				// Strip surrounding quotes.
				raw = strings.Trim(raw, `"'`)
				if raw != "" {
					ref, err := depreference.Parse(raw)
					if err == nil {
						refs = append(refs, ref)
					}
				}
			}
		}
	}
	return refs
}

// detectCircularDependencies performs DFS cycle detection on the tree.
func (r *Resolver) detectCircularDependencies(tree *depgraph.DependencyTree) []depgraph.CircularRef {
	var cycles []depgraph.CircularRef
	visited := make(map[string]bool)
	var currentPath []string
	currentPathSet := make(map[string]bool)

	var dfs func(node *depgraph.DependencyNode)
	dfs = func(node *depgraph.DependencyNode) {
		nodeID := node.GetID()
		uniqueKey := node.Ref.UniqueKey

		if currentPathSet[uniqueKey] {
			// Cycle detected.
			startIdx := -1
			for i, k := range currentPath {
				if k == uniqueKey {
					startIdx = i
					break
				}
			}
			if startIdx >= 0 {
				cyclePath := append([]string{}, currentPath[startIdx:]...)
				cyclePath = append(cyclePath, uniqueKey)
				cycles = append(cycles, depgraph.CircularRef{
					CyclePath:       cyclePath,
					DetectedAtDepth: node.Depth,
				})
			}
			return
		}

		visited[nodeID] = true
		currentPath = append(currentPath, uniqueKey)
		currentPathSet[uniqueKey] = true

		for _, child := range node.Children {
			childID := child.GetID()
			if !visited[childID] || currentPathSet[child.Ref.UniqueKey] {
				dfs(child)
			}
		}

		// Backtrack.
		currentPath = currentPath[:len(currentPath)-1]
		delete(currentPathSet, uniqueKey)
	}

	for _, node := range tree.GetNodesAtDepth(1) {
		if !visited[node.GetID()] {
			currentPath = nil
			currentPathSet = make(map[string]bool)
			dfs(node)
		}
	}
	return cycles
}

// flattenDependencies flattens the tree using BFS breadth-first, first-wins
// conflict resolution (NPM hoisting).
func (r *Resolver) flattenDependencies(tree *depgraph.DependencyTree) *depgraph.FlatDependencyMap {
	flat := depgraph.NewFlatDependencyMap()
	seen := make(map[string]bool)

	for depth := 1; depth <= tree.MaxDepth; depth++ {
		nodes := tree.GetNodesAtDepth(depth)
		// Deterministic ordering.
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].GetID() < nodes[j].GetID()
		})
		for _, node := range nodes {
			key := node.Ref.UniqueKey
			if !seen[key] {
				flat.AddDependency(node.Ref, false)
				seen[key] = true
			} else {
				flat.AddDependency(node.Ref, true)
			}
		}
	}
	return flat
}
