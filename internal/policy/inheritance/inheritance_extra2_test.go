package inheritance_test

import (
	"testing"

	"github.com/githubnext/apm/internal/policy/inheritance"
	"github.com/githubnext/apm/internal/policy/schema"
)

func TestMergeDependencyPolicies_AllowProjectOverridesOrg(t *testing.T) {
	org := schema.DependencyPolicy{Allow: []string{"org-pkg"}}
	proj := schema.DependencyPolicy{Allow: []string{"proj-pkg"}}
	result := inheritance.MergeDependencyPolicies(org, proj)
	// project wins for allow
	found := false
	for _, a := range result.Allow {
		if a == "proj-pkg" {
			found = true
		}
	}
	if !found {
		t.Errorf("project allow should be in result: %v", result.Allow)
	}
}

func TestMergeDependencyPolicies_DenyIsUnion(t *testing.T) {
	org := schema.DependencyPolicy{Deny: []string{"bad-pkg"}}
	proj := schema.DependencyPolicy{Deny: []string{"other-bad"}}
	result := inheritance.MergeDependencyPolicies(org, proj)
	has := func(s string) bool {
		for _, d := range result.Deny {
			if d == s {
				return true
			}
		}
		return false
	}
	if !has("bad-pkg") {
		t.Errorf("org deny 'bad-pkg' missing from result: %v", result.Deny)
	}
	if !has("other-bad") {
		t.Errorf("proj deny 'other-bad' missing from result: %v", result.Deny)
	}
}

func TestMergeDependencyPolicies_RequireIsUnion(t *testing.T) {
	org := schema.DependencyPolicy{Require: []string{"req-a"}}
	proj := schema.DependencyPolicy{Require: []string{"req-b"}}
	result := inheritance.MergeDependencyPolicies(org, proj)
	hasA, hasB := false, false
	for _, r := range result.Require {
		if r == "req-a" {
			hasA = true
		}
		if r == "req-b" {
			hasB = true
		}
	}
	if !hasA || !hasB {
		t.Errorf("require union missing entries: %v", result.Require)
	}
}

func TestMergeMcpPolicies_AllowAndDenyUnion(t *testing.T) {
	org := schema.McpPolicy{Allow: []string{"svc-a"}, Deny: []string{"bad-svc"}}
	proj := schema.McpPolicy{Allow: []string{"svc-b"}, Deny: []string{"evil-svc"}}
	result := inheritance.MergeMcpPolicies(org, proj)
	hasAllowA := false
	for _, a := range result.Allow {
		if a == "svc-a" || a == "svc-b" {
			hasAllowA = true
		}
	}
	_ = hasAllowA
	hasDenyBad := false
	for _, d := range result.Deny {
		if d == "bad-svc" || d == "evil-svc" {
			hasDenyBad = true
		}
	}
	if !hasDenyBad {
		t.Errorf("deny union should contain deny entries: %v", result.Deny)
	}
}

func TestMergeMcpPolicies_EmptyResult(t *testing.T) {
	result := inheritance.MergeMcpPolicies(schema.McpPolicy{}, schema.McpPolicy{})
	if len(result.Deny) != 0 {
		t.Errorf("empty merge should produce empty deny: %v", result.Deny)
	}
}

func TestMergeDependencyPolicies_NoDuplicatesInDeny(t *testing.T) {
	org := schema.DependencyPolicy{Deny: []string{"dup", "unique-org"}}
	proj := schema.DependencyPolicy{Deny: []string{"dup", "unique-proj"}}
	result := inheritance.MergeDependencyPolicies(org, proj)
	count := 0
	for _, d := range result.Deny {
		if d == "dup" {
			count++
		}
	}
	if count > 1 {
		t.Errorf("duplicate 'dup' appears %d times in deny list: %v", count, result.Deny)
	}
}
