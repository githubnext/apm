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

func TestIsDisabledByEnvVar_EmptyKey(t *testing.T) {
env := func(key string) string { return "" }
if policygate.IsDisabledByEnvVar(env) {
t.Fatal("expected false when env returns empty string for all keys")
}
}

func TestIsDisabledByEnvVar_TrueValue(t *testing.T) {
// IsDisabledByEnvVar only checks for "1"; other truthy values are not supported
env := func(key string) string {
if key == "APM_POLICY_DISABLE" {
return "true"
}
return ""
}
// "true" is not "1", so this should return false
if policygate.IsDisabledByEnvVar(env) {
t.Fatal("expected false when APM_POLICY_DISABLE=true (only '1' is accepted)")
}
}

func TestPolicyViolationError_EmptyMessage(t *testing.T) {
err := policygate.PolicyViolationError{}
if err.Error() != "" {
t.Errorf("empty message should give empty error string, got %q", err.Error())
}
}

func TestPolicyViolationError_WithSourceOnly(t *testing.T) {
err := policygate.PolicyViolationError{PolicySource: "https://example.com/pol.yaml"}
msg := err.Error()
// Message field is empty; Error() returns ""
if msg != "" {
t.Errorf("unexpected message: %q", msg)
}
if err.PolicySource == "" {
t.Error("PolicySource should be set")
}
}

func TestEnforcementResult_ZeroValue(t *testing.T) {
var r policygate.EnforcementResult
if r.EnforcementActive {
t.Error("zero value EnforcementActive should be false")
}
if r.HasBlocking {
t.Error("zero value HasBlocking should be false")
}
}

func TestEnforcementResult_AllFields(t *testing.T) {
r := policygate.EnforcementResult{
EnforcementActive: false,
HasBlocking:       false,
PolicySource:      "https://example.com/policy",
}
if r.PolicySource == "" {
t.Error("PolicySource should not be empty")
}
}
