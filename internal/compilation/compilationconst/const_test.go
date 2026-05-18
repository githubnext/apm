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

func TestConstitutionMarkersAreHTMLComments(t *testing.T) {
	if !strings.HasPrefix(ConstitutionMarkerBegin, "<!--") {
		t.Errorf("ConstitutionMarkerBegin should start with <!--")
	}
	if !strings.HasPrefix(ConstitutionMarkerEnd, "<!--") {
		t.Errorf("ConstitutionMarkerEnd should start with <!--")
	}
	if !strings.HasSuffix(ConstitutionMarkerBegin, "-->") {
		t.Errorf("ConstitutionMarkerBegin should end with -->")
	}
	if !strings.HasSuffix(ConstitutionMarkerEnd, "-->") {
		t.Errorf("ConstitutionMarkerEnd should end with -->")
	}
}

func TestBuildIDPlaceholderIsHTMLComment(t *testing.T) {
	if !strings.HasPrefix(BuildIDPlaceholder, "<!--") {
		t.Errorf("BuildIDPlaceholder should start with <!--")
	}
	if !strings.HasSuffix(BuildIDPlaceholder, "-->") {
		t.Errorf("BuildIDPlaceholder should end with -->")
	}
}

func TestConstitutionRelativePathNotAbsolute(t *testing.T) {
	if strings.HasPrefix(ConstitutionRelativePath, "/") {
		t.Error("ConstitutionRelativePath should not be an absolute path")
	}
}

func TestConstitutionRelativePathContainsMemory(t *testing.T) {
	if !strings.Contains(ConstitutionRelativePath, "memory") {
		t.Errorf("ConstitutionRelativePath %q should contain 'memory'", ConstitutionRelativePath)
	}
}

func TestAllConstantsNonEmpty(t *testing.T) {
	if ConstitutionMarkerBegin == "" {
		t.Error("ConstitutionMarkerBegin must not be empty")
	}
	if ConstitutionMarkerEnd == "" {
		t.Error("ConstitutionMarkerEnd must not be empty")
	}
	if ConstitutionRelativePath == "" {
		t.Error("ConstitutionRelativePath must not be empty")
	}
	if BuildIDPlaceholder == "" {
		t.Error("BuildIDPlaceholder must not be empty")
	}
}

func TestConstantsStability(t *testing.T) {
	// Calling constants multiple times returns identical values.
	if ConstitutionMarkerBegin != ConstitutionMarkerBegin {
		t.Error("ConstitutionMarkerBegin changed between accesses")
	}
	if BuildIDPlaceholder != BuildIDPlaceholder {
		t.Error("BuildIDPlaceholder changed between accesses")
	}
}

func TestConstitutionMarkerBeginContainsConstitution(t *testing.T) {
	if !strings.Contains(ConstitutionMarkerBegin, "CONSTITUTION") &&
		!strings.Contains(ConstitutionMarkerBegin, "constitution") &&
		!strings.Contains(ConstitutionMarkerBegin, "SPEC") {
		t.Errorf("ConstitutionMarkerBegin %q should contain a constitution-related keyword", ConstitutionMarkerBegin)
	}
}

func TestBuildIDPlaceholderContainsBuildID(t *testing.T) {
	if !strings.Contains(BuildIDPlaceholder, "Build ID") &&
		!strings.Contains(BuildIDPlaceholder, "BUILD_ID") {
		t.Errorf("BuildIDPlaceholder %q should mention Build ID", BuildIDPlaceholder)
	}
}
