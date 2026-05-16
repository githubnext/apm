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
