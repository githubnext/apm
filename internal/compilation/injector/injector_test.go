package injector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestInject_SkippedWhenNoConstitution(t *testing.T) {
	dir := t.TempDir()
	ci := &ConstitutionInjector{BaseDir: dir}
	content := "# AGENTS.md\n\nhello"
	out, status, _ := ci.Inject(content, false, filepath.Join(dir, "AGENTS.md"))
	if status != StatusSkipped {
		t.Errorf("want SKIPPED, got %s", status)
	}
	if out != content {
		t.Errorf("content should be unchanged")
	}
}

func TestInject_MissingConstitutionFile(t *testing.T) {
	dir := t.TempDir()
	ci := &ConstitutionInjector{BaseDir: dir}
	content := "# AGENTS.md"
	out, status, _ := ci.Inject(content, true, filepath.Join(dir, "AGENTS.md"))
	if status != StatusMissing {
		t.Errorf("want MISSING, got %s", status)
	}
	if out != content {
		t.Errorf("content should be unchanged when missing")
	}
}

func TestInject_Created(t *testing.T) {
	dir := t.TempDir()
	// Create constitution file at expected path
	constitPath := filepath.Join(dir, compilationconst.ConstitutionRelativePath)
	if err := os.MkdirAll(filepath.Dir(constitPath), 0o755); err != nil {
		t.Fatal(err)
	}
	constitContent := "you must always do X"
	if err := os.WriteFile(constitPath, []byte(constitContent), 0o644); err != nil {
		t.Fatal(err)
	}

	ci := &ConstitutionInjector{BaseDir: dir}
	content := "# AGENTS.md\n\nhello"
	outputPath := filepath.Join(dir, "AGENTS.md")
	out, status, _ := ci.Inject(content, true, outputPath)
	if status != StatusCreated {
		t.Errorf("want CREATED, got %s", status)
	}
	if out == content {
		t.Error("output should differ from input after injection")
	}
	if out == "" {
		t.Error("output should not be empty")
	}
}

func TestInject_Updated(t *testing.T) {
	dir := t.TempDir()
	constitPath := filepath.Join(dir, compilationconst.ConstitutionRelativePath)
	if err := os.MkdirAll(filepath.Dir(constitPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(constitPath, []byte("new content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write an existing AGENTS.md with an old constitution block
	outputPath := filepath.Join(dir, "AGENTS.md")
	oldBlock := compilationconst.ConstitutionMarkerBegin + "\nold content\n" + compilationconst.ConstitutionMarkerEnd
	if err := os.WriteFile(outputPath, []byte(oldBlock+"\n\nhello"), 0o644); err != nil {
		t.Fatal(err)
	}

	ci := &ConstitutionInjector{BaseDir: dir}
	_, status, _ := ci.Inject("# fresh", true, outputPath)
	if status != StatusUpdated {
		t.Errorf("want UPDATED, got %s", status)
	}
}

func TestInject_Unchanged(t *testing.T) {
	dir := t.TempDir()
	constitPath := filepath.Join(dir, compilationconst.ConstitutionRelativePath)
	if err := os.MkdirAll(filepath.Dir(constitPath), 0o755); err != nil {
		t.Fatal(err)
	}
	constitContent := "same content"
	if err := os.WriteFile(constitPath, []byte(constitContent), 0o644); err != nil {
		t.Fatal(err)
	}

	block := compilationconst.ConstitutionMarkerBegin + "\n" + constitContent + "\n" + compilationconst.ConstitutionMarkerEnd
	outputPath := filepath.Join(dir, "AGENTS.md")
	if err := os.WriteFile(outputPath, []byte(block+"\n\nhello"), 0o644); err != nil {
		t.Fatal(err)
	}

	ci := &ConstitutionInjector{BaseDir: dir}
	_, status, _ := ci.Inject("# fresh", true, outputPath)
	if status != StatusUnchanged {
		t.Errorf("want UNCHANGED, got %s", status)
	}
}

func TestInject_PreservesExistingBlock(t *testing.T) {
	dir := t.TempDir()
	block := compilationconst.ConstitutionMarkerBegin + "\nexisting\n" + compilationconst.ConstitutionMarkerEnd
	outputPath := filepath.Join(dir, "AGENTS.md")
	if err := os.WriteFile(outputPath, []byte(block+"\n\ncontent"), 0o644); err != nil {
		t.Fatal(err)
	}

	ci := &ConstitutionInjector{BaseDir: dir}
	out, status, _ := ci.Inject("# new", false, outputPath)
	if status != StatusUnchanged {
		t.Errorf("want UNCHANGED, got %s", status)
	}
	if out == "# new" {
		t.Error("existing block should have been preserved in output")
	}
}
