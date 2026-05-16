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
