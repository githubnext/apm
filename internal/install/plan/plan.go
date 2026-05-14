// Package plan provides the update-plan diff between a current lockfile and
// a fresh resolution. Mirrors src/apm_cli/install/plan.py.
package plan

import (
	"fmt"
	"sort"
	"strings"
)

const (
	ActionUpdate    = "update"
	ActionAdd       = "add"
	ActionRemove    = "remove"
	ActionUnchanged = "unchanged"
)

// actionOrder controls the sort order in RenderPlanText.
var actionOrder = map[string]int{
	ActionUpdate:    0,
	ActionAdd:       1,
	ActionRemove:    2,
	ActionUnchanged: 3,
}

// actionSymbols maps action constants to ASCII bracket symbols.
var actionSymbols = map[string]string{
	ActionUpdate:    "[~]",
	ActionAdd:       "[+]",
	ActionRemove:    "[-]",
	ActionUnchanged: "[=]",
}

// LockedDependency carries the minimal fields the plan builder reads from a
// lock file entry.
type LockedDependency struct {
	Key           string
	RepoURL       string
	VirtualPath   string
	ResolvedRef   string
	ResolvedCommit string
	ContentHash   string
	DeployedFiles []string
}

// DependencyReference carries the minimal fields the plan builder reads from
// a resolved manifest dependency.
type DependencyReference struct {
	RepoURL     string
	LocalPath   string
	VirtualPath string
	IsLocal     bool
	IsVirtual   bool
	Reference   string // manifest ref
	// ResolvedRefName and ResolvedCommit are populated by the resolve phase.
	ResolvedRefName string
	ResolvedCommit  string
}

// depRefKey returns the unique key for a manifest dependency, mirroring the
// Python _dep_ref_key helper.
func depRefKey(dep DependencyReference) string {
	if dep.IsLocal && dep.LocalPath != "" {
		return dep.LocalPath
	}
	if dep.IsVirtual && dep.VirtualPath != "" {
		return dep.RepoURL + "/" + dep.VirtualPath
	}
	return dep.RepoURL
}

func shortSHA(commit string, length int) string {
	if commit == "" {
		return "-"
	}
	if len(commit) <= length {
		return commit
	}
	return commit[:length]
}

// PlanEntry records one dependency's before/after state.
type PlanEntry struct {
	DepKey        string
	Action        string
	DisplayName   string
	OldResolvedRef    string
	OldResolvedCommit string
	OldContentHash    string
	NewResolvedRef    string
	NewResolvedCommit string
	DeployedFiles []string
}

// HasChanges returns true when the action is not "unchanged".
func (e PlanEntry) HasChanges() bool { return e.Action != ActionUnchanged }

// ShortOldCommit returns the 7-char abbreviated old commit SHA.
func (e PlanEntry) ShortOldCommit() string { return shortSHA(e.OldResolvedCommit, 7) }

// ShortNewCommit returns the 7-char abbreviated new commit SHA.
func (e PlanEntry) ShortNewCommit() string { return shortSHA(e.NewResolvedCommit, 7) }

// UpdatePlan is the structured diff between an existing lockfile and the
// freshly resolved dependencies.
type UpdatePlan struct {
	Entries []PlanEntry
}

// HasChanges returns true when at least one entry has a change.
func (p UpdatePlan) HasChanges() bool {
	for _, e := range p.Entries {
		if e.HasChanges() {
			return true
		}
	}
	return false
}

// ChangedEntries returns only the entries that represent a change.
func (p UpdatePlan) ChangedEntries() []PlanEntry {
	var out []PlanEntry
	for _, e := range p.Entries {
		if e.HasChanges() {
			out = append(out, e)
		}
	}
	return out
}

// SummaryCounts returns counts per action string.
func (p UpdatePlan) SummaryCounts() map[string]int {
	m := map[string]int{
		ActionUpdate:    0,
		ActionAdd:       0,
		ActionRemove:    0,
		ActionUnchanged: 0,
	}
	for _, e := range p.Entries {
		m[e.Action]++
	}
	return m
}

func displayName(key string, locked *LockedDependency) string {
	if locked != nil {
		name := locked.RepoURL
		if locked.VirtualPath != "" {
			name = name + "/" + locked.VirtualPath
		}
		return name
	}
	return key
}

