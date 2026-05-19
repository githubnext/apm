package compilationformatter_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/output/compilationformatter"
)

func TestProjectAnalysis_FileTypesSummary_EmptySlice(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{}
	got := p.FileTypesSummary()
	if got != "none" {
		t.Errorf("empty types summary = %q, want none", got)
	}
}

func TestProjectAnalysis_FileTypesSummary_One(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{FileTypesDetected: []string{".py"}}
	got := p.FileTypesSummary()
	if got != "py" {
		t.Errorf("single type = %q, want py", got)
	}
}

func TestProjectAnalysis_FileTypesSummary_Three(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{
		FileTypesDetected: []string{".py", ".go", ".ts"},
	}
	got := p.FileTypesSummary()
	if !strings.Contains(got, "py") || !strings.Contains(got, "go") || !strings.Contains(got, "ts") {
		t.Errorf("three types = %q, expected py, go, ts", got)
	}
}

func TestProjectAnalysis_FileTypesSummary_MoreThanThree(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{
		FileTypesDetected: []string{".py", ".go", ".ts", ".md", ".yaml"},
	}
	got := p.FileTypesSummary()
	if !strings.Contains(got, "and") {
		t.Errorf("more than 3 types should have 'and': %q", got)
	}
	if !strings.Contains(got, "2 more") {
		t.Errorf("5 types should show '2 more': %q", got)
	}
}

func TestPlacementStrategy_Constants(t *testing.T) {
	if compilationformatter.StrategySinglePoint == "" {
		t.Error("StrategySinglePoint should not be empty")
	}
	if compilationformatter.StrategySelectiveMulti == "" {
		t.Error("StrategySelectiveMulti should not be empty")
	}
	if compilationformatter.StrategyDistributed == "" {
		t.Error("StrategyDistributed should not be empty")
	}
	// All three should be distinct
	if compilationformatter.StrategySinglePoint == compilationformatter.StrategySelectiveMulti {
		t.Error("SinglePoint and SelectiveMulti should be distinct")
	}
}

func TestOptimizationStats_EfficiencyPercentage_Zero(t *testing.T) {
	s := compilationformatter.OptimizationStats{AverageContextEfficiency: 0}
	if s.EfficiencyPercentage() != 0 {
		t.Errorf("zero efficiency = %v", s.EfficiencyPercentage())
	}
}

func TestOptimizationStats_EfficiencyPercentage_One(t *testing.T) {
	s := compilationformatter.OptimizationStats{AverageContextEfficiency: 1.0}
	if s.EfficiencyPercentage() != 100.0 {
		t.Errorf("full efficiency = %v", s.EfficiencyPercentage())
	}
}

func TestCompilationResults_HasIssues_BothEmpty(t *testing.T) {
	r := &compilationformatter.CompilationResults{}
	if r.HasIssues() {
		t.Error("empty results should have no issues")
	}
}

func TestCompilationResults_HasIssues_OneWarning(t *testing.T) {
	r := &compilationformatter.CompilationResults{
		Warnings: []string{"some warning"},
	}
	if !r.HasIssues() {
		t.Error("results with warning should have issues")
	}
}

func TestCompilationResults_HasIssues_OneError(t *testing.T) {
	r := &compilationformatter.CompilationResults{
		Errors: []string{"some error"},
	}
	if !r.HasIssues() {
		t.Error("results with error should have issues")
	}
}

func TestNew_NoColor(t *testing.T) {
	f := compilationformatter.New(false)
	if f == nil {
		t.Error("New(false) should not return nil")
	}
}

func TestNew_WithColor(t *testing.T) {
	f := compilationformatter.New(true)
	if f == nil {
		t.Error("New(true) should not return nil")
	}
}

func TestFormatDefault_ReturnsString(t *testing.T) {
	f := compilationformatter.New(false)
	r := &compilationformatter.CompilationResults{
		TargetName: "vscode",
	}
	out := f.FormatDefault(r)
	if out == "" {
		t.Error("FormatDefault should return non-empty string")
	}
}

func TestFormatDryRun_ReturnsString(t *testing.T) {
	f := compilationformatter.New(false)
	r := &compilationformatter.CompilationResults{
		TargetName: "claude",
	}
	out := f.FormatDryRun(r)
	if out == "" {
		t.Error("FormatDryRun should return non-empty string")
	}
}

func TestProjectAnalysis_ConstitutionDetected(t *testing.T) {
	p := compilationformatter.ProjectAnalysis{
		ConstitutionDetected: true,
		ConstitutionPath:     "/project/.github/copilot-instructions.md",
	}
	if !p.ConstitutionDetected {
		t.Error("ConstitutionDetected should be true")
	}
	if p.ConstitutionPath == "" {
		t.Error("ConstitutionPath should be set")
	}
}
