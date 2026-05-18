package wfparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/workflow/wfparser"
)

func TestParseWorkflowFile_AuthorField(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: My flow\nauthor: Alice\n---\n# Body"
	f := filepath.Join(dir, "myflow.md")
	os.WriteFile(f, []byte(content), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if w.Author != "Alice" {
		t.Errorf("Author = %q, want Alice", w.Author)
	}
}

func TestParseWorkflowFile_LLMModel(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: flow\nllm: gpt-4o\n---\n"
	f := filepath.Join(dir, "flow.md")
	os.WriteFile(f, []byte(content), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if w.LLMModel != "gpt-4o" {
		t.Errorf("LLMModel = %q, want gpt-4o", w.LLMModel)
	}
}

func TestParseWorkflowFile_EmptyFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "---\n---\n# Just body"
	f := filepath.Join(dir, "empty.md")
	os.WriteFile(f, []byte(content), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if w.Description != "" {
		t.Errorf("expected empty description, got %q", w.Description)
	}
}

func TestParseWorkflowFile_BodyPreserved(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: test\n---\nHello body content"
	f := filepath.Join(dir, "wf.md")
	os.WriteFile(f, []byte(content), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if w.Content == "" {
		t.Error("body content should be non-empty")
	}
}

func TestParseWorkflowFile_Name(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "my-workflow.prompt.md")
	os.WriteFile(f, []byte("---\ndescription: x\n---\n"), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if w.Name != "my-workflow" {
		t.Errorf("Name = %q, want my-workflow", w.Name)
	}
}

func TestParseWorkflowFile_MCPList(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: flow\nmcp:\n  - tool-a\n  - tool-b\n---\n"
	f := filepath.Join(dir, "mcp-flow.md")
	os.WriteFile(f, []byte(content), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(w.MCPDependencies) != 2 {
		t.Errorf("MCPDependencies = %v, want 2", w.MCPDependencies)
	}
}

func TestParseWorkflowFile_InputList(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: flow\ninput:\n  - param1\n  - param2\n  - param3\n---\n"
	f := filepath.Join(dir, "input-flow.md")
	os.WriteFile(f, []byte(content), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(w.InputParameters) != 3 {
		t.Errorf("InputParameters = %v, want 3", w.InputParameters)
	}
}

func TestValidate_AllGood(t *testing.T) {
	w := &wfparser.WorkflowDefinition{
		Name:        "good",
		Description: "A valid workflow",
	}
	errs := w.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_EmptyDescription(t *testing.T) {
	w := &wfparser.WorkflowDefinition{Name: "noDesc"}
	errs := w.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for missing description")
	}
}

func TestParseWorkflowFile_FilePathSet(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "wf.md")
	os.WriteFile(f, []byte("---\ndescription: test\n---\n"), 0644)
	w, err := wfparser.ParseWorkflowFile(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if w.FilePath != f {
		t.Errorf("FilePath = %q, want %q", w.FilePath, f)
	}
}
