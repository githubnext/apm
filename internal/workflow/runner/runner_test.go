package runner

import (
	"testing"

	"github.com/githubnext/apm/internal/workflow/wfparser"
)

func TestSubstituteParameters_Basic(t *testing.T) {
	content := "Hello ${input:name}, you are ${input:age} years old."
	params := map[string]string{"name": "Alice", "age": "30"}
	result := SubstituteParameters(content, params)
	if result != "Hello Alice, you are 30 years old." {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestSubstituteParameters_MissingKey(t *testing.T) {
	content := "Hello ${input:name}"
	params := map[string]string{}
	result := SubstituteParameters(content, params)
	if result != "Hello ${input:name}" {
		t.Errorf("expected unchanged content, got: %q", result)
	}
}

func TestSubstituteParameters_EmptyContent(t *testing.T) {
	result := SubstituteParameters("", map[string]string{"k": "v"})
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func TestSubstituteParameters_EmptyParams(t *testing.T) {
	content := "no params here"
	result := SubstituteParameters(content, nil)
	if result != content {
		t.Errorf("expected unchanged content, got: %q", result)
	}
}

func TestSubstituteParameters_MultipleOccurrences(t *testing.T) {
	content := "${input:x} and ${input:x} again"
	params := map[string]string{"x": "hello"}
	result := SubstituteParameters(content, params)
	if result != "hello and hello again" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestCollectParameters_ProvidesDefaults(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{
		InputParameters: []string{"name", "age"},
	}
	provided := map[string]string{"name": "Alice"}
	result := CollectParameters(wf, provided)
	if result["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %q", result["name"])
	}
	if _, ok := result["age"]; !ok {
		t.Error("expected age key to be present")
	}
}

func TestCollectParameters_OverridesDefault(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{
		InputParameters: []string{"name"},
	}
	provided := map[string]string{"name": "Bob"}
	result := CollectParameters(wf, provided)
	if result["name"] != "Bob" {
		t.Errorf("expected name=Bob, got %q", result["name"])
	}
}

func TestCollectParameters_EmptyWorkflow(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{}
	result := CollectParameters(wf, map[string]string{"extra": "val"})
	if result["extra"] != "val" {
		t.Errorf("expected extra=val, got %q", result["extra"])
	}
}

func TestSubstituteParameters_AllReplaced(t *testing.T) {
	content := "${input:a}/${input:b}/${input:c}"
	params := map[string]string{"a": "1", "b": "2", "c": "3"}
	result := SubstituteParameters(content, params)
	if result != "1/2/3" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestRunWorkflow_NoExecutor(t *testing.T) {
	result := RunWorkflow("nonexistent", nil, "/tmp", nil)
	if result.Success {
		t.Error("expected failure with no executor")
	}
	if result.ErrorMsg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestPreviewWorkflow_NotFound(t *testing.T) {
	result := PreviewWorkflow("nonexistent-workflow-xyz", nil, "/tmp")
	if result.Success {
		t.Error("expected failure for nonexistent workflow")
	}
	if result.ErrorMsg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestFindWorkflowByName_NotFound(t *testing.T) {
	_, err := FindWorkflowByName("nonexistent-abc", "/tmp")
	if err == nil {
		t.Error("expected error for nonexistent workflow")
	}
}

func TestRunResult_Fields(t *testing.T) {
	r := RunResult{Success: true, Output: "hello", ErrorMsg: ""}
	if !r.Success {
		t.Error("expected Success=true")
	}
	if r.Output != "hello" {
		t.Errorf("unexpected output: %q", r.Output)
	}

	r2 := RunResult{ErrorMsg: "some error"}
	if r2.Success {
		t.Error("expected Success=false")
	}
	if r2.ErrorMsg == "" {
		t.Error("expected non-empty error")
	}
}

func TestRunWorkflow_ExecutorError(t *testing.T) {
	// No real workflow files in /tmp — expect "not found" error before executor
	result := RunWorkflow("bad-wf", nil, "/tmp", func(content, model string) (string, error) {
		return "", nil
	})
	if result.Success {
		t.Error("expected failure for missing workflow")
	}
}

func TestCollectParameters_NilProvided(t *testing.T) {
	wf := &wfparser.WorkflowDefinition{
		InputParameters: []string{"x", "y"},
	}
	result := CollectParameters(wf, nil)
	if _, ok := result["x"]; !ok {
		t.Error("expected x key")
	}
	if _, ok := result["y"]; !ok {
		t.Error("expected y key")
	}
}
