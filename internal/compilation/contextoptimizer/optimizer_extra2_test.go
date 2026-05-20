package contextoptimizer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/compilation/contextoptimizer"
)

func TestDirectoryAnalysis_RelevanceScore_OneOfTwo(t *testing.T) {
	d := &contextoptimizer.DirectoryAnalysis{
		TotalFiles:    2,
		PatternCounts: map[string]int{"*.go": 1},
	}
	score := d.RelevanceScore("*.go")
	if score != 0.5 {
		t.Errorf("expected 0.5, got %f", score)
	}
}

func TestDirectoryAnalysis_RelevanceScore_ZeroTotal(t *testing.T) {
	d := &contextoptimizer.DirectoryAnalysis{TotalFiles: 0}
	score := d.RelevanceScore("*.go")
	if score != 0 {
		t.Errorf("expected 0 for empty dir, got %f", score)
	}
}

func TestInheritanceAnalysis_EfficiencyRatio_ZeroTotal(t *testing.T) {
	a := &contextoptimizer.InheritanceAnalysis{TotalContextLoad: 0, RelevantContextLoad: 0}
	if a.EfficiencyRatio() != 1.0 {
		t.Errorf("expected 1.0 for zero load, got %f", a.EfficiencyRatio())
	}
}

func TestInheritanceAnalysis_EfficiencyRatio_HalfRelevant(t *testing.T) {
	a := &contextoptimizer.InheritanceAnalysis{TotalContextLoad: 10, RelevantContextLoad: 5}
	if a.EfficiencyRatio() != 0.5 {
		t.Errorf("expected 0.5, got %f", a.EfficiencyRatio())
	}
}

func TestNew_WithBaseDir(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, nil)
	if opt == nil {
		t.Fatal("New should not return nil")
	}
}

func TestEnableTiming_DoesNotPanic2(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, nil)
	opt.EnableTiming(false)
	opt.EnableTiming(true)
}

func TestOptimizeInstructionPlacement_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, nil)
	result := opt.OptimizeInstructionPlacement([]string{})
	if result == nil {
		t.Fatal("result should not be nil")
	}
}

func TestOptimizeInstructionPlacement_SingleGoFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "main.go")
	if err := os.WriteFile(f, []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	opt := contextoptimizer.New(tmp, nil)
	result := opt.OptimizeInstructionPlacement([]string{"*.go"})
	if result == nil {
		t.Fatal("result should not be nil")
	}
}

func TestAnalyzeContextInheritance_ReturnsAnalysis(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, nil)
	analysis := opt.AnalyzeContextInheritance(tmp)
	if analysis == nil {
		t.Fatal("AnalyzeContextInheritance should not return nil")
	}
}

func TestAnalyzeContextInheritance_WorkingDirSet(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, nil)
	analysis := opt.AnalyzeContextInheritance(tmp)
	if analysis.WorkingDirectory == "" {
		t.Error("WorkingDirectory should be set")
	}
}

func TestGetOptimizationStats_WithNoResults(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, nil)
	result := opt.OptimizeInstructionPlacement([]string{})
	stats := opt.GetOptimizationStats(result)
	_ = stats // OptimizationStats is a value type
}

func TestPlacementCandidate_FieldAccess(t *testing.T) {
	pc := contextoptimizer.PlacementCandidate{
		Directory:     "/some/dir",
		Score:         0.8,
		CoverageRatio: 0.9,
		Depth:         2,
		IsLeaf:        true,
	}
	if pc.Directory != "/some/dir" {
		t.Error("Directory field mismatch")
	}
	if pc.Score != 0.8 {
		t.Error("Score field mismatch")
	}
	if !pc.IsLeaf {
		t.Error("IsLeaf should be true")
	}
}

func TestDefaultExcludedDirnames_HasNodeModules(t *testing.T) {
	if !contextoptimizer.DefaultExcludedDirnames["node_modules"] {
		t.Error("DefaultExcludedDirnames should include 'node_modules'")
	}
}

func TestDefaultExcludedDirnames_HasGit(t *testing.T) {
	if !contextoptimizer.DefaultExcludedDirnames[".git"] {
		t.Error("DefaultExcludedDirnames should include '.git'")
	}
}

func TestNew_WithExcludePatterns_DoesNotPanic(t *testing.T) {
	tmp := t.TempDir()
	opt := contextoptimizer.New(tmp, []string{"*.log", "tmp/*"})
	if opt == nil {
		t.Fatal("New should not return nil")
	}
}
