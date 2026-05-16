package wfparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/workflow/wfparser"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return path
}

func TestParseWorkflowFile_WithFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: My workflow\nauthor: alice\nllm: gpt-4\n---\n# Body here\n"
	path := writeFile(t, dir, "my-workflow.md", content)
	w, err := wfparser.ParseWorkflowFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Description != "My workflow" {
		t.Errorf("Description: got %q", w.Description)
	}
	if w.Author != "alice" {
		t.Errorf("Author: got %q", w.Author)
	}
	if w.LLMModel != "gpt-4" {
		t.Errorf("LLMModel: got %q", w.LLMModel)
	}
	if w.Name != "my-workflow" {
		t.Errorf("Name: got %q", w.Name)
	}
}

func TestParseWorkflowFile_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "# Just content\nno frontmatter\n"
	path := writeFile(t, dir, "bare.md", content)
	w, err := wfparser.ParseWorkflowFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Description != "" {
		t.Errorf("Description should be empty, got %q", w.Description)
	}
	if w.Content != content {
		t.Errorf("Content mismatch")
	}
}

func TestParseWorkflowFile_MCPAndInput(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: Test\nmcp:\n  - server1\n  - server2\ninput:\n  - param_a\n  - param_b\n---\nBody\n"
	path := writeFile(t, dir, "tools.md", content)
	w, err := wfparser.ParseWorkflowFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.MCPDependencies) != 2 {
		t.Errorf("MCPDependencies: got %d, want 2", len(w.MCPDependencies))
	}
	if w.MCPDependencies[0] != "server1" {
		t.Errorf("MCPDependencies[0]: got %q", w.MCPDependencies[0])
	}
	if len(w.InputParameters) != 2 {
		t.Errorf("InputParameters: got %d, want 2", len(w.InputParameters))
	}
	if w.InputParameters[0] != "param_a" {
		t.Errorf("InputParameters[0]: got %q", w.InputParameters[0])
	}
}

func TestParseWorkflowFile_PromptMdExtension(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: Prompt workflow\n---\nContent\n"
	path := writeFile(t, dir, "myflow.prompt.md", content)
	w, err := wfparser.ParseWorkflowFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Name != "myflow" {
		t.Errorf("Name: got %q, want %q", w.Name, "myflow")
	}
}

func TestValidate_MissingDescription(t *testing.T) {
	w := &wfparser.WorkflowDefinition{}
	errs := w.Validate()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0] != "Missing 'description' in frontmatter" {
		t.Errorf("error message: %q", errs[0])
	}
}

func TestValidate_WithDescription(t *testing.T) {
	w := &wfparser.WorkflowDefinition{Description: "has a description"}
	errs := w.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestParseWorkflowFile_NotFound(t *testing.T) {
	_, err := wfparser.ParseWorkflowFile("/nonexistent/path/file.md")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
