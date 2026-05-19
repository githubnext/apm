package apmpackage_test

import (
	"testing"

	"github.com/githubnext/apm/internal/models/apmpackage"
)

func TestContentTypeString_AllKnown(t *testing.T) {
	cases := []struct {
		ct   apmpackage.PackageContentType
		want string
	}{
		{apmpackage.ContentTypeInstructions, "instructions"},
		{apmpackage.ContentTypeSkill, "skill"},
		{apmpackage.ContentTypeHybrid, "hybrid"},
		{apmpackage.ContentTypePrompts, "prompts"},
	}
	for _, tc := range cases {
		if got := tc.ct.String(); got != tc.want {
			t.Errorf("ContentType(%d).String() = %q, want %q", tc.ct, got, tc.want)
		}
	}
}

func TestParseContentType_MixedCase_Hybrid(t *testing.T) {
	got, err := apmpackage.ParseContentType("HyBrId")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != apmpackage.ContentTypeHybrid {
		t.Errorf("expected ContentTypeHybrid, got %v", got)
	}
}

func TestParseContentType_MixedCase_Prompts(t *testing.T) {
	got, err := apmpackage.ParseContentType("PROMPTS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != apmpackage.ContentTypePrompts {
		t.Errorf("expected ContentTypePrompts, got %v", got)
	}
}

func TestParseContentType_Whitespace_Errors(t *testing.T) {
	_, err := apmpackage.ParseContentType("  instructions  ")
	if err == nil {
		t.Error("expected error for whitespace-padded input")
	}
}

func TestAPMPackage_Fields_Assignable(t *testing.T) {
	ct := apmpackage.ContentTypeSkill
	pkg := &apmpackage.APMPackage{
		Name:        "mypkg",
		Version:     "1.2.3",
		Description: "A test package",
		Author:      "author",
		License:     "MIT",
		Source:      "owner/repo",
		Type:        &ct,
	}
	if pkg.Name != "mypkg" {
		t.Errorf("Name: got %q", pkg.Name)
	}
	if pkg.Type == nil || *pkg.Type != apmpackage.ContentTypeSkill {
		t.Error("Type field not set correctly")
	}
}

func TestAPMPackage_DependenciesNilByDefault(t *testing.T) {
	pkg := &apmpackage.APMPackage{}
	if pkg.Dependencies != nil {
		t.Error("expected nil Dependencies on zero-value APMPackage")
	}
	if pkg.DevDependencies != nil {
		t.Error("expected nil DevDependencies on zero-value APMPackage")
	}
}

func TestPackageInfo_InstallPathPreserved(t *testing.T) {
	info := &apmpackage.PackageInfo{InstallPath: "/opt/packages/foo"}
	if info.InstallPath != "/opt/packages/foo" {
		t.Errorf("InstallPath not preserved: %q", info.InstallPath)
	}
}

func TestPackageInfo_PackageTypeField(t *testing.T) {
	info := &apmpackage.PackageInfo{PackageType: "APM_PACKAGE"}
	if info.PackageType != "APM_PACKAGE" {
		t.Errorf("PackageType not preserved: %q", info.PackageType)
	}
}

func TestContentTypeString_UnknownMultiple(t *testing.T) {
	for _, v := range []apmpackage.PackageContentType{100, 200, 999} {
		if got := v.String(); got != "unknown" {
			t.Errorf("ContentType(%d).String() = %q, want 'unknown'", v, got)
		}
	}
}
