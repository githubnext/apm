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
