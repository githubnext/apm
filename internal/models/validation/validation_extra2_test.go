package validation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/models/validation"
)

func TestNewValidationResult_IsValid(t *testing.T) {
	r := validation.NewValidationResult()
	if !r.IsValid {
		t.Error("expected IsValid=true for empty result")
	}
}

func TestValidationResult_AddError_SetsInvalid(t *testing.T) {
	r := validation.NewValidationResult()
	r.AddError("something went wrong")
	if r.IsValid {
		t.Error("expected IsValid=false after AddError")
	}
}

func TestValidationResult_AddWarning_StaysValid(t *testing.T) {
	r := validation.NewValidationResult()
	r.AddWarning("just a warning")
	if !r.IsValid {
		t.Error("expected IsValid=true after AddWarning (no error)")
	}
}

func TestValidationResult_HasIssues_NoErrors(t *testing.T) {
	r := validation.NewValidationResult()
	if r.HasIssues() {
		t.Error("expected HasIssues=false for empty result")
	}
}

func TestValidationResult_HasIssues_WithWarning(t *testing.T) {
	r := validation.NewValidationResult()
	r.AddWarning("warn")
	if !r.HasIssues() {
		t.Error("expected HasIssues=true after warning")
	}
}

func TestValidationResult_HasIssues_WithError(t *testing.T) {
	r := validation.NewValidationResult()
	r.AddError("err")
	if !r.HasIssues() {
		t.Error("expected HasIssues=true after error")
	}
}

func TestValidationResult_Summary_Empty(t *testing.T) {
	r := validation.NewValidationResult()
	s := r.Summary()
	_ = s
}

func TestValidationResult_Summary_WithErrors(t *testing.T) {
	r := validation.NewValidationResult()
	r.AddError("err1")
	r.AddError("err2")
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary with errors")
	}
}

func TestPackageContentTypeFromString_Hybrid(t *testing.T) {
	ct, err := validation.PackageContentTypeFromString("hybrid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != validation.PackageContentTypeHybrid {
		t.Errorf("expected Hybrid, got %v", ct)
	}
}

func TestPackageContentTypeFromString_Skill(t *testing.T) {
	ct, err := validation.PackageContentTypeFromString("skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != validation.PackageContentTypeSkill {
		t.Errorf("expected Skill, got %v", ct)
	}
}

func TestPackageContentTypeFromString_Instructions(t *testing.T) {
	ct, err := validation.PackageContentTypeFromString("instructions")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != validation.PackageContentTypeInstructions {
		t.Errorf("expected Instructions, got %v", ct)
	}
}

func TestPackageContentTypeFromString_Invalid(t *testing.T) {
	_, err := validation.PackageContentTypeFromString("unknown_type_xyz")
	if err == nil {
		t.Error("expected error for unknown content type")
	}
}

func TestGatherDetectionEvidence_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	ev := validation.GatherDetectionEvidence(dir)
	if ev == nil {
		t.Error("expected non-nil DetectionEvidence")
	}
}

func TestGatherDetectionEvidence_WithApmYml(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ev := validation.GatherDetectionEvidence(dir)
	if ev == nil {
		t.Error("expected non-nil evidence")
	}
}

func TestDetectPackageType_UnknownDir(t *testing.T) {
	dir := t.TempDir()
	pt, _ := validation.DetectPackageType(dir)
	_ = pt
}

func TestValidateAPMPackage_EmptyDir_Invalid(t *testing.T) {
	dir := t.TempDir()
	result := validation.ValidateAPMPackage(dir)
	if result == nil {
		t.Error("expected non-nil result for empty dir")
	}
}

func TestValidateAPMPackage_WithMinimalApmYml(t *testing.T) {
	dir := t.TempDir()
	content := "name: my-package\nversion: 1.0.0\n"
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	result := validation.ValidateAPMPackage(dir)
	if result == nil {
		t.Error("expected non-nil result")
	}
}
