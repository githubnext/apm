// Package policygate implements the policy enforcement gate phase.
// Mirrors src/apm_cli/install/phases/policy_gate.py.
package policygate

// PolicyViolationError signals install blocked by org policy.
type PolicyViolationError struct {
	Message     string
	PolicySource string
}

func (e PolicyViolationError) Error() string {
	return e.Message
}

// EnforcementResult describes the outcome of a policy gate evaluation.
type EnforcementResult struct {
	// EnforcementActive is true when dep checks were run (policy found + enforcement != "off").
	EnforcementActive bool

	// HasBlocking is true when at least one check returned a "block" severity finding.
	HasBlocking bool

	// PolicySource is the URL or identifier of the policy that was fetched.
	PolicySource string
}

// IsDisabledByEnvVar returns true when APM_POLICY_DISABLE=1 is set.
func IsDisabledByEnvVar(env func(string) string) bool {
	return env("APM_POLICY_DISABLE") == "1"
}
