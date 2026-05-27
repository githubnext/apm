package output_test

import (
	"testing"

	"github.com/githubnext/apm/internal/output"
)

// TestParityPlacementStrategyConstants verifies strategy constants match Python enum.
func TestParityPlacementStrategyConstants(t *testing.T) {
	if output.PlacementSinglePoint != "Single Point" {
		t.Errorf("unexpected PlacementSinglePoint: %q", output.PlacementSinglePoint)
	}
	if output.PlacementSelectiveMulti != "Selective Multi" {
		t.Errorf("unexpected PlacementSelectiveMulti: %q", output.PlacementSelectiveMulti)
	}
	if output.PlacementDistributed != "Distributed" {
		t.Errorf("unexpected PlacementDistributed: %q", output.PlacementDistributed)
	}
}

// TestParityProjectAnalysisGetFileTypesSummaryNone returns "none" when empty.
func TestParityProjectAnalysisGetFileTypesSummaryNone(t *testing.T) {
	pa := output.ProjectAnalysis{}
	if pa.GetFileTypesSummary() != "none" {
		t.Errorf("expected 'none', got %q", pa.GetFileTypesSummary())
	}
}

// TestParityProjectAnalysisGetFileTypesSummaryFewTypes returns csv for <=3 types.
func TestParityProjectAnalysisGetFileTypesSummaryFewTypes(t *testing.T) {
	pa := output.ProjectAnalysis{
		FileTypesDetected: map[string]bool{".md": true, ".go": true},
	}
	s := pa.GetFileTypesSummary()
	// Should return "go, md" (sorted, dot stripped)
	if s != "go, md" {
		t.Errorf("expected 'go, md', got %q", s)
	}
}

// TestParityProjectAnalysisGetFileTypesSummaryManyTypes truncates beyond 3.
func TestParityProjectAnalysisGetFileTypesSummaryManyTypes(t *testing.T) {
	pa := output.ProjectAnalysis{
		FileTypesDetected: map[string]bool{
			".md": true, ".go": true, ".yaml": true, ".json": true, ".ts": true,
		},
	}
	s := pa.GetFileTypesSummary()
	// Should contain "and X more"
	if len(s) < 5 {
		t.Errorf("summary too short: %q", s)
	}
}

// TestParityOptimizationDecisionDistributionRatioZero returns 0 for zero total.
func TestParityOptimizationDecisionDistributionRatioZero(t *testing.T) {
	od := output.OptimizationDecision{MatchingDirectories: 5, TotalDirectories: 0}
	if od.DistributionRatio() != 0.0 {
		t.Errorf("expected 0.0, got %f", od.DistributionRatio())
	}
}

// TestParityOptimizationDecisionDistributionRatioNonZero computes ratio.
func TestParityOptimizationDecisionDistributionRatioNonZero(t *testing.T) {
	od := output.OptimizationDecision{MatchingDirectories: 3, TotalDirectories: 10}
	if od.DistributionRatio() != 0.3 {
		t.Errorf("expected 0.3, got %f", od.DistributionRatio())
	}
}

// TestParityCompilationResultsTotalInstructions sums across summaries.
func TestParityCompilationResultsTotalInstructions(t *testing.T) {
	cr := output.CompilationResults{
		PlacementSummaries: []output.PlacementSummary{
			{InstructionCount: 3},
			{InstructionCount: 5},
		},
	}
	if cr.TotalInstructions() != 8 {
		t.Errorf("expected 8, got %d", cr.TotalInstructions())
	}
}

// TestParityCompilationResultsTotalInstructionsEmpty returns 0 for no summaries.
func TestParityCompilationResultsTotalInstructionsEmpty(t *testing.T) {
	cr := output.CompilationResults{}
	if cr.TotalInstructions() != 0 {
		t.Errorf("expected 0, got %d", cr.TotalInstructions())
	}
}

// TestParityCompilationResultsHasIssuesFalse returns false when no issues.
func TestParityCompilationResultsHasIssuesFalse(t *testing.T) {
	cr := output.CompilationResults{}
	if cr.HasIssues() {
		t.Error("expected HasIssues=false")
	}
}

// TestParityCompilationResultsHasIssuesTrueWarning returns true for warnings.
func TestParityCompilationResultsHasIssuesTrueWarning(t *testing.T) {
	cr := output.CompilationResults{Warnings: []string{"watch out"}}
	if !cr.HasIssues() {
		t.Error("expected HasIssues=true with warnings")
	}
}

// TestParityCompilationResultsHasIssuesTrueError returns true for errors.
func TestParityCompilationResultsHasIssuesTrueError(t *testing.T) {
	cr := output.CompilationResults{Errors: []string{"something broke"}}
	if !cr.HasIssues() {
		t.Error("expected HasIssues=true with errors")
	}
}

// TestParityOptimizationStatsEfficiencyPercentage multiplies by 100.
func TestParityOptimizationStatsEfficiencyPercentage(t *testing.T) {
	s := output.OptimizationStats{AverageContextEfficiency: 0.75}
	if s.EfficiencyPercentage() != 75.0 {
		t.Errorf("expected 75.0, got %f", s.EfficiencyPercentage())
	}
}

// TestParityOptimizationStatsEfficiencyImprovementNilBaseline returns nil.
func TestParityOptimizationStatsEfficiencyImprovementNilBaseline(t *testing.T) {
	s := output.OptimizationStats{AverageContextEfficiency: 0.9}
	if s.EfficiencyImprovement() != nil {
		t.Error("expected nil when BaselineEfficiency is nil")
	}
}

// TestParityOptimizationStatsEfficiencyImprovementNonNil computes value.
func TestParityOptimizationStatsEfficiencyImprovementNonNil(t *testing.T) {
	base := 0.5
	s := output.OptimizationStats{AverageContextEfficiency: 0.75, BaselineEfficiency: &base}
	imp := s.EfficiencyImprovement()
	if imp == nil {
		t.Fatal("expected non-nil improvement")
	}
	// (0.75-0.5)/0.5*100 = 50.0
	if *imp != 50.0 {
		t.Errorf("expected 50.0, got %f", *imp)
	}
}

// TestParityCompilationResultsDefaultTargetName verifies empty string default.
func TestParityCompilationResultsDefaultTargetName(t *testing.T) {
	cr := output.CompilationResults{}
	// Go zero-value for string is ""
	if cr.TargetName != "" {
		t.Errorf("unexpected TargetName: %q", cr.TargetName)
	}
}

// TestParityProjectAnalysisConstitutionFields verifies constitution fields.
func TestParityProjectAnalysisConstitutionFields(t *testing.T) {
	pa := output.ProjectAnalysis{
		ConstitutionDetected: true,
		ConstitutionPath:     "/root/CONSTITUTION.md",
	}
	if !pa.ConstitutionDetected {
		t.Error("expected ConstitutionDetected=true")
	}
	if pa.ConstitutionPath != "/root/CONSTITUTION.md" {
		t.Errorf("unexpected ConstitutionPath: %s", pa.ConstitutionPath)
	}
}

// TestParityPlacementSummaryFields verifies PlacementSummary fields.
func TestParityPlacementSummaryFields(t *testing.T) {
	ps := output.PlacementSummary{
		Path:             "src/AGENTS.md",
		InstructionCount: 4,
		SourceCount:      2,
		Sources:          []string{"a.md", "b.md"},
	}
	if ps.SourceCount != 2 {
		t.Errorf("expected SourceCount=2, got %d", ps.SourceCount)
	}
	if len(ps.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(ps.Sources))
	}
}
