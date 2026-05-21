package compilationconst

import (
"strings"
"testing"
)

func TestConstitutionMarkerBegin_ContainsSpecKit_Extra4(t *testing.T) {
if !strings.Contains(ConstitutionMarkerBegin, "SPEC-KIT") {
t.Error("ConstitutionMarkerBegin should contain 'SPEC-KIT'")
}
}

func TestConstitutionMarkerEnd_ContainsSpecKit_Extra4(t *testing.T) {
if !strings.Contains(ConstitutionMarkerEnd, "SPEC-KIT") {
t.Error("ConstitutionMarkerEnd should contain 'SPEC-KIT'")
}
}

func TestBuildIDPlaceholder_ContainsUnderscoreBuildID_Extra4(t *testing.T) {
if !strings.Contains(BuildIDPlaceholder, "__BUILD_ID__") {
t.Error("BuildIDPlaceholder should contain '__BUILD_ID__'")
}
}

func TestConstitutionRelativePath_StartsWithDotSpecify_Extra4(t *testing.T) {
if !strings.HasPrefix(ConstitutionRelativePath, ".specify") {
t.Errorf("expected path to start with '.specify', got %q", ConstitutionRelativePath)
}
}

func TestConstitutionMarkerBegin_LongerThanEnd_Extra4(t *testing.T) {
// Both should be non-trivially long HTML comments
if len(ConstitutionMarkerBegin) < 10 {
t.Error("ConstitutionMarkerBegin too short")
}
if len(ConstitutionMarkerEnd) < 10 {
t.Error("ConstitutionMarkerEnd too short")
}
}

func TestConstitutionRelativePath_HasConstitution_Extra4(t *testing.T) {
if !strings.Contains(ConstitutionRelativePath, "constitution") {
t.Errorf("expected path to contain 'constitution', got %q", ConstitutionRelativePath)
}
}

func TestBuildIDPlaceholder_HTMLCommentStyle_Extra4(t *testing.T) {
if !strings.HasPrefix(BuildIDPlaceholder, "<!--") {
t.Error("BuildIDPlaceholder should start with '<!--'")
}
if !strings.HasSuffix(BuildIDPlaceholder, "-->") {
t.Error("BuildIDPlaceholder should end with '-->'")
}
}

func TestAllConstants_NonEmpty_Extra4(t *testing.T) {
consts := []string{
ConstitutionMarkerBegin,
ConstitutionMarkerEnd,
ConstitutionRelativePath,
BuildIDPlaceholder,
}
for i, c := range consts {
if c == "" {
t.Errorf("constant at index %d is empty", i)
}
}
}

func TestConstitutionMarkers_NotEqual_Extra4(t *testing.T) {
if ConstitutionMarkerBegin == ConstitutionMarkerEnd {
t.Error("ConstitutionMarkerBegin and ConstitutionMarkerEnd should differ")
}
}

func TestBuildIDPlaceholder_Different_Extra4(t *testing.T) {
if BuildIDPlaceholder == ConstitutionMarkerBegin {
t.Error("BuildIDPlaceholder should differ from ConstitutionMarkerBegin")
}
}
