package compilationconst

import (
	"strings"
	"testing"
)

func TestBuildIDPlaceholder_ContainsUnderscores(t *testing.T) {
	if !strings.Contains(BuildIDPlaceholder, "__") {
		t.Errorf("BuildIDPlaceholder should contain __ sentinel: %q", BuildIDPlaceholder)
	}
}

func TestConstitutionRelativePath_HasSlash(t *testing.T) {
	if !strings.Contains(ConstitutionRelativePath, "/") {
		t.Errorf("ConstitutionRelativePath should be a path with slashes: %q", ConstitutionRelativePath)
	}
}

func TestConstitutionRelativePath_EndsWithMD(t *testing.T) {
	if !strings.HasSuffix(ConstitutionRelativePath, ".md") {
		t.Errorf("ConstitutionRelativePath should end with .md: %q", ConstitutionRelativePath)
	}
}

func TestConstitutionMarkerBegin_EndsWithComment(t *testing.T) {
	if !strings.HasSuffix(ConstitutionMarkerBegin, "-->") {
		t.Errorf("ConstitutionMarkerBegin should end with -->: %q", ConstitutionMarkerBegin)
	}
}

func TestConstitutionMarkerEnd_StartsWithComment(t *testing.T) {
	if !strings.HasPrefix(ConstitutionMarkerEnd, "<!--") {
		t.Errorf("ConstitutionMarkerEnd should start with <!--: %q", ConstitutionMarkerEnd)
	}
}

func TestBuildIDPlaceholder_StartsAndEndsWithComment(t *testing.T) {
	if !strings.HasPrefix(BuildIDPlaceholder, "<!--") {
		t.Errorf("BuildIDPlaceholder should start with <!--: %q", BuildIDPlaceholder)
	}
	if !strings.HasSuffix(BuildIDPlaceholder, "-->") {
		t.Errorf("BuildIDPlaceholder should end with -->: %q", BuildIDPlaceholder)
	}
}

func TestConstitutionMarkerBegin_UniqueFromEnd(t *testing.T) {
	if ConstitutionMarkerBegin == ConstitutionMarkerEnd {
		t.Error("begin and end markers must differ")
	}
}

func TestAllConstants_PrintableASCII(t *testing.T) {
	for _, s := range []string{ConstitutionMarkerBegin, ConstitutionMarkerEnd, ConstitutionRelativePath, BuildIDPlaceholder} {
		for i, r := range s {
			if r > 127 {
				t.Errorf("non-ASCII character at position %d in %q: %q", i, s, r)
			}
		}
	}
}

func TestBuildIDPlaceholder_BuildIDSentinel(t *testing.T) {
	if !strings.Contains(BuildIDPlaceholder, "BUILD_ID") {
		t.Errorf("BuildIDPlaceholder should mention BUILD_ID: %q", BuildIDPlaceholder)
	}
}
