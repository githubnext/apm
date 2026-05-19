package base_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/base"
)

func TestInputVarRE_SpecialCharsInName(t *testing.T) {
	cases := []struct {
		input string
		match bool
	}{
		{"${input:MY_VAR_123}", true},
		{"${input:a-b}", true},  // hyphens allowed in input names
		{"${input:}", false},    // empty name - depends on regex
	}
	for _, c := range cases {
		m := base.InputVarRE.FindStringSubmatch(c.input)
		if c.match && m == nil {
			t.Errorf("InputVarRE: expected match for %q", c.input)
		}
	}
}

func TestEnvVarRE_UnderscoreOnly(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${_UNDERSCORE}")
	if m == nil {
		t.Fatal("expected match for underscore-prefixed var")
	}
	if m[1] != "_UNDERSCORE" {
		t.Errorf("expected _UNDERSCORE, got %q", m[1])
	}
}

func TestEnvVarRE_NumbersInName(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${VAR_123}")
	if m == nil {
		t.Fatal("expected match for var with numbers")
	}
	if m[1] != "VAR_123" {
		t.Errorf("expected VAR_123, got %q", m[1])
	}
}

func TestEnvVarRE_EnvPrefixWithUnderscore(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${env:_PRIVATE}")
	if m == nil {
		t.Fatal("expected match for env:_PRIVATE")
	}
	if m[1] != "_PRIVATE" {
		t.Errorf("expected _PRIVATE, got %q", m[1])
	}
}

func TestInputVarRE_NotMatchEnvPrefix(t *testing.T) {
	cases := []string{
		"${env:SECRET}",
		"${SECRET}",
		"${ENV:KEY}",
	}
	for _, c := range cases {
		if base.InputVarRE.MatchString(c) {
			t.Errorf("InputVarRE should not match %q", c)
		}
	}
}

func TestEnvVarRE_DoesNotMatchInputPrefix(t *testing.T) {
	cases := []string{
		"${input:VAR}",
		"${input:MY_VAR}",
	}
	for _, c := range cases {
		if base.EnvVarRE.MatchString(c) {
			t.Errorf("EnvVarRE should not match input placeholder %q", c)
		}
	}
}

func TestInputVarRE_MixedContent(t *testing.T) {
	input := "text ${input:VAR1} more text ${input:VAR2} end"
	matches := base.InputVarRE.FindAllStringSubmatch(input, -1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 input matches, got %d", len(matches))
	}
	names := []string{matches[0][1], matches[1][1]}
	if names[0] != "VAR1" || names[1] != "VAR2" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestEnvVarRE_MixedContent(t *testing.T) {
	input := "prefix ${FOO} middle ${env:BAR} suffix"
	matches := base.EnvVarRE.FindAllStringSubmatch(input, -1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 env matches, got %d", len(matches))
	}
	if matches[0][1] != "FOO" {
		t.Errorf("first match: want FOO, got %s", matches[0][1])
	}
	if matches[1][1] != "BAR" {
		t.Errorf("second match: want BAR, got %s", matches[1][1])
	}
}

func TestEnvVarRE_SingleLetter(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${A}")
	if m == nil {
		t.Fatal("expected match for single-letter var")
	}
	if m[1] != "A" {
		t.Errorf("expected A, got %q", m[1])
	}
}

func TestInputVarRE_SingleLetter(t *testing.T) {
	m := base.InputVarRE.FindStringSubmatch("${input:x}")
	if m == nil {
		t.Fatal("expected match for single-letter input var")
	}
	if m[1] != "x" {
		t.Errorf("expected x, got %q", m[1])
	}
}

func TestEnvVarRE_LowercaseName(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${my_var}")
	if m == nil {
		t.Fatal("expected match for lowercase var name")
	}
	if m[1] != "my_var" {
		t.Errorf("expected my_var, got %q", m[1])
	}
}

func TestInputVarRE_ExactCapture(t *testing.T) {
	// Verify we capture exactly the name part, not the delimiters
	m := base.InputVarRE.FindStringSubmatch("${input:EXACTLY_THIS}")
	if m == nil {
		t.Fatal("expected match")
	}
	if m[1] != "EXACTLY_THIS" {
		t.Errorf("expected EXACTLY_THIS, got %q", m[1])
	}
	// The full match should include the delimiters
	if m[0] != "${input:EXACTLY_THIS}" {
		t.Errorf("full match should be ${input:EXACTLY_THIS}, got %q", m[0])
	}
}

func TestEnvVarRE_ExactCapture(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${EXACTLY_THIS}")
	if m == nil {
		t.Fatal("expected match")
	}
	if m[1] != "EXACTLY_THIS" {
		t.Errorf("expected EXACTLY_THIS, got %q", m[1])
	}
	if m[0] != "${EXACTLY_THIS}" {
		t.Errorf("full match should be ${EXACTLY_THIS}, got %q", m[0])
	}
}

func TestEnvVarRE_EnvPrefixExactCapture(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${env:MY_SECRET}")
	if m == nil {
		t.Fatal("expected match")
	}
	if m[1] != "MY_SECRET" {
		t.Errorf("expected MY_SECRET, got %q", m[1])
	}
}
