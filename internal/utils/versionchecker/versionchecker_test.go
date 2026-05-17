package versionchecker_test

import (
"testing"

"github.com/githubnext/apm/internal/utils/versionchecker"
)

func TestParseVersion_Valid(t *testing.T) {
cases := []struct {
input string
major int
minor int
patch int
}{
{"1.2.3", 1, 2, 3},
{"0.0.1", 0, 0, 1},
{"10.20.30", 10, 20, 30},
}
for _, tc := range cases {
v := versionchecker.ParseVersion(tc.input)
if v == nil {
t.Errorf("ParseVersion(%q) returned nil", tc.input)
continue
}
if v.Major != tc.major || v.Minor != tc.minor || v.Patch != tc.patch {
t.Errorf("ParseVersion(%q) = %+v, want {%d,%d,%d}", tc.input, v, tc.major, tc.minor, tc.patch)
}
}
}

func TestParseVersion_Invalid(t *testing.T) {
for _, input := range []string{"", "not-a-version", "1.2"} {
v := versionchecker.ParseVersion(input)
if v != nil {
t.Errorf("ParseVersion(%q) expected nil, got %+v", input, v)
}
}
}

func TestIsNewerVersion_NewerLatest(t *testing.T) {
if !versionchecker.IsNewerVersion("1.0.0", "1.0.1") {
t.Error("1.0.1 should be newer than 1.0.0")
}
}

func TestIsNewerVersion_SameVersion(t *testing.T) {
if versionchecker.IsNewerVersion("1.0.0", "1.0.0") {
t.Error("same version should not be newer")
}
}

func TestIsNewerVersion_OlderLatest(t *testing.T) {
if versionchecker.IsNewerVersion("1.2.0", "1.1.0") {
t.Error("1.1.0 should not be newer than 1.2.0")
}
}

func TestIsNewerVersion_MajorBump(t *testing.T) {
if !versionchecker.IsNewerVersion("1.9.9", "2.0.0") {
t.Error("2.0.0 should be newer than 1.9.9")
}
}

func TestIsNewerVersion_MinorBump(t *testing.T) {
if !versionchecker.IsNewerVersion("1.0.0", "1.1.0") {
t.Error("1.1.0 should be newer than 1.0.0")
}
}

func TestIsNewerVersion_InvalidCurrent(t *testing.T) {
if versionchecker.IsNewerVersion("not-a-version", "1.0.0") {
t.Error("invalid current version should return false")
}
}

func TestIsNewerVersion_InvalidLatest(t *testing.T) {
if versionchecker.IsNewerVersion("1.0.0", "not-a-version") {
t.Error("invalid latest version should return false")
}
}

func TestIsNewerVersion_PreReleaseLower(t *testing.T) {
if !versionchecker.IsNewerVersion("1.0.0rc1", "1.0.0") {
t.Error("1.0.0 stable should be newer than 1.0.0rc1")
}
}

func TestIsNewerVersion_StableNotNewerThanPreRelease(t *testing.T) {
if versionchecker.IsNewerVersion("1.0.0", "1.0.0rc1") {
t.Error("1.0.0rc1 should not be newer than 1.0.0 stable")
}
}

func TestParseVersion_Prerelease(t *testing.T) {
v := versionchecker.ParseVersion("1.2.3rc1")
if v == nil {
t.Fatal("ParseVersion returned nil for 1.2.3rc1")
}
if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
t.Errorf("unexpected version: %+v", v)
}
if v.Prerelease != "rc1" {
t.Errorf("expected Prerelease=rc1, got %q", v.Prerelease)
}
}

func TestParseVersion_BetaPrerelease(t *testing.T) {
v := versionchecker.ParseVersion("0.5.0b2")
if v == nil {
t.Fatal("ParseVersion returned nil for 0.5.0b2")
}
if v.Prerelease != "b2" {
t.Errorf("expected Prerelease=b2, got %q", v.Prerelease)
}
}

func TestParseVersion_StableHasNoPrerelease(t *testing.T) {
v := versionchecker.ParseVersion("2.0.0")
if v == nil {
t.Fatal("ParseVersion returned nil")
}
if v.Prerelease != "" {
t.Errorf("expected empty Prerelease, got %q", v.Prerelease)
}
}

func TestVersionComponents_ZeroValues(t *testing.T) {
v := versionchecker.ParseVersion("0.0.0")
if v == nil {
t.Fatal("ParseVersion returned nil for 0.0.0")
}
if v.Major != 0 || v.Minor != 0 || v.Patch != 0 {
t.Errorf("expected all zeros, got %+v", v)
}
}
