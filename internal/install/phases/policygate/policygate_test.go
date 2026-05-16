package policygate_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/policygate"
)

func TestIsDisabledByEnvVar_Disabled(t *testing.T) {
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			return "1"
		}
		return ""
	}
	if !policygate.IsDisabledByEnvVar(env) {
		t.Fatal("expected true when APM_POLICY_DISABLE=1")
	}
}

func TestIsDisabledByEnvVar_Enabled(t *testing.T) {
	env := func(key string) string { return "" }
	if policygate.IsDisabledByEnvVar(env) {
		t.Fatal("expected false when APM_POLICY_DISABLE is not set")
	}
}

func TestIsDisabledByEnvVar_OtherValue(t *testing.T) {
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			return "0"
		}
		return ""
	}
	if policygate.IsDisabledByEnvVar(env) {
		t.Fatal("expected false when APM_POLICY_DISABLE=0")
	}
}

func TestPolicyViolationError_Error(t *testing.T) {
	err := policygate.PolicyViolationError{
		Message:      "policy blocked",
		PolicySource: "https://example.com/policy.yaml",
	}
	if err.Error() != "policy blocked" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestEnforcementResult_Fields(t *testing.T) {
	r := policygate.EnforcementResult{
		EnforcementActive: true,
		HasBlocking:       true,
		PolicySource:      "https://example.com/policy",
	}
	if !r.EnforcementActive {
		t.Fatal("EnforcementActive should be true")
	}
	if !r.HasBlocking {
		t.Fatal("HasBlocking should be true")
	}
	if r.PolicySource == "" {
		t.Fatal("PolicySource should not be empty")
	}
}
