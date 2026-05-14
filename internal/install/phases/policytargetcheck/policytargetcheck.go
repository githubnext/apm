// Package policytargetcheck implements the post-targets target-aware policy check phase.
// Mirrors src/apm_cli/install/phases/policy_target_check.py.
package policytargetcheck

// TargetCheckIDs lists the check names that are target/compilation-related.
// Only these are processed in this phase; all other check IDs already ran in
// the policy_gate phase and must not be double-emitted.
var TargetCheckIDs = map[string]bool{
	"compilation-target": true,
}

// CheckResult mirrors a single policy check result.
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
	Details []string
}

// PolicyViolationError signals a blocking policy enforcement failure.
type PolicyViolationError struct {
	Message string
}

func (e PolicyViolationError) Error() string {
	return e.Message
}

// ShouldRunCheck returns true when a check should be processed in this phase.
func ShouldRunCheck(name string) bool {
	return TargetCheckIDs[name]
}
