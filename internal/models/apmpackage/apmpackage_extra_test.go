package apmpackage_test

import (
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/models/apmpackage"
)

func TestParseContentType_CaseVariants(t *testing.T) {
	cases := []string{"INSTRUCTIONS", "Skill", "HYBRID", "Prompts"}
	for _, c := range cases {
		_, err := apmpackage.ParseContentType(c)
		if err != nil {
			t.Errorf("ParseContentType(%q) unexpected error: %v", c, err)
		}
	}
}

func TestParseContentType_EmptyString(t *testing.T) {
	_, err := apmpackage.ParseContentType("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestContentTypeString_Unknown(t *testing.T) {
	var ct apmpackage.PackageContentType = 999
	got := ct.String()
	if got != "unknown" {
		t.Errorf("expected 'unknown' for unrecognized type, got %q", got)
	}
}

func TestPackageInfo_GetPrimitivesPath(t *testing.T) {
	info := &apmpackage.PackageInfo{InstallPath: "/some/path"}
	got := info.GetPrimitivesPath()
	want := filepath.Join("/some/path", ".apm")
	if got != want {
		t.Errorf("GetPrimitivesPath: got %q want %q", got, want)
	}
}

func TestPackageInfo_HasPrimitives_NoDir(t *testing.T) {
	info := &apmpackage.PackageInfo{InstallPath: "/nonexistent/path/xyz"}
	if info.HasPrimitives() {
		t.Error("expected false for non-existent install path")
	}
}

func TestAPMPackage_ZeroValue(t *testing.T) {
	pkg := &apmpackage.APMPackage{}
	if pkg.Name != "" {
		t.Errorf("expected empty name, got %q", pkg.Name)
	}
	if pkg.Version != "" {
		t.Errorf("expected empty version, got %q", pkg.Version)
	}
}

func TestParseContentType_AllValidTypes(t *testing.T) {
	valid := []struct {
		s    string
		want apmpackage.PackageContentType
	}{
		{"instructions", apmpackage.ContentTypeInstructions},
		{"skill", apmpackage.ContentTypeSkill},
		{"hybrid", apmpackage.ContentTypeHybrid},
		{"prompts", apmpackage.ContentTypePrompts},
	}
	for _, tc := range valid {
		got, err := apmpackage.ParseContentType(tc.s)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc.s, err)
		}
		if got != tc.want {
			t.Errorf("ParseContentType(%q) = %v, want %v", tc.s, got, tc.want)
		}
		if got.String() != tc.s {
			t.Errorf("String() round-trip failed for %q: got %q", tc.s, got.String())
		}
	}
}
