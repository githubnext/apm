package compilationconst

import (
	"strings"
	"testing"
)

func TestConstitutionMarkers(t *testing.T) {
	if !strings.Contains(ConstitutionMarkerBegin, "BEGIN") {
		t.Error("ConstitutionMarkerBegin should contain BEGIN")
	}
	if !strings.Contains(ConstitutionMarkerEnd, "END") {
		t.Error("ConstitutionMarkerEnd should contain END")
	}
	if ConstitutionMarkerBegin == ConstitutionMarkerEnd {
		t.Error("begin and end markers should differ")
	}
}

func TestConstitutionRelativePath(t *testing.T) {
	if !strings.HasSuffix(ConstitutionRelativePath, "constitution.md") {
		t.Errorf("ConstitutionRelativePath = %q, want suffix constitution.md", ConstitutionRelativePath)
	}
}

func TestBuildIDPlaceholder(t *testing.T) {
	if !strings.Contains(BuildIDPlaceholder, "__BUILD_ID__") {
		t.Errorf("BuildIDPlaceholder = %q, want __BUILD_ID__ placeholder", BuildIDPlaceholder)
	}
}
