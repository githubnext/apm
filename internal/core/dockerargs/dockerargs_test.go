package dockerargs_test

import (
	"sort"
	"testing"

	"github.com/githubnext/apm/internal/core/dockerargs"
)

func TestProcessDockerArgs_AddsInteractiveAndRM(t *testing.T) {
	result := dockerargs.ProcessDockerArgs([]string{"docker", "run", "ubuntu"}, nil)
	hasI := false
	hasRM := false
	for _, a := range result {
		if a == "-i" {
			hasI = true
		}
		if a == "--rm" {
			hasRM = true
		}
	}
	if !hasI {
		t.Error("expected -i to be added")
	}
	if !hasRM {
		t.Error("expected --rm to be added")
	}
}

func TestProcessDockerArgs_NoopIfAlreadyPresent(t *testing.T) {
	args := []string{"docker", "run", "-i", "--rm", "ubuntu"}
	result := dockerargs.ProcessDockerArgs(args, nil)
	count := 0
	for _, a := range result {
		if a == "-i" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one -i, got %d", count)
	}
}

func TestProcessDockerArgs_EnvVarsInjected(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	result := dockerargs.ProcessDockerArgs([]string{"docker", "run", "ubuntu"}, env)
	found := false
	for i, a := range result {
		if a == "-e" && i+1 < len(result) && result[i+1] == "FOO=bar" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected -e FOO=bar in %v", result)
	}
}

func TestExtractEnvVars(t *testing.T) {
	args := []string{"docker", "run", "-e", "FOO=bar", "-e", "BAZ=qux", "ubuntu"}
	clean, env := dockerargs.ExtractEnvVars(args)
	if len(env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(env))
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
	if env["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", env["BAZ"])
	}
	for _, a := range clean {
		if a == "-e" {
			t.Error("clean args should not contain -e")
		}
	}
}

func TestExtractEnvVars_NoEqualsSign(t *testing.T) {
	_, env := dockerargs.ExtractEnvVars([]string{"-e", "MYVAR"})
	if env["MYVAR"] != "${MYVAR}" {
		t.Errorf("expected ${MYVAR}, got %q", env["MYVAR"])
	}
}

func TestMergeEnvVars(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "override", "C": "3"}
	merged := dockerargs.MergeEnvVars(a, b)
	if merged["A"] != "1" {
		t.Error("A should be 1")
	}
	if merged["B"] != "override" {
		t.Error("B should be overridden")
	}
	if merged["C"] != "3" {
		t.Error("C should be 3")
	}

	// original maps unchanged
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) != 2 {
		t.Error("original map should be unchanged")
	}
}
