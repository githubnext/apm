package coverage

import (
	"strings"
	"testing"
)

func TestCheckPrimitiveCoverage_missingHandlerMessage(t *testing.T) {
	prims := []string{"instructions", "prompts"}
	dispatch := map[string]DispatchEntry{
		"instructions": {Targets: []string{"copilot"}, Methods: []string{"integrate"}},
	}
	err := CheckPrimitiveCoverage(prims, dispatch, nil)
	if err == nil {
		t.Fatal("expected error for unhandled primitive")
	}
	if !strings.Contains(err.Error(), "prompts") {
		t.Errorf("error should mention missing primitive, got: %v", err)
	}
}

func TestCheckPrimitiveCoverage_extraDispatchMessage(t *testing.T) {
	prims := []string{"instructions"}
	dispatch := map[string]DispatchEntry{
		"instructions": {Targets: []string{"copilot"}, Methods: []string{"integrate"}},
		"unknown":      {Targets: []string{"x"}, Methods: []string{"y"}},
	}
	err := CheckPrimitiveCoverage(prims, dispatch, nil)
	if err == nil {
		t.Fatal("expected error for extra dispatch entry")
	}
	if !strings.Contains(err.Error(), "unknown") {
		t.Errorf("error should mention extra entry, got: %v", err)
	}
}

func TestCheckPrimitiveCoverage_specialCaseCoversDispatch(t *testing.T) {
	prims := []string{"instructions", "hooks"}
	dispatch := map[string]DispatchEntry{
		"instructions": {Targets: []string{"copilot"}, Methods: []string{"integrate"}},
	}
	special := map[string]bool{"hooks": true}
	err := CheckPrimitiveCoverage(prims, dispatch, special)
	if err != nil {
		t.Errorf("special case should cover hooks: %v", err)
	}
}

func TestCheckPrimitiveCoverage_emptySlices(t *testing.T) {
	err := CheckPrimitiveCoverage([]string{}, map[string]DispatchEntry{}, map[string]bool{})
	if err != nil {
		t.Errorf("empty slices should not error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_allSpecial(t *testing.T) {
	prims := []string{"hooks", "prompts"}
	special := map[string]bool{"hooks": true, "prompts": true}
	err := CheckPrimitiveCoverage(prims, nil, special)
	if err != nil {
		t.Errorf("all-special should not error: %v", err)
	}
}

func TestDispatchEntry_fields(t *testing.T) {
	d := DispatchEntry{
		Targets: []string{"copilot", "claude"},
		Methods: []string{"integrate", "copy"},
	}
	if len(d.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(d.Targets))
	}
	if d.Methods[0] != "integrate" {
		t.Errorf("expected first method integrate, got %s", d.Methods[0])
	}
}

func TestCheckPrimitiveCoverage_specialOverridesDispatch(t *testing.T) {
	prims := []string{"skills"}
	dispatch := map[string]DispatchEntry{}
	special := map[string]bool{"skills": true}
	err := CheckPrimitiveCoverage(prims, dispatch, special)
	if err != nil {
		t.Errorf("special should cover missing dispatch: %v", err)
	}
}

func TestCheckPrimitiveCoverage_extraDispatchCoveredBySpecial(t *testing.T) {
	prims := []string{"instructions"}
	dispatch := map[string]DispatchEntry{
		"instructions": {},
		"extra":        {},
	}
	special := map[string]bool{"extra": true}
	err := CheckPrimitiveCoverage(prims, dispatch, special)
	if err != nil {
		t.Errorf("special covers extra dispatch entry: %v", err)
	}
}
