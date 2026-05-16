package contextoptimizer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/compilation/contextoptimizer"
)

func TestDirectoryAnalysis_RelevanceScore_Empty(t *testing.T) {
	d := contextoptimizer.DirectoryAnalysis{
		Directory:     "/some/dir",
		TotalFiles:    0,
		PatternCounts: map[string]int{},
	}
	score := d.RelevanceScore("*.py")
	if score != 0 {
		t.Fatalf("expected 0 for empty directory, got %f", score)
	}
}

func TestDirectoryAnalysis_RelevanceScore_Full(t *testing.T) {
	d := contextoptimizer.DirectoryAnalysis{
		Directory:     "/some/dir",
		TotalFiles:    10,
		PatternCounts: map[string]int{"*.py": 5},
	}
	score := d.RelevanceScore("*.py")
	if score != 0.5 {
		t.Fatalf("expected 0.5, got %f", score)
	}
}

func TestDirectoryAnalysis_RelevanceScore_MissingPattern(t *testing.T) {
	d := contextoptimizer.DirectoryAnalysis{
		Directory:     "/some/dir",
		TotalFiles:    10,
		PatternCounts: map[string]int{},
	}
	score := d.RelevanceScore("*.go")
	if score != 0 {
		t.Fatalf("expected 0 for missing pattern, got %f", score)
	}
}

func TestInheritanceAnalysis_EfficiencyRatio_NoLoad(t *testing.T) {
	a := contextoptimizer.InheritanceAnalysis{
		TotalContextLoad: 0,
	}
	if a.EfficiencyRatio() != 1 {
		t.Fatal("expected 1 when TotalContextLoad is 0")
	}
}

func TestInheritanceAnalysis_EfficiencyRatio_Partial(t *testing.T) {
	a := contextoptimizer.InheritanceAnalysis{
		TotalContextLoad:    100,
		RelevantContextLoad: 40,
	}
	got := a.EfficiencyRatio()
	if got != 0.4 {
		t.Fatalf("expected 0.4, got %f", got)
	}
}

func TestNew_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	opt := contextoptimizer.New(dir, nil)
	if opt == nil {
		t.Fatal("expected non-nil optimizer")
	}
}

func TestOptimizeInstructionPlacement_NoPatterns(t *testing.T) {
	dir := t.TempDir()
	opt := contextoptimizer.New(dir, nil)
	result := opt.OptimizeInstructionPlacement(nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestOptimizeInstructionPlacement_WithFiles(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "src")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "main.py"), []byte("# py"), 0644); err != nil {
		t.Fatal(err)
	}
	opt := contextoptimizer.New(dir, nil)
	result := opt.OptimizeInstructionPlacement([]string{"*.py"})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	stats := opt.GetOptimizationStats(result)
	if stats.TotalInstructions != 1 {
		t.Fatalf("expected 1 instruction, got %d", stats.TotalInstructions)
	}
}
