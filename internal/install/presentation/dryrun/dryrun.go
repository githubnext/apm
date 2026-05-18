// Package dryrun renders the dry-run preview for apm install --dry-run.
// Mirrors src/apm_cli/install/presentation/dry_run.py.
package dryrun

import "fmt"

// Dep is a minimal dependency representation used for dry-run rendering.
type Dep interface {
	RepoURL() string
	Reference() string
	GetUniqueKey() string
}

// Logger is the subset of InstallLogger used by the dry-run renderer.
type Logger interface {
	Progress(msg string)
	DryRunNotice(msg string)
	Success(msg string)
}

// Options holds all inputs required for RenderAndExit.
type Options struct {
	Logger           Logger
	ShouldInstallAPM bool
	APMDeps          []Dep
	MCPDeps          []fmt.Stringer
	DevAPMDeps       []Dep
	ShouldInstallMCP bool
	Update           bool
	OnlyPackages     []string
	LockfileOrphans  []string // pre-computed orphan list; nil = skip
}

// RenderAndExit writes the dry-run preview to the logger.
// It does NOT exit; the caller must return after calling this function.
func RenderAndExit(opts Options) {
	opts.Logger.Progress("Dry run mode - showing what would be installed:")

	if opts.ShouldInstallAPM && len(opts.APMDeps) > 0 {
		opts.Logger.Progress(fmt.Sprintf("APM dependencies (%d):", len(opts.APMDeps)))
		for _, dep := range opts.APMDeps {
			action := "install"
			if opts.Update {
				action = "update"
			}
			ref := dep.Reference()
			if ref == "" {
				ref = "main"
			}
			opts.Logger.Progress(fmt.Sprintf("  - %s#%s -> %s", dep.RepoURL(), ref, action))
		}
	}

	if opts.ShouldInstallMCP && len(opts.MCPDeps) > 0 {
		opts.Logger.Progress(fmt.Sprintf("MCP dependencies (%d):", len(opts.MCPDeps)))
		for _, dep := range opts.MCPDeps {
			opts.Logger.Progress(fmt.Sprintf("  - %s", dep))
		}
	}

	if len(opts.APMDeps) == 0 && len(opts.DevAPMDeps) == 0 && len(opts.MCPDeps) == 0 {
		opts.Logger.Progress("No dependencies found in apm.yml")
	}

	// Orphan preview
	if len(opts.LockfileOrphans) > 0 {
		opts.Logger.Progress(
			fmt.Sprintf("Files that would be removed (packages no longer in apm.yml): %d",
				len(opts.LockfileOrphans)))
		limit := 10
		if len(opts.LockfileOrphans) < limit {
			limit = len(opts.LockfileOrphans)
		}
		for _, orphan := range opts.LockfileOrphans[:limit] {
			opts.Logger.Progress(fmt.Sprintf("  - %s", orphan))
		}
		if len(opts.LockfileOrphans) > 10 {
			opts.Logger.Progress(fmt.Sprintf("  ... and %d more", len(opts.LockfileOrphans)-10))
		}
	}

	if len(opts.APMDeps) > 0 || len(opts.DevAPMDeps) > 0 {
		opts.Logger.DryRunNotice(
			"Per-package stale-file cleanup (renames within a package) is " +
				"not previewed -- it requires running integration. Run without " +
				"--dry-run to apply.")
	}

	opts.Logger.Success("Dry run complete - no changes made")
}
