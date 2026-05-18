// Package installservice is the application service for the APM install pipeline.
// Migrated from src/apm_cli/install/service.py
package installservice

import "errors"

// InstallNotAvailableError is raised when the APM dependency subsystem is unavailable.
type InstallNotAvailableError struct {
	Cause error
}

func (e *InstallNotAvailableError) Error() string {
	if e.Cause != nil {
		return "APM install subsystem unavailable: " + e.Cause.Error()
	}
	return "APM install subsystem unavailable"
}

// FrozenInstallError is raised when the lockfile is missing or out of sync
// in a frozen install.
type FrozenInstallError struct {
	Reason string
}

func (e *FrozenInstallError) Error() string {
	if e.Reason != "" {
		return "frozen install failed: " + e.Reason
	}
	return "frozen install failed: lockfile missing or out of sync"
}

// IsFrozenInstallError reports whether err is a FrozenInstallError.
func IsFrozenInstallError(err error) bool {
	var fe *FrozenInstallError
	return errors.As(err, &fe)
}

// InstallRequest holds the parameters for one install invocation.
type InstallRequest struct {
	// Packages is the list of package specifiers to install.
	Packages []string
	// Frozen prevents resolve/download and requires the lockfile to be up-to-date.
	Frozen bool
	// UpdateRefs forces re-resolution of branch references.
	UpdateRefs bool
	// Scope restricts installation to a specific target scope.
	Scope string
	// Target overrides auto-detected integration targets.
	Target string
	// Verbose enables verbose output.
	Verbose bool
	// DryRun simulates the install without writing any files.
	DryRun bool
}

// InstallResult summarises the outcome of an install invocation.
type InstallResult struct {
	// Installed lists packages that were newly installed.
	Installed []string
	// Updated lists packages that were updated.
	Updated []string
	// Skipped lists packages that were already up-to-date.
	Skipped []string
	// Failed lists packages that could not be installed.
	Failed []string
	// ExitCode is the suggested process exit code (0 = success).
	ExitCode int
}

// InstallServicer is the interface implemented by InstallService.
type InstallServicer interface {
	Run(req *InstallRequest) (*InstallResult, error)
}

// InstallService orchestrates one install invocation.
// Stateless: a single instance can serve multiple Run calls.
type InstallService struct{}

// New creates a new InstallService.
func New() *InstallService {
	return &InstallService{}
}

// Run executes the install pipeline and returns the structured result.
// The actual pipeline implementation is injected at runtime; this
// skeleton validates inputs and returns a stub result.
func (s *InstallService) Run(req *InstallRequest) (*InstallResult, error) {
	if req == nil {
		return nil, &InstallNotAvailableError{Cause: errors.New("nil request")}
	}
	return &InstallResult{ExitCode: 0}, nil
}
