package policytargetcheck_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/policytargetcheck"
)

func TestTargetCheckIDs_ContainsCompilationTarget(t *testing.T) {
	if !policytargetcheck.TargetCheckIDs["compilation-target"] {
		t.Error("TargetCheckIDs must contain compilation-target")
	}
}

func TestTargetCheckIDs_DoesNotContainPolicyGate(t *testing.T) {
	if policytargetcheck.TargetCheckIDs["policy-gate"] {
		t.Error("TargetCheckIDs must not contain policy-gate")
	}
}

func TestTargetCheckIDs_DoesNotContainEmpty(t *testing.T) {
	if policytargetcheck.TargetCheckIDs[""] {
		t.Error("TargetCheckIDs must not contain empty string")
	}
}

func TestShouldRunCheck_KnownIDs(t *testing.T) {
	for id := range policytargetcheck.TargetCheckIDs {
		if !policytargetcheck.ShouldRunCheck(id) {
			t.Errorf("ShouldRunCheck(%q) = false, want true", id)
		}
	}
}

func TestShouldRunCheck_MultipleUnknown(t *testing.T) {
	ids := []string{"random", "build-check", "lint", "format", "test-run"}
	for _, id := range ids {
		if policytargetcheck.ShouldRunCheck(id) {
			t.Errorf("ShouldRunCheck(%q) = true, want false", id)
		}
	}
}

func TestPolicyViolationError_MultilineMessage(t *testing.T) {
	msg := "line1\nline2\nline3"
	err := policytargetcheck.PolicyViolationError{Message: msg}
	if err.Error() != msg {
		t.Errorf("Error() = %q, want %q", err.Error(), msg)
	}
}

func TestPolicyViolationError_SpecialChars(t *testing.T) {
	msg := "policy blocked: path/to/file (rule=no-secrets)"
	err := policytargetcheck.PolicyViolationError{Message: msg}
	if err.Error() != msg {
		t.Errorf("Error() = %q, want %q", err.Error(), msg)
	}
}

func TestCheckResult_ZeroValue(t *testing.T) {
	var cr policytargetcheck.CheckResult
	if cr.Name != "" {
		t.Errorf("zero value Name should be empty")
	}
	if cr.Passed {
		t.Error("zero value Passed should be false")
	}
	if cr.Message != "" {
		t.Error("zero value Message should be empty")
	}
	if len(cr.Details) != 0 {
		t.Error("zero value Details should be nil/empty")
	}
}

func TestCheckResult_MultipleDetails(t *testing.T) {
	cr := policytargetcheck.CheckResult{
		Name:    "compilation-target",
		Passed:  false,
		Message: "blocked",
		Details: []string{"d1", "d2", "d3", "d4"},
	}
	if len(cr.Details) != 4 {
		t.Errorf("expected 4 details, got %d", len(cr.Details))
	}
	if cr.Details[0] != "d1" {
		t.Errorf("Details[0] = %q, want d1", cr.Details[0])
	}
}

func TestCheckResult_EmptyDetails(t *testing.T) {
	cr := policytargetcheck.CheckResult{
		Name:   "compilation-target",
		Passed: true,
	}
	if len(cr.Details) != 0 {
		t.Errorf("expected empty details, got %v", cr.Details)
	}
}

func TestShouldRunCheck_TabAndSpace(t *testing.T) {
	// Whitespace variants should not match
	if policytargetcheck.ShouldRunCheck(" compilation-target") {
		t.Error("leading-space variant should not match")
	}
	if policytargetcheck.ShouldRunCheck("compilation-target ") {
		t.Error("trailing-space variant should not match")
	}
}

func TestShouldRunCheck_PrefixMatch(t *testing.T) {
	// Prefix of a known ID should not match
	if policytargetcheck.ShouldRunCheck("compilation") {
		t.Error("prefix should not match")
	}
}

func TestTargetCheckIDs_IsMapType(t *testing.T) {
	// Ensure we can range over the map without panicking
	count := 0
	for k, v := range policytargetcheck.TargetCheckIDs {
		if k == "" {
			t.Error("empty key in TargetCheckIDs")
		}
		if !v {
			t.Errorf("TargetCheckIDs[%q] = false, all values should be true", k)
		}
		count++
	}
	if count == 0 {
		t.Error("TargetCheckIDs should have at least one entry")
	}
}

func TestCheckResult_PassedFalseWithDetails(t *testing.T) {
	cr := policytargetcheck.CheckResult{
		Name:    "compilation-target",
		Passed:  false,
		Message: "compilation blocked by policy",
		Details: []string{"file: main.go", "rule: no-hardcoded-secrets"},
	}
	if cr.Passed {
		t.Error("Passed should be false")
	}
	if len(cr.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(cr.Details))
	}
}
