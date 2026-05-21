package contextoptimizer

import "testing"

func TestDirectoryAnalysis_TotalFilesField_Extra4(t *testing.T) {
d := DirectoryAnalysis{TotalFiles: 5}
if d.TotalFiles != 5 {
t.Errorf("expected 5, got %d", d.TotalFiles)
}
}

func TestDirectoryAnalysis_DepthField_Extra4(t *testing.T) {
d := DirectoryAnalysis{Depth: 3}
if d.Depth != 3 {
t.Errorf("expected 3, got %d", d.Depth)
}
}

func TestDirectoryAnalysis_DirectoryField_Extra4(t *testing.T) {
d := DirectoryAnalysis{Directory: "/a/b"}
if d.Directory != "/a/b" {
t.Errorf("expected '/a/b', got %q", d.Directory)
}
}

func TestInheritanceAnalysis_WorkingDirectoryField_Extra4(t *testing.T) {
a := InheritanceAnalysis{WorkingDirectory: "/x/y"}
if a.WorkingDirectory != "/x/y" {
t.Errorf("expected '/x/y', got %q", a.WorkingDirectory)
}
}

func TestInheritanceAnalysis_TotalContextLoad_Extra4(t *testing.T) {
a := InheritanceAnalysis{TotalContextLoad: 10, RelevantContextLoad: 4}
if a.TotalContextLoad != 10 {
t.Errorf("expected 10, got %d", a.TotalContextLoad)
}
}

func TestInheritanceAnalysis_EfficiencyRatio_HighLoad_Extra4(t *testing.T) {
a := InheritanceAnalysis{RelevantContextLoad: 9, TotalContextLoad: 10}
ratio := a.EfficiencyRatio()
if ratio < 0 || ratio > 1 {
t.Errorf("ratio out of range: %f", ratio)
}
}

func TestPlacementCandidate_DirectoryField_Extra4(t *testing.T) {
pc := PlacementCandidate{Directory: "/a/b"}
if pc.Directory != "/a/b" {
t.Errorf("expected '/a/b', got %q", pc.Directory)
}
}

func TestPlacementDecision_TargetDirectoryField_Extra4(t *testing.T) {
pd := PlacementDecision{TargetDirectory: "/c/d"}
if pd.TargetDirectory != "/c/d" {
t.Errorf("expected '/c/d', got %q", pd.TargetDirectory)
}
}

func TestOptimizationResult_ZeroValue_Extra4(t *testing.T) {
var r OptimizationResult
if r.Decisions != nil {
t.Error("zero Decisions should be nil")
}
}

func TestDefaultExcludedDirnames_IsMap_Extra4(t *testing.T) {
if DefaultExcludedDirnames == nil {
t.Error("DefaultExcludedDirnames should not be nil")
}
}
