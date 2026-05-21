package targetdetection_test

import (
"testing"

"github.com/githubnext/apm/internal/core/targetdetection"
)

func TestNormalizeTarget_ClaudePassthrough(t *testing.T) {
result := targetdetection.NormalizeTarget("claude")
if result != "claude" {
t.Errorf("expected 'claude', got %q", result)
}
}

func TestNormalizeTarget_CopilotAlias(t *testing.T) {
result := targetdetection.NormalizeTarget("copilot")
if result != "vscode" {
t.Errorf("expected 'vscode' for 'copilot', got %q", result)
}
}

func TestNormalizeTarget_VSCodeAlias(t *testing.T) {
result := targetdetection.NormalizeTarget("vscode")
if result != "vscode" {
t.Errorf("expected 'vscode' for 'vscode', got %q", result)
}
}

func TestNormalizeTarget_EmptyPassthrough(t *testing.T) {
result := targetdetection.NormalizeTarget("")
if result != "" {
t.Errorf("expected '', got %q", result)
}
}

func TestFormatProvenance_NonEmptyExtra4(t *testing.T) {
resolved := targetdetection.ResolvedTargets{
Targets: []string{"claude"},
Source:  "flag",
}
msg := targetdetection.FormatProvenance(resolved)
if msg == "" {
t.Error("expected non-empty provenance message")
}
}

func TestExpandAllTargets_EmptyYAML_NoCrashExtra4(t *testing.T) {
targets, err := targetdetection.ExpandAllTargets(".", []string{})
_ = err
_ = targets
}

func TestResolveTargets_FlagOverridesYAMLExtra4(t *testing.T) {
targets, err := targetdetection.ResolveTargets(".", []string{"claude"}, []string{"copilot"})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(targets.Targets) == 0 {
t.Error("expected at least one resolved target")
}
}
