package models

import (
	"testing"
)

func TestPlacementStrategy_Constants(t *testing.T) {
	if PlacementStrategySinglePoint == "" {
		t.Error("PlacementStrategySinglePoint should not be empty")
	}
	if PlacementStrategySelectiveMulti == "" {
		t.Error("PlacementStrategySelectiveMulti should not be empty")
	}
	if PlacementStrategyDistributed == "" {
		t.Error("PlacementStrategyDistributed should not be empty")
	}
	if PlacementStrategySinglePoint == PlacementStrategyDistributed {
		t.Error("strategies should be distinct")
	}
}

func TestProjectAnalysis_ZeroValue(t *testing.T) {
	var pa ProjectAnalysis
	if pa.DirectoriesScanned != 0 || pa.FilesAnalyzed != 0 {
		t.Error("expected zero value fields")
	}
	summary := pa.GetFileTypesSummary()
	if summary != "none" {
		t.Errorf("expected 'none' for empty file types, got %q", summary)
	}
}

func TestProjectAnalysis_GetFileTypesSummary_One(t *testing.T) {
	pa := &ProjectAnalysis{FileTypesDetected: []string{".go"}}
	s := pa.GetFileTypesSummary()
	if s != "go" {
		t.Errorf("expected 'go', got %q", s)
	}
}

func TestProjectAnalysis_GetFileTypesSummary_Two(t *testing.T) {
	pa := &ProjectAnalysis{FileTypesDetected: []string{".py", ".go"}}
	s := pa.GetFileTypesSummary()
	if s != "go, py" {
		t.Errorf("expected 'go, py', got %q", s)
	}
}

func TestProjectAnalysis_GetFileTypesSummary_Three(t *testing.T) {
	pa := &ProjectAnalysis{FileTypesDetected: []string{".ts", ".go", ".py"}}
	s := pa.GetFileTypesSummary()
	if s != "go, py, ts" {
		t.Errorf("expected 'go, py, ts', got %q", s)
	}
}

func TestProjectAnalysis_GetFileTypesSummary_FourPlus(t *testing.T) {
	pa := &ProjectAnalysis{FileTypesDetected: []string{".ts", ".go", ".py", ".rs", ".c"}}
	s := pa.GetFileTypesSummary()
	if s == "" {
		t.Error("expected non-empty summary for 5 types")
	}
	// Should contain "and N more"
	if len(s) < 10 {
		t.Errorf("expected longer summary for 5 types, got %q", s)
	}
}

func TestProjectAnalysis_GetFileTypesSummary_LeadingDot(t *testing.T) {
	pa := &ProjectAnalysis{FileTypesDetected: []string{"..go", ".py"}}
	s := pa.GetFileTypesSummary()
	// leading dots stripped
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestOptimizationDecision_Fields(t *testing.T) {
	od := OptimizationDecision{
		InstructionName:     "my-instr",
		Pattern:             "**/*.go",
		MatchingDirectories: 5,
		TotalDirectories:    20,
		Strategy:            PlacementStrategyDistributed,
		RelevanceScore:      0.9,
	}
	ratio := od.DistributionRatio()
	if ratio != 0.25 {
		t.Errorf("expected ratio=0.25, got %f", ratio)
	}
}

func TestOptimizationDecision_PlacementDirectories(t *testing.T) {
	od := OptimizationDecision{
		PlacementDirectories: []string{"src/", "tests/"},
		Reasoning:            "multi-dir coverage",
	}
	if len(od.PlacementDirectories) != 2 {
		t.Errorf("expected 2 dirs, got %d", len(od.PlacementDirectories))
	}
}

func TestPlacementSummary_Fields(t *testing.T) {
	ps := PlacementSummary{
		Path:             "src/AGENTS.md",
		InstructionCount: 3,
		SourceCount:      2,
		Sources:          []string{"pkg-a", "pkg-b"},
	}
	if ps.InstructionCount != 3 {
		t.Errorf("unexpected InstructionCount %d", ps.InstructionCount)
	}
	if len(ps.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(ps.Sources))
	}
}

func TestOptimizationStats_PollutionImprovement(t *testing.T) {
	v := 12.5
	os := &OptimizationStats{PollutionImprovement: &v}
	if os.PollutionImprovement == nil || *os.PollutionImprovement != 12.5 {
		t.Error("unexpected PollutionImprovement")
	}
}

func TestOptimizationStats_EfficiencyImprovement_Negative(t *testing.T) {
	baseline := 0.8
	os := &OptimizationStats{
		AverageContextEfficiency: 0.6,
		BaselineEfficiency:       &baseline,
	}
	imp := os.EfficiencyImprovement()
	if imp == nil {
		t.Fatal("expected non-nil improvement")
	}
	if *imp >= 0 {
		t.Errorf("expected negative improvement when efficiency dropped, got %f", *imp)
	}
}

func TestCompilationResults_NewDefaults(t *testing.T) {
	cr := NewCompilationResults()
	if cr.TargetName == "" {
		t.Error("expected non-empty TargetName default")
	}
	if cr.HasIssues() {
		t.Error("expected no issues for fresh CompilationResults")
	}
	if cr.TotalInstructions() != 0 {
		t.Error("expected 0 total instructions initially")
	}
}

func TestCompilationResults_HasIssues_Warnings(t *testing.T) {
	cr := NewCompilationResults()
	cr.Warnings = []string{"warn1"}
	if !cr.HasIssues() {
		t.Error("expected HasIssues=true with warnings")
	}
}

func TestCompilationResults_HasIssues_Errors(t *testing.T) {
	cr := NewCompilationResults()
	cr.Errors = []string{"err1", "err2"}
	if !cr.HasIssues() {
		t.Error("expected HasIssues=true with errors")
	}
}

func TestCompilationResults_TotalInstructions_Multiple(t *testing.T) {
	cr := NewCompilationResults()
	cr.PlacementSummaries = []PlacementSummary{
		{InstructionCount: 3},
		{InstructionCount: 7},
	}
	if cr.TotalInstructions() != 10 {
		t.Errorf("expected TotalInstructions=10, got %d", cr.TotalInstructions())
	}
}

func TestCompilationResults_DryRun(t *testing.T) {
	cr := NewCompilationResults()
	cr.IsDryRun = true
	if !cr.IsDryRun {
		t.Error("expected IsDryRun=true")
	}
}
