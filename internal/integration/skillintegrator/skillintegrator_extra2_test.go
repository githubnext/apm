package skillintegrator

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// NormalizeSkillName
// ---------------------------------------------------------------------------

func TestNormalizeSkillName_AlreadyValid(t *testing.T) {
	got := NormalizeSkillName("my-skill")
	if got != "my-skill" {
		t.Errorf("expected my-skill, got %q", got)
	}
}

func TestNormalizeSkillName_UpperCase(t *testing.T) {
	got := NormalizeSkillName("MySkill")
	if strings.ContainsAny(got, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Errorf("expected all lowercase, got %q", got)
	}
}

func TestNormalizeSkillName_Underscores(t *testing.T) {
	got := NormalizeSkillName("my_skill")
	if strings.Contains(got, "_") {
		t.Errorf("expected no underscores, got %q", got)
	}
}

func TestNormalizeSkillName_Spaces(t *testing.T) {
	got := NormalizeSkillName("my skill")
	if strings.Contains(got, " ") {
		t.Errorf("expected no spaces, got %q", got)
	}
}

func TestNormalizeSkillName_LeadingTrailingHyphens(t *testing.T) {
	got := NormalizeSkillName("-bad-")
	if strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
		t.Errorf("expected no leading/trailing hyphens, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// ValidateSkillName additional cases
// ---------------------------------------------------------------------------

func TestValidateSkillName_SingleChar(t *testing.T) {
	ok, _ := ValidateSkillName("a")
	if !ok {
		t.Error("single lowercase char should be valid")
	}
}

func TestValidateSkillName_AllDigits(t *testing.T) {
	ok, _ := ValidateSkillName("123")
	if !ok {
		t.Error("all-digits name should be valid")
	}
}

func TestValidateSkillName_MixedValid(t *testing.T) {
	ok, _ := ValidateSkillName("abc-123")
	if !ok {
		t.Error("abc-123 should be valid")
	}
}

func TestValidateSkillName_ConsecutiveHyphens(t *testing.T) {
	ok, _ := ValidateSkillName("a--b")
	if ok {
		t.Error("consecutive hyphens should be invalid")
	}
}

func TestValidateSkillName_UpperCase(t *testing.T) {
	ok, _ := ValidateSkillName("MySkill")
	if ok {
		t.Error("uppercase letters should be invalid")
	}
}

// ---------------------------------------------------------------------------
// SyncStats struct
// ---------------------------------------------------------------------------

func TestSyncStats_ZeroValue(t *testing.T) {
	var s SyncStats
	if s.FilesRemoved != 0 || s.Errors != 0 {
		t.Error("expected zero value SyncStats")
	}
}

func TestSyncStats_Fields(t *testing.T) {
	s := SyncStats{FilesRemoved: 5, Errors: 2}
	if s.FilesRemoved != 5 {
		t.Errorf("expected FilesRemoved=5, got %d", s.FilesRemoved)
	}
	if s.Errors != 2 {
		t.Errorf("expected Errors=2, got %d", s.Errors)
	}
}

// ---------------------------------------------------------------------------
// PackageInfo struct
// ---------------------------------------------------------------------------

func TestPackageInfo_ZeroValue(t *testing.T) {
	var p PackageInfo
	if p.IsVirtual || p.IsSubdir {
		t.Error("expected false booleans in zero value")
	}
}

func TestPackageInfo_Fields(t *testing.T) {
	p := PackageInfo{
		InstallPath: "/opt/apm/pkg",
		PackageType: "SKILL_BUNDLE",
		IsVirtual:   true,
		IsSubdir:    false,
		UniqueKey:   "owner/repo",
	}
	if p.InstallPath != "/opt/apm/pkg" {
		t.Errorf("expected /opt/apm/pkg, got %q", p.InstallPath)
	}
	if p.PackageType != "SKILL_BUNDLE" {
		t.Errorf("expected SKILL_BUNDLE, got %q", p.PackageType)
	}
}

// ---------------------------------------------------------------------------
// ToHyphenCase additional cases
// ---------------------------------------------------------------------------

func TestToHyphenCase_Numbers(t *testing.T) {
	got := ToHyphenCase("myPkg2")
	if !strings.Contains(got, "2") {
		t.Errorf("expected digits preserved, got %q", got)
	}
}

func TestToHyphenCase_AllLower(t *testing.T) {
	got := ToHyphenCase("already")
	if got != "already" {
		t.Errorf("expected already, got %q", got)
	}
}

func TestToHyphenCase_Empty(t *testing.T) {
	got := ToHyphenCase("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