// BuildUpdatePlan compares an existing lockfile against freshly-resolved
// dependencies and returns an UpdatePlan.
func BuildUpdatePlan(
	oldDeps map[string]*LockedDependency,
	resolvedDeps []DependencyReference,
) UpdatePlan {
	seenKeys := map[string]bool{}
	var entries []PlanEntry

	for _, dep := range resolvedDeps {
		key := depRefKey(dep)
		seenKeys[key] = true
		old := oldDeps[key]
		newRef := dep.ResolvedRefName
		if newRef == "" {
			newRef = dep.Reference
		}
		newCommit := dep.ResolvedCommit

		if old == nil {
			entries = append(entries, PlanEntry{
				DepKey:            key,
				Action:            ActionAdd,
				DisplayName:       dep.RepoURL,
				NewResolvedRef:    newRef,
				NewResolvedCommit: newCommit,
			})
			continue
		}

		oldRef := old.ResolvedRef
		oldCommit := old.ResolvedCommit

		if (oldCommit == newCommit || (oldCommit == "" && newCommit == "")) &&
			(oldRef == newRef || (oldRef == "" && newRef == "")) {
			entries = append(entries, PlanEntry{
				DepKey:            key,
				Action:            ActionUnchanged,
				DisplayName:       displayName(key, old),
				OldResolvedRef:    oldRef,
				OldResolvedCommit: oldCommit,
				OldContentHash:    old.ContentHash,
				NewResolvedRef:    newRef,
				NewResolvedCommit: newCommit,
				DeployedFiles:     old.DeployedFiles,
			})
			continue
		}

		entries = append(entries, PlanEntry{
			DepKey:            key,
			Action:            ActionUpdate,
			DisplayName:       displayName(key, old),
			OldResolvedRef:    oldRef,
			OldResolvedCommit: oldCommit,
			OldContentHash:    old.ContentHash,
			NewResolvedRef:    newRef,
			NewResolvedCommit: newCommit,
			DeployedFiles:     old.DeployedFiles,
		})
	}

	for key, old := range oldDeps {
		if seenKeys[key] {
			continue
		}
		entries = append(entries, PlanEntry{
			DepKey:            key,
			Action:            ActionRemove,
			DisplayName:       displayName(key, old),
			OldResolvedRef:    old.ResolvedRef,
			OldResolvedCommit: old.ResolvedCommit,
			OldContentHash:    old.ContentHash,
			DeployedFiles:     old.DeployedFiles,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		oi := actionOrder[entries[i].Action]
		oj := actionOrder[entries[j].Action]
		if oi != oj {
			return oi < oj
		}
		ni := entries[i].DisplayName
		if ni == "" {
			ni = entries[i].DepKey
		}
		nj := entries[j].DisplayName
		if nj == "" {
			nj = entries[j].DepKey
		}
		return ni < nj
	})

	return UpdatePlan{Entries: entries}
}

func formatRefChange(e PlanEntry) string {
	switch e.Action {
	case ActionAdd:
		ref := e.NewResolvedRef
		if ref == "" {
			ref = "-"
		}
		return fmt.Sprintf("%s (%s, new)", ref, e.ShortNewCommit())
	case ActionRemove:
		ref := e.OldResolvedRef
		if ref == "" {
			ref = "-"
		}
		return fmt.Sprintf("%s (%s, removed)", ref, e.ShortOldCommit())
	default:
		oldRef := e.OldResolvedRef
		if oldRef == "" {
			oldRef = "-"
		}
		newRef := e.NewResolvedRef
		if newRef == "" {
			newRef = oldRef
		}
		refPart := oldRef
		if oldRef != newRef {
			refPart = oldRef + " -> " + newRef
		}
		return fmt.Sprintf("%s (%s -> %s)", refPart, e.ShortOldCommit(), e.ShortNewCommit())
	}
}

// RenderPlanText returns an ASCII rendering of the UpdatePlan suitable for
// terminal display. Returns empty string when there are no changes (and
// verbose is false).
func RenderPlanText(plan UpdatePlan, verbose bool) string {
	if !plan.HasChanges() && !verbose {
		return ""
	}

	var lines []string
	lines = append(lines, "[i] Update plan for apm.yml", "")

	for _, e := range plan.Entries {
		if e.Action == ActionUnchanged && !verbose {
			continue
		}
		sym := actionSymbols[e.Action]
		if sym == "" {
			sym = "[?]"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", sym, e.DisplayName))
		lines = append(lines, fmt.Sprintf("      ref: %s", formatRefChange(e)))
		if len(e.DeployedFiles) > 0 {
			preview := strings.Join(e.DeployedFiles[:min(3, len(e.DeployedFiles))], ", ")
			if len(e.DeployedFiles) > 3 {
				preview += fmt.Sprintf(", +%d more", len(e.DeployedFiles)-3)
			}
			lines = append(lines, fmt.Sprintf("      files: %s", preview))
		}
		lines = append(lines, "")
	}

	counts := plan.SummaryCounts()
	var summaryParts []string
	if counts[ActionUpdate] > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d updated", counts[ActionUpdate]))
	}
	if counts[ActionAdd] > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d added", counts[ActionAdd]))
	}
	if counts[ActionRemove] > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d removed", counts[ActionRemove]))
	}
	if verbose && counts[ActionUnchanged] > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("%d unchanged", counts[ActionUnchanged]))
	}
	if len(summaryParts) > 0 {
		lines = append(lines, "  "+strings.Join(summaryParts, ", "))
	}

	result := strings.Join(lines, "\n")
	return strings.TrimRight(result, "\n")
}

// LockfileSatisfiesManifest checks that every manifest dep has a lockfile entry.
// Returns (satisfied, reasons).
func LockfileSatisfiesManifest(
	lockedKeys map[string]bool,
	manifestDeps []DependencyReference,
) (bool, []string) {
	var reasons []string
	for _, dep := range manifestDeps {
		if dep.IsLocal {
			continue
		}
		key := depRefKey(dep)
		if !lockedKeys[key] {
			reasons = append(reasons, fmt.Sprintf("  - %s is declared in apm.yml but missing from apm.lock.yaml", key))
		}
	}
	return len(reasons) == 0, reasons
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
