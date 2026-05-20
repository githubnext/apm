package policytargetcheck

import "testing"

func TestCheckResult_Fields(t *testing.T) {
	cr := CheckResult{
		Name:    "compilation-target",
		Passed:  true,
		Message: "ok",
		Details: []string{"detail1"},
	}
	if cr.Name != "compilation-target" {
		t.Errorf("expected compilation-target, got %q", cr.Name)
	}
	if !cr.Passed {
		t.Error("expected passed=true")
	}
	if len(cr.Details) != 1 {
		t.Errorf("expected 1 detail, got %d", len(cr.Details))
	}
}

func TestPolicyViolationError_Error(t *testing.T) {
	e := PolicyViolationError{Message: "blocked"}
	if e.Error() != "blocked" {
		t.Errorf("expected blocked, got %q", e.Error())
	}
}

func TestPolicyViolationError_Empty(t *testing.T) {
	e := PolicyViolationError{}
	if e.Error() != "" {
		t.Errorf("expected empty error, got %q", e.Error())
	}
}

func TestShouldRunCheck_CompilationTarget(t *testing.T) {
	if !ShouldRunCheck("compilation-target") {
		t.Error("compilation-target should be runnable")
	}
}

func TestShouldRunCheck_Unknown(t *testing.T) {
	if ShouldRunCheck("unknown-check") {
		t.Error("unknown check should not run")
	}
}

func TestTargetCheckIDs_Length(t *testing.T) {
	if len(TargetCheckIDs) == 0 {
		t.Error("TargetCheckIDs should not be empty")
	}
}

func TestShouldRunCheck_EmptyString(t *testing.T) {
	if ShouldRunCheck("") {
		t.Error("empty string should not be a valid check")
	}
}

func TestCheckResult_PassedTrue(t *testing.T) {
	cr := CheckResult{Passed: true}
	if !cr.Passed {
		t.Error("expected passed=true")
	}
}

func TestCheckResult_NilDetails(t *testing.T) {
	cr := CheckResult{Name: "test"}
	if cr.Details != nil {
		t.Error("zero value Details should be nil")
	}
}
