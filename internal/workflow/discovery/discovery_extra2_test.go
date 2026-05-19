package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverWorkflows_SingleFileReturnsOne(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "test.prompt.md"), []byte("---\nname: test\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	workflows, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_NoDotMd(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notworkflow.md"), []byte("# heading"), 0o644); err != nil {
		t.Fatal(err)
	}
	workflows, errs := DiscoverWorkflows(dir)
	_ = errs
	if len(workflows) != 0 {
		t.Errorf("expected 0 workflows for plain .md file, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_TxtFileIgnored(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "test.prompt.txt"), []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}
	workflows, _ := DiscoverWorkflows(dir)
	if len(workflows) != 0 {
		t.Errorf("expected .prompt.txt to be ignored, got %d workflows", len(workflows))
	}
}

func TestDiscoverWorkflows_NoFiles(t *testing.T) {
	dir := t.TempDir()
	workflows, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Errorf("expected no errors for empty dir, got %v", errs)
	}
	if len(workflows) != 0 {
		t.Errorf("expected 0 workflows in empty dir, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_AbsoluteDirReturnsResults(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "flow.prompt.md")
	if err := os.WriteFile(name, []byte("---\nname: flow\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	workflows, _ := DiscoverWorkflows(dir)
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow from absolute dir, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_SubdirFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "deep.prompt.md"), []byte("---\nname: deep\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	workflows, _ := DiscoverWorkflows(dir)
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow in subdir, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_MultipleInSameDir(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.prompt.md", "b.prompt.md", "c.prompt.md"} {
		if err := os.WriteFile(filepath.Join(dir, name),
			[]byte("---\nname: "+name+"\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	workflows, _ := DiscoverWorkflows(dir)
	if len(workflows) != 3 {
		t.Errorf("expected 3 workflows, got %d", len(workflows))
	}
}
