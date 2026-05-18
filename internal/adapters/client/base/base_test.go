package base_test

import (
"testing"

"github.com/githubnext/apm/internal/adapters/client/base"
)

func TestInputVarRE(t *testing.T) {
cases := []struct {
input string
match bool
name  string
}{
{"${input:MY_VAR}", true, "MY_VAR"},
{"${input:foo}", true, "foo"},
{"${env:BAR}", false, ""},
{"${BAR}", false, ""},
{"no placeholder", false, ""},
{"${input:a} and ${input:b}", true, "a"},
}
for _, c := range cases {
m := base.InputVarRE.FindStringSubmatch(c.input)
if c.match {
if m == nil {
t.Errorf("InputVarRE: expected match for %q", c.input)
} else if m[1] != c.name {
t.Errorf("InputVarRE: got name %q, want %q", m[1], c.name)
}
} else {
if m != nil {
t.Errorf("InputVarRE: expected no match for %q, got %v", c.input, m)
}
}
}
}

func TestEnvVarRE(t *testing.T) {
cases := []struct {
input string
match bool
name  string
}{
{"${MY_VAR}", true, "MY_VAR"},
{"${env:MY_VAR}", true, "MY_VAR"},
{"${input:foo}", false, ""},
{"${{ ctx.token }}", false, ""},
{"no placeholder", false, ""},
{"${A_1}", true, "A_1"},
}
for _, c := range cases {
m := base.EnvVarRE.FindStringSubmatch(c.input)
if c.match {
if m == nil {
t.Errorf("EnvVarRE: expected match for %q", c.input)
} else if m[1] != c.name {
t.Errorf("EnvVarRE: got name %q, want %q", m[1], c.name)
}
} else {
if m != nil {
t.Errorf("EnvVarRE: expected no match for %q, got %v", c.input, m)
}
}
}
}

func TestEnvVarREAllMatches(t *testing.T) {
input := "${FOO} and ${env:BAR} and ${input:skip}"
matches := base.EnvVarRE.FindAllStringSubmatch(input, -1)
if len(matches) != 2 {
t.Fatalf("expected 2 matches, got %d", len(matches))
}
if matches[0][1] != "FOO" {
t.Errorf("first match: want FOO, got %s", matches[0][1])
}
if matches[1][1] != "BAR" {
t.Errorf("second match: want BAR, got %s", matches[1][1])
}
}

func TestInputVarRE_MultipleMatches(t *testing.T) {
input := "${input:FOO} and ${input:BAR}"
matches := base.InputVarRE.FindAllStringSubmatch(input, -1)
if len(matches) != 2 {
t.Fatalf("expected 2 matches, got %d", len(matches))
}
if matches[0][1] != "FOO" {
t.Errorf("first match: want FOO, got %s", matches[0][1])
}
if matches[1][1] != "BAR" {
t.Errorf("second match: want BAR, got %s", matches[1][1])
}
}

func TestInputVarRE_EnvNotMatched(t *testing.T) {
cases := []string{"${MY_VAR}", "${env:MY_VAR}", "${{ secrets.TOKEN }}"}
for _, c := range cases {
if base.InputVarRE.MatchString(c) {
t.Errorf("InputVarRE should not match %q", c)
}
}
}

func TestEnvVarRE_NoMatchGitHubActions(t *testing.T) {
cases := []string{"${{ secrets.TOKEN }}", "${{ env.VAR }}", "literal"}
for _, c := range cases {
if base.EnvVarRE.MatchString(c) {
t.Errorf("EnvVarRE should not match %q", c)
}
}
}

func TestEnvVarRE_CaseSensitive(t *testing.T) {
// Variable names are case-sensitive in the regex
if !base.EnvVarRE.MatchString("${MY_VAR}") {
t.Error("expected match for ${MY_VAR}")
}
}

func TestEnvVarRE_WithPrefix(t *testing.T) {
m := base.EnvVarRE.FindStringSubmatch("${env:SECRET_KEY}")
if m == nil {
t.Fatal("expected match for ${env:SECRET_KEY}")
}
if m[1] != "SECRET_KEY" {
t.Errorf("expected SECRET_KEY, got %q", m[1])
}
}

func TestEnvVarRE_DigitStartNotMatched(t *testing.T) {
// Variable names cannot start with a digit
if base.EnvVarRE.MatchString("${1VAR}") {
t.Error("EnvVarRE should not match variable starting with digit")
}
}
