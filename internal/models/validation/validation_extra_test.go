package validation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/models/validation"
)

func TestPackageContentTypeString(t *testing.T) {
	cases := []struct {
		t    validation.PackageContentType
		want string
	}{
		{validation.PackageContentTypeInstructions, "instructions"},
		{validation.PackageContentTypeSkill, "skill"},
		{validation.PackageContentTypeHybrid, "hybrid"},
		{validation.PackageContentTypePrompts, "prompts"},
		{validation.PackageContentType(99), "hybrid"}, // unknown defaults to hybrid
	}
	for _, c := range cases {
		if got := c.t.String(); got != c.want {
			t.Errorf("PackageContentType(%d).String() = %q; want %q", c.t, got, c.want)
		}
	}
}

func TestValidationResultSummaryWithErrors(t *testing.T) {
	r := validation.NewValidationResult()
	r.AddError("first error")
	r.AddError("second error")
	r.AddWarning("a warning")
	s := r.Summary()
	if s == "" {
		t.Error("summary should not be empty")
	}
	if r.IsValid {
		t.Error("should be invalid with errors")
	}
}

func TestValidationResultNoIssues(t *testing.T) {
	r := validation.NewValidationResult()
	if r.HasIssues() {
		t.Error("fresh result should have no issues")
	}
}

func TestDetectPackageTypeAPMPackage(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: pkg\nversion: 1.0\n"), 0o644)
	os.MkdirAll(filepath.Join(dir, ".apm"), 0o755)
	pt, _ := validation.DetectPackageType(dir)
	if pt != validation.PackageTypeAPMPackage {
		t.Errorf("apm.yml + .apm/: got %v; want apm_package", pt)
	}
}

func TestDetectPackageTypeHybrid(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("name: pkg\nversion: 1.0\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: test\n---\n"), 0o644)
	pt, _ := validation.DetectPackageType(dir)
	if pt != validation.PackageTypeHybrid {
		t.Errorf("apm.yml+SKILL.md: got %v; want hybrid", pt)
	}
}

func TestDetectPackageTypeSkillBundle(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "myfeat")
	os.MkdirAll(skillDir, 0o755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\nname: feat\n---\n"), 0o644)
	pt, _ := validation.DetectPackageType(dir)
	if pt != validation.PackageTypeSkillBundle {
		t.Errorf("skill bundle: got %v; want skill_bundle", pt)
	}
}

func TestValidateAPMPackageValid(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".apm"), 0o755)
	apmYML := "name: my-pkg\nversion: 1.0.0\ndescription: A test package\n"
	os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(apmYML), 0o644)
	result := validation.ValidateAPMPackage(dir)
	if !result.IsValid {
		t.Errorf("valid package reported invalid: %v", result.Errors)
	}
}

func TestValidateAPMPackageMissingApmYML(t *testing.T) {
	dir := t.TempDir()
	result := validation.ValidateAPMPackage(dir)
	if result.IsValid {
		t.Error("expected invalid for missing apm.yml")
	}
}

func TestValidateAPMPackageMissingName(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "apm.yml"), []byte("version: 1.0.0\n"), 0o644)
	result := validation.ValidateAPMPackage(dir)
	if result.IsValid {
		t.Error("expected invalid when name missing")
	}
}

func TestPackageContentTypeFromStringCaseInsensitive(t *testing.T) {
	for _, s := range []string{"INSTRUCTIONS", "Instructions", "SKILL", "Skill", "PROMPTS"} {
		_, err := validation.PackageContentTypeFromString(s)
		if err != nil {
			t.Errorf("PackageContentTypeFromString(%q) error: %v", s, err)
		}
	}
}
