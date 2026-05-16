package models

import (
"testing"
)

func TestOptimizationDecisionDistributionRatio(t *testing.T) {
o := &OptimizationDecision{
MatchingDirectories: 3,
TotalDirectories:    10,
}
ratio := o.DistributionRatio()
if ratio < 0 || ratio > 1 {
t.Errorf("ratio should be between 0 and 1, got %f", ratio)
}
}

func TestOptimizationDecisionDistributionRatioZero(t *testing.T) {
o := &OptimizationDecision{}
ratio := o.DistributionRatio()
_ = ratio // should not panic
}

func TestOptimizationStatsEfficiencyPercentage(t *testing.T) {
o := &OptimizationStats{
AverageContextEfficiency: 0.8,
}
pct := o.EfficiencyPercentage()
_ = pct
}

func TestOptimizationStatsEfficiencyImprovementNil(t *testing.T) {
o := &OptimizationStats{}
result := o.EfficiencyImprovement()
if result != nil {
t.Error("expected nil when BaselineEfficiency is nil")
}
}

func TestOptimizationStatsEfficiencyImprovementWithBaseline(t *testing.T) {
baseline := 0.5
o := &OptimizationStats{
AverageContextEfficiency: 0.75,
BaselineEfficiency:       &baseline,
}
result := o.EfficiencyImprovement()
if result == nil {
t.Error("expected non-nil improvement")
}
if *result <= 0 {
t.Errorf("expected positive improvement, got %f", *result)
}
}

func TestCompilationResultsMethods(t *testing.T) {
c := NewCompilationResults()
if c == nil {
t.Fatal("expected non-nil CompilationResults")
}
if c.TotalInstructions() != 0 {
t.Errorf("expected 0 total instructions for empty results")
}
if c.HasIssues() {
t.Error("expected no issues for empty results")
}
}

func TestProjectAnalysisGetFileTypesSummary(t *testing.T) {
p := &ProjectAnalysis{
FileTypesDetected: []string{".go", ".py", ".md"},
}
summary := p.GetFileTypesSummary()
if summary == "" {
t.Error("expected non-empty summary")
}
}

func TestProjectAnalysisGetFileTypesSummaryEmpty(t *testing.T) {
p := &ProjectAnalysis{}
summary := p.GetFileTypesSummary()
if summary != "none" {
t.Errorf("expected 'none' for empty file types, got %q", summary)
}
}

func TestProjectAnalysisGetFileTypesSummaryMany(t *testing.T) {
p := &ProjectAnalysis{
FileTypesDetected: []string{".go", ".py", ".md", ".yaml", ".json"},
}
summary := p.GetFileTypesSummary()
if summary == "" {
t.Error("expected non-empty summary for many types")
}
}
