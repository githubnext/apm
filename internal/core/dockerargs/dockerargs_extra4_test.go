package dockerargs

import (
"testing"
)

func TestProcessDockerArgs_RunWithRM(t *testing.T) {
result := ProcessDockerArgs([]string{"run", "image"}, nil)
found := false
for _, a := range result {
if a == "--rm" {
found = true
}
}
if !found {
t.Error("expected --rm in result")
}
}

func TestProcessDockerArgs_EnvVarsInjectedAfterRun(t *testing.T) {
result := ProcessDockerArgs([]string{"run", "image"}, map[string]string{"FOO": "bar"})
found := false
for i, a := range result {
if a == "-e" && i+1 < len(result) {
found = true
}
}
if !found {
t.Error("expected -e flag for env vars")
}
}

func TestExtractEnvVars_NoFlags(t *testing.T) {
clean, env := ExtractEnvVars([]string{"image:latest"})
if len(env) != 0 {
t.Errorf("expected no env vars, got %v", env)
}
if len(clean) == 0 {
t.Error("expected non-empty clean args")
}
}

func TestExtractEnvVars_EqualSyntax(t *testing.T) {
_, env := ExtractEnvVars([]string{"-e", "KEY=val", "image"})
if env["KEY"] != "val" {
t.Errorf("expected KEY=val, got %v", env)
}
}

func TestMergeEnvVars_EmptyMaps(t *testing.T) {
result := MergeEnvVars(map[string]string{}, map[string]string{})
if result == nil {
t.Error("expected non-nil map")
}
}

func TestMergeEnvVars_ExistingPreserved(t *testing.T) {
result := MergeEnvVars(map[string]string{"A": "1"}, map[string]string{"B": "2"})
if result["A"] != "1" {
t.Errorf("expected A=1, got %q", result["A"])
}
if result["B"] != "2" {
t.Errorf("expected B=2, got %q", result["B"])
}
}
