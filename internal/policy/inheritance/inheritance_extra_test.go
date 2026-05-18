package inheritance_test

import (
	"testing"

	"github.com/githubnext/apm/internal/policy/inheritance"
	"github.com/githubnext/apm/internal/policy/schema"
)

func TestMergeDependencyPolicies_EmptyBoth(t *testing.T) {
	result := inheritance.MergeDependencyPolicies(schema.DependencyPolicy{}, schema.DependencyPolicy{})
	if len(result.Deny) != 0 {
		t.Errorf("expected empty deny, got %v", result.Deny)
	}
	if len(result.Require) != 0 {
		t.Errorf("expected empty require, got %v", result.Require)
	}
}

func TestMergeDependencyPolicies_ProjectWinsAllow(t *testing.T) {
	org := schema.DependencyPolicy{Allow: []string{"org-allowed"}}
	proj := schema.DependencyPolicy{Allow: []string{"proj-allowed"}}
	result := inheritance.MergeDependencyPolicies(org, proj)
	// Project values take precedence for allow
	for _, a := range result.Allow {
		if a == "proj-allowed" {
			return
		}
	}
	t.Errorf("expected 'proj-allowed' in allow list: %v", result.Allow)
}

func TestMergeDependencyPolicies_MaxDepthZeroOrgIgnored(t *testing.T) {
	org := schema.DependencyPolicy{MaxDepth: 0}
	proj := schema.DependencyPolicy{MaxDepth: 3}
	result := inheritance.MergeDependencyPolicies(org, proj)
	if result.MaxDepth != 3 {
		t.Errorf("org MaxDepth=0 should not override project MaxDepth=3, got %d", result.MaxDepth)
	}
}

func TestMergeDependencyPolicies_MaxDepthOrgStricter(t *testing.T) {
	org := schema.DependencyPolicy{MaxDepth: 1}
	proj := schema.DependencyPolicy{MaxDepth: 10}
	result := inheritance.MergeDependencyPolicies(org, proj)
	if result.MaxDepth != 1 {
		t.Errorf("expected MaxDepth=1 (org stricter), got %d", result.MaxDepth)
	}
}

func TestMergeDependencyPolicies_MaxDepthProjectWhenOrgNotSet(t *testing.T) {
	proj := schema.DependencyPolicy{MaxDepth: 5}
	result := inheritance.MergeDependencyPolicies(schema.DependencyPolicy{}, proj)
	if result.MaxDepth != 5 {
		t.Errorf("expected MaxDepth=5 (project), got %d", result.MaxDepth)
	}
}

func TestMergeDependencyPolicies_RequireResolutionDefaults(t *testing.T) {
	// Empty strings: escalation defaults to 0 for both; first arg returned
	result := inheritance.MergeDependencyPolicies(schema.DependencyPolicy{}, schema.DependencyPolicy{})
	// Both "" -> stricter("", "") -> "" (both have escalation 0, ai >= bi returns a)
	_ = result.RequireResolution // just verify no panic
}

func TestMergeDependencyPolicies_ProjectWinsResolution(t *testing.T) {
	org := schema.DependencyPolicy{RequireResolution: "project-wins"}
	proj := schema.DependencyPolicy{RequireResolution: "block"}
	result := inheritance.MergeDependencyPolicies(org, proj)
	if result.RequireResolution != "block" {
		t.Errorf("expected 'block', got %q", result.RequireResolution)
	}
}

func TestMergeMcpPolicies_EmptyBoth(t *testing.T) {
	result := inheritance.MergeMcpPolicies(schema.McpPolicy{}, schema.McpPolicy{})
	if len(result.Deny) != 0 {
		t.Errorf("expected empty deny, got %v", result.Deny)
	}
	if result.TrustTransitive {
		t.Error("expected TrustTransitive=false")
	}
}

func TestMergeMcpPolicies_BothTrustTransitive(t *testing.T) {
	org := schema.McpPolicy{TrustTransitive: true}
	proj := schema.McpPolicy{TrustTransitive: true}
	result := inheritance.MergeMcpPolicies(org, proj)
	if !result.TrustTransitive {
		t.Error("expected TrustTransitive=true when both are true")
	}
}

func TestMergeMcpPolicies_NoDuplicatesInDeny(t *testing.T) {
	org := schema.McpPolicy{Deny: []string{"a", "b"}}
	proj := schema.McpPolicy{Deny: []string{"b", "c"}}
	result := inheritance.MergeMcpPolicies(org, proj)
	count := 0
	for _, d := range result.Deny {
		if d == "b" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 'b' once in deny, got %d times: %v", count, result.Deny)
	}
}

func TestMergeMcpPolicies_OrgDenyOnlyProjectEmpty(t *testing.T) {
	org := schema.McpPolicy{Deny: []string{"org-only"}}
	result := inheritance.MergeMcpPolicies(org, schema.McpPolicy{})
	found := false
	for _, d := range result.Deny {
		if d == "org-only" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'org-only' in deny: %v", result.Deny)
	}
}
