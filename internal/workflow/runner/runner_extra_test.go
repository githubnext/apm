package runner

import (
	"testing"

	"github.com/githubnext/apm/internal/workflow/wfparser"
)

func TestSubstituteParameters_SpecialChars(t *testing.T) {
	content := "value=${input:key}&other=1"
	params := map[string]string{"key": "hello world"}
	result := SubstituteParameters(content, params)
	if result != "value=hello world&other=1" {
		t.Errorf("unexpected: %q", result)
	}
}

func TestSubstituteParameters_MultipleKeys(t *testing.T) {
	content := "${input:a}-${input:b}-${input:c}"
	params := map[string]string{"a": "x", "b": "y", "c": "z"}
	result := SubstituteParameters(content, params)
	if result != "x-y-z" {
		t.Errorf("unexpected: %q", result)
	}
}

func TestSubstituteParameters_NoPlaceholders(t *testing.T) {
	content := "plain text without placeholders"
	result := SubstituteParameters(content, map[string]string{"k": "v"})
	if result != content {
		t.Errorf("expected unchanged: %q", result)
	}
}

func TestSubstituteParameters_EmptyValue(t *testing.T) {
	content := "prefix-${input:key}-suffix"
	result := SubstituteParameters(content, map[string]string{"key": ""})
	if result != "prefix--suffix" {
		t.Errorf("unexpected: %q", result)
	}
}

func TestCollectParameters_AllProvided(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{
		InputParameters: []string{"a", "b"},
	}
	provided := map[string]string{"a": "1", "b": "2"}
	result := CollectParameters(wf, provided)
	if result["a"] != "1" || result["b"] != "2" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestCollectParameters_ExtraProvided(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{
		InputParameters: []string{"x"},
	}
	provided := map[string]string{"x": "1", "y": "2"}
	result := CollectParameters(wf, provided)
	if result["x"] != "1" {
		t.Error("expected x=1")
	}
	if result["y"] != "2" {
		t.Error("expected y=2 from provided")
	}
}

func TestCollectParameters_NoParams(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{}
	result := CollectParameters(wf, map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestRunResult_ZeroValue(t *testing.T) {
	var r RunResult
	if r.Success {
		t.Error("zero value should not be successful")
	}
	if r.Output != "" {
		t.Error("zero value output should be empty")
	}
	if r.ErrorMsg != "" {
		t.Error("zero value error should be empty")
	}
}

func TestFindWorkflowByName_EmptyBase(t *testing.T) {
	_, err := FindWorkflowByName("nonexistent-workflow-xyzabc", "")
	if err == nil {
		t.Error("expected error for nonexistent workflow in cwd")
	}
}

func TestPreviewWorkflow_EmptyBase(t *testing.T) {
	result := PreviewWorkflow("nonexistent-workflow-xyzabc", nil, "")
	if result.Success {
		t.Error("expected failure for nonexistent workflow")
	}
}

func TestRunWorkflow_NilExecutor(t *testing.T) {
	result := RunWorkflow("nonexistent-wf", nil, "/tmp/gh-aw/agent", nil)
	if result.Success {
		t.Error("expected failure with nil executor and missing workflow")
	}
}

func TestSubstituteParameters_OverlappingKeys(t *testing.T) {
	// Make sure key "ab" doesn't partially match "${input:a}"
	content := "${input:a} ${input:ab}"
	params := map[string]string{"a": "X", "ab": "Y"}
	result := SubstituteParameters(content, params)
	if result != "X Y" {
		t.Errorf("unexpected: %q", result)
	}
}
