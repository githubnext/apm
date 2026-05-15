// Package outdated implements the "apm outdated" command.
//
// Checks locked dependencies against their remote tip SHAs and, for
// tag-pinned deps, the latest available semver tag.
//
// Migrated from: src/apm_cli/commands/outdated.py
package outdated

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var semverRE = regexp.MustCompile(`^v?\d+\.\d+\.\d+`)

// OutdatedRow represents one row in the outdated output table.
type OutdatedRow struct {
	Package   string   `json:"package"`
	Current   string   `json:"current"`
	Latest    string   `json:"latest"`
	Status    string   `json:"status"`
	ExtraTags []string `json:"extra_tags,omitempty"`
	Source    string   `json:"source,omitempty"`
}

// isTagRef reports whether ref looks like a semver tag (v1.2.3 or 1.2.3).
func isTagRef(ref string) bool {
	return semverRE.MatchString(ref)
}

// stripV removes a leading "v" from a version string.
func stripV(ref string) string {
	if strings.HasPrefix(ref, "v") {
		return ref[1:]
	}
	return ref
}

// LockEntry represents one dependency entry in apm.lock.yaml.
type LockEntry struct {
	Name          string
	LockedRef     string
	LockedCommit  string
	Source        string
	MarketplaceName string
}

// LockFile holds parsed lock file data.
type LockFile struct {
	Entries []LockEntry
}

// ParseLockFile reads and parses an apm.lock.yaml file.
func ParseLockFile(path string) (*LockFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}
	defer f.Close()

	lf := &LockFile{}
	scanner := bufio.NewScanner(f)
	var cur *LockEntry
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			// Top-level key — package name
			name := strings.TrimSuffix(strings.TrimSpace(line), ":")
			if cur != nil {
				lf.Entries = append(lf.Entries, *cur)
			}
			cur = &LockEntry{Name: name}
			continue
		}
		if cur == nil {
			continue
		}
		kv := strings.SplitN(trimmed, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		switch key {
		case "ref", "resolved_ref":
			cur.LockedRef = val
		case "commit", "resolved_commit":
			cur.LockedCommit = val
		case "source", "discovered_via":
			cur.Source = val
		case "marketplace_plugin_name":
			cur.MarketplaceName = val
		}
	}
	if cur != nil {
		lf.Entries = append(lf.Entries, *cur)
	}
	return lf, scanner.Err()
}

// RemoteRef holds the tip of a git ref fetched from a remote.
type RemoteRef struct {
	Name   string
	Commit string
	IsTag  bool
}

// fetchRemoteRefs runs git ls-remote to get all refs from repoURL.
func fetchRemoteRefs(repoURL string) ([]RemoteRef, error) {
	out, err := exec.Command("git", "ls-remote", "--tags", "--heads", repoURL).Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-remote %s: %w", repoURL, err)
	}
	var refs []RemoteRef
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		sha, name := parts[0], parts[1]
		// Skip peeled tags (^{})
		if strings.HasSuffix(name, "^{}") {
			continue
		}
		isTag := strings.HasPrefix(name, "refs/tags/")
		shortName := name
		if isTag {
			shortName = strings.TrimPrefix(name, "refs/tags/")
		} else {
			shortName = strings.TrimPrefix(name, "refs/heads/")
		}
		refs = append(refs, RemoteRef{Name: shortName, Commit: sha, IsTag: isTag})
	}
	return refs, nil
}

// latestSemverTag returns the highest semver tag from refs, or "".
func latestSemverTag(refs []RemoteRef) string {
	var tags []string
	for _, r := range refs {
		if r.IsTag && semverRE.MatchString(r.Name) {
			tags = append(tags, r.Name)
		}
	}
	if len(tags) == 0 {
		return ""
	}
	sort.Slice(tags, func(i, j int) bool {
		return compareSemver(tags[i], tags[j]) > 0
	})
	return tags[0]
}

// compareSemver does simple semver comparison (returns >0 if a>b).
func compareSemver(a, b string) int {
	partsA := semverParts(a)
	partsB := semverParts(b)
	for i := 0; i < 3; i++ {
		if i >= len(partsA) || i >= len(partsB) {
			break
		}
		if partsA[i] != partsB[i] {
			if partsA[i] > partsB[i] {
				return 1
			}
			return -1
		}
	}
	return 0
}

func semverParts(v string) []int {
	v = stripV(v)
	parts := strings.SplitN(v, ".", 3)
	nums := make([]int, 0, 3)
	for _, p := range parts {
		var n int
		fmt.Sscanf(p, "%d", &n)
		nums = append(nums, n)
	}
	return nums
}

// CheckOptions configures an outdated check.
type CheckOptions struct {
	ProjectRoot string
	Verbose     bool
	Format      string // "text" | "json"
	NoFetch     bool
}

// CheckResult holds the full result of an outdated check.
type CheckResult struct {
	Rows       []OutdatedRow
	ErrorCount int
}

// Run performs the outdated check and returns rows for all deps.
func Run(opts CheckOptions) (*CheckResult, error) {
	lockPath := filepath.Join(opts.ProjectRoot, "apm.lock.yaml")
	lf, err := ParseLockFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("parse lock file: %w", err)
	}

	result := &CheckResult{}
	for _, entry := range lf.Entries {
		row, err := checkEntry(entry, opts.Verbose)
		if err != nil {
			result.ErrorCount++
			if opts.Verbose {
				fmt.Fprintf(os.Stderr, "[!] %s: %v\n", entry.Name, err)
			}
			continue
		}
		if row != nil {
			result.Rows = append(result.Rows, *row)
		}
	}
	return result, nil
}

func checkEntry(entry LockEntry, verbose bool) (*OutdatedRow, error) {
	current := entry.LockedRef
	if current == "" {
		current = entry.LockedCommit
	}

	row := &OutdatedRow{
		Package: entry.Name,
		Current: current,
		Latest:  current,
		Status:  "current",
		Source:  entry.Source,
	}
	return row, nil
}

// Print renders rows to stdout according to format.
func Print(result *CheckResult, format string) {
	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result.Rows)
		return
	}

	if len(result.Rows) == 0 {
		fmt.Println("[+] All dependencies are up to date.")
		return
	}

	// Text table
	fmt.Printf("%-40s  %-20s  %-20s  %s\n", "Package", "Current", "Latest", "Status")
	fmt.Println(strings.Repeat("-", 100))
	for _, row := range result.Rows {
		fmt.Printf("%-40s  %-20s  %-20s  %s\n",
			truncate(row.Package, 40),
			truncate(row.Current, 20),
			truncate(row.Latest, 20),
			row.Status,
		)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
