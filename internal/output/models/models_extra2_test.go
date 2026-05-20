package models

import (
	"testing"
)

func TestProjectAnalysis_ZeroValue_Extra2(t *testing.T) {
	var p ProjectAnalysis
	if p.DirectoriesScanned != 0 || p.FilesAnalyzed != 0 || p.ConstitutionDetected {
		t.Error("zero-value ProjectAnalysis should have empty fields")
	}
}

func TestProjectAnalysis_Fields_Extra2(t *testing.T) {
	p := ProjectAnalysis{
		DirectoriesScanned:          5,
		FilesAnalyzed:               20,
		FileTypesDetected:           []string{"py", "go", "md"},
		InstructionPatternsDetected: 3,
		MaxDepth:                    4,
		ConstitutionDetected:        true,
		ConstitutionPath:            "/proj/AGENTS.md",
	}
	if p.DirectoriesScanned != 5 {
		t.Errorf("DirectoriesScanned = %d", p.DirectoriesScanned)
	}
	if len(p.FileTypesDetected) != 3 {
		t.Errorf("FileTypesDetected len = %d", len(p.FileTypesDetected))
	}
	if !p.ConstitutionDetected {
		t.Error("ConstitutionDetected should be true")
	}
}

func TestProjectAnalysis_GetFileTypesSummary_ZeroTypes_Extra2(t *testing.T) {
	p := &ProjectAnalysis{}
	s := p.GetFileTypesSummary()
	if s == "" {
		t.Error("expected non-empty summary even with no types")
	}
}

func TestProjectAnalysis_GetFileTypesSummary_One_Extra2(t *testing.T) {
	p := &ProjectAnalysis{FileTypesDetected: []string{"py"}}
	s := p.GetFileTypesSummary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestOptimizationDecision_ZeroValue_Extra2(t *testing.T) {
	var d OptimizationDecision
	if d.InstructionName != "" || d.MatchingDirectories != 0 {
		t.Error("zero-value OptimizationDecision should have empty fields")
	}
}

func TestOptimizationDecision_Fields_Extra2(t *testing.T) {
	d := OptimizationDecision{
		InstructionName:     "AGENTS.md",
		Pattern:             "*.py",
		MatchingDirectories: 3,
		TotalDirectories:    10,
		DistributionScore:   0.75,
		Reasoning:           "broad match",
		RelevanceScore:      0.9,
	}
	if d.InstructionName != "AGENTS.md" {
		t.Errorf("InstructionName = %q", d.InstructionName)
	}
	if d.DistributionScore != 0.75 {
		t.Errorf("DistributionScore = %f", d.DistributionScore)
	}
}

func TestOptimizationDecision_DistributionRatio_Extra2(t *testing.T) {
	d := OptimizationDecision{MatchingDirectories: 4, TotalDirectories: 8}
	r := d.DistributionRatio()
	if r != 0.5 {
		t.Errorf("DistributionRatio = %f, want 0.5", r)
	}
}

func TestOptimizationDecision_DistributionRatio_ZeroTotal_Extra2(t *testing.T) {
	d := OptimizationDecision{MatchingDirectories: 0, TotalDirectories: 0}
	r := d.DistributionRatio()
	if r != 0.0 {
		t.Errorf("DistributionRatio with zero total = %f, want 0.0", r)
	}
}

func TestPlacementSummary_Fields_Extra2(t *testing.T) {
	ps := PlacementSummary{
		Path:             "/proj",
		InstructionCount: 2,
		SourceCount:      3,
		Sources:          []string{"a.md", "b.md", "c.md"},
	}
	if ps.Path != "/proj" {
		t.Errorf("Path = %q", ps.Path)
	}
	if len(ps.Sources) != 3 {
		t.Errorf("Sources len = %d", len(ps.Sources))
	}
}

func TestOptimizationStats_ZeroValue_Extra2(t *testing.T) {
	var s OptimizationStats
	if s.AverageContextEfficiency != 0.0 || s.TotalAgentsFiles != 0 {
		t.Error("zero-value OptimizationStats should have zero fields")
	}
}

func TestCompilationResults_ZeroValue_Extra2(t *testing.T) {
	var r CompilationResults
	if r.IsDryRun || r.TargetName != "" || len(r.Warnings) != 0 || len(r.Errors) != 0 {
		t.Error("zero-value CompilationResults should have empty fields")
	}
}

func TestCompilationResults_HasIssues_Extra2(t *testing.T) {
	r := CompilationResults{Warnings: []string{"w1"}}
	if !r.HasIssues() {
		t.Error("expected HasIssues=true when Warnings is non-empty")
	}
}

func TestCompilationResults_HasIssues_Errors_Extra2(t *testing.T) {
	r := CompilationResults{Errors: []string{"e1"}}
	if !r.HasIssues() {
		t.Error("expected HasIssues=true when Errors is non-empty")
	}
}

func TestCompilationResults_TotalInstructions_Extra2(t *testing.T) {
	r := CompilationResults{
		PlacementSummaries: []PlacementSummary{
			{InstructionCount: 1},
			{InstructionCount: 2},
			{InstructionCount: 3},
		},
	}
	if r.TotalInstructions() != 6 {
		t.Errorf("TotalInstructions = %d, want 6", r.TotalInstructions())
	}
}

func TestPlacementStrategy_Constants_Extra2(t *testing.T) {
	if PlacementStrategySinglePoint == PlacementStrategySelectiveMulti {
		t.Error("placement strategy constants should be distinct")
	}
	if PlacementStrategySinglePoint == PlacementStrategyDistributed {
		t.Error("placement strategy constants should be distinct")
	}
}
