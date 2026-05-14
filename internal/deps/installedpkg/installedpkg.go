// Package installedpkg defines InstalledPackage, a record of a successfully installed dependency.
package installedpkg

// InstalledPackage records a single successfully-installed dependency.
type InstalledPackage struct {
// DepRefURL is the repository URL of the installed dependency.
DepRefURL      string
ResolvedCommit string
Depth          int
ResolvedBy     string
IsDev          bool
RegistryHost   string
RegistryPrefix string
}
