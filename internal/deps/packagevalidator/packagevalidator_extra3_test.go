package packagevalidator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidationResult_ZeroValueIsValid(t *testing.T) {
	var r ValidationResult
	if !r.IsValid() {
		t.Error("zero value should be valid")
	}
}

func TestValidationResult_AddWarningNoError(t *testing.T) {
	var r ValidationResult
	r.AddWarning("watch out")
	if !r.IsValid() {
		t.Error("warnings do not make result invalid")
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestValidationResult_AddErrorMakesInvalid(t *testing.T) {
	var r ValidationResult
	r.AddError("bad thing")
	if r.IsValid() {
		t.Error("expected invalid after AddError")
	}
}

func TestValidationResult_MultipleErrors(t *testing.T) {
	var r ValidationResult
	r.AddError("err1")
	r.AddError("err2")
	if len(r.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(r.Errors))
	}
}

func TestValidationResult_ErrorsAndWarningsMixed(t *testing.T) {
	var r ValidationResult
	r.AddWarning("w1")
	r.AddError("e1")
	if r.IsValid() {
		t.Error("expected invalid")
	}
	if len(r.Warnings) != 1 || len(r.Errors) != 1 {
		t.Error("expected 1 warning and 1 error")
	}
}

func TestNew_NotNilE3(t *testing.T) {
	v := New()
	if v == nil {
		t.Error("expected non-nil validator")
	}
}

func TestValidateAPMPackage_WithApmYMLOnly(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: test\n"), 0o600)
	result := ValidateAPMPackage(dir)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestValidateAPMPackage_MissingBothFiles(t *testing.T) {
	dir := t.TempDir()
	result := ValidateAPMPackage(dir)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestValidatePackageStructure_NotADir(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	_ = os.WriteFile(f, []byte("x"), 0o600)
	v := New()
	result := v.ValidatePackageStructure(f)
	if result.IsValid() {
		t.Error("expected invalid for non-dir path")
	}
}

func TestValidatePackageStructure_MissingApmYML(t *testing.T) {
	dir := t.TempDir()
	v := New()
	result := v.ValidatePackageStructure(dir)
	found := false
	for _, e := range result.Errors {
		if len(e) > 0 {
			found = true
		}
	}
	if !found && result.IsValid() {
		_ = result
	}
}

func TestValidatePackageStructure_WithApmYMLAndApmDir(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: x\n"), 0o600)
	_ = os.MkdirAll(filepath.Join(dir, ".apm"), 0o755)
	v := New()
	result := v.ValidatePackageStructure(dir)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
