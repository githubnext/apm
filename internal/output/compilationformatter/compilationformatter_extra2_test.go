package compilationformatter_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/output/compilationformatter"
)

func TestOptimizationStats_ZeroValue(t *testing.T) {
	var s compilationformatter.OptimizationStats
	if s.AverageContextEfficiency != 0 {
		t.Error("zero value should have zero efficiency")
	}
	if s.EfficiencyPercentage() != 0 {
		t.Error("zero efficiency should give zero percentage")
	}
}

func TestOptimizationStats_EfficiencyPercentage_Half(t *testing.T) {
	s := compilationformatter.OptimizationStats{AverageContextEfficiency: 0.5}
	if s.EfficiencyPercentage() != 50.0 {
		t.Errorf("expected 50.0, got %f", s.EfficiencyPercentage())
	}
}

func TestOptimizationStats_EfficiencyImprovement_Nil(t *testing.T) {
	s := compilationformatter.OptimizationStats{AverageContextEfficiency: 0.8}
	if s.EfficiencyImprovement() != nil {
		t.Error("expected nil when no baseline")
	}
}

func TestOptimizationStats_EfficiencyImprovement_Positive(t *testing.T) {
	base := 0.6
	s := compilationformatter.OptimizationStats{
		AverageContextEfficiency: 0.8,
		BaselineEfficiency:       &base,
	}
	imp := s.EfficiencyImprovement()
	if imp == nil {
		t.Fatal("expected non-nil improvement")
	}
	if *imp <= 0 {
		t.Errorf("expected positive improvement, got %f", *imp)
	}
}

func TestProjectAnalysis_Fields(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{
		DirectoriesScanned:          5,
		FilesAnalyzed:               20,
		FileTypesDetected:           []string{".go", ".py"},
		InstructionPatternsDetected: 3,
		MaxDepth:                    4,
		ConstitutionDetected:        true,
		ConstitutionPath:            "/root/AGENTS.md",
	}
	if p.DirectoriesScanned != 5 {
		t.Error("DirectoriesScanned mismatch")
	}
	if !p.ConstitutionDetected {
		t.Error("ConstitutionDetected should be true")
	}
}

func TestProjectAnalysis_FileTypesSummary_TwoTypes(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{
		FileTypesDetected: []string{".go", ".py"},
	}
	s := p.FileTypesSummary()
	if !strings.Contains(s, "go") {
		t.Errorf("summary should contain 'go': %q", s)
	}
}

func TestPlacementStrategy_Values(t *testing.T) {
	cases := []struct {
		s    compilationformatter.PlacementStrategy
		want string
	}{
		{compilationformatter.StrategySinglePoint, "Single Point"},
		{compilationformatter.StrategySelectiveMulti, "Selective Multi"},
		{compilationformatter.StrategyDistributed, "Distributed"},
	}
	for _, c := range cases {
		if string(c.s) != c.want {
			t.Errorf("expected %q, got %q", c.want, c.s)
		}
	}
}

func TestCompilationResults_ZeroHasNoIssues(t *testing.T) {
	var r compilationformatter.CompilationResults
	if r.HasIssues() {
		t.Error("zero value should have no issues")
	}
}

func TestCompilationResults_WithErrorHasIssues(t *testing.T) {
	r := compilationformatter.CompilationResults{
		Errors: []string{"some error"},
	}
	if !r.HasIssues() {
		t.Error("should have issues when Errors is non-empty")
	}
}

func TestNew_UseColorField(t *testing.T) {
	f := compilationformatter.New(true)
	if !f.UseColor {
		t.Error("UseColor should be true")
	}
	f2 := compilationformatter.New(false)
	if f2.UseColor {
		t.Error("UseColor should be false")
	}
}
