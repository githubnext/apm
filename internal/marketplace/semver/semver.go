// Package semver provides dependency-free semver parsing and range matching.
package semver

import (
"fmt"
"regexp"
"strconv"
"strings"
)

var semverRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// SemVer is a parsed semantic version.
type SemVer struct {
Major      int
Minor      int
Patch      int
Prerelease string
BuildMeta  string
}

// Parse parses a semver string. Returns error if invalid.
func Parse(s string) (SemVer, error) {
m := semverRe.FindStringSubmatch(strings.TrimSpace(s))
if m == nil {
return SemVer{}, fmt.Errorf("invalid semver: %q", s)
}
major, _ := strconv.Atoi(m[1])
minor, _ := strconv.Atoi(m[2])
patch, _ := strconv.Atoi(m[3])
return SemVer{Major: major, Minor: minor, Patch: patch, Prerelease: m[4], BuildMeta: m[5]}, nil
}

// cmpTuple returns comparable representation (no prerelease = higher precedence).
func (v SemVer) cmpTuple() []int {
if v.Prerelease == "" {
return []int{v.Major, v.Minor, v.Patch, 1}
}
return []int{v.Major, v.Minor, v.Patch, 0}
}

// Compare returns -1, 0, or 1.
func (v SemVer) Compare(other SemVer) int {
a, b := v.cmpTuple(), other.cmpTuple()
for i := 0; i < len(a) && i < len(b); i++ {
if a[i] < b[i] {
return -1
}
if a[i] > b[i] {
return 1
}
}
if v.Prerelease != "" && other.Prerelease != "" {
if v.Prerelease < other.Prerelease {
return -1
}
if v.Prerelease > other.Prerelease {
return 1
}
}
return 0
}

// SatisfiesRange checks if v satisfies the given range string.
// Supports: exact, ^, ~, >=, >, <=, <, 1.2.x/*, AND (space-separated).
func SatisfiesRange(v SemVer, rangeStr string) bool {
parts := strings.Fields(rangeStr)
for _, part := range parts {
if !satisfiesSingle(v, part) {
return false
}
}
return true
}

func satisfiesSingle(v SemVer, r string) bool {
r = strings.TrimSpace(r)
if r == "" || r == "*" {
return true
}
// Wildcard: 1.2.x or 1.2.*
if strings.ContainsAny(r, "x*") && !strings.HasPrefix(r, "^") && !strings.HasPrefix(r, "~") {
r2 := strings.ReplaceAll(strings.ReplaceAll(r, ".x", ".0"), ".*", ".0")
base, err := Parse(r2)
if err != nil {
return false
}
if v.Major != base.Major {
return false
}
if !strings.HasSuffix(r, ".x") && !strings.HasSuffix(r, ".*") {
return v.Minor == base.Minor
}
return true
}
// Caret
if strings.HasPrefix(r, "^") {
base, err := Parse(r[1:])
if err != nil {
return false
}
if v.Major != base.Major {
return false
}
return v.Compare(base) >= 0
}
// Tilde
if strings.HasPrefix(r, "~") {
base, err := Parse(r[1:])
if err != nil {
return false
}
if v.Major != base.Major || v.Minor != base.Minor {
return false
}
return v.Compare(base) >= 0
}
// Comparison operators
for _, op := range []string{">=", "<=", ">", "<"} {
if strings.HasPrefix(r, op) {
other, err := Parse(r[len(op):])
if err != nil {
return false
}
cmp := v.Compare(other)
switch op {
case ">=":
return cmp >= 0
case "<=":
return cmp <= 0
case ">":
return cmp > 0
case "<":
return cmp < 0
}
}
}
// Exact
other, err := Parse(r)
if err != nil {
return false
}
return v.Compare(other) == 0
}
