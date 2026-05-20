package runner

import (
	"testing"

	"github.com/githubnext/apm/internal/workflow/wfparser"
)

func TestSubstituteParameters_EmptyMapStable(t *testing.T) {
	result := SubstituteParameters("hello ${input:name}", nil)
	if result != "hello ${input:name}" {
		t.Errorf("nil map should leave content unchanged, got %q", result)
	}
}

func TestSubstituteParameters_SingleReplacement(t *testing.T) {
	result := SubstituteParameters("run ${input:cmd}", map[string]string{"cmd": "go test"})
	if result != "run go test" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestCollectParameters_MissingOneMissing(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{InputParameters: []string{"x", "y"}}
	result := CollectParameters(wf, map[string]string{"x": "val"})
	if result["x"] != "val" {
		t.Errorf("expected x=val, got %q", result["x"])
	}
	if _, ok := result["y"]; !ok {
		t.Error("y should be present with empty value")
	}
	if result["y"] != "" {
		t.Errorf("missing param should be empty, got %q", result["y"])
	}
}

func TestCollectParameters_EmptyInput(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{}
	result := CollectParameters(wf, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestRunResult_SuccessTrue(t *testing.T) {
	r := RunResult{Success: true, Output: "ok"}
	if !r.Success {
		t.Error("expected success=true")
	}
	if r.ErrorMsg != "" {
		t.Error("expected empty error msg")
	}
}

func TestRunResult_FailedWithError(t *testing.T) {
	r := RunResult{Success: false, ErrorMsg: "something failed"}
	if r.Success {
		t.Error("expected success=false")
	}
	if r.ErrorMsg == "" {
		t.Error("expected non-empty error msg")
	}
}

func TestSubstituteParameters_PlaceholderRepeated(t *testing.T) {
	result := SubstituteParameters("${input:x} and ${input:x}", map[string]string{"x": "A"})
	if result != "A and A" {
		t.Errorf("expected repeated replacement, got %q", result)
	}
}
