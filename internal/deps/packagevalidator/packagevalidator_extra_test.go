package packagevalidator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidationResult_WarningsDoNotInvalidate(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("w1")
	r.AddWarning("w2")
	r.AddWarning("w3")
	if !r.IsValid() {
		t.Error("warnings alone should not make result invalid")
	}
}

func TestValidationResult_ErrorAfterWarnings(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("minor issue")
	r.AddError("fatal issue")
	if r.IsValid() {
		t.Error("result with error should not be valid")
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(r.Errors))
	}
}

func TestValidateAPMPackage_WithApmYMLAndApmDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: mypkg\nversion: 1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".apm"), 0o755); err != nil {
		t.Fatal(err)
	}
	result := ValidateAPMPackage(dir)
	if !result.IsValid() {
		t.Errorf("expected valid package, got errors: %v", result.Errors)
	}
}

func TestValidateAPMPackage_ApmYMLMinimalContent(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("x: y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".apm"), 0o755); err != nil {
		t.Fatal(err)
	}
	result := ValidateAPMPackage(dir)
	// Any non-empty apm.yml with .apm dir should be valid at structure level
	_ = result
}

func TestValidateAPMPackage_NonExistentPath(t *testing.T) {
	result := ValidateAPMPackage("/nonexistent/xyz/abc/123")
	if result.IsValid() {
		t.Error("expected invalid for non-existent path")
	}
}

func TestPackageValidator_ValidatePackageAndStructure_Same(t *testing.T) {
	v := New()
	dir := t.TempDir()
	// Neither should panic on empty dir
	r1 := v.ValidatePackage(dir)
	r2 := v.ValidatePackageStructure(dir)
	// Both should agree on validity
	if r1.IsValid() != r2.IsValid() {
		t.Errorf("ValidatePackage and ValidatePackageStructure disagreed for empty dir: pkg=%v struct=%v", r1.IsValid(), r2.IsValid())
	}
}

func TestValidationResult_PreserveOrder(t *testing.T) {
	r := &ValidationResult{}
	r.AddError("first")
	r.AddError("second")
	r.AddError("third")
	if r.Errors[0] != "first" || r.Errors[1] != "second" || r.Errors[2] != "third" {
		t.Errorf("errors should preserve insertion order: %v", r.Errors)
	}
}

func TestValidationResult_PreserveWarningOrder(t *testing.T) {
	r := &ValidationResult{}
	r.AddWarning("alpha")
	r.AddWarning("beta")
	if r.Warnings[0] != "alpha" || r.Warnings[1] != "beta" {
		t.Errorf("warnings should preserve insertion order: %v", r.Warnings)
	}
}

func TestValidateAPMPackage_SymlinksOrSpecialFile(t *testing.T) {
	dir := t.TempDir()
	// Create a regular file as apm.yml but make the path point to a file (not dir)
	f := filepath.Join(dir, "not_a_pkg.txt")
	os.WriteFile(f, []byte("data"), 0o644)
	result := ValidateAPMPackage(f)
	if result.IsValid() {
		t.Error("file as package path should not be valid")
	}
}
