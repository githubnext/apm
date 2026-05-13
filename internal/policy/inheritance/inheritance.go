// Package inheritance implements policy inheritance and merging logic.
package inheritance

import (
"github.com/githubnext/apm/internal/policy/schema"
)

// escalationOrder defines restriction severity for require_resolution.
var escalationOrder = map[string]int{
"project-wins": 0,
"policy-wins":  1,
"block":        2,
}

func stricter(a, b string) string {
ai, aok := escalationOrder[a]
bi, bok := escalationOrder[b]
if !aok {
ai = 0
}
if !bok {
bi = 0
}
if ai >= bi {
return a
}
return b
}

// MergeDependencyPolicies merges base (org) policy with project policy.
// Project values take precedence for allow; org values accumulate deny/require.
func MergeDependencyPolicies(org, project schema.DependencyPolicy) schema.DependencyPolicy {
result := project

// Merge deny lists (union)
deny := append([]string{}, org.Deny...)
deny = append(deny, project.Deny...)
result.Deny = unique(deny)

// Merge require lists (union)
require := append([]string{}, org.Require...)
require = append(require, project.Require...)
result.Require = unique(require)

// Escalate resolution
result.RequireResolution = stricter(org.RequireResolution, project.RequireResolution)

// MaxDepth: use the more restrictive (lower) value when org sets one.
if org.MaxDepth > 0 && (result.MaxDepth == 0 || org.MaxDepth < result.MaxDepth) {
result.MaxDepth = org.MaxDepth
}

return result
}

// MergeMcpPolicies merges base (org) McpPolicy with project McpPolicy.
func MergeMcpPolicies(org, project schema.McpPolicy) schema.McpPolicy {
result := project
deny := append([]string{}, org.Deny...)
deny = append(deny, project.Deny...)
result.Deny = unique(deny)
if org.TrustTransitive && !project.TrustTransitive {
result.TrustTransitive = org.TrustTransitive
}
return result
}

func unique(strs []string) []string {
seen := map[string]struct{}{}
out := []string{}
for _, s := range strs {
if _, ok := seen[s]; !ok {
seen[s] = struct{}{}
out = append(out, s)
}
}
return out
}
