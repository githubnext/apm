package dockerargs

import (
	"testing"
)

func TestProcessDockerArgs_RMAlreadyPresent(t *testing.T) {
	args := []string{"docker", "run", "--rm", "image"}
	result := ProcessDockerArgs(args, nil)
	count := 0
	for _, a := range result {
		if a == "--rm" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected --rm exactly once, got %d", count)
	}
}

func TestProcessDockerArgs_InteractiveShortAlreadyPresent(t *testing.T) {
	args := []string{"docker", "run", "-i", "image"}
	result := ProcessDockerArgs(args, nil)
	count := 0
	for _, a := range result {
		if a == "-i" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected -i exactly once, got %d", count)
	}
}

func TestProcessDockerArgs_NilEnvVars(t *testing.T) {
	args := []string{"docker", "run", "image"}
	result := ProcessDockerArgs(args, nil)
	if len(result) == 0 {
		t.Error("result should not be empty")
	}
	found := false
	for _, a := range result {
		if a == "image" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'image' in result")
	}
}

func TestExtractEnvVars_MultipleFlags(t *testing.T) {
	args := []string{"run", "-e", "FOO=bar", "-e", "BAZ=qux", "image"}
	clean, env := ExtractEnvVars(args)
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

func TestExtractEnvVars_TrailingEFlag(t *testing.T) {
	// -e at end with no value -- should not panic
	args := []string{"run", "-e"}
	clean, env := ExtractEnvVars(args)
	_ = clean
	_ = env
}

func TestMergeEnvVars_NewWins(t *testing.T) {
	existing := map[string]string{"K": "old"}
	newEnv := map[string]string{"K": "new"}
	merged := MergeEnvVars(existing, newEnv)
	if merged["K"] != "new" {
		t.Errorf("expected new value to win, got %q", merged["K"])
	}
}

func TestMergeEnvVars_BothNil(t *testing.T) {
	merged := MergeEnvVars(nil, nil)
	if merged == nil {
		t.Error("expected non-nil map")
	}
	if len(merged) != 0 {
		t.Errorf("expected empty map, got %v", merged)
	}
}

func TestMergeEnvVars_OnlyNew(t *testing.T) {
	merged := MergeEnvVars(nil, map[string]string{"X": "1"})
	if merged["X"] != "1" {
		t.Errorf("expected X=1, got %q", merged["X"])
	}
}

func TestProcessDockerArgs_EnvVarWithEqualsInValue(t *testing.T) {
	args := []string{"docker", "run", "image"}
	result := ProcessDockerArgs(args, map[string]string{"CONF": "a=b"})
	found := false
	for _, a := range result {
		if a == "CONF=a=b" {
			found = true
		}
	}
	if !found {
		t.Error("expected CONF=a=b in result args")
	}
}
