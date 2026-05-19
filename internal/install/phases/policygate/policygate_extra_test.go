package policygate_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/policygate"
)

func TestIsDisabledByEnvVar_TrueNumeric(t *testing.T) {
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			return "1"
		}
		return ""
	}
	if !policygate.IsDisabledByEnvVar(env) {
		t.Error("expected true for APM_POLICY_DISABLE=1")
	}
}

func TestIsDisabledByEnvVar_OtherEnvKeys(t *testing.T) {
	// Only APM_POLICY_DISABLE matters; other keys are irrelevant
	env := func(key string) string {
		if key == "OTHER_VAR" {
			return "1"
		}
		return ""
	}
	if policygate.IsDisabledByEnvVar(env) {
		t.Error("should not be disabled when APM_POLICY_DISABLE is unset")
	}
}

func TestPolicyViolationError_AsError(t *testing.T) {
	var err error = policygate.PolicyViolationError{Message: "blocked by policy"}
	if err.Error() != "blocked by policy" {
		t.Errorf("unexpected error: %q", err.Error())
	}
}

func TestPolicyViolationError_PolicySourceField(t *testing.T) {
	err := policygate.PolicyViolationError{
		Message:      "install blocked",
		PolicySource: "https://org.example.com/policy.yaml",
	}
	if err.PolicySource == "" {
		t.Error("PolicySource should not be empty")
	}
	if err.Message == "" {
		t.Error("Message should not be empty")
	}
}

func TestEnforcementResult_InactiveNonBlocking(t *testing.T) {
	r := policygate.EnforcementResult{
		EnforcementActive: false,
		HasBlocking:       false,
		PolicySource:      "",
	}
	if r.EnforcementActive {
		t.Error("EnforcementActive should be false")
	}
	if r.HasBlocking {
		t.Error("HasBlocking should be false")
	}
}

func TestEnforcementResult_ActiveBlocking(t *testing.T) {
	r := policygate.EnforcementResult{
		EnforcementActive: true,
		HasBlocking:       true,
		PolicySource:      "https://example.com/org-policy.yaml",
	}
	if !r.EnforcementActive {
		t.Error("expected EnforcementActive=true")
	}
	if !r.HasBlocking {
		t.Error("expected HasBlocking=true")
	}
}

func TestEnforcementResult_ActiveNonBlocking(t *testing.T) {
	r := policygate.EnforcementResult{
		EnforcementActive: true,
		HasBlocking:       false,
		PolicySource:      "https://example.com/lenient-policy.yaml",
	}
	if !r.EnforcementActive {
		t.Error("expected EnforcementActive=true")
	}
	if r.HasBlocking {
		t.Error("expected HasBlocking=false for lenient policy")
	}
}

func TestIsDisabledByEnvVar_WhitespaceNotDisabling(t *testing.T) {
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			return " 1 "
		}
		return ""
	}
	// " 1 " != "1", so should NOT be disabled
	if policygate.IsDisabledByEnvVar(env) {
		t.Error("whitespace-padded value should not trigger disable")
	}
}

func TestPolicyViolationError_Implements_Error_Interface(t *testing.T) {
	errs := []error{
		policygate.PolicyViolationError{Message: "first"},
		policygate.PolicyViolationError{Message: "second", PolicySource: "src"},
		policygate.PolicyViolationError{},
	}
	for _, e := range errs {
		_ = e.Error() // must not panic
	}
}

func TestEnforcementResult_PolicySourceEmpty(t *testing.T) {
	r := policygate.EnforcementResult{
		EnforcementActive: true,
		HasBlocking:       false,
		PolicySource:      "",
	}
	if r.PolicySource != "" {
		t.Errorf("expected empty PolicySource, got %q", r.PolicySource)
	}
}

func TestIsDisabledByEnvVar_ZeroValue(t *testing.T) {
	env := func(key string) string { return "0" }
	if policygate.IsDisabledByEnvVar(env) {
		t.Error("APM_POLICY_DISABLE=0 should not disable")
	}
}

func TestIsDisabledByEnvVar_CallsCorrectKey(t *testing.T) {
	called := false
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			called = true
			return "1"
		}
		return ""
	}
	policygate.IsDisabledByEnvVar(env)
	if !called {
		t.Error("env should be called with APM_POLICY_DISABLE key")
	}
}
