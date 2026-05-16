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
