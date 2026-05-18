package compilationconst

import (
	"strings"
	"testing"
)

func TestConstitutionMarkers_AreHTMLComments(t *testing.T) {
	if !strings.HasPrefix(ConstitutionMarkerBegin, "<!--") {
		t.Errorf("ConstitutionMarkerBegin should start with HTML comment: %q", ConstitutionMarkerBegin)
	}
	if !strings.HasSuffix(ConstitutionMarkerEnd, "-->") {
		t.Errorf("ConstitutionMarkerEnd should end with HTML comment: %q", ConstitutionMarkerEnd)
	}
}

func TestConstitutionMarkers_ContainSpecKit(t *testing.T) {
	upper := strings.ToUpper(ConstitutionMarkerBegin)
	if !strings.Contains(upper, "SPEC") && !strings.Contains(upper, "CONSTITUTION") {
		t.Errorf("ConstitutionMarkerBegin should reference constitution or spec-kit: %q", ConstitutionMarkerBegin)
	}
}

func TestConstitutionRelativePath_StartsFromRoot(t *testing.T) {
	if strings.HasPrefix(ConstitutionRelativePath, "/") {
		t.Errorf("ConstitutionRelativePath should be relative, not absolute: %q", ConstitutionRelativePath)
	}
}

func TestConstitutionRelativePath_HasMDExtension(t *testing.T) {
	if !strings.HasSuffix(ConstitutionRelativePath, ".md") {
		t.Errorf("ConstitutionRelativePath should be a .md file: %q", ConstitutionRelativePath)
	}
}

func TestBuildIDPlaceholder_IsHTMLComment(t *testing.T) {
	if !strings.HasPrefix(BuildIDPlaceholder, "<!--") || !strings.HasSuffix(BuildIDPlaceholder, "-->") {
		t.Errorf("BuildIDPlaceholder should be an HTML comment: %q", BuildIDPlaceholder)
	}
}

func TestBuildIDPlaceholder_ContainsBuildID(t *testing.T) {
	if !strings.Contains(BuildIDPlaceholder, "Build ID") && !strings.Contains(BuildIDPlaceholder, "BUILD_ID") {
		t.Errorf("BuildIDPlaceholder should reference Build ID: %q", BuildIDPlaceholder)
	}
}

func TestBuildIDPlaceholder_NonEmpty(t *testing.T) {
	if len(BuildIDPlaceholder) == 0 {
		t.Error("BuildIDPlaceholder must not be empty")
	}
}

func TestMarkerBeginNotEqualEnd(t *testing.T) {
	if ConstitutionMarkerBegin == ConstitutionMarkerEnd {
		t.Error("ConstitutionMarkerBegin and ConstitutionMarkerEnd must be distinct")
	}
}

func TestConstitutionRelativePath_IsNotEmpty(t *testing.T) {
	if ConstitutionRelativePath == "" {
		t.Error("ConstitutionRelativePath must not be empty")
	}
}

func TestMarkerBeginContainsBegin(t *testing.T) {
	if !strings.Contains(strings.ToUpper(ConstitutionMarkerBegin), "BEGIN") {
		t.Errorf("ConstitutionMarkerBegin should contain BEGIN: %q", ConstitutionMarkerBegin)
	}
}

func TestMarkerEndContainsEnd(t *testing.T) {
	if !strings.Contains(strings.ToUpper(ConstitutionMarkerEnd), "END") {
		t.Errorf("ConstitutionMarkerEnd should contain END: %q", ConstitutionMarkerEnd)
	}
}

func TestConstitutionRelativePath_NoBrokenSegments(t *testing.T) {
	// Should not have consecutive slashes
	if strings.Contains(ConstitutionRelativePath, "//") {
		t.Errorf("ConstitutionRelativePath has consecutive slashes: %q", ConstitutionRelativePath)
	}
}
