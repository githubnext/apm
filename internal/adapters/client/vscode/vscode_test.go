package vscode

import (
	"testing"
)

func TestTranslateEnvValueForVSCode_legacy_angle_var(t *testing.T) {
	got := translateEnvValueForVSCode("<MY_TOKEN>")
	if got != "${env:MY_TOKEN}" {
		t.Errorf("expected ${env:MY_TOKEN}, got %s", got)
	}
}

func TestTranslateEnvValueForVSCode_dollar_brace(t *testing.T) {
	got := translateEnvValueForVSCode("${MY_TOKEN}")
	if got != "${env:MY_TOKEN}" {
		t.Errorf("expected ${env:MY_TOKEN}, got %s", got)
	}
}

func TestTranslateEnvValueForVSCode_already_env(t *testing.T) {
	got := translateEnvValueForVSCode("${env:MY_TOKEN}")
	if got != "${env:MY_TOKEN}" {
		t.Errorf("already env: prefix should be preserved, got %s", got)
	}
}

func TestTranslateEnvValueForVSCode_plain_string(t *testing.T) {
	got := translateEnvValueForVSCode("no-vars-here")
	if got != "no-vars-here" {
		t.Errorf("plain string should be unchanged, got %s", got)
	}
}

func TestFilterOut_removes_target(t *testing.T) {
	ss := []string{"a", "b", "c", "b"}
	got := filterOut(ss, "b")
	if len(got) != 2 {
		t.Errorf("expected 2 items, got %d: %v", len(got), got)
	}
	for _, s := range got {
		if s == "b" {
			t.Error("filterOut should remove all occurrences of target")
		}
	}
}

func TestFilterOut_no_match(t *testing.T) {
	ss := []string{"a", "c"}
	got := filterOut(ss, "b")
	if len(got) != 2 {
		t.Errorf("no match should return same length, got %d", len(got))
	}
}

func TestFilterOut_empty(t *testing.T) {
	got := filterOut(nil, "x")
	if len(got) != 0 {
		t.Errorf("empty input should return empty, got %v", got)
	}
}

func TestStrField_present(t *testing.T) {
	m := map[string]interface{}{"key": "value"}
	if strField(m, "key") != "value" {
		t.Error("expected 'value'")
	}
}

func TestStrField_absent(t *testing.T) {
	m := map[string]interface{}{}
	if strField(m, "missing") != "" {
		t.Error("expected empty string for missing key")
	}
}

func TestToStringSlice_string_slice(t *testing.T) {
	v := []string{"a", "b"}
	got := toStringSlice(v)
	if len(got) != 2 || got[0] != "a" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestToStringSlice_interface_slice(t *testing.T) {
	v := []interface{}{"x", "y"}
	got := toStringSlice(v)
	if len(got) != 2 || got[0] != "x" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestToStringSlice_nil(t *testing.T) {
	got := toStringSlice(nil)
	if len(got) != 0 {
		t.Errorf("nil should return empty, got %v", got)
	}
}

func TestExtractPackageArgs_combined(t *testing.T) {
	pkg := map[string]interface{}{
		"runtime_arguments": []string{"--arg1"},
		"package_arguments": []string{"--pkg"},
	}
	got := extractPackageArgs(pkg)
	if len(got) != 2 {
		t.Errorf("expected 2 args, got %v", got)
	}
}

func TestExtractPackageArgs_empty(t *testing.T) {
	pkg := map[string]interface{}{}
	got := extractPackageArgs(pkg)
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestToInterfaceSlice(t *testing.T) {
	ss := []string{"a", "b", "c"}
	got := toInterfaceSlice(ss)
	if len(got) != 3 {
		t.Errorf("expected 3, got %d", len(got))
	}
}

func TestToSliceOfMaps(t *testing.T) {
	v := []interface{}{
		map[string]interface{}{"k": "v"},
		map[string]interface{}{"k2": "v2"},
	}
	got := toSliceOfMaps(v)
	if len(got) != 2 {
		t.Errorf("expected 2 maps, got %d", len(got))
	}
}

func TestToSliceOfMaps_non_slice(t *testing.T) {
	got := toSliceOfMaps("not-a-slice")
	if got != nil {
		t.Errorf("expected nil for non-slice input, got %v", got)
	}
}
