package install

import (
	"fmt"
	"strings"
)

// Action constants for PlanEntry.
const (
	ActionUpdate    = "update"
	ActionAdd       = "add"
	ActionRemove    = "remove"
	ActionUnchanged = "unchanged"
)

// actionSymbols maps actions to ASCII bracket symbols.
var actionSymbols = map[string]string{
	ActionUpdate:    "[~]",
	ActionAdd:       "[+]",
	ActionRemove:    "[-]",
	ActionUnchanged: "[=]",
}

// PlanEntry captures one dependency's before/after state in an UpdatePlan.
type PlanEntry struct {
	DepKey      string
	Action      string
	DisplayName string

	OldResolvedRef    string
	OldResolvedCommit string
	OldContentHash    string

	NewResolvedRef    string
	NewResolvedCommit string

	DeployedFiles []string
}

// HasChanges returns true when the entry represents a real change.
func (e PlanEntry) HasChanges() bool { return e.Action != ActionUnchanged }

// ShortOldCommit returns a 7-char prefix of OldResolvedCommit, or "-".
func (e PlanEntry) ShortOldCommit() string { return shortSHA(e.OldResolvedCommit) }

// ShortNewCommit returns a 7-char prefix of NewResolvedCommit, or "-".
func (e PlanEntry) ShortNewCommit() string { return shortSHA(e.NewResolvedCommit) }

func shortSHA(commit string) string {
	if commit == "" {
		return "-"
	}
	if len(commit) >= 7 {
		return commit[:7]
	}
	return commit
}

// UpdatePlan is a structured diff between the current lockfile and a fresh resolution.
type UpdatePlan struct {
	Entries []PlanEntry
}

// HasChanges returns true when any entry has a real change.
func (p UpdatePlan) HasChanges() bool {
	for _, e := range p.Entries {
		if e.HasChanges() {
			return true
		}
	}
	return false
}

// ChangedEntries returns entries that represent changes.
func (p UpdatePlan) ChangedEntries() []PlanEntry {
	var out []PlanEntry
	for _, e := range p.Entries {
		if e.HasChanges() {
			out = append(out, e)
		}
	}
	return out
}

// SummaryCounts returns action -> count map.
func (p UpdatePlan) SummaryCounts() map[string]int {
	counts := map[string]int{
		ActionUpdate:    0,
		ActionAdd:       0,
		ActionRemove:    0,
		ActionUnchanged: 0,
	}
	for _, e := range p.Entries {
		counts[e.Action]++
	}
	return counts
}

