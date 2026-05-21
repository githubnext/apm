package coverage

import (
	"testing"
)

func TestCheckPrimitiveCoverage_ExactMatch(t *testing.T) {
	prims := []string{"instructions"}
	dispatch := map[string]DispatchEntry{
		"instructions": {Targets: []string{"copilot"}, Methods: []string{"integrate"}},
	}
	if err := CheckPrimitiveCoverage(prims, dispatch, nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_SpecialCasesOnly(t *testing.T) {
	prims := []string{"system"}
	dispatch := map[string]DispatchEntry{}
	special := map[string]bool{"system": true}
	if err := CheckPrimitiveCoverage(prims, dispatch, special); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_ExtraInDispatch(t *testing.T) {
	prims := []string{"instructions"}
	dispatch := map[string]DispatchEntry{
		"instructions": {},
		"unknown":      {},
	}
	err := CheckPrimitiveCoverage(prims, dispatch, nil)
	if err == nil {
		t.Error("expected error for dispatch entry not in primitives")
	}
}

func TestDispatchEntry_ZeroValue(t *testing.T) {
	var d DispatchEntry
	if len(d.Targets) != 0 || len(d.Methods) != 0 {
		t.Error("zero value should have empty slices")
	}
}

func TestDispatchEntry_WithValues(t *testing.T) {
	d := DispatchEntry{
		Targets: []string{"copilot", "claude"},
		Methods: []string{"integrate", "remove"},
	}
	if len(d.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(d.Targets))
	}
	if len(d.Methods) != 2 {
		t.Errorf("expected 2 methods, got %d", len(d.Methods))
	}
}

func TestCheckPrimitiveCoverage_BothEmpty(t *testing.T) {
	if err := CheckPrimitiveCoverage(nil, nil, nil); err != nil {
		t.Errorf("empty inputs should not error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_DispatchAndSpecialBothCoverPrimitive(t *testing.T) {
	prims := []string{"p1", "p2"}
	dispatch := map[string]DispatchEntry{"p1": {}}
	special := map[string]bool{"p2": true}
	if err := CheckPrimitiveCoverage(prims, dispatch, special); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
