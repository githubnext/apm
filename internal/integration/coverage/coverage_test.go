package coverage

import "testing"

func TestCheckPrimitiveCoverage_ok(t *testing.T) {
	prims := []string{"instructions", "skills"}
	dispatch := map[string]DispatchEntry{
		"instructions": {Targets: []string{"t1"}, Methods: []string{"m1"}},
		"skills":       {Targets: []string{"t2"}, Methods: []string{"m2"}},
	}
	if err := CheckPrimitiveCoverage(prims, dispatch, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_missingHandler(t *testing.T) {
	prims := []string{"instructions", "skills"}
	dispatch := map[string]DispatchEntry{
		"instructions": {},
	}
	err := CheckPrimitiveCoverage(prims, dispatch, nil)
	if err == nil {
		t.Fatal("expected error for unhandled primitive")
	}
}

func TestCheckPrimitiveCoverage_extraDispatch(t *testing.T) {
	prims := []string{"instructions"}
	dispatch := map[string]DispatchEntry{
		"instructions": {},
		"unknown":      {},
	}
	err := CheckPrimitiveCoverage(prims, dispatch, nil)
	if err == nil {
		t.Fatal("expected error for extra dispatch entry")
	}
}

func TestCheckPrimitiveCoverage_specialCase(t *testing.T) {
	prims := []string{"instructions", "special"}
	dispatch := map[string]DispatchEntry{
		"instructions": {},
	}
	specials := map[string]bool{"special": true}
	if err := CheckPrimitiveCoverage(prims, dispatch, specials); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_empty(t *testing.T) {
	if err := CheckPrimitiveCoverage(nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
