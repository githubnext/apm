package policytargetcheck_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/policytargetcheck"
)

func TestTargetCheckIDs(t *testing.T) {
	if !policytargetcheck.TargetCheckIDs["compilation-target"] {
		t.Error("expected compilation-target to be in TargetCheckIDs")
	}
	if policytargetcheck.TargetCheckIDs["other-check"] {
		t.Error("expected other-check to not be in TargetCheckIDs")
	}
}

func TestShouldRunCheck(t *testing.T) {
	tests := []struct {
		name     string
		checkID  string
		expected bool
	}{
		{"compilation-target is included", "compilation-target", true},
		{"policy-gate is excluded", "policy-gate", false},
		{"empty string is excluded", "", false},
		{"unknown check is excluded", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policytargetcheck.ShouldRunCheck(tt.checkID)
			if got != tt.expected {
				t.Errorf("ShouldRunCheck(%q) = %v, want %v", tt.checkID, got, tt.expected)
			}
		})
	}
}

func TestPolicyViolationError(t *testing.T) {
	msg := "blocking policy violation"
	err := policytargetcheck.PolicyViolationError{Message: msg}
	if err.Error() != msg {
		t.Errorf("Error() = %q, want %q", err.Error(), msg)
	}

	// CheckResult struct
	cr := policytargetcheck.CheckResult{
		Name:    "compilation-target",
		Passed:  false,
		Message: "blocked",
		Details: []string{"detail1"},
	}
	if cr.Name != "compilation-target" {
		t.Errorf("CheckResult.Name = %q, want %q", cr.Name, "compilation-target")
	}
	if cr.Passed {
		t.Error("expected Passed to be false")
	}
}

func TestTargetCheckIDs_MapImmutability(t *testing.T) {
	// Verify map exists and contains expected keys
	ids := policytargetcheck.TargetCheckIDs
	if ids == nil {
		t.Fatal("TargetCheckIDs should not be nil")
	}
	if len(ids) == 0 {
		t.Fatal("TargetCheckIDs should not be empty")
	}
}

func TestShouldRunCheck_CaseSensitive(t *testing.T) {
	// Case sensitivity: "Compilation-Target" (capital) should not match
	got := policytargetcheck.ShouldRunCheck("Compilation-Target")
	if got {
		t.Error("ShouldRunCheck should be case-sensitive")
	}
}

func TestCheckResult_PassedTrue(t *testing.T) {
	cr := policytargetcheck.CheckResult{
		Name:    "compilation-target",
		Passed:  true,
		Message: "all good",
	}
	if !cr.Passed {
		t.Error("expected Passed to be true")
	}
	if cr.Message != "all good" {
		t.Errorf("Message = %q, want 'all good'", cr.Message)
	}
}

func TestCheckResult_Details(t *testing.T) {
	cr := policytargetcheck.CheckResult{
		Name:    "compilation-target",
		Passed:  false,
		Details: []string{"reason1", "reason2"},
	}
	if len(cr.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(cr.Details))
	}
}

func TestPolicyViolationError_EmptyMessage(t *testing.T) {
	err := policytargetcheck.PolicyViolationError{Message: ""}
	if err.Error() != "" {
		t.Errorf("Error() = %q, want empty string", err.Error())
	}
}
