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
