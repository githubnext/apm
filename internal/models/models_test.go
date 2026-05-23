package models_test

import (
	"testing"

	"github.com/githubnext/apm/internal/models"
)

// TestParityPackageTypeString mirrors PackageType enum string values
func TestParityPackageTypeString(t *testing.T) {
	cases := []struct {
		pt   models.PackageType
		want string
	}{
		{models.PackageTypeAPMPackage, "apm_package"},
		{models.PackageTypeClaudeSkill, "claude_skill"},
		{models.PackageTypeHookPackage, "hook_package"},
		{models.PackageTypeHybrid, "hybrid"},
		{models.PackageTypeMarketplacePlugin, "marketplace_plugin"},
		{models.PackageTypeSkillBundle, "skill_bundle"},
		{models.PackageTypeInvalid, "invalid"},
	}
	for _, c := range cases {
		if got := c.pt.String(); got != c.want {
			t.Errorf("PackageType(%d).String() = %s, want %s", c.pt, got, c.want)
		}
	}
}

// TestParityPackageContentTypeString mirrors PackageContentType enum string values
func TestParityPackageContentTypeString(t *testing.T) {
	cases := []struct {
		ct   models.PackageContentType
		want string
	}{
		{models.PackageContentTypeInstructions, "instructions"},
		{models.PackageContentTypeSkill, "skill"},
		{models.PackageContentTypeHybrid, "hybrid"},
		{models.PackageContentTypePrompts, "prompts"},
	}
	for _, c := range cases {
		if got := c.ct.String(); got != c.want {
			t.Errorf("PackageContentType.String() = %s, want %s", got, c.want)
		}
	}
}

// TestParityParsePackageContentTypeValid mirrors PackageContentType.from_string for valid values
func TestParityParsePackageContentTypeValid(t *testing.T) {
	cases := []struct {
		input string
		want  models.PackageContentType
	}{
		{"instructions", models.PackageContentTypeInstructions},
		{"skill", models.PackageContentTypeSkill},
		{"hybrid", models.PackageContentTypeHybrid},
		{"prompts", models.PackageContentTypePrompts},
	}
	for _, c := range cases {
		got, err := models.ParsePackageContentType(c.input)
		if err != nil {
			t.Errorf("ParsePackageContentType(%q) unexpected error: %v", c.input, err)
		}
		if got != c.want {
			t.Errorf("ParsePackageContentType(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

// TestParityParsePackageContentTypeEmpty mirrors PackageContentType.from_string("") raises ValueError
func TestParityParsePackageContentTypeEmpty(t *testing.T) {
	_, err := models.ParsePackageContentType("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

// TestParityParsePackageContentTypeInvalid mirrors PackageContentType.from_string("invalid")
func TestParityParsePackageContentTypeInvalid(t *testing.T) {
	_, err := models.ParsePackageContentType("bad_type")
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

// TestParityInstallResultDefaults mirrors InstallResult default field values
func TestParityInstallResultDefaults(t *testing.T) {
	r := models.InstallResult{PackageTypes: make(map[string]string)}
	if r.InstalledCount != 0 {
		t.Errorf("InstalledCount default should be 0, got %d", r.InstalledCount)
	}
	if r.PromptsIntegrated != 0 {
		t.Errorf("PromptsIntegrated default should be 0")
	}
	if r.AgentsIntegrated != 0 {
		t.Errorf("AgentsIntegrated default should be 0")
	}
}

// TestParityPrimitiveCounts mirrors PrimitiveCounts default zero values
func TestParityPrimitiveCounts(t *testing.T) {
	pc := models.PrimitiveCounts{}
	if pc.Prompts != 0 || pc.Agents != 0 || pc.Instructions != 0 ||
		pc.Skills != 0 || pc.Hooks != 0 || pc.Commands != 0 {
		t.Error("PrimitiveCounts defaults should all be 0")
	}
}

// TestParityValidationErrorCodes mirrors ValidationError enum values
func TestParityValidationErrorCodes(t *testing.T) {
	if models.ValidationErrMissingAPMYml != "missing_apm_yml" {
		t.Error("ValidationErrMissingAPMYml mismatch")
	}
	if models.ValidationErrMissingAPMDir != "missing_apm_dir" {
		t.Error("ValidationErrMissingAPMDir mismatch")
	}
	if models.ValidationErrInvalidYMLFormat != "invalid_yml_format" {
		t.Error("ValidationErrInvalidYMLFormat mismatch")
	}
	if models.ValidationErrMissingRequired != "missing_required_field" {
		t.Error("ValidationErrMissingRequired mismatch")
	}
}
