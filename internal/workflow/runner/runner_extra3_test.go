package runner

import (
	"testing"

	"github.com/githubnext/apm/internal/workflow/wfparser"
)

func newWF(params []string) *wfparser.WorkflowDefinition {
	return &wfparser.WorkflowDefinition{
		Description:     "test",
		InputParameters: params,
	}
}

func TestSubstituteParameters_OneKey(t *testing.T) {
	out := SubstituteParameters("hello ${input:name}", map[string]string{"name": "world"})
	if out != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestSubstituteParameters_TwoKeys(t *testing.T) {
	out := SubstituteParameters("${input:a} and ${input:b}", map[string]string{"a": "foo", "b": "bar"})
	if out != "foo and bar" {
		t.Errorf("unexpected result: %q", out)
	}
}

func TestSubstituteParameters_UnknownKey_Unchanged(t *testing.T) {
	out := SubstituteParameters("${input:missing}", map[string]string{})
	if out != "${input:missing}" {
		t.Errorf("unknown placeholder should remain, got %q", out)
	}
}

func TestSubstituteParameters_EmptyContentVariant(t *testing.T) {
	out := SubstituteParameters("", map[string]string{"k": "v"})
	if out != "" {
		t.Errorf("empty content should remain empty, got %q", out)
	}
}

func TestSubstituteParameters_NilParamsVariant(t *testing.T) {
	out := SubstituteParameters("hello", nil)
	if out != "hello" {
		t.Errorf("nil params should leave content unchanged, got %q", out)
	}
}

func TestCollectParameters_FillsDefaults(t *testing.T) {
	wf := newWF([]string{"a", "b"})
	res := CollectParameters(wf, map[string]string{"a": "val"})
	if res["a"] != "val" {
		t.Errorf("expected val, got %q", res["a"])
	}
	if _, ok := res["b"]; !ok {
		t.Error("missing b in result")
	}
}

func TestCollectParameters_NoParamsVariant(t *testing.T) {
	wf := newWF(nil)
	res := CollectParameters(wf, nil)
	if res == nil {
		t.Error("expected non-nil map")
	}
}

func TestCollectParameters_ExtraProvidedVariant(t *testing.T) {
	wf := newWF([]string{"x"})
	res := CollectParameters(wf, map[string]string{"x": "1", "extra": "2"})
	if res["extra"] != "2" {
		t.Errorf("extra provided parameter should be preserved, got %q", res["extra"])
	}
}

func TestRunResult_ZeroValue_NotSuccess(t *testing.T) {
	var r RunResult
	if r.Success {
		t.Error("zero-value RunResult should not be successful")
	}
}

func TestRunResult_SuccessFlag(t *testing.T) {
	r := RunResult{Success: true, Output: "done"}
	if !r.Success {
		t.Error("expected success")
	}
}

func TestRunResult_ErrorMsg(t *testing.T) {
	r := RunResult{ErrorMsg: "something went wrong"}
	if r.ErrorMsg == "" {
		t.Error("expected error message")
	}
}

func TestFindWorkflowByName_NotFoundVariant(t *testing.T) {
	dir := t.TempDir()
	_, err := FindWorkflowByName("nonexistent-workflow", dir)
	if err == nil {
		t.Error("expected error for non-existent workflow")
	}
}

func TestPreviewWorkflow_NotFoundVariant(t *testing.T) {
	dir := t.TempDir()
	r := PreviewWorkflow("nope", nil, dir)
	if r.Success {
		t.Error("expected failure for missing workflow")
	}
	if r.ErrorMsg == "" {
		t.Error("expected error message")
	}
}
