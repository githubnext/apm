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

func TestCompilationResultsTotalInstructions(t *testing.T) {
c := &CompilationResults{
PlacementSummaries: []PlacementSummary{
{InstructionCount: 3},
{InstructionCount: 7},
},
}
if c.TotalInstructions() != 10 {
t.Errorf("expected 10 total instructions, got %d", c.TotalInstructions())
}
}

func TestCompilationResultsHasIssuesBothEmpty(t *testing.T) {
c := &CompilationResults{}
if c.HasIssues() {
t.Error("expected HasIssues=false for empty warnings+errors")
}
}

func TestCompilationResultsHasIssuesBoth(t *testing.T) {
c := &CompilationResults{
Warnings: []string{"warn"},
Errors:   []string{"err"},
}
if !c.HasIssues() {
t.Error("expected HasIssues=true with both warnings and errors")
}
}

func TestNewCompilationResultsDefaults(t *testing.T) {
c := NewCompilationResults()
if c.TargetName != "AGENTS.md" {
t.Errorf("expected TargetName=AGENTS.md, got %q", c.TargetName)
}
}

func TestOptimizationDecisionDistributionRatioExact(t *testing.T) {
o := &OptimizationDecision{MatchingDirectories: 5, TotalDirectories: 10}
if o.DistributionRatio() != 0.5 {
t.Errorf("expected 0.5, got %f", o.DistributionRatio())
}
}

func TestOptimizationDecisionDistributionRatioFull(t *testing.T) {
o := &OptimizationDecision{MatchingDirectories: 10, TotalDirectories: 10}
if o.DistributionRatio() != 1.0 {
t.Errorf("expected 1.0, got %f", o.DistributionRatio())
}
}

func TestOptimizationStatsEfficiencyImprovementZeroBaseline(t *testing.T) {
baseline := 0.0
o := &OptimizationStats{
AverageContextEfficiency: 0.5,
BaselineEfficiency:       &baseline,
}
result := o.EfficiencyImprovement()
if result != nil {
t.Errorf("expected nil for zero baseline, got %v", result)
}
}

func TestProjectAnalysisAllFields(t *testing.T) {
p := &ProjectAnalysis{
DirectoriesScanned:          5,
FilesAnalyzed:               50,
FileTypesDetected:           []string{".go", ".py"},
InstructionPatternsDetected: 3,
MaxDepth:                    4,
ConstitutionDetected:        true,
ConstitutionPath:            "/root/AGENTS.md",
}
if !p.ConstitutionDetected {
t.Error("expected ConstitutionDetected=true")
}
if p.ConstitutionPath == "" {
t.Error("expected non-empty ConstitutionPath")
}
if p.MaxDepth != 4 {
t.Errorf("expected MaxDepth=4, got %d", p.MaxDepth)
}
}

func TestPlacementStrategies(t *testing.T) {
cases := []PlacementStrategy{
PlacementStrategySinglePoint,
PlacementStrategySelectiveMulti,
PlacementStrategyDistributed,
}
for _, s := range cases {
if string(s) == "" {
t.Errorf("strategy should not be empty: %v", s)
}
}
}

func TestOptimizationDecisionFields(t *testing.T) {
o := &OptimizationDecision{
InstructionName:      "my-inst",
Pattern:              "src/**",
MatchingDirectories:  2,
TotalDirectories:     8,
DistributionScore:    0.25,
Strategy:             PlacementStrategyDistributed,
PlacementDirectories: []string{"src/a", "src/b"},
Reasoning:            "matches pattern",
RelevanceScore:       0.9,
}
if o.InstructionName == "" {
t.Error("InstructionName should not be empty")
}
if len(o.PlacementDirectories) != 2 {
t.Errorf("expected 2 placement dirs, got %d", len(o.PlacementDirectories))
}
}
