// Package install provides types and utilities for the APM install pipeline.
package install

// DirectDependencyError is raised when one or more direct dependencies fail
// validation or integration.
type DirectDependencyError struct {
	msg string
}

func (e *DirectDependencyError) Error() string { return e.msg }

// NewDirectDependencyError creates a DirectDependencyError with the given message.
func NewDirectDependencyError(msg string) *DirectDependencyError {
	return &DirectDependencyError{msg: msg}
}

// AuthenticationError is raised when a remote host rejects credentials or
// none are available.  DiagnosticContext holds pre-rendered guidance.
type AuthenticationError struct {
	msg               string
	DiagnosticContext string
}

func (e *AuthenticationError) Error() string { return e.msg }

// NewAuthenticationError creates an AuthenticationError.
func NewAuthenticationError(msg, diagnosticContext string) *AuthenticationError {
	return &AuthenticationError{msg: msg, DiagnosticContext: diagnosticContext}
}

// FrozenInstallError is raised when `apm install --frozen` cannot proceed
// because the lockfile is missing or structurally out of sync with apm.yml.
type FrozenInstallError struct {
	msg     string
	Reasons []string
}

func (e *FrozenInstallError) Error() string { return e.msg }

// NewFrozenInstallError creates a FrozenInstallError.
func NewFrozenInstallError(msg string, reasons []string) *FrozenInstallError {
	if reasons == nil {
		reasons = []string{}
	}
	return &FrozenInstallError{msg: msg, Reasons: reasons}
}

// PolicyViolationError is raised when org-policy enforcement halts an install.
type PolicyViolationError struct {
	msg          string
	PolicySource string
}

func (e *PolicyViolationError) Error() string { return e.msg }

// NewPolicyViolationError creates a PolicyViolationError.
func NewPolicyViolationError(msg, policySource string) *PolicyViolationError {
	return &PolicyViolationError{msg: msg, PolicySource: policySource}
}
