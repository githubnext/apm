// Package packagevalidator validates APM package structure and content.
//
// Corresponds to src/apm_cli/deps/package_validator.py.
package packagevalidator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationResult holds the outcome of a package validation run.
type ValidationResult struct {
	Errors   []string
	Warnings []string
}

// IsValid returns true if the validation result has no errors.
func (r *ValidationResult) IsValid() bool { return len(r.Errors) == 0 }

// AddError appends an error message.
func (r *ValidationResult) AddError(msg string) { r.Errors = append(r.Errors, msg) }

// AddWarning appends a warning message.
func (r *ValidationResult) AddWarning(msg string) { r.Warnings = append(r.Warnings, msg) }

// PackageValidator validates APM package structure.
type PackageValidator struct{}

// New creates a new PackageValidator.
func New() *PackageValidator { return &PackageValidator{} }

// ValidatePackage validates that packagePath contains a valid APM package.
func (v *PackageValidator) ValidatePackage(packagePath string) *ValidationResult {
	return ValidateAPMPackage(packagePath)
}

// ValidatePackageStructure checks for required files and directories.
func (v *PackageValidator) ValidatePackageStructure(packagePath string) *ValidationResult {
	result := &ValidationResult{}

	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		result.AddError(fmt.Sprintf("Package directory does not exist: %s", packagePath))
		return result
	}
	fi, err := os.Stat(packagePath)
	if err != nil {
		result.AddError(fmt.Sprintf("Cannot stat package path: %s", packagePath))
		return result
	}
	if !fi.IsDir() {
		result.AddError(fmt.Sprintf("Package path is not a directory: %s", packagePath))
		return result
	}

	// Required: apm.yml at root
	apmYML := filepath.Join(packagePath, "apm.yml")
	if _, err := os.Stat(apmYML); os.IsNotExist(err) {
		result.AddError(fmt.Sprintf("Missing required file: apm.yml (looked in %s)", packagePath))
	}

	// Required: .apm/ directory
	apmDir := filepath.Join(packagePath, ".apm")
	if fi, err := os.Stat(apmDir); os.IsNotExist(err) || !fi.IsDir() {
		result.AddWarning(fmt.Sprintf("Missing .apm/ directory in %s", packagePath))
	}

	return result
}

// ValidateAPMPackage runs a full validation on packagePath.
func ValidateAPMPackage(packagePath string) *ValidationResult {
	v := &PackageValidator{}
	result := v.ValidatePackageStructure(packagePath)
	if !result.IsValid() {
		return result
	}

	// Additional content checks
	apmYML := filepath.Join(packagePath, "apm.yml")
	raw, err := os.ReadFile(apmYML)
	if err != nil {
		result.AddError(fmt.Sprintf("Cannot read apm.yml: %s", err))
		return result
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		result.AddError("apm.yml is empty")
		return result
	}

	return result
}
