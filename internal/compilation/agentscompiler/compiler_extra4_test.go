package agentscompiler

import (
"strings"
"testing"
)

func TestTargetConstants_AllNonEmpty_Extra4(t *testing.T) {
targets := []string{TargetAll, TargetVSCode, TargetAgents}
for _, target := range targets {
if target == "" {
t.Error("target constant should not be empty")
}
}
}

func TestTargetAll_Value_Extra4(t *testing.T) {
if TargetAll != "all" {
t.Errorf("expected 'all', got %q", TargetAll)
}
}

func TestStrategyDistributed_NonEmpty_Extra4(t *testing.T) {
if StrategyDistributed == "" {
t.Error("StrategyDistributed should not be empty")
}
}

func TestStrategySingleFile_NonEmpty_Extra4(t *testing.T) {
if StrategySingleFile == "" {
t.Error("StrategySingleFile should not be empty")
}
}

func TestStrategyConstants_AreDistinct_Extra4(t *testing.T) {
if StrategyDistributed == StrategySingleFile {
t.Error("strategy constants should differ")
}
}

func TestCompilationConfig_DefaultConfig_Extra4(t *testing.T) {
cfg := DefaultConfig()
if cfg.Target == "" {
t.Error("DefaultConfig should have non-empty Target")
}
}

func TestCompilationConfig_DefaultStrategy_Extra4(t *testing.T) {
cfg := DefaultConfig()
if cfg.Strategy == "" {
t.Error("DefaultConfig should have non-empty Strategy")
}
}

func TestBuildIDPlaceholder_ContainsAPM_Extra4(t *testing.T) {
if !strings.Contains(BuildIDPlaceholder, "APM") {
t.Errorf("expected BuildIDPlaceholder to contain 'APM', got %q", BuildIDPlaceholder)
}
}

func TestCopilotRootGeneratedMarker_NonEmpty_Extra4(t *testing.T) {
if CopilotRootGeneratedMarker == "" {
t.Error("CopilotRootGeneratedMarker should not be empty")
}
}

func TestCompilationResult_OKField_Extra4(t *testing.T) {
r := CompilationResult{Error: nil}
if !r.OK() {
t.Error("expected OK() to be true when Error is nil")
}
}
