// Package integration provides file-level integrator interfaces and types for APM.
package integration

import "errors"

// IntegrationResult holds the outcome of a file-integration operation.
type IntegrationResult struct {
	// FilesIntegrated is the number of files written.
	FilesIntegrated int
	// FilesUpdated is kept for CLI compat, always 0 today.
	FilesUpdated int
	// FilesSkipped is the number of files that were unchanged.
	FilesSkipped int
	// LinksResolved is the number of inter-file links resolved.
	LinksResolved int
	// ScriptsCopied is the number of hook scripts copied.
	ScriptsCopied int
	// SubSkillsPromoted is the number of sub-skills promoted.
	SubSkillsPromoted int
	// SkillCreated is true if a new skill was created.
	SkillCreated bool
	// FilesAdopted is the number of byte-identical files adopted silently.
	FilesAdopted int
}

// Integrator is the base interface for file-level integrators.
type Integrator interface {
	// Integrate performs the integration and returns a result.
	Integrate(opts IntegrateOptions) (IntegrationResult, error)
	// Name returns the integrator's identifier.
	Name() string
}

// IntegrateOptions carries options for an integration run.
type IntegrateOptions struct {
	// DryRun skips writing files to disk.
	DryRun bool
	// Force overwrites even byte-identical files.
	Force bool
	// Global applies the operation globally.
	Global bool
}

// ErrIntegrationConflict is returned when a file conflict cannot be resolved.
var ErrIntegrationConflict = errors.New("integration conflict")

// ErrIntegrationSkipped is returned when integration was intentionally skipped.
var ErrIntegrationSkipped = errors.New("integration skipped")

// Total returns the total number of files touched (integrated + skipped + adopted).
func (r IntegrationResult) Total() int {
	return r.FilesIntegrated + r.FilesSkipped + r.FilesAdopted
}
