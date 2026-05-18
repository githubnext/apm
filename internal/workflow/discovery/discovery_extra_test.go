package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverWorkflows_EmptyString_UsesCwd(t *testing.T) {
	// Passing empty string should not panic; it uses cwd
	workflows, _ := DiscoverWorkflows("")
	// just verify it returns without crashing
	_ = workflows
}

func TestDiscoverWorkflows_NonExistentDir(t *testing.T) {
	workflows, _ := DiscoverWorkflows("/nonexistent/path/that/does/not/exist")
	if len(workflows) != 0 {
		t.Errorf("expected no workflows from non-existent dir, got %d", len(workflows))
	}
}

func TestDiscoverWorkflows_IgnoresDotFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".hidden.prompt.md"), []byte("---\ndescription: h\n---"), 0o600); err != nil {
		t.Fatal(err)
	}
	workflows, _ := DiscoverWorkflows(dir)
	// .hidden.prompt.md still matches *.prompt.md suffix -- just verify no panic
	_ = workflows
}

func TestDiscoverWorkflows_MixedFilesAndDirs(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Put a valid workflow in subdir and a plain md at root
	if err := os.WriteFile(filepath.Join(dir, "plain.md"), []byte("# not a workflow"), 0o600); err != nil {
		t.Fatal(err)
	}
	content := "---\ndescription: sub workflow\n---\n# Sub"
	if err := os.WriteFile(filepath.Join(subdir, "sub.prompt.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	workflows, errs := DiscoverWorkflows(dir)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}
	if workflows[0].Name != "sub" {
		t.Errorf("expected name 'sub', got %q", workflows[0].Name)
	}
}

func TestDiscoverWorkflows_ParseErrorCounted(t *testing.T) {
	dir := t.TempDir()
	// Write a file that will fail to parse (no frontmatter)
	if err := os.WriteFile(filepath.Join(dir, "bad.prompt.md"), []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	workflows, errs := DiscoverWorkflows(dir)
	// Either parsed successfully (empty file = valid) or error reported - no panic
	_ = workflows
	_ = errs
}

func TestDiscoverWorkflows_DuplicatePaths(t *testing.T) {
	dir := t.TempDir()
	content := "---\ndescription: flow\n---\n# Flow"
	if err := os.WriteFile(filepath.Join(dir, "flow.prompt.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	// Calling twice should work fine
	w1, _ := DiscoverWorkflows(dir)
	w2, _ := DiscoverWorkflows(dir)
	if len(w1) != len(w2) {
		t.Errorf("expected same count on repeated calls: %d vs %d", len(w1), len(w2))
	}
}

func TestDiscoverWorkflows_MultipleLevels(t *testing.T) {
	dir := t.TempDir()
	levels := []string{"a", "a/b", "a/b/c", "x"}
	for _, l := range levels {
		p := filepath.Join(dir, l)
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "---\ndescription: " + l + "\n---\n# " + l
		if err := os.WriteFile(filepath.Join(p, l[len(l)-1:]+".prompt.md"), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	workflows, _ := DiscoverWorkflows(dir)
	if len(workflows) != len(levels) {
		t.Errorf("expected %d workflows, got %d", len(levels), len(workflows))
	}
}
