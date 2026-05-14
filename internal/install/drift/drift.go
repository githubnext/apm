// Package drift provides pure drift-detection helpers for diff-aware apm install.
// These functions are stateless and side-effect-free.
// Migrated from src/apm_cli/drift.py
package drift

// DependencyRef is a minimal interface for dependency references.
// Implementations provide the fields compared during drift detection.
type DependencyRef interface {
	// Reference returns the git ref pinned in apm.yml (may be "").
	Reference() string
	// UniqueKey returns the canonical deduplication key (repo_url or repo_url/virtual_path).
	UniqueKey() string
	// IsInsecure returns true when the dep was declared with an insecure HTTP URL.
	IsInsecure() bool
	// Host returns the registry proxy host when set, or "".
	Host() string
	// ArtifactoryPrefix returns the Artifactory prefix when set, or "".
	ArtifactoryPrefix() string
}

// LockedDep is a minimal interface for lockfile dependency entries.
type LockedDep interface {
	// ResolvedRef returns the ref recorded in the lockfile.
	ResolvedRef() string
	// ResolvedCommit returns the commit SHA recorded in the lockfile (may be "").
	ResolvedCommit() string
	// DeployedFiles returns the list of deployed file paths.
	DeployedFiles() []string
	// IsInsecure returns the stored insecure flag.
	IsInsecure() bool
	// AllowInsecure returns the stored allow_insecure flag.
	AllowInsecure() bool
	// RegistryPrefix returns the Artifactory prefix (may be "").
	RegistryPrefix() string
	// Host returns the locked host (may be "").
	Host() string
}

// LockFile is a minimal interface for lockfile operations.
type LockFile interface {
	// Dependencies returns all locked dependencies keyed by unique key.
	Dependencies() map[string]LockedDep
	// GetDependency returns the locked entry for the given unique key (nil if absent).
	GetDependency(uniqueKey string) LockedDep
}

// RefChangeResult holds the outcome of DetectRefChange.
type RefChangeResult struct {
	Changed bool
}

// DetectRefChange reports whether the manifest ref differs from the locked resolved_ref.
//
// Returns true for transitions: ref added (""  -> "v1.0.0"),
// ref removed ("main" -> ""), ref changed ("v1.0.0" -> "v2.0.0"),
// or HTTP-insecure flag toggle.
//
// Returns false when updateRefs is true (--update mode), when lockedDep is nil
// (new package), or when the ref is unchanged.
func DetectRefChange(depRef DependencyRef, lockedDep LockedDep, updateRefs bool) bool {
	if updateRefs {
		return false
	}
	if lockedDep == nil {
		return false
	}
	if depRef.Reference() != lockedDep.ResolvedRef() {
		return true
	}
	return depRef.IsInsecure() != lockedDep.IsInsecure()
}

// DetectOrphans returns the set of deployed file paths whose owning package
// left the manifest.
//
// Only relevant for full installs (onlyPackages empty). Partial installs
// preserve all existing lockfile entries unchanged.
func DetectOrphans(existing LockFile, intendedDepKeys map[string]struct{}, onlyPackages []string) map[string]struct{} {
	orphaned := map[string]struct{}{}
	if len(onlyPackages) > 0 || existing == nil {
		return orphaned
	}
	for depKey, dep := range existing.Dependencies() {
		if _, ok := intendedDepKeys[depKey]; !ok {
			for _, f := range dep.DeployedFiles() {
				orphaned[f] = struct{}{}
			}
		}
	}
	return orphaned
}

// DetectStaleFiles returns the set of paths that were deployed previously
// but are no longer produced by the current install.
//
// Pure set-difference: set(oldDeployed) - set(newDeployed).
func DetectStaleFiles(oldDeployed, newDeployed []string) map[string]struct{} {
	newSet := make(map[string]struct{}, len(newDeployed))
	for _, f := range newDeployed {
		newSet[f] = struct{}{}
	}
	stale := map[string]struct{}{}
	for _, f := range oldDeployed {
		if _, ok := newSet[f]; !ok {
			stale[f] = struct{}{}
		}
	}
	return stale
}

// DetectConfigDrift returns names of entries whose current config differs
// from the stored baseline.
//
// Only entries with a stored baseline that has changed are returned.
// Brand-new entries (absent from storedConfigs) are excluded.
func DetectConfigDrift(currentConfigs, storedConfigs map[string]interface{}) map[string]struct{} {
	drifted := map[string]struct{}{}
	for name, current := range currentConfigs {
		stored, ok := storedConfigs[name]
		if !ok {
			continue
		}
		if !configsEqual(current, stored) {
			drifted[name] = struct{}{}
		}
	}
	return drifted
}

// configsEqual performs a deep equality check on two config values.
func configsEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	switch av := a.(type) {
	case map[string]interface{}:
		bv, ok := b.(map[string]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, va := range av {
			vb, ok := bv[k]
			if !ok || !configsEqual(va, vb) {
				return false
			}
		}
		return true
	case []interface{}:
		bv, ok := b.([]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !configsEqual(av[i], bv[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

// DownloadRefOptions controls BuildDownloadRef behavior.
type DownloadRefOptions struct {
	UpdateRefs bool
	RefChanged bool
}

// SimpleDepRef is a concrete DependencyRef implementation for use in tests
// and pipeline wiring.
type SimpleDepRef struct {
	Ref                string
	Key                string
	Insecure           bool
	HostVal            string
	ArtifactoryPfx     string
}

func (s *SimpleDepRef) Reference() string       { return s.Ref }
func (s *SimpleDepRef) UniqueKey() string        { return s.Key }
func (s *SimpleDepRef) IsInsecure() bool         { return s.Insecure }
func (s *SimpleDepRef) Host() string             { return s.HostVal }
func (s *SimpleDepRef) ArtifactoryPrefix() string { return s.ArtifactoryPfx }
