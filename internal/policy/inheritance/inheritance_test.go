package inheritance_test

import (
"testing"

"github.com/githubnext/apm/internal/policy/inheritance"
"github.com/githubnext/apm/internal/policy/schema"
)

func TestMergeDependencyPolicies_DenyUnion(t *testing.T) {
org := schema.DependencyPolicy{Deny: []string{"bad-pkg", "evil-pkg"}}
proj := schema.DependencyPolicy{Deny: []string{"local-deny"}}
result := inheritance.MergeDependencyPolicies(org, proj)
seen := map[string]bool{}
for _, d := range result.Deny {
seen[d] = true
}
for _, want := range []string{"bad-pkg", "evil-pkg", "local-deny"} {
if !seen[want] {
t.Errorf("deny union missing %q", want)
}
}
}

func TestMergeDependencyPolicies_RequireUnion(t *testing.T) {
org := schema.DependencyPolicy{Require: []string{"org-req"}}
proj := schema.DependencyPolicy{Require: []string{"proj-req"}}
result := inheritance.MergeDependencyPolicies(org, proj)
seen := map[string]bool{}
for _, r := range result.Require {
seen[r] = true
}
if !seen["org-req"] || !seen["proj-req"] {
t.Errorf("require union incorrect: %v", result.Require)
}
}

func TestMergeDependencyPolicies_RequireResolutionEscalation(t *testing.T) {
org := schema.DependencyPolicy{RequireResolution: "block"}
proj := schema.DependencyPolicy{RequireResolution: "project-wins"}
result := inheritance.MergeDependencyPolicies(org, proj)
if result.RequireResolution != "block" {
t.Errorf("expected block, got %q", result.RequireResolution)
}
}

func TestMergeDependencyPolicies_MaxDepthOrgWins(t *testing.T) {
org := schema.DependencyPolicy{MaxDepth: 2}
proj := schema.DependencyPolicy{MaxDepth: 5}
result := inheritance.MergeDependencyPolicies(org, proj)
if result.MaxDepth != 2 {
t.Errorf("expected org MaxDepth=2, got %d", result.MaxDepth)
}
}

func TestMergeDependencyPolicies_DenyNoDuplicates(t *testing.T) {
org := schema.DependencyPolicy{Deny: []string{"shared"}}
proj := schema.DependencyPolicy{Deny: []string{"shared", "extra"}}
result := inheritance.MergeDependencyPolicies(org, proj)
count := 0
for _, d := range result.Deny {
if d == "shared" {
count++
}
}
if count != 1 {
t.Errorf("expected shared once, got %d times in %v", count, result.Deny)
}
}

func TestMergeMcpPolicies_DenyUnion(t *testing.T) {
org := schema.McpPolicy{Deny: []string{"bad-mcp"}}
proj := schema.McpPolicy{Deny: []string{"local-mcp"}}
result := inheritance.MergeMcpPolicies(org, proj)
seen := map[string]bool{}
for _, d := range result.Deny {
seen[d] = true
}
if !seen["bad-mcp"] || !seen["local-mcp"] {
t.Errorf("mcp deny union incorrect: %v", result.Deny)
}
}

func TestMergeMcpPolicies_TrustTransitiveOrgWins(t *testing.T) {
org := schema.McpPolicy{TrustTransitive: true}
proj := schema.McpPolicy{TrustTransitive: false}
result := inheritance.MergeMcpPolicies(org, proj)
if !result.TrustTransitive {
t.Error("org TrustTransitive=true should propagate when project is false")
}
}
