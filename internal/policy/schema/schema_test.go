package schema

import (
	"testing"
)

func TestDefaultDependencyPolicy(t *testing.T) {
	p := DefaultDependencyPolicy()
	if p.RequireResolution != "project-wins" {
		t.Errorf("want 'project-wins', got %q", p.RequireResolution)
	}
	if p.MaxDepth != 50 {
		t.Errorf("want MaxDepth=50, got %d", p.MaxDepth)
	}
}

func TestApmPolicyZeroValue(t *testing.T) {
	var p ApmPolicy
	if p.Enforcement != "" {
		t.Error("zero value Enforcement should be empty")
	}
	if p.Deps.MaxDepth != 0 {
		t.Error("zero value MaxDepth should be 0")
	}
}

func TestPolicyCacheZeroValue(t *testing.T) {
	var pc PolicyCache
	if pc.TTL != 0 {
		t.Error("zero value TTL should be 0")
	}
}

func TestMcpPolicyFields(t *testing.T) {
	p := McpPolicy{
		Allow:           []string{"stdio"},
		Deny:            []string{"sse"},
		SelfDefined:     "warn",
		TrustTransitive: true,
		Transport:       McpTransportPolicy{Allow: []string{"stdio", "sse"}},
	}
	if len(p.Allow) != 1 || p.Allow[0] != "stdio" {
		t.Error("Allow not set correctly")
	}
	if !p.TrustTransitive {
		t.Error("TrustTransitive should be true")
	}
	if len(p.Transport.Allow) != 2 {
		t.Errorf("want 2 transport allows, got %d", len(p.Transport.Allow))
	}
}

func TestCompilationPolicy(t *testing.T) {
	p := CompilationPolicy{
		Targets:  CompilationTargetPolicy{Allow: []string{"all"}, Enforce: "block"},
		Strategy: CompilationStrategyPolicy{Enforce: "distributed"},
	}
	if p.Targets.Enforce != "block" {
		t.Errorf("target enforce: want 'block', got %q", p.Targets.Enforce)
	}
	if p.Strategy.Enforce != "distributed" {
		t.Errorf("strategy enforce: want 'distributed', got %q", p.Strategy.Enforce)
	}
}
