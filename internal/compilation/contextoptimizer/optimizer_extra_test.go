package contextoptimizer_test

import (
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/compilation/contextoptimizer"
)

func TestDirectoryAnalysis_RelevanceScore_AllMatch(t *testing.T) {
d := contextoptimizer.DirectoryAnalysis{
Directory:     "/root",
TotalFiles:    5,
PatternCounts: map[string]int{"*.go": 5},
}
if score := d.RelevanceScore("*.go"); score != 1.0 {
t.Fatalf("expected 1.0 got %f", score)
}
}

func TestDirectoryAnalysis_RelevanceScore_NoMatch(t *testing.T) {
d := contextoptimizer.DirectoryAnalysis{
Directory:     "/root",
TotalFiles:    10,
PatternCounts: map[string]int{"*.go": 0},
}
if score := d.RelevanceScore("*.py"); score != 0 {
t.Fatalf("expected 0 for missing pattern, got %f", score)
}
}

func TestInheritanceAnalysis_EfficiencyRatio_Full(t *testing.T) {
a := contextoptimizer.InheritanceAnalysis{
TotalContextLoad:    100,
RelevantContextLoad: 100,
}
if r := a.EfficiencyRatio(); r != 1.0 {
t.Fatalf("expected 1.0 got %f", r)
}
}

func TestInheritanceAnalysis_EfficiencyRatio_Zero(t *testing.T) {
a := contextoptimizer.InheritanceAnalysis{
TotalContextLoad:    100,
RelevantContextLoad: 0,
}
if r := a.EfficiencyRatio(); r != 0.0 {
t.Fatalf("expected 0.0 got %f", r)
}
}

func TestInheritanceAnalysis_PollutionScore_Field(t *testing.T) {
a := contextoptimizer.InheritanceAnalysis{PollutionScore: 0.42}
if a.PollutionScore != 0.42 {
t.Fatalf("unexpected PollutionScore %f", a.PollutionScore)
}
}

func TestNew_WithExcludePatterns(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, []string{"*.log"})
if opt == nil {
t.Fatal("expected non-nil optimizer with exclude patterns")
}
}

func TestNew_NilExcludePatterns(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
if opt == nil {
t.Fatal("expected non-nil optimizer with nil exclude patterns")
}
}

func TestOptimizeInstructionPlacement_SingleFile(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644); err != nil {
t.Fatal(err)
}
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{"*.go"})
if result == nil {
t.Fatal("expected non-nil result")
}
if len(result.Decisions) == 0 {
t.Fatal("expected at least one decision")
}
}

func TestOptimizeInstructionPlacement_MultiplePatterns(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, "a.go"), []byte("package x"), 0644)
os.WriteFile(filepath.Join(dir, "b.py"), []byte("# py"), 0644)
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{"*.go", "*.py"})
if result == nil {
t.Fatal("expected non-nil result")
}
stats := opt.GetOptimizationStats(result)
if stats.TotalInstructions != 2 {
t.Fatalf("expected 2 instructions, got %d", stats.TotalInstructions)
}
}

func TestOptimizeInstructionPlacement_EmptyPatterns(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{})
if result == nil {
t.Fatal("expected non-nil result even with no patterns")
}
}

func TestGetOptimizationStats_EmptyResult(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement(nil)
stats := opt.GetOptimizationStats(result)
if stats.TotalInstructions != 0 {
t.Fatalf("expected 0, got %d", stats.TotalInstructions)
}
}

func TestGetOptimizationStats_AllUnchanged(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{"*.go"})
stats := opt.GetOptimizationStats(result)
total := stats.Optimized + stats.Unchanged
if stats.TotalInstructions > 0 && total != stats.TotalInstructions {
t.Fatalf("Optimized+Unchanged (%d) != TotalInstructions (%d)", total, stats.TotalInstructions)
}
}

func TestAnalyzeContextInheritance_RootDir(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
analysis := opt.AnalyzeContextInheritance(dir)
if analysis == nil {
t.Fatal("expected non-nil inheritance analysis")
}
}

func TestAnalyzeContextInheritance_SubDir(t *testing.T) {
dir := t.TempDir()
sub := filepath.Join(dir, "src", "pkg")
os.MkdirAll(sub, 0755)
opt := contextoptimizer.New(dir, nil)
analysis := opt.AnalyzeContextInheritance(sub)
if analysis == nil {
t.Fatal("expected non-nil inheritance analysis")
}
if len(analysis.InheritanceChain) == 0 {
t.Fatal("expected non-empty inheritance chain for subdirectory")
}
}

func TestAnalyzeContextInheritance_WorkingDirectory(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
analysis := opt.AnalyzeContextInheritance(dir)
if analysis.WorkingDirectory != dir {
t.Fatalf("expected working dir %q, got %q", dir, analysis.WorkingDirectory)
}
}

func TestOptimizationResult_Decisions_Field(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, "x.go"), []byte("package x"), 0644)
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{"*.go"})
for _, dec := range result.Decisions {
if dec.Strategy == "" {
t.Error("expected non-empty Strategy in PlacementDecision")
}
}
}

func TestPlacementDecision_StrategyValues(t *testing.T) {
valid := map[string]bool{
"single_point": true,
"distributed":  true,
"selective":    true,
"unchanged":    true,
}
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, "f.go"), []byte("package f"), 0644)
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{"*.go"})
for _, dec := range result.Decisions {
if !valid[dec.Strategy] {
t.Errorf("unexpected Strategy %q", dec.Strategy)
}
}
}

func TestOptimizeInstructionPlacement_NestedFiles(t *testing.T) {
dir := t.TempDir()
sub := filepath.Join(dir, "sub")
os.MkdirAll(sub, 0755)
os.WriteFile(filepath.Join(dir, "root.go"), []byte("package x"), 0644)
os.WriteFile(filepath.Join(sub, "leaf.go"), []byte("package x"), 0644)
opt := contextoptimizer.New(dir, nil)
result := opt.OptimizeInstructionPlacement([]string{"*.go"})
if result == nil {
t.Fatal("expected non-nil result")
}
}

func TestEnableTiming_DoesNotPanic(t *testing.T) {
dir := t.TempDir()
opt := contextoptimizer.New(dir, nil)
defer func() {
if r := recover(); r != nil {
t.Fatalf("EnableTiming panicked: %v", r)
}
}()
opt.EnableTiming(true)
opt.EnableTiming(false)
}

func TestOptimizeInstructionPlacement_ExcludePattern(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)
os.WriteFile(filepath.Join(dir, "skip.log"), []byte("log data"), 0644)
opt := contextoptimizer.New(dir, []string{"*.log"})
result := opt.OptimizeInstructionPlacement([]string{"*.go"})
if result == nil {
t.Fatal("expected non-nil result")
}
}
