package wfparser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidate_MissingDescription_Extra4(t *testing.T) {
	w := &WorkflowDefinition{}
	errs := w.Validate()
	if len(errs) == 0 {
		t.Error("expected error for missing description")
	}
}

func TestValidate_WithDescription_Extra4(t *testing.T) {
	w := &WorkflowDefinition{Description: "my desc"}
	errs := w.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestParseWorkflowFile_NotFound_Extra4(t *testing.T) {
	_, err := ParseWorkflowFile("/nonexistent/file.md")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseWorkflowFile_Empty_Extra4(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "workflow.md")
	_ = os.WriteFile(f, []byte(""), 0o644)
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Name != "workflow" {
		t.Errorf("expected workflow, got %s", w.Name)
	}
}

func TestParseWorkflowFile_WithFrontmatter_Extra4(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: hello\n---\nBody here"
	f := filepath.Join(dir, "my-wf.md")
	_ = os.WriteFile(f, []byte(content), 0o644)
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Description != "hello" {
		t.Errorf("expected hello, got %s", w.Description)
	}
	if !strings.Contains(w.Content, "Body") {
		t.Errorf("expected body content, got %s", w.Content)
	}
}

func TestParseWorkflowFile_PromptMdSuffix_Extra4(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "my-workflow.prompt.md")
	_ = os.WriteFile(f, []byte("---\ndescription: d\n---\n"), 0o644)
	w, _ := ParseWorkflowFile(f)
	if w.Name != "my-workflow" {
		t.Errorf("expected my-workflow, got %s", w.Name)
	}
}

func TestWorkflowDefinition_Fields_Extra4(t *testing.T) {
	w := &WorkflowDefinition{
		Name:        "n",
		FilePath:    "/f",
		Description: "d",
		Author:      "a",
		LLMModel:    "gpt4",
		Content:     "body",
	}
	if w.Name != "n" || w.FilePath != "/f" || w.Description != "d" {
		t.Error("unexpected field values")
	}
	if w.LLMModel != "gpt4" {
		t.Errorf("expected gpt4, got %s", w.LLMModel)
	}
}

func TestWorkflowDefinition_MCPDeps_Extra4(t *testing.T) {
	w := &WorkflowDefinition{MCPDependencies: []string{"a", "b"}}
	if len(w.MCPDependencies) != 2 {
		t.Errorf("expected 2 MCP deps, got %d", len(w.MCPDependencies))
	}
}

func TestWorkflowDefinition_InputParams_Extra4(t *testing.T) {
	w := &WorkflowDefinition{InputParameters: []string{"p1", "p2"}}
	if len(w.InputParameters) != 2 {
		t.Errorf("expected 2 input params, got %d", len(w.InputParameters))
	}
}

func TestParseWorkflowFile_NoFrontmatter_Extra4(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "simple.md")
	_ = os.WriteFile(f, []byte("Just content"), 0o644)
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Description != "" {
		t.Errorf("expected empty description, got %s", w.Description)
	}
}
