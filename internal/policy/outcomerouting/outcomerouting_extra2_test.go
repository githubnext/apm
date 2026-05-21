package outcomerouting

import (
	"testing"
)

func TestPolicyFetchResult_ZeroValue(t *testing.T) {
	var r PolicyFetchResult
	if r.Outcome != "" {
		t.Errorf("expected empty Outcome, got %q", r.Outcome)
	}
	if r.Cached {
		t.Error("expected Cached false")
	}
	if r.Policy != nil {
		t.Error("expected nil Policy")
	}
}

func TestPolicyFetchResult_FieldAssignment(t *testing.T) {
	r := PolicyFetchResult{
		Outcome:         "found",
		Source:          "https://example.com",
		Cached:          true,
		Error:           "",
		CacheAgeSeconds: 120,
	}
	if r.Outcome != "found" {
		t.Errorf("unexpected Outcome %q", r.Outcome)
	}
	if !r.Cached {
		t.Error("expected Cached true")
	}
	if r.CacheAgeSeconds != 120 {
		t.Errorf("expected CacheAgeSeconds 120, got %d", r.CacheAgeSeconds)
	}
}

func TestPolicyViolationError_Error(t *testing.T) {
	e := &PolicyViolationError{Message: "blocked", PolicySource: "org"}
	if e.Error() != "blocked" {
		t.Errorf("unexpected Error() %q", e.Error())
	}
	if e.PolicySource != "org" {
		t.Errorf("unexpected PolicySource %q", e.PolicySource)
	}
}

func TestPolicyViolationError_EmptyMessage(t *testing.T) {
	e := &PolicyViolationError{}
	if e.Error() != "" {
		t.Errorf("expected empty error string, got %q", e.Error())
	}
}

func TestRouteDiscoveryOutcome_Found_NilLogger(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "found", Source: "src", Policy: nil}
	policy, err := RouteDiscoveryOutcome(fr, nil, "warn", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if policy != nil {
		t.Error("expected nil policy for nil fetch result")
	}
}

func TestRouteDiscoveryOutcome_Disabled_ReturnsNil(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "disabled"}
	policy, err := RouteDiscoveryOutcome(fr, nil, "warn", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if policy != nil {
		t.Error("expected nil policy for disabled")
	}
}

func TestRouteDiscoveryOutcome_Absent_BlockMode(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "absent", Source: "org-url"}
	_, err := RouteDiscoveryOutcome(fr, nil, "block", true)
	if err == nil {
		t.Error("expected error for absent + block")
	}
}

func TestRouteDiscoveryOutcome_NoGitRemote_WarnMode(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "no_git_remote", Source: ""}
	policy, err := RouteDiscoveryOutcome(fr, nil, "warn", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if policy != nil {
		t.Error("expected nil policy")
	}
}

func TestRouteDiscoveryOutcome_UnknownOutcome_NoError(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "totally_unknown_outcome"}
	policy, err := RouteDiscoveryOutcome(fr, nil, "warn", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if policy != nil {
		t.Error("expected nil policy for unknown outcome")
	}
}

func TestRouteDiscoveryOutcome_Absent_WarnMode_NoError(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "absent", Source: "src"}
	_, err := RouteDiscoveryOutcome(fr, nil, "warn", true)
	if err != nil {
		t.Errorf("expected no error for absent+warn, got: %v", err)
	}
}

func TestRouteDiscoveryOutcome_GarbageResponse_Block(t *testing.T) {
	fr := PolicyFetchResult{Outcome: "garbage_response", Source: "s"}
	_, err := RouteDiscoveryOutcome(fr, nil, "block", true)
	if err == nil {
		t.Error("expected error for garbage_response + block")
	}
}
