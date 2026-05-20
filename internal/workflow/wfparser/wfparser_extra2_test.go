package wfparser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkflowDefinition_ZeroValue(t *testing.T) {
	w := WorkflowDefinition{}
	if w.Name != "" || w.Description != "" || w.Author != "" {
		t.Error("zero value should have empty string fields")
	}
	if len(w.MCPDependencies) != 0 || len(w.InputParameters) != 0 {
		t.Error("zero value should have nil/empty slices")
	}
}

func TestValidate_WithDescriptionNoError(t *testing.T) {
	w := &WorkflowDefinition{Description: "do something"}
	errs := w.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestParseWorkflowFile_MultipleInputs(t *testing.T) {
	content := "---\ndescription: multi\ninputs:\n  - name: foo\n  - name: bar\n---\nbody\n"
	f := filepath.Join(t.TempDir(), "multi.prompt.md")
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	wf, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if wf.Description != "multi" {
		t.Errorf("expected multi, got %q", wf.Description)
	}
}

func TestParseWorkflowFile_NoDashes(t *testing.T) {
	content := "# Just a workflow\n\nNo frontmatter here.\n"
	f := filepath.Join(t.TempDir(), "nodash.md")
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	wf, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if wf.Description != "" {
		t.Errorf("expected empty description, got %q", wf.Description)
	}
}

func TestParseWorkflowFile_NameFromPromptMd(t *testing.T) {
	content := "---\ndescription: x\n---\n"
	f := filepath.Join(t.TempDir(), "myworkflow.prompt.md")
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	wf, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if wf.Name != "myworkflow" {
		t.Errorf("expected myworkflow, got %q", wf.Name)
	}
}

func TestParseWorkflowFile_ContentAfterFrontmatter(t *testing.T) {
	content := "---\ndescription: test\n---\nHello World\n"
	f := filepath.Join(t.TempDir(), "test.md")
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	wf, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if wf.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestValidate_ReturnsSlice(t *testing.T) {
	w := &WorkflowDefinition{}
	errs := w.Validate()
	if len(errs) == 0 {
		t.Error("expected at least one error for missing description")
	}
}
