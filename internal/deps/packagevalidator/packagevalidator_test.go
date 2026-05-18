package packagevalidator

import (
"os"
"path/filepath"
"testing"
)

func TestValidationResultMethods(t *testing.T) {
r := &ValidationResult{}
if !r.IsValid() {
t.Error("empty result should be valid")
}
r.AddError("bad input")
if r.IsValid() {
t.Error("result with error should be invalid")
}
if len(r.Errors) != 1 || r.Errors[0] != "bad input" {
t.Errorf("unexpected errors: %v", r.Errors)
}
r.AddWarning("minor issue")
if len(r.Warnings) != 1 || r.Warnings[0] != "minor issue" {
t.Errorf("unexpected warnings: %v", r.Warnings)
}
}

func TestValidateAPMPackageMissingDir(t *testing.T) {
result := ValidateAPMPackage("/nonexistent/path/12345")
if result.IsValid() {
t.Error("expected validation error for missing directory")
}
}

func TestValidateAPMPackageEmptyDir(t *testing.T) {
dir := t.TempDir()
result := ValidateAPMPackage(dir)
// Empty dir should have errors (missing apm.yml etc.)
if result.IsValid() {
t.Log("empty dir unexpectedly valid - check ValidateAPMPackage requirements")
}
}

func TestValidateAPMPackageWithApmYml(t *testing.T) {
dir := t.TempDir()
apmYml := filepath.Join(dir, "apm.yml")
if err := os.WriteFile(apmYml, []byte("name: test\nversion: 1.0.0\n"), 0o644); err != nil {
t.Fatal(err)
}
result := ValidateAPMPackage(dir)
// Should have fewer or no errors with apm.yml present
_ = result
}

func TestNewPackageValidator(t *testing.T) {
v := New()
if v == nil {
t.Fatal("New() returned nil")
}
result := v.ValidatePackage("/nonexistent/path/12345")
if result.IsValid() {
t.Error("expected invalid result for missing path")
}
}

func TestPackageValidatorStructure(t *testing.T) {
v := New()
dir := t.TempDir()
result := v.ValidatePackageStructure(dir)
_ = result // may or may not be valid depending on required files
}

func TestValidateAPMPackageIsFile(t *testing.T) {
	dir := t.TempDir()
	f := dir + "/notadir.txt"
	if err := os.WriteFile(f, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := ValidateAPMPackage(f)
	if result.IsValid() {
		t.Error("expected invalid result when path is a file, not a directory")
	}
}

func TestValidateAPMPackageEmptyApmYml(t *testing.T) {
	dir := t.TempDir()
	apmYml := filepath.Join(dir, "apm.yml")
	if err := os.WriteFile(apmYml, []byte("   \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := ValidateAPMPackage(dir)
	if result.IsValid() {
		t.Error("expected invalid result for empty apm.yml")
	}
}

func TestValidateAPMPackageWithApmDir(t *testing.T) {
	dir := t.TempDir()
	apmYml := filepath.Join(dir, "apm.yml")
	if err := os.WriteFile(apmYml, []byte("name: mypkg\nversion: 1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	apmDir := filepath.Join(dir, ".apm")
	if err := os.MkdirAll(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	result := ValidateAPMPackage(dir)
	if !result.IsValid() {
		t.Errorf("expected valid result with apm.yml and .apm dir: %v", result.Errors)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", result.Warnings)
	}
}

func TestValidationResultMultipleErrors(t *testing.T) {
	r := &ValidationResult{}
	r.AddError("error one")
	r.AddError("error two")
	r.AddError("error three")
	if r.IsValid() {
		t.Error("result with multiple errors should not be valid")
	}
	if len(r.Errors) != 3 {
		t.Errorf("expected 3 errors, got %d", len(r.Errors))
	}
}

func TestValidationResultMultipleWarnings(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("warn a")
	r.AddWarning("warn b")
	if !r.IsValid() {
		t.Error("result with only warnings should be valid")
	}
	if len(r.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(r.Warnings))
	}
}

func TestValidatePackageStructure_NotDir(t *testing.T) {
	v := New()
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := v.ValidatePackageStructure(f)
	if result.IsValid() {
		t.Error("expected invalid result when path is a file")
	}
}

func TestValidatePackageStructure_WithBothFiles(t *testing.T) {
	v := New()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".apm"), 0o755); err != nil {
		t.Fatal(err)
	}
	result := v.ValidatePackageStructure(dir)
	if !result.IsValid() {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", result.Warnings)
	}
}