// RenderPlanText renders the UpdatePlan as ASCII terminal output.
// Returns empty string when plan has no changes and verbose is false.
func RenderPlanText(plan UpdatePlan, verbose bool) string {
	if !plan.HasChanges() && !verbose {
		return ""
	}

	var lines []string
	lines = append(lines, "[i] Update plan for apm.yml", "")

	for _, entry := range plan.Entries {
		if entry.Action == ActionUnchanged && !verbose {
			continue
		}
		symbol := actionSymbols[entry.Action]
		if symbol == "" {
			symbol = "[?]"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", symbol, entry.DisplayName))
		lines = append(lines, fmt.Sprintf("      ref: %s", formatRefChange(entry)))
		if len(entry.DeployedFiles) > 0 {
			preview := strings.Join(entry.DeployedFiles[:min3(len(entry.DeployedFiles), 3)], ", ")
			if len(entry.DeployedFiles) > 3 {
				preview += fmt.Sprintf(", +%d more", len(entry.DeployedFiles)-3)
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

func formatRefChange(entry PlanEntry) string {
	if entry.Action == ActionRemove {
		ref := entry.OldResolvedRef
		if ref == "" {
			ref = "-"
		}
		return fmt.Sprintf("%s (%s, removed)", ref, entry.ShortOldCommit())
	}
	oldRef := entry.OldResolvedRef
	if oldRef == "" {
		oldRef = "-"
	}
	newRef := entry.NewResolvedRef
	if newRef == "" {
		newRef = oldRef
	}
	var refPart string
	if oldRef == newRef {
		refPart = oldRef
	} else {
		refPart = fmt.Sprintf("%s -> %s", oldRef, newRef)
	}
	return fmt.Sprintf("%s (%s -> %s)", refPart, entry.ShortOldCommit(), entry.ShortNewCommit())
}

func min3(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LockfileEntry represents a minimal locked dependency record for plan building.
type LockfileEntry struct {
	RepoURL      string
	VirtualPath  string
	ResolvedRef  string
	ResolvedCommit string
	ContentHash  string
	DeployedFiles []string
}

// BuildUpdatePlan compares old locked entries against newly-resolved deps.
// oldEntries: map of dep_key -> LockfileEntry (nil or empty for new installs).
// resolvedDeps: list of resolved dependency infos.
func BuildUpdatePlan(oldEntries map[string]LockfileEntry, resolvedDeps []ResolvedDep) UpdatePlan {
	seen := map[string]bool{}
	var planEntries []PlanEntry

	for _, dep := range resolvedDeps {
		key := dep.Key
		seen[key] = true
		old, hasOld := oldEntries[key]
		newRef := dep.ResolvedRef
		newCommit := dep.ResolvedCommit
		displayName := dep.DisplayName
		if displayName == "" {
			displayName = key
		}

		if !hasOld {
			planEntries = append(planEntries, PlanEntry{
				DepKey:            key,
				Action:            ActionAdd,
				DisplayName:       displayName,
				NewResolvedRef:    newRef,
				NewResolvedCommit: newCommit,
			})
			continue
		}

		deployed := make([]string, len(old.DeployedFiles))
		copy(deployed, old.DeployedFiles)

		oldCommit := old.ResolvedCommit
		oldRef := old.ResolvedRef
		if oldCommit == newCommit && oldRef == newRef {
			dn := old.RepoURL
			if old.VirtualPath != "" {
				dn = old.RepoURL + "/" + old.VirtualPath
			}
			planEntries = append(planEntries, PlanEntry{
				DepKey:            key,
				Action:            ActionUnchanged,
				DisplayName:       dn,
				OldResolvedRef:    oldRef,
				OldResolvedCommit: oldCommit,
				OldContentHash:    old.ContentHash,
				NewResolvedRef:    newRef,
				NewResolvedCommit: newCommit,
				DeployedFiles:     deployed,
			})
			continue
		}

		dn := old.RepoURL
		if old.VirtualPath != "" {
			dn = old.RepoURL + "/" + old.VirtualPath
		}
		planEntries = append(planEntries, PlanEntry{
			DepKey:            key,
			Action:            ActionUpdate,
			DisplayName:       dn,
			OldResolvedRef:    oldRef,
			OldResolvedCommit: oldCommit,
			OldContentHash:    old.ContentHash,
			NewResolvedRef:    newRef,
			NewResolvedCommit: newCommit,
			DeployedFiles:     deployed,
		})
	}

	// Entries in lockfile not in resolved set -> removed
	for key, old := range oldEntries {
		if seen[key] {
			continue
		}
		dn := old.RepoURL
		if old.VirtualPath != "" {
			dn = old.RepoURL + "/" + old.VirtualPath
		}
		planEntries = append(planEntries, PlanEntry{
			DepKey:            key,
			Action:            ActionRemove,
			DisplayName:       dn,
			OldResolvedRef:    old.ResolvedRef,
			OldResolvedCommit: old.ResolvedCommit,
			OldContentHash:    old.ContentHash,
			DeployedFiles:     old.DeployedFiles,
		})
	}

	return UpdatePlan{Entries: planEntries}
}

// ResolvedDep is a minimal resolved dependency for plan building.
type ResolvedDep struct {
	Key            string
	DisplayName    string
	ResolvedRef    string
	ResolvedCommit string
}

// LockfileSatisfiesManifest checks if every direct dep in the manifest has a
// lockfile entry. Returns (satisfied, reasons).
func LockfileSatisfiesManifest(lockedKeys map[string]bool, manifestDeps []ManifestDep) (bool, []string) {
	var reasons []string
	for _, dep := range manifestDeps {
		if dep.IsLocal {
			continue
		}
		if !lockedKeys[dep.Key] {
			reasons = append(reasons, fmt.Sprintf("  - %s is declared in apm.yml but missing from apm.lock.yaml", dep.Key))
		}
	}
	return len(reasons) == 0, reasons
}

// ManifestDep is a minimal manifest dependency entry.
type ManifestDep struct {
	Key     string
	IsLocal bool
}
