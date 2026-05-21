package policygate

import "testing"

func TestPolicyViolationError_MessageAndSource(t *testing.T) {
	err := PolicyViolationError{Message: "blocked by policy", PolicySource: "https://example.com/policy"}
	if err.Error() != "blocked by policy" {
		t.Errorf("Error() mismatch: %s", err.Error())
	}
	if err.PolicySource != "https://example.com/policy" {
		t.Error("PolicySource field mismatch")
	}
}

func TestPolicyViolationError_EmptyMessage(t *testing.T) {
	err := PolicyViolationError{}
	if err.Error() != "" {
		t.Errorf("expected empty error string, got: %s", err.Error())
	}
}

func TestEnforcementResult_ZeroValue(t *testing.T) {
	var r EnforcementResult
	if r.EnforcementActive || r.HasBlocking || r.PolicySource != "" {
		t.Error("EnforcementResult zero value should have all zero fields")
	}
}

func TestEnforcementResult_PolicySourceSet(t *testing.T) {
	r := EnforcementResult{PolicySource: "https://my-policy.io"}
	if r.PolicySource != "https://my-policy.io" {
		t.Error("PolicySource field mismatch")
	}
}

func TestIsDisabledByEnvVar_Disabled(t *testing.T) {
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			return "1"
		}
		return ""
	}
	if !IsDisabledByEnvVar(env) {
		t.Error("should be disabled when APM_POLICY_DISABLE=1")
	}
}

func TestIsDisabledByEnvVar_NotDisabled(t *testing.T) {
	env := func(key string) string { return "" }
	if IsDisabledByEnvVar(env) {
		t.Error("should not be disabled when APM_POLICY_DISABLE is unset")
	}
}

func TestIsDisabledByEnvVar_Zero(t *testing.T) {
	env := func(key string) string {
		if key == "APM_POLICY_DISABLE" {
			return "0"
		}
		return ""
	}
	if IsDisabledByEnvVar(env) {
		t.Error("should not be disabled when APM_POLICY_DISABLE=0")
	}
}

func TestIsDisabledByEnvVar_OtherKeys_Extra2(t *testing.T) {
	env := func(key string) string {
		if key == "OTHER_VAR" {
			return "1"
		}
		return ""
	}
	if IsDisabledByEnvVar(env) {
		t.Error("should not be disabled when other env vars are set")
	}
}

func TestEnforcementResult_ActiveBlocking_Extra2(t *testing.T) {
	r := EnforcementResult{EnforcementActive: true, HasBlocking: true, PolicySource: "src"}
	if !r.EnforcementActive || !r.HasBlocking {
		t.Error("active blocking enforcement result fields not set correctly")
	}
}
