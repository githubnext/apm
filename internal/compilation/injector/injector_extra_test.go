package injector_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
	"github.com/githubnext/apm/internal/compilation/injector"
)

func TestInject_StatusConstants(t *testing.T) {
	statuses := []injector.InjectionStatus{
		injector.StatusCreated,
		injector.StatusUpdated,
		injector.StatusUnchanged,
		injector.StatusSkipped,
		injector.StatusMissing,
	}
	seen := map[injector.InjectionStatus]bool{}
	for _, s := range statuses {
		if s == "" {
			t.Error("status constant must not be empty")
		}
		if seen[s] {
			t.Errorf("duplicate status: %q", s)
		}
		seen[s] = true
	}
}

func TestInject_MissingConstitutionFile_ReturnsMissing(t *testing.T) {
	ci := &injector.ConstitutionInjector{BaseDir: t.TempDir()}
	_, status, _ := ci.Inject("content", true, filepath.Join(t.TempDir(), "AGENTS.md"))
	if status != injector.StatusMissing {
		t.Errorf("status = %q, want %q", status, injector.StatusMissing)
	}
}

func TestInject_WithoutConstitution_NoExistingBlock_Skipped(t *testing.T) {
	ci := &injector.ConstitutionInjector{BaseDir: t.TempDir()}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	// No existing AGENTS.md, no block => StatusSkipped
	result, status, _ := ci.Inject("# My content\n", false, outputPath)
	if status != injector.StatusSkipped {
		t.Errorf("status = %q, want %q", status, injector.StatusSkipped)
	}
	if result != "# My content\n" {
		t.Errorf("result = %q, want same content", result)
	}
}

func TestInject_WithConstitution_CreatesBlock(t *testing.T) {
	baseDir := t.TempDir()
	constitutionDir := filepath.Join(baseDir, ".specify", "memory")
	if err := os.MkdirAll(constitutionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	constitPath := filepath.Join(baseDir, compilationconst.ConstitutionRelativePath)
	if err := os.WriteFile(constitPath, []byte("# Rules\n\nBe helpful.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ci := &injector.ConstitutionInjector{BaseDir: baseDir}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	result, status, _ := ci.Inject("# Original\n", true, outputPath)

	if status != injector.StatusCreated {
		t.Errorf("status = %q, want %q", status, injector.StatusCreated)
	}
	if !strings.Contains(result, compilationconst.ConstitutionMarkerBegin) {
		t.Error("result should contain constitution marker begin")
	}
	if !strings.Contains(result, "# Rules") {
		t.Error("result should contain constitution content")
	}
}

func TestInject_WithConstitution_UpdatesBlock(t *testing.T) {
	baseDir := t.TempDir()
	constitutionDir := filepath.Join(baseDir, ".specify", "memory")
	if err := os.MkdirAll(constitutionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	constitPath := filepath.Join(baseDir, compilationconst.ConstitutionRelativePath)
	if err := os.WriteFile(constitPath, []byte("# New rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create an existing AGENTS.md with an old block
	oldBlock := compilationconst.ConstitutionMarkerBegin + "\n# Old rules\n" + compilationconst.ConstitutionMarkerEnd
	ci := &injector.ConstitutionInjector{BaseDir: baseDir}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	if err := os.WriteFile(outputPath, []byte(oldBlock), 0o644); err != nil {
		t.Fatal(err)
	}

	result, status, _ := ci.Inject("# Content\n", true, outputPath)
	if status != injector.StatusUpdated {
		t.Errorf("status = %q, want %q", status, injector.StatusUpdated)
	}
	if strings.Contains(result, "# Old rules") {
		t.Error("old constitution content should be replaced")
	}
	if !strings.Contains(result, "# New rules") {
		t.Error("new constitution content should be present")
	}
}

func TestInject_WithConstitution_UnchangedWhenSame(t *testing.T) {
	baseDir := t.TempDir()
	constitutionDir := filepath.Join(baseDir, ".specify", "memory")
	if err := os.MkdirAll(constitutionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	constitPath := filepath.Join(baseDir, compilationconst.ConstitutionRelativePath)
	constitContent := "# Rules\n"
	if err := os.WriteFile(constitPath, []byte(constitContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// First inject to discover what block is written
	ci := &injector.ConstitutionInjector{BaseDir: baseDir}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	result1, _, _ := ci.Inject("# Content\n", true, outputPath)

	// Save the result as the existing AGENTS.md
	if err := os.WriteFile(outputPath, []byte(result1), 0o644); err != nil {
		t.Fatal(err)
	}

	// Second inject with same constitution should be Unchanged
	_, status, _ := ci.Inject("# Content\n", true, outputPath)
	if status != injector.StatusUnchanged {
		t.Errorf("status = %q, want %q", status, injector.StatusUnchanged)
	}
}

func TestInject_WithoutConstitution_PreservesExistingBlock(t *testing.T) {
	baseDir := t.TempDir()
	existingBlock := compilationconst.ConstitutionMarkerBegin + "\n# Saved rules\n" + compilationconst.ConstitutionMarkerEnd
	ci := &injector.ConstitutionInjector{BaseDir: baseDir}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	if err := os.WriteFile(outputPath, []byte(existingBlock), 0o644); err != nil {
		t.Fatal(err)
	}

	result, status, _ := ci.Inject("# Fresh content\n", false, outputPath)
	if status != injector.StatusUnchanged {
		t.Errorf("status = %q, want %q", status, injector.StatusUnchanged)
	}
	if !strings.Contains(result, "# Saved rules") {
		t.Error("existing constitution block should be preserved")
	}
}

func TestConstitutionInjector_BaseDirField(t *testing.T) {
	ci := &injector.ConstitutionInjector{BaseDir: "/some/path"}
	if ci.BaseDir != "/some/path" {
		t.Errorf("BaseDir = %q", ci.BaseDir)
	}
}
