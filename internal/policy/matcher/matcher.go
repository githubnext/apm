// Package matcher implements pattern matching for policy allow/deny lists.
package matcher

import (
"regexp"
"strings"
"sync"
)

var (
patternCacheMu sync.Mutex
patternCache   = map[string]*regexp.Regexp{}
)

func compilePattern(pattern string) *regexp.Regexp {
patternCacheMu.Lock()
defer patternCacheMu.Unlock()
if re, ok := patternCache[pattern]; ok {
return re
}
parts := strings.Split(pattern, "**")
var sb strings.Builder
for i, part := range parts {
if i > 0 {
sb.WriteString(".*")
}
subParts := strings.Split(part, "*")
for j, sub := range subParts {
if j > 0 {
sb.WriteString("[^/]*")
}
sb.WriteString(regexp.QuoteMeta(sub))
}
}
re := regexp.MustCompile("^" + sb.String() + "$")
patternCache[pattern] = re
return re
}

// MatchesPattern checks if a canonical dependency ref matches a policy pattern.
func MatchesPattern(canonicalRef, pattern string) bool {
if pattern == "" || canonicalRef == "" {
return false
}
if canonicalRef == pattern {
return true
}
return compilePattern(pattern).MatchString(canonicalRef)
}

// CheckAllowDeny implements shared allow/deny logic.
// Returns (allowed bool, reason string).
func CheckAllowDeny(ref string, allow []string, deny []string) (bool, string) {
for _, p := range deny {
if MatchesPattern(ref, p) {
return false, "denied by pattern: " + p
}
}
if allow == nil {
return true, ""
}
if len(allow) == 0 {
return false, "allow list is empty: all refs blocked"
}
for _, p := range allow {
if MatchesPattern(ref, p) {
return true, ""
}
}
return false, "not in allowed sources"
}
