// Package cleanup orchestrates orphan and stale-file removal during install.
// Mirrors src/apm_cli/install/phases/cleanup.py.
package cleanup

// CleanupResult summarises the outcome of a stale-file cleanup pass.
type CleanupResult struct {
	Deleted        []string
	DeletedTargets []string
	Failed         []string
	SkippedUserEdit []string
}

// OrphanCleanupConfig holds the inputs for the orphan cleanup pass.
type OrphanCleanupConfig struct {
	// ExistingLockDeps maps dep_key -> deployed_files for deps in the prior lockfile.
	ExistingLockDeps map[string][]string
	// IntendedDepKeys is the set of dep keys still present in the manifest.
	IntendedDepKeys map[string]bool
	// SelfKey is the special lockfile self-entry key to skip.
	SelfKey string
}

// StaleCleanupConfig holds the inputs for the intra-package stale-file cleanup.
type StaleCleanupConfig struct {
	// OldDeployedFiles maps dep_key -> previously deployed files.
	OldDeployedFiles map[string][]string
	// NewDeployedFiles maps dep_key -> newly deployed files from integration.
	NewDeployedFiles map[string][]string
	// PackageErrorCounts maps dep_key -> count of errors during integration.
	PackageErrorCounts map[string]int
}

// DetectStaleFiles returns the set of paths that were deployed before but are
// not in the new deployment set.
func DetectStaleFiles(oldFiles, newFiles []string) []string {
	newSet := make(map[string]bool, len(newFiles))
	for _, f := range newFiles {
		newSet[f] = true
	}
	var stale []string
	for _, f := range oldFiles {
		if !newSet[f] {
			stale = append(stale, f)
		}
	}
	return stale
}

// CollectOrphanKeys returns dep keys in the existing lockfile that are no
// longer in the intended set (i.e. removed from the manifest).
func CollectOrphanKeys(cfg OrphanCleanupConfig) []string {
	var orphans []string
	for key := range cfg.ExistingLockDeps {
		if key == cfg.SelfKey {
			continue
		}
		if cfg.IntendedDepKeys[key] {
			continue
		}
		if len(cfg.ExistingLockDeps[key]) == 0 {
			continue
		}
		orphans = append(orphans, key)
	}
	return orphans
}

// CollectStalePerPackage returns, for each dep still in the manifest, the
// files that should be removed (present in old but not in new deployment).
// Packages with integration errors this run are skipped.
func CollectStalePerPackage(cfg StaleCleanupConfig) map[string][]string {
	result := map[string][]string{}
	for depKey, newDeployed := range cfg.NewDeployedFiles {
		if cfg.PackageErrorCounts[depKey] > 0 {
			continue
		}
		oldDeployed := cfg.OldDeployedFiles[depKey]
		if len(oldDeployed) == 0 {
			continue
		}
		stale := DetectStaleFiles(oldDeployed, newDeployed)
		if len(stale) > 0 {
			result[depKey] = stale
		}
	}
	return result
}
