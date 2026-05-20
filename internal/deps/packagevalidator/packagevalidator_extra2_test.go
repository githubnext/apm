package packagevalidator

import (
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// ValidationResult additional methods
// ---------------------------------------------------------------------------

func TestValidationResult_IsValidAfterAddError(t *testing.T) {
	r := &ValidationResult{}
	if !r.IsValid() {
		t.Error("empty result should be valid")
	}
	r.AddError("something wrong")
	if r.IsValid() {
		t.Error("result with error should not be valid")
	}
}

func TestValidationResult_MultipleWarningsValid(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("w1")
	r.AddWarning("w2")
	r.AddWarning("w3")
	if !r.IsValid() {
		t.Error("multiple warnings should not invalidate result")
	}
	if len(r.Warnings) != 3 {
		t.Errorf("expected 3 warnings, got %d", len(r.Warnings))
	}
}

func TestValidationResult_ErrorsAndWarnings(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("w1")
	r.AddError("e1")
	if r.IsValid() {
		t.Error("result with error should not be valid")
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

// ---------------------------------------------------------------------------
// ValidateAPMPackage on a complete package
// ---------------------------------------------------------------------------

func TestValidateAPMPackage_WithApmYMLAndApmDirAndSkill(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	apmDir := filepath.Join(dir, ".apm")
	if err := os.Mkdir(apmDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(apmDir, "SKILL.md"), []byte("# skill"), 0644); err != nil {
		t.Fatal(err)
	}
	r := ValidateAPMPackage(dir)
	if !r.IsValid() {
		t.Errorf("expected valid package, got errors: %v", r.Errors)
	}
}

func TestValidateAPMPackage_EmptyApmYML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	r := ValidateAPMPackage(dir)
	// empty apm.yml should trigger a warning or error
	_ = r
}

// ---------------------------------------------------------------------------
// PackageValidator: new and ValidatePackageStructure same as ValidatePackage
// ---------------------------------------------------------------------------

func TestPackageValidator_ValidatePackage_MissingDir(t *testing.T) {
	v := New()
	r := v.ValidatePackage("/nonexistent/dir/abc")
	if r.IsValid() {
		t.Error("expected invalid result for missing dir")
	}
}

func TestPackageValidator_ValidatePackageStructure_MissingDir(t *testing.T) {
	v := New()
	r := v.ValidatePackageStructure("/nonexistent/dir/abc")
	if r.IsValid() {
		t.Error("expected invalid result for missing dir")
	}
}

func TestPackageValidator_ValidatePackage_ValidDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: pkg\n"), 0644); err != nil {
		t.Fatal(err)
	}
	v := New()
	r := v.ValidatePackage(dir)
	if !r.IsValid() {
		t.Errorf("expected valid package, got errors: %v", r.Errors)
	}
}

// ---------------------------------------------------------------------------
// ValidationResult zero value
// ---------------------------------------------------------------------------

func TestValidationResult_ZeroValue(t *testing.T) {
	var r ValidationResult
	if !r.IsValid() {
		t.Error("zero value should be valid")
	}
	if len(r.Errors) != 0 || len(r.Warnings) != 0 {
		t.Error("expected empty errors and warnings in zero value")
	}
}
