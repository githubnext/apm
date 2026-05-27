// Package install: dependency source strategy types.
//
// Each DependencySource knows how to acquire one dependency: bring its files
// onto disk, build a PackageInfo, register it in the lockfile-bound state, and
// return the metadata the integration template needs.
//
// Mirrors src/apm_cli/install/sources.py.
package install

// SourceKind identifies which acquisition strategy a DependencySource uses.
// Mirrors the three concrete classes in sources.py.
type SourceKind int

const (
	// SourceKindLocal is a file:// dependency copied from the workspace.
	SourceKindLocal SourceKind = iota
	// SourceKindCached is a dependency already extracted in apm_modules/.
	SourceKindCached
	// SourceKindFresh is a dependency that needs a network download.
	SourceKindFresh
)

// IntegrateErrorPrefix is the default error prefix used when primitive
// integration fails. Mirrors DependencySource.INTEGRATE_ERROR_PREFIX.
const IntegrateErrorPrefix = "Failed to integrate primitives"

// IntegrateErrorPrefixLocal is the error prefix for local packages.
// Mirrors LocalDependencySource.INTEGRATE_ERROR_PREFIX.
const IntegrateErrorPrefixLocal = "Failed to integrate primitives from local package"

// IntegrateErrorPrefixCached is the error prefix for cached packages.
// Mirrors CachedDependencySource.INTEGRATE_ERROR_PREFIX.
const IntegrateErrorPrefixCached = "Failed to integrate primitives from cached package"

// Materialization is the outcome of DependencySource.Acquire().
//
// Carries everything the integration template needs to run the security
// gate and primitive integration on a freshly-acquired package.
// Mirrors the Materialization dataclass in sources.py.
type Materialization struct {
	// PackageInfo is the resolved package metadata (nil when integration
	// should be skipped, e.g. no targets configured).
	PackageInfo any

	// InstallPath is the on-disk directory where the package was extracted.
	InstallPath string

	// DepKey is the lockfile key for this dependency.
	DepKey string

	// Deltas tracks install-phase counters accumulated by this package.
	// "installed" is always 1; "unpinned" is 1 when the dep has no ref.
	Deltas map[string]int
}

// NewMaterialization creates a Materialization with the default delta
// (installed:1). Mirrors the default_factory on the Materialization dataclass.
func NewMaterialization(packageInfo any, installPath, depKey string) *Materialization {
	return &Materialization{
		PackageInfo: packageInfo,
		InstallPath: installPath,
		DepKey:      depKey,
		Deltas:      map[string]int{"installed": 1},
	}
}

// DependencySource is the strategy interface: acquire one dependency and
// prepare it for integration. Mirrors the DependencySource ABC in sources.py.
//
// Implementations:
//   - LocalDependencySource (file:// deps)
//   - CachedDependencySource (already in apm_modules/)
//   - FreshDependencySource (needs network download)
type DependencySource interface {
	// Acquire materialises the dependency on disk and builds PackageInfo.
	// Returns nil to skip integration entirely (e.g. local dep at user
	// scope, copy/download failure).
	Acquire() (*Materialization, error)

	// Kind returns which acquisition strategy this source uses.
	Kind() SourceKind

	// IntegrateErrorPrefix returns the error message prefix used when
	// integrate_package_primitives raises for this source type.
	IntegrateErrorPrefix() string
}
