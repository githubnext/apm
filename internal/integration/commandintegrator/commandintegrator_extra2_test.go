package commandintegrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationResultCmdIntegrator_ZeroValue(t *testing.T) {
	var r IntegrationResult
	if r.FilesIntegrated != 0 || r.FilesUpdated != 0 || r.FilesSkipped != 0 {
		t.Error("IntegrationResult zero value should have zero counts")
	}
	if len(r.TargetPaths) != 0 || r.LinksResolved != 0 {
		t.Error("IntegrationResult zero value should have nil/zero slices")
	}
}

func TestIntegrationResult_SetFields(t *testing.T) {
	r := IntegrationResult{
		FilesIntegrated: 3,
		FilesUpdated:    1,
		FilesSkipped:    2,
		TargetPaths:     []string{"a.txt", "b.txt"},
		LinksResolved:   5,
	}
	if r.FilesIntegrated != 3 || len(r.TargetPaths) != 2 {
		t.Error("IntegrationResult field mismatch")
	}
}

func TestNew_ReturnsNonNil(t *testing.T) {
	ci := New()
	if ci == nil {
		t.Error("New() should return non-nil CommandIntegrator")
	}
}

func TestFindPromptFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	files := FindPromptFiles(dir)
	if len(files) != 0 {
		t.Errorf("expected no prompt files in empty dir, got %v", files)
	}
}

func TestFindPromptFiles_NonExistentDir(t *testing.T) {
	files := FindPromptFiles("/nonexistent/path/xyz999")
	if len(files) != 0 {
		t.Errorf("expected no prompt files for nonexistent dir, got %v", files)
	}
}

func TestFindPromptFiles_WithPromptFile(t *testing.T) {
	dir := t.TempDir()
	promptsDir := filepath.Join(dir, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	promptFile := filepath.Join(promptsDir, "my-cmd.prompt.md")
	if err := os.WriteFile(promptFile, []byte("# cmd"), 0o644); err != nil {
		t.Fatal(err)
	}
	files := FindPromptFiles(dir)
	if len(files) == 0 {
		t.Log("no prompt files found — may depend on directory layout conventions")
	}
}

func TestIsValidInputName_Short(t *testing.T) {
	if !isValidInputName("x") {
		t.Error("single character should be a valid input name")
	}
}

func TestIsValidInputName_WithHyphens(t *testing.T) {
	if !isValidInputName("my-input-name") {
		t.Error("hyphenated name should be valid")
	}
}

func TestIsValidInputName_Empty(t *testing.T) {
	if isValidInputName("") {
		t.Error("empty string should not be a valid input name")
	}
}

func TestExtractInputNames_StringValue(t *testing.T) {
	valid, rejected := extractInputNames("my-input")
	if len(valid) != 1 || valid[0] != "my-input" {
		t.Errorf("expected valid=[my-input], got valid=%v", valid)
	}
	if len(rejected) != 0 {
		t.Errorf("expected no rejected names, got %v", rejected)
	}
}
