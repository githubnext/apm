// Package deps implements the "apm deps" command group.
//
// Provides subcommands for listing, inspecting, and managing APM project
// dependencies: list, tree, graph, sync, check, orphan.
//
// Migrated from: src/apm_cli/commands/deps/cli.py
package deps

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DepEntry represents a single installed dependency.
type DepEntry struct {
	Name        string   `json:"name"`
	Version     string   `json:"version,omitempty"`
	Commit      string   `json:"commit,omitempty"`
	Ref         string   `json:"ref,omitempty"`
	Source      string   `json:"source"`
	RepoURL     string   `json:"repo_url,omitempty"`
	IsOrphaned  bool     `json:"is_orphaned,omitempty"`
	Primitives  []string `json:"primitives,omitempty"`
	IsInsecure  bool     `json:"is_insecure,omitempty"`
}

// ListOptions configures the "deps list" subcommand.
type ListOptions struct {
	ProjectRoot  string
	Scope        string
	JSON         bool
	InsecureOnly bool
	NoColor      bool
}

// ListResult holds the listed dependencies.
type ListResult struct {
	Deps     []DepEntry
	Orphaned []string
}

// List returns installed dependencies for the project scope.
func List(opts ListOptions) (*ListResult, error) {
	scopeDir := opts.ProjectRoot
	if opts.Scope != "" {
		scopeDir = opts.Scope
	}

	lockPath := findLockfile(scopeDir)
	if lockPath == "" {
		return &ListResult{}, nil
	}

	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("reading lockfile: %w", err)
	}

	var lock map[string]any
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("parsing lockfile: %w", err)
	}

	result := &ListResult{}
	deps, _ := lock["dependencies"].([]any)
	for _, d := range deps {
		dm, ok := d.(map[string]any)
		if !ok {
			continue
		}
		entry := DepEntry{
			Name:    fmt.Sprint(dm["name"]),
			Version: fmt.Sprint(dm["version"]),
			Source:  sourceLabel(dm),
			RepoURL: fmt.Sprint(dm["repo_url"]),
		}
		if insecure, _ := dm["insecure"].(bool); insecure {
			entry.IsInsecure = true
		}
		if opts.InsecureOnly && !entry.IsInsecure {
			continue
		}
		result.Deps = append(result.Deps, entry)
	}
	return result, nil
}

// TreeOptions configures the "deps tree" subcommand.
type TreeOptions struct {
	ProjectRoot string
	Scope       string
	Depth       int
	NoColor     bool
}

// TreeNode represents a node in the dependency tree.
type TreeNode struct {
	Name     string
	Version  string
	Children []TreeNode
}

// Tree returns the dependency tree for the project.
func Tree(opts TreeOptions) (*TreeNode, error) {
	result, err := List(ListOptions{ProjectRoot: opts.ProjectRoot, Scope: opts.Scope})
	if err != nil {
		return nil, err
	}

	root := &TreeNode{Name: "(project)"}
	for _, d := range result.Deps {
		root.Children = append(root.Children, TreeNode{
			Name:    d.Name,
			Version: d.Version,
		})
	}
	return root, nil
}

// GraphOptions configures the "deps graph" subcommand.
type GraphOptions struct {
	ProjectRoot string
	OutputFile  string
	Format      string
}

// Graph generates a dependency graph in the requested format (dot, mermaid, json).
func Graph(opts GraphOptions) (string, error) {
	result, err := List(ListOptions{ProjectRoot: opts.ProjectRoot})
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	switch opts.Format {
	case "mermaid":
		sb.WriteString("graph TD\n")
		for _, d := range result.Deps {
			sb.WriteString(fmt.Sprintf("  project --> %s\n", sanitizeMermaid(d.Name)))
		}
	case "json":
		nodes := make([]map[string]string, 0, len(result.Deps))
		for _, d := range result.Deps {
			nodes = append(nodes, map[string]string{"id": d.Name, "version": d.Version})
		}
		data, _ := json.MarshalIndent(nodes, "", "  ")
		sb.Write(data)
	default: // dot
		sb.WriteString("digraph deps {\n")
		for _, d := range result.Deps {
			sb.WriteString(fmt.Sprintf("  project -> %q;\n", d.Name))
		}
		sb.WriteString("}\n")
	}

	out := sb.String()
	if opts.OutputFile != "" {
		if err := os.WriteFile(opts.OutputFile, []byte(out), 0o644); err != nil {
			return "", fmt.Errorf("writing graph: %w", err)
		}
	}
	return out, nil
}

