package contextoptimizer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_ReturnNotNil(t *testing.T) {
	o := New(".", nil)
	if o == nil {
		t.Fatal("expected non-nil optimizer")
	}
}

func TestOptimizeInstructionPlacement_NilPatterns(t *testing.T) {
	o := New(".", nil)
	result := o.OptimizeInstructionPlacement(nil)
	if result == nil {
		t.Fatal("expected non-nil result for nil patterns")
	}
}

func TestOptimizeInstructionPlacement_EmptyPatterns2(t *testing.T) {
	o := New(".", nil)
	result := o.OptimizeInstructionPlacement([]string{})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestDefaultExcludedDirnames_ContainsDotGit(t *testing.T) {
	if !DefaultExcludedDirnames[".git"] {
		t.Fatal("expected .git in DefaultExcludedDirnames")
	}
}

func TestDefaultExcludedDirnames_ContainsNodeModules(t *testing.T) {
	if !DefaultExcludedDirnames["node_modules"] {
		t.Fatal("expected node_modules in DefaultExcludedDirnames")
	}
}

func TestDefaultExcludedDirnames_ContainsBuild(t *testing.T) {
	if !DefaultExcludedDirnames["build"] {
		t.Fatal("expected build in DefaultExcludedDirnames")
	}
}

func TestDirectoryAnalysis_RelevanceScore_NonNegative(t *testing.T) {
	d := &DirectoryAnalysis{
		PatternCounts: map[string]int{"*.go": 1},
		TotalFiles:    1,
	}
	score := d.RelevanceScore("*.go")
	if score < 0 {
		t.Fatalf("expected non-negative score, got %f", score)
	}
}

func TestInheritanceAnalysis_EfficiencyRatio_NonNegative(t *testing.T) {
	a := &InheritanceAnalysis{
		TotalContextLoad:    5,
		RelevantContextLoad: 3,
	}
	ratio := a.EfficiencyRatio()
	if ratio < 0 {
		t.Fatalf("expected non-negative ratio, got %f", ratio)
	}
}

func TestContextOptimizer_WithRealDir(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "test.go"), []byte("package main"), 0o644)
	o := New(dir, nil)
	result := o.OptimizeInstructionPlacement([]string{"*.go"})
	if result == nil {
		t.Fatal("expected non-nil result with real directory")
	}
}

func TestAnalyzeContextInheritance_WorkingDirField(t *testing.T) {
	o := New(".", nil)
	a := o.AnalyzeContextInheritance(".")
	if a == nil {
		t.Fatal("expected non-nil analysis")
	}
}

func TestOptimizationResult_DecisionsField(t *testing.T) {
	r := &OptimizationResult{
		Decisions: []PlacementDecision{},
	}
	if r.Decisions == nil {
		t.Fatal("expected non-nil decisions slice")
	}
}

func TestPlacementDecision_StrategyField(t *testing.T) {
	d := PlacementDecision{Strategy: "global"}
	if d.Strategy != "global" {
		t.Fatalf("expected 'global', got %q", d.Strategy)
	}
}
