package wfparser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkflowDefinition_ZeroValue_Fields(t *testing.T) {
	var w WorkflowDefinition
	if w.Name != "" || w.FilePath != "" || w.Description != "" {
		t.Error("zero-value WorkflowDefinition should have empty strings")
	}
	if w.MCPDependencies != nil {
		t.Error("MCPDependencies should be nil for zero value")
	}
}

func TestValidate_DescriptionPresent_NoError(t *testing.T) {
	w := &WorkflowDefinition{Description: "some description"}
	errs := w.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_MissingDescription_ReturnsError(t *testing.T) {
	w := &WorkflowDefinition{}
	errs := w.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for missing description")
	}
}

func TestParseWorkflowFile_WithDescription(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.prompt.md")
	content := "---\ndescription: My workflow\n---\nDo something\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Description != "My workflow" {
		t.Errorf("expected 'My workflow', got %q", w.Description)
	}
}

func TestParseWorkflowFile_NameFromFilename(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "myflow.prompt.md")
	if err := os.WriteFile(f, []byte("---\ndescription: x\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Name != "myflow" {
		t.Errorf("expected 'myflow', got %q", w.Name)
	}
}

func TestParseWorkflowFile_MissingFile_ReturnsError(t *testing.T) {
	_, err := ParseWorkflowFile("/nonexistent/path/nope.prompt.md")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseWorkflowFile_NoFrontmatter_ContentSet(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "plain.prompt.md")
	if err := os.WriteFile(f, []byte("Just content\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestParseWorkflowFile_AuthorField(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "auth.prompt.md")
	content := "---\ndescription: x\nauthor: alice\n---\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.Author != "alice" {
		t.Errorf("expected 'alice', got %q", w.Author)
	}
}

func TestParseWorkflowFile_FilePathStored(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "fp.prompt.md")
	if err := os.WriteFile(f, []byte("---\ndescription: y\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := ParseWorkflowFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if w.FilePath != f {
		t.Errorf("FilePath mismatch: got %q, want %q", w.FilePath, f)
	}
}
