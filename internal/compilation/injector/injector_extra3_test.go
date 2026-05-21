package injector_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/injector"
)

func TestInjectionStatus_Constants_Extra3(t *testing.T) {
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
			t.Error("InjectionStatus constant should not be empty")
		}
		if seen[s] {
			t.Errorf("duplicate InjectionStatus value: %q", s)
		}
		seen[s] = true
	}
}

func TestConstitutionInjector_ZeroValue_Extra3(t *testing.T) {
	var ci injector.ConstitutionInjector
	if ci.BaseDir != "" {
		t.Errorf("BaseDir should be empty, got %q", ci.BaseDir)
	}
}

func TestInject_NoConstitution_NoOutputFile_Extra3(t *testing.T) {
	ci := injector.ConstitutionInjector{BaseDir: t.TempDir()}
	content, status, _ := ci.Inject("hello world", false, "/no/such/output.md")
	if content != "hello world" {
		t.Errorf("content should be unchanged, got %q", content)
	}
	if status != injector.StatusSkipped {
		t.Errorf("status = %q, want SKIPPED", status)
	}
}

func TestInject_WithConstitution_NewFile_Extra3(t *testing.T) {
	dir := t.TempDir()
	constitDir := filepath.Join(dir, ".specify", "memory")
	os.MkdirAll(constitDir, 0o755)
	os.WriteFile(filepath.Join(constitDir, "constitution.md"), []byte("# Constitution"), 0o644)

	ci := injector.ConstitutionInjector{BaseDir: dir}
	content, status, _ := ci.Inject("main content", true, filepath.Join(dir, "AGENTS.md"))
	if status != injector.StatusCreated {
		t.Errorf("status = %q, want CREATED", status)
	}
	if !strings.Contains(content, "# Constitution") {
		t.Error("content should contain constitution text")
	}
}

func TestInject_WithConstitution_Unchanged_Extra3(t *testing.T) {
	dir := t.TempDir()
	constitDir := filepath.Join(dir, ".specify", "memory")
	os.MkdirAll(constitDir, 0o755)
	os.WriteFile(filepath.Join(constitDir, "constitution.md"), []byte("# C"), 0o644)

	ci := injector.ConstitutionInjector{BaseDir: dir}
	outputPath := filepath.Join(dir, "AGENTS.md")

	// First inject to create the file
	content, _, _ := ci.Inject("main", true, outputPath)
	os.WriteFile(outputPath, []byte(content), 0o644)

	// Second inject should be unchanged
	_, status, _ := ci.Inject("main", true, outputPath)
	if status != injector.StatusUnchanged {
		t.Errorf("second inject status = %q, want UNCHANGED", status)
	}
}

func TestInject_BaseDirAssigned_Extra3(t *testing.T) {
	ci := injector.ConstitutionInjector{BaseDir: "/custom/dir"}
	if ci.BaseDir != "/custom/dir" {
		t.Errorf("BaseDir = %q, want /custom/dir", ci.BaseDir)
	}
}
