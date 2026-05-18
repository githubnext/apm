package compilationformatter_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/output/compilationformatter"
)

func TestOptimizationStats_EfficiencyPercentage(t *testing.T) {
	s := compilationformatter.OptimizationStats{
		AverageContextEfficiency: 0.75,
	}
	pct := s.EfficiencyPercentage()
	if pct != 75.0 {
		t.Fatalf("expected 75.0, got %v", pct)
	}
}

func TestOptimizationStats_EfficiencyImprovement_NoBaseline(t *testing.T) {
	s := compilationformatter.OptimizationStats{
		AverageContextEfficiency: 0.8,
	}
	imp := s.EfficiencyImprovement()
	if imp != nil {
		t.Fatalf("expected nil improvement when BaselineEfficiency=nil, got %v", imp)
	}
}

func TestOptimizationStats_EfficiencyImprovement_WithBaseline(t *testing.T) {
	base := 0.6
	s := compilationformatter.OptimizationStats{
		AverageContextEfficiency: 0.8,
		BaselineEfficiency:       &base,
	}
	imp := s.EfficiencyImprovement()
	if imp == nil {
		t.Fatal("expected non-nil improvement")
	}
	if *imp < 19.9 || *imp > 20.1 {
		t.Fatalf("expected ~20pp improvement, got %v", *imp)
	}
}

func TestCompilationResults_HasIssues_NoIssues(t *testing.T) {
	r := &compilationformatter.CompilationResults{}
	if r.HasIssues() {
		t.Fatal("HasIssues should be false when no warnings/errors")
	}
}

func TestCompilationResults_HasIssues_WithWarning(t *testing.T) {
	r := &compilationformatter.CompilationResults{
		Warnings: []string{"watch out"},
	}
	if !r.HasIssues() {
		t.Fatal("HasIssues should be true with warnings")
	}
}

func TestCompilationResults_HasIssues_WithError(t *testing.T) {
	r := &compilationformatter.CompilationResults{
		Errors: []string{"boom"},
	}
	if !r.HasIssues() {
		t.Fatal("HasIssues should be true with errors")
	}
}

func TestProjectAnalysis_FileTypesSummary_Empty(t *testing.T) {
	p := &compilationformatter.ProjectAnalysis{}
	if p.FileTypesSummary() != "none" {
		t.Fatalf("expected 'none', got %q", p.FileTypesSummary())
	}
}

func TestProjectAnalysis_FileTypesSummary_Few(t *testing.T) {
	p := &compilationformatter.ProjectAnalysis{
		FileTypesDetected: []string{".md", ".py"},
	}
	got := p.FileTypesSummary()
	if !strings.Contains(got, "md") || !strings.Contains(got, "py") {
		t.Fatalf("expected file types in summary, got %q", got)
	}
}

func TestProjectAnalysis_FileTypesSummary_Many(t *testing.T) {
	p := &compilationformatter.ProjectAnalysis{
		FileTypesDetected: []string{".md", ".py", ".go", ".ts", ".js"},
	}
	got := p.FileTypesSummary()
	if !strings.Contains(got, "more") {
		t.Fatalf("expected 'more' for many types, got %q", got)
	}
}

func TestCompilationFormatter_FormatDefault(t *testing.T) {
	f := compilationformatter.New(false)
	r := &compilationformatter.CompilationResults{
		OptimizationStats: compilationformatter.OptimizationStats{
			AverageContextEfficiency: 0.8,
		},
	}
	out := f.FormatDefault(r)
	if out == "" {
		t.Fatal("FormatDefault returned empty output")
	}
}

func TestCompilationFormatter_FormatDryRun(t *testing.T) {
	f := compilationformatter.New(false)
	r := &compilationformatter.CompilationResults{}
	out := f.FormatDryRun(r)
	if out == "" {
		t.Fatal("FormatDryRun returned empty output")
	}
}
