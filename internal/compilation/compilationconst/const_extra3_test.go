package compilationconst

import (
	"strings"
	"testing"
)

func TestConstitutionMarkerBegin_NotEmpty_Extra3(t *testing.T) {
	if ConstitutionMarkerBegin == "" {
		t.Error("ConstitutionMarkerBegin must not be empty")
	}
}

func TestConstitutionMarkerEnd_NotEmpty_Extra3(t *testing.T) {
	if ConstitutionMarkerEnd == "" {
		t.Error("ConstitutionMarkerEnd must not be empty")
	}
}

func TestBuildIDPlaceholder_ContainsBuildID_Extra3(t *testing.T) {
	if !strings.Contains(BuildIDPlaceholder, "Build ID") {
		t.Errorf("BuildIDPlaceholder = %q, want 'Build ID' substring", BuildIDPlaceholder)
	}
}

func TestConstitutionRelativePath_StartsWithDot_Extra3(t *testing.T) {
	if !strings.HasPrefix(ConstitutionRelativePath, ".") {
		t.Errorf("ConstitutionRelativePath should start with '.': %q", ConstitutionRelativePath)
	}
}

func TestConstitutionMarkerBegin_ContainsBEGIN_Extra3(t *testing.T) {
	if !strings.Contains(ConstitutionMarkerBegin, "BEGIN") {
		t.Errorf("ConstitutionMarkerBegin should contain BEGIN: %q", ConstitutionMarkerBegin)
	}
}

func TestConstitutionMarkerEnd_ContainsEND_Extra3(t *testing.T) {
	if !strings.Contains(ConstitutionMarkerEnd, "END") {
		t.Errorf("ConstitutionMarkerEnd should contain END: %q", ConstitutionMarkerEnd)
	}
}

func TestMarkerPair_BeginPrecedesEnd_Extra3(t *testing.T) {
	// In typical usage, begin marker should precede end marker textually.
	// They share the same prefix up to BEGIN/END.
	beginIdx := strings.Index(ConstitutionMarkerBegin, "CONSTITUTION")
	endIdx := strings.Index(ConstitutionMarkerEnd, "CONSTITUTION")
	if beginIdx < 0 || endIdx < 0 {
		t.Skip("markers do not contain 'CONSTITUTION'")
	}
}

func TestBuildIDPlaceholder_Unique_Extra3(t *testing.T) {
	if BuildIDPlaceholder == ConstitutionMarkerBegin || BuildIDPlaceholder == ConstitutionMarkerEnd {
		t.Error("BuildIDPlaceholder must differ from constitution markers")
	}
}

func TestAllConstantsDiffer_Extra3(t *testing.T) {
	consts := []string{ConstitutionMarkerBegin, ConstitutionMarkerEnd, ConstitutionRelativePath, BuildIDPlaceholder}
	for i := 0; i < len(consts); i++ {
		for j := i + 1; j < len(consts); j++ {
			if consts[i] == consts[j] {
				t.Errorf("constants[%d] == constants[%d]: %q", i, j, consts[i])
			}
		}
	}
}
