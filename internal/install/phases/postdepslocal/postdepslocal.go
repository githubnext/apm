// Package postdepslocal handles stale cleanup and lockfile persistence for
// local .apm/ content after the dependency integration phase.
// Mirrors src/apm_cli/install/phases/post_deps_local.py.
package postdepslocal

import "sort"

// LocalContentState holds the inputs and mutable outputs for this phase.
type LocalContentState struct {
	// LocalDeployedFiles is the list of files deployed by the local content
	// integration; mutated to append failed-cleanup paths.
	LocalDeployedFiles []string
	// OldLocalDeployed is the list from the pre-install lockfile.
	OldLocalDeployed []string
	// LocalContentErrorsBefore is the diagnostics error count before local
	// content integration started (used to detect new errors).
	LocalContentErrorsBefore int
	// CurrentErrorCount is the total diagnostics error count after integration.
	CurrentErrorCount int
}

// HasLocalContentErrors returns true when new errors occurred during local
// content integration.
func HasLocalContentErrors(s LocalContentState) bool {
	return s.CurrentErrorCount > s.LocalContentErrorsBefore
}

// DetectStaleLocalFiles returns files in OldLocalDeployed not present in
// LocalDeployedFiles, subject to the error guard.
func DetectStaleLocalFiles(s LocalContentState) []string {
	if HasLocalContentErrors(s) {
		return nil
	}
	if len(s.OldLocalDeployed) == 0 {
		return nil
	}
	newSet := make(map[string]bool, len(s.LocalDeployedFiles))
	for _, f := range s.LocalDeployedFiles {
		newSet[f] = true
	}
	var stale []string
	for _, f := range s.OldLocalDeployed {
		if !newSet[f] {
			stale = append(stale, f)
		}
	}
	return stale
}

// SortedLocalDeployedFiles returns a sorted copy of the deployed files for
// lockfile serialisation.
func SortedLocalDeployedFiles(files []string) []string {
	cp := make([]string, len(files))
	copy(cp, files)
	sort.Strings(cp)
	return cp
}

// ShouldRun returns false when the phase should be skipped (non-PROJECT scope
// or no local content to process).
func ShouldRun(isProjectScope bool, hasLocalContent bool, hasOldLocalContent bool) bool {
	if !isProjectScope {
		return false
	}
	return hasLocalContent || hasOldLocalContent
}
