package validation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/models/validation"
)

func TestPackageTypeString(t *testing.T) {
	cases := []struct {
		t    validation.PackageType
		want string
	}{
		{validation.PackageTypeAPMPackage, "apm_package"},
		{validation.PackageTypeClaudeSkill, "claude_skill"},
		{validation.PackageTypeHookPackage, "hook_package"},
		{validation.PackageTypeHybrid, "hybrid"},
		{validation.PackageTypeMarketplacePlugin, "marketplace_plugin"},
		{validation.PackageTypeSkillBundle, "skill_bundle"},
		{validation.PackageTypeInvalid, "invalid"},
	}
	for _, c := range cases {
		if got := c.t.String(); got != c.want {
			t.Errorf("PackageType(%d).String() = %q; want %q", c.t, got, c.want)
		}
	}
}

func TestPackageContentTypeFromString(t *testing.T) {
	cases := []struct {
		input   string
		want    validation.PackageContentType
		wantErr bool
	}{
		{"instructions", validation.PackageContentTypeInstructions, false},
		{"skill", validation.PackageContentTypeSkill, false},
		{"hybrid", validation.PackageContentTypeHybrid, false},
		{"prompts", validation.PackageContentTypePrompts, false},
		{"HYBRID", validation.PackageContentTypeHybrid, false},
		{"", 0, true},
		{"unknown", 0, true},
	}
	for _, c := range cases {
		got, err := validation.PackageContentTypeFromString(c.input)
		if c.wantErr {
			if err == nil {
				t.Errorf("PackageContentTypeFromString(%q) expected error", c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("PackageContentTypeFromString(%q) unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("PackageContentTypeFromString(%q) = %v; want %v", c.input, got, c.want)
		}
	}
}

func TestValidationResult(t *testing.T) {
	r := validation.NewValidationResult()
	if !r.IsValid {
		t.Error("new result should be valid")
	}
	r.AddWarning("test warning")
	if !r.IsValid {
		t.Error("warning should not make invalid")
	}
	if !r.HasIssues() {
		t.Error("has issues after warning")
	}
	r.AddError("test error")
	if r.IsValid {
		t.Error("should be invalid after error")
	}
	summary := r.Summary()
	if summary == "" {
		t.Error("summary should not be empty")
	}
}

func TestDetectPackageTypeInvalid(t *testing.T) {
	dir := t.TempDir()
	pt, _ := validation.DetectPackageType(dir)
	if pt != validation.PackageTypeInvalid {
		t.Errorf("empty dir: got %v; want invalid", pt)
	}
}

func TestDetectPackageTypeClaudeSkill(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: test\n---\n# Test"), 0o644)
	pt, _ := validation.DetectPackageType(dir)
	if pt != validation.PackageTypeClaudeSkill {
		t.Errorf("skill dir: got %v; want claude_skill", pt)
	}
}

func TestDetectPackageTypeHookPackage(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	os.MkdirAll(hooksDir, 0o755)
	os.WriteFile(filepath.Join(hooksDir, "hooks.json"), []byte("{}"), 0o644)
	pt, _ := validation.DetectPackageType(dir)
	if pt != validation.PackageTypeHookPackage {
		t.Errorf("hooks dir: got %v; want hook_package", pt)
	}
}
