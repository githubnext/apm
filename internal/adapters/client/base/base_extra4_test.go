package base_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/base"
)

func TestInputVarRE_MultipleMatches_Extra4(t *testing.T) {
	s := "${input:foo} and ${input:bar}"
	matches := base.InputVarRE.FindAllString(s, -1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestInputVarRE_NotMatchEnvVar_Extra4(t *testing.T) {
	s := "${MY_VAR}"
	if base.InputVarRE.MatchString(s) {
		t.Fatal("InputVarRE should not match ${MY_VAR}")
	}
}

func TestEnvVarRE_MatchBareVar_Extra4(t *testing.T) {
	s := "${MY_VAR}"
	if !base.EnvVarRE.MatchString(s) {
		t.Fatal("EnvVarRE should match ${MY_VAR}")
	}
}

func TestEnvVarRE_MatchEnvColon_Extra4(t *testing.T) {
	s := "${env:MY_VAR}"
	if !base.EnvVarRE.MatchString(s) {
		t.Fatal("EnvVarRE should match ${env:MY_VAR}")
	}
}

func TestEnvVarRE_NoMatchInput_Extra4(t *testing.T) {
	s := "${input:MY_VAR}"
	if base.EnvVarRE.MatchString(s) {
		t.Fatal("EnvVarRE should not match ${input:MY_VAR}")
	}
}

func TestEnvVarRE_CapturesBareVarName_Extra4(t *testing.T) {
	s := "${TOKEN}"
	m := base.EnvVarRE.FindStringSubmatch(s)
	if len(m) < 2 || m[1] != "TOKEN" {
		t.Fatalf("expected capture TOKEN, got %v", m)
	}
}

func TestEnvVarRE_CapturesEnvColonVarName_Extra4(t *testing.T) {
	s := "${env:SECRET_KEY}"
	m := base.EnvVarRE.FindStringSubmatch(s)
	if len(m) < 2 || m[1] != "SECRET_KEY" {
		t.Fatalf("expected capture SECRET_KEY, got %v", m)
	}
}

func TestInputVarRE_NoMatchGitHubActions_Extra4(t *testing.T) {
	s := "${{ secrets.TOKEN }}"
	if base.InputVarRE.MatchString(s) {
		t.Fatal("InputVarRE should not match GitHub Actions syntax")
	}
}

func TestEnvVarRE_MultipleVars_Extra4(t *testing.T) {
	s := "${A} and ${env:B}"
	all := base.EnvVarRE.FindAllStringSubmatch(s, -1)
	if len(all) != 2 {
		t.Fatalf("expected 2 env var matches, got %d", len(all))
	}
}

func TestInputVarRE_CapturesMultipleNames_Extra4(t *testing.T) {
	s := "${input:alpha}${input:beta}"
	all := base.InputVarRE.FindAllStringSubmatch(s, -1)
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}
	if all[0][1] != "alpha" {
		t.Errorf("expected alpha, got %s", all[0][1])
	}
	if all[1][1] != "beta" {
		t.Errorf("expected beta, got %s", all[1][1])
	}
}
