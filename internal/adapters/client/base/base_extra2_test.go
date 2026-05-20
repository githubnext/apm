package base_test

import (
	"regexp"
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/base"
)

func TestInputVarRE_NoMatch_PlainText(t *testing.T) {
	if base.InputVarRE.MatchString("plain text without vars") {
		t.Error("should not match plain text")
	}
}

func TestInputVarRE_NoMatch_EmptyBraces(t *testing.T) {
	if base.InputVarRE.MatchString("${}") {
		t.Error("should not match empty braces")
	}
}

func TestInputVarRE_Match_SimpleInput(t *testing.T) {
	if !base.InputVarRE.MatchString("${input:myvar}") {
		t.Error("should match input var")
	}
}

func TestInputVarRE_CaptureGroup(t *testing.T) {
	m := base.InputVarRE.FindStringSubmatch("${input:TOKEN}")
	if len(m) < 2 || m[1] != "TOKEN" {
		t.Errorf("capture group expected TOKEN, got %v", m)
	}
}

func TestInputVarRE_Multiword(t *testing.T) {
	m := base.InputVarRE.FindAllStringSubmatch("use ${input:a} and ${input:b}", -1)
	if len(m) != 2 {
		t.Errorf("expected 2 matches, got %d", len(m))
	}
	names := []string{m[0][1], m[1][1]}
	if names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected names %v", names)
	}
}

func TestEnvVarRE_NoMatchEmptyBraces(t *testing.T) {
	if base.EnvVarRE.MatchString("${}") {
		t.Error("should not match empty braces")
	}
}

func TestEnvVarRE_NoMatchInputPrefix(t *testing.T) {
	if base.EnvVarRE.MatchString("${input:MYVAR}") {
		t.Error("should not match input: prefix")
	}
}

func TestEnvVarRE_Match_PlainBraces(t *testing.T) {
	if !base.EnvVarRE.MatchString("${MYVAR}") {
		t.Error("should match ${MYVAR}")
	}
}

func TestEnvVarRE_Match_EnvPrefix(t *testing.T) {
	if !base.EnvVarRE.MatchString("${env:MYVAR}") {
		t.Error("should match ${env:MYVAR}")
	}
}

func TestEnvVarRE_NoMatchGitHubExpressions(t *testing.T) {
	// GitHub Actions ${{ ... }} should not be matched.
	if base.EnvVarRE.MatchString("${{ secrets.TOKEN }}") {
		t.Error("should not match GitHub Actions expressions")
	}
}

func TestEnvVarRE_CaptureGroupPlain(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${TOKEN}")
	if len(m) < 2 || m[1] != "TOKEN" {
		t.Errorf("expected TOKEN, got %v", m)
	}
}

func TestEnvVarRE_CaptureGroupEnvPrefix(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${env:TOKEN}")
	if len(m) < 2 || m[1] != "TOKEN" {
		t.Errorf("expected TOKEN, got %v", m)
	}
}

func TestInputVarRE_IsCompiled(t *testing.T) {
	if base.InputVarRE == nil {
		t.Fatal("InputVarRE should not be nil")
	}
	if _, ok := interface{}(base.InputVarRE).(*regexp.Regexp); !ok {
		t.Error("InputVarRE should be *regexp.Regexp")
	}
}

func TestEnvVarRE_IsCompiled(t *testing.T) {
	if base.EnvVarRE == nil {
		t.Fatal("EnvVarRE should not be nil")
	}
	if _, ok := interface{}(base.EnvVarRE).(*regexp.Regexp); !ok {
		t.Error("EnvVarRE should be *regexp.Regexp")
	}
}

func TestInputVarRE_WithHyphensInName(t *testing.T) {
	// hyphens are valid inside input var names
	if !base.InputVarRE.MatchString("${input:my-var}") {
		t.Error("should match input var with hyphens")
	}
	m := base.InputVarRE.FindStringSubmatch("${input:my-var}")
	if len(m) < 2 || m[1] != "my-var" {
		t.Errorf("capture expected 'my-var', got %v", m)
	}
}

func TestEnvVarRE_NoMatchDigitStart(t *testing.T) {
	if base.EnvVarRE.MatchString("${1VAR}") {
		t.Error("should not match var starting with digit")
	}
}

func TestEnvVarRE_MultipleMatches(t *testing.T) {
	ms := base.EnvVarRE.FindAllStringSubmatch("${A} text ${env:B}", -1)
	if len(ms) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(ms))
	}
	if ms[0][1] != "A" || ms[1][1] != "B" {
		t.Errorf("expected A,B got %v %v", ms[0][1], ms[1][1])
	}
}
