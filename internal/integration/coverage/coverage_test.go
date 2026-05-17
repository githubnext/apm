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

func TestCheckPrimitiveCoverage_multipleSpecials(t *testing.T) {
	prims := []string{"a", "b", "c"}
	dispatch := map[string]DispatchEntry{
		"a": {},
	}
	specials := map[string]bool{"b": true, "c": true}
	if err := CheckPrimitiveCoverage(prims, dispatch, specials); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_allSpecials(t *testing.T) {
	prims := []string{"x", "y"}
	specials := map[string]bool{"x": true, "y": true}
	if err := CheckPrimitiveCoverage(prims, nil, specials); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDispatchEntry_Fields(t *testing.T) {
	entry := DispatchEntry{
		Targets: []string{"target1", "target2"},
		Methods: []string{"install", "uninstall"},
	}
	if len(entry.Targets) != 2 {
		t.Errorf("Targets length: got %d, want 2", len(entry.Targets))
	}
	if entry.Methods[0] != "install" {
		t.Errorf("Methods[0]: got %q, want %q", entry.Methods[0], "install")
	}
}

func TestCheckPrimitiveCoverage_singlePrimSingleDispatch(t *testing.T) {
	prims := []string{"instructions"}
	dispatch := map[string]DispatchEntry{
		"instructions": {Targets: []string{"cursor"}, Methods: []string{"install"}},
	}
	if err := CheckPrimitiveCoverage(prims, dispatch, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPrimitiveCoverage_extraDispatchNotInSpecials(t *testing.T) {
	prims := []string{"a"}
	dispatch := map[string]DispatchEntry{
		"a": {},
		"b": {},
	}
	err := CheckPrimitiveCoverage(prims, dispatch, nil)
	if err == nil {
		t.Fatal("expected error for extra dispatch entry")
	}
}
