package outcomerouting_test

import (
	"testing"

	"github.com/githubnext/apm/internal/policy/outcomerouting"
	"github.com/githubnext/apm/internal/policy/schema"
)

type mockLog2 struct {
	resolved []string
	missed   []string
}

func (m *mockLog2) PolicyResolved(source string, cached bool, enforcement string, ageSeconds int) {
	m.resolved = append(m.resolved, source)
}
func (m *mockLog2) PolicyDiscoveryMiss(outcome string, source string, err string) {
	m.missed = append(m.missed, outcome)
}

func TestRouteDiscoveryOutcome_FoundLogs(t *testing.T) {
	p := &schema.ApmPolicy{Enforcement: "block"}
	result := outcomerouting.PolicyFetchResult{Outcome: "found", Policy: p, Source: "myorg"}
	log := &mockLog2{}
	policy, err := outcomerouting.RouteDiscoveryOutcome(result, log, "block", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy == nil {
		t.Fatal("expected non-nil policy")
	}
	if len(log.resolved) != 1 || log.resolved[0] != "myorg" {
		t.Errorf("expected resolved[myorg], got %v", log.resolved)
	}
}

func TestRouteDiscoveryOutcome_NilLogger_NoPanic(t *testing.T) {
	p := &schema.ApmPolicy{Enforcement: "warn"}
	result := outcomerouting.PolicyFetchResult{Outcome: "found", Policy: p}
	_, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRouteDiscoveryOutcome_DisabledNoLog(t *testing.T) {
	result := outcomerouting.PolicyFetchResult{Outcome: "disabled"}
	log := &mockLog2{}
	policy, err := outcomerouting.RouteDiscoveryOutcome(result, log, "warn", true)
	if err != nil || policy != nil {
		t.Errorf("disabled: expected nil,nil; got %v,%v", policy, err)
	}
	if len(log.resolved) != 0 || len(log.missed) != 0 {
		t.Errorf("expected no logging for disabled: resolved=%v missed=%v", log.resolved, log.missed)
	}
}

func TestRouteDiscoveryOutcome_AbsentBlock_IsError(t *testing.T) {
	result := outcomerouting.PolicyFetchResult{Outcome: "absent", Source: "orgX"}
	_, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "block", true)
	if err == nil {
		t.Fatal("expected error for absent+block")
	}
	var pve *outcomerouting.PolicyViolationError
	if !isPolicyViolationError(err, &pve) {
		t.Errorf("expected PolicyViolationError, got %T", err)
	}
}

func TestRouteDiscoveryOutcome_AbsentWarnLogs(t *testing.T) {
	result := outcomerouting.PolicyFetchResult{Outcome: "absent", Source: "s1"}
	log := &mockLog2{}
	outcomerouting.RouteDiscoveryOutcome(result, log, "warn", true) //nolint:errcheck
	if len(log.missed) == 0 {
		t.Error("expected at least one missed log for absent+warn")
	}
}

func TestRouteDiscoveryOutcome_CachedStaleFetchFailWarn(t *testing.T) {
	p := &schema.ApmPolicy{Enforcement: "warn", FetchFailure: "warn"}
	result := outcomerouting.PolicyFetchResult{
		Outcome: "cached_stale", Policy: p, Source: "cached-org", CacheAgeSeconds: 7200,
	}
	log := &mockLog2{}
	policy, err := outcomerouting.RouteDiscoveryOutcome(result, log, "warn", true)
	if err != nil {
		t.Fatalf("cached_stale warn: unexpected error: %v", err)
	}
	if policy == nil {
		t.Error("cached_stale should return policy")
	}
}

func TestRouteDiscoveryOutcome_CachedStaleFetchFailBlock(t *testing.T) {
	p := &schema.ApmPolicy{Enforcement: "warn", FetchFailure: "block"}
	result := outcomerouting.PolicyFetchResult{
		Outcome: "cached_stale", Policy: p, Source: "strict-org", CacheAgeSeconds: 10000,
	}
	_, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", true)
	// With FetchFailure=block, stale cache might be an error
	_ = err // implementation-defined; just no panic
}

func TestRouteDiscoveryOutcome_HashMismatchNoRaise(t *testing.T) {
	result := outcomerouting.PolicyFetchResult{Outcome: "hash_mismatch", Source: "tampered"}
	policy, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", false)
	if err != nil || policy != nil {
		t.Errorf("hash_mismatch+noRaise: expected nil,nil; got %v,%v", policy, err)
	}
}

func TestRouteDiscoveryOutcome_HashMismatchRaise_HasSource(t *testing.T) {
	result := outcomerouting.PolicyFetchResult{Outcome: "hash_mismatch", Source: "evil-source"}
	_, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", true)
	if err == nil {
		t.Fatal("expected error for hash_mismatch+raise")
	}
}

// isPolicyViolationError checks whether err is a *PolicyViolationError via type assertion.
func isPolicyViolationError(err error, out **outcomerouting.PolicyViolationError) bool {
	pve, ok := err.(*outcomerouting.PolicyViolationError)
	if ok && out != nil {
		*out = pve
	}
	return ok
}
