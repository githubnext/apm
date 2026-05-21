package packagevalidator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidationResult_NewHasNoErrors(t *testing.T) {
	r := &ValidationResult{}
	if !r.IsValid() {
		t.Error("expected new result to be valid")
	}
}

func TestValidationResult_AddErrorMakesItInvalid(t *testing.T) {
	r := &ValidationResult{}
	r.AddError("something broke")
	if r.IsValid() {
		t.Error("expected invalid after adding error")
	}
}

func TestValidationResult_AddWarningRemainsValid(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("heads up")
	if !r.IsValid() {
		t.Error("warnings should not invalidate")
	}
}

func TestValidationResult_ErrorCount(t *testing.T) {
	r := &ValidationResult{}
	r.AddError("e1")
	r.AddError("e2")
	if len(r.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(r.Errors))
	}
}

func TestValidationResult_WarningCount(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("w1")
	r.AddWarning("w2")
	r.AddWarning("w3")
	if len(r.Warnings) != 3 {
		t.Errorf("expected 3 warnings, got %d", len(r.Warnings))
	}
}

func TestNew_ReturnsValidator(t *testing.T) {
	v := New()
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
}

func TestValidateAPMPackage_TempDirWithYml(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: test\n"), 0o644)
	r := ValidateAPMPackage(dir)
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	if !r.IsValid() {
		t.Errorf("expected valid for dir with apm.yml: %v", r.Errors)
	}
}

func TestValidateAPMPackage_TempDirNoFiles(t *testing.T) {
	dir := t.TempDir()
	r := ValidateAPMPackage(dir)
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	if r.IsValid() {
		t.Error("expected invalid for empty dir")
	}
}

func TestValidateAPMPackage_NonExistentDir(t *testing.T) {
	r := ValidateAPMPackage("/nonexistent/path/xyz999")
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	if r.IsValid() {
		t.Error("expected invalid for nonexistent path")
	}
}

func TestValidationResult_MixedErrorsAndWarnings(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("w1")
	r.AddError("e1")
	if r.IsValid() {
		t.Error("expected invalid when errors present alongside warnings")
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}
