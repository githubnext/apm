package compilationconst_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestConstitutionMarkerBegin_ContainsBEGIN_Extra4(t *testing.T) {
	if !strings.Contains(compilationconst.ConstitutionMarkerBegin, "BEGIN") {
		t.Errorf("expected BEGIN in marker, got %q", compilationconst.ConstitutionMarkerBegin)
	}
}

func TestConstitutionMarkerEnd_ContainsEND_Extra4(t *testing.T) {
	if !strings.Contains(compilationconst.ConstitutionMarkerEnd, "END") {
		t.Errorf("expected END in marker, got %q", compilationconst.ConstitutionMarkerEnd)
	}
}

func TestBuildIDPlaceholder_ContainsBuildID_Extra4(t *testing.T) {
	if !strings.Contains(compilationconst.BuildIDPlaceholder, "Build ID") {
		t.Errorf("expected 'Build ID' in placeholder, got %q", compilationconst.BuildIDPlaceholder)
	}
}

func TestBuildIDPlaceholder_HasUnderscoreSentinel_Extra4(t *testing.T) {
	if !strings.Contains(compilationconst.BuildIDPlaceholder, "__BUILD_ID__") {
		t.Errorf("expected __BUILD_ID__ sentinel, got %q", compilationconst.BuildIDPlaceholder)
	}
}

func TestConstitutionRelativePath_StartsWithDot_Extra4(t *testing.T) {
	if !strings.HasPrefix(compilationconst.ConstitutionRelativePath, ".") {
		t.Errorf("expected dot-relative path, got %q", compilationconst.ConstitutionRelativePath)
	}
}

func TestConstitutionRelativePath_NoLeadingSlash_Extra4(t *testing.T) {
	if strings.HasPrefix(compilationconst.ConstitutionRelativePath, "/") {
		t.Errorf("expected no leading slash, got %q", compilationconst.ConstitutionRelativePath)
	}
}

func TestConstitutionMarkerBegin_IsHTMLComment_Extra4(t *testing.T) {
	if !strings.HasPrefix(compilationconst.ConstitutionMarkerBegin, "<!--") {
		t.Errorf("expected HTML comment prefix, got %q", compilationconst.ConstitutionMarkerBegin)
	}
}

func TestConstitutionMarkerEnd_IsHTMLComment_Extra4(t *testing.T) {
	if !strings.HasPrefix(compilationconst.ConstitutionMarkerEnd, "<!--") {
		t.Errorf("expected HTML comment prefix, got %q", compilationconst.ConstitutionMarkerEnd)
	}
}

func TestBuildIDPlaceholder_IsHTMLComment_Extra4(t *testing.T) {
	if !strings.HasPrefix(compilationconst.BuildIDPlaceholder, "<!--") {
		t.Errorf("expected HTML comment prefix, got %q", compilationconst.BuildIDPlaceholder)
	}
}

func TestConstitutionMarkerBegin_LongerThan10_Extra4(t *testing.T) {
	if len(compilationconst.ConstitutionMarkerBegin) <= 10 {
		t.Errorf("expected length > 10, got %d", len(compilationconst.ConstitutionMarkerBegin))
	}
}

func TestConstitutionRelativePath_HasMemorySegment_Extra4(t *testing.T) {
	if !strings.Contains(compilationconst.ConstitutionRelativePath, "memory") {
		t.Errorf("expected 'memory' in path, got %q", compilationconst.ConstitutionRelativePath)
	}
}
