package apmyml_test

import (
"testing"

"github.com/githubnext/apm/internal/core/apmyml"
)

func TestParseTargetsField_plural(t *testing.T) {
data := map[string]interface{}{"targets": []interface{}{"claude", "copilot"}}
got, err := apmyml.ParseTargetsField(data)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 2 {
t.Errorf("expected 2 targets, got %v", got)
}
}

func TestParseTargetsField_singular(t *testing.T) {
data := map[string]interface{}{"target": "claude"}
got, err := apmyml.ParseTargetsField(data)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 1 || got[0] != "claude" {
t.Errorf("expected [claude], got %v", got)
}
}

func TestParseTargetsField_csv(t *testing.T) {
data := map[string]interface{}{"target": "claude,copilot"}
got, err := apmyml.ParseTargetsField(data)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 2 {
t.Errorf("expected 2 targets, got %v", got)
}
}

func TestParseTargetsField_both_conflict(t *testing.T) {
data := map[string]interface{}{"targets": []interface{}{"claude"}, "target": "copilot"}
_, err := apmyml.ParseTargetsField(data)
if err == nil {
t.Fatal("expected conflict error")
}
if _, ok := err.(*apmyml.ConflictingTargetsError); !ok {
t.Errorf("expected ConflictingTargetsError, got %T", err)
}
}

func TestParseTargetsField_empty(t *testing.T) {
got, err := apmyml.ParseTargetsField(map[string]interface{}{})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 0 {
t.Errorf("expected empty, got %v", got)
}
}

func TestParseTargetsField_unknown_target(t *testing.T) {
data := map[string]interface{}{"target": "unknown-tool"}
_, err := apmyml.ParseTargetsField(data)
if err == nil {
t.Fatal("expected error for unknown target")
}
}

func TestParseTargetsField_list_under_singular(t *testing.T) {
data := map[string]interface{}{"target": []interface{}{"claude", "copilot"}}
got, err := apmyml.ParseTargetsField(data)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 2 {
t.Errorf("expected 2 targets, got %v", got)
}
}

func TestParseTargetsField_whitespace_csv(t *testing.T) {
data := map[string]interface{}{"target": "claude , copilot"}
got, err := apmyml.ParseTargetsField(data)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 2 {
t.Errorf("expected 2, got %v", got)
}
}

func TestParseTargetsField_all_canonical_targets(t *testing.T) {
all := []interface{}{"claude", "copilot", "cursor", "opencode", "codex", "gemini", "windsurf", "agent-skills"}
data := map[string]interface{}{"targets": all}
got, err := apmyml.ParseTargetsField(data)
if err != nil {
t.Fatalf("unexpected error for all canonical: %v", err)
}
if len(got) != len(all) {
t.Errorf("expected %d targets, got %d", len(all), len(got))
}
}

func TestConflictingTargetsError_message(t *testing.T) {
data := map[string]interface{}{"targets": []interface{}{"claude"}, "target": "cursor"}
_, err := apmyml.ParseTargetsField(data)
if err == nil {
t.Fatal("expected error")
}
if err.Error() == "" {
t.Error("expected non-empty error message")
}
}

func TestUnknownTargetError_message(t *testing.T) {
data := map[string]interface{}{"target": "vscode"}
_, err := apmyml.ParseTargetsField(data)
if err == nil {
t.Fatal("expected error for unknown target")
}
if _, ok := err.(*apmyml.UnknownTargetError); !ok {
t.Errorf("expected UnknownTargetError, got %T", err)
}
if err.Error() == "" {
t.Error("expected non-empty error message")
}
}

func TestParseTargetsField_targets_empty_list(t *testing.T) {
data := map[string]interface{}{"targets": []interface{}{}}
_, err := apmyml.ParseTargetsField(data)
if err == nil {
t.Fatal("expected error for empty targets list")
}
if _, ok := err.(*apmyml.EmptyTargetsListError); !ok {
t.Errorf("expected EmptyTargetsListError, got %T", err)
}
}

func TestCanonicalTargets_present(t *testing.T) {
for name := range apmyml.CanonicalTargets {
if name == "" {
t.Error("canonical target should not be empty string")
}
}
if !apmyml.CanonicalTargets["claude"] {
t.Error("claude should be in canonical targets")
}
}
