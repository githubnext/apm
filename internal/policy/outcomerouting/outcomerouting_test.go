package outcomerouting_test

import (
"testing"

"github.com/githubnext/apm/internal/policy/outcomerouting"
"github.com/githubnext/apm/internal/policy/schema"
)

type mockLogger struct {
resolved []string
missed   []string
}

func (m *mockLogger) PolicyResolved(source string, cached bool, enforcement string, ageSeconds int) {
m.resolved = append(m.resolved, source)
}
func (m *mockLogger) PolicyDiscoveryMiss(outcome string, source string, err string) {
m.missed = append(m.missed, outcome)
}

func TestRouteDiscoveryOutcome_Disabled(t *testing.T) {
result := outcomerouting.PolicyFetchResult{Outcome: "disabled"}
policy, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", true)
if err != nil || policy != nil {
t.Errorf("disabled: expected nil,nil; got %v,%v", policy, err)
}
}

func TestRouteDiscoveryOutcome_Found(t *testing.T) {
p := &schema.ApmPolicy{Enforcement: "warn"}
result := outcomerouting.PolicyFetchResult{Outcome: "found", Policy: p, Source: "org"}
log := &mockLogger{}
policy, err := outcomerouting.RouteDiscoveryOutcome(result, log, "warn", true)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if policy != p {
t.Error("expected the policy to be returned")
}
if len(log.resolved) != 1 {
t.Errorf("expected 1 resolved log, got %d", len(log.resolved))
}
}

func TestRouteDiscoveryOutcome_HashMismatch_Blocks(t *testing.T) {
result := outcomerouting.PolicyFetchResult{Outcome: "hash_mismatch", Source: "org"}
_, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", true)
if err == nil {
t.Error("expected PolicyViolationError for hash_mismatch with raiseBlockingErrors=true")
}
}

func TestRouteDiscoveryOutcome_HashMismatch_NoRaise(t *testing.T) {
result := outcomerouting.PolicyFetchResult{Outcome: "hash_mismatch", Source: "org"}
policy, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", false)
if err != nil || policy != nil {
t.Errorf("expected nil,nil; got %v,%v", policy, err)
}
}

func TestRouteDiscoveryOutcome_Absent_Warn(t *testing.T) {
result := outcomerouting.PolicyFetchResult{Outcome: "absent", Source: "org"}
log := &mockLogger{}
policy, err := outcomerouting.RouteDiscoveryOutcome(result, log, "warn", true)
if err != nil || policy != nil {
t.Errorf("absent+warn: expected nil,nil; got %v,%v", policy, err)
}
if len(log.missed) != 1 || log.missed[0] != "absent" {
t.Errorf("expected absent in missed, got %v", log.missed)
}
}

func TestRouteDiscoveryOutcome_Absent_Block(t *testing.T) {
result := outcomerouting.PolicyFetchResult{Outcome: "absent", Source: "org"}
_, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "block", true)
if err == nil {
t.Error("expected PolicyViolationError for absent+block")
}
}

func TestRouteDiscoveryOutcome_CachedStale_ReturnsPolicy(t *testing.T) {
p := &schema.ApmPolicy{Enforcement: "warn", FetchFailure: "warn"}
result := outcomerouting.PolicyFetchResult{
Outcome: "cached_stale", Policy: p, Source: "org", CacheAgeSeconds: 3600,
}
log := &mockLogger{}
policy, err := outcomerouting.RouteDiscoveryOutcome(result, log, "warn", true)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if policy != p {
t.Error("cached_stale should return the cached policy")
}
}

func TestRouteDiscoveryOutcome_Unknown_ReturnsNil(t *testing.T) {
result := outcomerouting.PolicyFetchResult{Outcome: "unknown_outcome"}
policy, err := outcomerouting.RouteDiscoveryOutcome(result, nil, "warn", true)
if err != nil || policy != nil {
t.Errorf("unknown outcome: expected nil,nil; got %v,%v", policy, err)
}
}
