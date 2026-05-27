// Package policy -- glob-style pattern matcher for policy allow/deny lists.
// Mirrors Python apm_cli.policy.matcher.
package policy

import "path/filepath"

// MatchesAny returns true if the candidate matches any of the given glob patterns.
// Patterns use filepath.Match semantics (shell glob, no path separator in *).
func MatchesAny(candidate string, patterns []string) bool {
	for _, p := range patterns {
		if ok, _ := filepath.Match(p, candidate); ok {
			return true
		}
	}
	return false
}

// MatchesAllowList returns true if candidate is permitted by the allowlist.
// nil allowlist = "no opinion" (allow all). Empty allowlist = "allow nothing".
func MatchesAllowList(candidate string, allow []string) bool {
	if allow == nil {
		return true
	}
	return MatchesAny(candidate, allow)
}

// MatchesDenyList returns true if candidate is blocked by the denylist.
// nil denylist = "no opinion" (deny nothing). Empty = deny nothing.
func MatchesDenyList(candidate string, deny []string) bool {
	if deny == nil || len(deny) == 0 {
		return false
	}
	return MatchesAny(candidate, deny)
}
