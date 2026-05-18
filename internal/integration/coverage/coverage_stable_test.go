package coverage

import (
"strings"
"testing"
)

func TestCheckPrimitiveCoverage_singleMatch(t *testing.T) {
prims := []string{"instructions"}
dispatch := map[string]DispatchEntry{
"instructions": {Targets: []string{"copilot"}, Methods: []string{"integrate"}},
}
err := CheckPrimitiveCoverage(prims, dispatch, nil)
if err != nil {
t.Errorf("expected no error for exact match: %v", err)
}
}

func TestCheckPrimitiveCoverage_multipleExact(t *testing.T) {
prims := []string{"instructions", "prompts", "hooks"}
dispatch := map[string]DispatchEntry{
"instructions": {Targets: []string{"copilot"}, Methods: []string{"integrate"}},
"prompts":      {Targets: []string{"claude"}, Methods: []string{"copy"}},
"hooks":        {Targets: []string{"vscode"}, Methods: []string{"write"}},
}
err := CheckPrimitiveCoverage(prims, dispatch, nil)
if err != nil {
t.Errorf("expected no error for exact multi-match: %v", err)
}
}

func TestCheckPrimitiveCoverage_nilSpecial_nil(t *testing.T) {
prims := []string{"instructions"}
dispatch := map[string]DispatchEntry{
"instructions": {},
}
err := CheckPrimitiveCoverage(prims, dispatch, nil)
if err != nil {
t.Errorf("nil special should work: %v", err)
}
}

func TestCheckPrimitiveCoverage_missingMultiple(t *testing.T) {
prims := []string{"instructions", "prompts", "skills"}
dispatch := map[string]DispatchEntry{
"instructions": {},
}
err := CheckPrimitiveCoverage(prims, dispatch, nil)
if err == nil {
t.Fatal("expected error for 2 unhandled primitives")
}
}

func TestCheckPrimitiveCoverage_errorMentionsMissing(t *testing.T) {
prims := []string{"instructions", "missing-primitive"}
dispatch := map[string]DispatchEntry{
"instructions": {},
}
err := CheckPrimitiveCoverage(prims, dispatch, nil)
if err == nil {
t.Fatal("expected error")
}
if !strings.Contains(err.Error(), "missing-primitive") {
t.Errorf("error should name the missing primitive: %v", err)
}
}

func TestDispatchEntry_emptyTargets(t *testing.T) {
d := DispatchEntry{Targets: []string{}, Methods: []string{"integrate"}}
if len(d.Targets) != 0 {
t.Errorf("expected empty targets, got %d", len(d.Targets))
}
}

func TestDispatchEntry_emptyMethods(t *testing.T) {
d := DispatchEntry{Targets: []string{"copilot"}, Methods: []string{}}
if len(d.Methods) != 0 {
t.Errorf("expected empty methods, got %d", len(d.Methods))
}
}

func TestCheckPrimitiveCoverage_specialAndDispatch_coexist(t *testing.T) {
prims := []string{"instructions", "hooks", "prompts"}
dispatch := map[string]DispatchEntry{
"instructions": {},
"prompts":      {},
}
special := map[string]bool{"hooks": true}
err := CheckPrimitiveCoverage(prims, dispatch, special)
if err != nil {
t.Errorf("expected no error when special covers hooks: %v", err)
}
}

func TestCheckPrimitiveCoverage_extraDispatch_notInSpecial(t *testing.T) {
prims := []string{"instructions"}
dispatch := map[string]DispatchEntry{
"instructions": {},
"extra-key":    {},
}
err := CheckPrimitiveCoverage(prims, dispatch, nil)
if err == nil {
t.Fatal("expected error for unrecognized dispatch entry")
}
}

func TestCheckPrimitiveCoverage_nilEverything(t *testing.T) {
err := CheckPrimitiveCoverage(nil, nil, nil)
if err != nil {
t.Errorf("nil everything should not error: %v", err)
}
}