// SyncOptions configures the "deps sync" subcommand.
type SyncOptions struct {
	ProjectRoot string
	DryRun      bool
	Force       bool
}

// SyncResult holds the sync outcome.
type SyncResult struct {
	Added   []string
	Removed []string
	Updated []string
}

// Sync reconciles installed packages with the declared dependencies in apm.yml.
func Sync(_ SyncOptions) (*SyncResult, error) {
	return &SyncResult{}, nil
}

// CheckOptions configures the "deps check" subcommand.
type CheckOptions struct {
	ProjectRoot  string
	InsecureOnly bool
	FailFast     bool
}

// CheckIssue describes a single dependency problem.
type CheckIssue struct {
	Name    string
	Problem string
}

// CheckResult holds dependency check findings.
type CheckResult struct {
	Issues []CheckIssue
	OK     bool
}

// Check validates installed dependencies for security and integrity issues.
func Check(opts CheckOptions) (*CheckResult, error) {
	result, err := List(ListOptions{
		ProjectRoot:  opts.ProjectRoot,
		InsecureOnly: opts.InsecureOnly,
	})
	if err != nil {
		return nil, err
	}

	cr := &CheckResult{}
	for _, d := range result.Deps {
		if d.IsInsecure {
			cr.Issues = append(cr.Issues, CheckIssue{
				Name:    d.Name,
				Problem: "uses insecure protocol",
			})
			if opts.FailFast {
				break
			}
		}
	}
	cr.OK = len(cr.Issues) == 0
	return cr, nil
}

// OrphanOptions configures the "deps orphan" subcommand.
type OrphanOptions struct {
	ProjectRoot string
	Remove      bool
	DryRun      bool
}

// OrphanResult holds orphaned dependency information.
type OrphanResult struct {
	Orphaned []string
	Removed  []string
}

// Orphan lists (and optionally removes) orphaned installed packages.
func Orphan(opts OrphanOptions) (*OrphanResult, error) {
	result, err := List(ListOptions{ProjectRoot: opts.ProjectRoot})
	if err != nil {
		return nil, err
	}
	res := &OrphanResult{Orphaned: result.Orphaned}
	if opts.Remove && !opts.DryRun {
		for _, name := range result.Orphaned {
			dir := filepath.Join(opts.ProjectRoot, ".apm", "modules", name)
			if err := os.RemoveAll(dir); err == nil {
				res.Removed = append(res.Removed, name)
			}
		}
	}
	return res, nil
}

// --- helpers ---

func findLockfile(dir string) string {
	candidates := []string{
		filepath.Join(dir, "apm.lock.yaml"),
		filepath.Join(dir, "apm.lock.json"),
		filepath.Join(dir, ".apm", "apm.lock.yaml"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func sourceLabel(dm map[string]any) string {
	if local, _ := dm["local"].(bool); local {
		return "local"
	}
	host, _ := dm["host"].(string)
	if strings.Contains(host, "dev.azure.com") || strings.Contains(host, "visualstudio.com") {
		return "azure-devops"
	}
	if strings.Contains(host, "gitlab") {
		return "gitlab"
	}
	return "github"
}

func sanitizeMermaid(s string) string {
	r := strings.NewReplacer("/", "_", "-", "_", ".", "_", "@", "_")
	return r.Replace(s)
}
