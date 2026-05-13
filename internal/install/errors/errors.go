// Package errors provides canonical exception types for the install pipeline.
//
// Centralises typed errors raised by the install machinery so call sites
// can handle a single class hierarchy.
package errors

// DirectDependencyError is raised when one or more direct dependencies fail
// validation or integration.
type DirectDependencyError struct {
	Msg string
}

func (e *DirectDependencyError) Error() string { return e.Msg }

// NewDirectDependencyError creates a DirectDependencyError.
func NewDirectDependencyError(msg string) *DirectDependencyError {
	return &DirectDependencyError{Msg: msg}
}

// AuthenticationError is raised when a remote host rejects credentials or
// none are available.
type AuthenticationError struct {
	Msg               string
	DiagnosticContext string
}

func (e *AuthenticationError) Error() string { return e.Msg }

// NewAuthenticationError creates an AuthenticationError.
func NewAuthenticationError(msg, diagnosticContext string) *AuthenticationError {
	return &AuthenticationError{Msg: msg, DiagnosticContext: diagnosticContext}
}

// FrozenInstallError is raised when apm install --frozen cannot proceed.
// Two trigger conditions:
//   - Lockfile (apm.lock.yaml) is missing entirely.
//   - Lockfile is structurally out of sync with apm.yml.
type FrozenInstallError struct {
	Msg     string
	Reasons []string
}

func (e *FrozenInstallError) Error() string { return e.Msg }

// NewFrozenInstallError creates a FrozenInstallError.
func NewFrozenInstallError(msg string, reasons []string) *FrozenInstallError {
	r := make([]string, len(reasons))
	copy(r, reasons)
	return &FrozenInstallError{Msg: msg, Reasons: r}
}

// PolicyViolationError is raised when org-policy enforcement halts an install.
type PolicyViolationError struct {
	Msg          string
	PolicySource string
}

func (e *PolicyViolationError) Error() string { return e.Msg }

// NewPolicyViolationError creates a PolicyViolationError.
func NewPolicyViolationError(msg, policySource string) *PolicyViolationError {
	return &PolicyViolationError{Msg: msg, PolicySource: policySource}
}

// IsDirect returns true if err is a DirectDependencyError.
func IsDirect(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*DirectDependencyError)
	return ok
}

// IsAuthentication returns true if err is an AuthenticationError.
func IsAuthentication(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*AuthenticationError)
	return ok
}

// IsFrozen returns true if err is a FrozenInstallError.
func IsFrozen(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*FrozenInstallError)
	return ok
}

// IsPolicy returns true if err is a PolicyViolationError.
func IsPolicy(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*PolicyViolationError)
	return ok
}
