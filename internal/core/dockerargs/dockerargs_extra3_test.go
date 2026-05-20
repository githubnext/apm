package dockerargs_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/dockerargs"
)

func TestProcessDockerArgs_AddsInteractiveAndRM_Extra3(t *testing.T) {
	args := []string{"docker", "run", "image"}
	result := dockerargs.ProcessDockerArgs(args, nil)
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
		t.Error("expected -i flag added")
	}
	if !hasRM {
		t.Error("expected --rm flag added")
	}
}

func TestExtractEnvVars_Basic(t *testing.T) {
	args := []string{"-e", "FOO=bar", "image"}
	clean, env := dockerargs.ExtractEnvVars(args)
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
	if len(clean) != 1 || clean[0] != "image" {
		t.Errorf("unexpected clean args: %v", clean)
	}
}

func TestExtractEnvVars_NoEnvVars(t *testing.T) {
	args := []string{"docker", "run", "image"}
	clean, env := dockerargs.ExtractEnvVars(args)
	if len(env) != 0 {
		t.Errorf("expected empty env, got %v", env)
	}
	if len(clean) != 3 {
		t.Errorf("expected 3 clean args, got %v", clean)
	}
}

func TestMergeEnvVars_MergesCorrectly(t *testing.T) {
	existing := map[string]string{"A": "1"}
	newEnv := map[string]string{"B": "2"}
	merged := dockerargs.MergeEnvVars(existing, newEnv)
	if merged["A"] != "1" {
		t.Errorf("expected A=1, got %q", merged["A"])
	}
	if merged["B"] != "2" {
		t.Errorf("expected B=2, got %q", merged["B"])
	}
}

func TestMergeEnvVars_NewOverridesExisting(t *testing.T) {
	existing := map[string]string{"A": "old"}
	newEnv := map[string]string{"A": "new"}
	merged := dockerargs.MergeEnvVars(existing, newEnv)
	if merged["A"] != "new" {
		t.Errorf("expected A=new, got %q", merged["A"])
	}
}

func TestExtractEnvVars_MultipleEnvVars(t *testing.T) {
	args := []string{"-e", "X=1", "-e", "Y=2"}
	_, env := dockerargs.ExtractEnvVars(args)
	if env["X"] != "1" || env["Y"] != "2" {
		t.Errorf("unexpected env: %v", env)
	}
}
