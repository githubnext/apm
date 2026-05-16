package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverWorkflows_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	workflows, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	if len(workflows) != 0 {
		t.Errorf("expected no workflows, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_FindsPromptMd(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: test\n---\n# Test workflow"
	if err := os.WriteFile(filepath.Join(dir, "myflow.prompt.md"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	workflows, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	if len(workflows) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(workflows))
	}
	if workflows[0].Name != "myflow" {
		t.Errorf("expected name=myflow, got %q", workflows[0].Name)
	}
}

func TestDiscoverWorkflows_IgnoresNonPromptMd(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# readme"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "flow.md"), []byte("# flow"), 0600); err != nil {
		t.Fatal(err)
	}

	workflows, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	if len(workflows) != 0 {
		t.Errorf("expected no workflows, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_Nested(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\ndescription: nested\n---\n# Nested"
	if err := os.WriteFile(filepath.Join(sub, "nested.prompt.md"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	workflows, _ := DiscoverWorkflows(dir)
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow from nested dir, got %d", len(workflows))
	}
}
