package dockerargs_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/dockerargs"
)

func TestProcessDockerArgs_EmptyArgs(t *testing.T) {
	result := dockerargs.ProcessDockerArgs(nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result for nil input")
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %v", result)
	}
}

func TestProcessDockerArgs_NoRunCommand(t *testing.T) {
	args := []string{"docker", "pull", "ubuntu"}
	result := dockerargs.ProcessDockerArgs(args, nil)
	for _, a := range result {
		if a == "-i" || a == "--rm" {
			t.Errorf("should not inject -i/--rm without 'run' command, got %v", result)
		}
	}
}

func TestProcessDockerArgs_InteractiveAlreadyLong(t *testing.T) {
	args := []string{"docker", "run", "--interactive", "ubuntu"}
	result := dockerargs.ProcessDockerArgs(args, nil)
	count := 0
	for _, a := range result {
		if a == "-i" || a == "--interactive" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one interactive flag, got %d", count)
	}
}

func TestProcessDockerArgs_MultipleEnvVars(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	result := dockerargs.ProcessDockerArgs([]string{"docker", "run", "ubuntu"}, env)
	envCount := 0
	for i, a := range result {
		if a == "-e" && i+1 < len(result) {
			envCount++
		}
	}
	if envCount != 3 {
		t.Errorf("expected 3 -e flags, got %d: %v", envCount, result)
	}
}

func TestProcessDockerArgs_OrderPreserved(t *testing.T) {
	args := []string{"docker", "run", "ubuntu", "bash"}
	result := dockerargs.ProcessDockerArgs(args, nil)
	if result[0] != "docker" {
		t.Errorf("expected 'docker' first, got %q", result[0])
	}
	last := result[len(result)-1]
	if last != "bash" {
		t.Errorf("expected 'bash' last, got %q", last)
	}
}

func TestExtractEnvVars_Empty(t *testing.T) {
	clean, env := dockerargs.ExtractEnvVars(nil)
	if len(clean) != 0 {
		t.Errorf("expected empty clean args, got %v", clean)
	}
	if len(env) != 0 {
		t.Errorf("expected empty env map, got %v", env)
	}
}

func TestExtractEnvVars_NoEnvFlags(t *testing.T) {
	args := []string{"docker", "run", "ubuntu"}
	clean, env := dockerargs.ExtractEnvVars(args)
	if len(env) != 0 {
		t.Errorf("expected no env vars, got %v", env)
	}
	if len(clean) != 3 {
		t.Errorf("expected 3 clean args, got %v", clean)
	}
}

func TestExtractEnvVars_ValueWithEquals(t *testing.T) {
	// Value itself contains '='
	_, env := dockerargs.ExtractEnvVars([]string{"-e", "FOO=a=b"})
	if env["FOO"] != "a=b" {
		t.Errorf("expected 'a=b', got %q", env["FOO"])
	}
}

func TestMergeEnvVars_Empty(t *testing.T) {
	merged := dockerargs.MergeEnvVars(nil, nil)
	if len(merged) != 0 {
		t.Errorf("expected empty map, got %v", merged)
	}
}

func TestMergeEnvVars_OnlyExisting(t *testing.T) {
	existing := map[string]string{"X": "10"}
	merged := dockerargs.MergeEnvVars(existing, nil)
	if merged["X"] != "10" {
		t.Errorf("expected X=10, got %q", merged["X"])
	}
	if len(merged) != 1 {
		t.Errorf("expected 1 entry, got %d", len(merged))
	}
}

func TestMergeEnvVars_DoesNotMutateOriginal(t *testing.T) {
	a := map[string]string{"K": "v1"}
	b := map[string]string{"K": "v2"}
	dockerargs.MergeEnvVars(a, b)
	if a["K"] != "v1" {
		t.Error("MergeEnvVars should not mutate original map")
	}
}
