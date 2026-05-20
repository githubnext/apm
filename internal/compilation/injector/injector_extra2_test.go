package injector_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
	"github.com/githubnext/apm/internal/compilation/injector"
)

func TestInject_StatusValues_Stable(t *testing.T) {
	if injector.StatusCreated != "CREATED" {
		t.Errorf("StatusCreated = %q, want CREATED", injector.StatusCreated)
	}
	if injector.StatusUpdated != "UPDATED" {
		t.Errorf("StatusUpdated = %q, want UPDATED", injector.StatusUpdated)
	}
	if injector.StatusUnchanged != "UNCHANGED" {
		t.Errorf("StatusUnchanged = %q, want UNCHANGED", injector.StatusUnchanged)
	}
	if injector.StatusSkipped != "SKIPPED" {
		t.Errorf("StatusSkipped = %q, want SKIPPED", injector.StatusSkipped)
	}
	if injector.StatusMissing != "MISSING" {
		t.Errorf("StatusMissing = %q, want MISSING", injector.StatusMissing)
	}
}

func TestInject_WithConstitution_UpdatesBlock_variants(t *testing.T) {
	baseDir := t.TempDir()
	constitDir := filepath.Join(baseDir, ".specify", "memory")
	if err := os.MkdirAll(constitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	constitPath := filepath.Join(baseDir, compilationconst.ConstitutionRelativePath)
	if err := os.WriteFile(constitPath, []byte("new constitution"), 0o644); err != nil {
		t.Fatal(err)
	}

	existingBlock := compilationconst.ConstitutionMarkerBegin + "\nold\n" + compilationconst.ConstitutionMarkerEnd
	content := "# Preamble\n" + existingBlock + "\n# Epilogue\n"

	ci := &injector.ConstitutionInjector{BaseDir: baseDir}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	result, status, _ := ci.Inject(content, true, outputPath)

	if status == injector.StatusMissing {
		t.Fatalf("expected non-missing status, got MISSING")
	}
	if !strings.Contains(result, "new constitution") {
		t.Errorf("result should contain new constitution content, got: %q", result)
	}
}

func TestInject_NoConstitutionDir_MissingStatus(t *testing.T) {
	ci := &injector.ConstitutionInjector{BaseDir: t.TempDir()}
	_, status, _ := ci.Inject("any content", true, filepath.Join(t.TempDir(), "out.md"))
	if status != injector.StatusMissing {
		t.Errorf("expected StatusMissing when constitution dir absent, got %q", status)
	}
}

func TestInject_WithConstitution_EmptyContent(t *testing.T) {
	baseDir := t.TempDir()
	constitDir := filepath.Join(baseDir, ".specify", "memory")
	if err := os.MkdirAll(constitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	constitPath := filepath.Join(baseDir, compilationconst.ConstitutionRelativePath)
	if err := os.WriteFile(constitPath, []byte("mission statement"), 0o644); err != nil {
		t.Fatal(err)
	}
	ci := &injector.ConstitutionInjector{BaseDir: baseDir}
	outputPath := filepath.Join(t.TempDir(), "AGENTS.md")
	result, status, _ := ci.Inject("", true, outputPath)
	// Empty content with no existing block => either created or skipped
	if status == injector.StatusMissing {
		t.Errorf("should not get MISSING when constitution file exists")
	}
	_ = result
}

func TestConstitutionInjector_BaseDir_Stored(t *testing.T) {
	ci := &injector.ConstitutionInjector{BaseDir: "/some/path"}
	if ci.BaseDir != "/some/path" {
		t.Errorf("BaseDir = %q, want /some/path", ci.BaseDir)
	}
}
